# Clustering Guide

https://www.rabbitmq.com/clustering.html

## Overview

This guide covers fundamental topics related to RabbitMQ clustering:  本指南涵盖了与 RabbitMQ 集群相关的基本主题：

- How RabbitMQ nodes are identified: node names

- Requirements for clustering

- What data is and isn't replicated between cluster nodes

- What clustering means for clients

- How clusters are formed

- How nodes authenticate to each other and with CLI tools

- Why it's important to use an odd number of nodes and **two cluster nodes are highly recommended against**

- Node restarts and how nodes rejoin their cluster

- Node readiness probes and how they can affect rolling cluster restarts

- How to remove a cluster node

- How to reset a cluster node to a pristine (blank) state

and more. [Cluster Formation and Peer Discovery](https://www.rabbitmq.com/cluster-formation.html) is a closely related guide that focuses on peer discovery and cluster formation automation-related topics. For queue contents (message) replication, see the [Quorum Queues](https://www.rabbitmq.com/quorum-queues.html) guide.

[Tanzu RabbitMQ](https://www.rabbitmq.com/tanzu) provides an [inter-node traffic compression](https://www.rabbitmq.com/clustering-compression.html) feature.

A RabbitMQ cluster is a logical grouping of one or several nodes, each sharing users, virtual hosts, queues, exchanges, bindings, runtime parameters and other distributed state.  RabbitMQ 集群是一个或多个节点的逻辑分组，每个节点共享用户、虚拟主机、队列、交换器、绑定、运行时参数和其他分布式状态。

## Cluster Formation

### Ways of Forming a Cluster

A RabbitMQ cluster can be formed in a number of ways:  RabbitMQ 集群可以通过多种方式形成：

- Declaratively by listing cluster nodes in [config file](https://www.rabbitmq.com/configure.html)

- Declaratively 声明式地 using DNS-based discovery

- Declaratively using [AWS (EC2) instance discovery](https://github.com/rabbitmq/rabbitmq-peer-discovery-aws) (via a plugin)

- Declaratively using [Kubernetes discovery](https://github.com/rabbitmq/rabbitmq-peer-discovery-k8s) (via a plugin)

- Declaratively using [Consul-based discovery](https://github.com/rabbitmq/rabbitmq-peer-discovery-consul) (via a plugin)

- Declaratively using [etcd-based discovery](https://github.com/rabbitmq/rabbitmq-peer-discovery-etcd) (via a plugin)

- Manually with rabbitmqctl

Please refer to the [Cluster Formation guide](https://www.rabbitmq.com/cluster-formation.html) for details.

The composition of a cluster can be altered dynamically. All RabbitMQ brokers start out as running on a single node. These nodes can be joined into clusters, and subsequently turned back into individual brokers again.  集群的组成可以动态改变。 所有 RabbitMQ 代理最初都在单个节点上运行。 这些节点可以加入集群，然后再次变回单个代理。

### Node Names (Identifiers)

RabbitMQ nodes are identified by node names. A node name consists of two parts, a prefix (usually `rabbit`) and hostname. For example, `rabbit@node1.messaging.svc.local` is a node name with the prefix of `rabbit` and hostname of `node1.messaging.svc.local`.

Node names in a cluster must be unique. If more than one node is running on a given host (this is usually the case in development and QA environments), they must use different prefixes, e.g. `rabbit1@hostname` and `rabbit2@hostname`.

In a cluster, nodes identify and contact each other using node names. This means that the hostname part of every node name must resolve. [CLI tools](https://www.rabbitmq.com/cli.html) also identify and address nodes using node names.  在集群中，节点使用节点名称相互识别和联系。 这意味着必须解析每个节点名称的主机名部分。 CLI 工具还使用节点名称识别和寻址节点。

When a node starts up, it checks whether it has been assigned a node name. This is done via the RABBITMQ_NODENAME [environment variable](https://www.rabbitmq.com/configure.html#supported-environment-variables). If no value was explicitly configured, the node resolves its hostname and prepends rabbit to it to compute its node name.  当一个节点启动时，它会检查它是否已经被分配了一个节点名称。 这是通过 RABBITMQ_NODENAME 环境变量完成的。 如果没有明确配置值，则节点解析其主机名并在其前面添加 rabbit 以计算其节点名称。

If a system uses fully qualified domain names (FQDNs) for hostnames, RabbitMQ nodes and CLI tools must be configured to use so called long node names. For server nodes this is done by setting the `RABBITMQ_USE_LONGNAME` [environment variable](https://www.rabbitmq.com/configure.html#supported-environment-variables) to true.  如果系统使用完全限定域名 (FQDN) 作为主机名，则 RabbitMQ 节点和 CLI 工具必须配置为使用所谓的长节点名称。 对于服务器节点，这是通过将 RABBITMQ_USE_LONGNAME 环境变量设置为 true 来完成的。

For CLI tools, either RABBITMQ_USE_LONGNAME must be set or the `--longnames` option must be specified.

## Cluster Formation Requirements

### Hostname Resolution

RabbitMQ nodes address each other using domain names, either short or fully-qualified (FQDNs). Therefore hostnames of all cluster members must be resolvable from all cluster nodes, as well as machines on which command line tools such as `rabbitmqctl` might be used.  RabbitMQ 节点使用域名相互寻址，无论是短域名还是完全限定域名 (FQDN)。 因此，所有集群成员的主机名必须可从所有集群节点以及可能使用 rabbitmqctl 等命令行工具的机器解析。

Hostname resolution can use any of the standard OS-provided methods:

- DNS records

- Local host files (e.g. /etc/hosts)

In more restrictive environments, where DNS record or hosts file modification is restricted, impossible or undesired, [Erlang VM can be configured to use alternative hostname resolution methods](http://erlang.org/doc/apps/erts/inet_cfg.html), such as an alternative DNS server, a local file, a non-standard hosts file location, or a mix of methods. Those methods can work in concert with the standard OS hostname resolution methods.  在更严格的环境中，DNS 记录或主机文件修改受到限制、不可能或不需要，Erlang VM 可以配置为使用替代主机名解析方法，例如替代 DNS 服务器、本地文件、非标准主机文件位置、 或混合方法。 这些方法可以与标准操作系统主机名解析方法协同工作。

To use FQDNs, see RABBITMQ_USE_LONGNAME in the [Configuration guide](https://www.rabbitmq.com/configure.html#supported-environment-variables). See Node Names above.

## Port Access

RabbitMQ nodes bind to ports (open server TCP sockets) in order to accept client and CLI tool connections. Other processes and tools such as SELinux may prevent RabbitMQ from binding to a port. When that happens, the node will fail to start.

CLI tools, client libraries and RabbitMQ nodes also open connections (client TCP sockets). Firewalls can prevent nodes and CLI tools from communicating with each other. Make sure the following ports are accessible:

- 4369: [epmd](https://www.rabbitmq.com/networking.html#epmd), a helper discovery daemon used by RabbitMQ nodes and CLI tools

- 5672, 5671: used by AMQP 0-9-1 and 1.0 clients without and with TLS

- 25672: used for inter-node and CLI tools communication (Erlang distribution server port) and is allocated from a dynamic range (limited to a single port by default, computed as AMQP port + 20000). Unless external connections on these ports are really necessary (e.g. the cluster uses [federation](https://www.rabbitmq.com/federation.html) or CLI tools are used on machines outside the subnet), these ports should not be publicly exposed. See [networking guide](https://www.rabbitmq.com/networking.html) for details.  用于节点间和 CLI 工具通信（Erlang 分发服务器端口）并从动态范围分配（默认限制为单个端口，计算为 AMQP 端口 + 20000）。 除非确实需要这些端口上的外部连接（例如，集群使用联合或在子网外的机器上使用 CLI 工具），否则不应公开这些端口。 有关详细信息，请参阅网络指南。

- 35672-35682: used by CLI tools (Erlang distribution client ports) for communication with nodes and is allocated from a dynamic range (computed as server distribution port + 10000 through server distribution port + 10010). See [networking guide](https://www.rabbitmq.com/networking.html) for details.

- 15672: [HTTP API](https://www.rabbitmq.com/management.html) clients, [management UI](https://www.rabbitmq.com/management.html) and [rabbitmqadmin](https://www.rabbitmq.com/management-cli.html) (only if the [management plugin](https://www.rabbitmq.com/management.html) is enabled)

- 61613, 61614: [STOMP clients](https://stomp.github.io/stomp-specification-1.2.html) without and with TLS (only if the [STOMP plugin](https://www.rabbitmq.com/stomp.html) is enabled)

- 1883, 8883: [MQTT clients](http://mqtt.org/) without and with TLS, if the [MQTT plugin](https://www.rabbitmq.com/mqtt.html) is enabled

- 15674: STOMP-over-WebSockets clients (only if the [Web STOMP plugin](https://www.rabbitmq.com/web-stomp.html) is enabled)

- 15675: MQTT-over-WebSockets clients (only if the [Web MQTT plugin](https://www.rabbitmq.com/web-mqtt.html) is enabled)

- 15692: Prometheus metrics (only if the [Prometheus plugin](https://www.rabbitmq.com/prometheus.html) is enabled)

It is possible to [configure RabbitMQ](https://www.rabbitmq.com/configure.html) to use [different ports and specific network interfaces](https://www.rabbitmq.com/networking.html).

## Nodes in a Cluster

### What is Replicated?

All data/state required for the operation of a RabbitMQ broker is replicated across all nodes. An exception to this are message queues, which by default reside on one node, though they are visible and reachable from all nodes. To replicate queues across nodes in a cluster, use a queue type that supports replication. This topic is covered in the [Quorum Queues](https://www.rabbitmq.com/quorum-queues.html) guide.  RabbitMQ 代理操作所需的所有数据/状态都在所有节点之间复制。 一个例外是消息队列，它们默认驻留在一个节点上，尽管它们对所有节点都是可见和可访问的。 要跨集群中的节点复制队列，请使用支持复制的队列类型。 Quorum Queues 指南中介绍了该主题。

### Nodes are Equal Peers

Some distributed systems have leader and follower nodes. This is generally not true for RabbitMQ. All nodes in a RabbitMQ cluster are equal peers: there are no special nodes in RabbitMQ core. This topic becomes more nuanced when [quorum queues](https://www.rabbitmq.com/quorum-queues.html) and plugins are taken into consideration but for most intents and purposes, all cluster nodes should be considered equal.  一些分布式系统具有领导者和追随者节点。 对于 RabbitMQ 来说，这通常不是真的。 RabbitMQ 集群中的所有节点都是对等的：RabbitMQ 核心中没有特殊节点。 当考虑仲裁队列和插件时，这个主题变得更加微妙，但对于大多数意图和目的，所有集群节点都应该被视为平等。

Many [CLI tool](https://www.rabbitmq.com/cli.html) operations can be executed against any node. An [HTTP API](https://www.rabbitmq.com/management.html) client can target any cluster node.  许多 CLI 工具操作可以针对任何节点执行。 HTTP API 客户端可以针对任何集群节点。

Individual plugins can designate (elect) certain nodes to be "special" for a period of time. For example, [federation links](https://www.rabbitmq.com/federation.html) are colocated on a particular cluster node. Should that node fail, the links will be restarted on a different node.  单个插件可以在一段时间内指定（选举）某些节点为“特殊”节点。 例如，联合链接位于特定的集群节点上。 如果该节点出现故障，链接将在不同的节点上重新启动。

In versions older than 3.6.7, [RabbitMQ management plugin](https://www.rabbitmq.com/management.html) used a dedicated node for stats collection and aggregation.

### How CLI Tools Authenticate to Nodes (and Nodes to Each Other): the Erlang Cookie

RabbitMQ nodes and CLI tools (e.g. rabbitmqctl) use a cookie to determine whether they are allowed to communicate with each other. For two nodes to be able to communicate they must have the same shared secret called the Erlang cookie. The cookie is just a string of alphanumeric characters up to 255 characters in size. It is usually stored in a local file. The file must be only accessible to the owner (e.g. have UNIX permissions of 600 or similar). Every cluster node must have the same cookie.  RabbitMQ 节点和 CLI 工具（例如 rabbitmqctl）使用 cookie 来确定它们是否被允许相互通信。 为了让两个节点能够通信，它们必须具有相同的共享秘密，称为 Erlang cookie。 cookie 只是一串最多 255 个字符的字母数字字符。 它通常存储在本地文件中。 该文件必须只能由所有者访问（例如，具有 600 或类似的 UNIX 权限）。 每个集群节点必须具有相同的 cookie。

If the file does not exist, Erlang VM will try to create one with a randomly generated value when the RabbitMQ server starts up. Using such generated cookie files are appropriate in development environments only. Since each node will generate its own value independently, this strategy is not really viable in a clustered environment.  如果该文件不存在，Erlang VM 会在 RabbitMQ 服务器启动时尝试创建一个随机生成的值。 使用此类生成的 cookie 文件仅适用于开发环境。 由于每个节点都会独立产生自己的价值，这种策略在集群环境中并不真正可行。

Erlang cookie generation should be done at cluster deployment stage, ideally using automation and orchestration tools.  Erlang cookie 生成应该在集群部署阶段完成，最好使用自动化和编排工具。

In distributed deployment  在分布式部署中

### Cookie File Locations

#### Linux, MacOS, *BSD

On UNIX systems, the cookie will be typically located in `/var/lib/rabbitmq/.erlang.cookie` (used by the server) and `$HOME/.erlang.cookie` (used by CLI tools). Note that since the value of `$HOME` varies from user to user, it's necessary to place a copy of the cookie file for each user that will be using the CLI tools. This applies to both non-privileged users and `root`.  在 UNIX 系统上，cookie 通常位于 /var/lib/rabbitmq/.erlang.cookie（由服务器使用）和 $HOME/.erlang.cookie（由 CLI 工具使用）。 请注意，由于 $HOME 的值因用户而异，因此有必要为将使用 CLI 工具的每个用户放置 cookie 文件的副本。 这适用于非特权用户和 root

RabbitMQ nodes will log its effective user's home directory location early on boot.  RabbitMQ 节点将在启动初期记录其有效用户的主目录位置。

#### Community Docker Image and Kubernetes

[Docker community RabbitMQ image](https://github.com/docker-library/rabbitmq/) uses `RABBITMQ_ERLANG_COOKIE` environment variable value to populate the cookie file.

Configuration management and container orchestration tools that use this image must make sure that every RabbitMQ node container in a cluster uses the same value.  使用此镜像的配置管理和容器编排工具必须确保集群中的每个 RabbitMQ 节点容器都使用相同的值。

In the context of Kubernetes, the value must be specified in the pod template specification of the [stateful set](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/). For instance, this can be seen in the [RabbitMQ on Kubernetes examples repository](https://github.com/rabbitmq/diy-kubernetes-examples).

#### Windows

On Windows, the cookie location depends on a few factors:

- Whether the `HOMEDRIVE` and `HOMEPATH` environment variables are both set

- Erlang version: prior to 20.2 (these are no longer supported by any [maintained release series of RabbitMQ](https://www.rabbitmq.com/versions.html)) or 20.2 and later

##### Erlang 20.2 or later

With Erlang versions starting with 20.2, the cookie file locations are:

- `%HOMEDRIVE%%HOMEPATH%\.erlang.cookie` (usually `C:\Users\%USERNAME%\.erlang.cookie` for user `%USERNAME%`) if both the `HOMEDRIVE` and `HOMEPATH` environment variables are set

- `%USERPROFILE%\.erlang.cookie` (usually `C:\Users\%USERNAME%\.erlang.cookie`) if `HOMEDRIVE` and `HOMEPATH` are not both set

- For the RabbitMQ Windows service - `%USERPROFILE%\.erlang.cookie` (usually `C:\WINDOWS\system32\config\systemprofile`)

If the Windows service is used, the cookie should be copied from `C:\Windows\system32\config\systemprofile\.erlang.cookie` to the expected location for users running commands like rabbitmqctl.bat.

### Overriding Using CLI and Runtime Command Line Arguments  使用 CLI 和运行时命令行参数覆盖

As an alternative, the option "-setcookie <value>" can be added to `RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS` [environment variable value](https://www.rabbitmq.com/configure.html) to override the cookie value used by a RabbitMQ node:  作为替代方案，可以将选项“-setcookie <value>”添加到 RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS 环境变量值以覆盖 RabbitMQ 节点使用的 cookie 值：

```bash
RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS="-setcookie cookie-value"
```

CLI tools can take a cookie value using a command line flag:

```bash
rabbitmq-diagnostics status --erlang-cookie "cookie-value"
```

Both are **the least secure options** and generally **not recommended**.  两者都是最不安全的选项，通常不推荐。

### Troubleshooting

When a node starts, it will [log](https://www.rabbitmq.com/logging.html) the home directory location of its effective user:

```plaintext
node           : rabbit@cdbf4de5f22d
home dir       : /var/lib/rabbitmq
```

Unless any [server directories](https://www.rabbitmq.com/relocate.html) were overridden, that's the directory where the cookie file will be looked for, and created by the node on first boot if it does not already exist.  除非任何服务器目录被覆盖，否则将在其中查找 cookie 文件的目录，如果它不存在，则由节点在第一次启动时创建。

In the example above, the cookie file location will be `/var/lib/rabbitmq/.erlang.cookie`.

### Authentication Failures

When the cookie is misconfigured (for example, not identical), RabbitMQ nodes will log errors such as "Connection attempt from disallowed node", "", "Could not auto-cluster".  当 cookie 配置错误（例如，不相同）时，RabbitMQ 节点将记录错误，例如“来自不允许的节点的连接尝试”、“”、“无法自动集群”。

For example, when a CLI tool connects and tries to authenticate using a mismatching secret value:

```plaintext
2020-06-15 13:03:33 [error] <0.1187.0> ** Connection attempt from node 'rabbitmqcli-99391-rabbit@warp10' rejected. Invalid challenge reply. **
```

When a CLI tool such as rabbitmqctl fails to authenticate with RabbitMQ, the message usually says

```plaintext
* epmd reports node 'rabbit' running on port 25672
* TCP connection succeeded but Erlang distribution failed
* suggestion: hostname mismatch?
* suggestion: is the cookie set correctly?
* suggestion: is the Erlang distribution using TLS?
```

An incorrectly placed cookie file or cookie value mismatch are most common scenarios for such failures.  错误放置的 cookie 文件或 cookie 值不匹配是此类故障的最常见情况。

When a recent Erlang/OTP version is used, authentication failures contain more information and cookie mismatches can be identified better:  当使用最新的 Erlang/OTP 版本时，身份验证失败包含更多信息，可以更好地识别 cookie 不匹配：

```ini
* connected to epmd (port 4369) on warp10
* epmd reports node 'rabbit' running on port 25672
* TCP connection succeeded but Erlang distribution failed

* Authentication failed (rejected by the remote node), please check the Erlang cookie
```

See the [CLI Tools guide](https://www.rabbitmq.com/cli.html) for more information.

#### Hostname Resolution

Since hostname resolution is a prerequisite for successful inter-node communication, starting with [RabbitMQ 3.8.6](https://www.rabbitmq.com/changelog.html), CLI tools provide two commands that help verify that hostname resolution on a node works as expected. The commands are not meant to replace [dig](https://en.wikipedia.org/wiki/Dig_(command)) and other specialised DNS tools but rather provide a way to perform most basic checks while taking [Erlang runtime hostname resolver features](https://erlang.org/doc/apps/erts/inet_cfg.html) into account.  由于主机名解析是节点间成功通信的先决条件，从 RabbitMQ 3.8.6 开始，CLI 工具提供了两个命令来帮助验证节点上的主机名解析是否按预期工作。 这些命令并不是要取代 dig 和其他专门的 DNS 工具，而是提供一种方法来执行最基本的检查，同时考虑 Erlang 运行时主机名解析器功能。

The commands are covered in the [Networking guide](https://www.rabbitmq.com/networking.html#dns-verify-resolution).

#### CLI Tools

Starting with [version 3.8.6](https://www.rabbitmq.com/changelog.html), rabbitmq-diagnostics includes a command that provides relevant information on the Erlang cookie file used by CLI tools:  从 3.8.6 版本开始，rabbitmq-diagnostics 包含一个命令，该命令提供有关 CLI 工具使用的 Erlang cookie 文件的相关信息：

```bash
rabbitmq-diagnostics erlang_cookie_sources
```

The command will report on the effective user, user home directory and the expected location of the cookie file:

```plaintext
Cookie File

Effective user: antares
Effective home directory: /home/cli-user
Cookie file path: /home/cli-user/.erlang.cookie
Cookie file exists? true
Cookie file type: regular
Cookie file access: read
Cookie file size: 20

Cookie CLI Switch

--erlang-cookie value set? false
--erlang-cookie value length: 0

Env variable  (Deprecated)

RABBITMQ_ERLANG_COOKIE value set? false
RABBITMQ_ERLANG_COOKIE value length: 0
```

## Node Counts and Quorum

Because several features (e.g. [quorum queues](https://www.rabbitmq.com/quorum-queues.html), [client tracking in MQTT](https://www.rabbitmq.com/mqtt.html)) require a consensus between cluster members, odd numbers of cluster nodes are highly recommended: 1, 3, 5, 7 and so on.  由于几个特性（例如仲裁队列、MQTT 中的客户端跟踪）需要集群成员之间达成共识，因此强烈建议使用奇数个集群节点：1、3、5、7 等。

Two node clusters are **highly recommended against** since it's impossible for cluster nodes to identify a majority and form a consensus in case of connectivity loss. For example, when the two nodes lose connectivity MQTT client connections won't be accepted, quorum queues would lose their availability, and so on.  强烈建议不要使用两个节点集群，因为在连接丢失的情况下，集群节点不可能识别多数并形成共识。 例如，当两个节点失去连接时，MQTT 客户端连接将不被接受，仲裁队列将失去其可用性，等等。

From the consensus point of view, four or six node clusters would have the same availability characteristics as three and five node clusters.  从共识的角度来看，四节点或六节点集群将具有与三节点和五节点集群相同的可用性特征。

The [Quorum Queues guide](https://www.rabbitmq.com/quorum-queues.html) covers this topic in more detail.

## Clustering and Clients

Assuming all cluster members are available, a client can connect to any node and perform any operation. Nodes will route operations to the [quorum queue leader](https://www.rabbitmq.com/quorum-queues.html) or [queue leader replica](https://www.rabbitmq.com/ha.html#leader-migration-data-locality) transparently to clients.  假设所有集群成员都可用，客户端可以连接到任何节点并执行任何操作。 节点将对客户端透明地将操作路由到仲裁队列领导者或队列领导者副本。

With all supported messaging protocols a client is only connected to one node at a time.  对于所有支持的消息传递协议，客户端一次只能连接到一个节点。

In case of a node failure, clients should be able to reconnect to a different node, recover their topology and continue operation. For this reason, most client libraries accept a list of endpoints (hostnames or IP addresses) as a connection option. The list of hosts will be used during initial connection as well as connection recovery, if the client supports it. See documentation guides for individual clients to learn more.  如果节点发生故障，客户端应该能够重新连接到不同的节点，恢复其拓扑结构并继续运行。 因此，大多数客户端库都接受端点列表（主机名或 IP 地址）作为连接选项。 如果客户端支持，主机列表将在初始连接和连接恢复期间使用。 请参阅各个客户的文档指南以了解更多信息。

With [quorum queues](https://www.rabbitmq.com/quorum-queues.html), clients will only be able to perform operations on queues that have a quorum of replicas online.  使用仲裁队列，客户端将只能对具有仲裁在线副本的队列执行操作。

With classic mirrored queues, there are scenarios where it may not be possible for a client to transparently continue operations after connecting to a different node. They usually involve [non-mirrored queues hosted on a failed node](https://www.rabbitmq.com/ha.html#non-mirrored-queue-behavior-on-node-failure).  对于经典的镜像队列，在某些情况下，客户端可能无法在连接到不同节点后透明地继续操作。 它们通常涉及托管在故障节点上的非镜像队列。

## Clustering and Observability

Client connections, channels and queues will be distributed across cluster nodes. Operators need to be able to inspect and [monitor](https://www.rabbitmq.com/monitoring.html) such resources across all cluster nodes.  客户端连接、通道和队列将分布在集群节点上。 运营商需要能够跨所有集群节点检查和监控此类资源。

RabbitMQ [CLI tools](https://www.rabbitmq.com/cli.html) such as rabbitmq-diagnostics and rabbitmqctl provide commands that inspect resources and cluster-wide state. Some commands focus on the state of a single node (e.g. rabbitmq-diagnostics environment and rabbitmq-diagnostics status), others inspect cluster-wide state. Some examples of the latter include rabbitmqctl list_connections, rabbitmqctl list_mqtt_connections, rabbitmqctl list_stomp_connections, rabbitmqctl list_users, rabbitmqctl list_vhosts and so on.

Such "cluster-wide" commands will often contact one node first, discover cluster members and contact them all to retrieve and combine their respective state. For example, rabbitmqctl list_connections will contact all nodes, retrieve their AMQP 0-9-1 and AMQP 1.0 connections, and display them all to the user. The user doesn't have to manually contact all nodes. Assuming a non-changing state of the cluster (e.g. no connections are closed or opened), two CLI commands executed against two different nodes one after another will produce identical or semantically identical results. "Node-local" commands, however, will not produce identical results since two nodes rarely have identical state: at the very least their node names will be different!  这种“集群范围”的命令通常会首先联系一个节点，发现集群成员并联系他们全部以检索和组合他们各自的状态。 例如，rabbitmqctl list_connections 将联系所有节点，检索它们的 AMQP 0-9-1 和 AMQP 1.0 连接，并将它们全部显示给用户。 用户不必手动联系所有节点。 假设集群的状态不变（例如没有关闭或打开连接），针对两个不同节点一个接一个执行的两个 CLI 命令将产生相同或语义相同的结果。 然而，“节点本地”命令不会产生相同的结果，因为两个节点很少有相同的状态：至少它们的节点名称会不同！

[Management UI](https://www.rabbitmq.com/management.html) works similarly: a node that has to respond to an HTTP API request will fan out to other cluster members and aggregate their responses. In a cluster with multiple nodes that have management plugin enabled, the operator can use any node to access management UI. The same goes for monitoring tools that use the HTTP API to collect data about the state of the cluster. There is no need to issue a request to every cluster node in turn.

### Node Failure Handling

RabbitMQ brokers tolerate the failure of individual nodes. Nodes can be started and stopped at will, as long as they can contact a cluster member node known at the time of shutdown.  RabbitMQ 代理可以容忍单个节点的故障。 节点可以随意启动和停止，只要它们可以联系在关闭时已知的集群成员节点。

[Quorum queue](https://www.rabbitmq.com/quorum-queues.html) allows queue contents to be replicated across multiple cluster nodes with parallel replication and a predictable leader election and [data safety](https://www.rabbitmq.com/quorum-queues.html#data-safety) behavior as long as a majority of replicas are online.  仲裁队列允许队列内容跨多个集群节点复制，只要大多数副本在线，并行复制和可预测的领导者选举和数据安全行为。

Non-replicated classic queues can also be used in clusters. Non-mirrored queue [behaviour in case of node failure](https://www.rabbitmq.com/ha.html#non-mirrored-queue-behavior-on-node-failure) depends on [queue durability](https://www.rabbitmq.com/queues.html#durability).  非复制经典队列也可用于集群。 节点故障时的非镜像队列行为取决于队列持久性。

RabbitMQ clustering has several modes of dealing with [network partitions](https://www.rabbitmq.com/partitions.html), primarily consistency oriented. Clustering is meant to be used across LAN. It is not recommended to run clusters that span WAN. The [Shovel](https://www.rabbitmq.com/shovel.html) or [Federation](https://www.rabbitmq.com/federation.html) plugins are better solutions for connecting brokers across a WAN. Note that [Shovel and Federation are not equivalent to clustering](https://www.rabbitmq.com/distributed.html).  RabbitMQ 集群有几种处理网络分区的模式，主要是面向一致性的。 群集旨在跨 LAN 使用。 不建议运行跨 WAN 的集群。 Shovel 或 Federation 插件是通过 WAN 连接代理的更好解决方案。 请注意，Shovel 和 Federation 并不等同于聚类。

### Metrics and Statistics

Every node stores and aggregates its own metrics and stats, and provides an API for other nodes to access it. Some stats are cluster-wide, others are specific to individual nodes. Node that responds to an [HTTP API](https://www.rabbitmq.com/management.html) request contacts its peers to retrieve their data and then produces an aggregated result.

In versions older than 3.6.7, [RabbitMQ management plugin](https://www.rabbitmq.com/management.html) used a dedicated node for stats collection and aggregation.

## Clustering Transcript with rabbitmqctl

The following several sections provide a transcript of manually setting up and manipulating a RabbitMQ cluster across three machines: rabbit1, rabbit2, rabbit3. It is recommended that the example is studied before [more automation-friendly](https://www.rabbitmq.com/cluster-formation.html) cluster formation options are used.  以下几节提供了跨三台机器手动设置和操作 RabbitMQ 集群的记录：rabbit1、rabbit2、rabbit3。 建议在使用更适合自动化的集群形成选项之前研究该示例。

We assume that the user is logged into all three machines, that RabbitMQ has been installed on the machines, and that the rabbitmq-server and rabbitmqctl scripts are in the user's PATH.  我们假设用户登录到所有三台机器上，机器上已经安装了 RabbitMQ，并且 rabbitmq-server 和 rabbitmqctl 脚本在用户的 PATH 中。

This transcript can be modified to run on a single host, as explained more details below.  可以修改此脚本以在单个主机上运行，详情如下所述。

## Starting Independent Nodes

Clusters are set up by re-configuring existing RabbitMQ nodes into a cluster configuration. Hence the first step is to start RabbitMQ on all nodes in the normal way:

```bash
# on rabbit1
rabbitmq-server -detached

# on rabbit2
rabbitmq-server -detached

# on rabbit3
rabbitmq-server -detached
```

This creates three *independent* RabbitMQ brokers, one on each node, as confirmed by the *cluster_status* command:

```bash
# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1]}]},{running_nodes,[rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit2]}]},{running_nodes,[rabbit@rabbit2]}]
# => ...done.

# on rabbit3
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit3 ...
# => [{nodes,[{disc,[rabbit@rabbit3]}]},{running_nodes,[rabbit@rabbit3]}]
# => ...done.
```

The node name of a RabbitMQ broker started from the rabbitmq-server shell script is rabbit@*shorthostname*, where the short node name is lower-case (as in rabbit@rabbit1, above). On Windows, if rabbitmq-server.bat batch file is used, the short node name is upper-case (as in rabbit@RABBIT1). When you type node names, case matters, and these strings must match exactly.

## Creating a Cluster

In order to link up our three nodes in a cluster, we tell two of the nodes, say rabbit@rabbit2 and rabbit@rabbit3, to join the cluster of the third, say rabbit@rabbit1. Prior to that both newly joining members must be [reset](https://www.rabbitmq.com/rabbitmqctl.8.html#reset).

We first join rabbit@rabbit2 in a cluster with rabbit@rabbit1. To do that, on rabbit@rabbit2 we stop the RabbitMQ application and join the rabbit@rabbit1 cluster, then restart the RabbitMQ application. Note that a node must be [reset](https://www.rabbitmq.com/rabbitmqctl.8.html#reset) before it can join an existing cluster. Resetting the node **removes all resources and data that were previously present on that node**. This means that a node cannot be made a member of a cluster and keep its existing data at the same time. When that's desired, using the [Blue/Green deployment strategy](https://www.rabbitmq.com/blue-green-upgrade.html) or [backup and restore](https://www.rabbitmq.com/backup.html) are the available options.

```bash
# on rabbit2
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit2 ...done.

rabbitmqctl reset
# => Resetting node rabbit@rabbit2 ...

rabbitmqctl join_cluster rabbit@rabbit1
# => Clustering node rabbit@rabbit2 with [rabbit@rabbit1] ...done.

rabbitmqctl start_app
# => Starting node rabbit@rabbit2 ...done.
```

We can see that the two nodes are joined in a cluster by running the *cluster_status* command on either of the nodes:

```bash
# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2]}]},
# =>  {running_nodes,[rabbit@rabbit1,rabbit@rabbit2]}]
# => ...done.
```

Now we join rabbit@rabbit3 to the same cluster. The steps are identical to the ones above, except this time we'll cluster to rabbit2 to demonstrate that the node chosen to cluster to does not matter - it is enough to provide one online node and the node will be clustered to the cluster that the specified node belongs to.

```bash
# on rabbit3
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit3 ...done.

# on rabbit3
rabbitmqctl reset
# => Resetting node rabbit@rabbit3 ...

rabbitmqctl join_cluster rabbit@rabbit2
# => Clustering node rabbit@rabbit3 with rabbit@rabbit2 ...done.

rabbitmqctl start_app
# => Starting node rabbit@rabbit3 ...done.
```

We can see that the three nodes are joined in a cluster by running the cluster_status command on any of the nodes:

```bash
# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit3,rabbit@rabbit2,rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit3,rabbit@rabbit1,rabbit@rabbit2]}]
# => ...done.

# on rabbit3
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit3 ...
# => [{nodes,[{disc,[rabbit@rabbit3,rabbit@rabbit2,rabbit@rabbit1]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit1,rabbit@rabbit3]}]
# => ...done.
```

By following the above steps we can add new nodes to the cluster at any time, while the cluster is running.

## Restarting Cluster Nodes

Nodes that have been joined to a cluster can be stopped at any time. They can also fail or be terminated by the OS.  已加入集群的节点可以随时停止。 它们也可能失败或被操作系统终止。

In general, if the majority of nodes is still online after a node is stopped, this does not affect the rest of the cluster, although client connection distribution, queue replica placement, and load distribution of the cluster will change.  一般来说，如果一个节点停止后大部分节点仍然在线，这不会影响集群的其余部分，尽管集群的客户端连接分布、队列副本放置和负载分布会发生变化。

### Schema Syncing from Online Peers  来自在线对等点的模式同步

A restarted node will sync the schema and other information from its peers on boot. Before this process completes, the node **won't be fully started and functional**.  重新启动的节点将在启动时从其对等方同步模式和其他信息。 在此过程完成之前，节点不会完全启动和运行。

It is therefore important to understand the process node go through when they are stopped and restarted.  因此，了解进程节点在停止和重新启动时经历的过程非常重要。

A stopping node picks an online cluster member (only disc nodes will be considered) to sync with after restart. Upon restart the node will try to contact that peer 10 times by default, with 30 second response timeouts.  停止节点选择一个在线集群成员（只考虑磁盘节点）在重启后同步。 重启后，默认情况下，节点将尝试联系该对等方 10 次，响应超时为 30 秒。

In case the peer becomes available in that time interval, the node successfully starts, syncs what it needs from the peer and keeps going.  如果对等点在该时间间隔内可用，则节点成功启动，从对等点同步所需的内容并继续运行。

If the peer does not become available, the restarted node will **give up and voluntarily stop**. Such condition can be identified by the timeout (timeout_waiting_for_tables) warning messages in the logs that eventually lead to node startup failure:  如果对等点不可用，重新启动的节点将放弃并自动停止。 这种情况可以通过日志中的超时（timeout_waiting_for_tables）警告消息来识别，最终导致节点启动失败：

```plaintext
2020-07-27 21:10:51.361 [warning] <0.269.0> Error while waiting for Mnesia tables: {timeout_waiting_for_tables,[rabbit@node2,rabbit@node1],[rabbit_durable_queue]}
2020-07-27 21:10:51.361 [info] <0.269.0> Waiting for Mnesia tables for 30000 ms, 1 retries left
2020-07-27 21:11:21.362 [warning] <0.269.0> Error while waiting for Mnesia tables: {timeout_waiting_for_tables,[rabbit@node2,rabbit@node1],[rabbit_durable_queue]}
2020-07-27 21:11:21.362 [info] <0.269.0> Waiting for Mnesia tables for 30000 ms, 0 retries left
2020-07-27 21:15:51.380 [info] <0.269.0> Waiting for Mnesia tables for 30000 ms, 1 retries left
2020-07-27 21:16:21.381 [warning] <0.269.0> Error while waiting for Mnesia tables: {timeout_waiting_for_tables,[rabbit@node2,rabbit@node1],[rabbit_user,rabbit_user_permission, …]}
2020-07-27 21:16:21.381 [info] <0.269.0> Waiting for Mnesia tables for 30000 ms, 0 retries left
2020-07-27 21:16:51.393 [info] <0.44.0> Application mnesia exited with reason: stopped
2020-07-27 21:16:51.397 [error] <0.269.0> BOOT FAILED
2020-07-27 21:16:51.397 [error] <0.269.0> ===========
2020-07-27 21:16:51.397 [error] <0.269.0> Timeout contacting cluster nodes: [rabbit@node1].
```

When a node has no online peers during shutdown, it will start without attempts to sync with any known peers. It does not start as a standalone node, however, and peers will be able to rejoin it.  当一个节点在关闭期间没有在线对等点时，它将在不尝试与任何已知对等点同步的情况下启动。 但是，它不是作为独立节点启动的，对等节点将能够重新加入它。

When the entire cluster is brought down therefore, the last node to go down is the only one that didn't have any running peers at the time of shutdown. That node can start without contacting any peers first. Since nodes will try to contact a known peer for up to 5 minutes (by default), nodes can be restarted in any order in that period of time. In this case they will rejoin each other one by one successfully. This window of time can be adjusted using two configuration settings:  因此，当整个集群关闭时，最后一个关闭的节点是唯一一个在关闭时没有任何正在运行的对等节点。 该节点可以在不首先联系任何对等方的情况下启动。 由于节点将尝试联系已知对等方长达 5 分钟（默认情况下），因此可以在该时间段内以任何顺序重新启动节点。 在这种情况下，他们将成功地一一重聚。 可以使用两个配置设置来调整此时间窗口：

```ini
# wait for 60 seconds instead of 30
mnesia_table_loading_retry_timeout = 60000

# retry 15 times instead of 10
mnesia_table_loading_retry_limit = 15
```

By adjusting these settings and tweaking the time window in which known peer has to come back it is possible to account for cluster-wide redeployment scenarios that can be longer than 5 minutes to complete.  通过调整这些设置并调整已知对等点必须返回的时间窗口，可以解决可能需要 5 分钟以上才能完成的集群范围重新部署方案。

During [upgrades](https://www.rabbitmq.com/upgrade.html), sometimes the last node to stop must be the first node to be started after the upgrade. That node will be designated to perform a cluster-wide schema migration that other nodes can sync from and apply when they rejoin.  在升级过程中，有时最后一个停止的节点必须是升级后要启动的第一个节点。 该节点将被指定执行集群范围的架构迁移，其他节点可以从中同步并在它们重新加入时应用。

### Restarts and Health Checks (Readiness Probes)

In some environments, node restarts are controlled with a designated [health check](https://www.rabbitmq.com/monitoring.html#health-checks). The checks verify that one node has started and the deployment process can proceed to the next one. If the check does not pass, the deployment of the node is considered to be incomplete and the deployment process will typically wait and retry for a period of time. One popular example of such environment is Kubernetes where an operator-defined [readiness probe](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-readiness-gate) can prevent a deployment from proceeding when the [OrderedReady pod management policy](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#deployment-and-scaling-guarantees) is used. Deployments that use the Parallel pod management policy will not be affected but must worry about the [natural race condition during initial cluster formation](https://www.rabbitmq.com/cluster-formation.html#initial-formation-race-condition).  在某些环境中，节点重新启动由指定的健康检查控制。 这些检查验证一个节点是否已启动，并且部署过程可以继续进行到下一个节点。 如果检查不通过，则认为节点的部署不完整，部署过程通常会等待并重试一段时间。 这种环境的一个流行示例是 Kubernetes，当使用 OrderedReady pod 管理策略时，操作员定义的就绪探针可以阻止部署继续进行。 使用 Parallel Pod 管理策略的部署不会受到影响，但必须担心初始集群形成期间的自然竞争条件。

Given the [peer syncing behavior described above](https://www.rabbitmq.com/clustering.html#restarting-schema-sync), such a health check can prevent a cluster-wide restart from completing in time. Checks that explicitly or implicitly assume a fully booted node that's rejoined its cluster peers will fail and block further node deployments.  鉴于上述对等同步行为，这样的健康检查可以防止集群范围的重启及时完成。 检查明确或隐含地假设重新加入其集群对等点的完全启动节点将失败并阻止进一步的节点部署。

[Most health check](https://www.rabbitmq.com/monitoring.html#health-checks), even relatively basic ones, implicitly assume that the node has finished booting. They are not suitable for nodes that are [awaiting schema table sync](https://www.rabbitmq.com/clustering.html#restarting-schema-sync) from a peer.

One very common example of such check is

```bash
# will exit with an error for the nodes that are currently waiting for
# a peer to sync schema tables from
rabbitmq-diagnostics check_running
```

One health check that does not expect a node to be fully booted and have schema tables synced is

```bash
# a very basic check that will succeed for the nodes that are currently waiting for
# a peer to sync schema from
rabbitmq-diagnostics ping
```

This basic check would allow the deployment to proceed and the nodes to eventually rejoin each other, assuming they are [compatible](https://www.rabbitmq.com/upgrade.html).

### Hostname Changes Between Restarts

A node rejoining after a node name or host name change can start as [a blank node](https://www.rabbitmq.com/cluster-formation.html#peer-discovery-how-does-it-work) if its data directory path changes as a result. Such nodes will fail to rejoin the cluster. While the node is offline, its peers can be reset or started with a blank data directory. In that case the recovering node will fail to rejoin its peer as well since internal data store cluster identity would no longer match.  如果节点名称或主机名称更改后重新加入的节点的数据目录路径因此更改，则该节点可以作为空白节点开始。 此类节点将无法重新加入集群。 当节点离线时，它的对等节点可以被重置或使用空白数据目录启动。 在这种情况下，恢复节点也将无法重新加入其对等节点，因为内部数据存储集群身份将不再匹配。

Consider the following scenario:

1. A cluster of 3 nodes, A, B and C is formed
2. Node A is shut down
3. Node B is reset
4. Node A is started
5. Node A tries to rejoin B but B's cluster identity has changed
6. Node B doesn't recognise A as a known cluster member because it's been reset

in this case node B will reject the clustering attempt from A with an appropriate error message in the log:

```plaintext
Node 'rabbit@node1.local' thinks it's clustered with node 'rabbit@node2.local', but 'rabbit@node2.local' disagrees
```

In this case B can be reset again and then will be able to join A, or A can be reset and will successfully join B.

### Cluster Node Restart Example

The below example uses CLI tools to shut down the nodes rabbit@rabbit1 and rabbit@rabbit3 and check on the cluster status at each step:

```bash
# on rabbit1
rabbitmqctl stop
# => Stopping and halting node rabbit@rabbit1 ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit3,rabbit@rabbit2]}]
# => ...done.

# on rabbit3
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit3 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit3]}]
# => ...done.

# on rabbit3
rabbitmqctl stop
# => Stopping and halting node rabbit@rabbit3 ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit2]}]
# => ...done.
```

In the below example, the nodes are started back, checking on the cluster status as we go along:

```bash
# on rabbit1
rabbitmq-server -detached
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit1,rabbit@rabbit2]}]
# => ...done.

# on rabbit3
rabbitmq-server -detached

# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit1,rabbit@rabbit3]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]
# => ...done.

# on rabbit3
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit3 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2,rabbit@rabbit3]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit1,rabbit@rabbit3]}]
# => ...done.
```

## Forcing Node Boot in Case of Unavailable Peers

In some cases the last node to go offline cannot be brought back up. It can be removed from the cluster using the forget_cluster_node [rabbitmqctl](https://www.rabbitmq.com/cli.html) command.  在某些情况下，最后一个离线的节点无法恢复。 可以使用forget_cluster_node rabbitmqctl 命令从集群中删除它。

Alternatively force_boot [rabbitmqctl](https://www.rabbitmq.com/cli.html) command can be used on a node to make it boot without trying to sync with any peers (as if they were last to shut down). This is usually only necessary if the last node to shut down or a set of nodes will never be brought back online.  或者，可以在节点上使用 force_boot rabbitmqctl 命令使其启动，而无需尝试与任何对等方同步（就像它们最后一次关闭一样）。 这通常只有在最后一个节点关闭或一组节点永远不会重新联机时才需要。

## [Breaking Up a Cluster](https://www.rabbitmq.com/clustering.html#removing-nodes)

Sometimes it is necessary to remove a node from a cluster. The operator has to do this explicitly using a rabbitmqctl command.

Some [peer discovery mechanisms](https://www.rabbitmq.com/cluster-formation.html) support node health checks and forced removal of nodes not known to the discovery backend. That feature is opt-in (disabled by default).

We first remove rabbit@rabbit3 from the cluster, returning it to independent operation. To do that, on rabbit@rabbit3 we stop the RabbitMQ application, reset the node, and restart the RabbitMQ application.

```bash
# on rabbit3
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit3 ...done.

rabbitmqctl reset
# => Resetting node rabbit@rabbit3 ...done.
rabbitmqctl start_app
# => Starting node rabbit@rabbit3 ...done.
```

Note that it would have been equally valid to list rabbit@rabbit3 as a node.

Running the *cluster_status* command on the nodes confirms that rabbit@rabbit3 now is no longer part of the cluster and operates independently:

```bash
# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2]}]},
# => {running_nodes,[rabbit@rabbit2,rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1,rabbit@rabbit2]}]},
# =>  {running_nodes,[rabbit@rabbit1,rabbit@rabbit2]}]
# => ...done.

# on rabbit3
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit3 ...
# => [{nodes,[{disc,[rabbit@rabbit3]}]},{running_nodes,[rabbit@rabbit3]}]
# => ...done.
```

We can also remove nodes remotely. This is useful, for example, when having to deal with an unresponsive node. We can for example remove rabbit@rabbit1 from rabbit@rabbit2.

```bash
# on rabbit1
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit1 ...done.

# on rabbit2
rabbitmqctl forget_cluster_node rabbit@rabbit1
# => Removing node rabbit@rabbit1 from cluster ...
# => ...done.
```

Note that rabbit1 still thinks it's clustered with rabbit2, and trying to start it will result in an error. We will need to reset it to be able to start it again.

```bash
# on rabbit1
rabbitmqctl start_app
# => Starting node rabbit@rabbit1 ...
# => Error: inconsistent_cluster: Node rabbit@rabbit1 thinks it's clustered with node rabbit@rabbit2, but rabbit@rabbit2 disagrees

rabbitmqctl reset
# => Resetting node rabbit@rabbit1 ...done.

rabbitmqctl start_app
# => Starting node rabbit@rabbit1 ...
# => ...done.
```

The cluster_status command now shows all three nodes operating as independent RabbitMQ brokers:

```bash
# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1]}]},{running_nodes,[rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit2]}]},{running_nodes,[rabbit@rabbit2]}]
# => ...done.

# on rabbit3
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit3 ...
# => [{nodes,[{disc,[rabbit@rabbit3]}]},{running_nodes,[rabbit@rabbit3]}]
# => ...done.
```

Note that rabbit@rabbit2 retains the residual state of the cluster, whereas rabbit@rabbit1 and rabbit@rabbit3 are freshly initialised RabbitMQ brokers. If we want to re-initialise rabbit@rabbit2 we follow the same steps as for the other nodes:

```bash
# on rabbit2
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit2 ...done.
rabbitmqctl reset
# => Resetting node rabbit@rabbit2 ...done.
rabbitmqctl start_app
# => Starting node rabbit@rabbit2 ...done.
```

Besides rabbitmqctl forget_cluster_node and the automatic cleanup of unknown nodes by some [peer discovery](https://www.rabbitmq.com/cluster-formation.html) plugins, there are no scenarios in which a RabbitMQ node will permanently remove its peer node from a cluster.

### [How to Reset a Node](https://www.rabbitmq.com/clustering.html#resetting-nodes)

Sometimes it may be necessary to reset a node (wipe all of its data) and later make it rejoin the cluster. Generally speaking, there are two possible scenarios: when the node is running, and when the node cannot start or won't respond to CLI tool commands e.g. due to an issue such as [ERL-430](https://bugs.erlang.org/browse/ERL-430).

Resetting a node will delete all of its data, cluster membership information, configured [runtime parameters](https://www.rabbitmq.com/parameters.html), users, virtual hosts and any other node data. It will also permanently remove the node from its cluster.

To reset a running and responsive node, first stop RabbitMQ on it using rabbitmqctl stop_app and then reset it using rabbitmqctl reset:

```bash
# on rabbit1
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit1 ...done.
rabbitmqctl reset
# => Resetting node rabbit@rabbit1 ...done.
```

In case of a non-responsive node, it must be stopped first using any means necessary. For nodes that fail to start this is already the case. Then [override](https://www.rabbitmq.com/relocate.html) the node's data directory location or [re]move the existing data store. This will make the node start as a blank one. It will have to be instructed to [rejoin its original cluster](https://www.rabbitmq.com/clustering.html#cluster-formation), if any.

A node that's been reset and rejoined its original cluster will sync all virtual hosts, users, permissions and topology (queues, exchanges, bindings), runtime parameters and policies. [Quorum queue](https://www.rabbitmq.com/quorum-queues.html) contents will be replicated if the node will be selected to host a replica. Non-replicated queue contents on a reset node will be lost.

## [Upgrading clusters](https://www.rabbitmq.com/clustering.html#upgrading)

You can find instructions for upgrading a cluster in [the upgrade guide](https://www.rabbitmq.com/upgrade.html#rabbitmq-cluster-configuration).

## [A Cluster on a Single Machine](https://www.rabbitmq.com/clustering.html#single-machine)

Under some circumstances it can be useful to run a cluster of RabbitMQ nodes on a single machine. This would typically be useful for experimenting with clustering on a desktop or laptop without the overhead of starting several virtual machines for the cluster.

In order to run multiple RabbitMQ nodes on a single machine, it is necessary to make sure the nodes have distinct node names, data store locations, log file locations, and bind to different ports, including those used by plugins. See RABBITMQ_NODENAME, RABBITMQ_NODE_PORT, and RABBITMQ_DIST_PORT in the [Configuration guide](https://www.rabbitmq.com/configure.html#supported-environment-variables), as well as RABBITMQ_MNESIA_DIR, RABBITMQ_CONFIG_FILE, and RABBITMQ_LOG_BASE in the [File and Directory Locations guide](https://www.rabbitmq.com/relocate.html).

You can start multiple nodes on the same host manually by repeated invocation of rabbitmq-server ( rabbitmq-server.bat on Windows). For example:

```bash
RABBITMQ_NODE_PORT=5672 RABBITMQ_NODENAME=rabbit rabbitmq-server -detached
RABBITMQ_NODE_PORT=5673 RABBITMQ_NODENAME=hare rabbitmq-server -detached
rabbitmqctl -n hare stop_app
rabbitmqctl -n hare join_cluster rabbit@`hostname -s`
rabbitmqctl -n hare start_app
```

will set up a two node cluster, both nodes as disc nodes. Note that if the node [listens on any ports](https://www.rabbitmq.com/networking.html) other than AMQP 0-9-1 and AMQP 1.0 ones, those must be configured to avoid a collision as well. This can be done via command line:

```bash
RABBITMQ_NODE_PORT=5672 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15672}]" RABBITMQ_NODENAME=rabbit rabbitmq-server -detached
RABBITMQ_NODE_PORT=5673 RABBITMQ_SERVER_START_ARGS="-rabbitmq_management listener [{port,15673}]" RABBITMQ_NODENAME=hare rabbitmq-server -detached
```

will start two nodes (which can then be clustered) when the management plugin is installed.

## [Hostname Changes](https://www.rabbitmq.com/clustering.html#issues-hostname)

RabbitMQ nodes use hostnames to communicate with each other. Therefore, all node names must be able to resolve names of all cluster peers. This is also true for tools such as rabbitmqctl.

In addition to that, by default RabbitMQ names the database directory using the current hostname of the system. If the hostname changes, a new empty database is created. To avoid data loss it's crucial to set up a fixed and resolvable hostname.

Whenever the hostname changes RabbitMQ node must be restarted.

A similar effect can be achieved by using rabbit@localhost as the broker nodename. The impact of this solution is that clustering will not work, because the chosen hostname will not resolve to a routable address from remote hosts. The rabbitmqctl command will similarly fail when invoked from a remote host. A more sophisticated solution that does not suffer from this weakness is to use DNS, e.g. [Amazon Route 53](http://aws.amazon.com/route53/) if running on EC2. If you want to use the full hostname for your nodename (RabbitMQ defaults to the short name), and that full hostname is resolvable using DNS, you may want to investigate setting the environment variable RABBITMQ_USE_LONGNAME=true.

See the section on [hostname resolution](https://www.rabbitmq.com/clustering.html#hostname-resolution-requirement) for more information.

## [Firewalled Nodes](https://www.rabbitmq.com/clustering.html#firewall)

Nodes can have a firewall enabled on them. In such case, traffic on certain ports must be allowed by the firewall in both directions, or nodes won't be able to join each other and perform all the operations they expect to be available on cluster peers.

Learn more in the [section on ports](https://www.rabbitmq.com/clustering.html#ports) above and dedicated [RabbitMQ Networking guide](https://www.rabbitmq.com/networking.html).

## [Erlang Versions Across the Cluster](https://www.rabbitmq.com/clustering.html#erlang)

All nodes in a cluster are *highly recommended* to run the same major [version of Erlang](https://www.rabbitmq.com/which-erlang.html): 22.2.0 and 22.2.8 can be mixed but 21.3.6 and 22.2.6 can potentially introduce breaking changes in inter-node communication protocols. While such breaking changes are relatively rare, they are possible.

Incompatibilities between patch releases of Erlang/OTP versions are very rare.

## [Connecting to Clusters from Clients](https://www.rabbitmq.com/clustering.html#clients)

A client can connect as normal to any node within a cluster. If that node should fail, and the rest of the cluster survives, then the client should notice the closed connection, and should be able to reconnect to some surviving member of the cluster.

Many clients support lists of hostnames that will be tried in order at connection time.

Generally it is not recommended to hardcode IP addresses into client applications: this introduces inflexibility and will require client applications to be edited, recompiled and redeployed should the configuration of the cluster change or the number of nodes in the cluster change.

Instead, consider a more abstracted approach: this could be a dynamic DNS service which has a very short TTL configuration, or a plain TCP load balancer, or a combination of them.

In general, this aspect of managing the connection to nodes within a cluster is beyond the scope of this guide, and we recommend the use of other technologies designed specifically to address these problems.

## [Disk and RAM Nodes](https://www.rabbitmq.com/clustering.html#cluster-node-types)

A node can be a *disk node* or a *RAM node*. (**Note:** *disk* and *disc* are used interchangeably). RAM nodes store internal database tables in RAM only. This does not include messages, message store indices, queue indices and other node state.

In the vast majority of cases you want all your nodes to be disk nodes; RAM nodes are a special case that can be used to improve the performance clusters with high queue, exchange, or binding churn. RAM nodes do not provide higher message rates. When in doubt, use disk nodes only.

Since RAM nodes store internal database tables in RAM only, they must sync them from a peer node on startup. This means that a cluster must contain at least one disk node. It is therefore not possible to manually remove the last remaining disk node in a cluster.

## [Clusters with RAM nodes](https://www.rabbitmq.com/clustering.html#ram-nodes)

RAM nodes keep their metadata only in memory. As RAM nodes don't have to write to disc as much as disc nodes, they can perform better. However, note that since persistent queue data is always stored on disc, the performance improvements will affect only resource management (e.g. adding/removing queues, exchanges, or vhosts), but not publishing or consuming speed.

RAM nodes are an advanced use case; when setting up your first cluster you should simply not use them. You should have enough disc nodes to handle your redundancy requirements, then if necessary add additional RAM nodes for scale.

A cluster containing only RAM nodes would be too volatile; if the cluster stops you will not be able to start it again and **will lose all data**. RabbitMQ will prevent the creation of a RAM-node-only cluster in many situations, but it can't absolutely prevent it.

The examples here show a cluster with one disc and one RAM node for simplicity only; such a cluster is a poor design choice.

### [Creating RAM nodes](https://www.rabbitmq.com/clustering.html#creating-ram)

We can declare a node as a RAM node when it first joins the cluster. We do this with rabbitmqctl join_cluster as before, but passing the --ram flag:

```bash
# on rabbit2
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit2 ...done.

rabbitmqctl join_cluster --ram rabbit@rabbit1
# => Clustering node rabbit@rabbit2 with [rabbit@rabbit1] ...done.

rabbitmqctl start_app
# => Starting node rabbit@rabbit2 ...done.
```

RAM nodes are shown as such in the cluster status:

```bash
# on rabbit1
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit1 ...
# => [{nodes,[{disc,[rabbit@rabbit1]},{ram,[rabbit@rabbit2]}]},
# =>  {running_nodes,[rabbit@rabbit2,rabbit@rabbit1]}]
# => ...done.

# on rabbit2
rabbitmqctl cluster_status
# => Cluster status of node rabbit@rabbit2 ...
# => [{nodes,[{disc,[rabbit@rabbit1]},{ram,[rabbit@rabbit2]}]},
# =>  {running_nodes,[rabbit@rabbit1,rabbit@rabbit2]}]
# => ...done.
```

### [Changing node types](https://www.rabbitmq.com/clustering.html#change-type)

We can change the type of a node from ram to disc and vice versa. Say we wanted to reverse the types of rabbit@rabbit2 and rabbit@rabbit1, turning the former from a ram node into a disc node and the latter from a disc node into a ram node. To do that we can use the change_cluster_node_type command. The node must be stopped first.

```bash
# on rabbit2
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit2 ...done.

rabbitmqctl change_cluster_node_type disc
# => Turning rabbit@rabbit2 into a disc node ...done.

rabbitmqctl start_app
# => Starting node rabbit@rabbit2 ...done.

# on rabbit1
rabbitmqctl stop_app
# => Stopping node rabbit@rabbit1 ...done.

rabbitmqctl change_cluster_node_type ram
# => Turning rabbit@rabbit1 into a ram node ...done.

rabbitmqctl start_app
# => Starting node rabbit@rabbit1 ...done.
```

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!