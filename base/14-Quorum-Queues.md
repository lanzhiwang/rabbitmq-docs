# Quorum Queues

https://www.rabbitmq.com/quorum-queues.html

## Overview

The quorum queue is a modern queue type for RabbitMQ implementing a durable, replicated FIFO queue based on the [Raft consensus algorithm](https://raft.github.io/). It is available as of RabbitMQ 3.8.0.

The quorum queue type is an alternative to durable [mirrored queues](https://www.rabbitmq.com/ha.html) purpose built for a [set of use cases](https://www.rabbitmq.com/quorum-queues.html#use-cases) where [data safety](https://www.rabbitmq.com/quorum-queues.html#data-safety) is a top priority. This is covered in [Motivation](https://www.rabbitmq.com/quorum-queues.html#motivation). They should be considered the default option for a replicated queue type.

Quorum queues also have important [differences in behaviour](https://www.rabbitmq.com/quorum-queues.html#behaviour) and some [limitations](https://www.rabbitmq.com/quorum-queues.html#feature-comparison) compared to classic mirrored queues, including workload-specific ones, e.g. when consumers [repeatedly requeue the same message](https://www.rabbitmq.com/quorum-queues.html#repeated-requeues).

Some features, such as [poison message handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling), are specific to quorum queues.

## [Overview](https://www.rabbitmq.com/quorum-queues.html#overview)

Topics covered in this guide include

- [What are quorum queues](https://www.rabbitmq.com/quorum-queues.html#motivation) and why they were introduced
- [How are they different](https://www.rabbitmq.com/quorum-queues.html#feature-comparison) from classic queues
- Primary [use cases](https://www.rabbitmq.com/quorum-queues.html#use-cases) of quorum queues and when not to use them
- How to [declare a quorum queue](https://www.rabbitmq.com/quorum-queues.html#usage)
- [Replication](https://www.rabbitmq.com/quorum-queues.html#replication)-related topics: replica management, [replica leader rebalancing](https://www.rabbitmq.com/quorum-queues.html#replica-rebalancing), optimal number of replicas, etc
- What guarantees quorum queues offer in terms of [leader failure handling](https://www.rabbitmq.com/quorum-queues.html#leader-election), [data safety](https://www.rabbitmq.com/quorum-queues.html#data-safety) and [availability](https://www.rabbitmq.com/quorum-queues.html#availability)
- [Performance](https://www.rabbitmq.com/quorum-queues.html#performance) characteristics
- [Poison message handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling) provided by quorum queues
- [Configurable settings](https://www.rabbitmq.com/quorum-queues.html#configuration) of quorum queues
- Resource use of quorum queues, most importantly their [memory footprint](https://www.rabbitmq.com/quorum-queues.html#resource-use)

and more.

This guide assumes general familiarity with [RabbitMQ clustering](https://www.rabbitmq.com/clustering.html).

## [Motivation](https://www.rabbitmq.com/quorum-queues.html#motivation)

Quorum queues are designed to be safer and provide simpler, well defined failure handling semantics that users should find easier to reason about when designing and operating their systems.

These design choices come with constraints. To reach this goal, quorum queues adopt a different replication and consensus protocol and give up support for certain "transient" in nature features. These constraints and limitations are covered later in this guide.

### [What is a Quorum?](https://www.rabbitmq.com/quorum-queues.html#what-is-quorum)

If intentionally simplified, [quorum](https://en.wikipedia.org/wiki/Quorum) in a distributed system can be defined as an agreement between the majority of nodes ((N/2)+1 where N is the total number of system participants).

When applied to queue mirroring in RabbitMQ [clusters](https://www.rabbitmq.com/clustering.html) this means that the majority of replicas (including the currently elected queue leader) agree on the state of the queue and its contents.

### Differences from Classic Mirrored Queues

Quorum queues share many of the fundamentals with [queues](https://www.rabbitmq.com/queues.html) of other types in RabbitMQ. However, they are more purpose-built, focus on data safety and predictable recovery, and do not support certain features.

The differences [are covered](https://www.rabbitmq.com/quorum-queues.html#feature-comparison) in this guide.

Classic mirrored queues in RabbitMQ have technical limitations that makes it difficult to provide comprehensible guarantees and clear failure handling semantics.

Certain failure scenarios can result in mirrored queues confirming messages too early, potentially resulting in a data loss.

## [Feature Comparison with Regular Queues](https://www.rabbitmq.com/quorum-queues.html#feature-comparison)

Quorum queues share most of the fundamentals with other [queue](https://www.rabbitmq.com/queues.html) types. A client library that can use regular mirrored queues will be able to use quorum queues.

The following operations work the same way for quorum queues as they do for regular queues:

- Consumption (subscription)
- [Consumer acknowledgements](https://www.rabbitmq.com/confirms.html) (except for global [QoS and prefetch](https://www.rabbitmq.com/quorum-queues.html#global-qos))
- Cancelling consumers
- Purging
- Deletion

With some queue operations there are minor differences:

- [Declaration](https://www.rabbitmq.com/quorum-queues.html#declaring)
- Setting prefetch for consumers

Some features are not currently supported by quorum queues.

### [Feature Matrix](https://www.rabbitmq.com/quorum-queues.html#feature-matrix)

| Feature                                                      | Classic Mirrored | Quorum                                                       |
| :----------------------------------------------------------- | :--------------- | :----------------------------------------------------------- |
| [Non-durable queues](https://www.rabbitmq.com/queues.html)   | yes              | no                                                           |
| [Exclusivity](https://www.rabbitmq.com/queues.html)          | yes              | no                                                           |
| Per message persistence                                      | per message      | always                                                       |
| Membership changes                                           | automatic        | manual                                                       |
| [Message TTL](https://www.rabbitmq.com/ttl.html)             | yes              | no                                                           |
| [Queue TTL](https://www.rabbitmq.com/ttl.html#queue-ttl)     | yes              | yes                                                          |
| [Queue length limits](https://www.rabbitmq.com/maxlength.html) | yes              | yes (except x-overflow: reject-publish-dlx)                  |
| [Lazy behaviour](https://www.rabbitmq.com/lazy-queues.html)  | yes              | yes through the [Memory Limit](https://www.rabbitmq.com/quorum-queues.html#memory-limit) feature |
| [Message priority](https://www.rabbitmq.com/priority.html)   | yes              | no                                                           |
| [Consumer priority](https://www.rabbitmq.com/consumer-priority.html) | yes              | yes                                                          |
| [Dead letter exchanges](https://www.rabbitmq.com/dlx.html)   | yes              | yes                                                          |
| Adheres to [policies](https://www.rabbitmq.com/parameters.html#policies) | yes              | yes (see policy support below)                               |
| Reacts to [memory alarms](https://www.rabbitmq.com/alarms.html) | yes              | no                                                           |
| Poison message handling                                      | no               | yes                                                          |
| Global [QoS Prefetch](https://www.rabbitmq.com/quorum-queues.html#global-qos) | yes              | no                                                           |

#### Non-durable Queues

Regular queues can be [non-durable](https://www.rabbitmq.com/queues.html). Quorum queues are always durable per their assumed [use cases](https://www.rabbitmq.com/quorum-queues.html#use-cases).

#### Exclusivity

[Exclusive queues](https://www.rabbitmq.com/queues.html#exclusive-queues) are tied to the lifecycle of their declaring connection. Quorum queues by design are replicated and durable, therefore the exclusive property makes no sense in their context. Therefore quorum queues cannot be exclusive.

Quorum queues are not meant to be used as [temporary queues](https://www.rabbitmq.com/queues.html#temporary-queues).

#### TTL

Quorum queues do not currently support Message TTL, but they do support [Queue TTL](https://www.rabbitmq.com/ttl.html#queue-ttl).

#### Length Limit

Quorum queues has support for [queue length limits](https://www.rabbitmq.com/maxlength.html).

The drop-head and reject-publish overflow behaviours are supported but they do not support reject-publish-dlx configurations as Quorum queues take a different implementation approach than classic queues.

When a quorum queue reaches the max-length limit and reject-publish is configured it notifies each publishing channel who from thereon will reject all messages back to the client. This means that quorum queues may overshoot their limit by some small number of messages as there may be messages in flight whilst the channels are notified. The number of additional messages that are accepted by the queue will vary depending on how many messages are in flight at the time.

#### Dead Lettering

Quorum queues do support [dead letter exchanges](https://www.rabbitmq.com/dlx.html) (DLXs).

#### Lazy Mode

Quorum queues store their content on disk (per Raft requirements) as well as in memory (up to the [in memory limit configured](https://www.rabbitmq.com/quorum-queues.html#memory-limit)).

The [lazy mode](https://www.rabbitmq.com/lazy-queues.html) does not apply to them.

It is possible to [limit how many messages a quorum queue keeps in memory](https://www.rabbitmq.com/quorum-queues.html#memory-limit) using a policy which can achieve a behaviour similar to lazy queues.

#### [Global QoS](https://www.rabbitmq.com/quorum-queues.html#global-qos)

Quorum queues do not support global [QoS prefetch](https://www.rabbitmq.com/confirms.html#channel-qos-prefetch) where a channel sets a single prefetch limit for all consumers using that channel. If an attempt is made to consume from a quorum queue from a channel with global QoS enabled a channel error will be returned.

Use [per-consumer QoS prefetch](https://www.rabbitmq.com/consumer-prefetch.html), which is the default in several popular clients.

#### Priorities

Quorum queues do not currently support [priorities](https://www.rabbitmq.com/priority.html), including [consumer priorities](https://www.rabbitmq.com/consumer-priority.html).

To achieve priority processing with Quorum Queues multiple queues should be used instead; one for each priority.

#### Poison Message Handling

Quorum queues [support poison message handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling) via a redelivery limit. This feature is currently unique to Quorum queues.

#### Policy Support

Quorum queues can be configured via RabbitMQ policies. The below table summarises the policy keys they adhere to.

| Definition Key          | Type                            |
| :---------------------- | :------------------------------ |
| max-length              | Number                          |
| max-length-bytes        | Number                          |
| overflow                | "drop-head" or "reject-publish" |
| expires                 | Number (milliseconds)           |
| dead-letter-exchange    | String                          |
| dead-letter-routing-key | String                          |
| max-in-memory-length    | Number                          |
| max-in-memory-bytes     | Number                          |
| delivery-limit          | Number                          |

## [Use Cases](https://www.rabbitmq.com/quorum-queues.html#use-cases)

Quorum queues are purpose built by design. They are *not* designed to be used for every problem. Their intended use is for topologies where queues exist for a long time and are critical to certain aspects of system operation, therefore fault tolerance and data safety is more important than, say, lowest possible latency and advanced queue features.

Examples would be incoming orders in a sales system or votes cast in an election system where potentially losing messages would have a significant impact on system correctness and function.

Stock tickers and instant messaging systems benefit less or not at all from quorum queues.

Publishers should use publisher confirms as this is how clients can interact with the quorum queue consensus system. Publisher confirms will only be issued once a published message has been successfully replicated to a quorum of nodes and is considered "safe" within the context of the system.

Consumers should use manual acknowledgements to ensure messages that aren't successfully processed are returned to the queue so that another consumer can re-attempt processing.

### When Not to Use Quorum Queues

In some cases quorum queues should not be used. They typically involve:

- Temporary nature of queues: transient or exclusive queues, high queue churn (declaration and deletion rates)
- Lowest possible latency: the underlying consensus algorithm has an inherently higher latency due to its data safety features
- When data safety is not a priority (e.g. applications do not use [manual acknowledgements and publisher confirms](https://www.rabbitmq.com/confirms.html) are not used)
- Very long queue backlogs (quorum queues currently keep all messages in memory at all times, up to a [limit](https://www.rabbitmq.com/quorum-queues.html#memory-limit))

## [Usage](https://www.rabbitmq.com/quorum-queues.html#usage)

As stated earlier, quorum queues share most of the fundamentals with other [queue](https://www.rabbitmq.com/queues.html) types. A client library that can specify [optional queue arguments](https://www.rabbitmq.com/queues.html#optional-arguments) will be able to use quorum queues.

First we will cover how to declare a quorum queue.

### [Declaring](https://www.rabbitmq.com/quorum-queues.html#declaring)

To declare a quorum queue set the x-queue-type queue argument to quorum (the default is classic). This argument must be provided by a client at queue declaration time; it cannot be set or changed using a [policy](https://www.rabbitmq.com/parameters.html#policies). This is because policy definition or applicable policy can be changed dynamically but queue type cannot. It must be specified at the time of declaration.

Declaring a queue with an x-queue-type argument set to quorum will declare a quorum queue with up to five replicas (default [replication factor](https://www.rabbitmq.com/quorum-queues.html#replication-factor)), one per each [cluster node](https://www.rabbitmq.com/clustering.html).

For example, a cluster of three nodes will have three replicas, one on each node. In a cluster of seven nodes, five nodes will have one replica each but two nodes won't host any replicas.

After declaration a quorum queue can be bound to any exchange just as any other RabbitMQ queue.

If declaring using [management UI](https://www.rabbitmq.com/management.html), queue type must be specified using the queue type drop down menu.

### Client Operations

The following operations work the same way for quorum queues as they do for classic queues:

- [Consumption](https://www.rabbitmq.com/consumers.html) (subscription)
- [Consumer acknowledgements](https://www.rabbitmq.com/confirms.html) (keep [QoS Prefetch Limitations](https://www.rabbitmq.com/quorum-queues.html#global-qos) in mind)
- Cancelation of consumers
- Purging of queue messages
- Queue deletion

With some queue operations there are minor differences:

- [Declaration](https://www.rabbitmq.com/quorum-queues.html#declaring) (covered above)
- Setting [QoS prefetch](https://www.rabbitmq.com/quorum-queues.html#global-qos) for consumers

## [Replication and Data Locality](https://www.rabbitmq.com/quorum-queues.html#replication)

When a quorum queue is declared, an initial number of replicas for it must be started in the cluster. By default the number of replicas to be started is up to three, one per RabbitMQ node in the cluster.

Three nodes is the practical minimum of replicas for a quorum queue. In RabbitMQ clusters with a larger number of nodes, adding more replicas than a [quorum](https://www.rabbitmq.com/quorum-queues.html#what-is-quorum) (majority) will not provide any improvements in terms of [quorum queue availability](https://www.rabbitmq.com/quorum-queues.html#quorum-requirements) but it will consume more cluster resources.

Therefore the **recommended number of replicas** for a quorum queue is the quorum of cluster nodes (but no fewer than three).

### [Controlling the Initial Replication Factor](https://www.rabbitmq.com/quorum-queues.html#replication-factor)

For example, a cluster of three nodes will have three replicas, one on each node. In a cluster of seven nodes, three nodes will have one replica each but four more nodes won't host any replicas of the newly declared queue.

Like with classic mirrored queues, the replication factor (number of replicas a queue has) can be configured for quorum queues.

The minimum factor value that makes practical sense is three. It is highly recommended for the factor to be an odd number. This way a clear quorum (majority) of nodes can be computed. For example, there is no "majority" of nodes in a two node cluster. This is covered with more examples below in the [Fault Tolerance and Minimum Number of Replicas Online](https://www.rabbitmq.com/quorum-queues.html#quorum-requirements) section.

This may not be desirable for larger clusters or for cluster with an even number of nodes. To control the number of quorum queue members set the x-quorum-initial-group-size queue argument when declaring the queue. The group size argument provided should be an integer that is greater than zero and smaller or equal to the current RabbitMQ cluster size. The quorum queue will be launched to run on a random subset of the RabbitMQ cluster.

### [Queue Leader Location](https://www.rabbitmq.com/quorum-queues.html#leader-placement)

Every quorum queue has a primary replica. That replica is called *queue leader* (originally "queue master"). All queue operations go through the leader first and then are replicated to followers (mirrors). This is necessary to guarantee FIFO ordering of messages.

To avoid some nodes in a cluster hosting the majority of queue leader replicas and thus handling most of the load, queue leaders should be reasonably evenly distributed across cluster nodes.

When a new quorum queue is declared, the set of nodes that will host its replicas is randomly picked. Which replica becomes the leader is decided by a [leader election](https://www.rabbitmq.com/quorum-queues.html#leader-election) process which is also based on randomization.

### [Managing Replicas](https://www.rabbitmq.com/quorum-queues.html#replica-management) (Quorum Group Members)

Replicas of a quorum queue are explicitly managed by the operator. When a new node is added to the cluster, it will host no quorum queue replicas unless the operator explicitly adds it to a member (replica) list of a quorum queue or a set of quorum queues.

When a node has to be decomissioned (permanently removed from the cluster), it must be explicitly removed from the member list of all quorum queues it currently hosts replicas for.

Several [CLI commands](https://www.rabbitmq.com/cli.html) are provided to perform the above operations:

```bash
rabbitmq-queues add_member [-p <vhost>] <queue-name> <node>
rabbitmq-queues delete_member [-p <vhost>] <queue-name> <node>
rabbitmq-queues grow <node> <all | even> [--vhost-pattern <pattern>] [--queue-pattern <pattern>]
rabbitmq-queues shrink <node> [--errors-only]
```

To successfully add and remove members a quorum of replicas in the cluster must be available because cluster membership changes are treated as queue state changes.

Care needs to be taken not to accidentally make a queue unavailable by losing the quorum whilst performing maintenance operations that involve membership changes.

When replacing a cluster node, it is safer to first add a new node and then decomission the node it replaces.

### [Rebalancing Replicas](https://www.rabbitmq.com/quorum-queues.html#replica-rebalancing)

Once declared, the RabbitMQ nodes a quorum queue resides on won't change even if the members of the RabbitMQ cluster change (e.g. a node is decomissioned or added). To re-balance after a RabbitMQ cluster change quorum queues will have to be manually adjusted using the rabbitmq-queues [command line tool](https://www.rabbitmq.com/cli.html):

```bash
# rebalances all quorum queues
rabbitmq-queues rebalance quorum
```

it is possible to rebalance a subset of queues selected by name:

```bash
# rebalances a subset of quorum queues
rabbitmq-queues rebalance quorum --queue-pattern "orders.*"
```

or quorum queues in a particular set of virtual hosts:

```bash
# rebalances a subset of quorum queues
rabbitmq-queues rebalance quorum --vhost-pattern "production.*"
```

## [Behaviour](https://www.rabbitmq.com/quorum-queues.html#behaviour)

A quorum queue relies on a consensus protocol called Raft to ensure data consistency and safety.

Every quorum queue has a primary replica (a *leader* in Raft parlance) and zero or more secondary replicas (called *followers*).

A leader is elected when the cluster is first formed and later if the leader becomes unavailable.

### [Leader Election and Failure Handling](https://www.rabbitmq.com/quorum-queues.html#leader-election)

A quorum queue requires a quorum of the declared nodes to be available to function. When a RabbitMQ node hosting a quorum queue's *leader* fails or is stopped another node hosting one of that quorum queue's *follower* will be elected leader and resume operations.

Failed and rejoining followers will re-synchronise ("catch up") with the leader. In contrast to classic mirrored queues, a temporary replica failure does not require a full re-synchronization from the currently elected leader. Only the delta will be transferred if a re-joining replica is behind the leader. This "catching up" process does not affect leader availability.

Except for the initial replica set selection, replicas must be explicitly added to a quorum queue. When a new replica is [added](https://www.rabbitmq.com/quorum-queues.html#replica-management), it will synchronise the entire queue state from the leader, similarly to classic mirrored queues.

### [Fault Tolerance and Minimum Number of Replicas Online](https://www.rabbitmq.com/quorum-queues.html#quorum-requirements)

Consensus systems can provide certain guarantees with regard to data safety. These guarantees do mean that certain conditions need to be met before they become relevant such as requiring a minimum of three cluster nodes to provide fault tolerance and requiring more than half of members to be available to work at all.

Failure tolerance characteristics of clusters of various size can be described in a table:

| Cluster node count | Tolerated number of node failures | Tolerant to a network partition      |
| :----------------- | :-------------------------------- | :----------------------------------- |
| 1                  | 0                                 | not applicable                       |
| 2                  | 0                                 | no                                   |
| 3                  | 1                                 | yes                                  |
| 4                  | 1                                 | yes if a majority exists on one side |
| 5                  | 2                                 | yes                                  |
| 6                  | 2                                 | yes if a majority exists on one side |
| 7                  | 3                                 | yes                                  |
| 8                  | 3                                 | yes if a majority exists on one side |
| 9                  | 4                                 | yes                                  |

As the table above shows RabbitMQ clusters with fewer than three nodes do not benefit fully from the quorum queue guarantees. RabbitMQ clusters with an even number of RabbitMQ nodes do not benefit from having quorum queue members spread over all nodes. For these systems the quorum queue size should be constrained to a smaller uneven number of nodes.

Performance tails off quite a bit for quorum queue node sizes larger than 5. We do not recommend running quorum queues on more than 7 RabbitMQ nodes. The default quorum queue size is 3 and is controllable using the x-quorum-initial-group-size [queue argument](https://www.rabbitmq.com/queues.html#optional-arguments).

### [Data Safety](https://www.rabbitmq.com/quorum-queues.html#data-safety)

Quorum queues are designed to provide data safety under network partition and failure scenarios. A message that was successfully confirmed back to the publisher using the [publisher confirms](https://www.rabbitmq.com/confirms.html) feature should not be lost as long as at least a majority of RabbitMQ nodes hosting the quorum queue are not permanently made unavailable.

Generally quorum queues favours data consistency over availability.

**No guarantees are provided for messages that have not been confirmed using the publisher confirm mechanism**. Such messages could be lost "mid-way", in an operating system buffer or otherwise fail to reach the queue leader.

### [Availability](https://www.rabbitmq.com/quorum-queues.html#availability)

A quorum queue should be able to tolerate a minority of queue members becoming unavailable with no or little effect on availability.

Note that depending on the [partition handling strategy](https://www.rabbitmq.com/partitions.html) used RabbitMQ may restart itself during recovery and reset the node but as long as that does not happen, this availability guarantee should hold true.

For example, a queue with three replicas can tolerate one node failure without losing availability. A queue with five replicas can tolerate two, and so on.

If a quorum of nodes cannot be recovered (say if 2 out of 3 RabbitMQ nodes are permanently lost) the queue is permanently unavailable and will need to be force deleted and recreated.

Quorum queue follower replicas that are disconnected from the leader or participating in a leader election will ignore queue operations sent to it until they become aware of a newly elected leader. There will be warnings in the log (received unhandled msg and similar) about such events. As soon as the replica discovers a newly elected leader, it will sync the queue operation log entries it does not have from the leader, including the dropped ones. Quorum queue state will therefore remain consistent.

## [Performance Characteristics](https://www.rabbitmq.com/quorum-queues.html#performance)

Quorum queues are designed to trade latency for throughput and have been tested and compared against durable [classic mirrored queues](https://www.rabbitmq.com/ha.html) in 3, 5 and 7 node configurations at several message sizes. In scenarios using both consumer acks and publisher confirms quorum queues have been observed to have equal or greater throughput to classic mirrored queues.

As quorum queues persist all data to disks before doing anything it is recommended to use the fastest disks possible. Quorum queues also benefit from consumers using higher prefetch values to ensure consumers aren't starved whilst acknowledgements are flowing through the system and allowing messages to be delivered in a timely fashion.

Due to the disk I/O-heavy nature of quorum queues, their throughput decreases as message sizes increase.

Just like mirrored queues, quorum queues are also affected by cluster sizes. The more replicas a quorum queue has, the lower its throughput generally will be since more work has to be done to replicate data and achieve consensus.

## [Configuration](https://www.rabbitmq.com/quorum-queues.html#configuration)

There are a few new configuration parameters that can be tweaked using the [advanced](https://www.rabbitmq.com/configure.html#advanced-config-file) config file.

Note that all settings related to [resource footprint](https://www.rabbitmq.com/quorum-queues.html#resource-use) are documented in a separate section.

The ra application (which is the Raft library that quorum queues use) has [its own set of tunable parameters](https://github.com/rabbitmq/ra#configuration).

The rabbit application has several quorum queue related configuration items available.

| advanced.config Configuration Key | Description                                                  | Default value |
| --------------------------------- | ------------------------------------------------------------ | ------------- |
| rabbit.quorum_cluster_size        | Sets the default quorum queue cluster size (can be over-ridden by the x-quorum-initial-group-size queue argument at declaration time. | 3             |
| rabbit.quorum_commands_soft_limit | This is a flow control related parameter defining the maximum number of unconfirmed messages a channel accepts before entering flow. The current default is configured to provide good performance and stability when there are multiple publishers sending to the same quorum queue. If the applications typically only have a single publisher per queue this limit could be increased to provide somewhat better ingress rates. | 32            |

### Example

The following advanced.config example modifies all values listed above:

```erlang
[
 %% five replicas by default, only makes sense for nine node clusters
 {rabbit, [{quorum_cluster_size, 5},
           {quorum_commands_soft_limit, 512}]}
]
```

## [Poison Message Handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling)

Quorum queue support handling of [poison messages](https://en.wikipedia.org/wiki/Poison_message), that is, messages that cause a consumer to repeatedly requeue a delivery (possibly due to a consumer failure) such that the message is never consumed completely and [positively acknowledged](https://www.rabbitmq.com/confirms.html) so that it can be marked for deletion by RabbitMQ.

Quorum queues keep track of the number of unsuccessful delivery attempts and expose it in the "x-delivery-count" header that is included with any redelivered message.

It is possible to set a delivery limit for a queue using a [policy](https://www.rabbitmq.com/parameters.html#policies) argument, delivery-limit.

When a message has been returned more times than the limit the message will be dropped or [dead-lettered](https://www.rabbitmq.com/dlx.html) (if a DLX is configured).

## [Resource Use](https://www.rabbitmq.com/quorum-queues.html#resource-use)

Quorum queues typically require more resources (disk and RAM) than classic mirrored queues. To enable fast election of a new leader and recovery, data safety as well as good throughput characteristics all members in a quorum queue "cluster" keep all messages in the queue in memory *and* on disk.

Quorum queues use a write-ahead-log (WAL) for all operations. WAL operations are stored both in memory and written to disk. When the current WAL file reaches a predefined limit, it is flushed to a WAL segment file on disk and the system will begin to release the memory used by that batch of log entries. The segment files are then compacted over time as consumers [acknowledge deliveries](https://www.rabbitmq.com/confirms.html). Compaction is the process that reclaims disk space.

The WAL file size limit at which it is flushed to disk can be controlled:

```ini
# Flush current WAL file to a segment file on disk once it reaches 64 MiB in size
raft.wal_max_size_bytes = 64000000
```

The value defaults to 512 MiB. This means that during steady load, the WAL table memory footprint can reach 512 MiB.

Because memory deallocation may take some time, we recommend that the RabbitMQ node is allocated at least 3 times the memory of the default WAL file size limit. More will be required in high-throughput systems. 4 times is a good starting point for those.

### [Configuring Per Queue Memory Limit](https://www.rabbitmq.com/quorum-queues.html#memory-limit)

It is possible to limit the amount of memory each quorum queue will use for the part of its log that is kept in memory. Note that these limits are different from those of the [in-memory Raft WAL table](https://www.rabbitmq.com/quorum-queues.html#resource-use) and [queue length limits](https://www.rabbitmq.com/maxlength.html).

The limit is controlled using [optional queue arguments](https://www.rabbitmq.com/queues.html#optional-arguments) that are best configured using a [policy](https://www.rabbitmq.com/parameters.html#policies).

- x-max-in-memory-length sets a limit as a number of messages. Must be a non-negative integer.
- x-max-in-memory-bytes sets a limit as the total size of message bodies (payloads), in bytes. Must be a non-negative integer.

### [Repeated Requeues](https://www.rabbitmq.com/quorum-queues.html#repeated-requeues)

Internally quorum queues are implemented using a log where all operations including messages are persisted. To avoid this log growing too large it needs to be truncated regularly. To be able to truncate a section of the log all messages in that section needs to be acknowledged. Usage patterns that continuously [reject or nack](https://www.rabbitmq.com/nack.html) the same message with the requeue flag set to true could cause the log to grow in an unbounded fashion and eventually fill up the disks.

### [Increased Atom Use](https://www.rabbitmq.com/quorum-queues.html#atom-use)

The internal implementation of quorum queues converts the queue name into an Erlang atom. If queues with arbitrary names are continuously created and deleted it *may* threaten the long term stability of the RabbitMQ system (if the size of the atom table reaches the maximum limit, about 1M by default). It is not recommended to use quorum queues in this manner at this point.

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!