# Memory and Disk Alarms

https://www.rabbitmq.com/alarms.html

## Overview

During operation, RabbitMQ nodes will consume varying amount of [memory](https://www.rabbitmq.com/memory-use.html) and disk space based on the workload. When usage spikes, both memory and free disk space can reach potentially dangerous levels. In case of memory, the node can be killed by the operating system's low-on-memory process termination mechanism (known as the "OOM killer" on Linux, for example). In case of free disk space, the node can run out of memory, which means it won't be able to perform many internal operations. 在运行期间，RabbitMQ 节点将根据工作负载消耗不同数量的内存和磁盘空间。当使用高峰时，内存和可用磁盘空间都可能达到潜在的危险水平。在内存的情况下，节点可以被操作系统的低内存进程终止机制（例如，在 Linux 上称为“OOM 杀手”）杀死。在可用磁盘空间的情况下，节点可能会耗尽内存，这意味着它将无法执行许多内部操作。

To reduce the likelihood of these scenarios, RabbitMQ has two configurable resource watermarks. When they are reached, RabbitMQ will block connections that publish messages.  为了减少这些场景的可能性，RabbitMQ 有两个可配置的资源水印。当它们到达时，RabbitMQ 将阻止发布消息的连接。

More specifically, RabbitMQ will block connections that publish messages in order to avoid being killed by the OS (out-of-memory killer) or exhausting all available free disk space:  更具体地说，RabbitMQ 将阻止发布消息的连接，以避免被操作系统杀死（内存不足杀手）或耗尽所有可用的可用磁盘空间：

- When [memory use](https://www.rabbitmq.com/memory-use.html) goes above the configured watermark (limit)  当内存使用超过配置的水印（限制）时

- When [free disk space](https://www.rabbitmq.com/disk-alarms.html) drops below the configured watermark (limit)  当可用磁盘空间低于配置的水印（限制）时

Nodes will temporarily *block* publishing connections by suspending reading from [client connection](https://www.rabbitmq.com/connections.html). Connections that are only used to *consume* messages will not be blocked.  节点将通过暂停从客户端连接读取来临时阻止发布连接。仅用于消费消息的连接不会被阻塞。

Connection [heartbeat monitoring](https://www.rabbitmq.com/heartbeats.html) will be disabled, too. All network connections will show in rabbitmqctl and the management UI as either blocking, meaning they have not attempted to publish and can thus continue, or blocked, meaning they have published and are now paused. Compatible clients will be notified when they are blocked.  连接心跳监控也将被禁用。所有网络连接都将在 rabbitmqctl 和管理 UI 中显示为阻塞，这意味着它们没有尝试发布，因此可以继续，或者阻塞，意味着它们已经发布并且现在被暂停。兼容的客户端将在被阻止时收到通知。

Connections that only consume are not blocked by resource alarms; deliveries to them continue as usual.  只消费的连接不会被资源告警阻塞；像往常一样继续向他们发货。

## Client Notifications  客户通知

Modern client libraries support [connection.blocked notification](https://www.rabbitmq.com/connection-blocked.html) (a protocol extension), so applications can monitor when they are blocked.  现代客户端库支持 connection.blocked 通知（协议扩展），因此应用程序可以监控它们何时被阻止。

## Alarms in Clusters

When running RabbitMQ in a cluster, the memory and disk alarms are cluster-wide; if one node goes over the limit then all nodes will block connections.  在集群中运行 RabbitMQ 时，内存和磁盘告警是集群范围的；如果一个节点超过限制，那么所有节点都会阻塞连接。

The intent here is to stop producers but let consumers continue unaffected. However, since the protocol permits producers and consumers to operate on the same channel, and on different channels of a single connection, this logic is necessarily imperfect. In practice that does not pose any problems for most applications since the throttling is observable merely as a delay. Nevertheless, other design considerations permitting, it is advisable to only use individual connections for either producing or consuming.  这里的目的是阻止生产者，但让消费者继续不受影响。但是，由于协议允许生产者和消费者在同一个通道上操作，而在单个连接的不同通道上操作，这个逻辑必然是不完善的。实际上，这对大多数应用程序不会造成任何问题，因为节流仅作为延迟可观察到。然而，在其他设计考虑允许的情况下，建议仅使用单独的连接进行生产或消费。

## Effects on Data Safety  对数据安全的影响

When an alarm is in effect, publishing connections will be blocked by TCP back pressure. In practice this means that publish operations will eventually time out of fail outright. Application developers must be prepared to handle such failures and use [publisher confirms](https://www.rabbitmq.com/confirms.html) to keep track of what messages have been successfully handled and processed by RabbitMQ.  当警报生效时，发布连接将被 TCP 背压阻塞。实际上，这意味着发布操作最终会完全超时。应用程序开发人员必须准备好处理此类故障，并使用发布者确认来跟踪 RabbitMQ 已成功处理和处理了哪些消息。

## Running Out of File Descriptors

When the server is close to using all the file descriptors that the OS has made available to it, it will refuse client connections. See [Networking guide](https://www.rabbitmq.com/networking.html) to learn more.  当服务器接近使用操作系统提供给它的所有文件描述符时，它将拒绝客户端连接。请参阅网络指南以了解更多信息。

## Transient Flow Control

When clients attempt to publish faster than the server can accept their messages, they go into transient [flow control](https://www.rabbitmq.com/flow-control.html).  当客户端尝试发布的速度快于服务器接受其消息的速度时，它们会进入瞬态流控制。

## Relevant Topics

- [Determining what uses memory](https://www.rabbitmq.com/memory-use.html) on a running node  确定正在运行的节点上使用内存的内容

- [Memory alarms](https://www.rabbitmq.com/memory.html)  内存警报

- [Free disk space alarms](https://www.rabbitmq.com/disk-alarms.html)  可用磁盘空间警报

- [How clients can determine if they are blocked](https://www.rabbitmq.com/connection-blocked.html)  客户如何确定他们是否被阻止

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!



