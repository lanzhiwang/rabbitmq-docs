# shovel federation

```yaml
apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: rabbitmq-source
  namespace: operators
spec:
  replicas: 3
  service:
    type: NodePort
  image: rabbitmq:3.8.12-management
  persistence:
    storageClassName: local-path
    storage: 2Gi
  rabbitmq:
    additionalPlugins:
      - rabbitmq_shovel
      - rabbitmq_shovel_management
      - rabbitmq_federation
      - rabbitmq_federation_management


apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: rabbitmq-destination
  namespace: operators
spec:
  replicas: 3
  service:
    type: NodePort
  image: rabbitmq:3.8.12-management
  persistence:
    storageClassName: local-path
    storage: 2Gi
  rabbitmq:
    additionalPlugins:
      - rabbitmq_shovel
      - rabbitmq_shovel_management
      - rabbitmq_federation
      - rabbitmq_federation_management

```

## shovel

* https://cloud.tencent.com/developer/article/1469332

### URI examples

* **amqp://**
connect to local server as default user

* **amqp://user@/my-vhost**
connect to local server with alternate user and virtual host (passwords are not required for local connections)

* **amqp://server-name**
connect to server-name, without SSL and default credentials

* **amqp://user:password@server-name/my-vhost**
connect to server-name, with credentials and overridden virtual host

* **amqps://user:password@server-name?cacertfile=/path/to/cacert.pem&certfile=/path/to/cert.pem&keyfile=/path/to/key.pem&verify=verify_peer**
connect to server-name, with credentials and SSL

* **amqps://server-name?cacertfile=/path/to/cacert.pem&certfile=/path/to/cert.pem&keyfile=/path/to/key.pem&verify=verify_peer&fail_if_no_peer_cert=true&auth_mechanism=external**
connect to server-name, with SSL and EXTERNAL authentication

### Add a new shovel

1. **Name**：Shovel 配置的名称

2. Source：Source 中需要指定协议类型、连接的源节点地址，源端的类型(队列、交换器，如果是交换器还需要填入 routingKey)
    AMQP 0.9.1：
        **URI**
        **Queue(name)**、**Exchange(name，Routing key)**
        **Prefetch count**：该参数表示 Shovel 内部缓存的消息条数，Shovel 的内部缓存是源端服务器和目的端服务器之间的中间缓存部分
        **Auto-delete**：Never、After initial length transferred  默认为 Never 表示不删除自己，如果设置为 After initial length transferred，则在消息转移完成后删除自己

    AMQP 1.0：
        **URI**
        **Address**：The AMQP 1.0 address representing the source or target node.
        **Prefetch count**
        **Auto-delete**：Never、After initial length transferred

3. Destination：Destination 需要指定协议类型，连接的目的节点地址，目的端的类型(队列、交换器，如果是交换器还需要填入 routingKey)
    AMQP 0.9.1：
        **URI**
        **Queue(name)**、**Exchange(name，Routing key)**
        **Add forwarding headers**：yes、no  如果设置为 true，则会在转发的消息内添加 x-shovelled 的 header 属性

    AMQP 1.0：
        **URI**
        **Address**：The AMQP 1.0 address representing the source or target node.
        **Add forwarding headers**：yes、no

4. **Reconnect delay**：Time in seconds to wait after a shovel goes down before attempting reconnection. Defaults to 1 if not set.  指定在 Shovel link 失效的情况下，重新建立连接前需要等待的时间，单位为秒。如果设置为 0，则不会进行重连动作，即 Shovel 会在首次连接失效时停止工作

5. **Acknowledgement mode**：On confirm、On publish、On ack  no ack 表示无须任何消息确认行为；on publish 表示 Shovel 会把每一条消息发送到目的端之后再向源端发送消息确认； on confirm 表示 Shovel 会使用 publisher confirm 机制，在收到目的端的消息确认之后再向源端发送消息确认。

### exchange to exchange

![](../images/shovel_01.png)

### exchange to queue

![](../images/shovel_02.png)

### queue to exchange

![](../images/shovel_03.png)

### queue to queue

![](../images/shovel_04.png)





```bash
Shovel plugin:

   delete_shovel                 Deletes a Shovel
   restart_shovel                Restarts a dynamic Shovel
   shovel_status                 Displays status of Shovel on a node


rabbitmqctl [--node <node>] [--longnames] [--quiet] delete_shovel [--vhost <vhost>] <name>
rabbitmqctl [--node <node>] [--longnames] [--quiet] restart_shovel <name>
rabbitmqctl [--node <node>] [--longnames] [--quiet] shovel_status

$ ./rabbitmqctl -n rabbit@rabbitmq-source-server-0.rabbitmq-source-nodes.operators -l shovel_status
$ ./rabbitmqctl -n rabbit@rabbitmq-source-server-1.rabbitmq-source-nodes.operators -l delete_shovel

$ ./rabbitmqctl -n rabbit@rabbitmq-source-server-0.rabbitmq-source-nodes.operators -l list_parameters
$ ./rabbitmqctl -n rabbit@rabbitmq-source-server-0.rabbitmq-source-nodes.operators -l clear_parameter shovel test5

```





## federation

* http://linyishui.top/2020101001.html

### URI examples

* **amqp://server-name**
connect to server-name, without SSL and default credentials

* **amqp://user:password@server-name/my-vhost**
connect to server-name, with credentials and overridden virtual host

* **amqps://user:password@server-name?cacertfile=/path/to/cacert.pem&certfile=/path/to/cert.pem&keyfile=/path/to/key.pem&verify=verify_peer**
connect to server-name, with credentials and SSL

* **amqps://server-name?cacertfile=/path/to/cacert.pem&certfile=/path/to/cert.pem&keyfile=/path/to/key.pem&verify=verify_peer&fail_if_no_peer_cert=true&auth_mechanism=external**
connect to server-name, with SSL and EXTERNAL authentication

### Add a new upstream

1. General parameters

* **Name**
定义这个 upstream 的名称

* **URI**
URI to connect to. If upstream is a cluster and can have several URIs, you can enter them here separated by spaces.  要连接的 URI。 如果 upstream 是一个集群并且可以有多个 URI，您可以在此处输入它们，用空格分隔。

* **Prefetch count**
Maximum number of unacknowledged messages that may be in flight over a federation link at one time. Defaults to 1000 if not set.  定义 Federation 内部缓存的消息条数，即在收到上游消息之后且在发送到下游之前缓存的消息条数。

* **Reconnect delay**
Time in seconds to wait after a network link goes down before attempting reconnection. Defaults to 5 if not set.  Federation link 由于某种原因断开之后， 需要等待多少秒开始重新建立连接

* **Acknowledgement Mode**

	* on-confirm  Messages are acknowledged to the upstream broker after they have been confirmed downstream. Handles network errors and broker failures without losing messages. The slowest option, and the default.

	* on-publish  Messages are acknowledged to the upstream broker after they have been published downstream. Handles network errors without losing messages, but may lose messages in the event of broker failures.

	* no-ack
Message acknowledgements are not used. The fastest option, but may lose messages in the event of network or broker failures.

* **Trust User-ID**
Set "Yes" to preserve the "user-id" field across a federation link, even if the user-id does not match that used to republish the message. Set to "No" to clear the "user-id" field when messages are federated. Only set this to "Yes" if you trust the upstream broker not to forge user-ids.  设定 Federation 是否使用 Validated User-ID 这个功能。如果设置为 false 或者没有设置，那么 Federation 会忽略消息的 user_id 这个属性；如果设置为 true ，则 Federation 只会转发 user_id 为上游任意有效的用户的消息。


2. Federated exchanges parameters

* **Exchange**
The name of the upstream exchange. Default is to use the same name as the federated exchange.  上游交易所的名称。 默认是使用与联合交易所相同的名称。

* **Max hops**
Maximum number of federation links that messages can traverse before being dropped. Defaults to 1 if not set.  指定消息被丢弃前在 Federation link 中最大的跳转次数，默认为 1。注意即使设置 max-hops 参数为大于 1 的值，同一条消息也不会在同一个 Broker 中出现 2 次，但是有可能会在多个节点中被复制。

* **Expires**
Time in milliseconds that the upstream should remember about this node for. After this time all upstream state will be removed. Leave this blank to mean "forever".  上游应记住此节点的时间（以毫秒为单位）。 在此之后，所有上游状态都将被删除。 将此留空表示“永远”。

* **Message TTL**
Time in milliseconds that undelivered messages should be held upstream when there is a network outage or backlog. Leave this blank to mean "forever".  当网络中断或积压时，未传递的消息应保留在上游的时间（以毫秒为单位）。 将此留空表示“永远”。

* **HA Policy**
Determines the "x-ha-policy" argument for the upstream queue for a federated exchange. Default is "none", meaning the queue is not HA.  确定联合交换的上游队列的“x-ha-policy”参数。 默认为“none”，表示队列不是 HA。

3. Federated queues parameter

* **Queue**
The name of the upstream queue. Default is to use the same name as the federated queue.

* **Consumer tag**
The consumer tag to use when consuming from upstream. Optional.



```bash
Federation plugin:

   federation_status             Displays federation link status
   restart_federation_link       Restarts a running federation link

rabbitmqctl [--node <node>] [--longnames] [--quiet] federation_status [--only-down]
rabbitmqctl [--node <node>] [--longnames] [--quiet] restart_federation_link <link_id>


# 定义 upstream
rabbitmqctl [--node <node>] [--longnames] [--quiet] set_parameter [--vhost <vhost>] <component_name> <name> <value>

$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-0.rabbitmq-destination-nodes.operators -l list_parameters
Listing runtime parameters for vhost "/" ...
component	name	value
federation-upstream	exchange-exchange	{"ack-mode":"on-confirm","exchange":"exchangeA","trust-user-id":false,"uri":"amqp://P2942TJUXt7Mb87uv6RwR06o6uYZue-L:LcJonmjvC9eLcquVHDCSZn8oEngAV0nD@10.0.128.239:31062"}

# policy
rabbitmqctl [--node <node>] [--longnames] [--quiet] set_policy [--vhost <vhost>] [--priority <priority>] [--apply-to <apply-to>] <name> <pattern> <definition>

rabbitmqctl [--node <node>] [--longnames] [--quiet] list_policies [--vhost <vhost>] [--no-table-headers] [--timeout <timeout>]

$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-0.rabbitmq-destination-nodes.operators -l list_policies
Listing policies for vhost "/" ...
vhost	name	pattern	apply-to	definition	priority
/	exchange-exchange	exchangeA	exchanges	{"federation-upstream":"exchange-exchange"}	1


$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators -l federation_status
Listing federation links on node rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators...
[#{error => <<>>,exchange => <<"exchangeA">>,id => <<"061eeb8a">>,
   last_changed => <<"2021-09-10 09:34:36">>,
   local_connection =>
       <<"<rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators.1631263961.4072.0>">>,
   queue => <<>>,status => running,type => exchange,
   upstream => <<"exchange-exchange">>,upstream_exchange => <<"exchangeA">>,
   upstream_queue => <<>>,uri => <<"amqp://10.0.128.239:31062">>,
   vhost => <<"/">>}]
$


$ ./rabbitmqctl -n rabbit@rabbitmq-source-server-0.rabbitmq-source-nodes.operators -l list_bindings
Listing bindings for vhost /...
source_name	source_kind	destination_name	destination_kind	routing_key	arguments
	exchange	federation: exchangeA -> rabbitmq-destination	queue	federation: exchangeA -> rabbitmq-destination	[]
exchangeA	exchange	federation: exchangeA -> rabbitmq-destination B	exchange	test	[{"x-bound-from",[[{"cluster-name","rabbitmq-destination"},{"vhost","/"},{"exchange","/:exchangeA B"},{"hops",1}]]}]
federation: exchangeA -> rabbitmq-destination B	exchange	federation: exchangeA -> rabbitmq-destination	queue		[]
$

$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-0.rabbitmq-destination-nodes.operators -l list_bindings
Listing bindings for vhost /...
source_name	source_kind	destination_name	destination_kind	routing_key	arguments
exchangeA	exchange	queueA	queue	test	[]
test	exchange	test	queue	test	[]
$




$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-0.rabbitmq-destination-nodes.operators -l list_parameters
Listing runtime parameters for vhost "/" ...
component	name	value
federation-upstream	queue-queue	{"ack-mode":"on-confirm","queue":"queueC","trust-user-id":false,"uri":"amqp://P2942TJUXt7Mb87uv6RwR06o6uYZue-L:LcJonmjvC9eLcquVHDCSZn8oEngAV0nD@10.0.128.239:31062"}
federation-upstream	exchange-exchange	{"ack-mode":"on-confirm","exchange":"exchangeA","trust-user-id":false,"uri":"amqp://P2942TJUXt7Mb87uv6RwR06o6uYZue-L:LcJonmjvC9eLcquVHDCSZn8oEngAV0nD@10.0.128.239:31062"}
$

$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-0.rabbitmq-destination-nodes.operators -l list_policies
Listing policies for vhost "/" ...
vhost	name	pattern	apply-to	definition	priority
/	exchange-exchange	exchangeA	exchanges	{"federation-upstream":"exchange-exchange"}	1
/	queue-queue	queueB	queues	{"federation-upstream":"queue-queue"}	1
$

$ ./rabbitmqctl -n rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators -l federation_status
Listing federation links on node rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators...
[#{error => <<>>,exchange => <<"exchangeA">>,id => <<"061eeb8a">>,
   last_changed => <<"2021-09-10 09:34:36">>,
   local_connection =>
       <<"<rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators.1631263961.4072.0>">>,
   queue => <<>>,status => running,type => exchange,
   upstream => <<"exchange-exchange">>,upstream_exchange => <<"exchangeA">>,
   upstream_queue => <<>>,uri => <<"amqp://10.0.128.239:31062">>,
   vhost => <<"/">>},
 #{consumer_tag => <<"federation-link-queue-queue">>,error => <<>>,
   exchange => <<>>,id => <<"24e4bf82">>,
   last_changed => <<"2021-09-10 12:17:21">>,
   local_connection =>
       <<"<rabbit@rabbitmq-destination-server-2.rabbitmq-destination-nodes.operators.1631263961.15899.0>">>,
   queue => <<"queueB">>,status => running,type => queue,
   upstream => <<"queue-queue">>,upstream_exchange => <<>>,
   upstream_queue => <<"queueC">>,uri => <<"amqp://10.0.128.239:31062">>,
   vhost => <<"/">>}]
$



docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://PCMg5gyAXaktfRDNnqePBVY9V_VLgFnj:gcr2Fw_esu0pTt_l75z31S6FwItjFdEk@10.0.128.239:31218 -x0 -y1 -u "queueB" -a --id "test 1"


```



