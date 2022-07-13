## Simulating High Loads

PerfTest can easily run hundreds of connections on a simple desktop machine. Each producer and consumer use a Java thread and a TCP connection though, so a PerfTest process can quickly run out of file descriptors, depending on the OS settings. A simple solution is to use several PerfTest processes, on the same machine or not. This is especially handy when combined with the [queue sequence](https://rabbitmq.github.io/rabbitmq-perf-test/stable/htmlsingle/#working-with-many-queues) feature.  PerfTest 可以在简单的台式机上轻松运行数百个连接。 不过，每个生产者和消费者都使用一个 Java 线程和一个 TCP 连接，因此 PerfTest 进程可能会迅速耗尽文件描述符，具体取决于操作系统设置。 一个简单的解决方案是在同一台机器上或不在同一台机器上使用多个 PerfTest 进程。 当与队列序列功能结合时，这特别方便。

The following command line launches a first PerfTest process that creates 500 queues (from `perf-test-1` to `perf-test-500`). Each queue will have 3 consumers and 1 producer sending messages to it:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 500 --producers 500 --consumers 1500
```

Then the following command line launches a second PerfTest process that creates 500 queues (from `perf-test-501` to `perf-test-1000`). Each queue will have 3 consumers and 1 producer sending messages to it:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' --queue-pattern-from 501 --queue-pattern-to 1000 --producers 500 --consumers 1500
```

Those 2 processes will simulate 1000 producers and 3000 consumers spread across 1000 queues.

A PerfTest process can exhaust its file descriptors limit and throw `java.lang.OutOfMemoryError: unable to create new native thread` exceptions. A first way to avoid this is to reduce the number of Java threads PerfTest uses with the `--heartbeat-sender-threads` option:  PerfTest 进程可能会耗尽其文件描述符限制并抛出 java.lang.OutOfMemoryError：无法创建新的本机线程异常。 避免这种情况的第一种方法是减少 PerfTest 使用 --heartbeat-sender-threads 选项使用的 Java 线程数：

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 1000 --producers 1000 --consumers 3000 --heartbeat-sender-threads 10
```

By default, each producer and consumer connection uses a dedicated thread to send heartbeats to the broker, so this is 4000 threads for heartbeats in the previous sample. Considering producers and consumers always communicate with the broker by publishing messages or sending acknowledgments, connections are never idle, so using 10 threads for heartbeats for the 4000 connections should be enough. Don’t hesitate to experiment to come up with the appropriate `--heartbeat-sender-threads` value for your use case.  默认情况下，每个生产者和消费者连接都使用一个专用线程向代理发送心跳，因此在上一个示例中，心跳是 4000 个线程。 考虑到生产者和消费者总是通过发布消息或发送确认与代理进行通信，连接永远不会空闲，因此对于 4000 个连接使用 10 个线程进行心跳应该足够了。 不要犹豫，尝试为您的用例提出适当的 --heartbeat-sender-threads 值。

Another way to avoid `java.lang.OutOfMemoryError: unable to create new native thread` exceptions is to tune the number of file descriptors allowed per process at the OS level, as some distributions use very low limits. Here the recommendations are the same as for the broker, so you can refer to our [networking guide](https://www.rabbitmq.com/networking.html#os-tuning).

```bash
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

```
