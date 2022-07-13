# Clustering and Network Partitions

https://www.rabbitmq.com/partitions.html

## Introduction

This guide covers one specific aspect of clustering: network failures between nodes, their effects and recovery options. For a general overview of clustering, see [Clustering](https://www.rabbitmq.com/clustering.html) and [Peer Discovery and Cluster Formation](https://www.rabbitmq.com/cluster-formation.html) guides.  本指南涵盖集群的一个特定方面：节点之间的网络故障、它们的影响和恢复选项。 有关集群的一般概述，请参阅集群和对等发现和集群形成指南。

Clustering can be used to achieve different goals: increased data safety through replication, increased availability for client operations, higher overall throughput and so on. Different configurations are optimal for different purposes.  集群可用于实现不同的目标：通过复制提高数据安全性、提高客户端操作的可用性、提高整体吞吐量等。 不同的配置对于不同的目的是最佳的。

Network connection failures between cluster members have an effect on data consistency and availability (as in the CAP theorem) to client operations. Since different applications have different requirements around consistency and can tolerate unavailability to a different extent, different partition handling strategies are available.  集群成员之间的网络连接故障会影响客户端操作的数据一致性和可用性（如 CAP 定理）。 由于不同的应用程序对一致性有不同的要求，并且可以在不同程度上容忍不可用性，因此可以使用不同的分区处理策略。

## Detecting Network Partitions  检测网络分区

Nodes determine if its peer is down if another node is unable to contact it for a [period of time](https://www.rabbitmq.com/nettick.html), 60 seconds by default. If two nodes come back into contact, both having thought the other is down, the nodes will determine that a partition has occurred. This will be written to the RabbitMQ log in a form like:  如果另一个节点在一段时间内（默认为 60 秒）无法与其联系，则节点会确定其对等方是否关闭。 如果两个节点重新接触，并且都认为另一个节点已关闭，则节点将确定发生了分区。 这将以如下形式写入 RabbitMQ 日志：

```
2020-05-18 06:55:37.324 [error] <0.341.0> Mnesia(rabbit@warp10): ** ERROR ** mnesia_event got {inconsistent_database, running_partitioned_network, rabbit@hostname2}
```

Partition presence can be identified via server [logs](https://www.rabbitmq.com/logging.html), [HTTP API](https://www.rabbitmq.com/management.html) (for [monitoring](https://www.rabbitmq.com/monitoring.html)) and a [CLI command](https://www.rabbitmq.com/cli.html):

```bash
rabbitmq-diagnostics cluster_status
```

rabbitmq-diagnostics cluster_status will normally show an empty list for partitions:

```bash
rabbitmq-diagnostics cluster_status
# => Cluster status of node rabbit@warp10 ...
# => Basics
# =>
# => Cluster name: local.1
# =>
# => ...edited out for brevity...
# =>
# => Network Partitions
# =>
# => (none)
# =>
# => ...edited out for brevity...
```

However, if a network partition has occurred then information about partitions will appear there:

```bash
rabbitmqctl cluster_status
# => Cluster status of node rabbit@warp10 ...
# => Basics
# =>
# => Cluster name: local.1
# =>
# => ...edited out for brevity...
# =>
# => Network Partitions
# =>
# => Node flopsy@warp10 cannot communicate with hare@warp10
# => Node rabbit@warp10 cannot communicate with hare@warp10
```

The HTTP API will return partition information for each node under partitions in GET /api/nodes endpoints.

The management UI will show a warning on the overview page if a partition has occurred.

## [Behavior During a Network Partition](https://www.rabbitmq.com/partitions.html#during)

While a network partition is in place, the two (or more!) sides of the cluster can evolve independently, with both sides thinking the other has crashed. This scenario is known as split-brain. Queues, bindings, exchanges can be created or deleted separately.

[Classic mirrored queues](https://www.rabbitmq.com/ha.html) which are split across the partition will end up with one leader on each side of the partition, again with both sides acting independently. [Quorum queues](https://www.rabbitmq.com/quorum-queues.html) will elect a new leader on the majority side. Quorum queue replicas on the minority side will no longer make progress (i.e. accept new messages, deliver to consumers, etc), all this work will be done by the new leader.

Unless a [partition handling strategy](https://www.rabbitmq.com/partitions.html#automatic-handling), such as pause_minority, is configured to be used, the split will continue even after network connectivity is restored.

## [Partitions Caused by Suspend and Resume](https://www.rabbitmq.com/partitions.html#suspend)

While we refer to "network" partitions, really a partition is any case in which the different nodes of a cluster can have communication interrupted without any node failing. In addition to network failures, suspending and resuming an entire OS can also cause partitions when used against running cluster nodes - as the suspended node will not consider itself to have failed, or even stopped, but the other nodes in the cluster will consider it to have done so.

While you could suspend a cluster node by running it on a laptop and closing the lid, the most common reason for this to happen is for a virtual machine to have been suspended by the hypervisor.

While it's fine to run RabbitMQ clusters in virtualised environments or containers, **make sure that VMs are not suspended while running**.

Note that some virtualisation features such as migration of a VM from one host to another will tend to involve the VM being suspended.

Partitions caused by suspend and resume will tend to be asymmetrical - the suspended node will not necessarily see the other nodes as having gone down, but will be seen as down by the rest of the cluster. This has particular implications for [pause_minority](https://www.rabbitmq.com/partitions.html#pause-minority) mode.

## [Recovering From a Split-Brain](https://www.rabbitmq.com/partitions.html#recovering)

To recover from a split-brain, first choose one partition which you trust the most. This partition will become the authority for the state of the system (schema, messages) to use; any changes which have occurred on other partitions will be lost.

Stop all nodes in the other partitions, then start them all up again. When they [rejoin the cluster](https://www.rabbitmq.com/clustering.html#restarting) they will restore state from the trusted partition.

Finally, you should also restart all the nodes in the trusted partition to clear the warning.

It may be simpler to stop the whole cluster and start it again; if so make sure that the **first** node you start is from the trusted partition.

## [Partition Handling Strategies](https://www.rabbitmq.com/partitions.html#automatic-handling)

RabbitMQ also offers three ways to deal with network partitions automatically: pause-minority mode, pause-if-all-down mode and autoheal mode. The default behaviour is referred to as ignore mode.

In pause-minority mode RabbitMQ will automatically pause cluster nodes which determine themselves to be in a minority (i.e. fewer or equal than half the total number of nodes) after seeing other nodes go down. It therefore chooses partition tolerance over availability from the CAP theorem. This ensures that in the event of a network partition, at most the nodes in a single partition will continue to run. The minority nodes will pause as soon as a partition starts, and will start again when the partition ends. This configuration prevents split-brain and is therefore able to automatically recover from network partitions without inconsistencies.

In pause-if-all-down mode, RabbitMQ will automatically pause cluster nodes which cannot reach any of the listed nodes. In other words, all the listed nodes must be down for RabbitMQ to pause a cluster node. This is close to the pause-minority mode, however, it allows an administrator to decide which nodes to prefer, instead of relying on the context. For instance, if the cluster is made of two nodes in rack A and two nodes in rack B, and the link between racks is lost, pause-minority mode will pause all nodes. In pause-if-all-down mode, if the administrator listed the two nodes in rack A, only nodes in rack B will pause. Note that it is possible the listed nodes get split across both sides of a partition: in this situation, no node will pause. That is why there is an additional *ignore*/*autoheal* argument to indicate how to recover from the partition.

In autoheal mode RabbitMQ will automatically decide on a winning partition if a partition is deemed to have occurred, and will restart all nodes that are not in the winning partition. Unlike pause_minority mode it therefore takes effect when a partition ends, rather than when one starts.

The winning partition is the one which has the most clients connected (or if this produces a draw, the one with the most nodes; and if that still produces a draw then one of the partitions is chosen in an unspecified way).

You can enable either mode by setting the configuration parameter cluster_partition_handling for the rabbit application in the [configuration file](https://www.rabbitmq.com/configure.html#configuration-file) to:

- autoheal
- pause_minority
- pause_if_all_down

If using the pause_if_all_down mode, additional parameters are required:

- nodes: nodes which should be unavailable to pause
- recover: recover action, can be ignore or autoheal

Example [config snippet](https://www.rabbitmq.com/configure.html#config-file) that uses pause_if_all_down:

```plaintext
cluster_partition_handling = pause_if_all_down

## Recovery strategy. Can be either 'autoheal' or 'ignore'
cluster_partition_handling.pause_if_all_down.recover = ignore

## Node names to check
cluster_partition_handling.pause_if_all_down.nodes.1 = rabbit@myhost1
cluster_partition_handling.pause_if_all_down.nodes.2 = rabbit@myhost2
```

### [Which Mode to Pick?](https://www.rabbitmq.com/partitions.html#options)

It's important to understand that allowing RabbitMQ to deal with network partitions automatically comes with trade offs.

As stated in the introduction, to connect RabbitMQ clusters over generally unreliable links, prefer [Federation](https://www.rabbitmq.com/federation.html) or the [Shovel](https://www.rabbitmq.com/shovel.html).

With that said, here are some guidelines to help the operator determine which mode may or may not be appropriate:

- ignore: use when network reliability is the highest practically possible and node availability is of topmost importance. For example, all cluster nodes can be in the same rack or equivalent, connected with a switch, and that switch is also the route to the outside world.
- pause_minority: appropriate when clustering across racks or availability zones in a single region, and the probability of losing a majority of nodes (zones) at once is considered to be very low. This mode trades off some availability for the ability to automatically recover if/when the lost node(s) come back.
- autoheal: appropriate when are more concerned with continuity of service than with data consistency across nodes.

### [More About Pause-minority Mode](https://www.rabbitmq.com/partitions.html#pause-minority)

The Erlang VM on the paused nodes will continue running but the nodes will not listen on any ports or be otherwise available. They will check once per second to see if the rest of the cluster has reappeared, and start up again if it has.

Note that nodes will not enter the paused state at startup, even if they are in a minority then. It is expected that any such minority at startup is due to the rest of the cluster not having been started yet.

Also note that RabbitMQ will pause nodes which are not in a *strict* majority of the cluster - i.e. containing more than half of all nodes. It is therefore not a good idea to enable pause-minority mode on a cluster of two nodes since in the event of any network partition **or node failure**, both nodes will pause. However, pause_minority mode is safer than ignore mode, with regards to integrity. For clusters of more than two nodes, especially if the most likely form of network partition is that a single minority of nodes drops off the network, the availability remains as good as with ignore mode.

Note that pause_minority mode will do nothing to defend against partitions caused by cluster nodes being [suspended](https://www.rabbitmq.com/partitions.html#suspend). This is because the suspended node will never see the rest of the cluster vanish, so will have no trigger to disconnect itself from the cluster.

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!