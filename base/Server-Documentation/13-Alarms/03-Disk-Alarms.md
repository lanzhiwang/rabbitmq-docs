# Free Disk Space Alarms

https://www.rabbitmq.com/disk-alarms.html

## Overview

When free disk space drops below a configured limit (50 MB by default), an alarm will be triggered and all producers will be blocked.  当可用磁盘空间低于配置的限制（默认为 50 MB）时，将触发警报并阻止所有生产者。

The goal is to avoid filling up the entire disk which will lead all write operations on the node to fail and can lead to RabbitMQ termination.  目标是避免填满整个磁盘，这将导致节点上的所有写入操作失败并可能导致 RabbitMQ 终止。

## How it Works

To reduce the risk of filling up the disk, all incoming messages are blocked. Transient messages, which aren't normally persisted, are still paged out to disk when under memory pressure, and will use up the already limited disk space.  为了降低磁盘被填满的风险，所有传入的消息都会被阻止。通常不会持久化的瞬态消息在内存压力下仍会被分页到磁盘，并且会耗尽已经有限的磁盘空间。

If the disk alarm is set too low and messages are paged out rapidly, it is possible to run out of disk space and crash RabbitMQ in between disk space checks (at least 10 seconds apart). A more conservative approach would be to set the limit to the same as the amount of memory installed on the system (see the configuration [below](https://www.rabbitmq.com/disk-alarms.html#configure)).  如果磁盘警报设置得太低并且消息被快速分页，则可能会耗尽磁盘空间并在磁盘空间检查之间（至少相隔 10 秒）使 RabbitMQ 崩溃。更保守的方法是将限制设置为与系统上安装的内存量相同（请参阅下面的配置）。

An alarm will be triggered if the amount of free disk space drops below a configured limit.  如果可用磁盘空间量低于配置的限制，将触发警报。

The free space of the drive or partition that the broker database uses will be monitored at least every 10 seconds to determine whether the disk alarm should be raised or cleared.  将至少每 10 秒监视一次代理数据库使用的驱动器或分区的可用空间，以确定是否应发出或清除磁盘警报。

Monitoring will begin on node start. It will leave a [log entry](https://www.rabbitmq.com/logging.html) like this:  监控将在节点启动时开始。它将留下这样的日志条目：

```bash
2019-04-01 12:02:11.564 [info] <0.329.0> Enabling free disk space monitoring
2019-04-01 12:02:11.564 [info] <0.329.0> Disk free limit set to 950MB
```

Free disk space monitoring will be disabled on unrecognised platforms, causing an entry such as the one below:  可用磁盘空间监控将在无法识别的平台上禁用，从而导致如下条目：

```bash
2019-04-01 11:04:54.002 [info] <0.329.0> Disabling disk free space monitoring
```

When running RabbitMQ in a cluster, the disk alarm is cluster-wide; if one node goes under the limit then all nodes will block incoming messages.  在集群中运行 RabbitMQ 时，磁盘告警是集群范围的；如果一个节点低于限制，那么所有节点都将阻止传入消息。

RabbitMQ periodically checks the amount of free disk space. The frequency with which disk space is checked is related to the amount of space at the last check. This is in order to ensure that the disk alarm goes off in a timely manner when space is exhausted. Normally disk space is checked every 10 seconds, but as the limit is approached the frequency increases. When very near the limit RabbitMQ will check as frequently as 10 times per second. This may have some effect on system load.  RabbitMQ 定期检查可用磁盘空间量。检查磁盘空间的频率与上次检查时的空间量有关。这是为了确保在空间耗尽时磁盘警报及时响起。通常每 10 秒检查一次磁盘空间，但随着接近限制，频率会增加。当非常接近极限时，RabbitMQ 会以每秒 10 次的频率进行检查。这可能对系统负载有一些影响。

When free disk space drops below the configured limit, RabbitMQ will block producers and prevent memory-based messages from being paged to disk. This will reduce the likelihood of a crash due to disk space being exhausted, but will not eliminate it entirely. In particular, if messages are being paged out rapidly it is possible to run out of disk space and crash in the time between two runs of the disk space monitor. A more conservative approach would be to set the limit to the same as the amount of memory installed on the system (see the configuration section below).  当可用磁盘空间低于配置的限制时，RabbitMQ 将阻止生产者并阻止基于内存的消息被分页到磁盘。这将减少由于磁盘空间耗尽而导致崩溃的可能性，但不会完全消除它。特别是，如果消息被快速调出，则可能会在两次运行磁盘空间监视器之间的时间内耗尽磁盘空间并崩溃。更保守的方法是将限制设置为与系统上安装的内存量相同（请参阅下面的配置部分）。

## Configuring Disk Free Space Limit  配置磁盘可用空间限制

The disk free space limit is configured with the disk_free_limit setting. By default 50MB is required to be free on the database partition (see the description of [file locations](https://www.rabbitmq.com/relocate.html) for the default database location). This configuration file sets the disk free space limit to 1GB:  磁盘可用空间限制使用 disk_free_limit 设置进行配置。默认情况下，数据库分区需要 50MB 可用空间（有关默认数据库位置，请参阅文件位置说明）。此配置文件将磁盘可用空间限制设置为 1GB：

```bash
disk_free_limit.absolute = 1000000000
```

Or you can use memory units (KB, MB GB etc.) like this:

```bash
disk_free_limit.absolute = 1GB
```

It is also possible to set a free space limit relative to the RAM in the machine. This configuration file sets the disk free space limit to the same as the amount of RAM on the machine:  也可以设置相对于机器 RAM 的可用空间限制。 此配置文件将磁盘可用空间限制设置为与机器上的 RAM 量相同：

```bash
disk_free_limit.relative = 1.0
```

The limit can be changed while the broker is running using the rabbitmqctl set_disk_free_limit command or rabbitmqctl set_disk_free_limit mem_relative command. This command will take effect until next node restart.  可以在代理运行时使用 rabbitmqctl set_disk_free_limit 命令或 rabbitmqctl set_disk_free_limit mem_relative 命令更改限制。 该命令将一直生效，直到下一个节点重新启动。

The corresponding configuration setting should also be changed when the effects should survive a node restart.  当效果应该在节点重新启动后仍然存在时，也应该更改相应的配置设置。

## Getting Help and Providing Feedback

If you have questions about the contents of this guide or any other topic related to RabbitMQ, don't hesitate to ask them on the [RabbitMQ mailing list](https://groups.google.com/forum/#!forum/rabbitmq-users).

## Help Us Improve the Docs <3

If you'd like to contribute an improvement to the site, its source is [available on GitHub](https://github.com/rabbitmq/rabbitmq-website). Simply fork the repository and submit a pull request. Thank you!

