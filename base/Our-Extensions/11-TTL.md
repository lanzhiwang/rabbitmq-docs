# Time-To-Live and Expiration

* https://www.rabbitmq.com/ttl.html

## Overview

RabbitMQ allows you to set TTL (time to live) for both messages and queues. Expired messages and queues will be deleted: the specifics will be covered in more detail later in this guide.
RabbitMQ 允许您为消息和队列设置 TTL（生存时间）。 过期的消息和队列将被删除：本指南稍后将更详细地介绍具体细节。

TTL behavior is controlled by [optional queue arguments](https://www.rabbitmq.com/queues.html) and best done using a [policy](https://www.rabbitmq.com/parameters.html).
TTL 行为由可选的队列参数控制，最好使用策略来完成。

Message TTL can be applied to a single queue, a group of queues or applied on the message-by-message basis.
消息 TTL 可以应用于单个队列、一组队列或逐条消息应用。

TTL settings also can be enforced by [operator policies](https://www.rabbitmq.com/parameters.html#operator-policies).
TTL 设置也可以由运营商策略强制执行。

## Per-Queue Message TTL in Queues

Message TTL can be set for a given queue by setting the `message-ttl` argument with a [policy](https://www.rabbitmq.com/parameters.html#policies) or by specifying the same argument at the time of queue declaration.
可以通过使用策略设置 message-ttl 参数或在队列声明时指定相同参数来为给定队列设置消息 TTL。

A message that has been in the queue for longer than the configured TTL is said to be *dead*. Note that a message routed to multiple queues can die at different times, or not at all, in each queue in which it resides. The death of a message in one queue has no impact on the life of the same message in other queues.
在队列中的时间超过配置的 TTL 的消息被称为死消息。 请注意，路由到多个队列的消息在它所在的每个队列中可能会在不同时间死亡，或者根本不会死亡。 一个队列中消息的死亡不会影响其他队列中同一消息的生命。

The server guarantees that dead messages will not be delivered using basic.deliver (to a consumer) or included into a basic.get-ok response (for one-off fetch operations). Further, the server will try to remove messages at or shortly after their TTL-based expiry.
服务器保证死消息不会使用 basic.deliver 传递（给消费者）或包含在 basic.get-ok 响应中（对于一次性获取操作）。 此外，服务器将尝试在基于 TTL 的到期时或之后不久删除消息。

The value of the TTL argument or policy must be a non-negative integer (0 <= n), describing the TTL period in milliseconds. Thus a value of 1000 means that a message added to the queue will live in the queue for 1 second or until it is delivered to a consumer. The argument can be of AMQP 0-9-1 type short-short-int, short-int, long-int, or long-long-int.
TTL 参数或策略的值必须是非负整数 (0 <= n)，以毫秒为单位描述 TTL 周期。 因此，值为 1000 意味着添加到队列的消息将在队列中存在 1 秒或直到它被传递给消费者。 参数可以是 AMQP 0-9-1 类型的 short-short-int、short-int、long-int 或 long-long-int。

### Define Message TTL for Queues Using a Policy

To specify a TTL using policy, add the key "message-ttl" to a policy definition:

|                       |                                                                             |
| --------------------- | --------------------------------------------------------------------------- |
| rabbitmqctl           | rabbitmqctl set_policy TTL ".*" '{"message-ttl":60000}' --apply-to queues   |
| rabbitmqctl (Windows) | rabbitmqctl set_policy TTL ".*" "{""message-ttl"":60000}" --apply-to queues |

This applies a TTL of 60 seconds to all queues.

### Define Message TTL for Queues Using x-arguments During Declaration

This example in Java creates a queue in which messages may reside for at most 60 seconds:

```
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-message-ttl", 60000);
channel.queueDeclare("myqueue", false, false, false, args);
```

The same example in C#:

```
var args = new Dictionary<string, object>();
args.Add("x-message-ttl", 60000);
model.QueueDeclare("myqueue", false, false, false, args);
```

It is possible to apply a message TTL policy to a queue which already has messages in it but this involves [some caveats](https://www.rabbitmq.com/ttl.html#per-message-ttl-caveats).
可以将消息 TTL 策略应用于其中已有消息的队列，但这涉及一些注意事项。

The original expiry time of a message is preserved if it is requeued (for example due to the use of an AMQP method that features a requeue parameter, or due to a channel closure).
如果消息重新排队（例如，由于使用具有重新排队参数的 AMQP 方法，或由于通道关闭），则会保留消息的原始到期时间。

Setting the TTL to 0 causes messages to be expired upon reaching a queue unless they can be delivered to a consumer immediately. Thus this provides an alternative to the immediate publishing flag, which the RabbitMQ server does not support. Unlike that flag, no basic.returns are issued, and if a dead letter exchange is set then messages will be dead-lettered.
将 TTL 设置为 0 会导致消息在到达队列后过期，除非它们可以立即传递给消费者。 因此，这提供了 RabbitMQ 服务器不支持的立即发布标志的替代方案。 与那个标志不同的是，没有 basic.returns 被发出，如果设置了死信交换，那么消息将是死信的。

## Per-Message TTL in Publishers

A TTL can be specified on a per-message basis, by setting the [expiration property](https://www.rabbitmq.com/publishers.html#message-properties) when publishing a message.
可以在每条消息的基础上指定 TTL，方法是在发布消息时设置过期属性。

The value of the `expiration` field describes the TTL period in milliseconds. The same constraints as for `x-message-ttl` apply. Since the expiration field must be a string, the broker will (only) accept the string representation of the number.
过期字段的值以毫秒为单位描述了 TTL 周期。 适用与 x-message-ttl 相同的约束。 由于过期字段必须是字符串，代理将（仅）接受数字的字符串表示。

When both a per-queue and a per-message TTL are specified, the lower value between the two will be chosen.
当指定了每个队列和每个消息的 TTL 时，将选择两者之间的较低值。

This example uses [RabbitMQ Java client](https://www.rabbitmq.com/api-guide.html) to publish a message which can reside in the queue for at most 60 seconds:

```
byte[] messageBodyBytes = "Hello, world!".getBytes();
AMQP.BasicProperties properties = new AMQP.BasicProperties.Builder()
                                   .expiration("60000")
                                   .build();
channel.basicPublish("my-exchange", "routing-key", properties, messageBodyBytes);
```

The same example in C#:

```
byte[] messageBodyBytes = System.Text.Encoding.UTF8.GetBytes("Hello, world!");
IBasicProperties props = model.CreateBasicProperties();
props.ContentType = "text/plain";
props.DeliveryMode = 2;
props.Expiration = "60000";
model.BasicPublish(exchangeName,
                   routingKey, props,
                   messageBodyBytes);
```

## Caveats

Queues that had a per-message TTL applied to them retroactively (when they already had messages) will discard the messages when specific events occur. Only when expired messages reach the head of a queue will they actually be discarded (or dead-lettered). Consumers will not have expired messages delivered to them. Keep in mind that there can be a natural race condition between message expiration and consumer delivery, e.g. a message can expire after it was written to the socket but before it has reached a consumer.
当特定事件发生时，追溯应用每条消息 TTL 的队列（当它们已经有消息时）将丢弃消息。 只有当过期消息到达队列头部时，它们才会真正被丢弃（或死信）。 消费者不会收到过期的消息。 请记住，消息过期和消费者交付之间可能存在自然竞争条件，例如 消息在写入套接字之后但在到达消费者之前可能会过期。

When setting per-message TTL expired messages can queue up behind non-expired ones until the latter are consumed or expired. Hence resources used by such expired messages will not be freed, and they will be counted in queue statistics (e.g. the number of messages in the queue).
设置每条消息 TTL 过期消息时，可以在未过期消息之后排队，直到后者被消耗或过期。 因此，此类过期消息使用的资源将不会被释放，并且它们将被计入队列统计信息（例如队列中的消息数）。

When retroactively applying a per-message TTL policy, it is recommended to have consumers online to make sure the messages are discarded quicker.
当追溯应用每条消息的 TTL 策略时，建议让消费者在线以确保更快地丢弃消息。

Given this behaviour of per-message TTL settings on existing queues, when the need to delete messages to free up resources arises, queue TTL should be used instead (or queue purging, or queue deletion).
考虑到现有队列上每条消息 TTL 设置的这种行为，当需要删除消息以释放资源时，应该改用队列 TTL（或队列清除或队列删除）。

## Queue TTL

TTL can also be set on queues, not just queue contents. This feature can be used together with the [auto-delete queue property](https://www.rabbitmq.com/queues.html).
TTL 也可以在队列上设置，而不仅仅是队列内容。 此功能可以与自动删除队列属性一起使用。

Setting TTL (expiration) on queues generally only makes sense for transient (non-durable) classic queues. Streams do not support expiration.
在队列上设置 TTL（到期）通常只对瞬态（非持久）经典队列有意义。 流不支持过期。

Queues will expire after a period of time only when they are not used (a queue is used if it has online consumers).
队列只有在不被使用的情况下才会在一段时间后过期（如果有在线消费者，就会使用队列）。

Expiry time can be set for a given queue by setting the x-expires argument to queue.declare, or by setting the expires [policy](https://www.rabbitmq.com/parameters.html#policies). This controls for how long a queue can be unused before it is automatically deleted. Unused means the queue has no consumers, the queue has not been recently redeclared (redeclaring renews the lease), and basic.get has not been invoked for a duration of at least the expiration period. This can be used, for example, for RPC-style reply queues, where many queues can be created which may never be drained.
可以通过将 x-expires 参数设置为 queue.declare 或通过设置过期策略来为给定队列设置过期时间。 这控制队列在自动删除之前可以闲置多长时间。 未使用意味着队列没有消费者，队列最近没有被重新声明（重新声明续租），并且 basic.get 在至少到期期限内没有被调用。 例如，这可以用于 RPC 样式的回复队列，其中可以创建许多可能永远不会耗尽的队列。

The server guarantees that the queue will be deleted, if unused for at least the expiration period. No guarantee is given as to how promptly the queue will be removed after the expiration period has elapsed.
服务器保证队列将被删除，如果至少在有效期内未使用。 无法保证在到期期限过后队列将以多快的速度被删除。

The value of the x-expires argument or expires policy describes the expiration period in milliseconds. It must be a positive integer (unlike message TTL it cannot be 0). Thus a value of 1000 means a queue which is unused for 1 second will be deleted.
x-expires 参数或过期策略的值以毫秒为单位描述过期时间。 它必须是正整数（与消息 TTL 不同，它不能为 0）。 因此，值为 1000 意味着将删除 1 秒未使用的队列。

### Define Queue TTL for Queues Using a Policy

The following policy makes all queues expire after 30 minutes since last use:

|                       |                                                                                  |
| --------------------- | -------------------------------------------------------------------------------- |
| rabbitmqctl           | rabbitmqctl set_policy expiry ".*" '{"expires":1800000}' --apply-to queues       |
| rabbitmqctl (Windows) | rabbitmqctl.bat set_policy expiry ".*" "{""expires"":1800000}" --apply-to queues |

### Define Queue TTL for Queues Using x-arguments During Declaration

This example in Java creates a queue which expires after it has been unused for 30 minutes.

```
Map<String, Object> args = new HashMap<String, Object>();
args.put("x-expires", 1800000);
channel.queueDeclare("myqueue", false, false, false, args);
```

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!
