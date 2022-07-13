## Workloads With a Large Number of Clients

A typical connected device workload (a.k.a "IoT workload") involves many producers and consumers (dozens or hundreds of thousands) that exchange messages at a low and mostly constant rate, usually a message every few seconds or minutes. Simulating such workloads requires a different set of settings compared to the workloads that have higher throughput and a small number of clients. With the appropriate set of flags, PerfTest can simulate IoT workloads without requiring too many resources, especially threads.  典型的连接设备工作负载（也称为“IoT 工作负载”）涉及许多生产者和消费者（数十或数十万），它们以低且几乎恒定的速率交换消息，通常每隔几秒或几分钟交换一次消息。 与具有更高吞吐量和少量客户端的工作负载相比，模拟此类工作负载需要一组不同的设置。 使用适当的标志集，PerfTest 可以模拟物联网工作负载，而无需太多资源，尤其是线程。

With an IoT workload, publishers usually don’t publish many messages per second, but rather a message every fixed period of time. This can be achieved by using the `--publishing-interval` flag instead of the `--rate` one. For example:  对于 IoT 工作负载，发布者通常不会每秒发布很多消息，而是每隔固定时间发布一条消息。 这可以通过使用 --publishing-interval 标志而不是 --rate 标志来实现。 例如：

```bash
bin/runjava com.rabbitmq.perf.PerfTest --publishing-interval 5

-x 1
-x,--producers
a single publisher without publisher confirms  生产者数量

-y 2
-y,--consumers
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

--queue
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

-st
--shutdown-timeout 20
The connection closing timeout can be set up with the --shutdown-timeout argument (or -st)  连接关闭的超时时间 

--queue-args x-max-length=10,x-dead-letter-exchange=some.exchange.name
PerfTest can create queues using provided queue arguments  定义队列参数

--quorum-queue
It is possible to use several arguments to create quorum queues  定义 quorum 类型的队列

--message-properties priority=10,header1=value1,header2=value2
You can specify message properties with key/value pairs separated by commas  定义消息属性

--body content1.json,content2.json  定义消息的来源的文件

--body-content-type application/json  定义消息类型

--json-body  生成随机的 json 格式的消息

--body-count  生成随机的 json 格式的消息时，预生成的随机的消息的数量

--body-field-count  生成随机的 json 格式的消息时，预生成的随机的字符串的数量

--variable-rate 100:60 --variable-rate 1000:10 --variable-rate 500:15
the publishing rate will be 100 messages per second for 60 seconds, then 1000 messages per second for 10 seconds, then 500 messages per second for 15 seconds, then back to 100 messages per second for 60 seconds, and so on.

--variable-size 1000:30 --variable-size 10000:20 --variable-size 5000:45
the size of published messages will be 1 kB for 30 seconds, then 10 kB for 20 seconds, then 5 kB for 45 seconds, then back to 1 kB for 30 seconds, and so on.

-L
--consumer-latency 1000  您可以使用固定或可变延迟值（以微秒为单位）模拟每条消息的处理时间。

--variable-latency 1000:60 --variable-latency 1000000:30
sets a variable consumer latency. it is set to 1 ms for 60 seconds then 1 second for 30 seconds

--queue-pattern 'perf-test-%d'
--queue-pattern-from 1
--queue-pattern-to 10
create the perf-test-1, perf-test-2, …, perf-test-10

--heartbeat-sender-threads 10  发送心跳检测的线程数量

--publishing-interval 5  每隔 5 秒发送一次消息

```

The command above makes the publisher publish a message every 5 seconds. To simulate a group of consumers, use the `--queue-pattern` flag to simulate many consumers across many queues:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 1000 --producers 1000 --consumers 1000 --heartbeat-sender-threads 10 --publishing-interval 5
```

To prevent publishers from publishing at roughly the same time and distribute the rate more evenly, use the `--producer-random-start-delay` option to add an random delay before the first published message:  为了防止发布者大致同时发布并更均匀地分配速率，请使用 --producer-random-start-delay 选项在第一条发布消息之前添加一个随机延迟：

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 1000 --producers 1000 --consumers 1000 --heartbeat-sender-threads 10 --publishing-interval 5 --producer-random-start-delay 120

-x 1
-x,--producers
a single publisher without publisher confirms  生产者数量

-y 2
-y,--consumers
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

--queue
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

-st
--shutdown-timeout 20
The connection closing timeout can be set up with the --shutdown-timeout argument (or -st)  连接关闭的超时时间 

--queue-args x-max-length=10,x-dead-letter-exchange=some.exchange.name
PerfTest can create queues using provided queue arguments  定义队列参数

--quorum-queue
It is possible to use several arguments to create quorum queues  定义 quorum 类型的队列

--message-properties priority=10,header1=value1,header2=value2
You can specify message properties with key/value pairs separated by commas  定义消息属性

--body content1.json,content2.json  定义消息的来源的文件

--body-content-type application/json  定义消息类型

--json-body  生成随机的 json 格式的消息

--body-count  生成随机的 json 格式的消息时，预生成的随机的消息的数量

--body-field-count  生成随机的 json 格式的消息时，预生成的随机的字符串的数量

--variable-rate 100:60 --variable-rate 1000:10 --variable-rate 500:15
the publishing rate will be 100 messages per second for 60 seconds, then 1000 messages per second for 10 seconds, then 500 messages per second for 15 seconds, then back to 100 messages per second for 60 seconds, and so on.

--variable-size 1000:30 --variable-size 10000:20 --variable-size 5000:45
the size of published messages will be 1 kB for 30 seconds, then 10 kB for 20 seconds, then 5 kB for 45 seconds, then back to 1 kB for 30 seconds, and so on.

-L
--consumer-latency 1000  您可以使用固定或可变延迟值（以微秒为单位）模拟每条消息的处理时间。

--variable-latency 1000:60 --variable-latency 1000000:30
sets a variable consumer latency. it is set to 1 ms for 60 seconds then 1 second for 30 seconds

--queue-pattern 'perf-test-%d'
--queue-pattern-from 1
--queue-pattern-to 10
create the perf-test-1, perf-test-2, …, perf-test-10

--heartbeat-sender-threads 10  发送心跳检测的线程数量

--publishing-interval 5  每隔 5 秒发送一次消息

--producer-random-start-delay 120  每次发送消息的随机延迟 each publisher will start with a random delay between 1 and 120 seconds.
```

With the command above, each publisher will start with a random delay between 1 and 120 seconds.

When using `--publishing-interval`, PerfTest will use one thread for scheduling publishing for all 50 producers. So 1000 producers should keep 20 threads busy for the publishing scheduling. This ratio can be decreased or increased with the `--producer-scheduler-threads` options depending on the load and the target environment. Very few threads can be used for very slow publishers:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' \
  --queue-pattern-from 1 --queue-pattern-to 1000 \
  --producers 1000 --consumers 1000 \
  --heartbeat-sender-threads 10 \
  --publishing-interval 60 --producer-random-start-delay 1800 \
  --producer-scheduler-threads 10

-x 1
-x,--producers
a single publisher without publisher confirms  生产者数量

-y 2
-y,--consumers
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

--queue
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

-st
--shutdown-timeout 20
The connection closing timeout can be set up with the --shutdown-timeout argument (or -st)  连接关闭的超时时间 

--queue-args x-max-length=10,x-dead-letter-exchange=some.exchange.name
PerfTest can create queues using provided queue arguments  定义队列参数

--quorum-queue
It is possible to use several arguments to create quorum queues  定义 quorum 类型的队列

--message-properties priority=10,header1=value1,header2=value2
You can specify message properties with key/value pairs separated by commas  定义消息属性

--body content1.json,content2.json  定义消息的来源的文件

--body-content-type application/json  定义消息类型

--json-body  生成随机的 json 格式的消息

--body-count  生成随机的 json 格式的消息时，预生成的随机的消息的数量

--body-field-count  生成随机的 json 格式的消息时，预生成的随机的字符串的数量

--variable-rate 100:60 --variable-rate 1000:10 --variable-rate 500:15
the publishing rate will be 100 messages per second for 60 seconds, then 1000 messages per second for 10 seconds, then 500 messages per second for 15 seconds, then back to 100 messages per second for 60 seconds, and so on.

--variable-size 1000:30 --variable-size 10000:20 --variable-size 5000:45
the size of published messages will be 1 kB for 30 seconds, then 10 kB for 20 seconds, then 5 kB for 45 seconds, then back to 1 kB for 30 seconds, and so on.

-L
--consumer-latency 1000  您可以使用固定或可变延迟值（以微秒为单位）模拟每条消息的处理时间。

--variable-latency 1000:60 --variable-latency 1000000:30
sets a variable consumer latency. it is set to 1 ms for 60 seconds then 1 second for 30 seconds

--queue-pattern 'perf-test-%d'
--queue-pattern-from 1
--queue-pattern-to 10
create the perf-test-1, perf-test-2, …, perf-test-10

--heartbeat-sender-threads 10  发送心跳检测的线程数量

--publishing-interval 5  每隔 5 秒发送一次消息

--producer-random-start-delay 120  每次发送消息的随机延迟 each publisher will start with a random delay between 1 and 120 seconds.

--producer-scheduler-threads 10  但使用 --publishing-interval 选项发送消息时，设置生产者使用的线程数量

```

In the example above, 1000 publishers will publish every 60 seconds with a random start-up delay between 1 second and 30 minutes (1800 seconds). They will be scheduled by only 10 threads (instead of 20 by default). Such delay values are suitable for long running tests.

Another option can be useful when simulating many consumers with a moderate message rate: `--consumers-thread-pools`. It allows to use a given number of thread pools for all the consumers, instead of one thread pool for each consumer by default. In the previous example, each consumer would use a 1-thread thread pool, which is overkill considering consumers processing is fast and producers publish one message every second. We can set the number of thread pools to use with `--consumers-thread-pools` and they will be shared by the consumers:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' \
  --queue-pattern-from 1 --queue-pattern-to 1000 \
  --producers 1000 --consumers 1000 \
  --heartbeat-sender-threads 10 \
  --publishing-interval 60 --producer-random-start-delay 1800 \
  --producer-scheduler-threads 10 \
  --consumers-thread-pools 10

-x 1
-x,--producers
a single publisher without publisher confirms  生产者数量

-y 2
-y,--consumers
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

--queue
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

-st
--shutdown-timeout 20
The connection closing timeout can be set up with the --shutdown-timeout argument (or -st)  连接关闭的超时时间 

--queue-args x-max-length=10,x-dead-letter-exchange=some.exchange.name
PerfTest can create queues using provided queue arguments  定义队列参数

--quorum-queue
It is possible to use several arguments to create quorum queues  定义 quorum 类型的队列

--message-properties priority=10,header1=value1,header2=value2
You can specify message properties with key/value pairs separated by commas  定义消息属性

--body content1.json,content2.json  定义消息的来源的文件

--body-content-type application/json  定义消息类型

--json-body  生成随机的 json 格式的消息

--body-count  生成随机的 json 格式的消息时，预生成的随机的消息的数量

--body-field-count  生成随机的 json 格式的消息时，预生成的随机的字符串的数量

--variable-rate 100:60 --variable-rate 1000:10 --variable-rate 500:15
the publishing rate will be 100 messages per second for 60 seconds, then 1000 messages per second for 10 seconds, then 500 messages per second for 15 seconds, then back to 100 messages per second for 60 seconds, and so on.

--variable-size 1000:30 --variable-size 10000:20 --variable-size 5000:45
the size of published messages will be 1 kB for 30 seconds, then 10 kB for 20 seconds, then 5 kB for 45 seconds, then back to 1 kB for 30 seconds, and so on.

-L
--consumer-latency 1000  您可以使用固定或可变延迟值（以微秒为单位）模拟每条消息的处理时间。

--variable-latency 1000:60 --variable-latency 1000000:30
sets a variable consumer latency. it is set to 1 ms for 60 seconds then 1 second for 30 seconds

--queue-pattern 'perf-test-%d'
--queue-pattern-from 1
--queue-pattern-to 10
create the perf-test-1, perf-test-2, …, perf-test-10

--heartbeat-sender-threads 10  发送心跳检测的线程数量

--publishing-interval 5  每隔 5 秒发送一次消息

--producer-random-start-delay 120  每次发送消息的随机延迟 each publisher will start with a random delay between 1 and 120 seconds.

--producer-scheduler-threads 10  但使用 --publishing-interval 选项发送消息时，设置生产者使用的线程数量

--consumers-thread-pools 10
默认情况下，每个消费者使用一个独立的线程池，通过设置 consumers-thread-pools 参数，1000 个消费者共享 10 个线程池

```

The previous example uses only 10 thread pools for all consumers instead of 1000 by default. These are 1-thread thread pools in this case, so this is 10 threads overall instead of 1000, another huge resource saving to simulate more clients with a single PerfTest instance for large IoT workloads.

By default, PerfTest uses blocking network socket I/O to communicate with the broker. This mode works fine for clients in many cases but the RabbitMQ Java client also supports an [asynchronous I/O mode](https://www.rabbitmq.com/api-guide.html#java-nio), where resources like threads can be easily tuned. The goal here is to use as few resources as possible to simulate as much load as possible with a single PerfTest instance. In the slow publisher example above, a handful of threads should be enough to handle the I/O. That’s what the `--nio-threads` flag is for:  默认情况下，PerfTest 使用阻塞网络套接字 I/O 与代理进行通信。 这种模式在很多情况下对客户端都适用，但 RabbitMQ Java 客户端也支持异步 I/O 模式，可以轻松调整线程等资源。 此处的目标是使用尽可能少的资源通过单个 PerfTest 实例模拟尽可能多的负载。 在上面的慢发布者示例中，少量线程应该足以处理 I/O。 这就是 --nio-threads 标志的用途：

```bash
bin/runjava com.rabbitmq.perf.PerfTest \
--queue-pattern 'perf-test-%d' \
--queue-pattern-from 1 \
--queue-pattern-to 1000 \
--producers 1000 \
--consumers 1000 \
--heartbeat-sender-threads 10 \
--publishing-interval 60 \
--producer-random-start-delay 1800 \
--producer-scheduler-threads 10 \
--nio-threads 10

-x 1
-x,--producers
a single publisher without publisher confirms  生产者数量

-y 2
-y,--consumers
two consumers (each receiving a copy of every message) that use automatic acknowledgement mode  消费者数量

--queue
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

-st
--shutdown-timeout 20
The connection closing timeout can be set up with the --shutdown-timeout argument (or -st)  连接关闭的超时时间 

--queue-args x-max-length=10,x-dead-letter-exchange=some.exchange.name
PerfTest can create queues using provided queue arguments  定义队列参数

--quorum-queue
It is possible to use several arguments to create quorum queues  定义 quorum 类型的队列

--message-properties priority=10,header1=value1,header2=value2
You can specify message properties with key/value pairs separated by commas  定义消息属性

--body content1.json,content2.json  定义消息的来源的文件

--body-content-type application/json  定义消息类型

--json-body  生成随机的 json 格式的消息

--body-count  生成随机的 json 格式的消息时，预生成的随机的消息的数量

--body-field-count  生成随机的 json 格式的消息时，预生成的随机的字符串的数量

--variable-rate 100:60 --variable-rate 1000:10 --variable-rate 500:15
the publishing rate will be 100 messages per second for 60 seconds, then 1000 messages per second for 10 seconds, then 500 messages per second for 15 seconds, then back to 100 messages per second for 60 seconds, and so on.

--variable-size 1000:30 --variable-size 10000:20 --variable-size 5000:45
the size of published messages will be 1 kB for 30 seconds, then 10 kB for 20 seconds, then 5 kB for 45 seconds, then back to 1 kB for 30 seconds, and so on.

-L
--consumer-latency 1000  您可以使用固定或可变延迟值（以微秒为单位）模拟每条消息的处理时间。

--variable-latency 1000:60 --variable-latency 1000000:30
sets a variable consumer latency. it is set to 1 ms for 60 seconds then 1 second for 30 seconds

--queue-pattern 'perf-test-%d'
--queue-pattern-from 1
--queue-pattern-to 10
create the perf-test-1, perf-test-2, …, perf-test-10

--heartbeat-sender-threads 10  发送心跳检测的线程数量

--publishing-interval 5  每隔 5 秒发送一次消息

--producer-random-start-delay 120  每次发送消息的随机延迟 each publisher will start with a random delay between 1 and 120 seconds.

--producer-scheduler-threads 10  但使用 --publishing-interval 选项发送消息时，设置生产者使用的线程数量

--consumers-thread-pools 10
默认情况下，每个消费者使用一个独立的线程池，通过设置 consumers-thread-pools 参数，1000 个消费者共享 10 个线程池

--nio-threads 10  使用异步 IO 时使用的线程数量

--disable-connection-recovery  使用 disable-connection-recovery 标志禁用连接恢复

```

This way PerfTest will use 12 threads for I/O over all the connections. With the default blocking I/O mode, each producer (or consumer) uses a thread for the I/O loop, that is 2000 threads to simulate 1000 producers and 1000 consumers. Using NIO in PerfTest can dramatically reduce the resources used to simulate workloads with a large number of connections with appropriate tuning.

Note that in NIO mode the number of threads used can increase temporarily when connections close unexpectedly and connection recovery kicks in. This is due to the NIO mode dispatching connection closing to non-I/O threads to avoid deadlocks. Connection recovery can be disabled with the `--disable-connection-recovery` flag.  请注意，在 NIO 模式下，当连接意外关闭和连接恢复启动时，使用的线程数可能会暂时增加。这是由于 NIO 模式将连接关闭分派给非 I/O 线程以避免死锁。 可以使用 --disable-connection-recovery 标志禁用连接恢复。