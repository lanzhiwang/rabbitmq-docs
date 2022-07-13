## Basic Usage

https://rabbitmq.github.io/rabbitmq-perf-test/stable/htmlsingle/#basic-usage

The most basic way of running PerfTest only specifies a URI to connect to, a number of publishers to use (say, 1) and a number of consumers to use (say, 2). Note that RabbitMQ Java client can achieve high rates for publishing (up to 80 to 90K messages per second per connection), given enough bandwidth and when some safety measures (publisher confirms) are disabled, so overprovisioning publishers is rarely necessary (unless that’s a specific objective of the test).  运行 PerfTest 的最基本方法仅指定要连接的 URI、要使用的发布者数量（例如，1）和要使用的消费者数量（例如，2）。 请注意，RabbitMQ Java 客户端可以实现高发布速率（每个连接每秒最多 80 到 90K 条消息），如果有足够的带宽并且当某些安全措施（发布者确认）被禁用时，因此很少需要过度配置发布者（除非这是特定的） 测试目的）。

------

The following command runs PerfTest with a single publisher without publisher confirms, two consumers (each receiving a copy of every message) that use automatic acknowledgement mode and a single queue named “throughput-test-x1-y2”. Publishers will publish as quickly as possible, without any rate limiting. Results will be prefixed with “test1” for easier identification and comparison:  以下命令运行 PerfTest，其中有一个没有发布者确认的发布者、两个使用自动确认模式的消费者（每个接收每条消息的副本）和一个名为“throughput-test-x1-y2”的队列。 出版商将尽快发布，没有任何速率限制。 结果将以“test1”为前缀，以便于识别和比较：

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-1" -a --id "test 1"

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

```

------

This modification will use 2 publishers and 4 consumers, typically yielding higher throughput given enough CPU cores on the machine and RabbitMQ nodes:  此修改将使用 2 个发布者和 4 个消费者，如果机器和 RabbitMQ 节点上有足够的 CPU 内核，通常会产生更高的吞吐量：

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 2 -y 4 -u "throughput-test-2" -a --id "test 2"

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

```

------

This modification switches consumers to manual acknowledgements:  此修改将消费者切换到手动确认：

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-3" --id "test 3"

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

```

------

This modification changes message size from default (12 bytes) to 4 kB:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-4" --id "test 4" -s 4000

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小
```

------

PerfTest can use durable queues and persistent messages:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-5" --id "test-5" -f persistent

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息
```

------

When PerfTest is running, it is important to monitor various publisher and consumer metrics provided by the [management UI](https://www.rabbitmq.com/management.html). For example, it is possible to see how much network bandwidth a publisher has been using recently on the connection page.  当 PerfTest 运行时，监控管理 UI 提供的各种发布者和消费者指标很重要。 例如，可以在连接页面上查看发布者最近使用了多少网络带宽。

Queue page demonstrates message rates, consumer count, acknowledgement mode used by the consumers, consumer utilisation and message location break down (disk, RAM, paged out transient messages, etc). When durable queues and persistent messages are used, node I/O and message store/queue index operation metrics become particularly important to monitor.  队列页面展示了消息速率、消费者数量、消费者使用的确认模式、消费者利用率和消息位置细分（磁盘、RAM、分页出的临时消息等）。 当使用持久队列和持久消息时，节点 I/O 和消息存储/队列索引操作指标变得尤为重要，需要监控。

------

Consumers can ack multiple messages at once, for example, 100 in this configuration:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-6" --id "test-6" -f persistent --multi-ack-every 100

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量
```

------

[Consumer prefetch (QoS)](https://www.rabbitmq.com/confirms.html) can be configured as well (in this example to 500):

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-7" --id "test-7" -f persistent --multi-ack-every 200 -q 500

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量
```

------

Publisher confirms can be used with a maximum of N outstanding publishes:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-8" --id "test-8" -f persistent -q 500 -c 500

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布
```

------

PerfTest can publish only a certain number of messages instead of running until shut down:  PerfTest 只能发布一定数量的消息，而不是运行直到关闭：

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-10" --id "test-10" -f persistent -q 500 -pmessages 100000

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布

-pmessages 100000
PerfTest can publish only a certain number of messages  生产消息的数量

```

------

Publisher rate can be limited:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-11" --id "test-11" -f persistent -q 500 --rate 5000

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布

-pmessages 100000
PerfTest can publish only a certain number of messages  生产消息的数量

--rate 5000
Publisher rate can be limited  生产者生产消息的速率

```

------

Consumer rate can be limited as well to simulate slower consumers or create a backlog:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-12" --id "test-12" -f persistent --rate 5000 --consumer-rate 2000

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布

-pmessages 100000
PerfTest can publish only a certain number of messages  生产消息的数量

--rate 5000
Publisher rate can be limited  生产者生产消息的速率

--consumer-rate 2000
Consumer rate can be limited  消费者速率

```
Note that the consumer rate limit is applied per consumer, so in the configuration above the limit is actually 2 * 2000 = 4000 deliveries/second.

------

PerfTest can be configured to run for a limited amount of time in seconds with the `-z` option:

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x 1 -y 2 -u "throughput-test-13" --id "test-13" -f persistent -z 30

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布

-pmessages 100000
PerfTest can publish only a certain number of messages  生产消息的数量

--rate 5000
Publisher rate can be limited  生产者生产消息的速率

--consumer-rate 2000
Consumer rate can be limited  消费者速率

-z 30
run for a limited amount of time in seconds  运行测试的总时间，单位：秒

```



------

Running PerfTest without consumers and with a limited number of messages can be used to pre-populate a queue, e.g. with 1M messages 1 kB in size each::

```bash
bin/runjava com.rabbitmq.perf.PerfTest -y0 -p -u "throughput-test-14" -s 1000 -C 1000000 --id "test-14" -f persistent

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布

-C
-pmessages 100000
PerfTest can publish only a certain number of messages  生产消息的数量

-r
--rate 5000
Publisher rate can be limited  生产者生产消息的速率

--consumer-rate 2000
Consumer rate can be limited  消费者速率

-z 30
run for a limited amount of time in seconds  运行测试的总时间，单位：秒

-p
pre-populate a queue  预先生产数据到队列中

-D,--cmessages <arg>
consumer message count  消费消息的数量

```

------

Use the `-D` option to limit the number of consumed messages. Note the `-z` (time limit), `-C` (number of published messages), and `-D` (number of consumed messages) options can be used together but their combination can lead to funny results. `-x 1 -y 1 -C 10 -D 20 -r 1` would for example stop the producer once 10 messages have been published, letting the consumer wait forever the remaining 10 messages (as the publisher is stopped).  使用 `-D` 选项来限制消费的消息数量。 请注意，`-z`（时间限制）、`-C`（已发布消息的数量）和`-D`（已消费消息的数量）选项可以一起使用，但它们的组合会导致有趣的结果。 例如，`-x 1 -y 1 -C 10 -D 20 -r 1` 会在发布 10 条消息后停止生产者，让消费者永远等待剩余的 10 条消息（因为发布者已停止）。

To consume from a pre-declared and pre-populated queue without starting any publishers, use  要在不启动任何发布者的情况下从预先声明和预先填充的队列中消费，请使用

```bash
bin/runjava com.rabbitmq.perf.PerfTest -x0 -y10 -p -u "throughput-test-14" --id "test-15"

-x 1
a single publisher without publisher confirms  生产者数量

-y 2
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

-u "throughput-test-1"
queue named  队列名称

-a
without any rate limiting 消费者自动 ack

--id "test 1"
Results will be prefixed with “test1”  结果标识符

-s 4000
message size from default (12 bytes) to 4 kB  消息大小

-f persistent
durable queues and persistent messages  持久化队列和持久化消息

--multi-ack-every 100
Consumers can ack multiple messages at once, for example, 100 in this configuration  消费者手动确认时，一次确认的消息数量

-q 500
Consumer prefetch (QoS) can be configured as well (in this example to 500)  Consumer prefetch 数量

-c 500
Publisher confirms can be used with a maximum of N outstanding publishes  发布者确认最多可用于 N 个未完成的发布

-C
-pmessages 100000
PerfTest can publish only a certain number of messages  生产消息的数量

-r
--rate 5000
Publisher rate can be limited  生产者生产消息的速率

--consumer-rate 2000
Consumer rate can be limited  消费者速率

-z 30
run for a limited amount of time in seconds  运行测试的总时间，单位：秒

-p
pre-populate a queue  预先生产数据到队列中

-D,--cmessages <arg>
consumer message count  消费消息的数量

```

PerfTest is useful for establishing baseline cluster throughput with various configurations but does not simulate many other aspects of real world applications. It is also biased towards very simplistic workloads that use a single queue, which provides [limited CPU utilisation](https://www.rabbitmq.com/queues.html) on RabbitMQ nodes and is not recommended for most cases.  PerfTest 可用于建立具有各种配置的基线集群吞吐量，但不会模拟现实世界应用程序的许多其他方面。 它还偏向于使用单个队列的非常简单的工作负载，这在 RabbitMQ 节点上提供有限的 CPU 利用率，并且在大多数情况下不推荐使用。

Multiple PerfTest instances running simultaneously can be used to simulate more realistic workloads.  同时运行的多个 PerfTest 实例可用于模拟更真实的工作负载。

