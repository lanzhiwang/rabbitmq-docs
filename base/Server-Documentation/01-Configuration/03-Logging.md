# Logging

https://www.rabbitmq.com/logging.html

## Overview

Log files is a very important aspect of system observability, much like [monitoring](https://www.rabbitmq.com/monitoring.html).  日志文件是系统可观察性的一个非常重要的方面，就像监控一样。

Developers and operators should inspect logs when troubleshooting an issue or assessing the state of the system.  在对问题进行故障排除或评估系统状态时，开发人员和操作员应检查日志。

RabbitMQ supports a number of features when it comes to logging.  RabbitMQ 在日志记录方面支持许多功能。

This guide covers topics such as:  本指南涵盖以下主题：

- Supported log outputs: file and standard streams (console)  支持的日志输出：文件和标准流（控制台）

- Log file location  日志文件位置

- Supported log levels  支持的日志级别

- How to enable debug logging  如何启用调试日志记录

- How to tail logs of a running node without having access to the log file  如何在不访问日志文件的情况下跟踪正在运行的节点的日志

- Watching internal events  观看内部事件

- Connection lifecycle events logged  记录的连接生命周期事件

- Log categories  日志类别

- How to inspect service logs on systemd-based Linux systems  如何检查基于 systemd 的 Linux 系统上的服务日志

- Log rotation  日志轮换

- Logging to Syslog  记录到系统日志

- Logging to a system topic exchange, **amq.rabbitmq.log**  记录到系统主题交换 amq.rabbitmq.log

and more.

## Log File Location

Modern RabbitMQ versions use a single log file by default.  现代 RabbitMQ 版本默认使用单个日志文件。

Please see the [File and Directory Location](https://www.rabbitmq.com/relocate.html) guide to find default log file location for various platforms.  请参阅文件和目录位置指南以查找各种平台的默认日志文件位置。

There are two ways to configure log file location. One is the [configuration file](https://www.rabbitmq.com/configure.html). The other is the **RABBITMQ_LOGS** environment variable.  有两种方法可以配置日志文件位置。 一是配置文件。 另一个是 RABBITMQ_LOGS 环境变量。

Use [RabbitMQ management UI](https://www.rabbitmq.com/management.html) or [rabbitmq-diagnostics status](https://www.rabbitmq.com/cli.html) to find when a node stores its log file(s).  使用 RabbitMQ 管理 UI 或 rabbitmq-diagnostics status 来查找节点何时存储其日志文件。

The **RABBITMQ_LOGS** variable value can be either a file path or a hyphen (-). `RABBITMQ_LOGS=-` will result in all log messages being sent to standard output. See Logging to Console (Standard Output).  RABBITMQ_LOGS 变量值可以是文件路径或连字符 (-)。 RABBITMQ_LOGS=- 将导致所有日志消息发送到标准输出。 请参阅记录到控制台（标准输出）。

The environment variable takes precedence over the configuration file. When in doubt, consider overriding log file location via the config file. As a consequence of the environment variable precedence, if the environment variable is set, the configuration key log.file will not have any effect.  环境变量优先于配置文件。 如有疑问，请考虑通过配置文件覆盖日志文件位置。 由于环境变量优先，如果设置了环境变量，配置键 log.file 将不会有任何影响。

## Configuration

RabbitMQ starts logging early on node start. See the [Configuration guide](https://www.rabbitmq.com/configure.html) for a general overview of how to configure RabbitMQ.  RabbitMQ 在节点启动时尽早开始记录。 有关如何配置 RabbitMQ 的一般概述，请参阅配置指南。

### Log Outputs

Default RabbitMQ logging configuration will direct log messages to a log file. Standard output is another option available out of the box.  默认 RabbitMQ 日志配置会将日志消息定向到日志文件。 标准输出是另一种开箱即用的选项。

Several outputs can be used at the same time. Log entries will be copied to all of them.  可以同时使用多个输出。 日志条目将被复制到所有这些条目。

Different outputs can have different log levels. For example, the console output can log all messages including debug information while the file output can only log error and higher severity messages.  不同的输出可以有不同的日志级别。 例如，控制台输出可以记录所有消息，包括调试信息，而文件输出只能记录错误和更高严重性的消息。

### Logging to a File

- **log.file**: log file path or false to disable the file output. Default value is taken from the **RABBITMQ_LOGS** [environment variable or configuration file](https://www.rabbitmq.com/configure.html).  日志文件路径或 false 禁用文件输出。 默认值取自 RABBITMQ_LOGS 环境变量或配置文件。

- **log.file.level**: log level for the file output. Default level is info.  文件输出的日志级别。 默认级别是信息。

- **log.file.rotation.date**, **log.file.rotation.size**, **log.file.rotation.count** for log file rotation settings.

The following example overrides log file name:  以下示例覆盖日志文件名：

```ini
log.file = rabbit.log
```

The following example overrides log file directory:

```ini
log.dir = /data/logs/rabbitmq
```

The following example instructs RabbitMQ to log to a file at the debug level:

```ini
log.file.level = debug
```

Logging to a file can be disabled with

```ini
log.file = false
```

Find supported log levels in the [example rabbitmq.conf file](https://github.com/rabbitmq/rabbitmq-server/blob/v3.8.x/deps/rabbit/docs/rabbitmq.conf.example).  在示例 rabbitmq.conf 文件中查找支持的日志级别。

The rest of this guide describes more options, including more advanced ones.  本指南的其余部分描述了更多选项，包括更高级的选项。

### Log Rotation

RabbitMQ nodes always append to the log files, so a complete log history is preserved. Log file rotation is not performed by default. [Debian](https://www.rabbitmq.com/install-debian.html) and [RPM](https://www.rabbitmq.com/install-rpm.html) packages will set up log rotation via logrotate after package installation.  RabbitMQ 节点始终附加到日志文件，因此保留了完整的日志历史记录。 默认情况下不执行日志文件轮换。 Debian 和 RPM 软件包将在软件包安装后通过 logrotate 设置日志轮换。

**log.file.rotation.date**, **log.file.rotation.size**, **log.file.rotation.count** settings control log file rotation for the file output.  设置控制文件输出的日志文件轮换。

#### Built-in Periodic Rotation  内置周期性旋转

Use **log.file.rotation.date** to set up minimalistic periodic rotation:  使用 log.file.rotation.date 设置简约的定期轮换：

```ini
# rotate every night at midnight
log.file.rotation.date = $D0

# keep up to 5 archived log files in addition to the current one
log.file.rotation.count = 5
```

```ini
# rotate every day at 23:00 (11:00 p.m.)
log.file.rotation.date = $D23
```

```ini
# rotate every night at midnight
log.file.rotation.date = $D0
```

#### Built-in File Size-based Rotation

**log.file.rotation.size** controls rotation based on the current log file size:

```ini
# rotate when the file reaches 10 MiB
log.file.rotation.size = 10485760

# keep up to 5 archived log files in addition to the current one
log.file.rotation.count = 5
```

#### Rotation Using Logrotate  使用 Logrotate 进行旋转

On Linux, BSD and other UNIX-like systems, [logrotate](https://linux.die.net/man/8/logrotate) is an alternative way of log file rotation and compression.  在 Linux、BSD 和其他类 UNIX 系统上，logrotate 是日志文件轮换和压缩的另一种方式。

RabbitMQ [Debian](https://www.rabbitmq.com/install-debian.html) and [RPM](https://www.rabbitmq.com/install-rpm.html) packages will set up logrotate to run weekly on files located in default **/var/log/rabbitmq** directory. Rotation configuration can be found in **/etc/logrotate.d/rabbitmq-server**.  RabbitMQ Debian 和 RPM 软件包将设置 logrotate 每周在位于默认 /var/log/rabbitmq 目录中的文件上运行。 轮换配置可以在 /etc/logrotate.d/rabbitmq-server 中找到。

### Logging to Console (Standard Output)  记录到控制台（标准输出）

Here are the main settings that control console (standard output) logging:  以下是控制台（标准输出）日志记录的主要设置：

- **log.console** (boolean): set to true to enable console output. Default is false

- **log.console.level**: log level for the console output. Default level is info.

To enable console logging, use the following config snippet:  要启用控制台日志记录，请使用以下配置片段：

```ini
log.console = true
```

The following example disables console logging

```ini
log.console = false
```

The following example instructs RabbitMQ to use the debug logging level when logging to console:

```ini
log.console.level = debug
```

When console output is enabled, the file output will also be enabled by default. To disable the file output, set **log.file** to false.  启用控制台输出时，默认情况下也会启用文件输出。 要禁用文件输出，请将 log.file 设置为 false。

Please note that **RABBITMQ_LOGS=-** will disable the file output even if log.file is configured.  请注意，即使配置了 log.file，RABBITMQ_LOGS=- 也会禁用文件输出。

### Logging to Syslog  记录到系统日志

RabbitMQ logs can be forwarded to a Syslog server via TCP or UDP. UDP is used by default and **requires Syslog service configuration**. TLS is also supported.  RabbitMQ 日志可以通过 TCP 或 UDP 转发到 Syslog 服务器。 默认使用 UDP，需要 Syslog 服务配置。 还支持 TLS。

Syslog output has to be explicitly configured:  必须显式配置 Syslog 输出：

```ini
log.syslog = true
```

#### Syslog Endpoint Configuration  系统日志端点配置

By default the Syslog logger will send log messages to UDP port 514 using the [RFC 3164](https://www.ietf.org/rfc/rfc3164.txt) protocol. [RFC 5424](https://tools.ietf.org/html/rfc5424) protocol also can be used.  默认情况下，Syslog 记录器将使用 RFC 3164 协议将日志消息发送到 UDP 端口 514。 也可以使用 RFC 5424 协议。

In order to use UDP the **Syslog service must have UDP input configured**.  为了使用 UDP，Syslog 服务必须配置 UDP 输入。

UDP and TCP transports can be used with both RFC 3164 and RFC 5424 protocols. TLS support requires the RFC 5424 protocol.  UDP 和 TCP 传输可以与 RFC 3164 和 RFC 5424 协议一起使用。 TLS 支持需要 RFC 5424 协议。

The following example uses TCP and the RFC 5424 protocol:  以下示例使用 TCP 和 RFC 5424 协议：

```ini
log.syslog = true
log.syslog.transport = tcp
log.syslog.protocol = rfc5424
```

To TLS, a standard set of [TLS options](https://www.rabbitmq.com/ssl.html) must be provided:

```ini
log.syslog = true
log.syslog.transport = tls
log.syslog.protocol = rfc5424

log.syslog.ssl_options.cacertfile = /path/to/ca_certificate.pem
log.syslog.ssl_options.certfile = /path/to/client_certificate.pem
log.syslog.ssl_options.keyfile = /path/to/client_key.pem
```

Syslog service IP address and port can be customised:

```ini
log.syslog = true
log.syslog.ip = 10.10.10.10
log.syslog.port = 1514
```

If a hostname is to be used rather than an IP address:

```ini
log.syslog = true
log.syslog.host = my.syslog-server.local
log.syslog.port = 1514
```

Syslog metadata identity and facility values also can be configured. By default identity will be set to the name part of the node name (for example, rabbitmq in rabbitmq@hostname) and facility will be set to daemon.

To set identity and facility of log messages:

```ini
log.syslog = true
log.syslog.identity = my_rabbitmq
log.syslog.facility = user
```

Less commonly used [Syslog client](https://github.com/schlagert/syslog) options can be configured using the [advanced config file](https://www.rabbitmq.com/configure.html#configuration-files).

## Log Message Categories

RabbitMQ has several categories of messages, which can be logged with different levels or to different files.  RabbitMQ 有几类消息，可以用不同的级别或不同的文件记录。

The categories replace the rabbit.log_levels configuration setting in versions earlier than 3.7.0.  这些类别替换了 3.7.0 之前版本中的 rabbit.log_levels 配置设置。

The categories are:  类别是：


- **connection**: connection lifecycle events for AMQP 0-9-1, AMQP 1.0, MQTT and STOMP.

- **channel**: channel logs. Mostly errors and warnings on AMQP 0-9-1 channels.

- **queue**: queue logs. Mostly debug messages.

- **mirroring**: queue mirroring logs. Queue mirrors status changes: starting/stopping/synchronizing.

- **federation**: federation plugin logs.

- **upgrade**: verbose upgrade logs. These can be excessive.  详细的升级日志。 这些可能是多余的。

- **default**: all other log entries. You cannot override file location for this category.  所有其他日志条目。 您不能覆盖此类别的文件位置。

It is possible to configure a different log level or file location for each message category using **log.\<category\>.level** and **log.\<category\>.file** configuration variables.

By default each category will not filter by level. If an is output configured to log debug messages, the debug messages will be printed for all categories. Configure a log level for a category to override.  默认情况下，每个类别不会按级别过滤。 如果输出配置为记录调试消息，则将为所有类别打印调试消息。 为要覆盖的类别配置日志级别。

For example, given debug level in the file output, the following will disable debug logging for connection events:  例如，给定文件输出中的调试级别，以下将禁用连接事件的调试日志记录：

```ini
log.file.level = debug
log.connection.level = info
```

To redirect all federation logs to the rabbit_federation.log file, use:

```ini
log.federation.file = rabbit_federation.log
```

To disable a log type, you can use the none log level. For example, to disable upgrade logs:

```ini
log.upgrade.level = none
```

### Log Levels

Log levels is another way to filter and tune logging. Each log level has a severity associated with it. More critical messages have lower severity number, while debug has the highest number.  日志级别是过滤和调整日志记录的另一种方式。 每个日志级别都有一个与之关联的严重性。 更关键的消息具有较低的严重性编号，而调试具有最高的编号。

The following log levels are used by RabbitMQ:  RabbitMQ 使用以下日志级别：

| Log level | Severity |
| --------- | -------- |
| debug     | 128      |
| info      | 64       |
| warning   | 16       |
| error     | 8        |
| critical  | 4        |
| none      | 0        |

Default log level is info.

If the level of a log message is higher than the category level, the message will be dropped and not sent to any output.

If a category level is not configured, its messages will always be sent to all outputs.

To make the default category log only errors or higher severity messages, use

```ini
log.default.level = error
```

The none level means no logging.

Each output can use its own log level. If a message level number is higher than the output level, the message will not be logged.  每个输出都可以使用自己的日志级别。 如果消息级别编号高于输出级别，则不会记录该消息。

For example, if no outputs are configured to log debug messages, even if the category level is set to debug, the debug messages will not be logged.  例如，如果没有配置输出来记录调试消息，即使类别级别设置为调试，也不会记录调试消息。

Although, if an output is configured to log debug messages, it will get them from all categories, unless a category level is configured.  虽然，如果将输出配置为记录调试消息，它将从所有类别中获取它们，除非配置了类别级别。

#### Changing Log Level

There are two ways of changing effective log levels:

- Via [configuration file(s)](https://www.rabbitmq.com/configure.html): this is more flexible but requires a node restart between changes

- Using [CLI tools](https://www.rabbitmq.com/cli.html), **rabbitmqctl set_log_level \<leve\l>**: the changes are transient (will not survive node restart) but can be used to enable and disable e.g. debug logging at runtime for a period of time.

To set log level to debug on a running node:

```bash
rabbitmqctl -n rabbit@target-host set_log_level debug
```

To set the level to info:

```bash
rabbitmqctl -n rabbit@target-host set_log_level info
```

## Tailing Logs Using CLI Tools

Modern releases support tailing logs of a node using [CLI tools](https://www.rabbitmq.com/cli.html). This is convenient when log file location is not known or is not easily accessible but CLI tool connectivity is allowed.  现代版本支持使用 CLI 工具跟踪节点的日志。 当日志文件位置未知或不易访问但允许 CLI 工具连接时，这很方便。

To tail three hundred last lines on a node **rabbitmq@target-host**, use **rabbitmq-diagnostics log_tail**:  要在节点 rabbitmq@target-host 上拖尾最后三百行，请使用 rabbitmq-diagnostics log_tail：

```bash
# This is semantically equivalent to using `tail -n 300 /path/to/rabbit@hostname.log`.
# Use -n to specify target node, -N is to specify the number of lines.
rabbitmq-diagnostics -n rabbit@target-host log_tail -N 300
```

This will load and print last lines from the log file. If only console logging is enabled, this command will fail with a "file not found" (enoent) error.  这将从日志文件中加载并打印最后几行。 如果仅启用控制台日志记录，则此命令将失败并出现“找不到文件”（enoent）错误。

To continuously inspect as a stream of log messages as they are appended to a file, similarly to **tail -f** or console logging, use **rabbitmq-diagnostics log_tail_stream**:  要在附加到文件时作为日志消息流持续检查，类似于 tail -f 或控制台日志记录，请使用 rabbitmq-diagnostics log_tail_stream：

```bash
# This is semantically equivalent to using `tail -f /path/to/rabbit@hostname.log`.
# Use Control-C to stop the stream.
rabbitmq-diagnostics -n rabbit@target-host log_tail_stream
```

This will continuously tail and stream lines added to the log file. If only console logging is enabled, this command will fail with a "file not found" (enoent) error.  这将不断地尾随和流式添加到日志文件中。 如果仅启用控制台日志记录，则此命令将失败并出现“找不到文件”（enoent）错误。

The **rabbitmq-diagnostics log_tail_stream** command can only be used against a running RabbitMQ node and will fail if the node is not running or the RabbitMQ application on it was stopped using **rabbitmqctl stop_app**.  rabbitmq-diagnostics log_tail_stream 命令只能用于正在运行的 RabbitMQ 节点，如果该节点未运行或者其上的 RabbitMQ 应用程序已使用 rabbitmqctl stop_app 停止，则该命令将失败。

## Enabling Debug Logging

To enable debug messages, you should have a debug output.

For example to log debug messages to a file:

```ini
log.file.level = debug
```

To print log messages to standard I/O streams:

```ini
log.console = true
log.console.level = debug
```

To switch to debug logging at runtime:

```bash
rabbitmqctl -n rabbit@target-host set_log_level debug
```

To set the level back to info:

```bash
rabbitmqctl -n rabbit@target-host set_log_level info
```

It is possible to disable debug logging for some categories:

```ini
log.file.level = debug

log.connection.level = info
log.channel.level = info
```

## Service Logs

On systemd-based Linux distributions, system service logs can be inspected using **journalctl --system**  在基于 systemd 的 Linux 发行版上，可以使用 journalctl --system 检查系统服务日志

```bash
journalctl --system
```

which requires superuser privileges. Its output can be filtered to narrow it down to RabbitMQ-specific entries:

```bash
sudo journalctl --system | grep rabbitmq
```

Service logs will include standard output and standard error streams of the node. The output of journalctl --system will look similar to this:

```plaintext
Dec 26 11:03:04 localhost rabbitmq-server[968]: ##  ##
Dec 26 11:03:04 localhost rabbitmq-server[968]: ##  ##      RabbitMQ 3.8.5. Copyright (c) 2007-2022 VMware, Inc. or its affiliates.
Dec 26 11:03:04 localhost rabbitmq-server[968]: ##########  Licensed under the MPL.  See https://www.rabbitmq.com/
Dec 26 11:03:04 localhost rabbitmq-server[968]: ######  ##
Dec 26 11:03:04 localhost rabbitmq-server[968]: ##########  Logs: /var/log/rabbitmq/rabbit@localhost.log
Dec 26 11:03:04 localhost rabbitmq-server[968]: /var/log/rabbitmq/rabbit@localhost_upgrade.log
Dec 26 11:03:04 localhost rabbitmq-server[968]: Starting broker...
Dec 26 11:03:05 localhost rabbitmq-server[968]: systemd unit for activation check: "rabbitmq-server.service"
Dec 26 11:03:06 localhost rabbitmq-server[968]: completed with 6 plugins.
```

## Logged Events

### Connection Lifecycle Events

Successful TCP connections that send at least 1 byte of data will be logged. Connections that do not send any data, such as health checks of certain load balancer products, will not be logged.  将记录发送至少 1 个字节数据的成功 TCP 连接。 不会记录不发送任何数据的连接，例如某些负载均衡器产品的健康检查。

Here's an example:

```plaintext
2018-11-22 10:44:33.654 [info] <0.620.0> accepting AMQP connection <0.620.0> (127.0.0.1:52771 -> 127.0.0.1:5672)
```

The entry includes client IP address and port (127.0.0.1:52771) as well as the target IP address and port of the server (127.0.0.1:5672). This information can be useful when troubleshooting client connections.

Once a connection successfully authenticates and is granted access to a [virtual host](https://www.rabbitmq.com/vhosts.html), that is also logged:

```plaintext
2018-11-22 10:44:33.663 [info] <0.620.0> connection <0.620.0> (127.0.0.1:52771 -> 127.0.0.1:5672): user 'guest' authenticated and granted access to vhost '/'
```

The examples above include two values that can be used as connection identifiers in various scenarios: connection name (127.0.0.1:57919 -> 127.0.0.1:5672) and an Erlang process ID of the connection (<0.620.0>). The latter is used by [rabbitmqctl](https://www.rabbitmq.com/cli.html) and the former is used by the [HTTP API](https://www.rabbitmq.com/management.html).

A [client connection](https://www.rabbitmq.com/connections.html) can be closed cleanly or abnormally. In the former case the client closes AMQP 0-9-1 (or 1.0, or STOMP, or MQTT) connection gracefully using a dedicated library function (method). In the latter case the client closes TCP connection or TCP connection fails. RabbitMQ will log both cases.

Below is an example entry for a successfully closed connection:

```plaintext
2018-06-17 06:23:29.855 [info] <0.634.0> closing AMQP connection <0.634.0> (127.0.0.1:58588 -> 127.0.0.1:5672, vhost: '/', user: 'guest')
```

Abruptly closed connections will be logged as warnings:

```plaintext
2018-06-17 06:28:40.868 [warning] <0.646.0> closing AMQP connection <0.646.0> (127.0.0.1:58667 -> 127.0.0.1:5672, vhost: '/', user: 'guest'):
client unexpectedly closed TCP connection
```

Abruptly closed connections can be harmless. For example, a short lived program can naturally stop and don't have a chance to close its connection. They can also hint at a genuine issue such as a failed application process or a proxy that closes TCP connections it considers to be idle.

## Upgrading From pre-3.7 Versions  从 3.7 之前的版本升级

RabbitMQ versions prior to 3.7.0 had a different logging subsystem.  RabbitMQ 3.7.0 之前的版本有一个不同的日志子系统。

Older installations use two log files: **\<nodename\>.log** and **\<nodename\>\_sasl.log** (\<nodename\> is rabbit@{hostname} by default).  较旧的安装使用两个日志文件

Where **\<nodename\>.log** contains RabbitMQ logs, while **\<nodename\>\_sasl.log** contains runtime logs, mostly unhandled exceptions.

Starting with 3.7.0 these two files were merged and all errors now can be found in the **\<nodename\>.log** file. So **RABBITMQ_SASL_LOGS** environment variable is not used anymore.

Log levels in versions before 3.7.0 were configured using the **log_levels** configuration key. Starting with 3.7.0 it's been replaced with categories, which are more descriptive and powerful.

If the log_levels key is present in rabbitmq.config file, it should be updated to use categories.

rabbit.log_levels will work in 3.7.0 **only** if no categories are defined.

## Watching Internal Events

RabbitMQ nodes have an internal mechanism. Some of its events can be of interest for monitoring, audit and troubleshooting purposes. They can be consumed as JSON objects using a **rabbitmq-diagnostics** command:  RabbitMQ 节点有一个内部机制。 它的一些事件可能对监控、审计和故障排除目的感兴趣。 可以使用 rabbitmq-diagnostics 命令将它们作为 JSON 对象使用：

```bash
# will emit JSON objects
rabbitmq-diagnostics consume_event_stream
```

When used interactively, results can be piped to a command line JSON processor such as [jq](https://stedolan.github.io/jq/):

```bash
rabbitmq-diagnostics consume_event_stream | jq
```

The events can also be exposed to applications for [consumption](https://www.rabbitmq.com/consumers.html) with a plugin, [rabbitmq-event-exchange](https://github.com/rabbitmq/rabbitmq-event-exchange/).  这些事件也可以通过插件 rabbitmq-event-exchange 暴露给应用程序以供使用。

Events are published as messages with blank bodies. All event metadata is stored in message metadata (properties, headers).  事件作为带有空白正文的消息发布。 所有事件元数据都存储在消息元数据（属性、标头）中。

Below is a list of published events.  以下是已发布事件的列表。

### Core Broker

[Queue](https://www.rabbitmq.com/queues.html), Exchange and Binding events:

- queue.deleted
- queue.created
- exchange.created
- exchange.deleted
- binding.created
- binding.deleted

[Connection](https://www.rabbitmq.com/connections.html) and [Channel](https://www.rabbitmq.com/channels.html) events:

- connection.created
- connection.closed
- channel.created
- channel.closed

[Consumer](https://www.rabbitmq.com/consumers.html) events:

- consumer.created
- consumer.deleted

[Policy and Parameter](https://www.rabbitmq.com/parameters.html) events:

- policy.set
- policy.cleared
- parameter.set
- parameter.cleared

[Virtual host](https://www.rabbitmq.com/vhosts.html) events:

- vhost.created
- vhost.deleted
- vhost.limits.set
- vhost.limits.cleared

User management events:

- user.authentication.success
- user.authentication.failure
- user.created
- user.deleted
- user.password.changed
- user.password.cleared
- user.tags.set

[Permission](https://www.rabbitmq.com/access-control.html) events:

- permission.created
- permission.deleted
- topic.permission.created
- topic.permission.deleted

[Alarm](https://www.rabbitmq.com/alarms.html) events:

- alarm.set
- alarm.cleared

### [Shovel Plugin](https://www.rabbitmq.com/shovel.html)

Worker events:

- shovel.worker.status
- shovel.worker.removed

### [Federation Plugin](https://www.rabbitmq.com/federation.html)

Link events:

- federation.link.status
- federation.link.removed

## Consuming Log Entries Using a System Log Exchange

RabbitMQ can forward log entries to a system exchange, **amq.rabbitmq.log**, which will be declared in the default [virtual host](https://www.rabbitmq.com/vhosts.html).  RabbitMQ 可以将日志条目转发到系统交换 amq.rabbitmq.log，它将在默认虚拟主机中声明。

This feature is disabled by default. To enable this logging, set the **log.exchange** configuration key to true:  默认情况下禁用此功能。 要启用此日志记录，请将 log.exchange 配置键设置为 true：

```ini
# enable log forwarding to amq.rabbitmq.log, a topic exchange
log.exchange = true
```

**log.exchange.level** can be used to control the [log level](https://www.rabbitmq.com/logging.html#log-levels) that will be used by this logging target:

```ini
log.exchange = true
log.exchange.level = warning
```

**amq.rabbitmq.log** is a regular topic exchange and can be used as such. Log entries are published as messages. Message body contains the logged message and routing key is set to the log level.  amq.rabbitmq.log 是一个常规的主题交换，可以这样使用。 日志条目作为消息发布。 消息正文包含记录的消息，并且路由键设置为日志级别。

Application that would like to consume log entries need to declare a queue and bind it to the exchange, using a routing key to filter a specific log level, or # to consume all log entries allowed by the configured log level.  想要消费日志条目的应用程序需要声明一个队列并将其绑定到交换器，使用路由键来过滤特定的日志级别，或者使用 # 来消费配置的日志级别允许的所有日志条目。

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!







```bash

$ rabbitmq-diagnostics --node rabbit@my-rabbitmq-server-0.my-rabbitmq-nodes.hz-rabbitmq --longnames environment

$ env | grep ^RABBITMQ_
# rabbitmq 镜像中设置
RABBITMQ_DATA_DIR=/var/lib/rabbitmq
RABBITMQ_VERSION=3.8.16
RABBITMQ_PGP_KEY_ID=0x0A9AF2115F4687BD29803A206B73A36E6026DFCA
RABBITMQ_HOME=/opt/rabbitmq
RABBITMQ_LOGS=-  # RABBITMQ_LOGS=- 将导致所有日志消息发送到标准输出

# operators 中设置
RABBITMQ_NODENAME=rabbit@my-rabbitmq-server-0.my-rabbitmq-nodes.hz-rabbitmq
RABBITMQ_USE_LONGNAME=true
RABBITMQ_ENABLED_PLUGINS_FILE=/operator/enabled_plugins


log.console = true
log.console.level = debug
log.connection.level = debug
log.channel.level = debug
log.queue.level = debug
log.mirroring.level = debug
log.federation.level = debug
log.upgrade.level = debug
log.default.level = debug

```

