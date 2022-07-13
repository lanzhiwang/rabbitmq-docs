## Logging

`PerfTest` depends transitively on SLF4J for logging (through RabbitMQ Java Client). `PerfTest` binary distribution ships with Logback as a SLF4J binding and uses Logback default configuration (printing logs to the console). If for any reason you need to use a specific Logback configuration file, you can do it this way:

```
bin/runjava -Dlogback.configurationFile=/path/to/logback.xml com.rabbitmq.perf.PerfTest
```

As of PerfTest 2.11.0, it is possible to define loggers directly from the command line. This is less powerful than using a configuration file, yet simpler to use and useful for quick debugging. Use the `rabbitmq.perftest.loggers` system property with `name=level` pairs, e.g.:

```
bin/runjava -Drabbitmq.perftest.loggers=com.rabbitmq.perf=debug com.rabbitmq.perf.PerfTest
```

It is possible to define several loggers by separating them with commas, e.g. `-Drabbitmq.perftest.loggers=com.rabbitmq.perf=debug,com.rabbitmq.perf.Producer=info`.

It is also possible to use an environment variable:

```
export RABBITMQ_PERF_TEST_LOGGERS=com.rabbitmq.perf=info
```

The system property takes precedence over the environment variable.

Use the environment variable with the Docker image:

```
docker run -it --rm --network perf-test \
  --env RABBITMQ_PERF_TEST_LOGGERS=com.rabbitmq.perf=debug,com.rabbitmq.perf.Producer=debug \
  pivotalrabbitmq/perf-test:latest --uri amqp://rabbitmq
```

If you use `PerfTest` as a standalone JAR in your project, please note it doesn't depend on any SLF4J binding, you can use your favorite one.

## 