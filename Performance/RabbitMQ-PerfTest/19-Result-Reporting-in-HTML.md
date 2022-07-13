# RabbitMQ Performance Tool

https://github.com/rabbitmq/rabbitmq-perf-test/blob/main/html/README.md

We have created a couple of tools to facilitate benchmarking RabbitMQ in different usage scenarios. One part of these tools is the `PerfTest` Java class, the other part is a couple of HTML/JS tools that will let you plot the results obtained from the benchmarks into nicely looking graphs.  我们创建了几个工具来促进在不同使用场景中对 RabbitMQ 进行基准测试。 这些工具的一部分是 PerfTest Java 类，另一部分是一些 HTML/JS 工具，它们可以让您将从基准测试获得的结果绘制成漂亮的图形。

The following blog posts show some examples of what can be done with this library:  以下博客文章显示了一些可以使用此库完成的操作的示例：

* [RabbitMQ Performance Measurements, part 1](https://www.rabbitmq.com/blog/2012/04/17/rabbitmq-performance-measurements-part-1/).

* [RabbitMQ Performance Measurements, part 2](https://www.rabbitmq.com/blog/2012/04/25/rabbitmq-performance-measurements-part-2/).

## Running benchmarks

Let's see how to run some benchmarks and then display the results in HTML using this tool.  让我们看看如何运行一些基准测试，然后使用此工具以 HTML 格式显示结果。

To run a benchmark we need to create a *benchmark specification file*, which is simply a JSON file like this one:  要运行基准测试，我们需要创建一个基准测试规范文件，它只是一个像这样的 JSON 文件：

```json
[
    {
        'name': 'consume',
        'type': 'simple',
        'params': [
            {'time-limit': 30, 'producer-count': 4, 'consumer-count': 2}
        ]
    }
]
```

Place this code in a file called `publish-consume-spec.js` and then go to the root folder of the binary distribution and run the following command to start the benchmark:  将此代码放在名为 publish-consume-spec.js 的文件中，然后转到二进制分发版的根文件夹并运行以下命令以启动基准测试：

```bash
bin/runjava com.rabbitmq.perf.PerfTestMulti publish-consume-spec.js publish-consume-result.js
```

This command will start a benchmark scenario where four producers will send messages to RabbitMQ over a period of thirty seconds. At the same time, two consumers will be consuming those messages.  此命令将启动一个基准场景，其中四个生产者将在 30 秒内向 RabbitMQ 发送消息。 同时，两个消费者将消费这些消息。

The results will be stored in the file `publish-consume-result.js` which we will now use to display a graph in our HTML page.  结果将存储在文件 publish-consume-result.js 中，我们现在将使用该文件在 HTML 页面中显示图表。

## Displaying benchmark results

Provided you have included our libraries (refer to the "Boilerplate HTML" section to know how to do that), the following HTML snippet will display the graph for the benchmark that we just ran:  如果您已包含我们的库（请参阅“样板 HTML”部分以了解如何执行此操作），以下 HTML 片段将显示我们刚刚运行的基准测试的图表：

```html
<div class="chart"
     data-type="time"
     data-latency="true"
     data-x-axis="time (s)"
     data-y-axis="rate (msg/s)"
     data-y-axis2="latency (μs)"
     data-scenario="consume"></div>
```

Here we use HTML's *data* attributes to tell the performance library how the graph should be displayed. We are telling it to load the `consume` scenario, showing time in seconds on the x-axis, the rate of messages per second on the y-axis and a second y-axis showing latency in microseconds; all of this displayed in a *time* kind of graph:  这里我们使用 HTML 的数据属性来告诉性能库应该如何显示图形。 我们告诉它加载消费场景，在 x 轴上以秒为单位显示时间，在 y 轴上显示每秒消息的速率，第二个 y 轴以微秒为单位显示延迟； 所有这些都显示在一种时间图表中：

![](../../images/perf_test_publish-consume-graph.png)

If instead of the CSS class `"chart"` we use the `"small-chart"` CSS class, then we can get a graph like the one below:

```html
<div class="small-chart"
     data-type="time"
     data-x-axis="time(s)"
     data-y-axis=""
     data-scenario="no-ack"></div>
```

![](../../images/perf_test_small_chart.png)

Finally, there's a type of graphs called `"summary"` that can show a summary of the whole benchmark. Here's the *HTML* for displaying them:

```html
<div class="summary"
     data-scenario="shard"></div>
```

And this is how they look like:

![](../../images/perf_test_summary.png)

## Types of graphs

We support several `types` of graphs, that you can specify using the `data-type` attribute:

- `time`: this graph can plot several variables on the y-axis while plotting the time on the x-axis. For example you could compare the send and receive rate over a period of time.  该图可以在 y 轴上绘制多个变量，同时在 x 轴上绘制时间。 例如，您可以比较一段时间内的发送和接收速率。

In the previous section we showed how to display these kind of graphs using HTML.

- `series`: will plot how changing a variable affects the results of the benchmark, for example, what's the difference in speed from sending small, medium and large messages?. This type of graph can show you that.  将绘制更改变量如何影响基准测试的结果，例如，发送小型、中型和大型消息的速度有何不同？。 这种类型的图表可以告诉你。

Here's an HTML example of a `series` graph:

```html
<div class="chart"
     data-type="series"
     data-scenario="message-sizes-and-producers"
     data-x-key="producerCount"
     data-x-axis="producers"
     data-y-axis="rate (msg/s)"
     data-plot-key="send-msg-rate"
     data-series-key="minMsgSize"></div>
```

- `x-y`: we can use this one to compare, for example, how message size affects the message rate per second. Refer to the second blogpost for an example of this kind of graph.  我们可以使用这个来比较，例如，消息大小如何影响每秒的消息速率。 有关此类图的示例，请参阅第二篇博文。

![](../../images/perf_test_1_1_sending_rates_msg_sizes.png)

Here's how to represent an `x-y` graph in HTML:

```html
<div class="chart"
     data-type="x-y"
     data-scenario="message-sizes-large"
     data-x-key="minMsgSize"
     data-plot-keys="send-msg-rate send-bytes-rate"
     data-x-axis="message size (bytes)"
     data-y-axis="rate (msg/s)"
     data-y-axis2="rate (bytes/s)"
     data-legend="ne"></div>
```

- `r-l`: This type of graph can help us compare the sending rate of messages vs. the latency. See scenario "1 -> 1 sending rate attempted vs latency" from the first blogpost for an example:  这种类型的图表可以帮助我们比较消息的发送速率与延迟。 有关示例，请参阅第一篇博文中的场景“1 -> 1 尝试发送速率与延迟”：

![](../../images/perf_test_1_1_sending_rates_latency.png)

Here how's to draw a `r-l` graph with HTML:

```html
<div class="chart"
     data-type="r-l"
     data-x-axis="rate attempted (msg/s)"
     data-y-axis="rate (msg/s)"
     data-scenario="rate-vs-latency"></div>
```

To see how all these benchmark specifications can be put together take a look at the `various-spec.js` file in the HTML examples directory, The `various-result.js` file in the same directory contains the results of the benchmark process run on a particular computer and `various.html` shows you how to display the results in an HTML page.  要了解如何将所有这些基准规范放在一起，请查看 HTML 示例目录中的 different-spec.js 文件，同一目录中的 different-result.js 文件包含在特定环境中运行的基准进程的结果。 computer 和 various.html 向您展示了如何在 HTML 页面中显示结果。

## Supported HTML attributes

We can use several HTML attributes to tell the library how to draw the chart. Here's the list of the ones we support.

- `data-file`: this specifies the file from where to load the benchmark results, for example `data-file="results-mini-2.7.1.js"`. This file will be loaded via AJAX. If you are loading the results on a local machine, you might need to serve this file via HTTP, since certain browsers refuse to perform the AJAX call otherwise.  这指定了从何处加载基准测试结果的文件，例如 data-file="results-mini-2.7.1.js"。 该文件将通过 AJAX 加载。 如果您在本地机器上加载结果，您可能需要通过 HTTP 提供此文件，因为某些浏览器拒绝执行 AJAX 调用。

- `data-scenario`: A results file can contain several scenarios. This attribute specifies which one to display in the graph.  一个结果文件可以包含多个场景。 该属性指定在图表中显示哪一个。

- `data-type`: The type of graph as explained above in "Types of Graphs".

- `data-mode`: Tells the library from where to get the message rate. Possible values are `send` or `recv`. If no value is specified, then the rate is the average of the send and receive rates added together.  告诉库从何处获取消息速率。 可能的值为发送或接收。 如果未指定任何值，则速率是发送和接收速率相加的平均值。

- `data-latency`: If we are creating a chart to display latency, then by specifying the `data-latency` as `true` the average latency will also be plotted alongside *send msg rate* and *receive msg rate*.  如果我们正在创建一个图表来显示延迟，那么通过将数据延迟指定为 true，平均延迟也将与发送 msg 速率和接收 msg 速率一起绘制。

- `data-x-axis`, `data-y-axis`, `data-y-axis2`: These attributes specify the label of the `x` and the `y` axes.

- `data-series-key`: If we want to specify from where which JSON key to pick our series data, then we can provide this attribute. For example: `data-series-key="minMsgSize"`.  如果我们想指定从哪里选择我们的系列数据的 JSON 键，那么我们可以提供这个属性。 例如：data-series-key="minMsgSize"。

- `data-x-key`: Same as the previous attributed, but for the x axis. Example: `data-x-key="minMsgSize"`.

## Boilerplate HTML

The file `../html/examples/sample.html` shows a full HTML page used to display some results. You should include the following Javascript Files:

```html
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "https://www.w3.org/TR/html4/loose.dtd">
<html>
 <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>RabbitMQ Performance</title>
    <link href="../perf.css" rel="stylesheet" type="text/css">
    <!--[if lte IE 8]><script language="javascript" type="text/javascript" src="../lib/excanvas.min.js"></script><![endif]-->
    <script language="javascript" type="text/javascript" src="../lib/jquery.min.js"></script>
    <script language="javascript" type="text/javascript" src="../lib/jquery.flot.min.js"></script>
    <script language="javascript" type="text/javascript" src="../perf.js"></script>
    <script language="javascript" type="text/javascript">
    $(document).ready(function() {
      var main_results;
        $.ajax({
            url: 'publish-consume-result.js',
            success: function(data) {
                render_graphs(JSON.parse(data));
            },
            fail: function() { alert('error loading publish-consume-result.js'); }
        });
    });
    </script>
 </head>
    <body>
    <h1>RabbitMQ Performance Example</h1>

    <h3>Consume</h3>
    <div class="chart"
      data-type="time"
      data-latency="true"
      data-x-axis="time (s)"
      data-y-axis="rate (msg/s)"
      data-y-axis2="latency (μs)"
      data-scenario="consume"></div>
  </body>
 </html>
```

Our `perf.js` library depends on the *jQuery* and *jQuery Flot* libraries for drawing graphs, and the *excanvas* library for supporting older browsers.

Once we load the libraries we can initialize our page with the following Javascript:

```html
<script language="javascript" type="text/javascript">
$(document).ready(function() {
  var main_results;
    $.ajax({
        url: 'publish-consume-result.js',
        success: function(data) {
            render_graphs(JSON.parse(data));
        },
        fail: function() { alert('error loading publish-consume-result.js'); }
    });
});
</script>
```

We can then load the file with the benchmark results and pass that to our `render_graphs` function, which will take care of the rest, provided we have defined the various `div`s where our graphs are going to be drawn.  然后，我们可以加载包含基准测试结果的文件，并将其传递给我们的 render_graphs 函数，该函数将负责其余部分，前提是我们已经定义了将要绘制图形的各种 div。

## Writing benchmark specifications

Benchmarks specifications should be written in JSON format. We can define an array containing one or more benchmark scenarios to run. For example:

```json
[
    {
        'name': 'no-ack-long',
        'type': 'simple',
        'interval': 10000,
        'params': [
            {'time-limit': 500}
        ]
    },
    {
        'name': 'headline-publish',
        'type': 'simple',
        'params': [
            {'time-limit': 30, 'producer-count': 10, 'consumer-count': 0}
        ]
    }
]
```

This JSON object specifies two scenarios `'no-ack-long'` and `'headline-publish'`, of the type `simple` and sets parameters, like `producer-count`, for the benchmarks.

There are three kind of benchmark scenarios:

- `simple`: runs a basic benchmark based on the parameters in the spec as seen in the example above.

- `rate-vs-latency`: compares message rate with latency.

- `varying`: can vary some variables during the benchmark, for example message size as shown in the following scenario snippet:  在基准测试期间可以改变一些变量，例如消息大小，如以下场景片段所示：

```json
{
    'name': 'message-sizes-small',
    'type': 'varying',
    'params': [
        {'time-limit': 30}
    ],
    'variables': [
        {
            'name': 'min-msg-size',
            'values': [0, 100, 200, 500, 1000, 2000, 5000]
        }
    ]
}
```

Note that `min-msg-size` gets converted to `minMsgSize`.

You can also set the AMQP URI. See the [URI Spec](https://www.rabbitmq.com/uri-spec.html). Default to `"amqp://localhost"` . For example:

```json
[
    {
        'name': 'consume',
        'type': 'simple',
        'uri': 'amqp://rabbitmq_uri',
        'params': [
            {'time-limit': 30, 'producer-count': 4, 'consumer-count': 2}
        ]
    }
]
```

### Supported scenario parameters

The following parameters can be specified for a scenario:

- exchange-type: exchange type to be used during the benchmark. Defaults to `'direct'`

- exchange-name: exchange name to be used during the benchmark. Defaults to whatever `exchangeType` was set to.

- queue-names: list of queue names to be used during the benchmark. Defaults to a single queue, letting RabbitMQ provide a random queue name.

- routing-key: routing key to be used during the benchmark. Defaults to an empty routing key.

- random-routing-key: allows the publisher to send a different routing key per published message. Useful when testing exchanges like the consistent hashing one. Defaults to `false`.

- producer-rate-limit: limit number of messages a producer will produce per second. Defaults to `0.0f`

- consumer-rate-limit: limit number of messages a consumer will consume per second. Defaults to 0.0f

- producer-count: number of producers to run for the benchmark. Defaults to 1

- consumer-count: number of consumers to run for the benchmark. Defaults to 1

- producer-tx-size: number of messages to send before committing the transaction. Defaults to 0, i.e.: no transactions

- consumer-tx-size: number of messages to consume before committing the transaction. Defaults to 0, i.e.: no transactions

- confirm: specifies whether to wait for publisher confirms during the benchmark. Defaults to -1. Any number >= 0 will make the benchmarks to use confirms.

- auto-ack: specifies whether the benchmarks should auto-ack messages. Defaults to `false`.

- multi-ack-every: specifies whether to send a multi-ack every X seconds. Defaults to `0`.

- channel-prefetch: sets the per-channel prefetch. Defaults to `0`.

- consumer-prefetch: sets the prefetch consumers. Defaults to `0`.

- min-msg-size: the size in bytes of the messages to be published. Defaults to `0`.

- time-limit: specifies how long the benchmark should be run. Defaults to`0`.

- producer-msg-count: number of messages to be published by the producers. Defaults to `0`.

- consumer-msg-count: number of messages to be consumed by the consumer. Defaults to `0`.

- msg-count: single flag to set the previous two counts to the same value.

- flags: flags to pass to the producer, like `"mandatory"`, or `"persistent"`. Defaults to an empty list.

- predeclared: tells the benchmark tool if the exchange/queue name provided already exists in the broker. Defaults to `false`.

## Starting a web server to display the results

Some browsers may need to use a web server (`file://` wouldn't work).

From the `html` directory, you can start a web server with Python:

$ python -m SimpleHTTPServer

As an alternative, from the root directory of the binary distribution, you can launch a Java-based web server:

```bash
bin/runjava com.rabbitmq.perf.WebServer
```

The latter command starts a web server listening on port 8080, with the `html` directory as its base directory. You can then see the included sample at http://localhost:8080/examples/sample.html. To change these defaults:

```bash
bin/runjava com.rabbitmq.perf.WebServer ./other-base-dir 9090
```

At last, if you want a quick preview of your results (same layout as the first 'consume' scenario above), ensure the scenario name is 'benchmark' in the result file and launch the following command:

```
$ bin/runjava com.rabbitmq.perf.BenchmarkResults my-result-file.js
```

The latter command will start a web server on port 8080 and open a browser window to display the results.



