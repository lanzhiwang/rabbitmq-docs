# Publish message

1. Routing key:

2. Headers:  Headers can have any name. Only long string headers can be set here.  标题可以有任何名称。 此处只能设置长字符串标题。

3. Properties:  You can set other message properties here (delivery mode and headers are pulled out as the most common cases).  您可以在此处设置其他消息属性（最常见的情况是删除传递模式和标题）。Invalid properties will be ignored. Valid properties are:  无效的属性将被忽略。 有效的属性是：
    - content_type
    - content_encoding
    - priority
    - correlation_id
    - reply_to
    - expiration
    - message_id
    - timestamp
    - type
    - user_id
    - app_id
    - cluster_id

4. Payload:

# Add a new exchange

1. Name:

2. Type:
    * direct
    * fanout
    * headers
    * topic

3. Durability:
    * Durable
    * Transient

4. Auto delete:
    * No
    * Yes

5. Internal:  If yes, clients cannot publish to this exchange directly. It can only be used with exchange to exchange bindings.
    * No
    * Yes

6. Arguments:  If messages to this exchange cannot otherwise be routed, send them to the alternate exchange named here. (Sets the "[alternate-exchange](https://rabbitmq.com/ae.html)" argument.)
    * alternate-exchange

# Add a new queue(Type:  Classic)

1. Type:
    * Classic

2. Name:

3. Durability:
    * Durable
    * Transient

4. Node:

5. Auto delete:
    * No
    * Yes

6. Arguments:
    * **x-message-ttl**  How long a message published to a queue can live before it is discarded (milliseconds). (Sets the "[x-message-ttl](https://rabbitmq.com/ttl.html#per-queue-message-ttl)" argument.)

    * **x-expires**  How long a queue can be unused for before it is automatically deleted (milliseconds). (Sets the "[x-expires](https://rabbitmq.com/ttl.html#queue-ttl)" argument.)

    * **x-max-length**  How many (ready) messages a queue can contain before it starts to drop them from its head. (Sets the "[x-max-length](https://rabbitmq.com/maxlength.html)" argument.)

    * **x-max-length-bytes**  Total body size for ready messages a queue can contain before it starts to drop them from its head. (Sets the "[x-max-length-bytes](https://rabbitmq.com/maxlength.html)" argument.)

    * **x-overflow**   Sets the [queue overflow behaviour](https://www.rabbitmq.com/maxlength.html#overflow-behaviour). This determines what happens to messages when the maximum length of a queue is reached. Valid values are `drop-head`, `reject-publish` or `reject-publish-dlx`. The quorum queue type only supports `drop-head` and `reject-publish`.

    * **x-dead-letter-exchange**  Optional name of an exchange to which messages will be republished if they are rejected or expire. (Sets the "[x-dead-letter-exchange](https://rabbitmq.com/dlx.html)" argument.)

    * **x-dead-letter-routing-key**  Optional replacement routing key to use when a message is dead-lettered. If this is not set, the message's original routing key will be used. (Sets the "[x-dead-letter-routing-key](https://rabbitmq.com/dlx.html)" argument.)

    * **x-single-active-consumer**  If set, makes sure only one consumer at a time consumes from the queue and fails over to another registered consumer in case the active one is cancelled or dies. (Sets the "[x-single-active-consumer](https://rabbitmq.com/consumers.html#single-active-consumer)" argument.)

    * **x-max-priority**  Maximum number of priority levels for the queue to support; if not set, the queue will not support message priorities. (Sets the "[x-max-priority](https://rabbitmq.com/priority.html)" argument.)

    * **x-queue-mode = lazy**  Set the queue into lazy mode, keeping as many messages as possible on disk to reduce RAM usage; if not set, the queue will keep an in-memory cache to deliver messages as fast as possible. (Sets the "[x-queue-mode](https://www.rabbitmq.com/lazy-queues.html)" argument.)

    * **x-queue-master-locator**  Set the queue into master location mode, determining the rule by which the queue master is located when declared on a cluster of nodes. (Sets the "[x-queue-master-locator](https://www.rabbitmq.com/ha.html)" argument.)

# Add a new queue(Type:  Quorum)

1. Type:
    * Quorum

2. Name:

3. Node:

4. Arguments:

    * **x-expires**  How long a queue can be unused for before it is automatically deleted (milliseconds). (Sets the "[x-expires](https://rabbitmq.com/ttl.html#queue-ttl)" argument.)

    * **x-max-length**  How many (ready) messages a queue can contain before it starts to drop them from its head. (Sets the "[x-max-length](https://rabbitmq.com/maxlength.html)" argument.)

    * **x-max-length-bytes**  Total body size for ready messages a queue can contain before it starts to drop them from its head. (Sets the "[x-max-length-bytes](https://rabbitmq.com/maxlength.html)" argument.)

    * **x-delivery-limit**  The number of allowed unsuccessful delivery attempts. Once a message has been delivered unsucessfully this many times it will be dropped or dead-lettered, depending on the queue configuration.

    * **x-overflow**  Sets the [queue overflow behaviour](https://www.rabbitmq.com/maxlength.html#overflow-behaviour). This determines what happens to messages when the maximum length of a queue is reached. Valid values are `drop-head`, `reject-publish` or `reject-publish-dlx`. The quorum queue type only supports `drop-head` and `reject-publish`.

    * **x-dead-letter-exchange**  Optional name of an exchange to which messages will be republished if they are rejected or expire. (Sets the "[x-dead-letter-exchange](https://rabbitmq.com/dlx.html)" argument.)

    * **x-dead-letter-routing-key**  Optional replacement routing key to use when a message is dead-lettered. If this is not set, the message's original routing key will be used. (Sets the "[x-dead-letter-routing-key](https://rabbitmq.com/dlx.html)" argument.)

    * **x-single-active-consumer**  If set, makes sure only one consumer at a time consumes from the queue and fails over to another registered consumer in case the active one is cancelled or dies. (Sets the "[x-single-active-consumer](https://rabbitmq.com/consumers.html#single-active-consumer)" argument.)

    * **x-max-in-memory-length**  How many (ready) messages a quorum queue can contain in memory before it starts storing them on disk only. (Sets the x-max-in-memory-length argument.)

    * **x-max-in-memory-bytes**  Total body size for ready messages a quorum queue can contain in memory before it starts storing them on disk only. (Sets the x-max-in-memory-bytes argument.)



# Add / update a policy

1. Name:

2. Pattern:

3. Apply to:
    * Exchanges and queues
    * Exchanges
    * Queues

4. Priority:

5. Definition:

    * Queues [All types]
      * **max-length**
      * **max-length-bytes**
      * **overflow**  Sets the [queue overflow behaviour](https://www.rabbitmq.com/maxlength.html#overflow-behaviour). This determines what happens to messages when the maximum length of a queue is reached. Valid values are `drop-head`, `reject-publish` or `reject-publish-dlx`. The quorum queue type only supports `drop-head` and `reject-publish`.
      * **expires**
      * **dead-letter-exchange**
      * **dead-letter-routing-key**

    * Queues [Classic]
      * **ha-mode**  One of `all` (mirror to all nodes in the cluster), `exactly` (mirror to a set number of nodes) or `nodes` (mirror to an explicit list of nodes). If you choose one of the latter two, you must also set `ha-params`.
      * **ha-params**  Absent if `ha-mode` is `all`, a number if `ha-mode` is `exactly`, or a list of strings if `ha-mode` is `nodes`.
      * **ha-sync-mode**  One of `manual` or `automatic`. [Learn more](https://www.rabbitmq.com/ha.html#unsynchronised-mirrors)
      * **ha-promote-on-shutdown**  One of `when-synced` or `always`. [Learn more](https://www.rabbitmq.com/ha.html#unsynchronised-mirrors)
      * **ha-promote-on-failure**  One of `when-synced` or `always`. [Learn more](https://www.rabbitmq.com/ha.html#unsynchronised-mirrors)
      * **message-ttl**
      * **queue-mode = lazy**
      * **queue-master-locator**

    * Queues [Quorum]
      * **max-in-memory-length**  How many (ready) messages a quorum queue can contain in memory before it starts storing them on disk only. (Sets the x-max-in-memory-length argument.)
      * **max-in-memory-bytes**  Total body size for ready messages a quorum queue can contain in memory before it starts storing them on disk only. (Sets the x-max-in-memory-bytes argument.)
      * **delivery-limit**  The number of allowed unsuccessful delivery attempts. Once a message has been delivered unsucessfully this many times it will be dropped or dead-lettered, depending on the queue configuration.

    * Exchanges
      * **alternate-exchange**  If messages to this exchange cannot otherwise be routed, send them to the alternate exchange named here. (Sets the "[alternate-exchange](https://rabbitmq.com/ae.html)" argument.)

    * Federation
      * **federation-upstream-set**  A string; only if the federation plugin is enabled. Chooses the name of a set of upstreams to use with federation, or "all" to use all upstreams. Incompatible with `federation-upstream`.
      * **federation-upstream**  A string; only if the federation plugin is enabled. Chooses a specific upstream set to use for federation. Incompatible with `federation-upstream-set`.

# Add / update an operator policy

1. Name:

2. Pattern:

3. Apply to:
    * Queues

4. Priority:

5. Definition:

    * Queues [All types]
      * **max-length**
      * **max-length-bytes**
      * **overflow**  Sets the [queue overflow behaviour](https://www.rabbitmq.com/maxlength.html#overflow-behaviour). This determines what happens to messages when the maximum length of a queue is reached. Valid values are `drop-head`, `reject-publish` or `reject-publish-dlx`. The quorum queue type only supports `drop-head` and `reject-publish`.

    * Queues [Classic]
      * **message-ttl**
      * **expires**

    * Queues [Quorum]
      * **max-in-memory-length**  How many (ready) messages a quorum queue can contain in memory before it starts storing them on disk only. (Sets the x-max-in-memory-length argument.)
      * **max-in-memory-bytes**  Total body size for ready messages a quorum queue can contain in memory before it starts storing them on disk only. (Sets the x-max-in-memory-bytes argument.)
      * **delivery-limit**  The number of allowed unsuccessful delivery attempts. Once a message has been delivered unsucessfully this many times it will be dropped or dead-lettered, depending on the queue configuration.

# Feature Flags

| Name | State | Description |
| ---- | ----- | ----------- |
| drop_unroutable_metric | Enabled | Count unroutable publishes to be dropped in stats |
| empty_basic_get_metric | Enabled | Count AMQP `basic.get` on empty queues in stats |
| implicit_default_bindings | Enabled | Default bindings are now implicit, instead of being stored in the database |
| maintenance_mode_status | Enabled | Maintenance mode status |
| quorum_queue | Enabled | Support queues of type `quorum`[[Learn more\]](https://www.rabbitmq.com/quorum-queues.html) |
| user_limits | Enabled | Configure connection and channel limits for a user |
| virtual_host_metadata | Enabled | Virtual host metadata (description, tags, etc) |


# Add a user

1. Username:

2. Password:

3. Tags:
    * Admin
    * Monitoring
    * Policymaker
    * Management
    * Impersonator
    * None

Comma-separated list of tags to apply to the user. Currently [supported by the management plugin](https://www.rabbitmq.com/management.html#permissions):

- **management**  User can access the management plugin

- **policymaker**  User can access the management plugin and manage policies and parameters for the vhosts they have access to.

- **monitoring**  User can access the management plugin and see all connections and channels as well as node-related information.

- **administrator**  User can do everything monitoring can do, manage users, vhosts and permissions, close other user's connections, and manage policies and parameters for all vhosts.

Note that you can set any tag here; the links for the above four tags are just for convenience.

# Add a new virtual host

1. Name:

2. Description:

3. Tags:

# Set / update a virtual host limit

1. Virtual host:

2. Limit:
    * max-connections
    * max-queues

3. Value:


# Set / update a user limit

1. User:

2. Limit:
    * max-connections
    * max-channels

3. Value:

# Connection

1. Details
    * Node
    * Client-provided name
    * Username
    * Protocol
    * Connected at
    * Authentication
    * State
    * Heartbeat
    * Frame max
    * Channel limit

2. Channels

3. Client properties
    * authentication_failure_close:	true
    * basic.nack:	true
    * connection.blocked:	true
    * consumer_cancel_notify:	true
    * exchange_exchange_bindings:	true
    * publisher_confirms:	true

# Channel

1. Details

    * Connection
    * Node
    * Username
    * Mode  Channel guarantee mode. Can be one of the following, or neither:
      * C – [confirm](https://www.rabbitmq.com/confirms.html) Channel will send streaming publish confirmations.
      * T – [transactional](https://www.rabbitmq.com/amqp-0-9-1-reference.html#class.tx) Channel is transactional.
    * State
    * Prefetch count
    * Global prefetch count
    * Messages unacknowledged
    * Messages unconfirmed
    * Messages uncommitted
    * Acks uncommitted
    * Pending Raft commands

2. Consumers

