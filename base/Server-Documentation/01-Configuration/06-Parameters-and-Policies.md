# Parameters and Policies

https://www.rabbitmq.com/parameters.html

## Overview

While much of the configuration for RabbitMQ lives in the [configuration file](https://www.rabbitmq.com/configure.html), some things do not mesh well with the use of a configuration file:  虽然 RabbitMQ 的大部分配置都存在于配置文件中，但有些事情与使用配置文件并不能很好地结合：

- If they need to be the same across all nodes in a cluster  如果它们需要在集群中的所有节点上相同

- If they are likely to change at run time  如果它们可能在运行时发生变化

RabbitMQ calls these items **parameters**. Parameters can be set by invoking [rabbitmqctl](https://www.rabbitmq.com/man/rabbitmqctl.8.html) or through [the management plugin](https://www.rabbitmq.com/management.html)'s HTTP API. There are 2 kinds of parameters: **vhost-scoped parameters** and **global parameters**. Vhost-scoped parameters are tied to a virtual host and consist of **a component name, a name and a value**. Global parameters are not tied to a particular virtual host and they consist of **a name and value**.  RabbitMQ 调用这些项目参数。可以通过调用 rabbitmqctl 或通过管理插件的 HTTP API 来设置参数。有 2 种参数：vhost 范围的参数和全局参数。虚拟主机范围的参数与虚拟主机相关联，由组件名称、名称和值组成。全局参数不绑定到特定的虚拟主机，它们由名称和值组成。

One special case of parameters usage is **policies**, which are used for specifying optional arguments for groups of queues and exchanges, as well as plugins such as [Federation](https://www.rabbitmq.com/federation.html) and [Shovel](https://www.rabbitmq.com/shovel.html). Policies are vhost-scoped.  参数使用的一种特殊情况是策略，它用于为队列和交换组以及Federation 和Shovel 等插件指定可选参数。策略是虚拟主机范围的。

## Global and Per-vhost Parameters  全局和每个虚拟主机参数

As stated above, there are vhost-scoped parameters and global parameters. An example of vhost-scoped parameter is a federation upstream: it targets a component (federation-upstream), it has a name that identifies it, it's tied to a virtual host (federation links will target some resources of this virtual host), and its value defines connection parameters to an upstream broker.  如上所述，有虚拟主机范围的参数和全局参数。 vhost-scoped 参数的一个示例是上游联合：它针对一个组件（federation-upstream），它有一个标识它的名称，它绑定到一个虚拟主机（联合链接将针对这个虚拟主机的一些资源），以及它的值定义到上游代理的连接参数。

Vhost-scoped parameters can be set, cleared and listed:  可以设置、清除和列出虚拟主机范围的参数：

* rabbitmqctl

```bash
set_parameter                 Sets a runtime parameter.
clear_parameter               Clears a runtime parameter.
list_parameters               Lists runtime parameters for a virtual host

$ rabbitmqctl [--node <node>] [--longnames] [--quiet] set_parameter [--vhost <vhost>] <component_name> <name> <value>

<component_name> component name

<name> parameter name (identifier)

<value> parameter value

$ rabbitmqctl set_parameter federation-upstream f1 '{"uri":"amqp://root:root123@192.168.0.2:5672","ack-mode":"on-confirm"}'

```

* HTTP API

```bash
PUT /api/parameters/*{component_name}*/*{vhost}*/*{name}*
DELETE /api/parameters/*{component_name}*/*{vhost}*/*{name}*
GET /api/parameters |
```

Global parameters is the other kind of parameters. An example of a global parameter is the name of the cluster. Global parameters can be set, cleared and listed:  全局参数是另一种参数。 全局参数的一个示例是集群的名称。 全局参数可以设置、清除和列出：

* rabbitmqctl

```bash
set_global_parameter          Sets a runtime parameter.
clear_global_parameter        Clears a global runtime parameter
list_global_parameters        Lists global runtime parameters

$ rabbitmqctl [--node <node>] [--longnames] [--quiet] set_global_parameter <name> <value>

<name> global parameter name (identifier)

<value> parameter value

$ rabbitmqctl list_global_parameters
Listing global runtime parameters ...
name	value
cluster_name	"rabbitmq1"
internal_cluster_id	"rabbitmq-cluster-id-HLVgEleDZw5Qy5d7zMje5A"
$
```

* HTTP API

```bash
PUT /api/global-parameters/*name*
DELETE /api/global-parameters/*name*
GET /api/global-parameters
```

Since a parameter value is a JSON document, you will usually need to quote it when creating one on the command line with rabbitmqctl. On Unix it is usually easiest to quote the whole document with single quotes, and use double quotes within it. On Windows you will have to escape every double quote. We give examples for both Unix and Windows for this reason.  由于参数值是 JSON 文档，因此在使用 rabbitmqctl 在命令行上创建参数值时通常需要引用它。在 Unix 上，通常最容易用单引号引用整个文档，并在其中使用双引号。在 Windows 上，您必须转义每个双引号。出于这个原因，我们给出了 Unix 和 Windows 的示例。

Parameters reside in the database used by RabbitMQ for definitions of virtual hosts, exchanges, queues, bindings, users and permissions. Parameters are exported along with other object definitions by the management plugin's export feature.  参数驻留在 RabbitMQ 使用的数据库中，用于定义虚拟主机、交换器、队列、绑定、用户和权限。管理插件的导出功能将参数与其他对象定义一起导出。

Vhost-scoped parameters are used by the federation and shovel plugins. Global parameters are used by the MQTT plugin.  联邦和铲子插件使用虚拟主机范围的参数。 MQTT 插件使用全局参数。

## Policies

### Why Policies Exist

Before we explain what policies are and how to use them it would be helpful to explain why they were introduced to RabbitMQ.  在我们解释什么是策略以及如何使用它们之前，解释一下为什么将它们引入 RabbitMQ 会很有帮助。

In addition to mandatory properties (e.g. durable or exclusive), queues and exchanges in RabbitMQ have optional parameters (arguments), sometimes referred to as x-arguments. Those are provided by clients when they declare queues (exchanges) and control various optional features, such as [queue length limit](https://www.rabbitmq.com/maxlength.html) or [TTL](https://www.rabbitmq.com/ttl.html).  除了强制属性（例如持久性或独占性）之外，RabbitMQ 中的队列和交换器还有可选参数（参数），有时称为 x-arguments。这些是由客户端在声明队列（交换）并控制各种可选功能（例如队列长度限制或 TTL）时提供的。

Client-controlled properties in some of the protocols RabbitMQ supports generally work well but they can be inflexible: updating TTL values or mirroring parameters that way required application changes, redeployment and queue re-declaration (which involves deletion). In addition, there is no way to control the extra arguments for groups of queues and exchanges. Policies were introduced to address the above pain points.  RabbitMQ 支持的某些协议中的客户端控制属性通常运行良好，但它们可能不灵活：更新 TTL 值或镜像参数需要应用程序更改、重新部署和队列重新声明（这涉及删除）。此外，无法控制队列和交换组的额外参数。引入了政策来解决上述痛点。

A policy matches one or more queues by name (using a regular expression pattern) and appends its definition (a map of optional arguments) to the x-arguments of the matching queues. In other words, it is possible to configure x-arguments for multiple queues at once with a policy, and update them all at once by updating policy definition.  策略按名称匹配一个或多个队列（使用正则表达式模式）并将其定义（可选参数的映射）附加到匹配队列的 x 参数。换句话说，可以使用策略一次为多个队列配置 x 参数，并通过更新策略定义一次更新它们。

In modern versions of RabbitMQ the set of features which can be controlled by policy is not the same as the set of features which can be controlled by client-provided arguments.  在 RabbitMQ 的现代版本中，可由策略控制的特性集与可由客户端提供的参数控制的特性集不同。

> 通常来说，我们会在创建一个 Queue 时指定了队列的各项属性和参数，例如 `message-ttl`、`x-dead-letter-exchange`、`x-dead-letter-route-key`、`x-max-length`等，一旦 RabbitMQ 创建 Queue 成功，则不允许再修改 Queue 的任何属性或参数，除非删除重建。
> 
> 我们在实际使用中，随着业务的变化，经常会出现需要调整 `message-ttl`、`x-max-length`等参数来改变 Queue 行为的场景，那这个时候怎么办呢，总不能暂停上下游的生产者和消费者来等 Queue 的删除重建吧？这个时候就需要使用到 RabbitMQ 的 Policy 机制了。
> 

### How Policies Work

![](../../../images/policy-01.png)

Key policy attributes are  关键政策属性是

- name: it can be anything but ASCII-based names without spaces are recommended  可以是任何名称，但建议使用不带空格的基于 ASCII 的名称

- pattern: a regular expression that matches one or more queue (exchange) names. Any regular expression can be used.  匹配一个或多个队列（交换）名称的正则表达式。可以使用任何正则表达式。

- definition: a set of key/value pairs (think a JSON document) that will be injected into the map of optional arguments of the matching queues and exchanges  一组键/值对（想象一个 JSON 文档）将被注入到匹配队列和交换的可选参数的映射中

- priority: see below  见下文

Policies automatically match against exchanges and queues, and help determine how they behave. Each exchange or queue will have at most one policy matching (see Combining Policy Definitions below), and each policy then injects a set of key-value pairs (policy definition) on to the matching queues (exchanges).  策略会自动匹配交换和队列，并帮助确定它们的行为方式。每个交换或队列最多有一个策略匹配（请参阅下面的组合策略定义），然后每个策略将一组键值对（策略定义）注入匹配的队列（交换）。

Policies can match only queues, only exchanges, or both. This is controlled using the apply-to flag when a policy is created.  策略可以只匹配队列、只匹配交换，或者两者都匹配。这在创建策略时使用 apply-to 标志进行控制。

Policies can change at any time. When a policy definition is updated, its effect on matching exchanges and queues will be reapplied. Usually it happens instantaneously but for very busy queues can take a bit of time (say, a few seconds).  政策可以随时更改。更新策略定义时，将重新应用其对匹配交换和队列的影响。通常它会立即发生，但对于非常繁忙的队列可能需要一些时间（比如几秒钟）。

Policies are matched and applied every time an exchange or queue is created, not just when the policy is created.  每次创建交换或队列时都会匹配和应用策略，而不仅仅是在创建策略时。

Policies can be used to configure  策略可用于配置

- [federation](https://www.rabbitmq.com/federation.html)  联邦
- [alternate exchanges](https://www.rabbitmq.com/ae.html)  使用备用交换器处理无法路由的消息
- [dead lettering](https://www.rabbitmq.com/dlx.html)

在 RabbitMQ 中，尽管听起来像是一种特殊的交换器，但死信交换器是一种正常的交换器。创建它时没有特别的要求也不需要执行特别的操作。使交换器成为死信交换器的唯一要做的事情是在创建队列时声明该交换器将被用作保存被拒绝的消息。一旦拒绝了一个不重新发送的消息，RabbitMQ 将把消息路由到队列的 **x-dead-letter-exchange** 参数中指定的交换器。

在发出 Queue.Declare RPC 请求时作为 x-dead-letter-exchange 参数进行传入。

- [per-queue TTLs](https://www.rabbitmq.com/ttl.html)

- [queue length limit](https://www.rabbitmq.com/maxlength.html)

and other features.

An example of defining a policy looks like:  定义策略的示例如下所示：

* rabbitmqctl

```bash
$ rabbitmqctl [--node <node>] [--longnames] [--quiet] set_policy [--vhost <vhost>] [--priority <priority>] [--apply-to <apply-to>] <name> <pattern> <definition>

<name> policy name (identifier)

<pattern> regular expression pattern that will be used to match queues, exchanges, etc

<definition> policy definition (arguments). Must be a valid JSON document

--priority <priority> policy priority

--apply-to <queues | exchanges | all> policy should only apply to queues, exchanges, or all entities (both of the above)

$ rabbitmqctl set_policy --vhost / --priority 1 --apply-to exchanges federate-me "^federated\." '{"federation-upstream-set":"all"}'

```

* rabbitmqctl (Windows)

```bash
rabbitmqctl.bat set_policy federate-me "^federated\." "{""federation-upstream-set"":""all""}" --priority 1 --apply-to exchanges
```

* HTTP API

```bash
PUT /api/policies/%2f/federate-me
{
"pattern": "^federated\.",
"definition": {"federation-upstream-set":"all"},
"priority": 1,
"apply-to": "exchanges"
}
```

* Web UI

	* Navigate to Admin > Policies > Add / update a policy.

	* Enter "federate-me" next to Name, "^federated\." next to Pattern, and select "Exchanges" next to Apply to.

	* Enter "federation-upstream-set" = "all" in the first line next to Policy.

	* Click Add policy.

This matches the value "all" with the key "federation-upstream-set" for all exchanges with names beginning with "federated.", in the virtual host "/".  这将值“all”与虚拟主机“/”中名称以“federated.”开头的所有交换的键“federation-upstream-set”匹配。

The "pattern" argument is a regular expression used to match exchange or queue names.  “模式”参数是用于匹配交换或队列名称的正则表达式。

In the event that more than one policy can match a given exchange or queue, the policy with the greatest priority applies.  如果多个策略可以匹配给定的交换或队列，则应用具有最高优先级的策略。

The "apply-to" argument can be "exchanges", "queues" or "all". The "apply-to" and "priority" settings are optional, in which case the defaults are "all" and "0" respectively.  “apply-to”参数可以是“exchanges”、“queues”或“all”。 “应用到”和“优先级”设置是可选的，在这种情况下，默认值分别为“全部”和“0”。

### Combining Policy Definitions  结合策略定义

In some cases we might want to apply more than one policy definition to a resource. For example we might need a queue to be federated and mirrored. At most one policy will apply to a resource at any given time, but we can apply multiple definitions in that policy.  在某些情况下，我们可能希望将多个策略定义应用于资源。例如，我们可能需要一个队列来联合和镜像。在任何给定时间最多将一个策略应用于资源，但我们可以在该策略中应用多个定义。

A federation policy definition would require an *upstream set* to be specified, so we would need the federation-upstream-set key in our definition. On the other hand to define some queues as mirrored, we would need the ha-mode key to be defined as well for the policy. The policy definition is just a JSON object and can have multiple keys combined in the same policy definition.  联合策略定义需要指定上游集，因此我们需要定义中的 federation-upstream-set 键。另一方面，要将某些队列定义为镜像，我们还需要为策略定义 ha-mode 键。策略定义只是一个 JSON 对象，可以在同一个策略定义中组合多个键。

Here's an example:  这是一个例子：

* rabbitmqctl

```bash
rabbitmqctl set_policy ha-fed "^hf\." '{"federation-upstream-set":"all", "ha-mode":"exactly", "ha-params":2}' --priority 1 --apply-to queues
```

* rabbitmqctl (Windows)

```bash
rabbitmqctl set_policy ha-fed "^hf\." "{""federation-upstream-set"":""all"", ""ha-mode"":""exactly"", ""ha-params"":2}" --priority 1 --apply-to queues
```

* HTTP API

```bash
PUT /api/policies/%2f/ha-fed
{
"pattern": "^hf\.",
"definition": {"federation-upstream-set":"all", "ha-mode":"exactly", "ha-params":2},
"priority": 1,
"apply-to": "queues"
}
```

* Web UI

	* Navigate to Admin > Policies > Add / update a policy.

	* Enter "ha-fed" next to Name, "^hf\." next to Pattern, and select "Queues" next to Apply to.

	* Enter "federation-upstream-set" = "all" in the first line next to Policy.

	* Enter "ha-mode" = "exactly" and "ha-params" = 2 on the following form lines.

	* Click Add policy.

By doing that all the queues matched by the pattern "^hf\." will have the "federation-upstream-set" and the policy definitions applied to them.  通过这样做，所有与模式“^hf\”匹配的队列。将有“federation-upstream-set”和应用于它们的策略定义。

## Operator Policies

### Difference From Regular Policies  与常规政策的区别

![](../../../images/policy-01.png)

Sometimes it is necessary for the operator to enforce certain policies. For example, it may be desirable to force [queue TTL](https://www.rabbitmq.com/ttl.html) but still let other users manage policies. Operator policies allow for that.  有时运营商需要强制执行某些策略。例如，可能需要强制队列 TTL，但仍让其他用户管理策略。运营商政策允许这样做。

Operator policies are much like regular ones but their definitions are used differently. They are merged with regular policy definitions before the result is applied to matching queues.  运营商政策与常规政策非常相似，但它们的定义使用不同。在将结果应用于匹配队列之前，它们会与常规策略定义合并。

Because operator policies can unexpectedly change queue attributes and, in turn, application assumptions and semantics, they are limited only to a few arguments:  由于操作员策略可能会意外更改队列属性，进而更改应用程序假设和语义，因此它们仅限于几个参数：

- expires

- message-ttl

- max-length

- max-length-bytes

The arguments above are all numeric. The reason for that is explained in the following section.  上面的参数都是数字的。下一节将解释其原因。

### Conflict Resolution with Regular Policies  使用常规策略解决冲突

An operator policy and a regular one can contain the same keys in their definitions. When it happens, the smaller value is chosen as effective. For example, if a matching operator policy definition sets max-length to 50 and a matching regular policy definition uses the value of 100, the value of 50 will be used. If, however, regular policy's value was 20, it would be used. Operator policies, therefore, don't just overwrite regular policy values. They enforce limits but try to not override user-provided policies where possible.  操作员策略和常规策略可以在其定义中包含相同的键。当它发生时，选择较小的值作为有效值。例如，如果匹配的操作员策略定义将 max-length 设置为 50，并且匹配的常规策略定义使用 100 的值，则将使用 50 的值。但是，如果常规策略的值为 20，则将使用它。因此，运营商策略不只是覆盖常规策略值。他们强制执行限制，但尽可能不覆盖用户提供的策略。

### Defining Operator Policies

Operator policies are defined in a way very similar to regular (user) policies. When rabbitmqctl is used, the command name is set_operator_policy instead of set_policy. In the HTTP API, /api/policies/ in request path becomes /api/operator-policies/:  操作员策略的定义方式与常规（用户）策略非常相似。使用rabbitmqctl时，命令名是set_operator_policy，而不是set_policy。在 HTTP API 中，请求路径中的 /api/policies/ 变为 /api/operator-policies/：

* rabbitmqctl

```bash
$ rabbitmqctl [--node <node>] [--longnames] [--quiet] set_operator_policy [--vhost <vhost>] [--priority <priority>] [--apply-to <apply-to>] <name> <pattern> <definition>

<name> policy name (identifier)

<pattern> a regular expression pattern that will be used to match queue, exchanges, etc

<definition> policy definition (arguments). Must be a valid JSON document

--priority <priority> policy priority

--apply-to <queues | exchanges | all> policy should only apply to queues, exchanges, or all entities (both of the above)

$ rabbitmqctl set_operator_policy transient-queue-ttl "^amq\." '{"expires":1800000}' --priority 1 --apply-to queues

```

* rabbitmqctl (Windows)

```bash
rabbitmqctl.bat set_operator_policy transient-queue-ttl "^amq\." "{""expires"": 1800000}" --priority 1 --apply-to queues
```

* HTTP API

```bash
PUT /api/operator-policies/%2f/transient-queue-ttl
{
"pattern": "^amq\.",
"definition": {"expires": 1800000},
"priority": 1,
"apply-to": "queues"
}
```

* Web UI

	* Navigate to Admin > Policies > Add / update an operator policy.

	* Enter "transient-queue-ttl" next to Name, "^amq\." next to Pattern, and select "Queues" next to Apply to.

	* Enter "expires" = 1800000 in the first line next to Policy.

	* Click Add policy.

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!









