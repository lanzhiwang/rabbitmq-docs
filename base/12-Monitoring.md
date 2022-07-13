# Monitoring with Prometheus & Grafana

https://www.rabbitmq.com/prometheus.html

## Overview

This guide covers RabbitMQ monitoring with two popular tools: [Prometheus](https://prometheus.io/docs/introduction/overview/), a monitoring toolkit; and [Grafana](https://grafana.com/grafana), a metrics visualisation system.  本指南涵盖了使用两个流行工具进行 RabbitMQ 监控：Prometheus，一个监控工具包； 和 Grafana，一个度量可视化系统。

These tools together form a powerful toolkit for long-term metric collection and monitoring of RabbitMQ clusters. While [RabbitMQ management UI](https://www.rabbitmq.com/management.html) also provides access to a subset of metrics, it by design doesn't try to be a long term metric collection solution.  这些工具共同构成了一个强大的工具包，用于 RabbitMQ 集群的长期指标收集和监控。 虽然 RabbitMQ 管理 UI 还提供对指标子集的访问，但它的设计并不试图成为长期的指标收集解决方案。

Please read through the main [guide on monitoring](https://www.rabbitmq.com/monitoring.html) first. Monitoring principles and available metrics are mostly relevant when Prometheus and Grafana are used.  请先通读有关监控的主要指南。 当使用 Prometheus 和 Grafana 时，监控原则和可用指标最相关。

Some key topics covered by this guide are  本指南涵盖的一些关键主题是

- Prometheus support basics

- Grafana support basics

- Quick Start for local experimentation

- Installation steps for production systems

- Two types of scraping endpoint responses: Aggregated vs. Individual Entity Metrics

Grafana dashboards follow a number of conventions to make the system more observable and anti-patterns easier to spot. Its design decisions are explained in a number of sections:  Grafana 仪表板遵循许多约定，以使系统更易于观察，并且更容易发现反模式。 它的设计决策在许多部分中进行了解释：

- RabbitMQ Overview Dashboard

- Health indicators on the Overview dashboard

- Graph colour labelling conventions

- Graph thresholds

- Relevant documentation for each graph (metric)

- Spotting Anti-patterns

- Other available dashboards

- TLS support for Prometheus scraping endpoint

### Built-in Prometheus Support

As of 3.8.0, RabbitMQ ships with built-in Prometheus & Grafana support.  从 3.8.0 开始，RabbitMQ 附带内置 Prometheus 和 Grafana 支持。

Support for Prometheus metric collector ships in the rabbitmq_prometheus plugin. The plugin exposes all RabbitMQ metrics on a dedicated TCP port, in Prometheus text format.  rabbitmq_prometheus 插件中提供了对 Prometheus 度量收集器的支持。 该插件以 Prometheus 文本格式在专用 TCP 端口上公开所有 RabbitMQ 指标。

These metrics provide deep insights into the state of RabbitMQ nodes and [the runtime](https://www.rabbitmq.com/runtime.html). They make reasoning about the behaviour of RabbitMQ, applications that use it and various infrastructure elements a lot more informed.  这些指标可以深入了解 RabbitMQ 节点和运行时的状态。 他们对 RabbitMQ 的行为、使用它的应用程序和各种基础设施元素进行了更深入的推理。

### Grafana Support

Collected metrics are not very useful unless they are visualised. Team RabbitMQ provides a [prebuilt set of Grafana dashboards](https://grafana.com/rabbitmq) that visualise a large number of available RabbitMQ and [runtime](https://www.rabbitmq.com/runtime.html) metrics in context-specific ways.

There is a number of dashboards available:

- an overview dashboard

- runtime [memory allocators](https://www.rabbitmq.com/runtime.html#allocators) dashboard

- an [inter-node communication](https://www.rabbitmq.com/clustering.html#cluster-membership) (Erlang distribution) dashboard

- a Raft metric dashboard

and others. Each is meant to provide an insight into a specific part of the system. When used together, they are able to explain RabbitMQ and application behaviour in detail.

Note that the Grafana dashboards are opinionated and use a number of conventions, for example, to spot system health issues quicker or make cross-graph referencing possible. Like all Grafana dashboards, they are also highly customizable. The conventions they assume are considered to be good practices and are thus recommended.

### An Example

When RabbitMQ is integrated with Prometheus and Grafana, this is what the RabbitMQ Overview dashboard looks like:

![](../images/rabbitmq-overview-dashboard.png)

## Quick Start

### Before We Start

This section explains how to set up a RabbitMQ cluster with Prometheus and Grafana dashboards, as well as some applications that will produce some activity and meaningful metrics.  本节介绍如何使用 Prometheus 和 Grafana 仪表板设置 RabbitMQ 集群，以及一些将产生一些活动和有意义的指标的应用程序。

With this setup you will be able to interact with RabbitMQ, Prometheus & Grafana running locally. You will also be able to try out different load profiles to see how it all fits together, make sense of the dashboards, panels and so on.  通过此设置，您将能够与本地运行的 RabbitMQ、Prometheus 和 Grafana 进行交互。 您还可以尝试不同的负载配置文件，看看它们是如何组合在一起的，了解仪表板、面板等。

This is merely an example; the rabbitmq_prometheus plugin and our Grafana dashboards do not require the use of Docker Compose demonstrated below.  这只是一个例子； rabbitmq_prometheus 插件和我们的 Grafana 仪表板不需要使用下面演示的 Docker Compose。

### Prerequisites

The instructions below assume a host machine that has a certain set of tools installed:

- A terminal to run the commands

- [Git](https://git-scm.com/) to clone the repository

- [Docker Desktop](https://www.docker.com/products/docker-desktop) to use Docker Compose locally

- A Web browser to browse the dashboards

Their installation is out of scope of this guide. Use

```bash
git version
docker info && docker-compose version
```

on the command line to verify that the necessary tools are available.

### Clone a Repository with Manifests

First step is to clone a Git repository, [rabbitmq-server](https://github.com/rabbitmq/rabbitmq-server), with the manifests and other components required to run a RabbitMQ cluster, Prometheus and a set of applications:

```bash
git clone https://github.com/rabbitmq/rabbitmq-server.git
cd rabbitmq-server/deps/rabbitmq_prometheus/docker
```

### Run Docker Compose

Next use Docker Compose manifests to run a pre-configured RabbitMQ cluster, a Prometheus instance and a basic workload that will produce the metrics displayed in the RabbitMQ overview dashboard:

```bash
docker-compose -f docker-compose-metrics.yml up -d
docker-compose -f docker-compose-overview.yml up -d
```

The docker-compose commands above can also be executed with a make target:

```bash
make metrics overview
```

When the above commands succeed, there will be a functional RabbitMQ cluster and a Prometheus instance collecting metrics from it running in a set of containers.

### Access RabbitMQ Overview Grafana Dashboard

Now navigate to http://localhost:3000/dashboards in a Web browser. It will bring up a login page. Use admin for both the username and the password. On the very first login Grafana will suggest changing your password. For the sake of this example, we suggest that this step is skipped.

Navigate to the **RabbitMQ-Overview** dashboard that will look like this:

![RabbitMQ Overview Dashboard Localhost](https://www.rabbitmq.com/img/rabbitmq-overview-dashboard-localhost.png)

Congratulations! You now have a 3-nodes RabbitMQ cluster integrated with Prometheus & Grafana running locally. This is a perfect time to learn more about the available dashboards.

## [RabbitMQ Overview Dashboard](https://www.rabbitmq.com/prometheus.html#rabbitmq-overview-dashboard)

All metrics available in the [management UI](https://www.rabbitmq.com/management.html) Overview page are available in the Overview Grafana dashboard. They are grouped by object type, with a focus on RabbitMQ nodes and message rates.

### [Health Indicators](https://www.rabbitmq.com/prometheus.html#health-indicators)

[Single stat metrics](https://grafana.com/docs/features/panels/singlestat/) at the top of the dashboard capture the health of a single RabbitMQ cluster. In this case, there's a single RabbitMQ cluster, **rabbitmq-overview**, as seen in the **Cluster** drop-down menu just below the dashboard title.

The panels on all RabbitMQ Grafana dashboards use different colours to capture the following metric states:

- **Green** means the value of the metric is within a healthy range
- **Blue** means under-utilisation or some form of degradation
- **Red** means the value of the metric is below or above the range that is considered healthy

![RabbitMQ Overview Dashboard Single-stat](https://www.rabbitmq.com/img/rabbitmq-overview-dashboard-single-stat.png)

Default ranges for the [single stat metrics](https://grafana.com/docs/features/panels/singlestat/) **will not be optimal for all** RabbitMQ deployments. For example, in environments with many consumers and/or high prefetch values, it may be perfectly fine to have over 1,000 unacknowledged messages. The default thresholds can be easily adjusted to suit the workload and system at hand. The users are encouraged to revisit these ranges and tweak them as they see fit for their workloads, monitoring and operational practices, and tolerance for false positives.

### [Metrics and Graphs](https://www.rabbitmq.com/prometheus.html#graphs)

Most RabbitMQ and runtime metrics are represented as graphs in Grafana: they are values that change over time. This is the simplest and clearest way of visualising how some aspect of the system changes. Time-based charting makes it easy to understand the change in key metrics: message rates, memory used by every node in the cluster, or the number of concurrent connections. All metrics except for health indicators are **node-specific**, that is, they represent values of a metric on a single node.

Some metrics, such as the panels grouped under **CONNECTIONS**, are stacked to capture the state of the **cluster as a whole**. These metrics are collected on individual nodes and grouped visually, which makes it easy to notice when, for example, one node serves a disproportionate number of connections.

We would refer to such a RabbitMQ cluster as **unbalanced**, meaning that at least in some ways, a minority of nodes perform the majority of work.

In the example below, connections are spread out evenly across all nodes most of the time:

![RabbitMQ Overview Dashboard CONNECTIONS](https://www.rabbitmq.com/img/rabbitmq-overview-dashboard-connections.png)

### [Colour Labelling in Graphs](https://www.rabbitmq.com/prometheus.html#graph-colour-labelling)

All metrics on all graphs are associated with specific node names. For example, all metrics drawn in green are for the node that contains 0 in its name, e.g. rabbit@rmq0. This makes it easy to correlate metrics of a specific node across graphs. Metrics for the first node, which is assumed to contain 0 in its name, will always appear as green across all graphs.

It is important to remember this aspect when using the RabbitMQ Overview dashboard. **If a different node naming convention is used, the colours will appear inconsistent across graphs**: green may represent e.g. rabbit@foo in one graph, and e.g. rabbit@bar in another graph.

When this is the case, the panels must be updated to use a different node naming scheme.

### [Thresholds in Graphs](https://www.rabbitmq.com/prometheus.html#graph-thresholds)

Most metrics have pre-configured thresholds. They define expected operating boundaries for the metric. On the graphs they appear as semi-transparent orange or red areas, as seen in the example below.

![RabbitMQ Overview Dashboard Single-stat](https://www.rabbitmq.com/img/rabbitmq-overview-dashboard-memory-threshold.png)

Metric values in the **orange** area signal that some pre-defined threshold has been exceeded. This may be acceptable, especially if the metric recovers. A metric that comes close to the orange area is considered to be in healthy state.

Metric values in the **red** area need attention and may identify some form of service degradation. For example, metrics in the red area can indicate that an [alarm](https://www.rabbitmq.com/alarms.html) in effect or when the node is [out of file descriptors](https://www.rabbitmq.com/networking.html) and cannot accept any more connections or open new files.

In the example above, we have a RabbitMQ cluster that runs at optimal memory capacity, which is just above the warning threshold. If there is a spike in published messages that should be stored in RAM, the amount of memory used by the node go up and the metric on the graph will go down (as it indicates the amount of **available** memory).

Because the system has more memory available than is allocated to the RabbitMQ node it hosts, we notice the dip below **0 B**. This emphasizes the importance of leaving spare memory available for the OS, housekeeping tasks that cause short-lived memory usage spikes, and other processes. When a RabbitMQ node exhausts all memory that it is allowed to use, a [memory alarm](https://www.rabbitmq.com/alarms.html) is triggered and publishers across the entire cluster will be blocked.

On the right side of the graph we can see that consumers catch up and the amount of memory used goes down. That clears the memory alarm on the node and, as a result, publishers become unblocked. This and related metrics across cluster then should return to their optimal state.

### [There is No "Right" Threshold for Many Metrics](https://www.rabbitmq.com/prometheus.html#graph-thresholds-are-system-specific)

Note that the thresholds used by the Grafana dashboards have to have a default value. No matter what value is picked by dashboard developers, they **will not be suitable for all environments and workloads**.

Some workloads may require higher thresholds, others may choose to lower them. While the defaults should be adequate in many cases, the **operator must review and adjust the thresholds** to suit their specific requirements.

### [Relevant Documentation for Graphs](https://www.rabbitmq.com/prometheus.html#graph-documentation)

Most metrics have a help icon in the top-left corner of the panel.

![RabbitMQ Overview Dashboard Single-stat](https://www.rabbitmq.com/img/rabbitmq-overview-dashboard-disk-documentation.png)

Some, like the available disk space metric, link to dedicated pages in [RabbitMQ documentation](https://www.rabbitmq.com/documentation.html). These pages contain information relevant to the metric. Getting familiar with the linked guides is highly recommended and will help the operator understand what the metric means better.

### [Using Graphs to Spot Anti-patterns](https://www.rabbitmq.com/prometheus.html#spot-anti-pattern)

Any metric drawn in red hints at an anti-pattern in the system. Such graphs try to highlight sub-optimal uses of RabbitMQ. A **red graphs with non-zero metrics should be investigated**. Such metrics might indicate an issue in RabbitMQ configuration or sub-optimal actions by clients ([publishers](https://www.rabbitmq.com/publishers.html) or [consumers](https://www.rabbitmq.com/consumers.html)).

In the example below we can see the usage of greatly inefficient [polling consumers](https://www.rabbitmq.com/consumers.html#fetching) that keep polling, even though most or even all polling operation return no messages. Like any polling-based algorithm, it is wasteful and should be avoided where possible.

It is a lot more and efficient to have RabbitMQ [push messages to the consumer](https://www.rabbitmq.com/consumers.html).

![RabbitMQ Overview Dashboard Antipatterns](https://www.rabbitmq.com/img/rabbitmq-overview-dashboard-antipattern.png)

### [Example Workloads](https://www.rabbitmq.com/prometheus.html#example-workloads)

The [Prometheus plugin repository](https://github.com/rabbitmq/rabbitmq-server/tree/master/deps/rabbitmq_prometheus/docker) contains example workloads that use [PerfTest](https://rabbitmq.github.io/rabbitmq-perf-test/stable/htmlsingle/) to simulate different workloads. Their goal is to exercise all metrics in the RabbitMQ Overview dashboard. These examples are meant to be edited and extended as developers and operators see fit when exploring various metrics, their thresholds and behaviour.

To deploy a workload app, run docker-compose -f docker-compose-overview.yml up -d. The same command will redeploy the app after the file has been updated.

To delete all workload containers, run docker-compose -f docker-compose-overview.yml down or

```bash
gmake down
```

## [More Dashboards: Raft and Erlang Runtime](https://www.rabbitmq.com/prometheus.html#other-dashboards)

There are two more Grafana dashboards available: **RabbitMQ-Raft** and **Erlang-Distribution**. They collect and visualise metrics related to the Raft consensus algorithm (used by [Quorum Queues](https://www.rabbitmq.com/quorum-queues.html) and other features) as well as more nitty-gritty [runtime metrics](https://www.rabbitmq.com/runtime.html) such as inter-node communication buffers.

The dashboards have corresponding RabbitMQ clusters and PerfTest instances which are started and stopped the same as the Overview one. Feel free to experiment with the other workloads that are included in [the same docker directory](https://github.com/rabbitmq/rabbitmq-server/tree/master/deps/rabbitmq_prometheus/docker).

For example, the docker-compose-dist-tls.yml Compose manifest is meant to stress the [inter-node communication links](https://www.rabbitmq.com/clustering.html). This workload uses a lot of system resources. docker-compose-qq.yml contains a quorum queue workload.

To stop and delete all containers used by the workloads, run docker-compose -f [file] down or

```bash
make down
```

## [Installation](https://www.rabbitmq.com/prometheus.html#installation)

Unlike the [Quick Start](https://www.rabbitmq.com/prometheus.html#quick-start) above, this section covers monitoring setup geared towards production usage.

We will assume that the following tools are provisioned and running:

- A [3-node RabbitMQ 3.8 cluster](https://www.rabbitmq.com/cluster-formation.html)
- Prometheus, including network connectivity with all RabbitMQ cluster nodes
- Grafana, including configuration that lists the above Prometheus instance as one of the data sources

### [RabbitMQ Configuration](https://www.rabbitmq.com/prometheus.html#rabbitmq-configuration)

#### Cluster Name

First step is to give the RabbitMQ cluster a descriptive name so that it can be distinguished from other clusters.

To find the current name of the cluster, use

```bash
rabbitmq-diagnostics -q cluster_status
```

This command can be executed against any cluster node. If the current cluster name is distinctive and appropriate, skip the rest of this paragraph. To change the name of the cluster, run the following command (the name used here is just an example):

```bash
rabbitmqctl -q set_cluster_name testing-prometheus
```

### Enable rabbitmq_prometheus

Next, enable the **rabbitmq_prometheus** [plugin](https://www.rabbitmq.com/plugins.html) on all nodes:

```bash
rabbitmq-plugins enable rabbitmq_prometheus
```

The output will look similar to this:

```ini
rabbitmq-plugins enable rabbitmq_prometheus

Enabling plugins on node rabbit@ed9618ea17c9:
rabbitmq_prometheus
The following plugins have been configured:
  rabbitmq_management_agent
  rabbitmq_prometheus
  rabbitmq_web_dispatch
Applying plugin configuration to rabbit@ed9618ea17c9...
The following plugins have been enabled:
  rabbitmq_management_agent
  rabbitmq_prometheus
  rabbitmq_web_dispatch

started 3 plugins.
```

To confirm that RabbitMQ now exposes metrics in Prometheus format, get the first couple of lines with curl or similar:

```bash
curl -s localhost:15692/metrics | head -n 3
# TYPE erlang_mnesia_held_locks gauge
# HELP erlang_mnesia_held_locks Number of held locks.
erlang_mnesia_held_locks{node="rabbit@65f1a10aaffa",cluster="rabbit@65f1a10aaffa"} 0
```

Notice that RabbitMQ exposes the metrics on a dedicated TCP port, 15692 by default.

### [Prometheus Configuration](https://www.rabbitmq.com/prometheus.html#prometheus)

Once RabbitMQ is configured to expose metrics to Prometheus, Prometheus should be made aware of where it should scrape RabbitMQ metrics from. There are a number of ways of doing this. Please refer to the official [Prometheus configuration documentation](https://prometheus.io/docs/prometheus/latest/configuration/configuration/). There's also a [first steps with Prometheus](https://prometheus.io/docs/introduction/first_steps/) guide for beginners.

#### Metric Collection and Scraping Intervals

Prometheus will periodically scrape (read) metrics from the systems it monitors, every 60 seconds by default. RabbitMQ metrics are updated periodically, too, every 5 seconds by default. Since [this value is configurable](https://www.rabbitmq.com/management.html#statistics-interval), check the metrics update interval by running the following command on any of the nodes:

```bash
rabbitmq-diagnostics environment | grep collect_statistics_interval
# => {collect_statistics_interval,5000}
```

The returned value will be **in milliseconds**.

For production systems, we recommend a minimum value of 15s for Prometheus scrape interval and a 10000 (10s) value for RabbitMQ's collect_statistics_interval. With these values, Prometheus doesn't scrape RabbitMQ too frequently, and RabbitMQ doesn't update metrics unnecessarily. If you configure a different value for Prometheus scrape interval, remember to set an appropriate interval when visualising metrics in Grafana with rate() - [4x the scrape interval is considered safe](https://www.robustperception.io/what-range-should-i-use-with-rate).

When using RabbitMQ's Management UI default 5 second auto-refresh, keeping the default collect_statistics_interval setting is optimal. Both intervals are 5000 ms (5 seconds) by default for this reason.

To confirm that Prometheus is scraping RabbitMQ metrics from all nodes, ensure that all RabbitMQ endpoints are **Up** on the Prometheus Targets page, as shown below:

![Prometheus RabbitMQ Targets](https://www.rabbitmq.com/img/monitoring/prometheus/prometheus-targets.png)

### [Network Interface and Port](https://www.rabbitmq.com/prometheus.html#listener)

The port is configured using the prometheus.tcp.port key:

```ini
prometheus.tcp.port = 15692
```

It is possible to configure what interface the Prometheus plugin API endpoint will use, similarly to [messaging protocol listeners](https://www.rabbitmq.com/networking.html#interfaces), using the prometheus.tcp.ip key:

```ini
prometheus.tcp.ip = 0.0.0.0
```

To check what interface and port is used by a running node, use rabbitmq-diagnostics:

```bash
rabbitmq-diagnostics -s listeners
# => Interface: [::], port: 15692, protocol: http/prometheus, purpose: Prometheus exporter API over HTTP
```

or [tools such as lsof, ss or netstat](https://www.rabbitmq.com/troubleshooting-networking.html#ports).

### [Aggregated and Per-Object Metrics](https://www.rabbitmq.com/prometheus.html#metric-aggregation)

The scraping HTTP endpoint can produce metrics as aggregated rows or individual rows.

By default, returned rows are aggregated by metric name. This significantly reduces the size of the output and makes it constant even as the number of objects (e.g. connections and queues) grows.

In the per-object mode, there will be an output row for **each object-metric pair**. With a large number of stats-emitting entities, e.g. a lot of connections and queues, this can result in very large payloads and a lot of CPU resources spent serialising data to output.

Metric aggregation is a more predictable and practical option for larger deployments. It scales very well with respect to the number of metric-emitting objects in the system (connections, channels, queues, consumers, etc) by keeping response size and time small. It is also predictably easy to visualise.

The downside of metric aggregation is that it loses data fidelity. Per-object metrics and alerting are not possible with aggregation. Individual object metrics, while very useful in some cases, are also hard to visualise. Consider what a chart with 200K connections charted on it would look like and whether an operator would be able to make sense of it.

To enable per-object (unaggregated) metrics, use the prometheus.return_per_object_metrics key:

```ini
# can result in a really excessive output produced,
# only suitable for environments with a relatively small
# number of metrics-emitting objects such as connections and queues
prometheus.return_per_object_metrics = true
```

For the sake of completeness, the default used by the plugin is

```ini
# Enables metric aggregation. Individual object (e.g. connection or queue) metrics
# will not be emitted to significantly reduce output size.
# This option is recommended for most environments.
prometheus.return_per_object_metrics = false
```

### [Scraping Endpoint Timeouts](https://www.rabbitmq.com/prometheus.html#timeout-configuration)

In some environments there aren't many stats-emitting entities (queues, connections, channels), in others the scraping HTTP endpoint may have to return a sizeable data set to the client (e.g. many thousands of rows). In those cases the amount of time it takes to process the request can exceed certain timeouts in the embedded HTTP server and the HTTP client used by Prometheus.

It is possible to bump plugin side HTTP request timeouts using the prometheus.tcp.idle_timeout, prometheus.tcp.inactivity_timeout, prometheus.tcp.request_timeout settings.

- prometheus.tcp.inactivity_timeout controls HTTP(S) client's TCP connection inactivity timeout. When it is reached, the connection will be closed by the HTTP server.
- prometheus.tcp.request_timeout controls the window of time in which the client has to send an HTTP request.
- prometheus.tcp.idle_timeout controls the window of time in which the client has to send more data (if any) within the context of an HTTP request.

If a load balancer or proxy is used between the Prometheus node and the RabbitMQ nodes it scrapes, the inactivity_timeout and idle_timeout values should be at least as large, and often greater than, the timeout and inactivity values used by the load balancer.

Here is an example configuration snippet that modifies the timeouts:

```ini
prometheus.tcp.idle_timeout = 120000
prometheus.tcp.inactivity_timeout = 120000
prometheus.tcp.request_timeout = 120000
```

### [Grafana Configuration](https://www.rabbitmq.com/prometheus.html#grafana-configuration)

The last component in this setup is Grafana. If this is your first time integrating Grafana with Prometheus, please follow the [official integration guide](https://prometheus.io/docs/visualization/grafana/).

After Grafana is integrated with the Prometheus instance that reads and stores RabbitMQ metrics, it is time to import the Grafana dashboards that Team RabbitMQ maintains. Please refer to the [the official Grafana tutorial](https://grafana.com/docs/reference/export_import/#importing-a-dashboard) on importing dashboards in Grafana.

Grafana dashboards for RabbitMQ and Erlang are open source and publicly from the [rabbitmq-server](https://github.com/rabbitmq/rabbitmq-server/tree/master/deps/rabbitmq_prometheus/docker/grafana/dashboards) GitHub repository.

To import **RabbitMQ-Overview** dashboard to Grafana:

1. Go to the [Grafana website](https://grafana.com/orgs/rabbitmq) to view the list of official RabbitMQ Grafana dashboards.

2. Select **RabbitMQ-Overview** [dashboard](https://grafana.com/grafana/dashboards/10991).

3. Click the **Download JSON** link or copy the dashboard ID.

4. Copy paste the file contents in Grafana, then click

    

   Load

   , as seen below:

   - Alternatively, paste the dashboard ID in the field **Grafana.com Dashboard**.

![Grafana Import Dashboard](https://www.rabbitmq.com/img/grafana-import-dashboard.png)

Repeat the process for all other Grafana dashboards that you would like to use with this RabbitMQ deployment.

Finally, switch the default data source used by Grafana to prometheus.

Congratulations! Your RabbitMQ is now monitored with Prometheus & Grafana!

## [Securing Prometheus Scraping Endpoint with TLS](https://www.rabbitmq.com/prometheus.html#tls)

The Prometheus metrics can be secured with TLS similar to the other listeners. For example, in the [configuration file](https://www.rabbitmq.com/configure.html#configuration-file)

```ini
prometheus.ssl.port       = 15691
prometheus.ssl.cacertfile = /full/path/to/ca_certificate.pem
prometheus.ssl.certfile   = /full/path/to/server_certificate.pem
prometheus.ssl.keyfile    = /full/path/to/server_key.pem
prometheus.ssl.password   = password-if-keyfile-is-encrypted
```

To enable TLS with [peer verification](https://www.rabbitmq.com/ssl.html#peer-verification), use a config similar to

```ini
prometheus.ssl.port       = 15691
prometheus.ssl.cacertfile = /full/path/to/ca_certificate.pem
prometheus.ssl.certfile   = /full/path/to/server_certificate.pem
prometheus.ssl.keyfile    = /full/path/to/server_key.pem
prometheus.ssl.password   = password-if-keyfile-is-encrypted
prometheus.ssl.verify     = verify_peer
prometheus.ssl.depth      = 2
prometheus.ssl.fail_if_no_peer_cert = true
```

## [Using Prometheus with RabbitMQ 3.7](https://www.rabbitmq.com/prometheus.html#3rd-party-plugin)

RabbitMQ versions prior to 3.8 can use a separate plugin, [prometheus_rabbitmq_exporter](https://github.com/deadtrickster/prometheus_rabbitmq_exporter), to expose metrics to Prometheus. The plugin uses [RabbitMQ HTTP API](https://www.rabbitmq.com/monitoring.html) internally and requires visualisation to be set up separately.

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!