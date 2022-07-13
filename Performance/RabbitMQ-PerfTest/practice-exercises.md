```bash

docker pull pivotalrabbitmq/perf-test
rabbitmq-perf-test-2.15.0-bin.tar.gz

bin/runjava com.rabbitmq.perf.PerfTest --help

docker run -it --rm pivotalrabbitmq/perf-test:latest --help
docker run -it --rm pivotalrabbitmq/perf-test:latest -v


apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: rabbitmq3816
  namespace: hz-rabbitmq
spec:
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchLabels:
              app.kubernetes.io/name: rabbitmq3816
          topologyKey: kubernetes.io/hostname
  image: rabbitmq:3.8.16-management
  override: {}
  persistence:
    storage: 5Gi
  rabbitmq:
    additionalPlugins:
      - rabbitmq_prometheu
    additionalConfig: |
      log.file.level = debug
  replicas: 3
  resources:
    limits:
      cpu: '2'
      memory: 2Gi
    requests:
      cpu: '2'
      memory: 2Gi
  service:
    type: NodePort
  terminationGracePeriodSeconds: 604800
  tls: {}




docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://2uNcpJiwNFvKmUDkTV6VMNpR3vufl_27:K54hO9UDhekI4C210Y7gbkZIpH845rVP@192.168.132.213:30267 -x 3 -y 3 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300





apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: sample-test
spec:
  persistence:
    storageClassName: local-path
    storage: 3Gi
  replicas: 3
  resources:
    limits:
      cpu: '4'
      memory: 8G
    requests:
      cpu: '4'
      memory: 8G
  service:
    type: NodePort
  terminationGracePeriodSeconds: 604800
  rabbitmq:
    additionalConfig: |
      vm_memory_high_watermark.relative=0.8
      vm_memory_high_watermark_paging_ratio=0.8
      log.file.level = debug



default_user = QltIBLxLgPgx5h7PyX_P3qghDTfhpD-f
default_pass = i7iz1Jz1o1YvdM0tomwmjuk7lag_b_aS
192.168.132.213
5672:31306/TCP,15672:32682/TCP,15692:31538/TCP

producer: 1 
consumer: 1
message size: 4 bytes
docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://QltIBLxLgPgx5h7PyX_P3qghDTfhpD-f:i7iz1Jz1o1YvdM0tomwmjuk7lag_b_aS@192.168.132.213:31306 -x 1 -y 1 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300

id: test 1, sending rate avg: 19659 msg/s
id: test 1, receiving rate avg: 19473 msg/s
[root@dataservice-master huzhi]#


producer: 1 
consumer: 0
message size: 4 bytes
docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -x 1 -y 0 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300

id: test 1, sending rate avg: 87979 msg/s
id: test 1, receiving rate avg: 0 msg/s
[root@dataservice-master huzhi]#


producer: 10 
consumer: 0
message size: 4 bytes
docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -x 10 -y 0 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300

test stopped (Reached time limit)
id: test 1, sending rate avg: 109179 msg/s
id: test 1, receiving rate avg: 0 msg/s
[root@dataservice-master huzhi]#


producer: 0 
consumer: 1
message size: 4 bytes

docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -y0 -p -u "throughput-test-1" -s 4 -C 1000000 -f persistent --id "test-1"

docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -x 0 -y 1 -u "throughput-test-1" --id "test 1" -s 4 -z 50 -f persistent

test stopped (Reached time limit)
id: test 1, sending rate avg: 0 msg/s
id: test 1, receiving rate avg: 11015 msg/s


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -x 0 -y 10 -u "throughput-test-1" --id "test 1" -s 4 -z 50 -f persistent

test stopped (Reached time limit)
id: test 1, sending rate avg: 0 msg/s
id: test 1, receiving rate avg: 18081 msg/s
[root@dataservice-master huzhi]#



找到瓶颈

docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -x 10 -y 0 -u "throughput-test-1" -a --id "test 1" -s 5000 -f persistent




消息大小 50k
cpu: '4'
memory: 8G
producer: 1 
consumer: 1

docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 -x 3 -y 3 -u "throughput-test-1" -a --id "test 1" -s 50000 --multi-ack-every 1000

docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 --producers 3 --consumers 3 --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 3 --size 50000 --autoack --time 30 --flag persistent

test stopped (Reached time limit)
id: test-094054-226, sending rate avg: 3197 msg/s
id: test-094054-226, receiving rate avg: 2182 msg/s


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 --producers 3 --consumers 3 --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 3 --size 50000 --autoack --time 30 --flag mandatory

test stopped (Reached time limit)
id: test-094158-312, sending rate avg: 5921 msg/s
id: test-094158-312, receiving rate avg: 4454 msg/s

docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 --producers 3 --consumers 3 --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 30 --size 50000 --autoack --time 30 --flag mandatory

test stopped (Reached time limit)
id: test-094514-652, sending rate avg: 6162 msg/s
id: test-094514-652, receiving rate avg: 4657 msg/s


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://BnhrYvsBQ2RJ-ZpkDgpDNDv1N4Ip4ldh:n82LEzFiXFFKkZoksNw5AwHyXg_WZyMG@10.0.132.141:31821 --producers 30 --consumers 30 --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 30 --size 50000 --autoack --time 30 --flag mandatory





M2EYgz9RLB868Z4gfYOI1eDA9UH9e8p5
NSntW1yhlWVsQ_gop6z6uqeNgRhJQceS
10.0.132.141
15672:31210/TCP,5672:31687/TCP



1.持久化性能测试
sh runjava com.rabbitmq.perf.PerfTest -h amqp:// USERNAME:PASSWORD@IPADDR:PORT -x4 -y4 -t"fanout" -u"testque" -k"kk01" -z 30 -s 50 -f persistent
2.非持久化性能测试
sh runjava com.rabbitmq.perf.PerfTest -h amqp://USERNAME:PASSWORD@IPADDR:PORT -x4 -y4 -t"fanout" -u"testque" -k"kk01" -z 30 -s 50


1. 持久化性能测试
docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://M2EYgz9RLB868Z4gfYOI1eDA9UH9e8p5:NSntW1yhlWVsQ_gop6z6uqeNgRhJQceS@10.0.132.141:31687 --producers 4 --consumers 4 -t fanout -k kk01 -u testque -z 30 -s 50 -f persistent

id: test-103616-830, sending rate avg: 13600 msg/s
id: test-103616-830, receiving rate avg: 11266 msg/s

id: test-104030-747, sending rate avg: 13661 msg/s
id: test-104030-747, receiving rate avg: 11333 msg/s

id: test-104148-605, sending rate avg: 13117 msg/s
id: test-104148-605, receiving rate avg: 11267 msg/s

id: test-104311-998, sending rate avg: 13280 msg/s
id: test-104311-998, receiving rate avg: 11164 msg/s

id: test-104439-074, sending rate avg: 13431 msg/s
id: test-104439-074, receiving rate avg: 11354 msg/s

2. 非持久化性能测试
docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://M2EYgz9RLB868Z4gfYOI1eDA9UH9e8p5:NSntW1yhlWVsQ_gop6z6uqeNgRhJQceS@10.0.132.141:31687 --producers 4 --consumers 4 -t fanout -k kk01 -u testque -z 30 -s 50

id: test-103808-835, sending rate avg: 26326 msg/s
id: test-103808-835, receiving rate avg: 25746 msg/s

id: test-104559-413, sending rate avg: 26182 msg/s
id: test-104559-413, receiving rate avg: 25488 msg/s

id: test-104709-902, sending rate avg: 25116 msg/s
id: test-104709-902, receiving rate avg: 24505 msg/s

id: test-104814-794, sending rate avg: 26225 msg/s
id: test-104814-794, receiving rate avg: 25492 msg/s

id: test-104916-557, sending rate avg: 26042 msg/s
id: test-104916-557, receiving rate avg: 25576 msg/s




docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://M2EYgz9RLB868Z4gfYOI1eDA9UH9e8p5:NSntW1yhlWVsQ_gop6z6uqeNgRhJQceS@10.0.132.141:31687 --producers 4 --consumers 4 -t fanout -k kk01 -u testque -z 30 -s 50 -f persistent -a

id: test-105046-473, sending rate avg: 43665 msg/s
id: test-105046-473, receiving rate avg: 43660 msg/s

id: test-105224-947, sending rate avg: 43808 msg/s
id: test-105224-947, receiving rate avg: 43807 msg/s

id: test-105315-433, sending rate avg: 42534 msg/s
id: test-105315-433, receiving rate avg: 42462 msg/s

id: test-105406-049, sending rate avg: 43025 msg/s
id: test-105406-049, receiving rate avg: 43024 msg/s

id: test-105504-889, sending rate avg: 42543 msg/s
id: test-105504-889, receiving rate avg: 42542 msg/s


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://M2EYgz9RLB868Z4gfYOI1eDA9UH9e8p5:NSntW1yhlWVsQ_gop6z6uqeNgRhJQceS@10.0.132.141:31687 --producers 4 --consumers 4 -t fanout -k kk01 -u testque -z 30 -s 50 -a

id: test-105645-132, sending rate avg: 50274 msg/s
id: test-105645-132, receiving rate avg: 50150 msg/s

id: test-105742-212, sending rate avg: 50281 msg/s
id: test-105742-212, receiving rate avg: 50210 msg/s

id: test-105834-151, sending rate avg: 48764 msg/s
id: test-105834-151, receiving rate avg: 48746 msg/s

id: test-105939-828, sending rate avg: 48605 msg/s
id: test-105939-828, receiving rate avg: 48447 msg/s

id: test-110043-038, sending rate avg: 50611 msg/s
id: test-110043-038, receiving rate avg: 50559 msg/s


bin/runjava -Drabbitmq.perftest.loggers=com.rabbitmq.perf=debug com.rabbitmq.perf.PerfTest --uri amqp://M2EYgz9RLB868Z4gfYOI1eDA9UH9e8p5:NSntW1yhlWVsQ_gop6z6uqeNgRhJQceS@10.0.132.141:31687 --producers 4 --consumers 4 -t fanout -k kk01 -u testque -z 30 -s 50 -a






###########################




docker run -d --name=netdata \
  -p 19999:19999 \
  -v netdataconfig:/etc/netdata \
  -v netdatalib:/var/lib/netdata \
  -v netdatacache:/var/cache/netdata \
  -v /etc/passwd:/host/etc/passwd:ro \
  -v /etc/group:/host/etc/group:ro \
  -v /proc:/host/proc:ro \
  -v /sys:/host/sys:ro \
  -v /etc/os-release:/host/etc/os-release:ro \
  --restart unless-stopped \
  --cap-add SYS_PTRACE \
  --security-opt apparmor=unconfined \
  netdata/netdata
  
http://localhost:19999, or http://NODE:19999.

http://10.0.131.89:19999


apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: sample-test
spec:
  persistence:
    storageClassName: top
    storage: 3Gi
  replicas: 3
  resources:
    limits:
      cpu: '4'
      memory: 8G
    requests:
      cpu: '4'
      memory: 8G
  service:
    type: NodePort
  terminationGracePeriodSeconds: 604800


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://tBgTflUr71x5cQsqySxvEYfTwIXoRHe8:QarhYvljFHBCmyi24WDJf4xSG7xt9iBN@10.0.131.89:31204 --producers 4 --consumers 4 -t fanout -k kk01 -u testque -z 30 -s 50




docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://tBgTflUr71x5cQsqySxvEYfTwIXoRHe8:QarhYvljFHBCmyi24WDJf4xSG7xt9iBN@10.0.131.89:31543 -x 1 -y 1 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300

tBgTflUr71x5cQsqySxvEYfTwIXoRHe8:QarhYvljFHBCmyi24WDJf4xSG7xt9iBN







docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://ERKbSoaeWNqbh3TK4RjfdwKw7bpu5p50:qqgGDdfAC3pELVirO0TLGomFzOmPGj84@10.0.128.224:30059 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300


docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://ERKbSoaeWNqbh3TK4RjfdwKw7bpu5p50:qqgGDdfAC3pELVirO0TLGomFzOmPGj84@10.0.128.224:30059 -x 30 -y 30 -a --id "test 1" -s 4 -z 300 --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 10


d6_UGYWJRBP1zxk9TQZZ_qCpHw_EESx5:H-nbnLXdH-HIQ-BnUQ99JBtf8V_9beZH



docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://d6_UGYWJRBP1zxk9TQZZ_qCpHw_EESx5:H-nbnLXdH-HIQ-BnUQ99JBtf8V_9beZH@10.0.131.89:30204 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300
...


LbTgLTv1EjV-pCXEPGm8--DRZmxSNo_e:ia5anbnXFKbeQZMbdphz0EwnJ-7Dxt4R


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://LbTgLTv1EjV-pCXEPGm8--DRZmxSNo_e:ia5anbnXFKbeQZMbdphz0EwnJ-7Dxt4R@10.0.131.89:30758 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300





docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://NDTG07DMaSMIGK7_AXcYm7LYuRWFlZgV:6CIEawEpI58y9cbWO1VcIixbfWzFqQOX@10.0.128.224:30059 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300
...


docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://NDTG07DMaSMIGK7_AXcYm7LYuRWFlZgV:6CIEawEpI58y9cbWO1VcIixbfWzFqQOX@10.0.128.224:30374 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 5000 -z 300


docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://6gVo10_U-_H_zHF0JeX1VPVJfggS2N4j:OFZCs2IuqWHUePXzTNprIxAnOFiE0mro@10.0.128.224:32368 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 5000 -z 300



docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://4vbWEVSVfqFTqIjubU6fn_SxTLGdRvie:D-py7tMRm7qwu9UkTjfr6KMp8XTq1a_P@10.0.129.181:31184 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 5000 -z 300
...




docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://GhKitp2WZ5QiXhj6VuR871IQVnQQKCeX:d5laKMwO7iX4lnx1QUnYhWjhEUNI1FrM@10.0.129.181:30558 -x 15 -y 15 -u "throughput-test-1" -a --id "test 1" -s 5000 -z 300
...



docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://iBxxQaCFXz2k2vKPENZp1CRYF_w41v30:L82Q2edBP4ExFC6l-nTnOdg6wOeztvI9@10.0.129.181:30786 -x 30 -y 1 -u "throughput-test-1" -a --id "test 1" -s 5000 -z 300
...
id: test 1, time: 286.144s, se




docker run -it --rm 10.0.130.200:60080/pivotalrabbitmq/perf-test:latest --uri amqp://5kCaNid7FBhr2ITumPbs6ch00nVLWMLd:-OF9FeCQCLjB1kgVzrThYLYP-X8yOO2B@10.0.129.181:30517 -x 10 -y 10 -u "throughput-test-1" -a --id "test 1" -s 4 -z 300
...


bin/runjava com.rabbitmq.perf.PerfTest \
--time 300 \
--queue-pattern 'perf-test-%d' \
--queue-pattern-from 1 \
--queue-pattern-to 1\
--producers 5 \
--consumers 150 \
--size 100 \
--autoack \
--flag persistent \
--uri amqp://fvQrym9IRe35quYnhYOGYMn1ydvgsfJV:PEMp7-n26Xv0pOAIQc0DY84PAi2WG_GV@10.4.82.62:5672?failover=failover_exchange


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://aTJ6TrMWlUBVqG_RV_X79h2Yfky-XaP0:S_TK9vQecglWlCJ_amUb8vxe1oCI2Of_@192.168.34.112:30428 -x 5 -y 150 -a --flag persistent --id "test 1" -s 100 -z 300 --queue-pattern 'perf-test-%d' --queue-pattern-from 1 --queue-pattern-to 1

 
 



docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://OTZnmSUv8RWSNo1TOjyOoRio_3ADyn38:sRBfs6X6KLVAk_YoIiXKjnZH7E-i3-m5@192.168.34.104:30699 --producers 1 --consumers 0 -u "perf-test" --size 500 --autoack --time 5 --flag persistent --message-properties message-ttl=60




docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://OTZnmSUv8RWSNo1TOjyOoRio_3ADyn38:sRBfs6X6KLVAk_YoIiXKjnZH7E-i3-m5@192.168.34.104:30699 --producers 1 --consumers 0 -u "perf-test" --size 500 --autoack --time 5 --flag persistent --message-properties expiration=60000


docker run -it --rm pivotalrabbitmq/perf-test:latest --uri amqp://OTZnmSUv8RWSNo1TOjyOoRio_3ADyn38:sRBfs6X6KLVAk_YoIiXKjnZH7E-i3-m5@192.168.34.104:30699 --producers 1 --consumers 0 -u "perf-test" --size 500 --autoack --time 5 --flag persistent --queue-args x-message-ttl=60000


```





|      |       |      |
| ---- | ----- | ---- |
| 1    | 16891 |      |
| 2    | 19146 |      |
| 3    | 21067 |      |
| 4    | 19599 |      |
| 5    | 19994 |      |
| 6    | 20765 |      |
| 7    | 20284 |      |
| 8    | 18910 |      |
| 9    | 20665 |      |
| 10   | 20274 |      |
|      |       |      |



2. 

