# Configuration

https://www.rabbitmq.com/configure.html

## Overview

RabbitMQ comes with default built-in settings. Those can be entirely sufficient in some environment (e.g. development and QA). For all other cases, as well as [production deployment tuning](https://www.rabbitmq.com/production-checklist.html), there is a way to configure many things in the broker as well as [plugins](https://www.rabbitmq.com/plugins.html).  RabbitMQ 带有默认的内置设置。 在某些环境（例如开发和 QA）中，这些可能是完全足够的。 对于所有其他情况，以及生产部署调整，有一种方法可以在代理和插件中配置许多东西。

This guide covers a number of topics related to configuration:

- Different ways in which various settings of the server and plugins are configured  配置服务器和插件的各种设置的不同方式

- Configuration file(s): primary **rabbitmq.conf** and optional **advanced.config**  配置文件：主rabbitmq.conf 和可选的advanced.config

- Default configuration file location(s) on various platforms  各种平台上的默认配置文件位置

- Configuration troubleshooting: how to find config file location and inspect and verify effective configuration  配置故障排除：如何查找配置文件位置并检查和验证有效配置

- Environment variables  环境变量

- Operating system (kernel) limits  操作系统（内核）限制

- Available core server settings  可用的核心服务器设置

- Available environment variables  可用的环境变量

- How to encrypt sensitive configuration values  如何加密敏感的配置值

and more.

Since configuration affects many areas of the system, including plugins, individual [documentation guides](https://www.rabbitmq.com/documentation.html) dive deeper into what can be configured. [Runtime Tuning](https://www.rabbitmq.com/runtime.html) is a companion to this guide that focuses on the configurable parameters in the runtime. [Production Checklist](https://www.rabbitmq.com/production-checklist.html) is a related guide that outlines what settings will likely need tuning in most production environments.  由于配置影响系统的许多领域，包括插件，单独的文档指南深入探讨了可以配置的内容。 运行时调优是本指南的配套，重点介绍运行时中的可配置参数。 生产清单是一个相关指南，概述了在大多数生产环境中可能需要调整的设置。

## Means of Configuration

A RabbitMQ node can be configured using a number of mechanisms responsible for different areas:

| Mechanism | Description |
| --------- | ----------- |
| Configuration File(s) | contains server and plugin settings for [TCP listeners and other networking-related settings](https://www.rabbitmq.com/networking.html)、 [TLS](https://www.rabbitmq.com/ssl.html)、 [resource constraints (alarms)](https://www.rabbitmq.com/alarms.html)、 [authentication and authorisation backends](https://www.rabbitmq.com/access-control.html)、 [message store settings](https://www.rabbitmq.com/persistence-conf.html) and so on. |
| Environment Variables | define [node name](https://www.rabbitmq.com/cli.html#node-names), file and directory locations, runtime flags taken from the shell, or set in the environment configuration file, **rabbitmq-env.conf** (Linux, MacOS, BSD) and rabbitmq-env-conf.bat (Windows) |
| [rabbitmqctl](https://www.rabbitmq.com/cli.html) | When [internal authentication/authorisation backend](https://www.rabbitmq.com/access-control.html) is used, rabbitmqctl is the tool that manages virtual hosts, users and permissions. It is also used to manage [runtime parameters and policies](https://www.rabbitmq.com/parameters.html). |
| [rabbitmq-queues](https://www.rabbitmq.com/cli.html) | rabbitmq-queues is the tool that manages settings specific to [quorum queues](https://www.rabbitmq.com/quorum-queues.html). |
| [rabbitmq-plugins](https://www.rabbitmq.com/cli.html) | rabbitmq-plugins is the tool that manages [plugins](https://www.rabbitmq.com/plugins.html). |
| [rabbitmq-diagnostics](https://www.rabbitmq.com/cli.html) | rabbitmq-diagnostics allows for inspection of node state, including effective configuration, as well as many other metrics and [health checks](https://www.rabbitmq.com/monitoring.html). |
| [Parameters and Policies](https://www.rabbitmq.com/parameters.html) | defines cluster-wide settings which can change at run time as well as settings that are convenient to configure for groups of queues (exchanges, etc) such as including optional queue arguments. |
| [Runtime (Erlang VM) Flags](https://www.rabbitmq.com/runtime.html) | Control lower-level aspects of the system: memory allocation settings, inter-node communication buffer size, runtime scheduler settings and more. |
| [Operating System Kernel Limits](https://www.rabbitmq.com/configure.html#kernel-limits) | Control process limits enforced by the kernel: [max open file handle limit](https://www.rabbitmq.com/networking.html#open-file-handle-limit), max number of processes and kernel threads, max resident set size and so on. |

Most settings are configured using the first two methods. This guide, therefore, focuses on them.

## Configuration File(s)

### Introduction

While some settings in RabbitMQ can be tuned using environment variables, most are configured using a main configuration file, usually named **rabbitmq.conf**. This includes configuration for the core server as well as plugins. An additional configuration file can be used to configure settings that cannot be expressed in the main file's configuration format. This is covered in more details below.  虽然 RabbitMQ 中的某些设置可以使用环境变量进行调整，但大多数设置使用主配置文件进行配置，通常名为 rabbitmq.conf。 这包括核心服务器和插件的配置。 附加配置文件可用于配置无法以主文件的配置格式表示的设置。 这在下面有更详细的介绍。

The sections below cover the syntax and location of both files, where to find examples, and more.

### Config File Locations

Default config file locations vary between operating systems and [package types](https://www.rabbitmq.com/download.html).

This topic is covered in more detail in the rest of this guide.

When in doubt about RabbitMQ config file location, consult the log file and/or management UI as explained in the following section.

### How to Find Config File Location

The active configuration file can be verified by inspecting the RabbitMQ log file. It will show up in the [log file](https://www.rabbitmq.com/logging.html) at the top, along with the other broker boot log entries. For example:

```bash
node           : rabbit@example
home dir       : /var/lib/rabbitmq
config file(s) : /etc/rabbitmq/advanced.config
               : /etc/rabbitmq/rabbitmq.conf
```

If the configuration file cannot be found or read by RabbitMQ, the log entry will say so:

```
node           : rabbit@example
home dir       : /var/lib/rabbitmq
config file(s) : /var/lib/rabbitmq/hare.conf (not found)
```

Alternatively, the location of configuration files used by a local node, use the [rabbitmq-diagnostics status](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) command:

```bash
# displays key
rabbitmq-diagnostics status
```

and look for the Config files section that would look like this:

```bash
Config files

 * /etc/rabbitmq/advanced.config
 * /etc/rabbitmq/rabbitmq.conf
```

To inspect the locations of a specific node, including nodes running remotely, use the **-n** (short for **--node**) switch:

```bash
rabbitmq-diagnostics status -n [node name]
```

Finally, config file location can be found in the [management UI](https://www.rabbitmq.com/management.html), together with other details about nodes.

When troubleshooting configuration settings, it is very useful to verify that the config file path is correct, exists and can be loaded (e.g. the file is readable) before verifying effective node configuration. Together, these steps help quickly narrow down most common misconfiguration problems.  在对配置设置进行故障排除时，在验证有效的节点配置之前验证配置文件路径是否正确、存在和可以加载（例如文件可读）非常有用。 总之，这些步骤有助于快速缩小最常见的错误配置问题的范围。

### The New and Old Config File Formats

All [supported RabbitMQ versions](https://www.rabbitmq.com/versions.html) use an ini-like, sysctl configuration file format for the main configuration file. The file is typically named **rabbitmq.conf**.

The new config format is much simpler, easier for humans to read and machines to generate. It is also relatively limited compared to the classic config format used prior to RabbitMQ 3.7.0. For example, when configuring [LDAP support](https://www.rabbitmq.com/ldap.html), it may be necessary to use deeply nested data structures to express desired configuration.  新的配置格式更简单，更易于人类阅读和机器生成。 与 RabbitMQ 3.7.0 之前使用的经典配置格式相比，它也相对有限。 例如，在配置 LDAP 支持时，可能需要使用深度嵌套的数据结构来表达所需的配置。

To accommodate this need, modern RabbitMQ versions allow for both formats to be used at the same time in separate files: **rabbitmq.conf** uses the new style format and is recommended for most settings, and **advanced.config** covers more advanced settings that the ini-style configuration cannot express. This is covered in more detail in the following sections.  为了满足这种需求，现代 RabbitMQ 版本允许在单独的文件中同时使用这两种格式：rabbitmq.conf 使用新样式格式并推荐用于大多数设置，advanced.config 涵盖更高级的设置，而不是 风格配置无法表达。 以下各节将对此进行更详细的介绍。

| Configuration File | Format Used | Purpose |
| ------------------ | ----------- | ------- |
| rabbitmq.conf | New style format (sysctl or ini-like) | Primary configuration file. Should be used for most settings. It is easier for humans to read and machines (deployment tools) to generate. Not every setting can be expressed in this format. |
| advanced.config | Classic (Erlang terms) | A limited number of settings that cannot be expressed in the new style configuration format, such as [LDAP queries](https://www.rabbitmq.com/ldap.html). Only should be used when necessary. |
| rabbitmq-env.conf (rabbitmq-env.conf.bat on Windows) | Environment variable pairs | Used to set environment variables relevant to RabbitMQ in one place. |

Compare this examplary **rabbitmq.conf** file

```ini
# A new style format snippet. This format is used by rabbitmq.conf files.
ssl_options.cacertfile           = /path/to/ca_certificate.pem
ssl_options.certfile             = /path/to/server_certificate.pem
ssl_options.keyfile              = /path/to/server_key.pem
ssl_options.verify               = verify_peer
ssl_options.fail_if_no_peer_cert = true
```

to

```erlang
%% A classic format snippet, now used by advanced.config files.
[
  {rabbit, [{ssl_options, [{cacertfile,           "/path/to/ca_certificate.pem"},
                           {certfile,             "/path/to/server_certificate.pem"},
                           {keyfile,              "/path/to/server_key.pem"},
                           {verify,               verify_peer},
                           {fail_if_no_peer_cert, true}]}]}
].
```

### The Main Configuration File, rabbitmq.conf

The configuration file **rabbitmq.conf** allows the RabbitMQ server and plugins to be configured. Starting with RabbitMQ 3.7.0, the format is in the [sysctl format](https://github.com/basho/cuttlefish/wiki/Cuttlefish-for-Application-Users).

The syntax can be briefly explained in 3 lines:

- One setting uses one line
- Lines are structured Key = Value
- Any line starting with a # character is a comment

A minimalistic example configuration file follows:

```ini
# this is a comment
listeners.tcp.default = 5673
```

The same example in the classic config format:

```erlang
%% this is a comment
[
  {rabbit, [
      {tcp_listeners, [5673]}
    ]
  }
].
```

This example will alter the [port RabbitMQ listens on](https://www.rabbitmq.com/networking.html#ports) for AMQP 0-9-1 and AMQP 1.0 client connections from 5672 to 5673.

The RabbitMQ server source repository contains [an example rabbitmq.conf file](https://github.com/rabbitmq/rabbitmq-server/blob/v3.8.x/deps/rabbit/docs/rabbitmq.conf.example) named rabbitmq.conf.example. It contains examples of most of the configuration items you might want to set (with some very obscure ones omitted), along with documentation for those settings.

Documentation guides such as [Networking](https://www.rabbitmq.com/networking.html), [TLS](https://www.rabbitmq.com/ssl.html), or [Access Control](https://www.rabbitmq.com/access-control.html) contain many examples in relevant formats.

Note that this configuration file is not to be confused with the environment variable configuration files, **rabbitmq-env.conf** and **rabbitmq-env-conf.bat**.  请注意，不要将此配置文件与环境变量配置文件 rabbitmq-env.conf 和 rabbitmq-env-conf.bat 混淆。

To override the main RabbitMQ config file location, use the **RABBITMQ_CONFIG_FILE** environment variable. Use .conf as file extension for the new style config format, e.g. **/etc/rabbitmq/rabbitmq.conf** or **/data/configuration/rabbitmq/rabbitmq.conf**

### The advanced.config File

Some configuration settings are not possible or are difficult to configure using the sysctl format. As such, it is possible to use an additional config file in the Erlang term format (same as rabbitmq.config). That file is commonly named **advanced.config**. It will be merged with the configuration provided in rabbitmq.conf.

The RabbitMQ server source repository contains [an example advanced.config file](https://github.com/rabbitmq/rabbitmq-server/blob/v3.8.x/deps/rabbit/docs/advanced.config.example) named advanced.config.example. It focuses on the options that are typically set using the advanced config.

To override the advanced config file location, use the **RABBITMQ_ADVANCED_CONFIG_FILE ** environment variable.

### Location of rabbitmq.conf, advanced.config and rabbitmq-env.conf

Default configuration file location is distribution-specific. RabbitMQ packages or nodes will not create any configuration files. Users and deployment tool should use the following locations when creating the files:

| Platform | Default Configuration File Directory | Example Configuration File Paths |
| -------- | ------------------------------------ | ------------------------------- |
| Generic binary package | $RABBITMQ_HOME/etc/rabbitmq/ | $RABBITMQ_HOME/etc/rabbitmq/rabbitmq.conf, $RABBITMQ_HOME/etc/rabbitmq/advanced.config |
| Debian and Ubuntu | /etc/rabbitmq | /etc/rabbitmq/rabbitmq.conf, /etc/rabbitmq/advanced.config |
| RPM-based Linux | /etc/rabbitmq/ | /etc/rabbitmq/rabbitmq.conf, /etc/rabbitmq/advanced.config |
| Windows | %APPDATA%\RabbitMQ\ | %APPDATA%\RabbitMQ\rabbitmq.conf, %APPDATA%\RabbitMQ\advanced.config |
| MacOS Homebrew Formula | ${install_prefix}/etc/rabbitmq/, and the Homebrew cellar prefix is usually /usr/local | ${install_prefix}/etc/rabbitmq/rabbitmq.conf, ${install_prefix}/etc/rabbitmq/advanced.config |

Environment variables can be used to override the location of the configuration file:

```ini
# overrides primary config file location
RABBITMQ_CONFIG_FILE=/path/to/a/custom/location/rabbitmq.conf

# overrides advanced config file location
RABBITMQ_ADVANCED_CONFIG_FILE=/path/to/a/custom/location/advanced.config

# overrides environment variable file location
RABBITMQ_CONF_ENV_FILE=/path/to/a/custom/location/rabbitmq-env.conf
```

### When Will Configuration File Changes Be Applied

rabbitmq.conf and advanced.config changes take effect after a node restart.

If rabbitmq-env.conf doesn't exist, it can be created manually in the location specified by the **RABBITMQ_CONF_ENV_FILE** variable. On Windows systems, it is named rabbitmq-env-conf.bat.  如果 rabbitmq-env.conf 不存在，则可以在 RABBITMQ_CONF_ENV_FILE 变量指定的位置手动创建它。 在 Windows 系统上，它被命名为 rabbitmq-env-conf.bat。

Windows service users will need to re-install the service if configuration file location or any values in rabbitmq-env-conf.bat have changed. Environment variables used by the service would not be updated otherwise.  如果配置文件位置或 rabbitmq-env-conf.bat 中的任何值发生更改，Windows 服务用户将需要重新安装该服务。 否则不会更新服务使用的环境变量。

In the context of deployment automation this means that environment variables such as RABBITMQ_BASE and RABBITMQ_CONFIG_FILE should ideally be set before RabbitMQ is installed. This would help avoid unnecessary confusion and Windows service re-installations.  在部署自动化的上下文中，这意味着最好在安装 RabbitMQ 之前设置 RABBITMQ_BASE 和 RABBITMQ_CONFIG_FILE 等环境变量。 这将有助于避免不必要的混淆和 Windows 服务重新安装。

### How to Inspect and Verify Effective Configuration of a Running Node

It is possible to print effective configuration (user provided values from all configuration files merged into defaults) using the [rabbitmq-diagnostics environment](https://www.rabbitmq.com/rabbitmq-diagnostics.8.html) command:

```bash
# inspect effective configuration on a node
rabbitmq-diagnostics environment
```

to check effective configuration of a specific node, including nodes running remotely, use the **-n** (short for **--node**) switch:

```bash
rabbitmq-diagnostics environment -n [node name]
```

The command above will print applied configuration for every application (RabbitMQ, plugins, libraries) running on the node. Effective configuration is computed using the following steps:

- rabbitmq.conf is translated into the internally used (advanced) config format. These configuration is merged into the defaults

- advanced.config is loaded if present, and merged into the result of the step above

Effective configuration should be verified together with config file location. Together, these steps help quickly narrow down most common misconfiguration problems.

### The rabbitmq.config (Classic Format) File

Prior to RabbitMQ 3.7.0, RabbitMQ config file was named rabbitmq.config and used [the same Erlang term format](http://www.erlang.org/doc/man/config.html) used by advanced.config today. That format is still supported for backwards compatibility.

The classic format is **deprecated**. Please prefer the new style config format in rabbitmq.conf accompanied by an advanced.config file as needed.

To use a config file in the classic format, export RABBITMQ_CONFIG_FILE to point to the file with a .config extension. The extension will indicate to RabbitMQ that it should treat the file as one in the classic config format.

[An example configuration file](https://github.com/rabbitmq/rabbitmq-server/blob/v3.7.x/deps/rabbit/docs/rabbitmq.config.example) named rabbitmq.config.example. It contains an example of most of the configuration items in the classic config format.

To override the main RabbitMQ config file location, use the RABBITMQ_CONFIG_FILE environment variable. Use .config as file extension for the classic config format.

The use of classic config format should only be limited to the advanced.config file and settings that cannot be configured using the ini-style config file.

### Example Configuration Files

The RabbitMQ server source repository contains examples for the configuration files:

- [rabbitmq.conf.example](https://github.com/rabbitmq/rabbitmq-server/blob/master/deps/rabbit/docs/rabbitmq.conf.example)

- [advanced.config.example](https://github.com/rabbitmq/rabbitmq-server/blob/master/deps/rabbit/docs/advanced.config.example)

These files contain examples of most of the configuration keys along with a brief explanation for those settings. All configuration items are commented out in the example, so you can uncomment what you need. Note that the example files are meant to be used as, well, examples, and should not be treated as a general recommendation.  这些文件包含大多数配置键的示例以及这些设置的简要说明。 示例中所有配置项都已注释掉，因此您可以取消注释您需要的内容。 请注意，示例文件旨在用作示例，不应将其视为一般建议。

In most distributions the example file is placed into the same location as the real file should be placed (see above). On Debian and RPM distributions policy forbids doing so; instead find the file under /usr/share/doc/rabbitmq-server/ or /usr/share/doc/rabbitmq-server-3.8.19/, respectively.

### Core Server Variables Configurable in rabbitmq.conf

These variables are the most common. The list is not complete, as some settings are quite obscure.  这些变量是最常见的。 该列表并不完整，因为有些设置相当模糊。

| Key | Documentation |
| --- | ------------- |
| listeners | Ports or hostname/pair on which to listen for "plain" AMQP 0-9-1 and AMQP 1.0 connections (without [TLS](https://www.rabbitmq.com/ssl.html)). See the [Networking guide](https://www.rabbitmq.com/networking.html) for more details and examples.Default:`listeners.tcp.default = 5672` |
| num_acceptors.tcp | Number of Erlang processes that will accept connections for the TCP listeners.Default:`num_acceptors.tcp = 10` |
| handshake_timeout | Maximum time for AMQP 0-9-1 handshake (after socket connection and TLS handshake), in milliseconds.Default:`handshake_timeout = 10000` |
| listeners.ssl | Ports or hostname/pair on which to listen for TLS-enabled AMQP 0-9-1 and AMQP 1.0 connections. See the [TLS guide](https://www.rabbitmq.com/ssl.html) for more details and examples.Default: none (not set) |
| num_acceptors.ssl | Number of Erlang processes that will accept TLS connections from clients.Default:`num_acceptors.ssl = 10` |
| ssl_options | TLS configuration. See the [TLS guide](https://www.rabbitmq.com/ssl.html#enabling-ssl).Default:`ssl_options = none` |
| ssl_handshake_timeout | TLS handshake timeout, in milliseconds.Default:`ssl_handshake_timeout = 5000` |
| vm_memory_high_watermark | Memory threshold at which the flow control is triggered. Can be absolute or relative to the amount of RAM available to the OS:`vm_memory_high_watermark.relative = 0.6` `vm_memory_high_watermark.absolute = 2GB` See the [memory-based flow control](https://www.rabbitmq.com/memory.html) and [alarms](https://www.rabbitmq.com/alarms.html) documentation.Default:`vm_memory_high_watermark.relative = 0.4` |
| vm_memory_calculation_strategy | Strategy for memory usage reporting. Can be one of the following: **allocated**: uses Erlang memory allocator statistics，**rss**: uses operating system RSS memory reporting. This uses OS-specific means and may start short lived child processes，**legacy**: uses legacy memory reporting (how much memory is considered to be used by the runtime). This strategy is fairly inaccurate，**erlang**: same as legacy, preserved for backwards compatibility。Default:`vm_memory_calculation_strategy = allocated ` |
| vm_memory_high_watermark_paging_ratio | Fraction of the high watermark limit at which queues start to page messages out to disc to free up memory. See the [memory-based flow control](https://www.rabbitmq.com/memory.html) documentation. Default:`vm_memory_high_watermark_paging_ratio = 0.5` |
| total_memory_available_override_value | Makes it possible to override the total amount of memory available, as opposed to inferring it from the environment using OS-specific means. This should only be used when actual maximum amount of RAM available to the node doesn't match the value that will be inferred by the node, e.g. due to containerization or similar constraints the node cannot be aware of. The value may be set to an integer number of bytes or, alternatively, in information units (e.g `8GB`). For example, when the value is set to 4 GB, the node will believe it is running on a machine with 4 GB of RAM.Default: undefined (not set or used). 使覆盖可用内存总量成为可能，而不是使用特定于操作系统的方法从环境中推断出来。 仅当节点可用的实际最大 RAM 量与节点推断的值不匹配时才应使用此方法，例如 由于容器化或类似的限制，节点无法意识到。 该值可以设置为整数字节，或者以信息单位（例如`8GB`）。 例如，当该值设置为 4 GB 时，节点将认为它在具有 4 GB RAM 的机器上运行。|
| disk_free_limit | Disk free space limit of the partition on which RabbitMQ is storing data. When available disk space falls below this limit, flow control is triggered. The value can be set relative to the total amount of RAM or as an absolute value in bytes or, alternatively, in information units (e.g `50MB` or `5GB`):`disk_free_limit.relative = 3.0` `disk_free_limit.absolute = 2GB` By default free disk space must exceed 50MB. See the [Disk Alarms](https://www.rabbitmq.com/disk-alarms.html) documentation.Default:`disk_free_limit.absolute = 50MB` |
| log.file.level | Controls the granularity of logging. The value is a list of log event category and log level pairs.The level can be one of **error** (only errors are logged), **warning** (only errors and warning are logged), **info** (errors, warnings and informational messages are logged), or **debug** (errors, warnings, informational messages and debugging messages are logged). Default:`log.file.level = info` |
| channel_max | Maximum permissible number of channels to negotiate with clients, not including a special channel number 0 used in the protocol. Setting to 0 means "unlimited", a dangerous value since applications sometimes have channel leaks. Using more channels increases memory footprint of the broker.Default:`channel_max = 2047` |
| channel_operation_timeout | Channel operation timeout in milliseconds (used internally, not directly exposed to clients due to messaging protocol differences and limitations).Default:`channel_operation_timeout = 15000` |
| max_message_size | The largest allowed message payload size in bytes. Messages of larger size will be rejected with a suitable channel exception.Default: 134217728  Max value: 536870912 |
| heartbeat | Value representing the heartbeat timeout suggested by the server during connection parameter negotiation. If set to 0 on both ends, heartbeats are disabled (this is not recommended). See the [Heartbeats guide](https://www.rabbitmq.com/heartbeats.html) for details.Default:`heartbeat = 60 ` |
| default_vhost | Virtual host to create when RabbitMQ creates a new database from scratch. The exchange `amq.rabbitmq.log` will exist in this virtual host.Default:`default_vhost = / ` |
| default_user | User name to create when RabbitMQ creates a new database from scratch.Default:`default_user = guest ` |
| default_pass | Password for the default user.Default:`default_pass = guest ` |
| default_user_tags | Tags for the default user.Default:`default_user_tags.administrator = true ` |
| default_permissions | [Permissions](https://www.rabbitmq.com/access-control.html) to assign to the default user when creating it.Default:`default_permissions.configure = .* default_permissions.read = .* default_permissions.write = .* ` |
| loopback_users | List of users which are only permitted to connect to the broker via a loopback interface (i.e. `localhost`).To allow the default `guest` user to connect remotely (a security practice [unsuitable for production use](https://www.rabbitmq.com/production-checklist.html)), set this to `none`:`# awful security practice, # consider creating a new # user with secure generated credentials! loopback_users = none `To restrict another user to localhost-only connections, do it like so (`monitoring` is the name of the user):`loopback_users.monitoring = true `Default:`# guest uses well known # credentials and can only # log in from localhost # by default loopback_users.guest = true ` |
| cluster_formation.classic_config.nodes | Classic [peer discovery](https://www.rabbitmq.com/cluster-formation.html) backend's list of nodes to contact. For example, to cluster with nodes `rabbit@hostname1` and `rabbit@hostname2` on first boot:`cluster_formation.classic_config.nodes.1 = rabbit@hostname1 cluster_formation.classic_config.nodes.2 = rabbit@hostname2 `Default: `none` (not set) |
| collect_statistics                     | Statistics collection mode. Primarily relevant for the management plugin. Options are:`none` (do not emit statistics events)`coarse` (emit per-queue / per-channel / per-connection statistics)`fine` (also emit per-message statistics)Default:`collect_statistics = none ` |
| collect_statistics_interval            | Statistics collection interval in milliseconds. Primarily relevant for the [management plugin](https://www.rabbitmq.com/management.html#statistics-interval).Default:`collect_statistics_interval = 5000 ` |
| management_db_cache_multiplier         | Affects the amount of time the [management plugin](https://www.rabbitmq.com/management.html#statistics-interval) will cache expensive management queries such as queue listings. The cache will multiply the elapsed time of the last query by this value and cache the result for this amount of time.Default:`management_db_cache_multiplier = 5 ` |
| auth_mechanisms                        | [SASL authentication mechanisms](https://www.rabbitmq.com/authentication.html) to offer to clients.Default:`auth_mechanisms.1 = PLAIN auth_mechanisms.2 = AMQPLAIN ` |
| auth_backends                          | List of [authentication and authorisation backends](https://www.rabbitmq.com/access-control.html) to use. See the [access control guide](https://www.rabbitmq.com/access-control.html) for details and examples.Other databases than `rabbit_auth_backend_internal` are available through [plugins](https://www.rabbitmq.com/plugins.html).Default:`auth_backends.1 = internal` |
| reverse_dns_lookups                    | Set to `true` to have RabbitMQ perform a reverse DNS lookup on client connections, and present that information through `rabbitmqctl` and the management plugin.Default:`reverse_dns_lookups = false` |
| delegate_count                         | Number of delegate processes to use for intra-cluster communication. On a machine which has a very large number of cores and is also part of a cluster, you may wish to increase this value.Default:`delegate_count = 16` |
| tcp_listen_options                     | Default socket options. You probably don't want to change this.Default:`tcp_listen_options.backlog = 128 tcp_listen_options.nodelay = true tcp_listen_options.linger.on = true tcp_listen_options.linger.timeout = 0 tcp_listen_options.exit_on_close = false ` |
| hipe_compile                           | Do not use. This option is no longer supported. HiPE supported has been dropped starting with Erlang 22.Default:`hipe_compile = false` |
| cluster_partition_handling             | How to handle network partitions. Available modes are:ignoreautohealpause_minoritypause_if_all_downpause_if_all_down mode requires additional parameters:nodesrecoverSee the [documentation on partitions](https://www.rabbitmq.com/partitions.html#automatic-handling) for more information.Default:`cluster_partition_handling = ignore` |
| cluster_keepalive_interval             | How frequently nodes should send keepalive messages to other nodes (in milliseconds). Note that this is not the same thing as [net_ticktime](https://www.rabbitmq.com/nettick.html); missed keepalive messages will not cause nodes to be considered down.Default:`cluster_keepalive_interval = 10000 ` |
| queue_index_embed_msgs_below           | Size in bytes of message below which messages will be embedded directly in the queue index. You are advised to read the [persister tuning](https://www.rabbitmq.com/persistence-conf.html) documentation before changing this.Default:`queue_index_embed_msgs_below = 4096 ` |
| mnesia_table_loading_retry_timeout     | Timeout used when waiting for Mnesia tables in a cluster to become available.Default:`mnesia_table_loading_retry_timeout = 30000 ` |
| mnesia_table_loading_retry_limit       | Retries when waiting for Mnesia tables in the cluster startup. Note that this setting is not applied to Mnesia upgrades or node deletions.Default:`mnesia_table_loading_retry_limit = 10 ` |
| mirroring_sync_batch_size              | Batch size used to transfer messages to an unsynchronised replica (queue mirror). See [documentation on eager batch synchronization](https://www.rabbitmq.com/ha.html#batch-sync).Default:`mirroring_sync_batch_size = 4096 ` |
| queue_master_locator                   | queue leader location strategy. Available strategies are: **min-masters**、**client-local**、**random**。See the [documentation on queue leader location](https://www.rabbitmq.com/ha.html#queue-master-location) for more information.Default:`queue_master_locator = client-local ` |
| proxy_protocol                         | If set to true, RabbitMQ will expect a [proxy protocol](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) header to be sent first when an AMQP connection is opened. This implies to set up a proxy protocol-compliant reverse proxy (e.g. [HAproxy](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) or [AWS ELB](http://docs.aws.amazon.com/elasticloadbalancing/latest/classic/enable-proxy-protocol.html)) in front of RabbitMQ. Clients can't directly connect to RabbitMQ when proxy protocol is enabled, so all connections must go through the reverse proxy.See [the networking guide](https://www.rabbitmq.com/networking.html#proxy-protocol) for more information.Default:`proxy_protocol = false ` |
| cluster_name                           | Operator-controlled cluster name. This name is used to identify a cluster, and by the federation and Shovel plugins to record the origin or path of transferred messages. Can be set to any arbitrary string to help identify the cluster (eg. london). This name can be inspected by AMQP 0-9-1 clients in the server properties map. Default: by default the name is derived from the first (seed) node in the cluster. |

The following configuration settings can be set in the [advanced config file](https://www.rabbitmq.com/configure.html#advanced-config-file) only, under the rabbit section.

| **Key**                                        | **Documentation**                                            |
| :--------------------------------------------- | :----------------------------------------------------------- |
| msg_store_index_module                         | Implementation module for queue indexing. You are advised to read the [message store tuning](https://www.rabbitmq.com/persistence-conf.html) documentation before changing this.Default: rabbit_msg_store_ets_index`{rabbit, [ {msg_store_index_module, rabbit_msg_store_ets_index} ]} ` |
| backing_queue_module                           | Implementation module for queue contents.Default:`{rabbit, [ {backing_queue_module, rabbit_variable_queue} ]} ` |
| msg_store_file_size_limit                      | Message store segment file size. Changing this for a node with an existing (initialised) database is dangerous can lead to data loss!Default: 16777216`{rabbit, [ %% Changing this for a node %% with an existing (initialised) database is dangerous can %% lead to data loss! {msg_store_file_size_limit, 16777216} ]} ` |
| trace_vhosts                                   | Used internally by the [tracer](https://www.rabbitmq.com/firehose.html). You shouldn't change this.Default:`{rabbit, [ {trace_vhosts, []} ]} ` |
| msg_store_credit_disc_bound                    | The credits that a queue process is given by the message store.By default, a queue process is given 4000 message store credits, and then 800 for every 800 messages that it processes.Messages which need to be paged out due to memory pressure will also use this credit.The Message Store is the last component in the credit flow chain. [Learn about credit flow.](https://blog.rabbitmq.com/posts/2015/10/new-credit-flow-settings-on-rabbitmq-3-5-5/)This value only takes effect when messages are persisted to the message store. If messages are embedded on the queue index, then modifying this setting has no effect because credit_flow is NOT used when writing to the queue index.Default:`{rabbit, [ {msg_store_credit_disc_bound, {4000, 800}} ]} ` |
| queue_index_max_journal_entries                | After how many queue index journal entries it will be flushed to disk.Default:`{rabbit, [ {queue_index_max_journal_entries, 32768} ]} ` |
| lazy_queue_explicit_gc_run_operation_threshold | Tunable value only for lazy queues when under memory pressure. This is the threshold at which the garbage collector and other memory reduction activities are triggered. A low value could reduce performance, and a high one can improve performance, but cause higher memory consumption. You almost certainly should not change this.Default:`{rabbit, [ {lazy_queue_explicit_gc_run_operation_threshold, 1000} ]} ` |
| queue_explicit_gc_run_operation_threshold      | Tunable value only for normal queues when under memory pressure. This is the threshold at which the garbage collector and other memory reduction activities are triggered. A low value could reduce performance, and a high one can improve performance, but cause higher memory consumption. You almost certainly should not change this.Default:`{rabbit, [ {queue_explicit_gc_run_operation_threshold, 1000} ]} ` |

Several [plugins](https://www.rabbitmq.com/plugins.html) that ship with RabbitMQ have dedicated documentation guides that cover plugin configuration:

- [rabbitmq_management](https://www.rabbitmq.com/management.html#configuration)
- [rabbitmq_management_agent](https://www.rabbitmq.com/management.html#configuration)
- [rabbitmq_stomp](https://www.rabbitmq.com/stomp.html)
- [rabbitmq_mqtt](https://www.rabbitmq.com/mqtt.html)
- [rabbitmq_shovel](https://www.rabbitmq.com/shovel.html)
- [rabbitmq_federation](https://www.rabbitmq.com/federation.html)
- [rabbitmq_auth_backend_ldap](https://www.rabbitmq.com/ldap.html)

### [Configuration Value Encryption](https://www.rabbitmq.com/configure.html#configuration-encryption)

Sensitive configuration entries (e.g. password, URL containing credentials) can be encrypted in the RabbitMQ configuration file. The broker decrypts encrypted entries on start.

Note that encrypted configuration entries don't make the system meaningfully more secure. Nevertheless, they allow deployments of RabbitMQ to conform to regulations in various countries requiring that no sensitive data should appear in plain text in configuration files.

Encrypted values must be inside an Erlang encrypted tuple: {encrypted, ...}. Here is an example of a configuration file with an encrypted password for the default user:

```erlang
[
  {rabbit, [
      {default_user, <<"guest">>},
      {default_pass,
        {encrypted,
         <<"cPAymwqmMnbPXXRVqVzpxJdrS8mHEKuo2V+3vt1u/fymexD9oztQ2G/oJ4PAaSb2c5N/hRJ2aqP/X0VAfx8xOQ==">>
        }
      },
      {config_entry_decoder, [
             {passphrase, <<"mypassphrase">>}
         ]}
    ]}
].
```

Note the config_entry_decoder key with the passphrase that RabbitMQ will use to decrypt encrypted values.

The passphrase doesn't have to be hardcoded in the configuration file, it can be in a separate file:

```erlang
[
  {rabbit, [
      %% ...
      {config_entry_decoder, [
             {passphrase, {file, "/path/to/passphrase/file"}}
         ]}
    ]}
].
```

RabbitMQ can also request an operator to enter the passphrase when it starts by using {passphrase, prompt}.

Use [rabbitmqctl](https://www.rabbitmq.com/cli.html) and the encode command to encrypt values:

```bash
rabbitmqctl encode '<<"guest">>' mypassphrase
{encrypted,<<"... long encrypted value...">>}
rabbitmqctl encode '"amqp://fred:secret@host1.domain/my_vhost"' mypassphrase
{encrypted,<<"... long encrypted value...">>}
```

Or, on Windows:

```powershell
rabbitmqctl encode "<<""guest"">>" mypassphrase
{encrypted,<<"... long encrypted value...">>}
rabbitmqctl encode '"amqp://fred:secret@host1.domain/my_vhost"' mypassphrase
{encrypted,<<"... long encrypted value...">>}
```

Add the decode command if you want to decrypt values:

```bash
rabbitmqctl decode '{encrypted, <<"...">>}' mypassphrase
<<"guest">>
rabbitmqctl decode '{encrypted, <<"...">>}' mypassphrase
"amqp://fred:secret@host1.domain/my_vhost"
```

Or, on Windows:

```powershell
rabbitmqctl decode "{encrypted, <<""..."">>}" mypassphrase
<<"guest">>
rabbitmqctl decode "{encrypted, <<""..."">>}" mypassphrase
"amqp://fred:secret@host1.domain/my_vhost"
```

Values of different types can be encoded. The example above encodes both binaries (<<"guest">>) and strings ("amqp://fred:secret@host1.domain/my_vhost").

The encryption mechanism uses PBKDF2 to produce a derived key from the passphrase. The default hash function is SHA512 and the default number of iterations is 1000. The default cipher is AES 256 CBC.

These defaults can be changed in the configuration file:

```erlang
[
  {rabbit, [
      ...
      {config_entry_decoder, [
             {passphrase, "mypassphrase"},
             {cipher, blowfish_cfb64},
             {hash, sha256},
             {iterations, 10000}
         ]}
    ]}
].
```

Or using [CLI tools](https://www.rabbitmq.com/cli.html):

```bash
rabbitmqctl encode --cipher blowfish_cfb64 --hash sha256 --iterations 10000 \
                     '<<"guest">>' mypassphrase
```

Or, on Windows:

```powershell
rabbitmqctl encode --cipher blowfish_cfb64 --hash sha256 --iterations 10000 \
                     "<<""guest"">>" mypassphrase
```

## [Configuration Using Environment Variables](https://www.rabbitmq.com/configure.html#customise-environment)

Certain server parameters can be configured using environment variables: [node name](https://www.rabbitmq.com/cli.html#node-names), RabbitMQ [configuration file location](https://www.rabbitmq.com/configure.html#configuration-files), [inter-node communication ports](https://www.rabbitmq.com/networking.html#ports), Erlang VM flags, and so on.

### [Path and Directory Name Restrictions](https://www.rabbitmq.com/configure.html#directory-and-path-restrictions)

Some of the environment variable configure paths and locations (node's base or data directory, [plugin source and expansion directories](https://www.rabbitmq.com/plugins.html), and so on). Those paths have must exclude a number of characters:

- \* and ? (on Linux, macOS, BSD and other UNIX-like systems)
- ^ and ! (on Windows)
- [ and ]
- { and }

The above characters will render the node unable to start or function as expected (e.g. expand plugins and load their metadata).

### [Linux, MacOS, BSD](https://www.rabbitmq.com/configure.html#environment-env-file-unix)

On UNIX-based systems (Linux, MacOS and flavours of BSD) it is possible to use a file named rabbitmq-env.conf to define environment variables that will be used by the broker. Its [location](https://www.rabbitmq.com/configure.html#config-location) is configurable using the RABBITMQ_CONF_ENV_FILE environment variable.

rabbitmq-env.conf uses the standard environment variable names but without the RABBITMQ_ prefix. For example, the RABBITMQ_CONFIG_FILE variable appears below as CONFIG_FILE and RABBITMQ_NODENAME becomes NODENAME:

```bash
# Example rabbitmq-env.conf file entries. Note that the variables
# do not have the RABBITMQ_ prefix.
#
# Overrides node name
NODENAME=bunny@myhost

# Specifies new style config file location
CONFIG_FILE=/etc/rabbitmq/rabbitmq.conf

# Specifies advanced config file location
ADVANCED_CONFIG_FILE=/etc/rabbitmq/advanced.config
```

See the [rabbitmq-env.conf man page](https://www.rabbitmq.com/man/rabbitmq-env.conf.5.html) for details.

### [Windows](https://www.rabbitmq.com/configure.html#rabbitmq-env-file-windows)

The easiest option to customise names, ports or locations is to configure environment variables in the Windows dialogue: Start > Settings > Control Panel > System > Advanced > Environment Variables. Then create or edit the system variable name and value.

Alternatively it is possible to use a file named rabbitmq-env-conf.bat to define environment variables that will be used by the broker. Its [location](https://www.rabbitmq.com/configure.html#config-location) is configurable using the RABBITMQ_CONF_ENV_FILE environment variable.

Windows service users will need to **re-install the service** if configuration file location or any values in `rabbitmq-env-conf.bat changed. Environment variables used by the service would not be updated otherwise.

This can be done using the installer or on the command line with administrator permissions:

- Start an [administrative command prompt](https://technet.microsoft.com/en-us/library/cc947813(v=ws.10).aspx)

- cd into the sbin folder under the RabbitMQ server installation directory (such as C:\Program Files (x86)\RabbitMQ Server\rabbitmq_server-{version}\sbin)

- Run rabbitmq-service.bat stop to stop the service

- Run rabbitmq-service.bat remove to remove the Windows service (this will *not* remove RabbitMQ or its data directory)

- Set environment variables via command line, i.e. run commands like the following:

  ```powershell
  set RABBITMQ_BASE=C:\Data\RabbitMQ
  ```

- Run rabbitmq-service.bat install

- Run rabbitmq-service.bat start

This will restart the node in a way that makes the environment variable and rabbitmq-env-conf.bat changes to be observable to it.

## [Environment Variables Used by RabbitMQ](https://www.rabbitmq.com/configure.html#supported-environment-variables)

All environment variables used by RabbitMQ use the prefix RABBITMQ_ (except when defined in [rabbitmq-env.conf](https://www.rabbitmq.com/configure.html#environment-env-file-unix) or [rabbitmq-env-conf.bat](https://www.rabbitmq.com/configure.html#environment-env-file-windows)).

Environment variables set in the shell environment take priority over those set in [rabbitmq-env.conf](https://www.rabbitmq.com/configure.html#environment-env-file-unix) or [rabbitmq-env-conf.bat](https://www.rabbitmq.com/configure.html#environment-env-file-windows), which in turn override RabbitMQ built-in defaults.

The table below describes key environment variables that can be used to configure RabbitMQ. More variables are covered in the [File and Directory Locations guide](https://www.rabbitmq.com/relocate.html).

| Name                                | Description                                                  |
| :---------------------------------- | :----------------------------------------------------------- |
| RABBITMQ_NODE_IP_ADDRESS            | Change this if you only want to bind to one network interface. Binding to two or more interfaces can be set up in the configuration file.**Default**: an empty string, meaning "bind to all network interfaces". |
| RABBITMQ_NODE_PORT                  | See [Networking guide](https://www.rabbitmq.com/networking.html) for more information on ports used by various parts of RabbitMQ.**Default**: 5672. |
| RABBITMQ_DIST_PORT                  | Port used for inter-node and CLI tool communication. Ignored if node config file sets kernel.inet_dist_listen_min or kernel.inet_dist_listen_max keys. See [Networking](https://www.rabbitmq.com/networking.html) for details, and [Windows Quirks](https://www.rabbitmq.com/windows-quirks.html) for Windows-specific details.**Default**: RABBITMQ_NODE_PORT + 20000 |
| ERL_EPMD_ADDRESS                    | Interface(s) used by [epmd](https://www.rabbitmq.com/networking.html#epmd), a component in inter-node and CLI tool communication.**Default**: all available interfaces, both IPv6 and IPv4. |
| ERL_EPMD_PORT                       | Port used by [epmd](https://www.rabbitmq.com/networking.html#epmd), a component in inter-node and CLI tool communication.**Default**: 4369 |
| RABBITMQ_DISTRIBUTION_BUFFER_SIZE   | [Outgoing data buffer size limit](https://erlang.org/doc/man/erl.html#+zdbbl) to use for inter-node communication connections, in kilobytes. Values lower than 64 MB are not recommended.**Default**: 128000 |
| RABBITMQ_NODENAME                   | The node name should be unique per Erlang-node-and-machine combination. To run multiple nodes, see the [clustering guide](https://www.rabbitmq.com/clustering.html).**Default**:**Unix\*:** rabbit@$HOSTNAME**Windows:** rabbit@%COMPUTERNAME% |
| RABBITMQ_CONFIG_FILE                | Main RabbitMQ config file path, for example, /etc/rabbitmq/rabbitmq.conf or /data/configuration/rabbitmq.conf for new style configuration format files. If classic config format it used, the extension must be .config**Default**:**Generic UNIX**: $RABBITMQ_HOME/etc/rabbitmq/rabbitmq.conf**Debian**: /etc/rabbitmq/rabbitmq.conf**RPM**: /etc/rabbitmq/rabbitmq.conf**MacOS(Homebrew)**: ${install_prefix}/etc/rabbitmq/rabbitmq.conf, the Homebrew prefix is usually /usr/local**Windows**: %APPDATA%\RabbitMQ\rabbitmq.conf |
| RABBITMQ_ADVANCED_CONFIG_FILE       | "Advanced" (Erlang term-based) RabbitMQ config file path with a .config file extension. For example, /data/rabbitmq/advanced.config.**Default**:**Generic UNIX**: $RABBITMQ_HOME/etc/rabbitmq/advanced.config**Debian**: /etc/rabbitmq/advanced.config**RPM**: /etc/rabbitmq/advanced.config**MacOS (Homebrew)**: ${install_prefix}/etc/rabbitmq/advanced.config, the Homebrew prefix is usually /usr/local**Windows**: %APPDATA%\RabbitMQ\advanced.config |
| RABBITMQ_CONF_ENV_FILE              | Location of the file that contains environment variable definitions (without the RABBITMQ_ prefix). Note that the file name on Windows is different from other operating systems.**Default**:**Generic UNIX package**: $RABBITMQ_HOME/etc/rabbitmq/rabbitmq-env.conf**Ubuntu and Debian**: /etc/rabbitmq/rabbitmq-env.conf**RPM**: /etc/rabbitmq/rabbitmq-env.conf**MacOS (Homebrew)**: ${install_prefix}/etc/rabbitmq/rabbitmq-env.conf, the Homebrew prefix is usually /usr/local**Windows**: %APPDATA%\RabbitMQ\rabbitmq-env-conf.bat |
| RABBITMQ_MNESIA_BASE                | This base directory contains sub-directories for the RabbitMQ server's node database, message store and cluster state files, one for each node, unless **RABBITMQ_MNESIA_DIR** is set explicitly. It is important that effective RabbitMQ user has sufficient permissions to read, write and create files and subdirectories in this directory at any time. This variable is typically not overridden. Usually RABBITMQ_MNESIA_DIR is overridden instead.**Default**:**Generic UNIX package**: $RABBITMQ_HOME/var/lib/rabbitmq/mnesia**Ubuntu and Debian** packages: /var/lib/rabbitmq/mnesia/**RPM**: /var/lib/rabbitmq/plugins**MacOS (Homebrew)**: ${install_prefix}/var/lib/rabbitmq/mnesia, the Homebrew prefix is usually /usr/local**Windows**: %APPDATA%\RabbitMQ |
| RABBITMQ_MNESIA_DIR                 | The directory where this RabbitMQ node's data is stored. This includes a schema database, message stores, cluster member information and other persistent node state.**Default**:**Generic UNIX package**: $RABBITMQ_MNESIA_BASE/$RABBITMQ_NODENAME**Ubuntu and Debian** packages: $RABBITMQ_MNESIA_BASE/$RABBITMQ_NODENAME**RPM**: $RABBITMQ_MNESIA_BASE/$RABBITMQ_NODENAME**MacOS (Homebrew)**: ${install_prefix}/var/lib/rabbitmq/mnesia/$RABBITMQ_NODENAME, the Homebrew prefix is usually /usr/local**Windows**: %APPDATA%\RabbitMQ\$RABBITMQ_NODENAME |
| RABBITMQ_PLUGINS_DIR                | The list of directories where [plugin](https://www.rabbitmq.com/plugins.html) archive files are located and extracted from. This is PATH-like variable, where different paths are separated by an OS-specific separator (: for Unix, ; for Windows). Plugins can be [installed](https://www.rabbitmq.com/plugins.html) to any of the directories listed here. Must not contain any characters mentioned in the [path restriction section](https://www.rabbitmq.com/configure.html#directory-and-path-restrictions).**Default**:**Generic UNIX package**: $RABBITMQ_HOME/plugins**Ubuntu and Debian** packages: /var/lib/rabbitmq/plugins**RPM**: /var/lib/rabbitmq/plugins**MacOS (Homebrew)**: ${install_prefix}/Cellar/rabbitmq/${version}/plugins, the Homebrew prefix is usually /usr/local**Windows**: %RABBITMQ_HOME%\plugins |
| RABBITMQ_PLUGINS_EXPAND_DIR         | The directory the node expand (unpack) [plugins](https://www.rabbitmq.com/plugins.html) to and use it as a code path location. Must not contain any characters mentioned in the [path restriction section](https://www.rabbitmq.com/configure.html#directory-and-path-restrictions).**Default**:**Generic UNIX package**: $RABBITMQ_MNESIA_BASE/$RABBITMQ_NODENAME-plugins-expand**Ubuntu and Debian** packages: $RABBITMQ_MNESIA_BASE/$RABBITMQ_NODENAME-plugins-expand**RPM**: $RABBITMQ_MNESIA_BASE/$RABBITMQ_NODENAME-plugins-expand**MacOS (Homebrew)**: ${install_prefix}/var/lib/rabbitmq/mnesia/$RABBITMQ_NODENAME-plugins-expand**Windows**: %APPDATA%\RabbitMQ\$RABBITMQ_NODENAME-plugins-expand |
| RABBITMQ_USE_LONGNAME               | When set to true this will cause RabbitMQ to use fully qualified names to identify nodes. This may prove useful in environments that use fully-qualified domain names or use IP addresses as hostnames or part of node names. Note that it is not possible to switch a node from short name to long name without resetting it.**Default**: false |
| RABBITMQ_SERVICENAME                | The name of the installed Windows service. This will appear in services.msc.**Default**: RabbitMQ. |
| RABBITMQ_CONSOLE_LOG                | Set this variable to new or reuse to redirect console output from the server to a file named %RABBITMQ_SERVICENAME% in the default RABBITMQ_BASE directory.If not set, console output from the server will be discarded (default).new: a new file will be created each time the service starts.reuse: the file will be overwritten each time the service starts.**Default**: (none) |
| RABBITMQ_SERVER_CODE_PATH           | Extra code path (a directory) to be specified when starting the runtime. Will be passed to the erl command when a node is started.**Default**: (none) |
| RABBITMQ_CTL_ERL_ARGS               | Parameters for the erl command used when invoking rabbitmqctl. This could be set to specify a range of ports to use for Erlang distribution: -kernel inet_dist_listen_min 35672 -kernel inet_dist_listen_max 35680**Default**: (none) |
| RABBITMQ_SERVER_ERL_ARGS            | Standard parameters for the erl command used when invoking the RabbitMQ Server. This should be overridden for debugging purposes only. Overriding this variable *replaces* the default value.**Default**:**Unix\*:** +P 1048576 +t 5000000 +stbt db +zdbbl 128000**Windows:** None |
| RABBITMQ_SERVER_ADDITIONAL_ERL_ARGS | Additional parameters for the erl command used when invoking the RabbitMQ Server. The value of this variable is appended to the default list of arguments (RABBITMQ_SERVER_ERL_ARGS).**Default**:**Unix\*:** None**Windows:** None |
| RABBITMQ_SERVER_START_ARGS          | Extra parameters for the erl command used when invoking the RabbitMQ Server. This will not override RABBITMQ_SERVER_ERL_ARGS.**Default**: (none) |

Besides the variables listed above, there are several environment variables which tell RabbitMQ [where to locate its database, log files, plugins, configuration and so on](https://www.rabbitmq.com/relocate.html).

Finally, some environment variables are operating system-specific.

| Name                        | Description                                                  |
| :-------------------------- | :----------------------------------------------------------- |
| HOSTNAME                    | The name of the current machine.**Default**:Unix, Linux: env hostnameMacOS: env hostname -s |
| COMPUTERNAME                | The name of the current machine.**Default**:Windows: localhost |
| ERLANG_SERVICE_MANAGER_PATH | This path is the location of erlsrv.exe, the Erlang service wrapper script.**Default**:Windows Service: %ERLANG_HOME%\erts-x.x.x\bin |

## [Operating System Kernel Limits](https://www.rabbitmq.com/configure.html#kernel-limits)

Most operating systems enforce limits on kernel resources: virtual memory, stack size, open file handles and more. To Linux users these limits can be known as "ulimit limits".

RabbitMQ nodes are most commonly affected by the maximum [open file handle limit](https://www.rabbitmq.com/networking.html#open-file-handle-limit). Default limit value on most Linux distributions is usually 1024, which is very low for a messaging broker (or generally, any data service). See [Production Checklist](https://www.rabbitmq.com/production-checklist.html) for recommended values.

### Modifying Limits

#### With systemd (Modern Linux Distributions)

On distributions that use systemd, the OS limits are controlled via a configuration file at /etc/systemd/system/rabbitmq-server.service.d/limits.conf. For example, to set the max open file handle limit (nofile) to 64000:

```plaintext
[Service]
LimitNOFILE=64000
```

See [systemd documentation](https://www.freedesktop.org/software/systemd/man/systemd.exec.html) to learn about the supported limits and other directives.

#### With Docker

To configure kernel limits for Docker contains, use the "default-ulimits" key in [Docker daemon configuration file](https://docs.docker.com/engine/reference/commandline/dockerd/#daemon-configuration-file). The file has to be installed on Docker hosts at /etc/docker/daemon.json:

```json
{
  "default-ulimits": {
    "nofile": {
      "Name": "nofile",
      "Hard": 64000,
      "Soft": 64000
    }
  }
}
```

#### Without systemd (Older Linux Distributions)

The most straightforward way to adjust the per-user limit for RabbitMQ on distributions that do not use systemd is to edit the /etc/default/rabbitmq-server (provided by the RabbitMQ Debian package) or [rabbitmq-env.conf](https://www.rabbitmq.com/configure.html#config-file) to invoke ulimit before the service is started.

```plaintext
ulimit -S -n 4096
```

This *soft* limit cannot go higher than the *hard* limit (which defaults to 4096 in many distributions). [The hard limit can be increased](https://github.com/basho/basho_docs/blob/master/content/riak/kv/2.2.3/using/performance/open-files-limit.md) via /etc/security/limits.conf. This also requires enabling the [pam_limits.so](http://askubuntu.com/a/34559) module and re-login or reboot.

Note that limits cannot be changed for running OS processes.

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!