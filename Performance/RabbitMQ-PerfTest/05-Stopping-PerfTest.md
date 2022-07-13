## Stopping PerfTest

There are 2 reasons for a PerfTest run to stop:

- one of the limits has been reached (time limit, producer or consumer message count)

- the process is stopped by the user, e.g. by using Ctrl-C in the terminal

In both cases, PerfTest tries to exit as cleanly as possible, in a reasonable amount of time. Nevertheless, when PerfTest AMQP connections are throttled by the broker, because they’re publishing too fast or because broker [alarms](https://www.rabbitmq.com/alarms.html) have kicked in, it can take time to close them (several seconds or more for one connection).  在这两种情况下，PerfTest 都会尝试在合理的时间内尽可能干净地退出。 然而，当 PerfTest AMQP 连接被代理限制时，因为它们发布太快或因为代理警报已经启动，关闭它们可能需要时间（一个连接需要几秒钟或更长时间）。

If closing connections in the gentle way takes too long (5 seconds by default), PerfTest will move on to the most important resources to free and terminates. This can result in `client unexpectedly closed TCP connection` messages in the broker logs. Note this means the AMQP connection hasn’t been closed with the right sequence of AMQP frames, but the socket has been closed properly. There’s no resource leakage here.  如果以温和的方式关闭连接花费的时间太长（默认为 5 秒），PerfTest 将转移到最重要的资源以释放并终止。 这可能会导致客户端意外关闭代理日志中的 TCP 连接消息。 请注意，这意味着 AMQP 连接尚未使用正确的 AMQP 帧序列关闭，但套接字已正确关闭。 这里没有资源泄漏。

The connection closing timeout can be set up with the `--shutdown-timeout` argument (or `-st`). The default timeout can be increased to let more time to close connections, e.g. the command below uses a shutdown timeout of 20 seconds:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --shutdown-timeout 20
```

The connection closing sequence can also be skipped by setting the timeout to 0 or any negative value:

```bash
bin/runjava com.rabbitmq.perf.PerfTest --shutdown-timeout -1
```

With the previous command, PerfTest won’t even try to close AMQP connections, it will exit as fast as possible, freeing only the most important resources. This is perfectly acceptable when performing runs on a test environment.


```bash
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

-st
--shutdown-timeout 20
The connection closing timeout can be set up with the --shutdown-timeout argument (or -st)  连接关闭的超时时间 

```

