# Memory Alarms

https://www.rabbitmq.com/memory.html

## Overview

This guide covers RabbitMQ memory threshold and paging settings, running nodes on 64-bit and 32-bit systems, and other related topics.  本指南涵盖 RabbitMQ 内存阈值和分页设置、在 64 位和 32 位系统上运行节点以及其他相关主题。

A separate guide, [Reasoning About Memory Use](https://www.rabbitmq.com/memory-use.html), covers how to determine what consumes memory on a running RabbitMQ node for the purpose of [monitoring](https://www.rabbitmq.com/monitoring.html) or troubleshooting.  一个单独的指南，关于内存使用的推理，介绍了如何确定正在运行的 RabbitMQ 节点上消耗内存的内容，以便进行监控或故障排除。

## Memory Threshold: What it is and How it Works  内存阈值：它是什么以及它是如何工作的

The RabbitMQ server detects the total amount of RAM installed in the computer on startup and when  RabbitMQ 服务器在启动时和何时检测计算机中安装的 RAM 总量

**rabbitmqctl set_vm_memory_high_watermark** *fraction* is executed. By default, when the RabbitMQ server uses above 40% of the available RAM, it raises a memory [alarm](https://www.rabbitmq.com/alarms.html) and blocks all connections that are publishing messages. Once the memory alarm has cleared (e.g. due to the server paging messages to disk or delivering them to clients that consume and [acknowledge the deliveries](https://www.rabbitmq.com/confirms.html)) normal service resumes.  rabbitmqctl set_vm_memory_high_watermark 分数被执行。默认情况下，当 RabbitMQ 服务器使用超过 40% 的可用 RAM 时，它会引发内存警报并阻止所有正在发布消息的连接。一旦内存警报清除（例如，由于服务器将消息分页到磁盘或将它们交付给消费并确认交付的客户端），正常服务就会恢复。

The default memory threshold is set to 40% of installed RAM. Note that this does not prevent the RabbitMQ server from using more than 40%, it is merely the point at which publishers are throttled. Erlang's garbage collector can, in the worst case, cause double the amount of memory to be used (by default, 80% of RAM). It is strongly recommended that OS swap or page files are enabled.  默认内存阈值设置为已安装 RAM 的 40%。请注意，这不会阻止 RabbitMQ 服务器使用超过 40%，它只是限制发布者的点。在最坏的情况下，Erlang 的垃圾收集器会导致使用双倍的内存量（默认情况下，80% 的 RAM）。强烈建议启用操作系统交换或页面文件。

32-bit architectures tend to impose a per process memory limit of 2GB. Common implementations of 64-bit architectures (i.e. AMD64 and Intel EM64T) permit only a paltry 256TB per process. 64-bit Windows further limits this to 8TB. However, note that even under 64-bit OSes, a 32-bit process frequently only has a maximum address space of 2GB.  32 位架构倾向于将每个进程的内存限制为 2GB。 64 位架构（即 AMD64 和 Intel EM64T）的常见实现每个进程只允许微不足道的 256TB。 64 位 Windows 进一步将其限制为 8TB。但是，请注意，即使在 64 位操作系统下，32 位进程通常也只有 2GB 的最大地址空间。

## Configuring the Memory Threshold  配置内存阈值

The memory threshold at which the flow control is triggered can be adjusted by editing the [configuration file](https://www.rabbitmq.com/configure.html#configuration-files).  可以通过编辑配置文件来调整触发流控的内存阈值。

The example below sets the threshold to the default value of 0.4:  下面的示例将阈值设置为默认值 0.4：

```bash
# new style config format, recommended
vm_memory_high_watermark.relative = 0.4
```

The default value of 0.4 stands for 40% of available (detected) RAM or 40% of available virtual address space, whichever is smaller. E.g. on a 32-bit platform with 4 GiB of RAM installed, 40% of 4 GiB is 1.6 GiB, but 32-bit Windows normally limits processes to 2 GiB, so the threshold is actually to 40% of 2 GiB (which is 820 MiB).  默认值 0.4 代表可用（检测到的）RAM 的 40% 或可用虚拟地址空间的 40%，以较小者为准。 例如。 在安装了 4 GiB RAM 的 32 位平台上，4 GiB 的 40% 是 1.6 GiB，但是 32 位 Windows 通常将进程限制为 2 GiB，因此阈值实际上是 2 GiB 的 40%（即 820 MiB ）。

Alternatively, the memory threshold can be adjusted by setting an absolute limit of RAM used by the node. The example below sets the threshold to 1073741824 bytes (1024 MiB):  或者，可以通过设置节点使用的 RAM 的绝对限制来调整内存阈值。 下面的示例将阈值设置为 1073741824 字节 (1024 MiB)：

```bash
vm_memory_high_watermark.absolute = 1073741824
```

Same example, but using memory units:

```bash
vm_memory_high_watermark.absolute = 1024MiB
```

If the absolute limit is larger than the installed RAM or available virtual address space, the threshold is set to whichever limit is smaller.  如果绝对限制大于安装的 RAM 或可用的虚拟地址空间，则将阈值设置为较小的限制。

The memory limit is appended to the [log file](https://www.rabbitmq.com/logging.html) when the RabbitMQ node starts:  RabbitMQ 节点启动时将内存限制附加到日志文件中：

```bash
2019-06-10 23:17:05.976 [info] <0.308.0> Memory high watermark set to 1024 MiB (1073741824 bytes) of 8192 MiB (8589934592 bytes) total
```

The memory limit may also be queried using the rabbitmq-diagnostics memory_breakdown and rabbitmq-diagnostics status commands.  也可以使用 rabbitmq-diagnostics memory_breakdown 和 rabbitmq-diagnostics status 命令来查询内存限制。

The threshold can be changed while the broker is running using the  可以在代理运行时使用

```bash
rabbitmqctl set_vm_memory_high_watermark <fraction>
```

command or

```bash
rabbitmqctl set_vm_memory_high_watermark absolute <memory_limit>
```

For example:

```bash
rabbitmqctl set_vm_memory_high_watermark 0.6
```

and

```bash
rabbitmqctl set_vm_memory_high_watermark absolute "4G"
```

When using the absolute mode, it is possible to use one of the following memory units:

- M, MiB for mebibytes (2^20 bytes)
- MB for megabytes (10^6 bytes)
- G, GiB for gibibytes (2^30 bytes)
- GB for gigabytes (10^9 bytes)

Both commands will have an effect until the node stops. To make the setting survive node restart, use the configuration setting instead.  在节点停止之前，这两个命令都会生效。要使设置在节点重新启动后仍然存在，请改用配置设置。

The memory limit may change on systems with hot-swappable RAM when this command is executed without altering the threshold, due to the fact that the total amount of system RAM is queried.  由于查询系统 RAM 总量这一事实，在不更改阈值的情况下执行此命令时，内存限制可能会在具有热插拔 RAM 的系统上发生变化。

### Disabling All Publishing  禁用所有发布

When the threshold or absolute limit is set to 0, it makes the memory alarm go off immediately and thus eventually blocks all publishing connections. This may be useful if you wish to disable publishing globally:  当阈值或绝对限制设置为 0 时，它会使内存警报立即响起，从而最终阻止所有发布连接。如果您希望全局禁用发布，这可能很有用：

```bash
rabbitmqctl set_vm_memory_high_watermark 0
```

## Limited Address Space

When running RabbitMQ inside a 32 bit Erlang VM in a 64 bit OS (or a 32 bit OS with PAE), the addressable memory is limited. The server will detect this and log a message like:  在 64 位操作系统（或带有 PAE 的 32 位操作系统）的 32 位 Erlang VM 中运行 RabbitMQ 时，可寻址内存是有限的。服务器将检测到这一点并记录如下消息：

```bash
2018-11-22 10:44:33.654 [warning] Only 2048MB of 12037MB memory usable due to limited address space.
```

The memory alarm system is not perfect. While stopping publishing will usually prevent any further memory from being used, it is quite possible for other things to continue to increase memory use. Normally when this happens and the physical memory is exhausted the OS will start to swap. But when running with a limited address space, running over the limit will cause the VM to terminate or killed by an out-of-memory mechanism of the operating system.  记忆警报系统并不完美。虽然停止发布通常会阻止使用任何进一步的内存，但其他事情很可能会继续增加内存使用量。通常，当这种情况发生并且物理内存耗尽时，操作系统将开始交换。但是当使用有限的地址空间运行时，超过限制运行会导致虚拟机被操作系统的内存不足机制终止或杀死。

It is therefore strongly recommended to run RabbitMQ on a 64 bit OS and a 64-bit [Erlang runtime](https://www.rabbitmq.com/which-erlang.html).  因此强烈建议在 64 位操作系统和 64 位 Erlang 运行时上运行 RabbitMQ。

## Configuring the Paging Threshold

Before the broker hits the high watermark and blocks publishers, it will attempt to free up memory by instructing queues to page their contents out to disc. Both persistent and transient messages will be paged out (the persistent messages will already be on disc but will be evicted from memory).  在代理达到高水位并阻止发布者之前，它将尝试通过指示队列将其内容分页到磁盘来释放内存。持久消息和临时消息都将被分页（持久消息将已经在磁盘上，但将从内存中逐出）。

By default this starts to happen when the broker is 50% of the way to the high watermark (i.e. with a default high watermark of 0.4, this is when 20% of memory is used). To change this value, modify the **vm_memory_high_watermark_paging_ratio** configuration from its default value of 0.5. For example:  默认情况下，当代理到达高水位线的 50% 时（即默认高水位线 0.4，这是使用 20% 的内存时），这开始发生。要更改此值，请将 vm_memory_high_watermark_paging_ratio 配置从其默认值 0.5 修改。例如：

```bash
vm_memory_high_watermark_paging_ratio = 0.75
vm_memory_high_watermark.relative = 0.4
```

The above configuration starts paging at 30% of memory used, and blocks publishers at 40%.  上面的配置在使用内存的 30% 时开始分页，并在 40% 时阻止发布者。

It is possible to set vm_memory_high_watermark_paging_ratio to a greater value than 1.0. In this case queues will not page their contents to disc. If this causes the memory alarm to go off, then producers will be blocked as explained above.  可以将 vm_memory_high_watermark_paging_ratio 设置为大于 1.0 的值。在这种情况下，队列不会将其内容分页到磁盘。如果这导致内存警报响起，那么生产者将被阻止，如上所述。

## Unrecognised Platforms

If the RabbitMQ server is unable to detect the operating system it is running on, it will append a warning to the [log file](https://www.rabbitmq.com/logging.html). It then assumes than 1GB of RAM is installed:  如果 RabbitMQ 服务器无法检测到它正在运行的操作系统，它将在日志文件中附加一个警告。然后假设安装了超过 1GB 的 RAM：

```bash
2018-11-22 10:44:33.654 [warning] Unknown total memory size for your OS {unix,magic_homegrown_os}. Assuming memory size is 1024MB.
```

In this case, the vm_memory_high_watermark configuration value is used to scale the assumed 1GB RAM. With the default value of vm_memory_high_watermark set to 0.4, RabbitMQ's memory threshold is set to 410MB, thus it will throttle producers whenever RabbitMQ is using more than 410MB memory. Thus when RabbitMQ can't recognize your platform, if you actually have 8GB RAM installed and you want RabbitMQ to throttle producers when the server is using above 3GB, set vm_memory_high_watermark to 3.  在这种情况下，vm_memory_high_watermark 配置值用于扩展假定的 1GB RAM。 vm_memory_high_watermark 的默认值设置为 0.4，RabbitMQ 的内存阈值设置为 410MB，因此当 RabbitMQ 使用超过 410MB 内存时，它会限制生产者。因此，当 RabbitMQ 无法识别您的平台时，如果您实际安装了 8GB RAM，并且您希望 RabbitMQ 在服务器使用超过 3GB 时限制生产者，请将 vm_memory_high_watermark 设置为 3。

For guidelines on recommended RAM watermark settings, see [Production Checklist](https://www.rabbitmq.com/production-checklist.html#resource-limits-ram).  有关推荐的 RAM 水印设置的指南，请参阅生产清单。

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!

