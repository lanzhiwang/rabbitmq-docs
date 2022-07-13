The following configuration settings can be set in the advanced config file only, under the **rabbit** section.

* msg_store_index_module

Implementation module for queue indexing. You are advised to read the [message store tuning](https://www.rabbitmq.com/persistence-conf.html) documentation before changing this.  队列索引的实现模块。 建议您在更改之前阅读消息存储调整文档。

Default: rabbit_msg_store_ets_index

```
{rabbit, [ 
{msg_store_index_module, rabbit_msg_store_ets_index} 
]}
```

* backing_queue_module

Implementation module for queue contents.

Default:
```
{rabbit, [ 
{backing_queue_module, rabbit_variable_queue} 
]}
```

* msg_store_file_size_limit

Message store segment file size. Changing this for a node with an existing (initialised) database is dangerous can lead to data loss!  消息存储段文件大小。 为具有现有（已初始化）数据库的节点更改此设置是危险的，可能会导致数据丢失！

Default: 16777216

```
{rabbit, [ 
%% Changing this for a node 
%% with an existing (initialised) database is dangerous can 
%% lead to data loss! 
{msg_store_file_size_limit, 16777216} 
]}
```

* trace_vhosts

Used internally by the [tracer](https://www.rabbitmq.com/firehose.html). You shouldn't change this.

Default:
```
{rabbit, [ 
{trace_vhosts, []} 
]}
```

* msg_store_credit_disc_bound

The credits that a queue process is given by the message store.  队列进程的信用由消息存储提供。

By default, a queue process is given 4000 message store credits, and then 800 for every 800 messages that it processes.  默认情况下，队列进程获得 4000 个消息存储信用，然后每处理 800 条消息获得 800 个信用。

Messages which need to be paged out due to memory pressure will also use this credit.  由于内存压力而需要分页的消息也将使用此信用。

The Message Store is the last component in the credit flow chain. [Learn about credit flow.](https://blog.rabbitmq.com/posts/2015/10/new-credit-flow-settings-on-rabbitmq-3-5-5/)  消息存储是信用流链中的最后一个组件。 了解信贷流动。

This value only takes effect when messages are persisted to the message store. If messages are embedded on the queue index, then modifying this setting has no effect because credit_flow is NOT used when writing to the queue index.  该值仅在消息持久化到消息存储时生效。 如果消息嵌入在队列索引中，则修改此设置无效，因为在写入队列索引时不使用 credit_flow。

Default:
```
{rabbit, [ 
{msg_store_credit_disc_bound, {4000, 800}} 
]}
```

* queue_index_max_journal_entries

After how many queue index journal entries it will be flushed to disk.  在多少个队列索引日志条目之后，它将被刷新到磁盘。

Default:
```
{rabbit, [ 
{queue_index_max_journal_entries, 32768} 
]}
```

* lazy_queue_explicit_gc_run_operation_threshold

Tunable value only for lazy queues when under memory pressure. This is the threshold at which the garbage collector and other memory reduction activities are triggered. A low value could reduce performance, and a high one can improve performance, but cause higher memory consumption. You almost certainly should not change this.  仅适用于在内存压力下的惰性队列的可调值。 这是触发垃圾收集器和其他内存减少活动的阈值。 较低的值可能会降低性能，较高的值可以提高性能，但会导致更高的内存消耗。 你几乎肯定不应该改变这一点。

Default:
```
{rabbit, [ 
{lazy_queue_explicit_gc_run_operation_threshold, 1000} 
]}
```

* queue_explicit_gc_run_operation_threshold

Tunable value only for normal queues when under memory pressure. This is the threshold at which the garbage collector and other memory reduction activities are triggered. A low value could reduce performance, and a high one can improve performance, but cause higher memory consumption. You almost certainly should not change this.  仅适用于内存压力下的正常队列的可调值。 这是触发垃圾收集器和其他内存减少活动的阈值。 较低的值可能会降低性能，较高的值可以提高性能，但会导致更高的内存消耗。 你几乎肯定不应该改变这一点。

Default:
```
{rabbit, [ 
{queue_explicit_gc_run_operation_threshold, 1000} 
]}
```



