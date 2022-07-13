## Customising messages

It is possible to customise messages that PerfTest publishes. This allows getting as close as possible to the target traffic or to populate queues with messages that real consumers will process.  可以自定义 PerfTest 发布的消息。 这允许尽可能接近目标流量或使用真实消费者将处理的消息填充队列。

### Message Properties

You can specify message properties with key/value pairs separated by commas:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --message-properties priority=5,timestamp=2007-12-03T10:15:30+01:00

```

The supported property keys are: `contentType`, `contentEncoding`, `deliveryMode`, `priority`, `correlationId`, `replyTo`, `expiration`, `messageId`, `timestamp`, `type`, `userId`, `appId`, `clusterId`. If some provided keys do not belong to the previous list, the pairs will be considered as headers (arbitrary key/value pairs):

```bash
bin/runjava com.rabbitmq.perf.PerfTest --message-properties priority=10,header1=value1,header2=value2


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

```

### Message Payload from Files

You can mimic real messages by specifying their content and content type. This can be useful when plugging real application consumers downstream. The content can come from one or several files and the content-type can be specified:  您可以通过指定内容和内容类型来模拟真实消息。 这在将真正的应用程序使用者插入下游时非常有用。 内容可以来自一个或多个文件，并且可以指定内容类型

```bash
bin/runjava com.rabbitmq.perf.PerfTest --consumers 0 --body content1.json,content2.json --body-content-type application/json


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

```

### Random JSON Payload

PerfTest can generate random JSON payload for messages. This is useful to experiment with traffic that (almost) always changes. To generate random JSON payloads, use the `--json-body` flag and the `--size` argument to specify the size in bytes:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --json-body --size 16000


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


--body-count
--body-field-count

```

Generate random values is costly, so PerfTest generates a pool of payloads upfront and uses them randomly in published messages. This way the generation of payloads does not impede publishing rate. There are 2 options to change the pre-generation of random JSON payload:  生成随机值的成本很高，因此 PerfTest 会预先生成一个有效负载池，并在发布的消息中随机使用它们。 这样，有效载荷的生成不会妨碍发布速率。 有 2 个选项可以更改随机 JSON 有效负载的预生成：

- `--body-count`: the size of the pool of payloads PerfTest will generate and use in published messages. The default size is 100. Increase this value if you want more randomness in published messages.  PerfTest 将在发布的消息中生成和使用的有效负载池的大小。 默认大小为 100。如果您希望发布的消息具有更多随机性，请增加此值。
- `--body-field-count`: the size of the pool of random strings used for field names and values in the JSON document. Before generating JSON payloads, PerfTest generates random strings and will use them randomly for field names and values in the JSON documents. The default value is 1,000. Increasing this value can be useful for "large" payloads (a few hundreds of kilobytes or more), which can "exhaust" the pool of random strings and then end up with duplicated field names. Duplicated field names are fine if the random JSON payloads are used to simulate traffic, but can be problematic if real consumers are plugged in and try to parse the JSON documents (JSON parsers do not always tolerate duplicated fields).  用于 JSON 文档中的字段名称和值的随机字符串池的大小。 在生成 JSON 有效负载之前，PerfTest 会生成随机字符串，并将它们随机用于 JSON 文档中的字段名称和值。 默认值为 1,000。 增加此值对于“大”有效负载（数百千字节或更多）很有用，这可能会“耗尽”随机字符串池，然后以重复的字段名称结束。 如果使用随机 JSON 有效负载来模拟流量，则重复的字段名称很好，但如果插入真正的消费者并尝试解析 JSON 文档（JSON 解析器并不总是容忍重复的字段），则可能会出现问题。

The defaults for `--body-count` and `--body-field-count` are usually fine, but can be increased for more randomness, at the cost of slower startup time and higher memory consumption.

Bear in mind that a large cache of generated payloads combined with a moderately large size can easily take up a significant amount of memory. As an example, `--json-body --body-count 50000 --size 100000` (50,000 payloads of 100 kB) will use about 5 GB of memory.

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

```
