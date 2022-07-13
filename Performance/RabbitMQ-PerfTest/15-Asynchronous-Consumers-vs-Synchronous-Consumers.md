## Asynchronous Consumers vs Synchronous Consumers

Consumers are asynchronous by default in PerfTest. This means they are registered with the AMQP `basic.consume` method and the broker pushes messages to them. This is the optimal way to consume messages. PerfTest also provides the `--polling` and `--polling-interval` options to consume messages by polling the broker with the AMQP `basic.get` method. These options are available to evaluate the performance and the effects of `basic.get`, but real applications should avoid using `basic.get` as much as possible because it has several drawbacks compared to asynchronous consumers: it needs a network round trip for each message, it typically keeps a thread busy for polling in the application, and it intrinsically increases latency.  在 PerfTest 中，消费者默认是异步的。 这意味着它们使用 AMQP basic.consume 方法注册，并且代理将消息推送给它们。 这是消费消息的最佳方式。 PerfTest 还提供了 --polling 和 --polling-interval 选项，通过使用 AMQP basic.get 方法轮询代理来使用消息。 这些选项可用于评估 basic.get 的性能和效果，但实际应用程序应尽可能避免使用 basic.get，因为与异步消费者相比，它有几个缺点：每条消息都需要网络往返， 通常使线程忙于在应用程序中进行轮询，并且它本质上会增加延迟。

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

--publishing-interval 5  每隔 5 秒发送一次消息

--producer-random-start-delay 120  每次发送消息的随机延迟 each publisher will start with a random delay between 1 and 120 seconds.

--producer-scheduler-threads 10  但使用 --publishing-interval 选项发送消息时，设置生产者使用的线程数量

--consumers-thread-pools 10
默认情况下，每个消费者使用一个独立的线程池，通过设置 consumers-thread-pools 参数，1000 个消费者共享 10 个线程池

--nio-threads 10  使用异步 IO 时使用的线程数量

--disable-connection-recovery  使用 disable-connection-recovery 标志禁用连接恢复

--polling
--polling-interval
使用 basic.get 同步获取消息

```
