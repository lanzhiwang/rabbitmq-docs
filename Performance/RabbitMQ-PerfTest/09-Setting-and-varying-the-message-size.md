## Setting and varying the message size

The default size of the messages that PerfTest publishes is 12 bytes (PerfTest stores some data in the message to calculate latency on the consumer side).

It is possible to make messages bigger with the `--size` (`-s`) option, e.g. to publish 4 kB messages:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --size 4000
```

The `--variable-size` (`-vs`) option allows to specify different message sizes for periods of time, e.g.:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --variable-size 1000:30 --variable-size 10000:20 --variable-size 5000:45
```

The variable rate option uses the `[SIZE]:[DURATION]` syntax, where `SIZE` is in bytes and `DURATION` is in seconds. In the previous example, the size of published messages will be 1 kB for 30 seconds, then 10 kB for 20 seconds, then 5 kB for 45 seconds, then back to 1 kB for 30 seconds, and so on.

```bash
-x 1
a single publisher without publisher confirms  生产者数量

-y 2
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

```
