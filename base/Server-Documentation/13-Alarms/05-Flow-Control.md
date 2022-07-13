# Flow Control

https://www.rabbitmq.com/flow-control.html

## Overview

This guide covers a back pressure mechanism applied by RabbitMQ nodes to publishing connections in order to avoid runaway [memory usage](https://www.rabbitmq.com/memory-use.html) growth. It is necessary because some components in a node can fall behind particularly fast publishers as they have to do significantly more work than publishing clients (e.g. replicate data to N peer nodes or store it on disk).  本指南介绍了 RabbitMQ 节点用于发布连接的背压机制，以避免内存使用量增长失控。这是必要的，因为节点中的某些组件可能会落后于特别快的发布者，因为它们必须比发布客户端做更多的工作（例如，将数据复制到 N 个对等节点或将其存储在磁盘上）。

## How Does Flow Control Work

RabbitMQ will reduce the speed of connections which are publishing too quickly for queues to keep up. No configuration is required.  RabbitMQ 将降低发布太快而队列无法跟上的连接速度。无需配置。

A flow-controlled connection will show a state of flow in rabbitmqctl, management UI and HTTP API responses. This means the connection is experiencing blocking and unblocking several times a second, in order to keep the rate of message ingress at one that the rest of the server (e.g. queues those messages are route to) can handle.  流控制的连接将在 rabbitmqctl、管理 UI 和 HTTP API 响应中显示流状态。这意味着连接每秒经历多次阻塞和解除阻塞，以便将消息进入速率保持在服务器其余部分（例如，这些消息路由到的队列）可以处理的速度。

In general, a connection which is in flow control should not see any difference from normal running; the flow state is there to inform the sysadmin that the publishing rate is restricted, but from the client's perspective it should just look like the network bandwidth to the server is lower than it actually is.  一般来说，处于流量控制状态的连接与正常运行应该没有任何区别；流状态通知系统管理员发布速率受到限制，但从客户端的角度来看，它应该只是看起来服务器的网络带宽低于实际值。

Other components than connections can be in the flow state. Channels, queues and other parts of the system can apply flow control that eventually propagates back to publishing connections.  连接以外的其他组件可以处于流动状态。通道、队列和系统的其他部分可以应用最终传播回发布连接的流控制。

To find out if consumers and [prefetch settings](https://www.rabbitmq.com/confirms.html) can be key limiting factors, [take a look at relevant metrics](https://blog.rabbitmq.com/posts/2014/04/finding-bottlenecks-with-rabbitmq-3-3/). See [Monitoring and Health Checks](https://www.rabbitmq.com/monitoring.html) guide to learn more.  要了解消费者和预取设置是否可能是关键限制因素，请查看相关指标。请参阅监控和健康检查指南了解更多信息。

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!


