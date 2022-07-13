
```bash
$ kubectl -n hz-rabbitmq run rabbitmq-test -ti --image=192.168.134.214:60080/3rdparty/rabbitmq:3.8.16-management --rm=true --restart=Never bash

$ echo -n 'FQCbHuoNb8IpcU0Cpqz9pNwuYNQdM98Q' > /var/lib/rabbitmq/.erlang.cookie && chown root:root /var/lib/rabbitmq/.erlang.cookie && chmod 400 /var/lib/rabbitmq/.erlang.cookie && cd /opt/rabbitmq/sbin

$ ./rabbitmqctl cluster_status -n rabbit@destination-server-0.destination-nodes.hz-rabbitmq -l

$ ./rabbitmqctl list_users -n rabbit@sample-memory-8-server-0.sample-memory-8-nodes.operators -l


$ ./rabbitmqctl cluster_status -n rabbit@lanzhiwang-server-0.lanzhiwang-nodes.operators -l
$ ./rabbitmq-plugins -n rabbit@lanzhiwang-server-0.lanzhiwang-nodes.operators -l list

$ apt-get update && apt-get install curl
$ curl -s 10.3.2.193:15692/metrics




apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: my-rabbitmq
spec:
  persistence:
    storage: 1Gi
  replicas: 3
  resources:
    limits:
      cpu: '1'
      memory: 2000Mi
    requests:
      cpu: '1'
      memory: 2000Mi
  service:
    type: NodePort
  terminationGracePeriodSeconds: 604800
  rabbitmq:
    additionalPlugins:
      - rabbitmq_prometheus



prometheus.io/port: 15692
prometheus.io/scrape: "true"

[root@dataservice-master ~]# kubectl -n operators get PodMonitor rabbitmq -o yaml
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  labels:
    prometheus: kube-prometheus
  name: rabbitmq
  namespace: operators
spec:
  jobLabel: rabbitmq-monitor
  podMetricsEndpoints:
  - honorLabels: true
    interval: 15s
    port: prometheus
  selector:
    matchLabels:
      app.kubernetes.io/component: rabbitmq
[root@dataservice-master ~]#

$ cat rabbitmq-monitor.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: rabbitmq-monitor
  namespace: operators
  labels:
    prometheus: kube-prometheus
spec:
  endpoints:
  - interval: 15s
    port: prometheus
    honorLabels: true
  jobLabel: rabbitmq-monitor
  selector:
    matchLabels:
      app.kubernetes.io/component: rabbitmq






```







```bash
$ ls -al /opt/rabbitmq/sbin
total 44
drwxr-xr-x 2 rabbitmq rabbitmq  192 Sep  8  2021 .
drwxr-xr-x 7 rabbitmq rabbitmq 4096 Sep  8  2021 ..
-rwxr-xr-x 1 rabbitmq rabbitmq  983 Sep  8  2021 rabbitmq-defaults
-rwxr-xr-x 1 rabbitmq rabbitmq 1254 May  4  2021 rabbitmq-diagnostics  用于诊断和健康检查
-rwxr-xr-x 1 rabbitmq rabbitmq 7277 May  4  2021 rabbitmq-env
-rwxr-xr-x 1 rabbitmq rabbitmq 1250 May  4  2021 rabbitmq-plugins  用于插件管理
-rwxr-xr-x 1 rabbitmq rabbitmq 1249 May  4  2021 rabbitmq-queues  用于队列维护任务，特别是仲裁队列
-rwxr-xr-x 1 rabbitmq rabbitmq 6382 May  4  2021 rabbitmq-server  启动一个 RabbitMQ 服务器节点
-rwxr-xr-x 1 rabbitmq rabbitmq 1250 May  4  2021 rabbitmq-upgrade  用于与升级相关的维护任务
-rwxr-xr-x 1 rabbitmq rabbitmq 1245 May  4  2021 rabbitmqctl  用于服务管理和一般操作员任务
$
$ ls -al /usr/local/bin/rabbitmqadmin
-rwxr-xr-x 1 root root 42225 Sep  8  2021 /usr/local/bin/rabbitmqadmin
$

##############################################################################################


$ rabbitmqctl

Usage

rabbitmqctl [--node <node>] [--timeout <timeout>] [--longnames] [--quiet] <command> [<command options>]

Available commands:

###### Help:

autocomplete
Provides command name autocomplete variants  提供命令名称自动完成变体

help
Displays usage information for a command  显示命令的使用信息

version
Displays CLI tools version  显示 CLI 工具版本

###### Nodes:

await_startup
Waits for the RabbitMQ application to start on the target node  等待 RabbitMQ 应用程序在目标节点上启动

reset
Instructs a RabbitMQ node to leave the cluster and return to its virgin state  指示 RabbitMQ 节点离开集群并返回其原始状态

rotate_logs
Instructs the RabbitMQ node to perform internal log rotation  指示 RabbitMQ 节点执行内部日志轮换

shutdown
Stops RabbitMQ and its runtime (Erlang VM). Monitors progress for local nodes. Does not require a PID file path.  停止 RabbitMQ 及其运行时 (Erlang VM)。 监控本地节点的进度。 不需要 PID 文件路径。

start_app
Starts the RabbitMQ application but leaves the runtime (Erlang VM) running  启动 RabbitMQ 应用程序但保持运行时（Erlang VM）运行

stop
Stops RabbitMQ and its runtime (Erlang VM). Requires a local node pid file path to monitor progress.  停止 RabbitMQ 及其运行时（Erlang VM）。 需要本地节点 pid 文件路径来监控进度。

stop_app
Stops the RabbitMQ application, leaving the runtime (Erlang VM) running  停止 RabbitMQ 应用程序，让运行时（Erlang VM）继续运行

wait
Waits for RabbitMQ node startup by monitoring a local PID file. See also 'rabbitmqctl await_online_nodes'  通过监视本地 PID 文件等待 RabbitMQ 节点启动。 另见'rabbitmqctl await_online_nodes'

###### Cluster:

await_online_nodes
Waits for <count> nodes to join the cluster  等待 <count> 个节点加入集群

change_cluster_node_type
Changes the type of the cluster node  更改集群节点的类型

cluster_status
Displays all the nodes in the cluster grouped by node type, together with the currently running nodes  显示集群中按节点类型分组的所有节点，以及当前运行的节点

force_boot
Forces node to start even if it cannot contact or rejoin any of its previously known peers  强制节点启动，即使它无法联系或重新加入任何先前已知的对等点

force_reset
Forcefully returns a RabbitMQ node to its virgin state  强制将 RabbitMQ 节点返回到其原始状态

forget_cluster_node
Removes a node from the cluster  从集群中删除一个节点

join_cluster
Instructs the node to become a member of the cluster that the specified node is in  指示节点成为指定节点所在集群的成员

rename_cluster_node
Renames cluster nodes in the local database  重命名本地数据库中的集群节点

update_cluster_nodes
Instructs a cluster member node to sync the list of known cluster members from <seed_node>  指示集群成员节点同步来自 <seed_node> 的已知集群成员列表

###### Replication:

cancel_sync_queue
Instructs a synchronising mirrored queue to stop synchronising itself  指示同步镜像队列停止同步自身

sync_queue
Instructs a mirrored queue with unsynchronised mirrors (follower replicas) to synchronise them  指示具有未同步镜像（跟随者副本）的镜像队列同步它们

###### Users:

add_user
Creates a new user in the internal database. This user will have no permissions for any virtual hosts by default.  在内部数据库中创建一个新用户。 默认情况下，此用户将没有任何虚拟主机的权限。

authenticate_user
Attempts to authenticate a user. Exits with a non-zero code if authentication fails.  尝试对用户进行身份验证。 如果身份验证失败，则以非零代码退出。

change_password
Changes the user password

clear_password
Clears (resets) password and disables password login for a user

clear_user_limits
Clears user connection/channel limits

delete_user
Removes a user from the internal database. Has no effect on users provided by external backends such as LDAP

list_user_limits
Displays configured user limits

list_users
List user names and tags

set_user_limits
Sets user limits

set_user_tags
Sets user tags

###### Access Control:

clear_permissions
Revokes user permissions for a vhost

clear_topic_permissions
Clears user topic permissions for a vhost or exchange

list_permissions
Lists user permissions in a virtual host

list_topic_permissions
Lists topic permissions in a virtual host

list_user_permissions
Lists permissions of a user across all virtual hosts

list_user_topic_permissions
Lists user topic permissions

list_vhosts
Lists virtual hosts

set_permissions
Sets user permissions for a vhost

set_topic_permissions
Sets user topic permissions for an exchange

###### Monitoring, observability and health checks:

list_bindings
Lists all bindings on a vhost

list_channels
Lists all channels in the node

list_ciphers
Lists cipher suites supported by encoding commands

list_connections
Lists AMQP 0.9.1 connections for the node

list_consumers
Lists all consumers for a vhost

list_exchanges
Lists exchanges

list_hashes
Lists hash functions supported by encoding commands

list_node_auth_attempt_stats
Lists authentication attempts on the target node

list_queues
Lists queues and their properties

list_unresponsive_queues
Tests queues to respond within timeout. Lists those which did not respond

ping
Checks that the node OS process is up, registered with EPMD and CLI tools can authenticate with it

report
Generate a server status report containing a concatenation of all server status information for support purposes

schema_info
Lists schema database tables and their properties

status
Displays status of a node

###### Parameters:

clear_global_parameter
Clears a global runtime parameter

clear_parameter
Clears a runtime parameter.

list_global_parameters
Lists global runtime parameters

list_parameters
Lists runtime parameters for a virtual host

set_global_parameter
Sets a runtime parameter.

set_parameter
Sets a runtime parameter.

###### Policies:

clear_operator_policy
Clears an operator policy

clear_policy
Clears (removes) a policy

list_operator_policies
Lists operator policy overrides for a virtual host

list_policies
Lists all policies in a virtual host

set_operator_policy
Sets an operator policy that overrides a subset of arguments in user policies

set_policy
Sets or updates a policy

###### Virtual hosts:

add_vhost
Creates a virtual host

clear_vhost_limits
Clears virtual host limits

delete_vhost
Deletes a virtual host

list_vhost_limits
Displays configured virtual host limits

restart_vhost
Restarts a failed vhost data stores and queues

set_vhost_limits
Sets virtual host limits

trace_off

trace_on

###### Configuration and Environment:

decode
Decrypts an encrypted configuration value

encode
Encrypts a sensitive configuration value

environment
Displays the name and value of each variable in the application environment for each running application

set_cluster_name
Sets the cluster name

set_disk_free_limit
Sets the disk_free_limit setting

set_log_level
Sets log level in the running node

set_vm_memory_high_watermark
Sets the vm_memory_high_watermark setting

###### Definitions:

export_definitions
Exports definitions in JSON or compressed Erlang Term Format.

import_definitions
Imports definitions in JSON or compressed Erlang Term Format.

###### Feature flags:

enable_feature_flag
Enables a feature flag or all supported feature flags on the target node

list_feature_flags
Lists feature flags

###### Operations:

close_all_connections
Instructs the broker to close all connections for the specified vhost or entire RabbitMQ node

close_all_user_connections
Instructs the broker to close all connections of the specified user

close_connection
Instructs the broker to close the connection associated with the Erlang process id

eval
Evaluates a snippet of Erlang code on the target node

eval_file
Evaluates a file that contains a snippet of Erlang code on the target node

exec
Evaluates a snippet of Elixir code on the CLI node

force_gc
Makes all Erlang processes on the target node perform/schedule a full sweep garbage collection

resume_listeners
Resumes client connection listeners making them accept client connections again

suspend_listeners
Suspends client connection listeners so that no new client connections are accepted

###### Queues:

delete_queue
Deletes a queue

purge_queue
Purges a queue (removes all messages in it)

###### Shovel plugin:

delete_shovel
Deletes a Shovel

restart_shovel
Restarts a dynamic Shovel

shovel_status
Displays status of Shovel on a node

###### Deprecated:

hipe_compile
DEPRECATED. This command is a no-op. HiPE is no longer supported by modern Erlang versions

node_health_check
DEPRECATED. Performs intrusive, opinionated health checks on a fully booted node. See https://www.rabbitmq.com/monitoring.html#health-checks instead

Use 'rabbitmqctl help <command>' to learn more about a specific command
$

##############################################################################################


$ rabbitmq-diagnostics --help

Usage

rabbitmq-diagnostics [--node <node>] [--timeout <timeout>] [--longnames] [--quiet] <command> [<command options>]

Available commands:

###### Help:

autocomplete
Provides command name autocomplete variants

help
Displays usage information for a command

version
Displays CLI tools version

###### Nodes:

wait
Waits for RabbitMQ node startup by monitoring a local PID file. See also 'rabbitmqctl await_online_nodes'

###### Cluster:

cluster_status
Displays all the nodes in the cluster grouped by node type, together with the currently running nodes

###### Users:

list_user_limits
Displays configured user limits

list_users
List user names and tags

###### Access Control:

list_permissions
Lists user permissions in a virtual host

list_topic_permissions
Lists topic permissions in a virtual host

list_user_permissions
Lists permissions of a user across all virtual hosts

list_user_topic_permissions
Lists user topic permissions

list_vhosts
Lists virtual hosts

###### Monitoring, observability and health checks:

alarms
Lists resource alarms (local or cluster-wide) in effect on the target node  列出对目标节点有效的资源警报（本地或集群范围）

check_alarms
Health check that exits with a non-zero code if the target node reports any alarms, local or cluster-wide.  如果目标节点报告任何本地或集群范围的警报，则以非零代码退出的运行状况检查。

check_certificate_expiration
Checks the expiration date on the certificates for every listener configured to use TLS  检查配置为使用 TLS 的每个侦听器的证书到期日期

check_if_node_is_mirror_sync_critical
Health check that exits with a non-zero code if there are classic mirrored queues without online synchronised mirrors (queues that would potentially lose data if the target node is shut down)  如果存在没有在线同步镜像的经典镜像队列（目标节点关闭时可能丢失数据的队列），则以非零代码退出的运行状况检查

check_if_node_is_quorum_critical
Health check that exits with a non-zero code if there are queues with minimum online quorum (queues that would lose their quorum if the target node is shut down)  如果存在具有最小在线仲裁的队列（如果目标节点关闭，队列将失去其仲裁），则以非零代码退出的运行状况检查

check_local_alarms
Health check that exits with a non-zero code if the target node reports any local alarms  如果目标节点报告任何本地警报，则以非零代码退出的健康检查

check_port_connectivity
Basic TCP connectivity health check for each listener's port on the target node  目标节点上每个侦听器端口的基本 TCP 连接健康检查

check_port_listener
Health check that exits with a non-zero code if target node does not have an active listener for given port  如果目标节点没有给定端口的活动侦听器，则以非零代码退出的运行状况检查

check_protocol_listener
Health check that exits with a non-zero code if target node does not have an active listener for given protocol  如果目标节点没有给定协议的活动侦听器，则以非零代码退出的运行状况检查

check_running
Health check that exits with a non-zero code if the RabbitMQ app on the target node is not running  如果目标节点上的 RabbitMQ 应用程序未运行，则以非零代码退出的运行状况检查

check_virtual_hosts
Health check that checks if all vhosts are running in the target node

cipher_suites
Lists cipher suites enabled by default. To list all available cipher suites, add the --all argument.

consume_event_stream
Streams internal events from a running node. Output is jq-compatible.

discover_peers
Performs peer discovery and lists discovered nodes, if any

erlang_version
Displays Erlang/OTP version on the target node

is_booting
Checks if RabbitMQ is still booting on the target node

is_running
Checks if RabbitMQ is fully booted and running on the target node

list_bindings
Lists all bindings on a vhost

list_channels
Lists all channels in the node

list_ciphers
Lists cipher suites supported by encoding commands

list_connections
Lists AMQP 0.9.1 connections for the node

list_consumers
Lists all consumers for a vhost

list_exchanges
Lists exchanges

list_hashes
Lists hash functions supported by encoding commands

list_network_interfaces
Lists network interfaces (NICs) on the target node

list_node_auth_attempt_stats
Lists authentication attempts on the target node

list_queues
Lists queues and their properties

list_unresponsive_queues
Tests queues to respond within timeout. Lists those which did not respond

listeners
Lists active connection listeners (bound interface, port, protocol) on the target node

log_tail
Prints the last N lines of the log on the node

log_tail_stream
Streams logs from a running node for a period of time

maybe_stuck
Detects Erlang processes ("lightweight threads") potentially not making progress on the target node

memory_breakdown
Provides a memory usage breakdown on the target node.

observer
Starts a CLI observer interface on the target node

ping
Checks that the node OS process is up, registered with EPMD and CLI tools can authenticate with it

quorum_status
Displays quorum status of a quorum queue

remote_shell
Starts an interactive Erlang shell on the target node

report
Generate a server status report containing a concatenation of all server status information for support purposes

runtime_thread_stats
Provides a breakdown of runtime thread activity stats on the target node

schema_info
Lists schema database tables and their properties

server_version
Displays server version on the target node

status
Displays status of a node

tls_versions
Lists TLS versions supported (but not necessarily allowed) on the target node


###### Parameters:

list_global_parameters
Lists global runtime parameters

list_parameters
Lists runtime parameters for a virtual host

###### Policies:

list_operator_policies
Lists operator policy overrides for a virtual host

list_policies
Lists all policies in a virtual host

###### Virtual hosts:

list_vhost_limits
Displays configured virtual host limits

###### Configuration and Environment:

certificates
Displays certificates (public keeys) for every listener on target node that is configured to use TLS

command_line_arguments
Displays target node's command-line arguments and flags as reported by the runtime

disable_auth_attempt_source_tracking
Disables the tracking of peer IP address and username of authentication attempts

enable_auth_attempt_source_tracking
Enables the tracking of peer IP address and username of authentication attempts

environment
Displays the name and value of each variable in the application environment for each running application

erlang_cookie_hash
Displays a hash of the Erlang cookie (shared secret) used by the target node

erlang_cookie_sources
Display Erlang cookie source (e.g. $HOME/.erlang.cookie file) information useful for troubleshooting

log_location
Shows log file location(s) on target node

os_env
Lists RabbitMQ-specific environment variables set on target node

reset_node_auth_attempt_metrics
Resets auth attempt metrics on the target node

resolve_hostname
Resolves a hostname to a set of addresses. Takes Erlang's inetrc file into account.

resolver_info
Displays effective hostname resolver (inetrc) configuration on target node

###### Feature flags:

list_feature_flags
Lists feature flags

###### Operations:

reclaim_quorum_memory
Flushes quorum queue processes WAL, performs a full sweep GC on all of its local Erlang processes

###### Shovel plugin:

shovel_status
Displays status of Shovel on a node

###### Deprecated:

node_health_check
DEPRECATED. Performs intrusive, opinionated health checks on a fully booted node. See https://www.rabbitmq.com/monitoring.html#health-checks instead

Use 'rabbitmq-diagnostics help <command>' to learn more about a specific command
$



##############################################################################################

##############################################################################################

##############################################################################################

##############################################################################################

##############################################################################################

##############################################################################################



```

