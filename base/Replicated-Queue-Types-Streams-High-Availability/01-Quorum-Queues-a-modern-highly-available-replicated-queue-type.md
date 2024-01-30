# Quorum Queues

* https://www.rabbitmq.com/quorum-queues.html

## Overview

The RabbitMQ quorum queue is a modern queue type, which implements a durable, replicated FIFO queue based on the [Raft consensus algorithm](https://raft.github.io/).
RabbitMQ 仲裁队列是一种现代队列类型，它基于 Raft 共识算法实现了持久的、可复制的 FIFO 队列。

Quorum queues are designed to be safer and provide simpler, well defined failure handling semantics that users should find easier to reason about when designing and operating their systems.
仲裁队列的设计更加安全，并提供更简单、定义明确的故障处理语义，用户在设计和操作系统时应该更容易推理这些语义。

Quorum queues and [streams](https://www.rabbitmq.com/streams.html) now replace the original, replicated [mirrored classic queue](https://www.rabbitmq.com/ha.html). Mirrored classic queues are [now deprecated and scheduled for removal](https://blog.rabbitmq.com/posts/2021/08/4.0-deprecation-announcements/). Use the [Migrate your RabbitMQ Mirrored Classic Queues to Quorum Queues](https://www.rabbitmq.com/migrate-mcq-to-qq.html) guide for migrating RabbitMQ installations that currently use classic mirrored queues.
仲裁队列和流现在取代了原始的复制镜像经典队列。 镜像经典队列现已弃用并计划删除。 使用将 RabbitMQ 镜像经典队列迁移到仲裁队列指南来迁移当前使用经典镜像队列的 RabbitMQ 安装。

Quorum queues are optimized for [set of use cases](https://www.rabbitmq.com/quorum-queues.html#use-cases) where [data safety](https://www.rabbitmq.com/quorum-queues.html#data-safety) is a top priority. This is covered in [Motivation](https://www.rabbitmq.com/quorum-queues.html#motivation). Quorum queues should be considered the default option for a replicated queue type.
仲裁队列针对数据安全是重中之重的一组用例进行了优化。 这在动机中有所涉及。 仲裁队列应被视为复制队列类型的默认选项。

Quorum queues also have important [differences in behaviour](https://www.rabbitmq.com/quorum-queues.html#behaviour) and some [limitations](https://www.rabbitmq.com/quorum-queues.html#feature-comparison) compared to classic mirrored queues, including workload-specific ones, e.g. when consumers [repeatedly requeue the same message](https://www.rabbitmq.com/quorum-queues.html#repeated-requeues).
与经典的镜像队列相比，仲裁队列在行为上也有重要的差异和一些限制，包括特定于工作负载的队列，例如 当消费者重复重新排队相同的消息时。

Some features, such as [poison message handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling), are specific to quorum queues.
某些功能（例如有害消息处理）特定于仲裁队列。

For cases that would benefit from replication and repeatable reads, [streams](https://www.rabbitmq.com/streams.html) may be a better option than quorum queues.
对于受益于复制和可重复读取的情况，流可能是比仲裁队列更好的选择。

## Topics Covered

Topics covered in this information include:
此信息涵盖的主题包括：

- [What are quorum queues](https://www.rabbitmq.com/quorum-queues.html#motivation) and why they were introduced
  什么是仲裁队列以及引入它们的原因

- [How are they different](https://www.rabbitmq.com/quorum-queues.html#feature-comparison) from classic queues
  它们与经典队列有何不同

- Primary [use cases](https://www.rabbitmq.com/quorum-queues.html#use-cases) of quorum queues and when not to use them
  仲裁队列的主要用例以及何时不使用它们

- How to [declare a quorum queue](https://www.rabbitmq.com/quorum-queues.html#usage)
  如何声明仲裁队列

- [Replication](https://www.rabbitmq.com/quorum-queues.html#replication)-related topics: [replica management](https://www.rabbitmq.com/quorum-queues.html#replica-management), [replica leader rebalancing](https://www.rabbitmq.com/quorum-queues.html#replica-rebalancing), optimal number of replicas, etc
  复制相关主题：副本管理、副本领导者重新平衡、最佳副本数量等

- What guarantees quorum queues offer in terms of [leader failure handling](https://www.rabbitmq.com/quorum-queues.html#leader-election), [data safety](https://www.rabbitmq.com/quorum-queues.html#data-safety) and [availability](https://www.rabbitmq.com/quorum-queues.html#availability)
  在领导者故障处理、数据安全性和可用性方面保证仲裁队列提供什么

- [Performance](https://www.rabbitmq.com/quorum-queues.html#performance) characteristics of quorum queues and [performance tuning](https://www.rabbitmq.com/quorum-queues.html#performance-tuning) relevant to them
  仲裁队列的性能特征以及与其相关的性能调优

- [Poison message handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling) provided by quorum queues
  仲裁队列提供的有害消息处理

- [Configurable settings](https://www.rabbitmq.com/quorum-queues.html#configuration) of quorum queues
  仲裁队列的可配置设置

- Resource use of quorum queues, most importantly their [memory footprint](https://www.rabbitmq.com/quorum-queues.html#resource-use)
  仲裁队列的资源使用，最重要的是它们的内存占用

and more.

General familiarity with [RabbitMQ clustering](https://www.rabbitmq.com/clustering.html) would be helpful here when learning more about quorum queues.
在了解有关仲裁队列的更多信息时，对 RabbitMQ 集群的总体熟悉会很有帮助。

## Motivation
动机

Quorum queues adopt a different replication and consensus protocol and give up support for certain "transient" in nature features, which results in some limitations. These limitations are covered later in this information.
仲裁队列采用不同的复制和共识协议，并放弃对某些“瞬态”本质特征的支持，这导致了一些限制。 本信息稍后将介绍这些限制。

Quorum queues pass a [refactored and more demanding version](https://github.com/rabbitmq/jepsen#jepsen-tests-for-rabbitmq) of the [original Jepsen test](https://aphyr.com/posts/315-jepsen-rabbitmq#rabbit-as-a-queue). This ensures they behave as expected under network partitions and failure scenarios. The new test runs continuously to spot possible regressions and is enhanced regularly to test new features (e.g. [dead lettering](https://www.rabbitmq.com/quorum-queues.html#dead-lettering)).
仲裁队列通过了原始 Jepsen 测试的重构且要求更高的版本。 这可确保它们在网络分区和故障情况下按预期运行。 新测试持续运行以发现可能的回归，并定期增强以测试新功能（例如死字）。

### What is a Quorum?
什么是法定人数？

If intentionally simplified, [quorum](https://en.wikipedia.org/wiki/Quorum_(distributed_computing)) in a distributed system can be defined as an agreement between the majority of nodes ((N/2)+1 where N is the total number of system participants).
如果有意简化，分布式系统中的仲裁可以定义为大多数节点之间的协议（(N/2)+1，其中 N 是系统参与者的总数）。

When applied to queue mirroring in RabbitMQ [clusters](https://www.rabbitmq.com/clustering.html) this means that the majority of replicas (including the currently elected queue leader) agree on the state of the queue and its contents.
当应用于 RabbitMQ 集群中的队列镜像时，这意味着大多数副本（包括当前选举的队列领导者）都同意队列的状态及其内容。

### Differences between Quorum Queues and Classic Mirrored Queues
仲裁队列和经典镜像队列之间的差异

Quorum queues share many of the fundamentals with [queues](https://www.rabbitmq.com/queues.html) of other types in RabbitMQ. However, they are more purpose-built, focus on data safety and predictable recovery, and do not support certain features.
仲裁队列与 RabbitMQ 中其他类型的队列共享许多基本原理。 然而，它们更具针对性，注重数据安全和可预测恢复，并且不支持某些功能。

The differences [are covered](https://www.rabbitmq.com/quorum-queues.html#feature-comparison) in this guide.
本指南介绍了这些差异。

Classic mirrored queues in RabbitMQ have technical limitations that makes it difficult to provide comprehensible guarantees and clear failure handling semantics.
RabbitMQ 中的经典镜像队列具有技术限制，因此很难提供可理解的保证和清晰的故障处理语义。

Certain failure scenarios can result in mirrored queues confirming messages too early, potentially resulting in a data loss.
某些故障情况可能会导致镜像队列过早确认消息，从而可能导致数据丢失。

## Feature Comparison with Regular Queues
与常规队列的功能比较

Quorum queues share most of the fundamentals with other [queue](https://www.rabbitmq.com/queues.html) types. A client library that can use regular mirrored queues will be able to use quorum queues.
仲裁队列与其他队列类型共享大部分基本原理。 可以使用常规镜像队列的客户端库将能够使用仲裁队列。

The following operations work the same way for quorum queues as they do for regular queues:
以下操作对于仲裁队列的工作方式与对于常规队列的工作方式相同：

- Consumption (subscription)
  消费（订阅）

- [Consumer acknowledgements](https://www.rabbitmq.com/confirms.html) (except for global [QoS and prefetch](https://www.rabbitmq.com/quorum-queues.html#global-qos))
  消费者确认（全局 QoS 和预取除外）

- Cancelling consumers
  取消消费者

- Purging
  净化

- Deletion
  删除

With some queue operations there are minor differences:
某些队列操作存在细微差别：

- [Declaration](https://www.rabbitmq.com/quorum-queues.html#declaring)

- Setting prefetch for consumers
  为消费者设置预取

Some features are not currently supported by quorum queues.
仲裁队列当前不支持某些功能。

### Feature Matrix

| Feature | Classic Mirrored | Quorum |
| ------- | ---------------- | ------ |
| [Non-durable queues](https://www.rabbitmq.com/queues.html) | yes | no |
| [Exclusivity](https://www.rabbitmq.com/queues.html) | yes | no |
| Per message persistence | per message | always |
| Membership changes | automatic | manual |
| [Message TTL (Time-To-Live)](https://www.rabbitmq.com/ttl.html) | yes | yes ([since 3.10](https://blog.rabbitmq.com/posts/2022/05/rabbitmq-3.10-release-overview/)) |
| [Queue TTL](https://www.rabbitmq.com/ttl.html#queue-ttl) | yes | partially (lease is not renewed on queue re-declaration) |
| [Queue length limits](https://www.rabbitmq.com/maxlength.html) | yes | yes (except x-overflow: reject-publish-dlx) |
| [Lazy behaviour](https://www.rabbitmq.com/lazy-queues.html) | yes | always (since 3.10) |
| [Message priority](https://www.rabbitmq.com/priority.html) | yes | no |
| [Consumer priority](https://www.rabbitmq.com/consumer-priority.html) | yes | yes |
| [Dead letter exchanges](https://www.rabbitmq.com/dlx.html) | yes | yes |
| Adheres to [policies](https://www.rabbitmq.com/parameters.html#policies) | yes | yes (see [Policy support](https://www.rabbitmq.com/quorum-queues.html#policy-support)) |
| Poison message handling | no | yes |
| Global [QoS Prefetch](https://www.rabbitmq.com/quorum-queues.html#global-qos) | yes | no |

Modern quorum queues also offer [higher throughput and less latency variability](https://blog.rabbitmq.com/posts/2022/05/rabbitmq-3.10-performance-improvements/) for many workloads.
现代仲裁队列还为许多工作负载提供更高的吞吐量和更少的延迟变化。

#### Non-durable Queues
非持久队列

Classic queues can be [non-durable](https://www.rabbitmq.com/queues.html). Quorum queues are always durable per their assumed [use cases](https://www.rabbitmq.com/quorum-queues.html#use-cases).
经典队列可能是非持久的。 根据其假设的用例，仲裁队列始终是持久的。

#### Exclusivity
排他性

[Exclusive queues](https://www.rabbitmq.com/queues.html#exclusive-queues) are tied to the lifecycle of their declaring connection. Quorum queues by design are replicated and durable, therefore the exclusive property makes no sense in their context. Therefore quorum queues cannot be exclusive.
独占队列与其声明连接的生命周期相关。 仲裁队列的设计是可复制且持久的，因此独占属性在其上下文中没有任何意义。 因此仲裁队列不能是排他的。

Quorum queues are not meant to be used as [temporary queues](https://www.rabbitmq.com/queues.html#temporary-queues).
仲裁队列并不意味着用作临时队列。

#### Queue and Per-Message TTL (since RabbitMQ 3.10)
队列和每条消息的 TTL（自 RabbitMQ 3.10 起）

Quorum queues support both [Queue TTL](https://www.rabbitmq.com/ttl.html#queue-ttl) and message TTL (including [Per-Queue Message TTL in Queues](https://www.rabbitmq.com/ttl.html#per-queue-message-ttl) and [Per-Message TTL in Publishers](https://www.rabbitmq.com/ttl.html#per-message-ttl-in-publishers)). When using any form of message TTL, the memory overhead increases by 2 bytes per message.
仲裁队列支持队列 TTL 和消息 TTL（包括队列中的每队列消息 TTL 和发布者中的每消息 TTL）。 当使用任何形式的消息 TTL 时，每条消息的内存开销都会增加 2 个字节。

#### Length Limit
长度限制

Quorum queues has support for [queue length limits](https://www.rabbitmq.com/maxlength.html).
仲裁队列支持队列长度限制。

The drop-head and reject-publish overflow behaviours are supported but they do not support reject-publish-dlx configurations as Quorum queues take a different implementation approach than classic queues.
支持 drop-head 和拒绝发布溢出行为，但它们不支持拒绝发布 dlx 配置，因为仲裁队列采用与经典队列不同的实现方法。

The current implementation of reject-publish overflow behaviour does not strictly enforce the limit and allows a quorum queue to overshoot its limit by at least one message, therefore it should be taken with care in scenarios where a precise limit is required.
目前拒绝发布溢出行为的实现并没有严格执行限制，并且允许仲裁队列超出其限制至少一条消息，因此在需要精确限制的场景中应小心谨慎。

When a quorum queue reaches the max-length limit and reject-publish is configured it notifies each publishing channel who from thereon will reject all messages back to the client. This means that quorum queues may overshoot their limit by some small number of messages as there may be messages in flight whilst the channels are notified. The number of additional messages that are accepted by the queue will vary depending on how many messages are in flight at the time.
当仲裁队列达到最大长度限制并且配置了拒绝发布时，它会通知每个发布通道，该通道将拒绝所有返回客户端的消息。 这意味着仲裁队列可能会因少量消息而超出其限制，因为在通知通道时可能有消息正在传输。 队列接受的附加消息数量将根据当时正在传输的消息数量而变化。

#### Dead Lettering
死字

Quorum queues support [dead letter exchanges](https://www.rabbitmq.com/dlx.html) (DLXs).
仲裁队列支持死信交换 (DLX)。

Traditionally, using DLXs in a clustered environment has not been [safe](https://www.rabbitmq.com/dlx.html#safety).
传统上，在集群环境中使用 DLX 并不安全。

Since RabbitMQ 3.10 quorum queues support a safer form of dead-lettering that uses at-least-once guarantees for the message transfer between queues (with the limitations and caveats outlined below).
自 RabbitMQ 3.10 仲裁队列支持更安全的死信形式，它使用至少一次保证队列之间的消息传输（具有下面列出的限制和警告）。

This is done by implementing a special, internal dead-letter consumer process that works similarly to a normal queue consumer with manual acknowledgements apart from it only consumes messages that have been dead-lettered.
这是通过实现一个特殊的内部死信消费者进程来完成的，该进程的工作方式与具有手动确认的普通队列消费者类似，但它仅消费已死信的消息。

This means that the source quorum queue will retain the dead-lettered messages until they have been acknowledged. The internal consumer will consume dead-lettered messages and publish them to the target queue(s) using publisher confirms. It will only acknowledge once publisher confirms have been received, hence providing at-least-once guarantees.
这意味着源仲裁队列将保留死信消息，直到它们被确认为止。 内部消费者将使用死信消息并使用发布者确认将它们发布到目标队列。 它只会在收到发布者确认后才会确认，因此提供至少一次保证。

at-most-once remains the default dead-letter-strategy for quorum queues and is useful for scenarios where the dead lettered messages are more of an informational nature and where it does not matter so much if they are lost in transit between queues or when the overflow configuration restriction outlined below is not suitable.
至多一次仍然是仲裁队列的默认死信策略，对于死信消息更多地具有信息性质的情况非常有用，并且如果它们在队列之间传输或在以下情况下丢失并不那么重要。 下面概述的溢出配置限制是不合适的。

##### Activating at-least-once dead-lettering
激活至少一次死信

To activate or turn on at-least-once dead-lettering for a source quorum queue, apply all of the following policies (or the equivalent queue arguments starting with x-):
要激活或打开源仲裁队列的至少一次死信，请应用以下所有策略（或以 x- 开头的等效队列参数）：

- Set dead-letter-strategy to at-least-once (default is at-most-once).
  将死信策略设置为至少一次（默认为最多一次）。

- Set overflow to reject-publish (default is drop-head).
  将溢出设置为拒绝发布（默认为 drop-head）。

- Configure a dead-letter-exchange.
  配置死信交换。

- Turn on [feature flag](https://www.rabbitmq.com/feature-flags.html) stream_queue (turned on by default for RabbitMQ clusters created in 3.9 or later).
  打开功能标志stream_queue（对于在3.9或更高版本中创建的RabbitMQ集群默认打开）。

It is recommended to additionally configure max-length or max-length-bytes to prevent excessive message buildup in the source quorum queue (see caveats below).
建议另外配置 max-length 或 max-length-bytes 以防止源仲裁队列中的消息堆积过多（请参阅下面的注意事项）。

Optionally, configure a dead-letter-routing-key.
（可选）配置死信路由密钥。

##### Limitations
局限性

at-least-once dead lettering does not work with the default drop-head overflow strategy even if a queue length limit is not set. Hence if drop-head is configured the dead-lettering will fall back to at-most-once. Use the overflow strategy reject-publish instead.
即使未设置队列长度限制，至少一次死信也不适用于默认的丢弃头溢出策略。 因此，如果配置了 drop-head，死信将回退到最多一次。 请改用溢出策略拒绝发布。

##### Caveats
注意事项

at-least-once dead-lettering will require more system resources such as memory and CPU. Therefore, turn on at-least-once only if dead lettered messages should not be lost.
至少一次死信将需要更多的系统资源，例如内存和 CPU。 因此，只有在死信消息不应该丢失的情况下才应至少打开一次。

at-least-once guarantees opens up some specific failure cases that needs handling. As dead-lettered messages are now retained by the source quorum queue until they have been safely accepted by the dead-letter target queue(s) this means they have to contribute to the queue resource limits, such as max length limits so that the queue can refuse to accept more messages until some have been removed. Theoretically it is then possible for a queue to *only* contain dead-lettered messages, in the case where, say a target dead-letter queue isn't available to accept messages for a long time and normal queue consumers consume most of the messages.
至少一次保证会引发一些需要处理的特定故障情况。 由于死信消息现在由源仲裁队列保留，直到它们被死信目标队列安全地接受为止，这意味着它们必须对队列资源限制做出贡献，例如最大长度限制，以便队列 可以拒绝接受更多消息，直到某些消息被删除。 理论上，队列有可能仅包含死信消息，例如目标死信队列长时间无法接受消息并且普通队列使用者消耗了大部分消息的情况。

Dead-lettered messages are considered "live" until they have been confirmed by the dead-letter target queue(s).
死信消息被视为“活动”消息，直到它们被死信目标队列确认为止。

There are few cases for which dead lettered messages will not be removed from the source queue in a timely manner:
在极少数情况下，死信消息不会被及时从源队列中删除：

- The configured dead-letter exchange does not exist.
  配置的死信交换不存在。

- The messages cannot be routed to any queue (equivalent to the mandatory message property).
  消息不能路由到任何队列（相当于强制消息属性）。

- One (of possibly many) routed target queues does not confirm receipt of the message. This can happen when a target queue is not available or when a target queue rejects a message (e.g. due to exceeded queue length limit).
  一个（可能是多个）路由目标队列不确认消息的接收。 当目标队列不可用或目标队列拒绝消息时（例如，由于超出队列长度限制），可能会发生这种情况。

The dead-letter consumer process will retry periodically if either of the scenarios above occur which means there is a possibility of duplicates appearing at the DLX target queue(s).
如果发生上述任一情况，死信使用者进程将定期重试，这意味着 DLX 目标队列中可能会出现重复项。

For each quorum queue with at-least-once dead-lettering turned on, there will be one internal dead-letter consumer process. The internal dead-letter consumer process is co-located on the quorum queue leader node. It keeps all dead-lettered message bodies in memory. It uses a prefetch size of 32 messages to limit the amount of message bodies kept in memory if no confirms are received from the target queues.
对于每个启用至少一次死信的仲裁队列，将有一个内部死信使用者进程。 内部死信使用者进程位于仲裁队列领导节点上。 它将所有死信消息体保留在内存中。 如果没有从目标队列收到确认，它会使用 32 条消息的预取大小来限制内存中保存的消息体数量。

That prefetch size can be increased by the dead_letter_worker_consumer_prefetch setting in the rabbit app section of the [advanced config file](https://www.rabbitmq.com/configure.html#advanced-config-file) if high dead-lettering throughput (thousands of messages per second) is required.
如果需要高死信吞吐量（每秒数千条消息），则可以通过高级配置文件的rabbit应用程序部分中的dead_letter_worker_consumer_prefetch设置来增加预取大小。

For a source quorum queue, it is possible to switch dead-letter strategy dynamically from at-most-once to at-least-once and vice versa. If the dead-letter strategy is changed either directly from at-least-once to at-most-once or indirectly, for example by changing overflow from reject-publish to drop-head, any dead-lettered messages that have not yet been confirmed by all target queues will be deleted.
对于源仲裁队列，可以动态地将死信策略从最多一次切换到至少一次，反之亦然。 如果死信策略直接从至少一次更改为最多一次或间接更改，例如通过将溢出从拒绝发布更改为丢弃头，则任何尚未确认的死信消息 所有目标队列将被删除。

Messages published to the source quorum queue are persisted on disk regardless of the message delivery mode (transient or persistent). However, messages that are dead lettered by the source quorum queue will keep the original message delivery mode. This means if dead lettered messages in the target queue should survive a broker restart, the target queue must be durable and the message delivery mode must be set to persistent when publishing messages to the source quorum queue.
无论消息传递模式如何（瞬时或持久），发布到源仲裁队列的消息都会保留在磁盘上。 但是，源仲裁队列中的死信消息将保留原始消息传递模式。 这意味着，如果目标队列中的死信消息在代理重新启动后仍然存在，则目标队列必须是持久的，并且在将消息发布到源仲裁队列时，消息传递模式必须设置为持久。

#### Lazy Mode (since RabbitMQ 3.10)
惰性模式（自 RabbitMQ 3.10 起）

Quorum queues store their message content on disk (per Raft requirements) and only keep a small metadata record of each message in memory. This is a change from prior versions of quorum queues where there was an option to keep the message bodies in memory as well. This never proved to be beneficial especially when the queue length was large.
仲裁队列将其消息内容存储在磁盘上（根据 Raft 要求），并且仅在内存中保留每条消息的一小部分元数据记录。 这是与之前版本的仲裁队列的一个变化，之前版本可以选择将消息正文保留在内存中。 这从未被证明是有益的，尤其是当队列长度很大时。

The [memory limit](https://www.rabbitmq.com/quorum-queues.html#memory-limit) configuration is still permitted but has no effect. The only option now is effectively the same as configuring: x-max-in-memory-length=0
内存限制配置仍然允许，但没有效果。 现在唯一的选项实际上与配置相同：x-max-in-memory-length=0

The [lazy mode configuration](https://www.rabbitmq.com/lazy-queues.html#configuration) does not apply.
惰性模式配置不适用。

#### Lazy Mode (before RabbitMQ 3.10)
惰性模式（RabbitMQ 3.10 之前）

Quorum queues store their content on disk (per Raft requirements) as well as in memory (up to the [in memory limit configured](https://www.rabbitmq.com/quorum-queues.html#memory-limit)).
仲裁队列将其内容存储在磁盘上（根据 Raft 要求）以及内存中（最多达到配置的内存限制）。

The [lazy mode configuration](https://www.rabbitmq.com/lazy-queues.html#configuration) does not apply.
惰性模式配置不适用。

It is possible to [limit how many messages a quorum queue keeps in memory](https://www.rabbitmq.com/quorum-queues.html#memory-limit) using a policy which can achieve a behaviour similar to lazy queues.
可以使用可以实现类似于惰性队列的行为的策略来限制仲裁队列在内存中保留的消息数量。

#### Global QoS
全局服务质量

Quorum queues do not support global [QoS prefetch](https://www.rabbitmq.com/confirms.html#channel-qos-prefetch) where a channel sets a single prefetch limit for all consumers using that channel. If an attempt is made to consume from a quorum queue from a channel with global QoS activated a channel error will be returned.
仲裁队列不支持全局 QoS 预取，其中通道为使用该通道的所有使用者设置单个预取限制。 如果尝试从激活了全局 QoS 的通道的仲裁队列中进行消费，将返回通道错误。

Use [per-consumer QoS prefetch](https://www.rabbitmq.com/consumer-prefetch.html), which is the default in several popular clients.
使用每个消费者的 QoS 预取，这是几个流行客户端的默认设置。

#### Priorities

Quorum queues support [consumer priorities](https://www.rabbitmq.com/consumer-priority.html), but not [message priorities](https://www.rabbitmq.com/priority.html).
仲裁队列支持消费者优先级，但不支持消息优先级。

To prioritize messages with Quorum Queues, use multiple queues; one for each priority.
要使用仲裁队列对消息进行优先级排序，请使用多个队列； 每个优先级都有一个。

#### Poison Message Handling
有毒消息处理

Quorum queues [support poison message handling](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling) via a redelivery limit. This feature is currently unique to quorum queues.
仲裁队列通过重新传递限制支持有害消息处理。 此功能目前是仲裁队列所独有的。

#### Policy Support

Quorum queues can be configured via RabbitMQ policies. The below table summarises the policy keys they adhere to.
仲裁队列可以通过 RabbitMQ 策略进行配置。 下表总结了他们遵守的政策要点。

| Definition Key          | Type                            |
| ----------------------- | ------------------------------- |
| max-length              | Number                          |
| max-length-bytes        | Number                          |
| overflow                | "drop-head" or "reject-publish" |
| expires                 | Number (milliseconds)           |
| dead-letter-exchange    | String                          |
| dead-letter-routing-key | String                          |
| max-in-memory-length    | Number                          |
| max-in-memory-bytes     | Number                          |
| delivery-limit          | Number                          |

## Use Cases

Quorum queues are purpose built by design. They are *not* designed to be used for every problem. Their intended use is for topologies where queues exist for a long time and are critical to certain aspects of system operation, therefore fault tolerance and data safety is more important than, say, lowest possible latency and advanced queue features.
仲裁队列是专门设计的。 它们并非旨在用于解决所有问题。 它们的预期用途是队列长期存在并且对系统操作的某些方面至关重要的拓扑，因此容错和数据安全比最低可能的延迟和高级队列功能更重要。

Examples would be incoming orders in a sales system or votes cast in an election system where potentially losing messages would have a significant impact on system correctness and function.
例如，销售系统中的传入订单或选举系统中的投票，其中潜在的消息丢失将对系统的正确性和功能产生重大影响。

Stock tickers and instant messaging systems benefit less or not at all from quorum queues.
股票行情自动收录器和即时通讯系统从仲裁队列中获益较少或根本没有获益。

Publishers should use publisher confirms as this is how clients can interact with the quorum queue consensus system. Publisher confirms will only be issued once a published message has been successfully replicated to a quorum of nodes and is considered "safe" within the context of the system.
发布者应该使用发布者确认，因为这是客户端与仲裁队列共识系统交互的方式。 仅当发布的消息已成功复制到法定数量的节点并在系统上下文中被视为“安全”时，发布者确认才会发出。

Consumers should use manual acknowledgements to ensure messages that aren't successfully processed are returned to the queue so that another consumer can re-attempt processing.
消费者应使用手动确认来确保未成功处理的消息返回到队列，以便另一个消费者可以重新尝试处理。

### When Not to Use Quorum Queues

In some cases quorum queues should not be used. They typically involve:
在某些情况下，不应使用仲裁队列。 它们通常涉及：

- Temporary nature of queues: transient or exclusive queues, high queue churn (declaration and deletion rates)
  队列的临时性：瞬态或独占队列、高队列流失率（声明和删除率）

- Lowest possible latency: the underlying consensus algorithm has an inherently higher latency due to its data safety features
  尽可能低的延迟：由于其数据安全特性，底层共识算法本身具有较高的延迟

- When data safety is not a priority (e.g. applications do not use [manual acknowledgements and publisher confirms](https://www.rabbitmq.com/confirms.html) are not used)
  当数据安全不是优先考虑的时候（例如应用程序不使用手动确认并且不使用发布者确认）

- Very long queue backlogs ([streams](https://www.rabbitmq.com/stream.html) are likely to be a better fit)
  队列积压非常长（流可能更适合）

## Usage

As stated earlier, quorum queues share most of the fundamentals with other [queue](https://www.rabbitmq.com/queues.html) types. A client library that can specify [optional queue arguments](https://www.rabbitmq.com/queues.html#optional-arguments) will be able to use quorum queues.
如前所述，仲裁队列与其他队列类型共享大部分基本原理。 可以指定可选队列参数的客户端库将能够使用仲裁队列。

First we will cover how to declare a quorum queue.
首先我们将介绍如何声明仲裁队列。

### Declaring

To declare a quorum queue set the x-queue-type queue argument to quorum (the default is classic). This argument must be provided by a client at queue declaration time; it cannot be set or changed using a [policy](https://www.rabbitmq.com/parameters.html#policies). This is because policy definition or applicable policy can be changed dynamically but queue type cannot. It must be specified at the time of declaration.
要声明仲裁队列，请将 x-queue-type 队列参数设置为仲裁（默认为 classic）。 该参数必须由客户端在队列声明时提供； 不能使用策略来设置或更改它。 这是因为策略定义或适用的策略可以动态更改，但队列类型不能。 必须在申报时注明。

Declaring a queue with an x-queue-type argument set to quorum will declare a quorum queue with up to five replicas (default [replication factor](https://www.rabbitmq.com/quorum-queues.html#replication-factor)), one per each [cluster node](https://www.rabbitmq.com/clustering.html).
声明一个将 x-queue-type 参数设置为 quorum 的队列将声明一个最多包含五个副本（默认复制因子）的仲裁队列，每个集群节点一个。

For example, a cluster of three nodes will have three replicas, one on each node. In a cluster of seven nodes, five nodes will have one replica each but two nodes won't host any replicas.
例如，一个包含三个节点的集群将具有三个副本，每个节点上一个。 在七个节点的集群中，五个节点各有一个副本，但两个节点不会托管任何副本。

After declaration a quorum queue can be bound to any exchange just as any other RabbitMQ queue.
声明后，仲裁队列可以像任何其他 RabbitMQ 队列一样绑定到任何交换。

If declaring using [management UI](https://www.rabbitmq.com/management.html), queue type must be specified using the queue type drop down menu.
如果使用管理 UI 声明，则必须使用队列类型下拉菜单指定队列类型。

### Client Operations for Quorum Queues

The following operations work the same way for quorum queues as they do for classic queues:
以下操作对于仲裁队列的工作方式与对于经典队列的工作方式相同：

- [Consumption](https://www.rabbitmq.com/consumers.html) (subscription)
  消费（订阅）

- [Consumer acknowledgements](https://www.rabbitmq.com/confirms.html) (keep [QoS Prefetch Limitations](https://www.rabbitmq.com/quorum-queues.html#global-qos) in mind)
  消费者确认（牢记 QoS 预取限制）

- Cancellation of consumers
  消费者取消

- Purging of queue messages
  清除队列消息

- Queue deletion
  队列删除

With some queue operations there are minor differences:
某些队列操作存在细微差别：

- [Declaration](https://www.rabbitmq.com/quorum-queues.html#declaring) (covered above)
  声明（如上所述）

- Setting [QoS prefetch](https://www.rabbitmq.com/quorum-queues.html#global-qos) for consumers
  为消费者设置 QoS 预取

## Quorum Queue Replication and Data Locality
仲裁队列复制和数据局部性

When a quorum queue is declared, an initial number of replicas for it must be started in the cluster. By default the number of replicas to be started is up to three, one per RabbitMQ node in the cluster.
当声明仲裁队列时，必须在集群中启动其初始数量的副本。 默认情况下，要启动的副本数量最多为三个，集群中每个 RabbitMQ 节点一个。

Three nodes is the **practical minimum** of replicas for a quorum queue. In RabbitMQ clusters with a larger number of nodes, adding more replicas than a [quorum](https://www.rabbitmq.com/quorum-queues.html#what-is-quorum) (majority) will not provide any improvements in terms of [quorum queue availability](https://www.rabbitmq.com/quorum-queues.html#quorum-requirements) but it will consume more cluster resources.
三个节点是仲裁队列的实际最少副本数。 在具有大量节点的 RabbitMQ 集群中，添加比仲裁（多数）更多的副本不会在仲裁队列可用性方面提供任何改进，但会消耗更多集群资源。

Therefore the **recommended number of replicas** for a quorum queue is the quorum of cluster nodes (but no fewer than three). This assumes a [fully formed](https://www.rabbitmq.com/cluster-formation.html) cluster of at least three nodes.
因此，仲裁队列的建议副本数量是集群节点的仲裁数量（但不少于三个）。 这假设一个完整的集群至少包含三个节点。

### Controlling the Initial Replication Factor
控制初始复制因子

For example, a cluster of three nodes will have three replicas, one on each node. In a cluster of seven nodes, three nodes will have one replica each but four more nodes won't host any replicas of the newly declared queue.
例如，一个包含三个节点的集群将具有三个副本，每个节点上一个。 在由七个节点组成的集群中，三个节点各有一个副本，但另外四个节点将不会托管新声明的队列的任何副本。

Like with classic mirrored queues, the replication factor (number of replicas a queue has) can be configured for quorum queues.
与经典镜像队列一样，可以为仲裁队列配置复制因子（队列具有的副本数）。

The minimum factor value that makes practical sense is three. It is highly recommended for the factor to be an odd number. This way a clear quorum (majority) of nodes can be computed. For example, there is no "majority" of nodes in a two node cluster. This is covered with more examples below in the [Fault Tolerance and Minimum Number of Replicas Online](https://www.rabbitmq.com/quorum-queues.html#quorum-requirements) section.
具有实际意义的最小因子值为三。 强烈建议该因子为奇数。 这样就可以计算出明确的节点法定人数（多数）。 例如，在两节点集群中不存在“大多数”节点。 下面的容错和最小在线副本数部分中的更多示例对此进行了介绍。

This may not be desirable for larger clusters or for cluster with an even number of nodes. To control the number of quorum queue members set the x-quorum-initial-group-size queue argument when declaring the queue. The group size argument provided should be an integer that is greater than zero and smaller or equal to the current RabbitMQ cluster size. The quorum queue will be launched to run on a random subset of RabbitMQ nodes present in the cluster at declaration time.
对于较大的集群或具有偶数个节点的集群来说，这可能是不可取的。 要控制仲裁队列成员的数量，请在声明队列时设置 x-quorum-initial-group-size 队列参数。 提供的组大小参数应该是大于零且小于或等于当前 RabbitMQ 集群大小的整数。 仲裁队列将在声明时启动，在集群中存在的 RabbitMQ 节点的随机子集上运行。

In case a quorum queue is declared before all cluster nodes have joined the cluster, and the initial replica count is greater than the total number of cluster members, the effective value used will be equal to the total number of cluster nodes. When more nodes join the cluster, the replica count will not be automatically increased but it can be [increased by the operator](https://www.rabbitmq.com/quorum-queues.html#replica-management).
如果在所有集群节点加入集群之前声明仲裁队列，并且初始副本计数大于集群成员总数，则使用的有效值将等于集群节点总数。 当更多节点加入集群时，副本数不会自动增加，但可以由操作员增加。

### Queue Leader Location
队列领导位置

Every quorum queue has a primary replica. That replica is called *queue leader*. All queue operations go through the leader first and then are replicated to followers (mirrors). This is necessary to guarantee FIFO ordering of messages.
每个仲裁队列都有一个主副本。 该副本称为队列领导者。 所有队列操作首先经过领导者，然后复制到跟随者（镜像）。 这对于保证消息的 FIFO 排序是必要的。

To avoid some nodes in a cluster hosting the majority of queue leader replicas and thus handling most of the load, queue leaders should be reasonably evenly distributed across cluster nodes.
为了避免集群中的某些节点托管大多数队列领导者副本并因此处理大部分负载，队列领导者应合理均匀地分布在集群节点上。

When a new quorum queue is declared, the set of nodes that will host its replicas is randomly picked, but will always include the node the client that declares the queue is connected to.
当声明一个新的仲裁队列时，将随机选择将托管其副本的节点集，但将始终包括声明该队列连接到的客户端的节点。

Which replica becomes the initial leader can controlled using three options:
哪个副本成为初始领导者可以使用三个选项进行控制：

1. Setting the queue-leader-locator [policy](https://www.rabbitmq.com/parameters.html#policies) key (recommended)
   设置queue-leader-locator策略键（推荐）

2. By defining the queue_leader_locator key in [the configuration file](https://www.rabbitmq.com/configure.html#configuration-files) (recommended)
   通过在配置文件中定义queue_leader_locator键（推荐）

3. Using the x-queue-leader-locator [optional queue argument](https://www.rabbitmq.com/queues.html#optional-arguments)
   使用 x-queue-leader-locator 可选队列参数

Supported queue leader locator values are
支持的队列领导者定位符值是

- client-local: Pick the node the client that declares the queue is connected to. This is the default value.
  client-local：选择声明队列连接到的客户端的节点。 这是默认值。

- balanced: If there are overall less than 1000 queues (classic queues, quorum queues, and streams), pick the node hosting the minimum number of quorum queue leaders. If there are overall more than 1000 queues, pick a random node.
  平衡：如果队列（经典队列、仲裁队列和流）总数少于 1000 个，则选择托管最少数量仲裁队列领导者的节点。 如果队列总数超过 1000 个，则随机选择一个节点。

### Managing Replicas(Quorum Group Members)
管理副本（仲裁组成员）

Replicas of a quorum queue are explicitly managed by the operator. When a new node is added to the cluster, it will host no quorum queue replicas unless the operator explicitly adds it to a member (replica) list of a quorum queue or a set of quorum queues.
仲裁队列的副本由操作员显式管理。 将新节点添加到集群时，它将不会托管任何仲裁队列副本，除非操作员显式将其添加到仲裁队列或一组仲裁队列的成员（副本）列表中。

When a node has to be decommissioned (permanently removed from the cluster), it must be explicitly removed from the member list of all quorum queues it currently hosts replicas for.
当节点必须停用（从集群中永久删除）时，必须将其从当前为其托管副本的所有仲裁队列的成员列表中显式删除。

Several [CLI commands](https://www.rabbitmq.com/cli.html) are provided to perform the above operations:
提供了几个 CLI 命令来执行上述操作：

```bash
rabbitmq-queues add_member [-p <vhost>] <queue-name> <node>

rabbitmq-queues delete_member [-p <vhost>] <queue-name> <node>

rabbitmq-queues grow <node> <all | even> [--vhost-pattern <pattern>] [--queue-pattern <pattern>]

rabbitmq-queues shrink <node> [--errors-only]
```

To successfully add and remove members a quorum of replicas in the cluster must be available because cluster membership changes are treated as queue state changes.
要成功添加和删除成员，集群中的副本仲裁必须可用，因为集群成员资格更改被视为队列状态更改。

Care needs to be taken not to accidentally make a queue unavailable by losing the quorum whilst performing maintenance operations that involve membership changes.
需要注意的是，在执行涉及成员资格更改的维护操作时，不要意外地因失去法定人数而导致队列不可用。

When replacing a cluster node, it is safer to first add a new node and then decomission the node it replaces.
更换集群节点时，更安全的做法是先添加新节点，然后停用其所更换的节点。

### Rebalancing Replicas for Quorum Queues
重新平衡仲裁队列的副本

Once declared, the RabbitMQ quorum queue leaders may be unevenly distributed across the RabbitMQ cluster. To re-balance use the rabbitmq-queues rebalance command. It is important to know that this does not change the nodes which the quorum queues span. To modify the membership instead see [managing replicas](https://www.rabbitmq.com/quorum-queues.html#replica-management).
一旦声明，RabbitMQ 仲裁队列领导者可能会在 RabbitMQ 集群中分布不均匀。 要重新平衡，请使用rabbitmq-queues rebalance命令。 重要的是要知道这不会改变仲裁队列所跨越的节点。 要修改成员资格，请参阅管理副本。

```bash
# rebalances all quorum queues
rabbitmq-queues rebalance quorum
```

it is possible to rebalance a subset of queues selected by name:

```bash
# rebalances a subset of quorum queues
rabbitmq-queues rebalance quorum --queue-pattern "orders.*"
```

or quorum queues in a particular set of virtual hosts:

```bash
# rebalances a subset of quorum queues
rabbitmq-queues rebalance quorum --vhost-pattern "production.*"
```

## Quorum Queue Behaviour

A quorum queue relies on a consensus protocol called Raft to ensure data consistency and safety.
仲裁队列依赖于名为 Raft 的共识协议来确保数据的一致性和安全性。

Every quorum queue has a primary replica (a *leader* in Raft parlance) and zero or more secondary replicas (called *followers*).
每个仲裁队列都有一个主副本（Raft 术语中的领导者）和零个或多个辅助副本（称为追随者）。

A leader is elected when the cluster is first formed and later if the leader becomes unavailable.
当集群首次形成时，将选举领导者；如果领导者不可用，则稍后选举领导者。

### Leader Election and Failure Handling of Quorum Queues
仲裁队列的领导者选举和失败处理

A quorum queue requires a quorum of the declared nodes to be available to function. When a RabbitMQ node hosting a quorum queue's *leader* fails or is stopped another node hosting one of that quorum queue's *follower* will be elected leader and resume operations.
仲裁队列需要声明的节点达到仲裁数量才能运行。 当托管仲裁队列的领导者的 RabbitMQ 节点发生故障或停止时，托管该仲裁队列的跟随者之一的另一个节点将被选举为领导者并恢复操作。

Failed and rejoining followers will re-synchronise ("catch up") with the leader. In contrast to classic mirrored queues, a temporary replica failure does not require a full re-synchronization from the currently elected leader. Only the delta will be transferred if a re-joining replica is behind the leader. This "catching up" process does not affect leader availability.
失败和重新加入的追随者将与领导者重新同步（“赶上”）。 与经典的镜像队列相比，临时副本故障不需要当前选出的领导者完全重新同步。 如果重新加入的副本位于领导者后面，则仅传输增量。 这种“追赶”过程不会影响领导者的可用性。

Except for the initial replica set selection, replicas must be explicitly added to a quorum queue. When a new replica is [added](https://www.rabbitmq.com/quorum-queues.html#replica-management), it will synchronise the entire queue state from the leader, similarly to classic mirrored queues.
除了初始副本集选择之外，副本必须显式添加到仲裁队列中。 当添加新的副本时，它将从领导者同步整个队列状态，类似于经典的镜像队列。

### Fault Tolerance and Minimum Number of Replicas Online
容错和最小在线副本数

Consensus systems can provide certain guarantees with regard to data safety. These guarantees do mean that certain conditions need to be met before they become relevant such as requiring a minimum of three cluster nodes to provide fault tolerance and requiring more than half of members to be available to work at all.
共识系统可以为数据安全提供一定的保障。 这些保证确实意味着在它们变得相关之前需要满足某些条件，例如要求至少三个集群节点提供容错能力，并要求一半以上的成员可以正常工作。

Failure tolerance characteristics of clusters of various size can be described in a table:
不同规模集群的容错特性可以用表格描述：

| Cluster node count | Tolerated number of node failures | Tolerant to a network partition      |
| ------------------ | --------------------------------- | ------------------------------------ |
| 1                  | 0                                 | not applicable                       |
| 2                  | 0                                 | no                                   |
| 3                  | 1                                 | yes                                  |
| 4                  | 1                                 | yes if a majority exists on one side |
| 5                  | 2                                 | yes                                  |
| 6                  | 2                                 | yes if a majority exists on one side |
| 7                  | 3                                 | yes                                  |
| 8                  | 3                                 | yes if a majority exists on one side |
| 9                  | 4                                 | yes                                  |

As the table above shows RabbitMQ clusters with fewer than three nodes do not benefit fully from the quorum queue guarantees. RabbitMQ clusters with an even number of RabbitMQ nodes do not benefit from having quorum queue members spread over all nodes. For these systems the quorum queue size should be constrained to a smaller uneven number of nodes.
如上表所示，节点少于三个的 RabbitMQ 集群无法充分受益于仲裁队列保证。 具有偶数个 RabbitMQ 节点的 RabbitMQ 集群不会因为仲裁队列成员分布在所有节点上而受益。 对于这些系统，仲裁队列大小应限制为较小的奇数节点数。

Performance tails off quite a bit for quorum queue node sizes larger than 5. We do not recommend running quorum queues on more than 7 RabbitMQ nodes. The default quorum queue size is 3 and is controllable using the x-quorum-initial-group-size [queue argument](https://www.rabbitmq.com/queues.html#optional-arguments).
对于大于 5 的仲裁队列节点大小，性能会大幅下降。我们不建议在超过 7 个 RabbitMQ 节点上运行仲裁队列。 默认仲裁队列大小为 3，并且可以使用 x-quorum-initial-group-size 队列参数进行控制。

### Data Safety provided with Quorum Queues
仲裁队列提供数据安全

Quorum queues are designed to provide data safety under network partition and failure scenarios. A message that was successfully confirmed back to the publisher using the [publisher confirms](https://www.rabbitmq.com/confirms.html) feature should not be lost as long as at least a majority of RabbitMQ nodes hosting the quorum queue are not permanently made unavailable.
仲裁队列旨在在网络分区和故障场景下提供数据安全。 只要托管仲裁队列的至少大多数 RabbitMQ 节点没有永久不可用，使用发布者确认功能成功确认回发布者的消息就不应丢失。

Generally quorum queues favours data consistency over availability.
一般来说，仲裁队列更注重数据一致性而不是可用性。

**No guarantees are provided for messages that have not been confirmed using the publisher confirm mechanism**. Such messages could be lost "mid-way", in an operating system buffer or otherwise fail to reach the queue leader.
对于尚未使用发布者确认机制确认的消息，不提供任何保证。 此类消息可能会在操作系统缓冲区中“中途”丢失，或者无法到达队列领导者。

### Quorum Queue Availability
仲裁队列可用性

A quorum queue should be able to tolerate a minority of queue members becoming unavailable with no or little effect on availability.
仲裁队列应该能够容忍少数队列成员变得不可用，而对可用性没有影响或影响很小。

Note that depending on the [partition handling strategy](https://www.rabbitmq.com/partitions.html) used RabbitMQ may restart itself during recovery and reset the node but as long as that does not happen, this availability guarantee should hold true.
请注意，根据所使用的分区处理策略，RabbitMQ 可能会在恢复期间自行重启并重置节点，但只要这种情况没有发生，这种可用性保证就应该成立。

For example, a queue with three replicas can tolerate one node failure without losing availability. A queue with five replicas can tolerate two, and so on.
例如，具有三个副本的队列可以容忍一个节点故障而不会失去可用性。 具有五个副本的队列可以容忍两个副本，依此类推。

If a quorum of nodes cannot be recovered (say if 2 out of 3 RabbitMQ nodes are permanently lost) the queue is permanently unavailable and will need to be force deleted and recreated.
如果无法恢复法定数量的节点（例如，如果 3 个 RabbitMQ 节点中有 2 个永久丢失），则队列将永久不可用，需要强制删除并重新创建。

Quorum queue follower replicas that are disconnected from the leader or participating in a leader election will ignore queue operations sent to it until they become aware of a newly elected leader. There will be warnings in the log (received unhandled msg and similar) about such events. As soon as the replica discovers a newly elected leader, it will sync the queue operation log entries it does not have from the leader, including the dropped ones. Quorum queue state will therefore remain consistent.
与领导者断开连接或参与领导者选举的仲裁队列跟随者副本将忽略发送给它的队列操作，直到它们意识到新选举的领导者。 日志中将会有关于此类事件的警告（收到未处理的消息和类似消息）。 一旦副本发现新当选的领导者，它就会从领导者那里同步它没有的队列操作日志条目，包括删除的日志条目。 因此，仲裁队列状态将保持一致。

## Quorum Queue Performance Characteristics
仲裁队列性能特征

Quorum queues are designed to trade latency for throughput and have been tested and compared against durable [classic mirrored queues](https://www.rabbitmq.com/ha.html) in 3, 5 and 7 node configurations at several message sizes.
仲裁队列旨在以延迟换取吞吐量，并且已经过测试，并与 3、5 和 7 个节点配置中多种消息大小的持久经典镜像队列进行了比较。

In scenarios using both consumer acks and publisher confirms quorum queues have been observed to have superior throughput to classic mirrored queues. For example, take a look at [these benchmarks with 3.10](https://blog.rabbitmq.com/posts/2022/05/rabbitmq-3.10-performance-improvements/) and [another with 3.12](https://blog.rabbitmq.com/posts/2023/05/rabbitmq-3.12-performance-improvements/#significant-improvements-to-quorum-queues).
在同时使用消费者确认和发布者确认的场景中，观察到仲裁队列比经典镜像队列具有更高的吞吐量。 例如，看一下这些基准测试的 3.10 和另一个 3.12。

As quorum queues persist all data to disks before doing anything it is recommended to use the fastest disks possible and certain [Performance Tuning](https://www.rabbitmq.com/quorum-queues.html#performance-tuning) settings.
由于仲裁队列在执行任何操作之前都会将所有数据保存到磁盘，因此建议使用尽可能最快的磁盘和某些性能调整设置。

Quorum queues also benefit from consumers using higher prefetch values to ensure consumers aren't starved whilst acknowledgements are flowing through the system and allowing messages to be delivered in a timely fashion.
仲裁队列还受益于消费者使用更高的预取值，以确保消费者在确认流经系统并允许及时传递消息时不会饥饿。

Due to the disk I/O-heavy nature of quorum queues, their throughput decreases as message sizes increase.
由于仲裁队列的磁盘 I/O 密集型特性，它们的吞吐量会随着消息大小的增加而降低。

Quorum queue throughput is also affected by the number of replicas. The more replicas a quorum queue has, the lower its throughput generally will be since more work has to be done to replicate data and achieve consensus.
仲裁队列吞吐量还受副本数量的影响。 仲裁队列的副本越多，其吞吐量通常越低，因为需要做更多的工作来复制数据并达成共识。

## Configurable Settings

There are a few new configuration parameters that can be tweaked using the [advanced](https://www.rabbitmq.com/configure.html#advanced-config-file) config file.
有一些新的配置参数可以使用高级配置文件进行调整。

Note that all settings related to [resource footprint](https://www.rabbitmq.com/quorum-queues.html#resource-use) are documented in a separate section.
请注意，与资源占用相关的所有设置都记录在单独的部分中。

The ra application (which is the Raft library that quorum queues use) has [its own set of tunable parameters](https://github.com/rabbitmq/ra#configuration).
ra 应用程序（仲裁队列使用的 Raft 库）有自己的一组可调参数。

The rabbit application has several quorum queue related configuration items available.
兔子应用程序有几个与仲裁队列相关的可用配置项。

| advanced.config Configuration Key | Description | Default value |
| --------------------------------- | ----------- | ------------- |
| rabbit.quorum_cluster_size        | Sets the default quorum queue cluster size (can be over-ridden by the x-quorum-initial-group-size queue argument at declaration time. | 3             |
| rabbit.quorum_commands_soft_limit | This is a flow control related parameter defining the maximum number of unconfirmed messages a channel accepts before entering flow. The current default is configured to provide good performance and stability when there are multiple publishers sending to the same quorum queue. If the applications typically only have a single publisher per queue this limit could be increased to provide somewhat better ingress rates. | 32            |

* rabbit.quorum_cluster_size
  Sets the default quorum queue cluster size (can be over-ridden by the x-quorum-initial-group-size queue argument at declaration time.
  设置默认仲裁队列集群大小（可以在声明时由 x-quorum-initial-group-size 队列参数覆盖）。

* rabbit.quorum_commands_soft_limit
  This is a flow control related parameter defining the maximum number of unconfirmed messages a channel accepts before entering flow. The current default is configured to provide good performance and stability when there are multiple publishers sending to the same quorum queue. If the applications typically only have a single publisher per queue this limit could be increased to provide somewhat better ingress rates.
  这是与流控制相关的参数，定义通道在进入流之前接受的未确认消息的最大数量。 当前的默认配置是当有多个发布者发送到同一仲裁队列时提供良好的性能和稳定性。 如果应用程序通常每个队列只有一个发布者，则可以增加此限制以提供更好的进入率。

### Example of a Quorum Queue Configuration

The following advanced.config example modifies all values listed above:

```bash
[
 %% five replicas by default, only makes sense for nine node clusters
 {rabbit, [{quorum_cluster_size, 5},
           {quorum_commands_soft_limit, 512}]}
].
```

## Poison Message Handling for Quorum Queues
仲裁队列的有害消息处理

Quorum queue support handling of [poison messages](https://en.wikipedia.org/wiki/Poison_message), that is, messages that cause a consumer to repeatedly requeue a delivery (possibly due to a consumer failure) such that the message is never consumed completely and [positively acknowledged](https://www.rabbitmq.com/confirms.html) so that it can be marked for deletion by RabbitMQ.
仲裁队列支持处理有毒消息，即导致消费者重复重新排队交付的消息（可能是由于消费者故障），使得该消息永远不会被完全消费并得到肯定确认，以便可以被 RabbitMQ 标记为删除 。

Quorum queues keep track of the number of unsuccessful delivery attempts and expose it in the "x-delivery-count" header that is included with any redelivered message.
仲裁队列跟踪不成功的传递尝试的次数，并将其公开在任何重新传递的消息中包含的“x-delivery-count”标头中。

It is possible to set a delivery limit for a queue using a [policy](https://www.rabbitmq.com/parameters.html#policies) argument, delivery-limit.
可以使用策略参数传递限制来设置队列的传递限制。

When a message has been returned more times than the limit the message will be dropped or [dead-lettered](https://www.rabbitmq.com/dlx.html) (if a DLX is configured).
当消息返回的次数超过限制时，该消息将被丢弃或成为死信（如果配置了 DLX）。

## Resources that Quorum Queues Use
仲裁队列使用的资源

Quorum queues are optimised for data safety and performance and typically require more resources (disk and RAM) than classic mirrored queues under a steady workload. Each quorum queue process maintains an in-memory index of the messages in the queue, which requires at least 32 bytes of metadata for each message (more, if the message was returned or has a TTL set). A quorum queue process will therefore use at least 1MB for every 30000 messages in the queue (message size is irrelevant). You can perform back-of-the-envelope calculations based on the number of queues and expected or maximum number of messages in them). Keeping the queues short is the best way to maintain low memory usage. [Setting the maximum queue length](https://www.rabbitmq.com/maxlength.html) for all queues is a good way to limit the total memory usage if the queues become long for any reason.
仲裁队列针对数据安全和性能进行了优化，并且在稳定的工作负载下通常需要比经典镜像队列更多的资源（磁盘和 RAM）。 每个仲裁队列进程都维护队列中消息的内存索引，这需要每条消息至少 32 字节的元数据（如果消息已返回或设置了 TTL，则需要更多字节）。 因此，仲裁队列进程将为队列中的每 30000 条消息使用至少 1MB（消息大小无关）。 您可以根据队列数量以及队列中的预期或最大消息数量执行粗略计算。 保持队列短是保持低内存使用率的最佳方法。 如果队列因任何原因变长，设置所有队列的最大队列长度是限制总内存使用量的好方法。

Additionally, quorum queues on a given node share a write-ahead-log (WAL) for all operations. WAL operations are stored both in memory and written to disk. When the current WAL file reaches a predefined limit, it is flushed to a WAL segment file on disk and the system will begin to release the memory used by that batch of log entries. The segment files are then compacted over time as consumers [acknowledge deliveries](https://www.rabbitmq.com/confirms.html). Compaction is the process that reclaims disk space.
此外，给定节点上的仲裁队列为所有操作共享预写日志 (WAL)。 WAL 操作既存储在内存中，又写入磁盘。 当当前的WAL文件达到预定义的限制时，它会被刷新到磁盘上的WAL段文件，系统将开始释放该批日志条目使用的内存。 当消费者确认交付时，段文件会随着时间的推移而被压缩。 压缩是回收磁盘空间的过程。

The WAL file size limit at which it is flushed to disk can be controlled:
可以控制刷新到磁盘的 WAL 文件大小限制：

```bash
# Flush current WAL file to a segment file on disk once it reaches 64 MiB in size
raft.wal_max_size_bytes = 64000000
```

The value defaults to 512 MiB. This means that during steady load, the WAL table memory footprint can reach 512 MiB. You can expect your memory usage to look like this: ![Quorum Queues WAL memory usage pattern](https://www.rabbitmq.com/img/memory/quorum-queue-memory-usage-pattern.png "Quorum Queues WAL memory usage pattern")
该值默认为 512 MiB。 这意味着在稳定负载期间，WAL 表内存占用量可以达到 512 MiB。 您可以预期您的内存使用情况如下所示：Quorum Queues WAL 内存使用模式

Because memory deallocation may take some time, we recommend that the RabbitMQ node is allocated at least 3 times the memory of the default WAL file size limit. More will be required in high-throughput systems. 4 times is a good starting point for those.
由于内存释放可能需要一些时间，因此我们建议 RabbitMQ 节点分配的内存至少是默认 WAL 文件大小限制的 3 倍。 高通量系统将需要更多。 4 次对于那些人来说是一个很好的起点。

### Repeated Requeues
重复请求

Internally quorum queues are implemented using a log where all operations including messages are persisted. To avoid this log growing too large it needs to be truncated regularly. To be able to truncate a section of the log all messages in that section needs to be acknowledged. Usage patterns that continuously [reject or nack](https://www.rabbitmq.com/nack.html) the same message with the requeue flag set to true could cause the log to grow in an unbounded fashion and eventually fill up the disks.
内部仲裁队列是使用日志实现的，其中包括消息在内的所有操作都被持久化。 为了避免该日志变得太大，需要定期截断它。 为了能够截断日志的一部分，需要确认该部分中的所有消息。 在重新排队标志设置为 true 的情况下连续拒绝或拒绝同一消息的使用模式可能会导致日志以无限制的方式增长并最终填满磁盘。

Messages that are rejected or nacked back to a quorum queue will be returned to the *back* of the queue *if* no [delivery-limit](https://www.rabbitmq.com/quorum-queues.html#poison-message-handling) is set. This avoids the above scenario where repeated re-queues causes the Raft log to grow in an unbounded manner. If a delivery-limit is set it will use the original behaviour of returning the message near the head of the queue.
如果未设置传递限制，则被拒绝或返回仲裁队列的消息将返回到队列的末尾。 这避免了上述重复重新排队导致 Raft 日志无限制增长的情况。 如果设置了传递限制，它将使用在队列头部附近返回消息的原始行为。

### Increased Atom Use
增加原子使用

The internal implementation of quorum queues converts the queue name into an Erlang atom. If queues with arbitrary names are continuously created and deleted it *may* threaten the long term stability of the RabbitMQ system (if the size of the atom table reaches the maximum limit, about 1M by default). It is not recommended to use quorum queues in this manner at this point.
仲裁队列的内部实现将队列名称转换为 Erlang 原子。 如果不断创建和删除任意名称的队列，可能会威胁到RabbitMQ系统的长期稳定性（如果原子表的大小达到最大限制，默认约为1M）。 目前不建议以这种方式使用仲裁队列。

## Quorum Queue Performance Tuning
仲裁队列性能调优

This section aims to cover a couple of tunable parameters that may increase throughput of quorum queues for **some workloads**. Other workloads may not see any increases, or observe decreases in throughput, with these settings.
本节旨在介绍几个可调参数，这些参数可能会增加某些工作负载的仲裁队列的吞吐量。 使用这些设置，其他工作负载可能不会看到吞吐量有任何增加，或者观察到吞吐量有所下降。

Use the values and recommendations here as a **starting point** and conduct your own benchmark (for example, [using PerfTest](https://rabbitmq.github.io/rabbitmq-perf-test/stable/htmlsingle/)) to conclude what combination of values works best for a particular workloads.
使用此处的值和建议作为起点，并执行您自己的基准测试（例如，使用 PerfTest）来得出最适合特定工作负载的值组合。

### Tuning: Raft Segment File Entry Count
调整：Raft 段文件条目计数

Workloads with small messages and higher message rates can benefit from the following configuration change that increases the number of Raft log entries (such as enqueued messages) that are allowed in a single write-ahead log file:
具有小消息和较高消息速率的工作负载可以受益于以下配置更改，该更改增加了单个预写日志文件中允许的 Raft 日志条目（例如排队消息）的数量：

```bash
# Positive values up to 65535 are allowed, the default is 4096.
raft.segment_max_entries = 32768
```

Values greater than 65535 are **not supported**.

### Tuning: Linux Readahead
调优：Linux 预读

In addition, the aforementioned workloads with a higher rate of small messages can benefit from a higher readahead, a configurable block device parameter of storage devices on Linux.
此外，上述具有较高小消息率的工作负载可以受益于较高的预读（Linux 上存储设备的可配置块设备参数）。

To inspect the effective readahead value, use [blockdev --getra](https://man7.org/linux/man-pages/man8/blockdev.8.html) and specify the block device that hosts RabbitMQ node data directory:
要检查有效预读值，请使用 blockdev --getra 并指定托管 RabbitMQ 节点数据目录的块设备：

```bash
# This is JUST AN EXAMPLE.
# The name of the block device in your environment will be different.
#
# Displays effective readahead value device /dev/sda.
sudo blockdev --getra /dev/sda
```

To configure readahead, use [blockdev --setra](https://man7.org/linux/man-pages/man8/blockdev.8.html) for the block device that hosts RabbitMQ node data directory:
要配置预读，请使用 blockdev --setra 作为托管 RabbitMQ 节点数据目录的块设备：

```bash
# This is JUST AN EXAMPLE.
# The name of the block device in your environment will be different.
# Values between 256 and 4096 in steps of 256 are most commonly used.
#
# Sets readahead for device /dev/sda to 4096.
sudo blockdev --setra 4096 /dev/sda
```

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them using [GitHub Discussions](https://github.com/rabbitmq/rabbitmq-server/discussions) or our community [Discord server](https://rabbitmq.com/discord).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!

