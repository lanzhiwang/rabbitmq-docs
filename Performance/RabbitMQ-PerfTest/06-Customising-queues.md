## Customising queues

PerfTest can create queues using provided queue arguments:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-args x-max-length=10
```

The previous command will create a [queue with a length limit](https://www.rabbitmq.com/maxlength.html) of 10. You can also provide several queue arguments by separating the key/value pairs with commas:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --queue-args x-max-length=10,x-dead-letter-exchange=some.exchange.name
```

It is possible to use several arguments to create [quorum queues](https://rabbitmq.com/quorum-queues.html), but PerfTest provides a `--quorum-queue` flag to do that:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --quorum-queue --queue name
```

`--quorum-queue` is a shortcut for `--flag persistent --queue-args x-queue-type=quorum --auto-delete false`. Note a quorum queue cannot have a server-generated name, so the `--queue` argument must be used to specify the name of the queue(s).

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

```
