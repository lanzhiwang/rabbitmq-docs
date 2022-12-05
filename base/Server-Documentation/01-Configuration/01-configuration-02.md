### Core Server Variables Configurable in rabbitmq.conf  在 rabbitmq.conf 中可配置的核心服务器变量

These variables are the most common. The list is not complete, as some settings are quite obscure.  这些变量是最常见的。 该列表并不完整，因为某些设置非常晦涩。

* listeners

Ports or hostname/pair on which to listen for "plain" AMQP 0-9-1 and AMQP 1.0 connections (without [TLS](https://www.rabbitmq.com/ssl.html)). See the [Networking guide](https://www.rabbitmq.com/networking.html) for more details and examples.

Default: `listeners.tcp.default = 5672`

* num_acceptors.tcp

Number of Erlang processes that will accept connections for the TCP listeners.

Default: `num_acceptors.tcp = 10`

* handshake_timeout

Maximum time for AMQP 0-9-1 handshake (after socket connection and TLS handshake), in milliseconds.

Default: `handshake_timeout = 10000`

* listeners.ssl

Ports or hostname/pair on which to listen for TLS-enabled AMQP 0-9-1 and AMQP 1.0 connections. See the [TLS guide](https://www.rabbitmq.com/ssl.html) for more details and examples.

Default: none (not set)

* num_acceptors.ssl

Number of Erlang processes that will accept TLS connections from clients.

Default: `num_acceptors.ssl = 10`

* ssl_options

TLS configuration. See the [TLS guide](https://www.rabbitmq.com/ssl.html#enabling-ssl).

Default: `ssl_options = none`

* ssl_handshake_timeout

TLS handshake timeout, in milliseconds.

Default: `ssl_handshake_timeout = 5000`

* vm_memory_high_watermark

Memory threshold at which the flow control is triggered. Can be absolute or relative to the amount of RAM available to the OS: 

`vm_memory_high_watermark.relative = 0.6` 

`vm_memory_high_watermark.absolute = 2GB`

See the [memory-based flow control](https://www.rabbitmq.com/memory.html) and [alarms](https://www.rabbitmq.com/alarms.html) documentation.

Default: `vm_memory_high_watermark.relative = 0.4`

* vm_memory_calculation_strategy

Strategy for memory usage reporting. Can be one of the following:

`allocated`: uses Erlang memory allocator statistics  使用 Erlang 内存分配器统计信息

`rss`: uses operating system RSS memory reporting. This uses OS-specific means and may start short lived child processes.  使用操作系统的 RSS 内存报告。 这使用特定于操作系统的方法，并且可能会启动短暂的子进程。

`legacy`: uses legacy memory reporting (how much memory is considered to be used by the runtime). This strategy is fairly inaccurate.  使用遗留内存报告（运行时认为使用了多少内存）。 这种策略是相当不准确的。

`erlang`: same as legacy, preserved for backwards compatibility  与 legacy 相同，为向后兼容而保留

Default: `vm_memory_calculation_strategy = allocated`

* vm_memory_high_watermark_paging_ratio

Fraction of the high watermark limit at which queues start to page messages out to disc to free up memory. See the [memory-based flow control](https://www.rabbitmq.com/memory.html) documentation.  队列开始将消息分页到磁盘以释放内存的高水位线限制的一部分。 请参阅基于内存的流控制文档。

Default: `vm_memory_high_watermark_paging_ratio = 0.5`

* total_memory_available_override_value

Makes it possible to override the total amount of memory available, as opposed to inferring it from the environment using OS-specific means. This should only be used when actual maximum amount of RAM available to the node doesn't match the value that will be inferred by the node, e.g. due to containerization or similar constraints the node cannot be aware of. The value may be set to an integer number of bytes or, alternatively, in information units (e.g `8GB`). For example, when the value is set to 4 GB, the node will believe it is running on a machine with 4 GB of RAM.  可以覆盖可用内存总量，而不是使用特定于操作系统的方法从环境中推断出来。 这仅应在节点可用的实际最大 RAM 量与节点推断的值不匹配时使用，例如 由于容器化或类似的限制，节点无法意识到。 该值可以设置为整数字节数，或者以信息单位（例如“8GB”）。 例如，当该值设置为 4 GB 时，节点会认为它运行在具有 4 GB RAM 的机器上。

Default: undefined (not set or used).

* disk_free_limit

Disk free space limit of the partition on which RabbitMQ is storing data. When available disk space falls below this limit, flow control is triggered. The value can be set relative to the total amount of RAM or as an absolute value in bytes or, alternatively, in information units (e.g `50MB` or `5GB`):  RabbitMQ 存储数据的分区的磁盘可用空间限制。 当可用磁盘空间低于此限制时，将触发流量控制。 该值可以相对于 RAM 的总量进行设置，也可以设置为以字节为单位的绝对值，或者以信息单位（例如“50MB”或“5GB”）为单位：

`disk_free_limit.relative = 3.0`

`disk_free_limit.absolute = 2GB`

By default free disk space must exceed 50MB. See the [Disk Alarms](https://www.rabbitmq.com/disk-alarms.html) documentation.  默认情况下，可用磁盘空间必须超过 50MB。 请参阅磁盘警报文档。

Default: `disk_free_limit.absolute = 50MB`

* log.file.level

Controls the granularity of logging. The value is a list of log event category and log level pairs.

The level can be one of `error` (only errors are logged), `warning` (only errors and warning are logged), `info` (errors, warnings and informational messages are logged), or `debug` (errors, warnings, informational messages and debugging messages are logged).  控制日志记录的粒度。 该值是日志事件类别和日志级别对的列表。级别可以是错误（仅记录错误）、警告（仅记录错误和警告）、信息（记录错误、警告和信息性消息）或调试（记录错误、警告、信息性消息和调试消息）之一 ）。

Default:`log.file.level = info`

* channel_max

Maximum permissible number of channels to negotiate with clients, not including a special channel number 0 used in the protocol. Setting to 0 means "unlimited", a dangerous value since applications sometimes have channel leaks. Using more channels increases memory footprint of the broker.  与客户端协商的最大允许通道数，不包括协议中使用的特殊通道号 0。 设置为 0 意味着“无限制”，这是一个危险的值，因为应用程序有时会出现通道泄漏。 使用更多通道会增加代理的内存占用。

Default: `channel_max = 2047`

* channel_operation_timeout

Channel operation timeout in milliseconds (used internally, not directly exposed to clients due to messaging protocol differences and limitations).  通道操作超时，以毫秒为单位（内部使用，由于消息传递协议的差异和限制，不直接暴露给客户端）。

Default:`channel_operation_timeout = 15000`

* max_message_size

The largest allowed message payload size in bytes. Messages of larger size will be rejected with a suitable channel exception.

Default: 134217728

Max value: 536870912

* heartbeat

Value representing the heartbeat timeout suggested by the server during connection parameter negotiation. If set to 0 on both ends, heartbeats are disabled (this is not recommended). See the [Heartbeats guide](https://www.rabbitmq.com/heartbeats.html) for details.  表示服务器在连接参数协商过程中建议的心跳超时值。 如果两端都设置为 0，则禁用心跳（不建议这样做）。 有关详细信息，请参阅 Heartbeats 指南。

Default:`heartbeat = 60`

* default_vhost

Virtual host to create when RabbitMQ creates a new database from scratch. The exchange `amq.rabbitmq.log` will exist in this virtual host.

Default:`default_vhost = /`

* default_user

User name to create when RabbitMQ creates a new database from scratch.

Default: `default_user = guest`

* default_pass

Password for the default user.

Default: `default_pass = guest`

* default_user_tags

Tags for the default user.

Default: `default_user_tags.administrator = true`

* default_permissions

[Permissions](https://www.rabbitmq.com/access-control.html) to assign to the default user when creating it.

Default:
```
default_permissions.configure = .*
default_permissions.read = .*
default_permissions.write = .*
````

* loopback_users

List of users which are only permitted to connect to the broker via a loopback interface (i.e. `localhost`).  仅允许通过环回接口（即 `localhost`）连接到代理的用户列表。

To allow the default `guest` user to connect remotely (a security practice [unsuitable for production use](https://www.rabbitmq.com/production-checklist.html)), set this to `none`:  要允许默认的 `guest` 用户远程连接（不适合生产使用的安全做法），请将其设置为 `none`：

```
# awful security practice,
# consider creating a new
# user with secure generated credentials!
loopback_users = none
```

To restrict another user to localhost-only connections, do it like so (`monitoring` is the name of the user):  要将另一个用户限制为仅限本地主机的连接，请这样做（`monitoring` 是用户的名称）：

`loopback_users.monitoring = true`

Default:

```
# guest uses well known
# credentials and can only
# log in from localhost
# by default
loopback_users.guest = true
```

* cluster_formation.classic_config.nodes

Classic [peer discovery](https://www.rabbitmq.com/cluster-formation.html) backend's list of nodes to contact. For example, to cluster with nodes `rabbit@hostname1` and `rabbit@hostname2` on first boot:

```
cluster_formation.classic_config.nodes.1 = rabbit@hostname1 cluster_formation.classic_config.nodes.2 = rabbit@hostname2
```

Default: `none` (not set)

* collect_statistics

Statistics collection mode. Primarily relevant for the management plugin. Options are:  统计收集模式。 主要与管理插件相关。 选项是：

`none` (do not emit statistics events)  不发出统计事件

`coarse` (emit per-queue / per-channel / per-connection statistics)  发送每个队列/每个通道/每个连接的统计信息

`fine` (also emit per-message statistics)  也发出每条消息的统计信息

Default:`collect_statistics = none`

* collect_statistics_interval

Statistics collection interval in milliseconds. Primarily relevant for the [management plugin](https://www.rabbitmq.com/management.html#statistics-interval).

Default: `collect_statistics_interval = 5000`

* management_db_cache_multiplier

Affects the amount of time the [management plugin](https://www.rabbitmq.com/management.html#statistics-interval) will cache expensive management queries such as queue listings. The cache will multiply the elapsed time of the last query by this value and cache the result for this amount of time.  影响管理插件缓存昂贵的管理查询（例如队列列表）的时间量。 缓存会将最后一次查询的经过时间乘以该值，并将结果缓存这段时间。

Default: `management_db_cache_multiplier = 5`

* auth_mechanisms

[SASL authentication mechanisms](https://www.rabbitmq.com/authentication.html) to offer to clients.

Default:
```
auth_mechanisms.1 = PLAIN
auth_mechanisms.2 = AMQPLAIN
```

* auth_backends

List of [authentication and authorisation backends](https://www.rabbitmq.com/access-control.html) to use. See the [access control guide](https://www.rabbitmq.com/access-control.html) for details and examples.

Other databases than `rabbit_auth_backend_internal` are available through [plugins](https://www.rabbitmq.com/plugins.html).

Default: `auth_backends.1 = internal`

* reverse_dns_lookups

Set to `true` to have RabbitMQ perform a reverse DNS lookup on client connections, and present that information through `rabbitmqctl` and the management plugin.  设置为 `true` 让 RabbitMQ 在客户端连接上执行反向 DNS 查找，并通过 `rabbitmqctl` 和管理插件显示该信息。

Default: `reverse_dns_lookups = false`

* delegate_count

Number of delegate processes to use for intra-cluster communication. On a machine which has a very large number of cores and is also part of a cluster, you may wish to increase this value.  用于集群内通信的委托进程数。 在具有大量内核并且也是集群一部分的机器上，您可能希望增加此值。

Default: `delegate_count = 16`

* tcp_listen_options

Default socket options. You probably don't want to change this.

Default:
```
tcp_listen_options.backlog = 128
tcp_listen_options.nodelay = true
tcp_listen_options.linger.on = true
tcp_listen_options.linger.timeout = 0
tcp_listen_options.exit_on_close = false
```

* hipe_compile

Do not use. This option is no longer supported. HiPE supported has been dropped starting with Erlang 22.

Default:`hipe_compile = false`

* cluster_partition_handling

How to handle network partitions. Available modes are:

`ignore`

`autoheal`

`pause_minority`

`pause_if_all_down`

pause_if_all_down mode requires additional parameters:

`nodes`

`recover`

See the [documentation on partitions](https://www.rabbitmq.com/partitions.html#automatic-handling) for more information.

Default: `cluster_partition_handling = ignore`



RabbitMQ提供了4种处理网络分区的方式，在rabbitmq.config中配置cluster_partition_handling参数即可，分别为：ignore、pause_minority、pause_if_all_down、autoheal

* ignore: 默认是 ignore, ignore 的配置是当网络分区的时候，RabbitMQ 不会自动做任何处理，即需要手动处理

* pause_minority: 当发生网络分区时，集群中的节点在观察到某些节点 down 掉时，会自动检测其自身是否处于少数派（小于或者等于集群中一半的节点数）。少数派中的节点在分区发生时会自动关闭（类似于执行了 rabbitmqctl stop_app 命令），当分区结束时又会启动。处于关闭的节点会每秒检测一次是否可连通到剩余集群中，如果可以则启动自身的应用，相当于执行 rabbitmqctl start_app 命令。这种处理方式适合集群节点数大于 2 个且最好为奇数的情况。

* pause_if_all_down: 在 pause_if_all_down 模式下，RabbitMQ 会自动关闭不能和 list 中节点通信的节点。语法为 {pause_if_all_down, [nodes], ignore|autoheal}，其中 [nodes] 即为前面所说的 list。如果一个节点与 list 中的所有节点都无法通信时，自关闭其自身。如果 list 中的所有节点都 down 时，其余节点如果是 ok 的话，也会根据这个规则去关闭其自身，此时集群中所有的节点会关闭。如果某节点能够与 list 中的节点恢复通信，那么会启动其自身的 RabbitMQ 应用，慢慢的集群可以恢复。为什么这里会有 ignore 和 autoheal 两种不同的配置，考虑这样一种情况：有两个节点 node1 和 node2 在机架 A 上，node3 和 node4 在机架 B 上，此时机架 A 和机架 B 的通信出现异常，如果此时使用 pause-minority 的话会关闭所有的节点，如果此时采用 pause-if-all-down，list 中配置成 ['node1', 'node3'] 的话，集群中的 4 个节点都不会关闭，但是会形成两个分区，此时就需要 ignore 或者 autoheal 来指引如何处理此种分区的情形。

* autoheal: 在 autoheal 模式下，当认为发生网络分区时，RabbitMQ 会自动决定一个获胜的（winning）分区，然后重启不在这个分区中的节点以恢复网络分区。一个获胜的分区是指客户端连接最多的一个分区。如果产生一个平局，既有两个或者多个分区的客户端连接数一样多，那么节点数最多的一个分区就是获胜的分区。如果此时节点数也一样多，将会以参数输入的顺序来挑选获胜分区。


* cluster_keepalive_interval

How frequently nodes should send keepalive messages to other nodes (in milliseconds). Note that this is not the same thing as [net_ticktime](https://www.rabbitmq.com/nettick.html); missed keepalive messages will not cause nodes to be considered down.  节点应该多久向其他节点发送保活消息（以毫秒为单位）。 请注意，这与 net_ticktime 不同； 错过的 keepalive 消息不会导致节点被视为关闭。

Default: `cluster_keepalive_interval = 10000`

* queue_index_embed_msgs_below

Size in bytes of message below which messages will be embedded directly in the queue index. You are advised to read the [persister tuning](https://www.rabbitmq.com/persistence-conf.html) documentation before changing this.  以字节为单位的消息大小，低于该消息将直接嵌入到队列索引中。 建议您在更改之前阅读持久性调整文档。

Default: `queue_index_embed_msgs_below = 4096`

* mnesia_table_loading_retry_timeout

Timeout used when waiting for Mnesia tables in a cluster to become available.

Default: `mnesia_table_loading_retry_timeout = 30000`

* mnesia_table_loading_retry_limit

Retries when waiting for Mnesia tables in the cluster startup. Note that this setting is not applied to Mnesia upgrades or node deletions.  在集群启动中等待 Mnesia 表时重试。 请注意，此设置不适用于 Mnesia 升级或节点删除。

Default: `mnesia_table_loading_retry_limit = 10`

* mirroring_sync_batch_size

Batch size used to transfer messages to an unsynchronised replica (queue mirror). See [documentation on eager batch synchronization](https://www.rabbitmq.com/ha.html#batch-sync).

Default: `mirroring_sync_batch_size = 4096`

* queue_master_locator

queue leader location strategy. Available strategies are:  队列领导者定位策略。 可用的策略是：

`min-masters`

`client-local`

`random`

See the [documentation on queue leader location](https://www.rabbitmq.com/ha.html#queue-master-location) for more information.

Default: `queue_master_locator = client-local`

* proxy_protocol

If set to true, RabbitMQ will expect a [proxy protocol](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) header to be sent first when an AMQP connection is opened. This implies to set up a proxy protocol-compliant reverse proxy (e.g. [HAproxy](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) or [AWS ELB](http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-proxy-protocol.html)) in front of RabbitMQ. Clients can't directly connect to RabbitMQ when proxy protocol is enabled, so all connections must go through the reverse proxy.  如果设置为 true，RabbitMQ 将期望在打开 AMQP 连接时首先发送代理协议头。 这意味着在 RabbitMQ 前面设置一个符合代理协议的反向代理（例如 HAproxy 或 AWS ELB）。 启用代理协议后，客户端无法直接连接到 RabbitMQ，因此所有连接都必须经过反向代理。

See [the networking guide](https://www.rabbitmq.com/networking.html#proxy-protocol) for more information.

Default: `proxy_protocol = false`

* cluster_name

Operator-controlled cluster name. This name is used to identify a cluster, and by the federation and Shovel plugins to record the origin or path of transferred messages. Can be set to any arbitrary string to help identify the cluster (eg. london). This name can be inspected by AMQP 0-9-1 clients in the server properties map.  操作员控制的集群名称。 此名称用于标识集群，并由 federation 和 Shovel 插件记录传输消息的来源或路径。 可以设置为任意字符串以帮助识别集群（例如 london）。 AMQP 0-9-1 客户端可以在服务器属性映射中检查此名称。

Default: by default the name is derived from the first (seed) node in the cluster.



