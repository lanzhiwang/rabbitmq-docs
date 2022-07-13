# Monitoring

https://www.rabbitmq.com/monitoring.html

## Overview

This document provides an overview of topics related to RabbitMQ monitoring. Monitoring RabbitMQ and applications that use it is critically important. Monitoring helps detect issues before they affect the rest of the environment and, eventually, the end users.  本文档概述了与 RabbitMQ 监控相关的主题。 监控 RabbitMQ 和使用它的应用程序至关重要。 监控有助于在问题影响其他环境并最终影响最终用户之前发现问题。

Many aspects of the system can be monitored. This guide will group them into a handful of categories:  可以监控系统的许多方面。 本指南将它们分为几类：

- What is monitoring, what common approaches to it exist and why it is important.  什么是监控，存在哪些常见的方法以及为什么它很重要。

- Built-in and external monitoring options  内置和外部监控选项

- What infrastructure and kernel metrics are important to monitor  监控哪些基础设施和内核指标很重要

- What RabbitMQ metrics are available:  有哪些 RabbitMQ 指标可用：

    - Node metrics
    - Queue metrics
    - Cluster-wide metrics

- How frequently should monitoring checks be performed?  应多久执行一次监控检查？

- Application-level metrics

- How to approach node health checking and why it's more involved than a single CLI command  如何进行节点健康检查以及为什么它比单个 CLI 命令更复杂

- Health checks use as node readiness probes during deployment or upgrades  运行状况检查在部署或升级期间用作节点准备情况探测

- Log aggregation  日志聚合

- Command-line based observer tool  基于命令行的观察者工具

Log aggregation across all nodes and applications is closely related to monitoring and also mentioned in this guide.  跨所有节点和应用程序的日志聚合与监控密切相关，在本指南中也有提及。

A number of popular tools, both open source and commercial, can be used to monitor RabbitMQ. [Prometheus and Grafana](https://www.rabbitmq.com/prometheus.html) are one highly recommended option.  许多流行的工具，无论是开源的还是商业的，都可用于监控 RabbitMQ。 Prometheus 和 Grafana 是一种强烈推荐的选择。

## What is Monitoring?

In this guide we define monitoring as a process of capturing the behaviour of a system via health checks and metrics over time. This helps detect anomalies: when the system is unavailable, experiences an unusual load, exhausted of certain resources or otherwise does not behave within its normal (expected) parameters. Monitoring involves collecting and storing metrics for the long term, which is important for more than anomaly detection but also root cause analysis, trend detection and capacity planning.  在本指南中，我们将监控定义为通过健康检查和指标随时间捕获系统行为的过程。这有助于检测异常情况：当系统不可用、遇到异常负载、某些资源耗尽或以其他方式不在其正常（预期）参数范围内运行时。监控涉及长期收集和存储指标，这不仅对异常检测很重要，而且对根本原因分析、趋势检测和容量规划也很重要。

Monitoring systems typically integrate with alerting systems. When an anomaly is detected by a monitoring system an alarm of some sort is typically passed to an alerting system, which notifies interested parties such as the technical operations team.  监控系统通常与警报系统集成。当监控系统检测到异常时，通常会将某种警报传递给警报系统，警报系统会通知相关方，例如技术运营团队。

Having monitoring in place means that important deviations in system behavior, from degraded service in some areas to complete unavailability, is easier to detect and the root cause takes much less time to find. Operating a distributed system without monitoring data is a bit like trying to get out of a forest without a GPS navigator device or compass. It doesn't matter how brilliant or experienced the person is, having relevant information is very important for a good outcome.  监控到位意味着系统行为中的重要偏差，从某些区域的服务降级到完全不可用，更容易检测到，查找根本原因所需的时间要少得多。在没有监控数据的情况下运行分布式系统有点像在没有 GPS 导航设备或指南针的情况下试图离开森林。不管这个人有多聪明或有多么有经验，拥有相关信息对于一个好的结果是非常重要的。

### Health Checks' Role in Monitoring

A Health check is the most basic aspect of monitoring. It involves a command or set of commands that collect a few essential metrics of the monitored system over time and test them. For example, whether RabbitMQ's Erlang VM is running is one such check. The metric in this case is "is an OS process running?". The normal operating parameters are "the process must be running". Finally, there is an evaluation step.  健康检查是监控的最基本方面。它涉及一个或一组命令，这些命令或一组命令随着时间的推移收集受监控系统的一些基本指标并对其进行测试。例如，RabbitMQ 的 Erlang VM 是否正在运行就是这样一种检查。这种情况下的指标是“操作系统进程是否正在运行？”。正常运行参数是“进程必须运行”。最后，还有一个评估步骤。

Of course, there are more varieties of health checks. Which ones are most appropriate depends on the definition of a "healthy node" used. So, it is a system- and team-specific decision. [RabbitMQ CLI tools](https://www.rabbitmq.com/cli.html) provide commands that can serve as useful health checks. They will be covered later in this guide.  当然，还有更多种类的健康检查。哪些最合适取决于所使用的“健康节点”的定义。因此，这是一个系统和团队特定的决定。 RabbitMQ CLI 工具提供了可用作有用的健康检查的命令。本指南稍后将介绍它们。

While health checks are a useful tool, they only provide so much insight into the state of the system because they are by design focused on one or a handful of metrics, usually check a single node and can only reason about the state of that node at a particular moment in time. For a more comprehensive assessment, collect more metrics over time. This detects more types of anomalies as some can only be identified over longer periods of time. This is usually done by tools known as monitoring tools of which there are a grand variety. This guides covers some tools used for RabbitMQ monitoring.  虽然健康检查是一个有用的工具，但它们只能提供对系统状态的深入了解，因为它们在设计上专注于一个或少数几个指标，通常检查单个节点并且只能推断该节点的状态特定的时刻。要进行更全面的评估，请随着时间的推移收集更多指标。这可以检测更多类型的异常，因为有些异常只能在更长的时间内识别。这通常通过称为监视工具的工具来完成，这些工具种类繁多。本指南涵盖了一些用于 RabbitMQ 监控的工具。

### System and RabbitMQ Metrics

Some metrics are RabbitMQ-specific: they are collected and reported by RabbitMQ nodes. In this guide we refer to them as "RabbitMQ metrics". Examples include the number of socket descriptors used, total number of enqueued messages or inter-node communication traffic rates. Others metrics are collected and reported by the OS kernel. Such metrics are often called system metrics or infrastructure metrics. System metrics are not specific to RabbitMQ. Examples include CPU utilisation rate, amount of memory used by processes, network packet loss rate, et cetera. Both types are important to track. Individual metrics are not always useful but when analysed together, they can provide a more complete insight into the state of the system. Then operators can form a hypothesis about what's going on and needs addressing.  一些指标是 RabbitMQ 特定的：它们由 RabbitMQ 节点收集和报告。 在本指南中，我们将它们称为“RabbitMQ 指标”。 示例包括使用的套接字描述符的数量、入队消息的总数或节点间通信流量率。 其他指标由操作系统内核收集和报告。 此类指标通常称为系统指标或基础设施指标。 系统指标并非特定于 RabbitMQ。 示例包括 CPU 利用率、进程使用的内存量、网络丢包率等。 这两种类型对于跟踪都很重要。 单个指标并不总是有用，但当一起分析时，它们可以提供对系统状态的更完整的洞察。 然后操作员可以对正在发生的事情和需要解决的问题形成假设。

### Infrastructure and Kernel Metrics

First step towards a useful monitoring system starts with infrastructure and kernel metrics. There are quite a few of them but some are more important than others. Collect the following metrics on all hosts that run RabbitMQ nodes or applications:  迈向有用监控系统的第一步从基础设施和内核指标开始。 其中有很多，但有些比其他更重要。 在运行 RabbitMQ 节点或应用程序的所有主机上收集以下指标：

- CPU stats (user, system, iowait & idle percentages)

- Memory usage (used, buffered, cached & free percentages)

- [Virtual Memory](https://www.kernel.org/doc/Documentation/sysctl/vm.txt) statistics (dirty page flushes, writeback volume)

- Disk I/O (operations & amount of data transferred per unit time, time to service operations)

- Free disk space on the mount used for the [node data directory](https://www.rabbitmq.com/relocate.html)

- File descriptors used by beam.smp vs. [max system limit](https://www.rabbitmq.com/networking.html#open-file-handle-limit)

- TCP connections by state (ESTABLISHED, CLOSE_WAIT, TIME_WAIT)

- Network throughput (bytes received, bytes sent) & maximum network throughput

- Network latency (between all RabbitMQ nodes in a cluster as well as to/from clients)

There is no shortage of existing tools (such as Prometheus or Datadog) that collect infrastructure and kernel metrics, store and visualise them over periods of time.  不乏现有工具（例如 Prometheus 或 Datadog）来收集基础设施和内核指标，并在一段时间内存储和可视化它们。

### Frequency of Monitoring

Many monitoring systems poll their monitored services periodically. How often that's done varies from tool to tool but usually can be configured by the operator.  许多监控系统定期轮询它们的监控服务。完成的频率因工具而异，但通常可由操作员配置。

Very frequent polling can have negative consequences on the system under monitoring. For example, excessive load balancer checks that open a test TCP connection to a node can lead to a [high connection churn](https://www.rabbitmq.com/networking.html#dealing-with-high-connection-churn). Excessive checks of channels and queues in RabbitMQ will increase its CPU consumption. When there are many (say, 10s of thousands) of them on a node, the difference can be significant.  非常频繁的轮询会对受监控的系统产生负面影响。例如，打开与节点的测试 TCP 连接的过多负载均衡器检查会导致高连接流失。 RabbitMQ 中对通道和队列的过多检查会增加其 CPU 消耗。当一个节点上有很多（例如，数千个）它们时，差异可能会很大。

The recommended metric collection interval is 15 second. To collect at an interval which is closer to real-time, use 5 second - but not lower. For rate metrics, use a time range that spans 4 metric collection intervals so that it can tolerate race-conditions and is resilient to scrape failures.  建议的指标收集间隔为 15 秒。要以更接近实时的时间间隔收集，请使用 5 秒 - 但不能更低。对于速率指标，请使用跨越 4 个指标收集间隔的时间范围，以便它可以容忍竞争条件并能够灵活应对抓取失败。

For production systems a collection interval of 30 or even 60 seconds is recommended. [Prometheus](https://www.rabbitmq.com/prometheus.html) exporter API is designed to be scraped every 15 seconds, including production systems.  对于生产系统，建议收集间隔为 30 甚至 60 秒。 Prometheus 导出器 API 旨在每 15 秒抓取一次，包括生产系统。

## Management UI and External Monitoring Systems

RabbitMQ comes with a management UI and HTTP API which exposes a number of RabbitMQ metrics for nodes, connections, queues, message rates and so on. This is a convenient option for development and in environments where external monitoring is difficult or impossible to introduce.  RabbitMQ 带有管理 UI 和 HTTP API，它公开了节点、连接、队列、消息速率等的许多 RabbitMQ 指标。 对于开发和难以或不可能引入外部监控的环境来说，这是一个方便的选择。

However, the management UI has a number of limitations:  但是，管理 UI 有许多限制：

- The monitoring system is intertwined with the system being monitored  监控系统与被监控系统交织在一起

- A certain amount of overhead  一定的开销

- It only stores recent data (think hours, not days or months)  它只存储最近的数据（想想几小时，而不是几天或几个月）

- It has a basic user interface  它有一个基本的用户界面

- Its design [emphasizes ease of use over best possible availability](https://www.rabbitmq.com/management.html#clustering).  它的设计强调易用性而不是最佳可用性。

- Management UI access is controlled via the [RabbitMQ permission tags system](https://www.rabbitmq.com/access-control.html) (or a convention on JWT token scopes)  管理 UI 访问通过 RabbitMQ 权限标签系统（或 JWT 令牌范围的约定）控制

Long term metric storage and visualisation services such as [Prometheus and Grafana](https://www.rabbitmq.com/prometheus.html) or the [ELK stack](https://www.elastic.co/what-is/elk-stack) are more suitable options for production systems. They offer:  Prometheus 和 Grafana 或 ELK 堆栈等长期指标存储和可视化服务更适合生产系统。 他们提供：

- Decoupling of the monitoring system from the system being monitored  监控系统与被监控系统的解耦

- Lower overhead  降低开销

- Long term metric storage  长期指标存储

- Access to additional related metrics such as [Erlang runtime](https://www.rabbitmq.com/runtime.html) ones  访问其他相关指标，例如 Erlang 运行时指标

- More powerful and customizable user interface  更强大和可定制的用户界面

- Ease of metric data sharing: both metric state and dashboards  易于共享指标数据：指标状态和仪表板

- Metric access permissions are not specific to RabbitMQ  度量访问权限并非特定于 RabbitMQ

- Collection and aggregation of node-specific metrics which is more resilient to individual node failures  收集和聚合特定于节点的指标，对单个节点故障更具弹性

RabbitMQ provides first class support for [Prometheus and Grafana](https://www.rabbitmq.com/prometheus.html) as of 3.8. It is recommended for production environments.  从 3.8 开始，RabbitMQ 为 Prometheus 和 Grafana 提供了一流的支持。 推荐用于生产环境。

## RabbitMQ Metrics

The RabbitMQ [management plugin](https://www.rabbitmq.com/management.html) provides an API for accessing RabbitMQ metrics. The plugin will store up to one day's worth of metric data. Longer term monitoring should be accomplished with an external tool.  RabbitMQ 管理插件提供了一个用于访问 RabbitMQ 指标的 API。 该插件最多可存储一天的指标数据。 应使用外部工具完成更长期的监测。

This section will cover multiple RabbitMQ-specific aspects of monitoring.  本节将涵盖监控的多个 RabbitMQ 特定方面。

### Monitoring of Clusters

When monitoring clusters it is important to understand the guarantees provided by the HTTP API. In a clustered environment every node can serve metric endpoint requests. Cluster-wide metrics can be fetched from any node that [can contact its peers](https://www.rabbitmq.com/management.html#clustering). That node will collect and combine data from its peers as needed before producing a response.  在监控集群时，了解 HTTP API 提供的保证很重要。 在集群环境中，每个节点都可以为指标端点请求提供服务。 可以从可以联系其对等方的任何节点获取集群范围的指标。 在产生响应之前，该节点将根据需要从其对等方收集和组合数据。

Every node also can serve requests to endpoints that provide node-specific metrics for itself as well as other cluster nodes. Like with infrastructure and OS metrics, node-specific metrics must be collected for each node. Monitoring tools can execute HTTP API requests against any node.  每个节点还可以向端点提供请求，这些端点为自身以及其他集群节点提供特定于节点的指标。 与基础设施和操作系统指标一样，必须为每个节点收集特定于节点的指标。 监控工具可以针对任何节点执行 HTTP API 请求。

As mentioned earlier, inter-node connectivity issues will [affect HTTP API behaviour](https://www.rabbitmq.com/management.html#clustering). Choose a random online node for monitoring requests. For example, using a load balancer or [round-robin DNS](https://en.wikipedia.org/wiki/Round-robin_DNS).  如前所述，节点间连接问题将影响 HTTP API 行为。 选择一个随机的在线节点来监控请求。 例如，使用负载平衡器或循环 DNS。

Some endpoints perform operations on the target node. Node-local health checks is the most common example. Those are an exception, not the rule.  一些端点在目标节点上执行操作。 节点本地健康检查是最常见的例子。 这些是例外，而不是规则。

### Cluster-wide Metrics

Cluster-wide metrics provide a high level view of cluster state. Some of them describe interaction between nodes. Examples of such metrics are cluster link traffic and detected network partitions. Others combine metrics across all cluster members. A complete list of connections to all nodes would be one example. Both types are complimentary to infrastructure and node metrics.  集群范围的指标提供集群状态的高级视图。 其中一些描述了节点之间的交互。 此类指标的示例是集群链接流量和检测到的网络分区。 其他人将所有集群成员的指标结合起来。 到所有节点的完整连接列表就是一个例子。 这两种类型都是对基础设施和节点指标的补充。

GET /api/overview is the [HTTP API](https://www.rabbitmq.com/management.html#http-api) endpoint that returns cluster-wide metrics.

| Metric                                                       | JSON field name                                              |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| Cluster name                                                 | cluster_name                                                 |
| Cluster-wide message rates                                   | message_stats                                                |
| Total number of connections                                  | object_totals.connections                                    |
| Total number of channels                                     | object_totals.channels                                       |
| Total number of queues                                       | object_totals.queues                                         |
| Total number of consumers                                    | object_totals.consumers                                      |
| Total number of messages (ready plus unacknowledged)         | queue_totals.messages                                        |
| Number of messages ready for delivery                        | queue_totals.messages_ready                                  |
| Number of [unacknowledged](https://www.rabbitmq.com/confirms.html) messages | queue_totals.messages_unacknowledged                         |
| Messages published recently                                  | message_stats.publish                                        |
| Message publish rate                                         | message_stats.publish_details.rate                           |
| Messages delivered to consumers recently                     | message_stats.deliver_get                                    |
| Message delivery rate                                        | message_stats.deliver_get_details.rate                       |
| Other message stats                                          | message_stats.* (see [HTTP API reference](https://rawcdn.githack.com/rabbitmq/rabbitmq-server/v3.9.1/deps/rabbitmq_management/priv/www/api/index.html)) |

### Node Metrics

There are two [HTTP API](https://www.rabbitmq.com/management.html#http-api) endpoints that provide access to node-specific metrics:

- GET /api/nodes/{node} returns stats for a single node

- GET /api/nodes returns stats for all cluster members

The latter endpoint returns an array of objects. Monitoring tools that support (or can support) that as an input should prefer that endpoint since it reduces the number of requests. When that's not the case, use the former endpoint to retrieve stats for every cluster member in turn. That implies that the monitoring system is aware of the list of cluster members.  后一个端点返回一个对象数组。 支持（或可以支持）作为输入的监控工具应该更喜欢该端点，因为它减少了请求的数量。 如果不是这种情况，请使用前一个端点依次检索每个集群成员的统计信息。 这意味着监控系统知道集群成员列表。

Most of the metrics represent point-in-time absolute values. Some, represent activity over a recent period of time (for example, GC runs and bytes reclaimed). The latter metrics are most useful when compared to their previous values and historical mean/percentile values.  大多数指标表示时间点绝对值。 有些代表最近一段时间的活动（例如，GC 运行和字节回收）。 与之前的值和历史平均值/百分位值相比，后一个指标最有用。

| Metric                                                       | JSON field name                   |
| ------------------------------------------------------------ | --------------------------------- |
| Total amount of [memory used](https://www.rabbitmq.com/memory-use.html) | mem_used                          |
| Memory usage high watermark                                  | mem_limit                         |
| Is a [memory alarm](https://www.rabbitmq.com/memory.html) in effect? | mem_alarm                         |
| Free disk space low watermark                                | disk_free_limit                   |
| Is a [disk alarm](https://www.rabbitmq.com/disk-alarms.html) in effect? | disk_free_alarm                   |
| [File descriptors available](https://www.rabbitmq.com/networking.html#open-file-handle-limit) | fd_total                          |
| File descriptors used                                        | fd_used                           |
| File descriptor open attempts                                | io_file_handle_open_attempt_count |
| Sockets available                                            | sockets_total                     |
| Sockets used                                                 | sockets_used                      |
| Message store disk reads                                     | message_stats.disk_reads          |
| Message store disk writes                                    | message_stats.disk_writes         |
| Inter-node communication links                               | cluster_links                     |
| GC runs                                                      | gc_num                            |
| Bytes reclaimed by GC                                        | gc_bytes_reclaimed                |
| Erlang process limit                                         | proc_total                        |
| Erlang processes used                                        | proc_used                         |
| Runtime run queue                                            | run_queue                         |

### Individual Queue Metrics

Individual queue metrics are made available through the [HTTP API](https://www.rabbitmq.com/management.html#http-api) via the GET /api/queues/{vhost}/{qname} endpoint.

| Metric                                                       | JSON field name                                              |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| Memory                                                       | memory                                                       |
| Total number of messages (ready plus unacknowledged)         | messages                                                     |
| Number of messages ready for delivery                        | messages_ready                                               |
| Number of [unacknowledged](https://www.rabbitmq.com/confirms.html) messages | messages_unacknowledged                                      |
| Messages published recently                                  | message_stats.publish                                        |
| Message publishing rate                                      | message_stats.publish_details.rate                           |
| Messages delivered recently                                  | message_stats.deliver_get                                    |
| Message delivery rate                                        | message_stats.deliver_get.rate                               |
| Other message stats                                          | message_stats.* (see [HTTP API reference](https://rawcdn.githack.com/rabbitmq/rabbitmq-management/v3.9.1/priv/www/api/index.html)) |

## Application-level Metrics

A system that uses messaging is almost always distributed. In such systems it is often not immediately obvious which component is misbehaving. Every single part of the system, including applications, should be monitored and investigated.  使用消息传递的系统几乎总是分布式的。在这样的系统中，通常不是很明显哪个组件行为不正常。系统的每个部分，包括应用程序，都应该受到监控和调查。

Some infrastructure-level and RabbitMQ metrics can show presence of an unusual system behaviour or issue but can't pinpoint the root cause. For example, it is easy to tell that a node is running out of disk space but not always easy to tell why. This is where application metrics come in: they can help identify a run-away publisher, a repeatedly failing consumer, a consumer that cannot keep up with the rate, even a downstream service that's experiencing a slowdown (e.g. a missing index in a database used by the consumers).  一些基础设施级别和 RabbitMQ 指标可以显示异常系统行为或问题的存在，但无法查明根本原因。例如，很容易判断一个节点的磁盘空间不足，但并不总是很容易判断原因。这就是应用程序指标的用武之地：它们可以帮助识别失控的发布者、反复失败的消费者、无法跟上速度的消费者，甚至是遇到减速的下游服务（例如，所使用的数据库中缺少索引）由消费者）。

Some client libraries and frameworks provide means of registering metrics collectors or collect metrics out of the box. [RabbitMQ Java client](https://www.rabbitmq.com/api-guide.html) and [Spring AMQP](http://spring.io/projects/spring-amqp) are two examples. With others developers have to track metrics in their application code.  一些客户端库和框架提供注册指标收集器或开箱即用收集指标的方法。 RabbitMQ Java 客户端和 Spring AMQP 就是两个例子。对于其他开发人员，他们必须在他们的应用程序代码中跟踪指标。

What metrics applications track can be system-specific but some are relevant to most systems:  应用程序跟踪的指标可能是特定于系统的，但有些与大多数系统相关：

- Connection opening rate

- Channel opening rate

- Connection failure (recovery) rate

- Publishing rate

- Delivery rate

- Positive delivery acknowledgement rate

- Negative delivery acknowledgement rate

- Mean/95th percentile delivery processing latency

## Health Checks

A health check is a command that tests whether an aspect of the RabbitMQ service is operating as expected. Health checks are executed periodically by machines or interactively by operators.  健康检查是测试 RabbitMQ 服务的某个方面是否按预期运行的命令。 健康检查由机器定期执行或由操作员交互执行。

Health checks can be used to both assess the state and liveness of a node but also as readiness probes by deployment automation and orchestration tools, including during upgrades.  运行状况检查既可用于评估节点的状态和活跃度，也可用作部署自动化和编排工具（包括升级期间）的准备情况探测。

There is a series of health checks that can be performed, starting with the most basic and very rarely producing [false positives](https://en.wikipedia.org/wiki/False_positives_and_false_negatives), to increasingly more comprehensive, intrusive, and opinionated that have a higher probability of false positives. In other words, the more comprehensive a health check is, the less conclusive the result will be.  可以执行一系列健康检查，从最基本且极少产生误报的检查开始，到越来越全面、侵入性和自以为是的误报概率更高的检查。 换句话说，健康检查越全面，结果就越不确定。

Health checks can verify the state of an individual node (node health checks), or the entire cluster (cluster health checks).  健康检查可以验证单个节点（节点健康检查）或整个集群（集群健康检查）的状态。

### Individual Node Checks

This section covers several examples of node health check. They are organised in stages. Higher stages perform more comprehensive and opinionated checks. Such checks will have a higher probability of false positives. Some stages have dedicated RabbitMQ CLI tool commands, others can involve extra tools.  本节介绍节点健康检查的几个示例。 他们是分阶段组织的。 更高的阶段执行更全面和更固执的检查。 这样的检查将有更高的误报概率。 某些阶段具有专用的 RabbitMQ CLI 工具命令，其他阶段可能涉及额外的工具。

While the health checks are ordered, a higher number does not mean a check is "better".  虽然健康检查是有序的，但数字越大并不意味着检查“更好”。

The health checks can be used selectively and combined. Unless noted otherwise, the checks should follow the same monitoring frequency recommendation as metric collection.  健康检查可以有选择地和组合使用。 除非另有说明，否则检查应遵循与指标收集相同的监控频率建议。

Earlier versions of RabbitMQ used an intrusive health check that has since been deprecated and should be avoided. Use one of the checks covered in this section (or their combination).  较早版本的 RabbitMQ 使用了侵入式健康检查，该检查已被弃用，应避免使用。 使用本节中介绍的检查之一（或它们的组合）。

#### Stage 1

The most basic check ensures that the [runtime](https://www.rabbitmq.com/runtime.html) is running and (indirectly) that CLI tools can [authenticate](https://www.rabbitmq.com/cli.html#erlang-cookie) with it.

Except for the CLI tool authentication part, the probability of false positives can be considered approaching 0 except for upgrades and maintenance windows.

[rabbitmq-diagnostics ping](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) performs this check:

```bash
rabbitmq-diagnostics -q ping
# => Ping succeeded if exit code is 0
```

#### Stage 2

A slightly more comprehensive check is executing [rabbitmq-diagnostics status](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) status:

This includes the stage 1 check plus retrieves some essential system information which is useful for other checks and should always be available if RabbitMQ is running on the node (see below).

```bash
rabbitmq-diagnostics -q status
# => [output elided for brevity]
```

This is a common way of sanity checking a node. The probability of false positives can be considered approaching 0 except for upgrades and maintenance windows.

#### Stage 3

Includes previous checks and also verifies that the RabbitMQ application is running (not stopped with [rabbitmqctl stop_app](https://www.rabbitmq.com/rabbitmqctl.8.html#stop_app) or the [Pause Minority partition handling strategy](https://www.rabbitmq.com/partitions.html)) and there are no resource alarms.

```bash
# lists alarms in effect across the cluster, if any
rabbitmq-diagnostics -q alarms
```

[rabbitmq-diagnostics check_running](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) is a check that makes sure that the runtime is running and the RabbitMQ application on it is not stopped or paused.

[rabbitmq-diagnostics check_local_alarms](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) checks that there are no local alarms in effect on the node. If there are any, it will exit with a non-zero status.

The two commands in combination deliver the stage 3 check:

```bash
rabbitmq-diagnostics -q check_running && rabbitmq-diagnostics -q check_local_alarms
# if both checks succeed, the exit code will be 0
```

The probability of false positives is low. Systems hovering around their [high runtime memory watermark](https://www.rabbitmq.com/alarms.html) will have a high probability of false positives. During upgrades and maintenance windows can raise significantly.

Specifically for memory alarms, the GET /api/nodes/{node}/memory HTTP API endpoint can be used for additional checks. In the following example its output is piped to [jq](https://stedolan.github.io/jq/manual/):

```bash
curl --silent -u guest:guest -X GET http://127.0.0.1:15672/api/nodes/rabbit@hostname/memory | jq
# => {
# =>     "memory": {
# =>         "connection_readers": 24100480,
# =>         "connection_writers": 1452000,
# =>         "connection_channels": 3924000,
# =>         "connection_other": 79830276,
# =>         "queue_procs": 17642024,
# =>         "queue_slave_procs": 0,
# =>         "plugins": 63119396,
# =>         "other_proc": 18043684,
# =>         "metrics": 7272108,
# =>         "mgmt_db": 21422904,
# =>         "mnesia": 1650072,
# =>         "other_ets": 5368160,
# =>         "binary": 4933624,
# =>         "msg_index": 31632,
# =>         "code": 24006696,
# =>         "atom": 1172689,
# =>         "other_system": 26788975,
# =>         "allocated_unused": 82315584,
# =>         "reserved_unallocated": 0,
# =>         "strategy": "rss",
# =>         "total": {
# =>             "erlang": 300758720,
# =>             "rss": 342409216,
# =>             "allocated": 383074304
# =>         }
# =>     }
# => }
```

The [breakdown information](https://www.rabbitmq.com/memory-use.html) it produces can be reduced down to a single value using [jq](https://stedolan.github.io/jq/manual/) or similar tools:

```bash
curl --silent -u guest:guest -X GET http://127.0.0.1:15672/api/nodes/rabbit@hostname/memory | jq ".memory.total.allocated"
# => 397365248
```

[rabbitmq-diagnostics -q memory_breakdown](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) provides access to the same per category data and supports various units:

```bash
rabbitmq-diagnostics -q memory_breakdown --unit "MB"
# => connection_other: 50.18 mb (22.1%)
# => allocated_unused: 43.7058 mb (19.25%)
# => other_proc: 26.1082 mb (11.5%)
# => other_system: 26.0714 mb (11.48%)
# => connection_readers: 22.34 mb (9.84%)
# => code: 20.4311 mb (9.0%)
# => queue_procs: 17.687 mb (7.79%)
# => other_ets: 4.3429 mb (1.91%)
# => connection_writers: 4.068 mb (1.79%)
# => connection_channels: 4.012 mb (1.77%)
# => metrics: 3.3802 mb (1.49%)
# => binary: 1.992 mb (0.88%)
# => mnesia: 1.6292 mb (0.72%)
# => atom: 1.0826 mb (0.48%)
# => msg_index: 0.0317 mb (0.01%)
# => plugins: 0.0119 mb (0.01%)
# => queue_slave_procs: 0.0 mb (0.0%)
# => mgmt_db: 0.0 mb (0.0%)
# => reserved_unallocated: 0.0 mb (0.0%)
```

#### Stage 4

Includes all checks in stage 3 plus a check on all enabled listeners (using a temporary TCP connection).

To inspect all listeners enabled on a node, use [rabbitmq-diagnostics listeners](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html):

```bash
rabbitmq-diagnostics -q listeners
# => Interface: [::], port: 25672, protocol: clustering, purpose: inter-node and CLI tool communication
# => Interface: [::], port: 5672, protocol: amqp, purpose: AMQP 0-9-1 and AMQP 1.0
# => Interface: [::], port: 5671, protocol: amqp/ssl, purpose: AMQP 0-9-1 and AMQP 1.0 over TLS
# => Interface: [::], port: 15672, protocol: http, purpose: HTTP API
# => Interface: [::], port: 15671, protocol: https, purpose: HTTP API over TLS (HTTPS)
```

[rabbitmq-diagnostics check_port_connectivity](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) is a command that performs the basic TCP connectivity check mentioned above:

```bash
rabbitmq-diagnostics -q check_port_connectivity
# If the check succeeds, the exit code will be 0
```

The probability of false positives is generally low but during upgrades and maintenance windows can raise significantly.

#### Stage 5

Includes all checks in stage 4 plus checks that there are no failed [virtual hosts](https://www.rabbitmq.com/vhosts.html).

[rabbitmq-diagnostics check_virtual_hosts](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) is a command checks whether any virtual host dependencies may have failed. This is done for all virtual hosts.

```bash
rabbitmq-diagnostics -q check_virtual_hosts
# if the check succeeded, exit code will be 0
```

The probability of false positives is generally low except for systems that are under high CPU load.

### [Health Checks as Readiness Probes](https://www.rabbitmq.com/monitoring.html#readiness-probes)

In some environments, node restarts are controlled with a designated [health check](https://www.rabbitmq.com/monitoring.html#health-checks). The checks verify that one node has started and the deployment process can proceed to the next one. If the check does not pass, the deployment of the node is considered to be incomplete and the deployment process will typically wait and retry for a period of time. One popular example of such environment is Kubernetes where an operator-defined [readiness probe](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-readiness-gate) can prevent a deployment from proceeding when the [OrderedReady pod management policy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#deployment-and-scaling-guarantees) is used.

Given the [peer syncing behavior during node restarts](https://www.rabbitmq.com/clustering.html#restarting-schema-sync), such a health check can prevent a cluster-wide restart from completing in time. Checks that explicitly or implicitly assume a fully booted node that's rejoined its cluster peers will fail and block further node deployments.

Most health check, even relatively basic ones, implicitly assume that the node has finished booting. They are not suitable for nodes that are [awaiting schema table sync](https://www.rabbitmq.com/clustering.html#restarting-schema-sync) from a peer.

One very common example of such check is

```bash
# will exit with an error for the nodes that are currently waiting for
# a peer to sync schema tables from
rabbitmq-diagnostics check_running
```

One health check that does not expect a node to be fully booted and have schema tables synced is

```bash
# a very basic check that will succeed for the nodes that are currently waiting for
# a peer to sync schema from
rabbitmq-diagnostics ping
```

This basic check would allow the deployment to proceed and the nodes to eventually rejoin each other, assuming they are [compatible](https://www.rabbitmq.com/upgrade.html).

#### Optional Check 1

This check verifies that an expected set of plugins is enabled. It is orthogonal to the primary checks.

[rabbitmq-plugins list --enabled](https://www.rabbitmq.com/rabbitmq-plugins.8.html#list) is the command that lists enabled plugins on a node:

```bash
rabbitmq-plugins -q list --enabled --minimal
# => Configured: E = explicitly enabled; e = implicitly enabled
# => | Status: * = running on rabbit@mercurio
# => |/
# => [E*] rabbitmq_auth_mechanism_ssl       3.8.0
# => [E*] rabbitmq_consistent_hash_exchange 3.8.0
# => [E*] rabbitmq_management               3.8.0
# => [E*] rabbitmq_management_agent         3.8.0
# => [E*] rabbitmq_shovel                   3.8.0
# => [E*] rabbitmq_shovel_management        3.8.0
# => [E*] rabbitmq_top                      3.8.0
# => [E*] rabbitmq_tracing                  3.8.0
```

A health check that verifies that a specific plugin, [rabbitmq_shovel](https://www.rabbitmq.com/shovel.html) is enabled and running:

```bash
rabbitmq-plugins -q is_enabled rabbitmq_shovel
# if the check succeeded, exit code will be 0
```

The probability of false positives is generally low but raises in environments where environment variables that can affect [rabbitmq-plugins](https://www.rabbitmq.com/cli.html) are overridden.

## [Deprecated Health Checks and Monitoring Features](https://www.rabbitmq.com/monitoring.html#deprecations)

### Legacy Intrusive Health Check

Earlier versions of RabbitMQ provided a single opinionated and intrusive health check command (and its respective HTTP API endpoint):

```bash
# DO NOT USE: this health check is very intrusive, resource-intensive, prone to false positives
#             and as such, deprecated
rabbitmq-diagnostics node_health_check
```

The above command is **deprecated** will be **removed in a future version** of RabbitMQ and is to be avoided. Systems that use it should adopt one of the [fine grained modern health checks](https://www.rabbitmq.com/monitoring.html#health-checks) instead.

The above check forced every connection, queue leader replica, and channel in the system to emit certain metrics. With a large number of concurrent connections and queues, this can be very resource-intensive and too likely to produce false positives.

The above check is also not suitable to be used as a [readiness probe](https://www.rabbitmq.com/monitoring.html#readiness-probes) as it implicitly assumes a fully booted node.

## [Command-line Based Observer Tool](https://www.rabbitmq.com/monitoring.html#diagnostics-observer)

rabbitmq-diagnostics observer is a command-line tool similar to top, htop, vmstat. It is a command line alternative to [Erlang's Observer application](http://erlang.org/doc/man/observer.html). It provides access to many metrics, including detailed state of individual [runtime](https://www.rabbitmq.com/runtime.html) processes:

- Runtime version information
- CPU and schedule stats
- Memory allocation and usage stats
- Top processes by CPU (reductions) and memory usage
- Network link stats
- Detailed process information such as basic TCP socket stats

and more, in an interactive [ncurses](https://en.wikipedia.org/wiki/Ncurses)-like command line interface with periodic updates.

Here are some screenshots that demonstrate what kind of information the tool provides.

An overview page with key runtime metrics:

[![rabbitmq-diagnostics observer overview](https://www.rabbitmq.com/img/monitoring/observer_cli/diagnostics-observer-overview.png)](https://www.rabbitmq.com/img/monitoring/observer_cli/diagnostics-observer-overview.png)

Memory allocator stats:

[![rabbitmq-diagnostics memory breakdown](https://www.rabbitmq.com/img/monitoring/observer_cli/diagnostics-observer-heap-inspector.png)](https://www.rabbitmq.com/img/monitoring/observer_cli/diagnostics-observer-heap-inspector.png)

A client connection process metrics:

[![rabbitmq-diagnostics connection process](https://www.rabbitmq.com/img/monitoring/observer_cli/diagnostics-observer-connection-process.png)](https://www.rabbitmq.com/img/monitoring/observer_cli/diagnostics-observer-connection-process.png)

## [Monitoring Tools](https://www.rabbitmq.com/monitoring.html#monitoring-tools)

The following is an alphabetised list of third-party tools commonly used to collect RabbitMQ metrics. These tools vary in capabilities but usually can collect both [infrastructure-level](https://www.rabbitmq.com/monitoring.html#system-metrics) and [RabbitMQ metrics](https://www.rabbitmq.com/monitoring.html#rabbitmq-metrics).

Note that this list is by no means complete.

| Monitoring Tool  | Online Resource(s)                                           |
| ---------------- | ------------------------------------------------------------ |
| AppDynamics      | [AppDynamics](https://www.appdynamics.com/community/exchange/extension/rabbitmq-monitoring-extension/), [GitHub](https://github.com/Appdynamics/rabbitmq-monitoring-extension) |
| AWS CloudWatch   | [GitHub](https://github.com/noxdafox/rabbitmq-cloudwatch-exporter) |
| collectd         | [GitHub](https://github.com/signalfx/integrations/tree/master/collectd-rabbitmq) |
| DataDog          | [DataDog RabbitMQ integration](https://docs.datadoghq.com/integrations/rabbitmq/), [GitHub](https://github.com/DataDog/integrations-core/tree/master/rabbitmq) |
| Dynatrace        | [Dynatrace RabbitMQ monitoring](https://www.dynatrace.com/technologies/rabbitmq-monitoring/) |
| Ganglia          | [GitHub](https://github.com/ganglia/gmond_python_modules/tree/master/rabbit) |
| Graphite         | [Tools that work with Graphite](http://graphite.readthedocs.io/en/latest/tools.html) |
| Munin            | [Munin docs](http://munin-monitoring.org/), [GitHub](https://github.com/ask/rabbitmq-munin) |
| Nagios           | [GitHub](https://github.com/nagios-plugins-rabbitmq/nagios-plugins-rabbitmq) |
| Nastel AutoPilot | [Nastel RabbitMQ Solutions](https://www.nastel.com/rabbitmq/) |
| New Relic        | [NewRelic Plugins](https://newrelic.com/plugins/vmware-29/95), [GitHub](https://github.com/pivotalsoftware/newrelic_pivotal_agent) |
| Prometheus       | [Prometheus guide](https://www.rabbitmq.com/prometheus.html), [GitHub](https://github.com/rabbitmq/rabbitmq-prometheus) |
| Sematext         | [Sematext RabbitMQ monitoring integration](https://sematext.com/docs/integration/rabbitmq/), [Sematext RabbitMQ logs integration](https://sematext.com/docs/integration/rabbitmq-logs/) |
| Zabbix           | [Zabbix by HTTP](https://git.zabbix.com/projects/ZBX/repos/zabbix/browse/templates/app/rabbitmq_http), [Zabbix by Agent](https://git.zabbix.com/projects/ZBX/repos/zabbix/browse/templates/app/rabbitmq_agent), [Blog article](http://blog.thomasvandoren.com/monitoring-rabbitmq-queues-with-zabbix.html) |
| Zenoss           | [RabbitMQ ZenPack](https://www.zenoss.com/product/zenpacks/rabbitmq), [Instructional Video](http://www.youtube.com/watch?v=CAak2ayFcV0) |

## [Log Aggregation](https://www.rabbitmq.com/monitoring.html#log-aggregation)

[Logs](https://www.rabbitmq.com/logging.html) are also very important in troubleshooting a distributed system. Like metrics, logs can provide important clues that will help identify the root cause. Collect logs from all RabbitMQ nodes as well as all applications (if possible).

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!


