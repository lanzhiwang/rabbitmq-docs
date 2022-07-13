## Monitoring

PerfTest can gather metrics and make them available to various monitoring systems. Metrics include messaging-centric metrics (message latency, number of connections and channels, number of published messages, etc) as well as OS process and JVM metrics (memory, CPU usage, garbage collection, JVM heap, etc).  PerfTest 可以收集指标并将其提供给各种监控系统。 指标包括以消息为中心的指标（消息延迟、连接和通道数、已发布消息的数量等）以及操作系统进程和 JVM 指标（内存、CPU 使用率、垃圾收集、JVM 堆等）。

Here is how to list the available metrics options:

```bash
$ ./runjava com.rabbitmq.perf.PerfTest --metrics-help
$ docker run -it --rm pivotalrabbitmq/perf-test:latest --metrics-help
usage: <program>

-mc,--metrics-client
enable client metrics

-mcl,--metrics-class-loader
enable JVM class loader metrics

-mda,--metrics-datadog
enable Datadog metrics

-mdak,--metrics-datadog-application-key <arg>
Datadog application key

-mdd,--metrics-datadog-descriptions
if meter descriptions should be sent to Datadog

-mdh,--metrics-datadog-host-tag <arg>
tag that will be mapped to "host" when shipping metrics to datadog

-mdk,--metrics-datadog-api-key <arg>
Datadog API key

-mds,--metrics-datadog-step-size <arg>
step size (reporting frequency) to use in seconds, default is 10 seconds

-mdu,--metrics-datadog-uri <arg>
URI to ship metrics, useful when using a proxy, default is https://app.datadoghq.com

-mjgc,--metrics-jvm-gc
enable JVM GC metrics

-mjm,--metrics-jvm-memory
enable JVM memory metrics

-mjp,--metrics-processor
enable processor metrics (gathered by JVM)

-mjt,--metrics-jvm-thread
enable JVM thread metrics

-mjx,--metrics-jmx
enable JMX metrics

-mpe,--metrics-prometheus-endpoint <arg>
the HTTP metrics endpoint, default is /metrics

-mpp,--metrics-prometheus-port <arg>
the port to launch the HTTP metrics endpoint on, default is 8080

-mpr,--metrics-prometheus
enable Prometheus metrics

-mpx,--metrics-prefix <arg>
prefix for PerfTest metrics, default is perftest_

-mt,--metrics-tags <arg>
metrics tags as key-value pairs separated by commas
```

This command displays the available flags to enable the various metrics PerfTest can gather, as well as options to configure the exposure to the monitoring systems PerfTest supports.  此命令显示可用标志以启用 PerfTest 可以收集的各种指标，以及配置暴露于 PerfTest 支持的监控系统的选项。

### Supported Metrics

Here are the metrics PerfTest can gather:

- default metrics: number of published, returned, confirmed, nacked, and consumed messages, message latency, publisher confirm latency. Message latency is a major concern in many types of workload, it can be easily monitored here. [Publisher confirm](https://www.rabbitmq.com/confirms.html#publisher-confirms) latency reflects the time a message can be considered unsafe. It is calculated as soon as the `--confirm`/`-c` option is used. Default metrics are available as long as PerfTest support for a monitoring system is enabled.  默认指标：发布、返回、确认、nack 和消费消息的数量、消息延迟、发布者确认延迟。 消息延迟是许多类型工作负载中的主要问题，可以在这里轻松监控。 发布者确认延迟反映了消息被视为不安全的时间。 一旦使用 --confirm/-c 选项，它就会被计算出来。 只要启用了对监控系统的 PerfTest 支持，默认指标就可用。

- client metrics: these are the [Java Client metrics](https://www.rabbitmq.com/api-guide.html#metrics). Enabling these metrics shouldn’t bring much compared to the default PerfTest metrics, except to see how PerfTest behaves with regards to number of open connections and channels for instance. Client metrics are enabled with the `-mc` or `--metrics-client` flag.  客户端指标：这些是 Java 客户端指标。 与默认的 PerfTest 指标相比，启用这些指标应该不会带来太多好处，除了查看 PerfTest 在打开连接数和通道数方面的表现如何。 客户端指标通过 -mc 或 --metrics-client 标志启用。

- JVM memory metrics: these metrics report memory usage of the JVM, e.g. current heap size, etc. They can be useful to have a better understanding of the client behavior, e.g. heap memory fluctuation could be due to frequent garbage collection that could explain high latency numbers. These metrics are enabled with the `-mjm` or `--metrics-jvm-memory` flag.  JVM 内存指标：这些指标报告 JVM 的内存使用情况，例如 当前堆大小等。它们对于更好地理解客户端行为很有用，例如 堆内存波动可能是由于频繁的垃圾收集可以解释高延迟数字。 这些指标通过 -mjm 或 --metrics-jvm-memory 标志启用。

- JVM thread metrics: these metrics report the number of JVM threads used in the PerfTest process, as well as their state. This can be useful to optimize the usage of PerfTest to simulate [high loads with fewer resources](https://rabbitmq.github.io/rabbitmq-perf-test/stable/htmlsingle/#workloads-with-a-large-number-of-clients). These metrics are enabled with the `-mjt` or `--metrics-jvm-thread` flag.  JVM 线程指标：这些指标报告 PerfTest 进程中使用的 JVM 线程数及其状态。 这对于优化 PerfTest 的使用以用更少的资源模拟高负载很有用。 这些指标通过 -mjt 或 --metrics-jvm-thread 标志启用。

- JVM GC metrics: these metrics reports garbage collection activity. They can vary depending on the JVM used, its version, and the GC settings. They can be useful to correlate the GC activity with PerfTest behavior, e.g. abnormal low throughput because of very frequent garbage collection. These metrics are enabled with the `-mjgc` or `--metrics-jvm-gc` flag.  JVM GC 指标：这些指标报告垃圾收集活动。 它们可能因使用的 JVM、其版本和 GC 设置而异。 它们可用于将 GC 活动与 PerfTest 行为相关联，例如 由于非常频繁的垃圾收集，异常的低吞吐量。 这些指标通过 -mjgc 或 --metrics-jvm-gc 标志启用。

- JVM class loader metrics: the number of loaded and unloaded classes. These metrics are enabled with the `-mcl` or `--metrics-class-loader` flag.  JVM 类加载器指标：加载和卸载的类的数量。 这些指标通过 -mcl 或 --metrics-class-loader 标志启用。

- Processor metrics: there metrics report CPU activity as gathered by the JVM. They can be enabled with the `-mjp` or `--metrics-processor` flag.  处理器指标：指标报告 JVM 收集的 CPU 活动。 可以使用 -mjp 或 --metrics-processor 标志启用它们。

> The JVM-related metrics are not available when using the [native executable](https://rabbitmq.github.io/rabbitmq-perf-test/stable/htmlsingle/#native-executable).

### Tags

One can specify metrics tags with the `-mt` or `--metrics-tags` options, e.g. `--metrics-tags env=performance,datacenter=eu` to tell monitoring systems that those metrics are from the `performance` environment located in the `eu` data center. Monitoring systems that support dimensions can then make it easier to navigate across metrics (group by, drill down). See [Micrometer](https://micrometer.io/) documentation for more information about tags and dimensions.  可以使用 -mt 或 --metrics-tags 选项指定指标标签，例如 --metrics-tags env=performance,datacenter=eu 告诉监控系统这些指标来自位于欧盟数据中心的性能环境。 然后，支持维度的监控系统可以更轻松地跨指标导航（分组、向下钻取）。 有关标签和尺寸的更多信息，请参阅千分尺文档。

### Supported Monitoring Systems

PerfTest builds on top [Micrometer](https://micrometer.io/) to report gathered metrics to various monitoring systems. Nevertheless, not all systems supported by Micrometer are actually supported by PerfTest. PerfTest currently supports [Datadog](https://www.datadoghq.com/), [JMX](https://en.wikipedia.org/wiki/Java_Management_Extensions), and [Prometheus](https://prometheus.io/). Don’t hesitate to [request support for other monitoring systems](https://github.com/rabbitmq/rabbitmq-perf-test/issues).

#### Datadog

The API key is the only required option to send metrics to Datadog:

```bash
$ ./runjava com.rabbitmq.perf.PerfTest --metrics-datadog-api-key YOUR_API_KEY
```

Another useful option is the step size or reporting frequency. The default value is 10 seconds.

```bash
$ ./runjava com.rabbitmq.perf.PerfTest --metrics-datadog-api-key YOUR_API_KEY --metrics-datadog-step-size 20
```

#### JMX

JMX support provides a simple way to view metrics locally. Use the `--metrics-jmx` flag to export metrics to JMX:

```bash
$ ./runjava com.rabbitmq.perf.PerfTest --metrics-jmx
```

#### Prometheus

Use the `-mpr` or `--metrics-prometheus` flag to enable metrics reporting to Prometheus:

```bash
$ ./runjava com.rabbitmq.perf.PerfTest --metrics-prometheus
```

Prometheus expects to scrape or poll individual app instances for metrics, so PerfTest starts up a web server listening on port 8080 and exposes metrics on the `/metrics` endpoint. These defaults can be changed:

```bash
$ ./runjava com.rabbitmq.perf.PerfTest --metrics-prometheus --metrics-prometheus-port 8090 --metrics-prometheus-endpoint perf-test-metrics
```

