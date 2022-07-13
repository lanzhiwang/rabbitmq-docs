# Networking and RabbitMQ

https://www.rabbitmq.com/networking.html

## Overview

Clients communicate with RabbitMQ over the network. All protocols supported by the broker are TCP-based. Both RabbitMQ and the operating system provide a number of knobs that can be tweaked. Some of them are directly related to TCP and IP operations, others have to do with application-level protocols such as TLS. This guide covers multiple topics related to networking in the context of RabbitMQ. This guide is not meant to be an extensive reference but rather an overview. Some tuneable parameters discussed are OS-specific. This guide focuses on Linux when covering OS-specific subjects, as it is the most common platform RabbitMQ is deployed on.  客户端通过网络与 RabbitMQ 通信。代理支持的所有协议都是基于 TCP 的。 RabbitMQ 和操作系统都提供了许多可以调整的旋钮。其中一些与 TCP 和 IP 操作直接相关，另一些与应用层协议（如 TLS）有关。本指南涵盖了与 RabbitMQ 上下文中的网络相关的多个主题。本指南不是广泛的参考，而是一个概述。讨论的一些可调参数是特定于操作系统的。本指南在涵盖特定于操作系统的主题时侧重于 Linux，因为它是部署 RabbitMQ 的最常见平台。

There are several areas which can be configured or tuned. Each has a section in this guide:  有几个区域可以配置或调整。每个在本指南中都有一个部分：

- Interfaces the node listens on for client connections  节点监听客户端连接的接口

- IP version preferences: dual stack, IPv6-only and IPv4-only  IP 版本首选项：双栈、仅 IPv6 和仅 IPv4

- Ports used by clients, inter-node traffic in clusters and [CLI tools](https://www.rabbitmq.com/cli.html)  客户端使用的端口、集群中的节点间流量和 CLI 工具

- IPv6 support for inter-node traffic  IPv6 支持节点间流量

- TLS for client connections  用于客户端连接的 TLS

- Tuning for a large number of concurrent connections  调整大量并发连接

- High connection churn scenarios and resource exhaustion  高连接流失场景和资源耗尽

- TCP buffer size (affects throughput and how much memory is used per connection)  TCP 缓冲区大小（影响吞吐量和每个连接使用多少内存）

- Hostname resolution-related topics such as reverse DNS lookups  与主机名解析相关的主题，例如反向 DNS 查找

- The interface and port used by epmd  epmd使用的接口和端口

- How to suspend and resume listeners to temporarily stop and resume new client connections  如何暂停和恢复侦听器以临时停止和恢复新的客户端连接

- Other TCP socket settings  其他 TCP 套接字设置

- Proxy protocol support for client connections  客户端连接的代理协议支持

- Kernel TCP settings and limits (e.g. TCP keepalives and open file handle limit)  内核 TCP 设置和限制（例如 TCP keepalives 和打开文件句柄限制）

- How to allow Erlang runtime to accept inbound connections when MacOS Application Firewall is enabled  启用 MacOS 应用程序防火墙时如何允许 Erlang 运行时接受入站连接

This guide also covers a few topics closely related to networking:  本指南还涵盖了一些与网络密切相关的主题：

Except for OS kernel parameters and DNS, all RabbitMQ settings are [configured via RabbitMQ configuration file(s)](https://www.rabbitmq.com/configure.html).  除了操作系统内核参数和 DNS，所有 RabbitMQ 设置都是通过 RabbitMQ 配置文件配置的。

Networking is a broad topic. There are many configuration options that can have positive or negative effect on certain workloads. As such, this guide does not try to be a complete reference but rather offer an index of key tunable parameters and serve as a starting point.  网络是一个广泛的话题。有许多配置选项会对某些工作负载产生积极或消极的影响。因此，本指南并不试图成为完整的参考，而是提供关键可调参数的索引并作为起点。

In addition, this guide touches on a few topics closely related to networking, such as  此外，本指南还涉及一些与网络密切相关的主题，例如

- Hostnames, hostname resolution and DNS  主机名、主机名解析和 DNS

- connection lifecycle logging  连接生命周期日志

- Heartbeats (a.k.a. keepalives)  心跳（又名keepalives）

- proxies and load balancers  代理和负载均衡器

[Tanzu RabbitMQ](https://www.rabbitmq.com/tanzu) provides an [inter-node traffic compression](https://www.rabbitmq.com/clustering-compression.html) feature.  Tanzu RabbitMQ 提供了节点间流量压缩功能。

A methodology for [troubleshooting of networking-related issues](https://www.rabbitmq.com/troubleshooting-networking.html) is covered in a separate guide.  单独的指南中介绍了网络相关问题的故障排除方法。

## Network Interfaces for Client Connections  客户端连接的网络接口

For RabbitMQ to accept client connections, it needs to bind to one or more interfaces and listen on (protocol-specific) ports. One such interface/port pair is called a listener in RabbitMQ parlance. Listeners are configured using the **listeners.tcp.*** configuration option(s).  为了让 RabbitMQ 接受客户端连接，它需要绑定到一个或多个接口并监听（特定于协议的）端口。一个这样的接口/端口对在 RabbitMQ 用语中称为侦听器。使用 listeners.tcp.* 配置选项配置监听器。

TCP listeners configure both an interface and port. The following example demonstrates how to configure AMQP 0-9-1 and AMQP 1.0 listener to use a specific IP and the standard port:  TCP 侦听器配置接口和端口。以下示例演示如何配置 AMQP 0-9-1 和 AMQP 1.0 侦听器以使用特定 IP 和标准端口：

```bash
listeners.tcp.1 = 192.168.1.99:5672
```

By default, RabbitMQ will listen on port 5672 on **all available interfaces**. It is possible to limit client connections to a subset of the interfaces or even just one, for example, IPv6-only interfaces. The following few sections demonstrate how to do it.  默认情况下，RabbitMQ 将在所有可用接口上侦听端口 5672。 可以将客户端连接限制到接口的一个子集，或者甚至只是一个，例如，仅限 IPv6 的接口。 以下几节演示了如何做到这一点。

### Listening on Dual Stack (Both IPv4 and IPv6) Interfaces  监听双栈（IPv4 和 IPv6）接口

The following example demonstrates how to configure RabbitMQ to listen on localhost only for both IPv4 and IPv6:  以下示例演示了如何将 RabbitMQ 配置为仅在 localhost 上侦听 IPv4 和 IPv6：

```bash
listeners.tcp.1 = 127.0.0.1:5672
listeners.tcp.2 = ::1:5672
```

With modern Linux kernels and Windows releases, when a port is specified and RabbitMQ is configured to listen on all IPv6 addresses but IPv4 is not disabled explicitly, IPv4 address will be included, so  在现代 Linux 内核和 Windows 版本中，当指定端口并且将 RabbitMQ 配置为侦听所有 IPv6 地址但未明确禁用 IPv4 时，将包含 IPv4 地址，因此

```bash
listeners.tcp.1 = :::5672
```

is equivalent to

```bash
listeners.tcp.1 = 0.0.0.0:5672
listeners.tcp.2 = :::5672
```

### Listening on IPv6 Interfaces Only

In this example RabbitMQ will listen on an IPv6 interface only:

```bash
listeners.tcp.1 = fe80::2acf:e9ff:fe17:f97b:5672
```

In IPv6-only environments the node must also be configured to use IPv6 for inter-node communication and CLI tool connections.  在纯 IPv6 环境中，还必须将节点配置为使用 IPv6 进行节点间通信和 CLI 工具连接。

### Listening on IPv4 Interfaces Only  仅侦听 IPv4 接口

In this example RabbitMQ will listen on an IPv4 interface only:  在此示例中，RabbitMQ 将仅侦听 IPv4 接口：

```bash
listeners.tcp.1 = 192.168.1.99:5672
```

It is possible to disable non-TLS connections by disabling all regular TCP listeners. Only [TLS-enabled](https://www.rabbitmq.com/ssl.html) clients will be able to connect:  可以通过禁用所有常规 TCP 侦听器来禁用非 TLS 连接。 只有启用 TLS 的客户端才能连接：

```bash
# disables non-TLS listeners, only TLS-enabled clients will be able to connect
listeners.tcp = none

listeners.ssl.default = 5671

ssl_options.cacertfile = /path/to/ca_certificate.pem
ssl_options.certfile   = /path/to/server_certificate.pem
ssl_options.keyfile    = /path/to/server_key.pem
ssl_options.verify     = verify_peer
ssl_options.fail_if_no_peer_cert = false
```

## Port Access

RabbitMQ nodes bind to ports (open server TCP sockets) in order to accept client and CLI tool connections. Other processes and tools such as SELinux may prevent RabbitMQ from binding to a port. When that happens, the node will fail to start.  RabbitMQ 节点绑定到端口（打开服务器 TCP 套接字）以接受客户端和 CLI 工具连接。 SELinux 等其他进程和工具可能会阻止 RabbitMQ 绑定到端口。发生这种情况时，节点将无法启动。

CLI tools, client libraries and RabbitMQ nodes also open connections (client TCP sockets). Firewalls can prevent nodes and CLI tools from communicating with each other. Make sure the following ports are accessible:  CLI 工具、客户端库和 RabbitMQ 节点也打开连接（客户端 TCP 套接字）。防火墙可以阻止节点和 CLI 工具相互通信。确保可以访问以下端口：

- 4369: [epmd](http://erlang.org/doc/man/epmd.html), a peer discovery service used by RabbitMQ nodes and CLI tools  RabbitMQ 节点和 CLI 工具使用的对等发现服务

- 5672, 5671: used by AMQP 0-9-1 and AMQP 1.0 clients without and with TLS  由不带和带 TLS 的 AMQP 0-9-1 和 AMQP 1.0 客户端使用

- 5552, 5551: used by the [RabbitMQ Stream protocol](https://www.rabbitmq.com/streams.html) clients without and with TLS  由 RabbitMQ 流协议客户端使用，不带和带 TLS

- 6000 through 6500 (usually 6000, 6001, 6002, and so on through 6005): used by [RabbitMQ Stream](https://www.rabbitmq.com/stream.html) replication  由 RabbitMQ 流复制使用

- 25672: used for inter-node and CLI tools communication (Erlang distribution server port) and is allocated from a dynamic range (limited to a single port by default, computed as AMQP port + 20000). Unless external connections on these ports are really necessary (e.g. the cluster uses [federation](https://www.rabbitmq.com/federation.html) or CLI tools are used on machines outside the subnet), these ports should not be publicly exposed. See networking guide for details.  用于节点间和 CLI 工具通信（Erlang 分发服务器端口），从动态范围分配（默认限制为单个端口，计算为 AMQP 端口 + 20000）。除非确实需要这些端口上的外部连接（例如，集群使用联合或 CLI 工具在子网外的机器上使用），否则这些端口不应公开。有关详细信息，请参阅网络指南。

- 35672-35682: used by CLI tools (Erlang distribution client ports) for communication with nodes and is allocated from a dynamic range (computed as server distribution port + 10000 through server distribution port + 10010). See networking guide for details.  由 CLI 工具（Erlang 分发客户端端口）用于与节点通信，并从动态范围分配（计算为服务器分发端口 + 10000 到服务器分发端口 + 10010）。有关详细信息，请参阅网络指南。

- 15672, 15671: [HTTP API](https://www.rabbitmq.com/management.html) clients, [management UI](https://www.rabbitmq.com/management.html) and [rabbitmqadmin](https://www.rabbitmq.com/management-cli.html), without and with TLS (only if the [management plugin](https://www.rabbitmq.com/management.html) is enabled)  HTTP API 客户端、管理 UI 和 rabbitmqadmin，不带和带 TLS（仅在启用管理插件的情况下）

- 61613, 61614: [STOMP clients](https://stomp.github.io/stomp-specification-1.2.html) without and with TLS (only if the [STOMP plugin](https://www.rabbitmq.com/stomp.html) is enabled)  不带和带 TLS 的 STOMP 客户端（仅在启用 STOMP 插件的情况下）

- 1883, 8883: [MQTT clients](http://mqtt.org/) without and with TLS, if the [MQTT plugin](https://www.rabbitmq.com/mqtt.html) is enabled  如果启用了 MQTT 插件，则不带和带 TLS 的 MQTT 客户端

- 15674: STOMP-over-WebSockets clients (only if the [Web STOMP plugin](https://www.rabbitmq.com/web-stomp.html) is enabled)  STOMP-over-WebSockets 客户端（仅当 Web STOMP 插件启用时）

- 15675: MQTT-over-WebSockets clients (only if the [Web MQTT plugin](https://www.rabbitmq.com/web-mqtt.html) is enabled)  MQTT-over-WebSockets 客户端（仅当启用 Web MQTT 插件时）

- 15692: Prometheus metrics (only if the [Prometheus plugin](https://www.rabbitmq.com/prometheus.html) is enabled)  Prometheus 指标（仅在启用 Prometheus 插件的情况下）

It is possible to [configure RabbitMQ](https://www.rabbitmq.com/configure.html) to use different ports and specific network interfaces.  可以将 RabbitMQ 配置为使用不同的端口和特定的网络接口。

## How to Temporarily Stop New Client Connections  如何临时停止新的客户端连接

Starting with RabbitMQ 3.8.8, client connection listeners can be *suspended* to prevent new client connections from being accepted. Existing connections will not be affected in any way.  从 RabbitMQ 3.8.8 开始，可以暂停客户端连接侦听器以防止接受新的客户端连接。现有连接不会受到任何影响。

This can be useful during node operations and is one of the steps performed when a node is [put into maintenance mode](https://www.rabbitmq.com/upgrade.html#maintenance-mode).  这在节点操作期间很有用，并且是节点进入维护模式时执行的步骤之一。

To suspend all listeners on a node and prevent new client connections to it, use **rabbitmqctl suspend_listeners**:  要挂起节点上的所有侦听器并阻止新的客户端连接到它，请使用 rabbitmqctl suspend_listeners：

```bash
rabbitmqctl suspend_listeners
```

As all other CLI commands, this command can be invoked against an arbitrary node (including remote ones) using the **-n** switch:  与所有其他 CLI 命令一样，可以使用 -n 开关对任意节点（包括远程节点）调用此命令：

```bash
# suspends listeners on node rabbit@node2.cluster.rabbitmq.svc: it won't accept any new client connections
rabbitmqctl suspend_listeners -n rabbit@node2.cluster.rabbitmq.svc
```

To resume all listeners on a node and make it accept new client connections again, use **rabbitmqctl resume_listeners**:  要恢复节点上的所有侦听器并使其再次接受新的客户端连接，请使用 rabbitmqctl resume_listeners：

```bash
rabbitmqctl resume_listeners
# resumes listeners on node rabbit@node2.cluster.rabbitmq.svc: it will accept new client connections again
rabbitmqctl resume_listeners -n rabbit@node2.cluster.rabbitmq.svc
```

Both operations will leave [log entries](https://www.rabbitmq.com/logging.html) in the node's log.

## EPMD and Inter-node Communication  EPMD 和节点间通信

### What is EPMD and How is It Used?  什么是 EPMD，它是如何使用的？

[epmd](http://www.erlang.org/doc/man/epmd.html) (for Erlang Port Mapping Daemon) is a small additional daemon that runs alongside every RabbitMQ node and is used by the [runtime](https://www.rabbitmq.com/runtime.html) to discover what port a particular node listens on for inter-node communication. The port is then used by peer nodes and [CLI tools](https://www.rabbitmq.com/cli.html).  epmd（用于 Erlang 端口映射守护程序）是一个小型附加守护程序，它与每个 RabbitMQ 节点一起运行，运行时使用它来发现特定节点侦听哪个端口以进行节点间通信。 然后，对等节点和 CLI 工具会使用该端口。

When a node or CLI tool needs to contact node rabbit@hostname2 it will do the following:  当节点或 CLI 工具需要联系节点 rabbit@hostname2 时，它将执行以下操作：

- Resolve hostname2 to an IPv4 or IPv6 address using the standard OS resolver or a custom one specified in the [inetrc file](http://erlang.org/doc/apps/erts/inet_cfg.html)  使用标准操作系统解析器或 inetrc 文件中指定的自定义解析器将 hostname2 解析为 IPv4 或 IPv6 地址

- Contact epmd running on hostname2 using the above address  使用上述地址联系在 hostname2 上运行的 epmd

- Ask epmd for the port used by node rabbit on it  向 epmd 询问 node rabbit 在其上使用的端口

- Connect to the node using the resolved IP address and the discovered port  使用解析的 IP 地址和发现的端口连接到节点

- Proceed with communication  继续沟通

### EPMD Interface  EPMD 接口

epmd will listen on all interfaces by default. It can be limited to a number of interfaces using the ERL_EPMD_ADDRESS environment variable:  默认情况下，epmd 将侦听所有接口。 它可以限制为使用 ERL_EPMD_ADDRESS 环境变量的多个接口：

```bash
# makes epmd listen on loopback IPv6 and IPv4 interfaces
export ERL_EPMD_ADDRESS="::1"
```

When ERL_EPMD_ADDRESS is changed, both RabbitMQ node and epmd on the host must be stopped. For epmd, use  当 ERL_EPMD_ADDRESS 改变时，RabbitMQ 节点和主机上的 epmd 都必须停止。 对于 epmd，使用

```bash
# Stops local epmd process.
# Use after shutting down RabbitMQ.
epmd -kill
```

to terminate it. The service will be started by the local RabbitMQ node automatically on boot.  终止它。 该服务将在启动时由本地 RabbitMQ 节点自动启动。

The loopback interface will be implicitly added to that list (in other words, epmd will always bind to the loopback interface).  loopback 接口将被隐式添加到该列表中（换句话说，epmd 将始终绑定到 loopback 接口）。

### EPMD Port  EPMD 端口

The default epmd port is 4369, but this can be changed using the ERL_EPMD_PORT environment variable:  默认的 epmd 端口是 4369，但可以使用 ERL_EPMD_PORT 环境变量进行更改：

```bash
# makes epmd bind to port 4369
export ERL_EPMD_PORT="4369"
```

All hosts in a [cluster](https://www.rabbitmq.com/clustering.html) must use the same port.  集群中的所有主机必须使用相同的端口。

When ERL_EPMD_PORT is changed, both RabbitMQ node and epmd on the host must be stopped. For epmd, use  当 ERL_EPMD_PORT 更改时，必须停止主机上的 RabbitMQ 节点和 epmd。 对于 epmd，使用

```bash
# Stops local epmd process.
# Use after shutting down RabbitMQ.
epmd -kill
```

to terminate it. The service will be started by the local RabbitMQ node automatically on boot.  终止它。该服务将在启动时由本地 RabbitMQ 节点自动启动。

### Inter-node Communication Port Range  节点间通信端口范围

RabbitMQ nodes will use a port from a certain range known as the inter-node communication port range. The same port is used by CLI tools when they need to contact the node. The range can be modified.  RabbitMQ 节点将使用某个范围内的端口，该范围称为节点间通信端口范围。 CLI 工具在需要联系节点时使用相同的端口。范围可以修改。

RabbitMQ nodes communicate with CLI tools and other nodes using a port known as the **distribution port**. It is dynamically allocated from a range of values. For RabbitMQ, the default range is limited to a single value computed as RABBITMQ_NODE_PORT (AMQP 0-9-1 and AMQP 1.0 port) + 20000, which results in using port 25672. This single port can be [configured](https://www.rabbitmq.com/configure.html) using the RABBITMQ_DIST_PORT environment variable.  RabbitMQ 节点使用称为分发端口的端口与 CLI 工具和其他节点通信。它是从一系列值中动态分配的。对于 RabbitMQ，默认范围仅限于计算为 RABBITMQ_NODE_PORT（AMQP 0-9-1 和 AMQP 1.0 端口）+ 20000 的单个值，这导致使用端口 25672。可以使用 RABBITMQ_DIST_PORT 环境变量配置这个单个端口。

RabbitMQ [command line tools](https://www.rabbitmq.com/cli.html) also use a range of ports. The default range is computed by taking the RabbitMQ distribution port value and adding 10000 to it. The next 10 ports are also part of this range. Thus, by default, this range is 35672 through 35682. This range can be configured using the RABBITMQ_CTL_DIST_PORT_MIN and RABBITMQ_CTL_DIST_PORT_MAX environment variables. Note that limiting the range to a single port will prevent more than one CLI tool from running concurrently on the same host and may affect CLI commands that require parallel connections to multiple cluster nodes. A port range of 10 is therefore a recommended value.  RabbitMQ 命令行工具也使用一系列端口。默认范围是通过获取 RabbitMQ 分发端口值并将 10000 添加到它来计算的。接下来的 10 个端口也在此范围内。因此，默认情况下，此范围为 35672 到 35682。可以使用 RABBITMQ_CTL_DIST_PORT_MIN 和 RABBITMQ_CTL_DIST_PORT_MAX 环境变量配置此范围。请注意，将范围限制为单个端口将阻止多个 CLI 工具在同一主机上同时运行，并且可能会影响需要并行连接到多个集群节点的 CLI 命令。因此，建议使用端口范围 10。

When configuring firewall rules it is highly recommended to allow remote connections on the inter-node communication port from every cluster member and every host where CLI tools might be used. epmd port must be open for CLI tools and clustering to function.  在配置防火墙规则时，强烈建议允许来自每个集群成员和可能使用 CLI 工具的每个主机的节点间通信端口上的远程连接。必须打开 epmd 端口，CLI 工具和集群才能运行。

On Windows, the following settings have no effect when RabbitMQ runs as a service. Please see [Windows Quirks](https://www.rabbitmq.com/windows-quirks.html) for details.  在 Windows 上，当 RabbitMQ 作为服务运行时，以下设置无效。有关详细信息，请参阅 Windows 怪癖。

The range used by RabbitMQ can also be controlled via two configuration keys:  RabbitMQ 使用的范围也可以通过两个配置键来控制：

- kernel.inet_dist_listen_min in the **classic** config format *only*  仅采用经典配置格式

- kernel.inet_dist_listen_max in the **classic** config format *only*  仅采用经典配置格式

They define the range's lower and upper bounds, inclusive.  它们定义了范围的下限和上限，包括在内。

The example below uses a range with a single port but a value different from default:  下面的示例使用具有单个端口但值不同于默认值的范围：

```erlang
[
  {kernel, [
    {inet_dist_listen_min, 33672},
    {inet_dist_listen_max, 33672}
  ]},
  {rabbit, [
    ...
  ]}
].
```

To verify what port is used by a node for inter-node and CLI tool communication, run  要验证节点使用哪个端口进行节点间和 CLI 工具通信，请运行

```bash
epmd -names
```

on that node's host. It will produce output that looks like this:

```ini
epmd: up and running on port 4369 with data:
name rabbit at port 25672
```

### Inter-node Communication Buffer Size Limit  节点间通信缓冲区大小限制

Inter-node connections use a buffer for data pending to be sent. Temporary throttling on inter-node traffic is applied when the buffer is at max allowed capacity. The limit is controlled via the **RABBITMQ_DISTRIBUTION_BUFFER_SIZE** [environment variable](https://www.rabbitmq.com/configure.html#supported-environment-variables) in kilobytes. Default value is 128 MB (128000 kB).  节点间连接使用缓冲区来存储待发送的数据。当缓冲区处于最大允许容量时，将应用对节点间流量的临时限制。该限制通过 RABBITMQ_DISTRIBUTION_BUFFER_SIZE 环境变量以千字节为单位进行控制。默认值为 128 MB (128000 kB)。

In clusters with heavy inter-node traffic increasing this value may have a positive effect on throughput. Values lower than 64 MB are not recommended.  在节点间流量较大的集群中，增加此值可能会对吞吐量产生积极影响。不建议使用低于 64 MB 的值。

## Using IPv6 for Inter-node Communication (and CLI Tools)  使用 IPv6 进行节点间通信（和 CLI 工具）

In addition to exclusive IPv6 use for client connections for client connections, a node can also be configured to use IPv6 exclusively for inter-node and CLI tool connectivity.  除了将 IPv6 用于客户端连接的客户端连接之外，还可以将节点配置为将 IPv6 专用于节点间和 CLI 工具连接。

This involves configuration in a few places:  这涉及到几个地方的配置：

- Inter-node communication protocol setting in the [runtime](https://www.rabbitmq.com/runtime.html)  运行时节点间通信协议设置

- Configuring IPv6 to be used by CLI tools  配置 CLI 工具使用的 IPv6

- epmd, a service involved in inter-node communication (discovery)  涉及节点间通信（发现）的服务

It is possible to use IPv6 for inter-node and CLI tool communication but use IPv4 for client connections or vice versa. Such configurations can be hard to troubleshoot and reason about, so using the same IP version (e.g. IPv6) across the board or a dual stack setup is recommended.  可以使用 IPv6 进行节点间和 CLI 工具通信，但使用 IPv4 进行客户端连接，反之亦然。此类配置可能难以排除故障和推理，因此建议全面使用相同的 IP 版本（例如 IPv6）或双堆栈设置。

### Inter-node Communication Protocol  节点间通信协议

To instruct the runtime to use IPv6 for inter-node communication and related tasks, use the **RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS** environment variable to pass a couple of flags:  要指示运行时使用 IPv6 进行节点间通信和相关任务，请使用 RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS 环境变量传递几个标志：

```bash
# these flags will be used by RabbitMQ nodes
RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS="-kernel inetrc '/etc/rabbitmq/erl_inetrc' -proto_dist inet6_tcp"
# these flags will be used by CLI tools
RABBITMQ_CTL_ERL_ARGS="-proto_dist inet6_tcp"
```

RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS above uses two closely related flags:  上面的 RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS 使用了两个密切相关的标志：

- **-kernel inetrc** to configure a path to an [inetrc file](http://erlang.org/doc/apps/erts/inet_cfg.html) that controls hostname resolution  配置控制主机名解析的 inetrc 文件的路径

- **-proto_dist inet6_tcp** to tell the node to use IPv6 when connecting to peer nodes and listening for CLI tool connections  告诉节点在连接到对等节点并侦听 CLI 工具连接时使用 IPv6

The **erl_inetrc** file at **/etc/rabbitmq/erl_inetrc** will control hostname resolution settings. For IPv6-only environments, it must include the following line:

/etc/rabbitmq/erl_inetrc 中的 erl_inetrc 文件将控制主机名解析设置。 对于纯 IPv6 环境，它必须包含以下行：

```bash
%% Tells DNS client on RabbitMQ nodes and CLI tools to resolve hostnames to IPv6 addresses.
%% The trailing dot is not optional.
{inet6,true}.
```

### CLI Tools

With CLI tools, use the same runtime flag as used for RabbitMQ nodes above but provide it using a different environment variable, **RABBITMQ_CTL_ERL_ARGS**:  使用 CLI 工具，使用与上述 RabbitMQ 节点相同的运行时标志，但使用不同的环境变量 RABBITMQ_CTL_ERL_ARGS 提供它：

```bash
RABBITMQ_CTL_ERL_ARGS="-proto_dist inet6_tcp"
```

Note that once instructed to use IPv6, CLI tools won't be able to connect to nodes that do not use IPv6 for inter-node communication. This involves the epmd service running on the same host as target RabbitMQ node.  请注意，一旦指示使用 IPv6，CLI 工具将无法连接到不使用 IPv6 进行节点间通信的节点。 这涉及与目标 RabbitMQ 节点在同一主机上运行的 epmd 服务。

### epmd

epmd is a small helper daemon that runs next to a RabbitMQ node and lets its peers and CLI tools discover what port they should use to communicate to it. It can be configured to bind to a specific interface, much like RabbitMQ listeners. This is done using the **ERL_EPMD_ADDRESS** environment variable:  epmd 是一个小型助手守护程序，它在 RabbitMQ 节点旁边运行，并让其对等节点和 CLI 工具发现它们应该使用哪个端口与之通信。 它可以配置为绑定到特定接口，就像 RabbitMQ 侦听器一样。 这是使用 ERL_EPMD_ADDRESS 环境变量完成的：

```bash
export ERL_EPMD_ADDRESS="::1"
```

By default RabbitMQ nodes will use an IPv4 interface when connecting to epmd. Nodes that are configured to use IPv6 for inter-node communication (see above) will also use IPv6 to connect to epmd.  默认情况下，RabbitMQ 节点在连接到 epmd 时将使用 IPv4 接口。 配置为使用 IPv6 进行节点间通信的节点（见上文）也将使用 IPv6 连接到 epmd。

When epmd is configured to use IPv6 exclusively but RabbitMQ nodes are not, RabbitMQ will log an error message similar to this:  当 epmd 配置为独占使用 IPv6 但 RabbitMQ 节点不是时，RabbitMQ 将记录类似以下的错误消息：

```bash
Protocol 'inet_tcp': register/listen error: econnrefused
```

#### systemd Unit File

On distributions that use systemd, the **epmd.socket service** controls network settings of epmd. It is possible to configure epmd to only listen on IPv6 interfaces:  在使用 systemd 的发行版上，epmd.socket 服务控制 epmd 的网络设置。 可以将 epmd 配置为仅侦听 IPv6 接口：

```bash
ListenStream=[::1]:4369
```

The service will need reloading after its unit file has been updated:

```bash
systemctl daemon-reload
systemctl restart epmd.socket epmd.service
```

## Intermediaries: Proxies and Load Balancers  中介：代理和负载均衡器

Proxies and load balancers are fairly commonly used to distribute client connections between [cluster nodes](https://www.rabbitmq.com/clustering.html). Proxies can also be useful to make it possible for clients to access RabbitMQ nodes without exposing them publicly. Intermediaries can also have side effects on connections.  代理和负载平衡器通常用于在集群节点之间分配客户端连接。代理也可以用于使客户端可以访问 RabbitMQ 节点而不公开它们。中介也可能对连接产生副作用。

### Proxy Effects  代理效果

Proxies and load balancers introduce an extra network hop (or even multiple ones) between client and its target node. Intermediaries also can become a network contention point: their throughput will then become a limiting factor for the entire system. Network bandwidth overprovisioning and throughput monitoring for proxies and load balancers are therefore very important.  代理和负载均衡器在客户端和它的目标节点之间引入了一个额外的网络跃点（甚至多个）。中介也可能成为网络竞争点：它们的吞吐量将成为整个系统的限制因素。因此，代理和负载均衡器的网络带宽过度配置和吞吐量监控非常重要。

Intermediaries also may terminate "idle" TCP connections when there's no activity on them for a certain period of time. Most of the time it is not desirable. Such events will result in [abrupt connection closure log messages](https://www.rabbitmq.com/logging.html#connection-lifecycle-events) on the server end and I/O exceptions on the client end.  中介也可以在一段时间内没有活动时终止“空闲” TCP 连接。大多数时候这是不可取的。此类事件将导致服务器端突然的连接关闭日志消息和客户端的 I/O 异常。

When [heartbeats](https://www.rabbitmq.com/heartbeats.html) are enabled on a connection, it results in periodic light network traffic. Therefore heartbeats have a side effect of guarding client connections that can go idle for periods of time against premature closure by proxies and load balancers.  在连接上启用心跳时，会导致周期性的少量网络流量。因此，心跳具有保护客户端连接的副作用，客户端连接可以空闲一段时间以防止代理和负载平衡器过早关闭。

Heartbeat timeouts from 10 to 30 seconds will produce periodic network traffic often enough (roughly every 5 to 15 seconds) to satisfy defaults of most proxy tools and load balancers. Values that are too low will produce false positives.  10 到 30 秒的心跳超时将产生足够频繁的周期性网络流量（大约每 5 到 15 秒）以满足大多数代理工具和负载平衡器的默认设置。值太低会产生误报。

### Proxy Protocol

RabbitMQ supports [Proxy protocol](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) versions 1 (text header format) and 2 (binary header format).  RabbitMQ 支持代理协议版本 1（文本头格式）和 2（二进制头格式）。

The protocol makes servers such as RabbitMQ aware of the actual client IP address when connections go over a proxy (e.g. [HAproxy](http://cbonte.github.io/haproxy-dconv/1.8/configuration.html#send-proxy) or [AWS ELB](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-target-groups.html#proxy-protocol)). This makes it easier for the operator to inspect connection origins in the management UI or CLI tools.  当连接通过代理（例如 HAproxy 或 AWS ELB）时，该协议使 RabbitMQ 等服务器知道实际的客户端 IP 地址。这使操作员更容易在管理 UI 或 CLI 工具中检查连接来源。

The protocol spec dictates that either it must be applied to all connections or none of them for security reasons, this feature is disabled by default and needs to be enabled for individual protocols supported by RabbitMQ. To enable it for AMQP 0-9-1 and AMQP 1.0 clients:

协议规范规定，出于安全原因，它必须应用于所有连接或不应用于所有连接，此功能默认禁用，需要为 RabbitMQ 支持的各个协议启用。要为 AMQP 0-9-1 和 AMQP 1.0 客户端启用它：

```bash
proxy_protocol = true
```

When proxy protocol is enabled, clients won't be able to connect to RabbitMQ directly unless they themselves support the protocol. Therefore, when this option is enabled, all client connections must go through a proxy that also supports the protocol and is configured to send a Proxy protocol header. [HAproxy](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) and [AWS ELB](https://docs.aws.amazon.com/elasticloadbalancing/latest/network/load-balancer-target-groups.html#proxy-protocol) documentation explains how to do it.  启用代理协议后，客户端将无法直接连接到 RabbitMQ，除非它们自己支持该协议。因此，启用此选项后，所有客户端连接都必须通过也支持该协议并配置为发送 Proxy 协议标头的代理。 HAproxy 和 AWS ELB 文档解释了如何做到这一点。

When proxy protocol is enabled and connections go through a compatible proxy, no action or modifications are required from client libraries. The communication is entirely transparent to them.  当启用代理协议并且连接通过兼容的代理时，客户端库不需要任何操作或修改。沟通对他们来说是完全透明的。

[STOMP](https://www.rabbitmq.com/stomp.html#proxy-protocol) and [MQTT](https://www.rabbitmq.com/mqtt.html#proxy-protocol), as well as [Web STOMP](https://www.rabbitmq.com/web-stomp.html#proxy-protocol) and [Web MQTT](https://www.rabbitmq.com/web-mqtt.html#proxy-protocol) have their own settings that enable support for the proxy protocol.  STOMP 和 MQTT，以及 Web STOMP 和 Web MQTT 都有自己的设置，可以支持代理协议。

## TLS (SSL) Support

It is possible to encrypt connections using TLS with RabbitMQ. Authentication using peer certificates is also possible. Please refer to the [TLS/SSL guide](https://www.rabbitmq.com/ssl.html) for more information.  可以使用带有 RabbitMQ 的 TLS 加密连接。也可以使用对等证书进行身份验证。有关详细信息，请参阅 TLS/SSL 指南。

## Tuning for Throughput  调整吞吐量

Tuning for throughput is a common goal. Improvements can be achieved by  调整吞吐量是一个共同的目标。可以通过以下方式进行改进

- Increasing TCP buffer sizes  增加 TCP 缓冲区大小

- Ensuring Nagle's algorithm is disabled  确保禁用 Nagle 算法

- Enabling optional TCP features and extensions  启用可选的 TCP 功能和扩展

For the latter two, see the OS-level tuning section below.  对于后两者，请参阅下面的操作系统级别调整部分。

Note that tuning for throughput will involve trade-offs. For example, increasing TCP buffer sizes will increase the amount of RAM used by every connection, which can be a significant total server RAM use increase.  请注意，吞吐量调整将涉及权衡。例如，增加 TCP 缓冲区大小将增加每个连接使用的 RAM 量，这可能会显着增加服务器 RAM 的总使用量。

### TCP Buffer Size

This is one of the key tunable parameters. Every TCP connection has buffers allocated for it. Generally speaking, the larger these buffers are, the more RAM is used per connection and better the throughput. On Linux, the OS will automatically tune TCP buffer size by default, typically settling on a value between 80 and 120 KB.  这是关键的可调参数之一。每个 TCP 连接都有为其分配的缓冲区。一般来说，这些缓冲区越大，每个连接使用的 RAM 越多，吞吐量也越高。在 Linux 上，默认情况下操作系统会自动调整 TCP 缓冲区大小，通常设置在 80 到 120 KB 之间。

For maximum throughput, it is possible to increase buffer size using a group of config options:  为了获得最大吞吐量，可以使用一组配置选项来增加缓冲区大小：

- **tcp_listen_options** for AMQP 0-9-1 and AMQP 1.0
AMQP 0-9-1 和 AMQP 1.0 的 tcp_listen_options

- **mqtt.tcp_listen_options** for MQTT
MQTT 的 mqtt.tcp_listen_options

- **stomp.tcp_listen_options** for STOMP
STOMP 的 stomp.tcp_listen_options

Note that increasing TCP buffer size will increase how much [RAM the node uses](https://www.rabbitmq.com/memory-use.html) for every client connection.  请注意，增加 TCP 缓冲区大小将增加节点用于每个客户端连接的 RAM 量。

The following example sets TCP buffers for AMQP 0-9-1 connections to 192 KiB:  以下示例将 AMQP 0-9-1 连接的 TCP 缓冲区设置为 192 KiB：

```bash
tcp_listen_options.backlog = 128
tcp_listen_options.nodelay = true
tcp_listen_options.linger.on = true
tcp_listen_options.linger.timeout = 0
tcp_listen_options.sndbuf = 196608
tcp_listen_options.recbuf = 196608
```

The same example for MQTT:

```bash
mqtt.tcp_listen_options.backlog = 128
mqtt.tcp_listen_options.nodelay = true
mqtt.tcp_listen_options.linger.on = true
mqtt.tcp_listen_options.linger.timeout = 0
mqtt.tcp_listen_options.sndbuf = 196608
mqtt.tcp_listen_options.recbuf = 196608
```

and STOMP:

```bash
stomp.tcp_listen_options.backlog = 128
stomp.tcp_listen_options.nodelay = true
stomp.tcp_listen_options.linger.on = true
stomp.tcp_listen_options.linger.timeout = 0
stomp.tcp_listen_options.sndbuf = 196608
stomp.tcp_listen_options.recbuf = 196608
```

Note that setting send and receive buffer sizes to different values can be dangerous and **not recommended**.  请注意，将发送和接收缓冲区大小设置为不同的值可能很危险，不建议这样做。

## Tuning for a Large Number of Connections  调整大量连接

Some workloads, often referred to as "the Internet of Things", assume a large number of client connections per node, and a relatively low volume of traffic from each node. One such workload is sensor networks: there can be hundreds of thousands or millions of sensors deployed, each emitting data every several minutes. Optimising for the maximum number of concurrent clients can be more important than for total throughput.  一些工作负载，通常被称为“物联网”，假设每个节点有大量的客户端连接，并且每个节点的流量相对较低。其中一种工作负载是传感器网络：可能部署了数十万或数百万个传感器，每个传感器每隔几分钟就会发出一次数据。优化并发客户端的最大数量可能比总吞吐量更重要。

Several factors can limit how many concurrent connections a single node can support:  有几个因素可以限制单个节点可以支持多少并发连接：

- Maximum number of open file handles (including sockets) as well as other kernel-enforced resource limits  打开文件句柄（包括套接字）的最大数量以及其他内核强制资源限制

- Amount of [RAM used by each connection](https://www.rabbitmq.com/memory-use.html)  每个连接使用的 RAM 量

- Amount of CPU resources used by each connection  每个连接使用的 CPU 资源量

- Maximum number of Erlang processes the VM is configured to allow.  虚拟机配置允许的最大 Erlang 进程数。

### Open File Handle Limit  打开文件句柄限制

Most operating systems limit the number of file handles that can be opened at the same time. When an OS process (such as RabbitMQ's Erlang VM) reaches the limit, it won't be able to open any new files or accept any more TCP connections.  大多数操作系统限制可以同时打开的文件句柄的数量。当操作系统进程（例如 RabbitMQ 的 Erlang VM）达到限制时，它将无法打开任何新文件或接受更多 TCP 连接。

How the limit is configured [varies from OS to OS](https://github.com/basho/basho_docs/blob/master/content/riak/kv/2.2.3/using/performance/open-files-limit.md) and distribution to distribution, e.g. depending on whether systemd is used. For Linux, Controlling System Limits on Linux in our [Debian](https://www.rabbitmq.com/install-debian.html#kernel-resource-limits) and [RPM](https://www.rabbitmq.com/install-rpm.html#kernel-resource-limits) installation guides provides. Linux kernel limit management is covered by many resources on the Web, including the [open file handle limit](https://ro-che.info/articles/2017-03-26-increase-open-files-limit).  限制的配置方式因操作系统和分发而异，例如取决于是否使用 systemd。对于 Linux，我们的 Debian 和 RPM 安装指南中提供了控制 Linux 上的系统限制。 Linux 内核限制管理被 Web 上的许多资源所覆盖，包括打开文件句柄限制。

With Docker, [Docker daemon configuration file](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file) in the host controls the limits.  使用 Docker，主机中的 Docker 守护程序配置文件控制限制。

MacOS uses a [similar system](https://superuser.com/questions/433746/is-there-a-fix-for-the-too-many-open-files-in-system-error-on-os-x-10-7-1).  MacOS 使用类似的系统。

On Windows, the limit for the Erlang runtime is controlled using the ERL_MAX_PORTS environment variable.  在 Windows 上，Erlang 运行时的限制是使用 ERL_MAX_PORTS 环境变量控制的。

When optimising for the number of concurrent connections, make sure your system has enough file descriptors to support not only client connections but also files the node may use. To calculate a ballpark limit, multiply the number of connections per node by 1.5. For example, to support 100,000 connections, set the limit to 150,000.  在优化并发连接数时，请确保您的系统有足够的文件描述符来支持客户端连接以及节点可能使用的文件。要计算大致限制，请将每个节点的连接数乘以 1.5。例如，要支持 100,000 个连接，请将限制设置为 150,000。

Increasing the limit slightly increases the amount of RAM idle machine uses but this is a reasonable trade-off.  增加限制会略微增加 RAM 空闲机器的使用量，但这是一个合理的权衡。

### Per Connection Memory Consumption: TCP Buffer Size  每个连接的内存消耗：TCP 缓冲区大小

See the section above for an overview.  有关概述，请参见上面的部分。

For maximum number of concurrent client connections, it is possible to decrease TCP buffer size using a group of config options:  对于最大并发客户端连接数，可以使用一组配置选项来减小 TCP 缓冲区大小：

- **tcp_listen_options** for AMQP 0-9-1 and AMQP 1.0
AMQP 0-9-1 和 AMQP 1.0 的 tcp_listen_options

- **mqtt.tcp_listen_options** for MQTT
MQTT 的 mqtt.tcp_listen_options

- **stomp.tcp_listen_options** for STOMP
STOMP 的 stomp.tcp_listen_options

Decreasing TCP buffer size will decrease how much [RAM the node uses](https://www.rabbitmq.com/memory-use.html) for every client connection.  减小 TCP 缓冲区大小将减少节点用于每个客户端连接的 RAM 量。

This is often necessary in environments where the number of concurrent connections sustained per node is more important than throughput.  在每个节点维持的并发连接数比吞吐量更重要的环境中，这通常是必要的。

The following example sets TCP buffers for AMQP 0-9-1 connections to 32 KiB:  以下示例将 AMQP 0-9-1 连接的 TCP 缓冲区设置为 32 KiB：

```bash
tcp_listen_options.backlog = 128
tcp_listen_options.nodelay = true
tcp_listen_options.linger.on = true
tcp_listen_options.linger.timeout = 0
tcp_listen_options.sndbuf  = 32768
tcp_listen_options.recbuf  = 32768
```

The same example for MQTT:

```bash
mqtt.tcp_listen_options.backlog = 128
mqtt.tcp_listen_options.nodelay = true
mqtt.tcp_listen_options.linger.on = true
mqtt.tcp_listen_options.linger.timeout = 0
mqtt.tcp_listen_options.sndbuf  = 32768
mqtt.tcp_listen_options.recbuf  = 32768
```

and for STOMP:

```bash
stomp.tcp_listen_options.backlog = 128
stomp.tcp_listen_options.nodelay = true
stomp.tcp_listen_options.linger.on = true
stomp.tcp_listen_options.linger.timeout = 0
stomp.tcp_listen_options.sndbuf  = 32768
stomp.tcp_listen_options.recbuf  = 32768
```

Note that lowering TCP buffer sizes will result in a proportional throughput drop, so an optimal value between throughput and per-connection RAM use needs to be found for every workload.  请注意，降低 TCP 缓冲区大小将导致吞吐量按比例下降，因此需要为每个工作负载找到吞吐量和每个连接 RAM 使用之间的最佳值。

Setting send and receive buffer sizes to different values is dangerous and is not recommended. Values lower than 8 KiB are not recommended.  将发送和接收缓冲区大小设置为不同的值是危险的，不建议这样做。不建议使用低于 8 KiB 的值。

### Reducing CPU Footprint of Stats Emission  减少 Stats Emission 的 CPU 占用

A large number of concurrent connections will generate a lot of metric (stats) emission events. This increases CPU consumption even with mostly idle connections. To reduce this footprint, increase the statistics collection interval using the **collect_statistics_interval** key:  大量的并发连接会产生大量的metric（stats）发射事件。即使在大部分空闲连接的情况下，这也会增加 CPU 消耗。要减少此占用，请使用 collect_statistics_interval 键增加统计信息收集间隔：

```bash
# sets the interval to 60 seconds
collect_statistics_interval = 60000
```

The default is 5 seconds (5000 milliseconds).  默认值为 5 秒（5000 毫秒）。

Increasing the interval value to 30-60s will reduce CPU footprint and peak memory consumption. This comes with a downside: with the value in the example above, metrics of said entities will refresh every 60 seconds.  将间隔值增加到 30-60 秒将减少 CPU 占用和峰值内存消耗。这有一个缺点：使用上面示例中的值，所述实体的指标将每 60 秒刷新一次。

This can be perfectly reasonable in an [externally monitored](https://www.rabbitmq.com/monitoring.html#monitoring-frequency) production system but will make management UI less convenient to use for operators.  这在外部监控的生产系统中是完全合理的，但会使操作员使用管理 UI 不太方便。

### Limiting Number of Channels on a Connection  限制连接上的通道数

Channels also consume RAM. By optimising how many channels applications use, that amount can be decreased. It is possible to cap the max number of channels on a connection using the **channel_max** configuration setting:  通道也消耗 RAM。 通过优化应用程序使用的通道数量，可以减少该数量。 可以使用 channel_max 配置设置来限制连接上的最大通道数：

```bash
channel_max = 16
```

Note that some libraries and tools that build on top of RabbitMQ clients may implicitly require a certain number of channels. Values above 200 are rarely necessary. Finding an optimal value is usually a matter of trial and error.  请注意，一些构建在 RabbitMQ 客户端之上的库和工具可能隐含地需要一定数量的通道。 很少需要高于 200 的值。 找到一个最佳值通常是一个反复试验的问题。

### Nagle's Algorithm ("nodelay")

Disabling [Nagle's algorithm](http://en.wikipedia.org/wiki/Nagle's_algorithm) is primarily useful for reducing latency but can also improve throughput.  禁用 Nagle 算法主要用于减少延迟，但也可以提高吞吐量。

**kernel.inet_default_connect_options** and **kernel.inet_default_listen_options** must include {nodelay, true} to disable Nagle's algorithm for inter-node connections.

kernel.inet_default_connect_options 和 kernel.inet_default_listen_options 必须包含 {nodelay, true} 以禁用 Nagle 的节点间连接算法。

When configuring sockets that serve client connections, **tcp_listen_options** must include the same option. This is the default.  配置服务于客户端连接的套接字时，tcp_listen_options 必须包含相同的选项。 这是默认设置。

The following example demonstrates that. First, rabbitmq.conf:  下面的例子证明了这一点。 首先，rabbitmq.conf：

```bash
tcp_listen_options.backlog = 4096
tcp_listen_options.nodelay = true
```

which should be used together with the following bits in the [advanced config file](https://www.rabbitmq.com/configure.html#advanced-config-file):

```erlang
[
  {kernel, [
    {inet_default_connect_options, [{nodelay, true}]},
    {inet_default_listen_options,  [{nodelay, true}]}
  ]}].
```

When using the [classic config format](https://www.rabbitmq.com/configure.html#erlang-term-config-file), everything is configured in a single file:

```erlang
[
  {kernel, [
    {inet_default_connect_options, [{nodelay, true}]},
    {inet_default_listen_options,  [{nodelay, true}]}
  ]},
  {rabbit, [
    {tcp_listen_options, [
                          {backlog,       4096},
                          {nodelay,       true},
                          {linger,        {true,0}},
                          {exit_on_close, false}
                         ]}
  ]}
].
```

### Erlang VM I/O Thread Pool Tuning      Erlang VM I/O 线程池调优

Adequate Erlang VM I/O thread pool size is also important when tuning for a large number of concurrent connections. See the section above.  在调整大量并发连接时，足够的 Erlang VM I/O 线程池大小也很重要。 请参阅上面的部分。

### Connection Backlog  连接积压

With a low number of clients, new connection rate is very unevenly distributed but is also small enough to not make much difference. When the number reaches tens of thousands or more, it is important to make sure that the server can accept inbound connections. Unaccepted TCP connections are put into a queue with bounded length. This length has to be sufficient to account for peak load hours and possible spikes, for instance, when many clients disconnect due to a network interruption or choose to reconnect. This is configured using the tcp_listen_options.backlog option:  在客户端数量较少的情况下，新连接速率分布非常不均匀，但也足够小，不会产生太大影响。 当数量达到数万或更多时，确保服务器可以接受入站连接很重要。 未接受的 TCP 连接被放入一个有界长度的队列中。 此长度必须足以考虑峰值负载时间和可能的峰值，例如，当许多客户端由于网络中断而断开连接或选择重新连接时。 这是使用 tcp_listen_options.backlog 选项配置的：

```bash
tcp_listen_options.backlog = 4096
tcp_listen_options.nodelay = true
```

In the [classic config format](https://www.rabbitmq.com/configure.html#erlang-term-config-file):

```erlang
[
  {rabbit, [
    {tcp_listen_options, [
                          {backlog,       4096},
                          {nodelay,       true},
                          {linger,        {true, 0}},
                          {exit_on_close, false}
                         ]}
  ]}
].
```

Default value is 128. When pending connection queue length grows beyond this value, connections will be rejected by the operating system. See also net.core.somaxconn in the kernel tuning section.  默认值为 128。当挂起的连接队列长度超过此值时，连接将被操作系统拒绝。另请参阅内核调优部分中的 net.core.somaxconn。

## Dealing with High Connection Churn  处理高连接流失

### Why is High Connection Churn Problematic?  为什么高连接流失有问题？

Workloads with high connection churn (a high rate of connections being opened and closed) will require TCP setting tuning to avoid exhaustion of certain resources: **max number of file handles**, **Erlang processes** on RabbitMQ nodes, **kernel's ephemeral port range** (for hosts that *open* a lot of connections, including [Federation](https://www.rabbitmq.com/federation.html) links and [Shovel](https://www.rabbitmq.com/shovel.html) connections), and others. Nodes that are exhausted of those resources **won't be able to accept new connections**, which will negatively affect overall system availability.  具有高连接流失率（打开和关闭连接的速率很高）的工作负载将需要 TCP 设置调整以避免耗尽某些资源：文件句柄的最大数量、RabbitMQ 节点上的 Erlang 进程、内核的临时端口范围（对于打开很多连接，包括联邦链接和铲子连接）等。用尽这些资源的节点将无法接受新连接，这将对整个系统的可用性产生负面影响。

Due to a combination of certain TCP features and defaults of most modern Linux distributions, closed connections can be detected after a prolonged period of time. This is covered in the [heartbeats guide](https://www.rabbitmq.com/heartbeats.html). This can be one contributing factor to connection build-up. Another is the **TIME_WAIT** TCP connection state. The state primarily exists to make sure that retransmitted segments from closed connections won't "reappear" on a different (newer) connection with the same client host and port. Depending on the OS and TCP stack configuration connections can spend minutes in this state, which on a busy system is guaranteed to lead to a connection build-up.  由于某些 TCP 功能和大多数现代 Linux 发行版的默认设置相结合，可以在很长一段时间后检测到关闭的连接。这在心跳指南中有介绍。这可能是连接建立的一个促成因素。另一个是 TIME_WAIT TCP 连接状态。该状态的存在主要是为了确保从关闭的连接重新传输的段不会“重新出现”在具有相同客户端主机和端口的不同（较新）连接上。根据操作系统和 TCP 堆栈配置，连接可能会在此状态下花费数分钟，这在繁忙的系统上肯定会导致连接建立。

See [Coping with the TCP TIME_WAIT connections on busy servers](http://vincent.bernat.im/en/blog/2014-tcp-time-wait-state-linux.html) for details.  有关详细信息，请参阅处理繁忙服务器上的 TCP TIME_WAIT 连接。

TCP stack configuration can reduce peak number of connection in closing states and avoid resource exhaustion, in turn allowing nodes to accept new connections at all times.  TCP 堆栈配置可以减少关闭状态下的峰值连接数并避免资源耗尽，从而允许节点始终接受新连接。

High connection churn can also mean developer mistakes or incorrect assumptions about how the messaging protocols supported by RabbitMQ are meant to be used. All supported protocols assume long lived connections. Applications that open and almost immediately close connections unnecessarily waste resources (network bandwidth, CPU, RAM) and contribute to the problem described in this section.  高连接流失也可能意味着开发人员错误或对 RabbitMQ 支持的消息传递协议的使用方式有错误的假设。所有支持的协议都假定连接是长寿命的。打开和几乎立即关闭连接的应用程序会不必要地浪费资源（网络带宽、CPU、RAM）并导致本节中描述的问题。

### Inspecting Connections and Gathering Evidence  检查连接和收集证据

If a node fails to accept connections it is important to first gather data (metrics, evidence) to determine the state of the system and the limiting factor (exhausted resource). Tools such as [netstat](https://en.wikipedia.org/wiki/Netstat), [ss](https://linux.die.net/man/8/ss), [lsof](https://en.wikipedia.org/wiki/Lsof) can be used to inspect TCP connections of a node. See [Troubleshooting Networking](https://www.rabbitmq.com/troubleshooting-networking.html) for examples.  如果节点无法接受连接，则首先收集数据（指标、证据）以确定系统状态和限制因素（资源耗尽）非常重要。 netstat、ss、lsof 等工具可用于检查节点的 TCP 连接。有关示例，请参阅网络故障排除。

While [heartbeats](https://www.rabbitmq.com/heartbeats.html) are sufficient for detecting defunct connections, they are not going to be sufficient in high connection churn scenarios. In those cases heartbeats should be combined with TCP keepalives to speed up disconnected client detection.  虽然心跳足以检测失效的连接，但在高连接流失情况下它们将不够用。在这些情况下，心跳应该与 TCP keepalive 结合起来，以加快断开连接的客户端检测。

### Reducing Amount of Time Spent in TIME_WAIT  减少在 TIME_WAIT 中花费的时间

TCP stack tuning can also reduce the amount of time connections spend in the TIME_WAIT state. The **net.ipv4.tcp_fin_timeout** setting specifically can help here:  TCP 堆栈调整还可以减少连接花费在 TIME_WAIT 状态的时间量。 net.ipv4.tcp_fin_timeout 设置特别可以在这里提供帮助：

```bash
net.ipv4.tcp_fin_timeout = 30
```

Note that like other settings prefixed with net.ipv4., this one applies to both IPv4 and IPv6 connections despite the name.  请注意，与其他以 net.ipv4. 为前缀的设置一样，尽管名称如此，但它同时适用于 IPv4 和 IPv6 连接。

If inbound connections (from clients, plugins, CLI tools and so on) do not rely on NAT, **net.ipv4.tcp_tw_reuse** can be set to 1 (enabled) to allow the kernel to reuse sockets in the TIME_WAIT state for outgoing connections. This setting can be applied on client hosts or intermediaries such as proxies and load balancers. Note that if NAT is used the setting is not safe and can lead to hard to track down issues.  如果入站连接（来自客户端、插件、CLI 工具等）不依赖于 NAT，则 net.ipv4.tcp_tw_reuse 可以设置为 1（启用）以允许内核重用处于 TIME_WAIT 状态的套接字用于传出连接。 此设置可应用于客户端主机或代理和负载平衡器等中介。 请注意，如果使用 NAT，则该设置不安全，并且可能导致难以追踪问题。

The settings above generally should be combined with reduced TCP keepalive values, for example:  上述设置通常应与减少的 TCP keepalive 值结合使用，例如：

```bash
net.ipv4.tcp_fin_timeout = 30

net.ipv4.tcp_keepalive_time=30
net.ipv4.tcp_keepalive_intvl=10
net.ipv4.tcp_keepalive_probes=4

net.ipv4.tcp_tw_reuse = 1
```

## OS Level Tuning  操作系统级别调整

Operating system settings can affect operation of RabbitMQ. Some are directly related to networking (e.g. TCP settings), others affect TCP sockets as well as other things (e.g. open file handles limit).  操作系统设置会影响 RabbitMQ 的运行。 有些与网络直接相关（例如 TCP 设置），有些影响 TCP 套接字以及其他事物（例如打开文件句柄限制）。

Understanding these limits is important, as they may change depending on the workload.  了解这些限制很重要，因为它们可能会根据工作负载而变化。

A few important configurable kernel options include (note that despite option names they are effective for both IPv4 and IPv6 connections):  一些重要的可配置内核选项包括（请注意，尽管有选项名称，但它们对 IPv4 和 IPv6 连接都有效）：

| Kernel setting | Description |
| -------------- | ----------- |
| fs.file-max | Max number of files the kernel will allocate. Limits and current value can be inspected using **/proc/sys/fs/file-nr**.  内核将分配的最大文件数。 可以使用 /proc/sys/fs/file-nr 检查限制和当前值。 |
| net.ipv4.ip_local_port_range | Local IP port range, define as a pair of values. The range must provide enough entries for the peak number of concurrent connections.  本地 IP 端口范围，定义为一对值。 该范围必须为并发连接的峰值数量提供足够的条目。 |
| net.ipv4.tcp_tw_reuse | When enabled, allows the kernel to reuse sockets in TIME_WAIT state when it's safe to do so. See Dealing with High Connection Churn. This option is dangerous when clients and peers connect using NAT.  启用后，允许内核在安全的情况下重用处于 TIME_WAIT 状态的套接字。 请参阅处理高连接流失。 当客户端和对等方使用 NAT 连接时，此选项很危险。 |
| net.ipv4.tcp_fin_timeout | Lowering this timeout to a value in the 15-30 second range reduces the amount of time closed connections will stay in the TIME_WAIT state. See Dealing with High Connection Churn.  将此超时值降低到 15-30 秒范围内的值会减少关闭的连接将保持在 TIME_WAIT 状态的时间量。 请参阅处理高连接流失。 |
| net.core.somaxconn | Size of the listen queue (how many connections are in the process of being established at the same time). Default is 128. Increase to 4096 or higher to support inbound connection bursts, e.g. when clients reconnect en masse.  侦听队列的大小（有多少连接正在同时建立过程中）。 默认值为 128。增加到 4096 或更高以支持入站连接突发，例如当客户端重新连接时。 |
| net.ipv4.tcp_max_syn_backlog | Maximum number of remembered connection requests which did not receive an acknowledgment yet from connecting client. Default is 128, max value is 65535. 4096 and 8192 are recommended starting values when optimising for throughput.  尚未从连接客户端收到确认的已记住连接请求的最大数量。 默认值为 128，最大值为 65535。在优化吞吐量时，建议使用 4096 和 8192 的起始值。 |
| net.ipv4.tcp_keepalive_* | **net.ipv4.tcp_keepalive_time**, **net.ipv4.tcp_keepalive_intvl**, and **net.ipv4.tcp_keepalive_probes** configure TCP keepalive. AMQP 0-9-1 and STOMP have [Heartbeats](https://www.rabbitmq.com/heartbeats.html) which partially undo its effect, namely that it can take minutes to detect an unresponsive peer, e.g. in case of a hardware or power failure. MQTT also has its own keepalives mechanism which is the same idea under a different name. When enabling TCP keepalive with default settings, we recommend setting heartbeat timeout to 8-20 seconds. Also see a note on TCP keepalives later in this guide.            net.ipv4.tcp_keepalive_time、net.ipv4.tcp_keepalive_intvl 和 net.ipv4.tcp_keepalive_probes 配置 TCP keepalive。 AMQP 0-9-1 和 STOMP 有 Heartbeats 可以部分撤销它的效果，即它可能需要几分钟来检测一个无响应的对等体，例如。 在硬件或电源故障的情况下。 MQTT 也有自己的 keepalives 机制，这是相同的想法，但名称不同。 使用默认设置启用 TCP keepalive 时，我们建议将心跳超时设置为 8-20 秒。 另请参阅本指南后面关于 TCP keepalives 的说明。 |
| net.ipv4.conf.default.rp_filter | Enabled reverse path filtering. If [IP address spoofing](http://en.wikipedia.org/wiki/IP_address_spoofing) is not a concern for your system, disable it.  启用反向路径过滤。 如果您的系统不关心 IP 地址欺骗，请将其禁用。 |

Note that default values for these vary between Linux kernel releases and distributions. Using a recent kernel (3.9 or later) is recommended.  请注意，这些默认值因 Linux 内核版本和发行版而异。 建议使用最新的内核（3.9 或更高版本）。

Kernel parameter tuning differs from OS to OS. This guide focuses on Linux. To configure a kernel parameter interactively, use **sysctl -w** (requires superuser privileges), for example:  内核参数调整因操作系统而异。 本指南侧重于 Linux。 要以交互方式配置内核参数，请使用 sysctl -w（需要超级用户权限），例如：

```bash
sysctl -w fs.file-max 200000
```

To make the changes permanent (stick between reboots), they need to be added to **/etc/sysctl.conf**. See [sysctl(8)](http://man7.org/linux/man-pages/man8/sysctl.8.html) and [sysctl.conf(5)](http://man7.org/linux/man-pages/man5/sysctl.conf.5.html) for more details.  要使更改永久生效（在重新启动之间保持不变），需要将它们添加到 /etc/sysctl.conf。 有关详细信息，请参阅 sysctl(8) 和 sysctl.conf(5)。

TCP stack tuning is a broad topic that is covered in much detail elsewhere:  TCP 堆栈调整是一个广泛的主题，在其他地方有更详细的介绍：

- [Enabling High Performance Data Transfers](https://psc.edu/index.php/services/networking/68-research/networking/641-tcp-tune)  启用高性能数据传输

- [Network Tuning Guide](https://fasterdata.es.net/network-tuning/)  网络调优指南

## TCP Socket Options  TCP 套接字选项

### Common Options  常用选项

| Kernel setting | Description |
| -------------- | ----------- |
| tcp_listen_options.nodelay | When set to true, disables [Nagle's algorithm](http://en.wikipedia.org/wiki/Nagle's_algorithm). Default is true. Highly recommended for most users.  当设置为 true 时，禁用 Nagle 算法。 默认为真。 强烈推荐给大多数用户。 |
| tcp_listen_options.sndbuf | See TCP buffers discussion earlier in this guide. Default value is automatically tuned by the OS, typically in the 88 KiB to 128 KiB range on modern Linux versions. Increasing buffer size improves consumer throughput and RAM use for every connection. Decreasing has the opposite effect.  请参阅本指南前面的 TCP 缓冲区讨论。 默认值由操作系统自动调整，在现代 Linux 版本上通常在 88 KiB 到 128 KiB 范围内。 增加缓冲区大小可以提高每个连接的消费者吞吐量和 RAM 使用率。 减少会产生相反的效果。 |
| tcp_listen_options.recbuf | See TCP buffers discussion earlier in this guide. Default value effects are similar to that of tcp_listen_options.sndbuf but for publishers and protocol operations in general.  请参阅本指南前面的 TCP 缓冲区讨论。 默认值效果与 tcp_listen_options.sndbuf 的效果类似，但一般用于发布者和协议操作。 |
| tcp_listen_options.backlog | Maximum size of the unaccepted TCP connections queue. When this size is reached, new connections will be rejected. Set to 4096 or higher for environments with thousands of concurrent connections and possible bulk client reconnections.  未接受的 TCP 连接队列的最大大小。 当达到这个大小时，新的连接将被拒绝。 对于具有数千个并发连接和可能的批量客户端重新连接的环境，设置为 4096 或更高。 |
| tcp_listen_options.keepalive | When set to true, enables TCP keepalives (see above). Default is false. Makes sense for environments where connections can go idle for a long time (at least 10 minutes), although using [heartbeats](https://www.rabbitmq.com/heartbeats.html) is still recommended over this option.  当设置为 true 时，启用 TCP keepalives（见上文）。 默认为假。 对于连接可能长时间处于空闲状态（至少 10 分钟）的环境是有意义的，尽管仍然建议在此选项上使用心跳。 |

### Defaults

Below is the default TCP socket option configuration used by RabbitMQ:  以下是 RabbitMQ 使用的默认 TCP 套接字选项配置：

- TCP connection backlog is limited to 128 connections  TCP 连接积压限制为 128 个连接

- Nagle's algorithm is disabled  Nagle 算法被禁用

- Server socket lingering is enabled with the timeout of 0  服务器套接字延迟启用，超时时间为 0

## Heartbeats

Some protocols supported by RabbitMQ, including AMQP 0-9-1, support *heartbeats*, a way to detect dead TCP peers quicker. Please refer to the [Heartbeats guide](https://www.rabbitmq.com/heartbeats.html) for more information.  RabbitMQ 支持的一些协议，包括 AMQP 0-9-1，支持心跳，这是一种更快检测死 TCP 对等点的方法。有关详细信息，请参阅 Heartbeats 指南。

## Net Tick Time

[Heartbeats](https://www.rabbitmq.com/heartbeats.html) are used to detect peer or connection failure between clients and RabbitMQ nodes. [net_ticktime](https://www.rabbitmq.com/nettick.html) serves the same purpose but for cluster node communication. Values lower than 5 (seconds) may result in false positive and are not recommended.  心跳用于检测客户端和 RabbitMQ 节点之间的对等或连接故障。 net_ticktime 具有相同的目的，但用于集群节点通信。小于 5（秒）的值可能会导致误报，不推荐使用。

## TCP Keepalives

TCP contains a mechanism similar in purpose to the heartbeat (a.k.a. keepalive) one in messaging protocols and net tick timeout covered above: TCP keepalives. Due to inadequate defaults, TCP keepalives often don't work the way they are supposed to: it takes a very long time (say, an hour or more) to detect a dead peer. However, with tuning they can serve the same purpose as heartbeats and clean up stale TCP connections e.g. with clients that opted to not use heartbeats, intentionally or not.  TCP 包含一种机制，其目的类似于消息传递协议中的心跳（a.k.a.keepalive）机制和上面介绍的网络滴答超时：TCP keepalives。由于默认值不足，TCP keepalives 通常不能按预期方式工作：检测死对等点需要很长时间（例如，一个小时或更长时间）。但是，通过调优，它们可以起到与心跳相同的目的，并清理过时的 TCP 连接，例如与有意或无意选择不使用心跳的客户一起使用。

Below is an example sysctl configuration for TCP keepalives that considers TCP connections dead or unreachable after 70 seconds (4 attempts every 10 seconds after connection idle for 30 seconds):  下面是 TCP keepalives 的示例 sysctl 配置，它认为 TCP 连接在 70 秒后死亡或无法访问（连接空闲 30 秒后每 10 秒尝试 4 次）：

```bash
net.ipv4.tcp_keepalive_time=30
net.ipv4.tcp_keepalive_intvl=10
net.ipv4.tcp_keepalive_probes=4
```

TCP keepalives can be a useful additional defense mechanism in environments where RabbitMQ operator has no control over application settings or client libraries used.  在 RabbitMQ 操作员无法控制所使用的应用程序设置或客户端库的环境中，TCP keepalives 可以成为一种有用的附加防御机制。

## Connection Handshake Timeout  连接握手超时

RabbitMQ has a timeout for connection handshake, 10 seconds by default. When clients run in heavily constrained environments, it may be necessary to increase the timeout. This can be done via the rabbit.handshake_timeout (in milliseconds):  RabbitMQ 的连接握手超时，默认为 10 秒。当客户端在严重受限的环境中运行时，可能需要增加超时时间。这可以通过 rabbit.handshake_timeout（以毫秒为单位）来完成：

```bash
handshake_timeout = 20000
```

It should be pointed out that this is only necessary with very constrained clients and networks. Handshake timeouts in other circumstances indicate a problem elsewhere.  应该指出的是，这只有在客户端和网络非常受限的情况下才需要。其他情况下的握手超时表明其他地方存在问题。

### TLS (SSL) Handshake

If TLS/SSL is enabled, it may be necessary to increase also the TLS/SSL handshake timeout. This can be done via the rabbit.ssl_handshake_timeout (in milliseconds):  如果启用了 TLS/SSL，则可能还需要增加 TLS/SSL 握手超时。这可以通过 rabbit.ssl_handshake_timeout（以毫秒为单位）来完成：

```bash
ssl_handshake_timeout = 10000
```

## Hostname Resolution and DNS  主机名解析和 DNS

In many cases, RabbitMQ relies on the Erlang runtime for inter-node communication (including tools such as rabbitmqctl, rabbitmq-plugins, etc). Client libraries also perform hostname resolution when connecting to RabbitMQ nodes. This section briefly covers most common issues associated with that.  在很多情况下，RabbitMQ 依赖 Erlang 运行时进行节点间通信（包括 rabbitmqctl、rabbitmq-plugins 等工具）。客户端库在连接到 RabbitMQ 节点时也会执行主机名解析。本节简要介绍了与此相关的最常见问题。

### Performed by Client Libraries  由客户端库执行

If a client library is configured to connect to a hostname, it performs hostname resolution. Depending on DNS and local resolver (/etc/hosts and similar) configuration, this can take some time. Incorrect configuration may lead to resolution timeouts, e.g. when trying to resolve a local hostname such as my-dev-machine, over DNS. As a result, client connections can take a long time (from tens of seconds to a few minutes).  如果客户端库配置为连接到主机名，它会执行主机名解析。根据 DNS 和本地解析器（/etc/hosts 和类似）配置，这可能需要一些时间。不正确的配置可能会导致解析超时，例如尝试通过 DNS 解析本地主机名（例如 my-dev-machine）时。因此，客户端连接可能需要很长时间（从几十秒到几分钟）。

### Short and Fully-qualified RabbitMQ Node Names  简短且完全限定的 RabbitMQ 节点名称

RabbitMQ relies on the Erlang runtime for inter-node communication. Erlang nodes include a hostname, either short (rmq1) or fully-qualified (rmq1.dev.megacorp.local). Mixing short and fully-qualified hostnames is not allowed by the runtime. Every node in a cluster must be able to resolve every other node's hostname, short or fully-qualified.  RabbitMQ 依赖 Erlang 运行时进行节点间通信。 Erlang 节点包括一个主机名，可以是短的 (rmq1) 或完全限定的 (rmq1.dev.megacorp.local)。运行时不允许混合短主机名和全限定主机名。集群中的每个节点都必须能够解析每个其他节点的主机名，无论是短主机名还是全限定主机名。

By default RabbitMQ will use short hostnames. Set the RABBITMQ_USE_LONGNAME environment variable to make RabbitMQ nodes use fully-qualified names, e.g. rmq1.dev.megacorp.local.  默认情况下，RabbitMQ 将使用短主机名。设置 RABBITMQ_USE_LONGNAME 环境变量以使 RabbitMQ 节点使用完全限定的名称，例如rmq1.dev.megacorp.local。

### Reverse DNS Lookups  反向 DNS 查找

If the reverse_dns_lookups configuration option is set to true, RabbitMQ will perform reverse DNS lookups for client IP addresses and list hostnames in connection information (e.g. in the [Management UI](https://www.rabbitmq.com/management.html)).  如果 reverse_dns_lookups 配置选项设置为 true，RabbitMQ 将对客户端 IP 地址执行反向 DNS 查找，并在连接信息中列出主机名（例如在管理 UI 中）。

Reverse DNS lookups can potentially take a long time if node's hostname resolution is not optimally configured. This can increase latency when accepting client connections.  如果节点的主机名解析没有最佳配置，反向 DNS 查找可能需要很长时间。这会在接受客户端连接时增加延迟。

To explicitly enable reverse DNS lookups:  要显式启用反向 DNS 查找：

```bash
reverse_dns_lookups = true
```

To disable reverse DNS lookups:

```bash
reverse_dns_lookups = false
```

### Verify Hostname Resolution on a Node or Locally  在节点上或本地验证主机名解析

Since hostname resolution is a [prerequisite for successful inter-node communication](https://www.rabbitmq.com/clustering.html#hostname-resolution-requirement), starting with [RabbitMQ 3.8.6](https://www.rabbitmq.com/changelog.html), CLI tools provide two commands that help verify that hostname resolution on a node works as expected. The commands are not meant to replace [dig](https://en.wikipedia.org/wiki/Dig_(command)) and other specialised DNS tools but rather provide a way to perform most basic checks while taking [Erlang runtime hostname resolver features](https://erlang.org/doc/apps/erts/inet_cfg.html) into account.  由于主机名解析是成功的节点间通信的先决条件，从 RabbitMQ 3.8.6 开始，CLI 工具提供了两个命令来帮助验证节点上的主机名解析是否按预期工作。这些命令并不是要取代 dig 和其他专门的 DNS 工具，而是提供一种在考虑 Erlang 运行时主机名解析器功能的同时执行最基本检查的方法。

The first command is rabbitmq-diagnostics resolve_hostname:  第一个命令是 rabbitmq-diagnostics resolve_hostname：

```bash
# resolves node2.cluster.local.svc to IPv6 addresses on node rabbit@node1.cluster.local.svc
rabbitmq-diagnostics resolve_hostname node2.cluster.local.svc --address-family IPv6 -n rabbit@node1.cluster.local.svc

# makes local CLI tool resolve node2.cluster.local.svc to IPv4 addresses
rabbitmq-diagnostics resolve_hostname node2.cluster.local.svc --address-family IPv4 --offline
```

The second one is rabbitmq-diagnostics resolver_info:

```bash
rabbitmq-diagnostics resolver_info
```

It will report key resolver settings such as the lookup order (whether CLI tools should prefer the OS resolver, inetrc file, and so on) as well as inetrc hostname entries, if any:  它将报告关键解析器设置，例如查找顺序（CLI 工具是否应该首选操作系统解析器、inetrc 文件等）以及 inetrc 主机名条目（如果有）：

```bash
Runtime Hostname Resolver (inetrc) Settings

Lookup order: native
Hosts file: /etc/hosts
Resolver conf file: /etc/resolv.conf
Cache size:

inetrc File Host Entries

(none)
```

## Connection Event Logging  连接事件记录

See [Connection Lifecycle Events](https://www.rabbitmq.com/logging.html#connection-lifecycle-events) in the logging guide.  请参阅日志记录指南中的连接生命周期事件。

## Troubleshooting Network Connectivity  网络连接故障排除

A methodology for [troubleshooting of networking-related issues](https://www.rabbitmq.com/troubleshooting-networking.html) is covered in a separate guide.  单独的指南中介绍了网络相关问题的故障排除方法。

## MacOS Application Firewall  MacOS 应用程序防火墙

On MacOS systems with [Application Firewall](https://support.apple.com/en-us/HT201642) enabled, Erlang runtime processes must be allowed to bind to ports and accept connections. Without this, RabbitMQ nodes won't be able to bind to their [ports](https://www.rabbitmq.com/networking.html#ports) and will fail to start.  在启用了应用程序防火墙的 MacOS 系统上，必须允许 Erlang 运行时进程绑定到端口并接受连接。 没有这个，RabbitMQ 节点将无法绑定到它们的端口并且将无法启动。

A list of blocked applications can be seen under Security and Privacy => Firewall in system settings.  在系统设置中的安全和隐私 => 防火墙下可以看到被阻止的应用程序列表。

To "unblock" a command line tool, use sudo /usr/libexec/ApplicationFirewall/socketfilterfw. The examples below assume that Erlang is installed under /usr/local/Cellar/erlang/{version}, used by the Homebrew Erlang formula:  要“取消阻止”命令行工具，请使用 sudo /usr/libexec/ApplicationFirewall/socketfilterfw。 下面的示例假设 Erlang 安装在 Homebrew Erlang 公式使用的 /usr/local/Cellar/erlang/{version} 下：

```bash
# allow CLI tools and shell to bind to ports and accept inbound connections
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/local/Cellar/erlang/{version}/lib/erlang/bin/erl
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /usr/local/Cellar/erlang/{version}/lib/erlang/bin/erl
# allow server nodes (Erlang VM) to bind to ports and accept inbound connections
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --add /usr/local/Cellar/erlang/{version}/lib/erlang/erts-{erts version}/bin/beam.smp
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --unblockapp /usr/local/Cellar/erlang/{version}/lib/erlang/erts-{erts version}/bin/beam.smp
```

Note that socketfilterfw command line arguments can vary between MacOS releases. To see supports command line arguments, use  请注意，socketfilterfw 命令行参数可能因 MacOS 版本而异。 要查看支持命令行参数，请使用

```bash
sudo /usr/libexec/ApplicationFirewall/socketfilterfw --help
```

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!

