### rabbitmq 重启之后无法正常重新组建集群，报错如下：

```
2021-10-25 08:09:21.964 [warning] <0.273.0> Error while waiting for Mnesia tables: {timeout_waiting_for_tables,['rabbit@e2e-rabbitmq-server-2.e2e-rabbitmq-nodes.local-midautons','rabbit@e2e-rabbitmq-server-1.e2e-rabbitmq-nodes.local-midautons','rabbit@e2e-rabbitmq-server-0.e2e-rabbitmq-nodes.local-midautons'],[rabbit_durable_queue]}
2021-10-25 08:09:21.989 [info] <0.273.0> Waiting for Mnesia tables for 30000 ms, 7 retries left
```

这是因为在 mnesia 数据库中存储有集群节点信息，在重启 rabbitmq 时没有正确清理这些状态信息，并且现在在 k8s 上的 rabbitmq sts 的 pod 启动策略上按顺序启动，这导致先启动的 pod 会一直其他 pod启动完成，这样形成了死循环。目前的解决办法是强制删除集群状态然后重新启动。

```bash
kubectl -n local-midautons exec e2e-rabbitmq-server-0 -- rabbitmqctl force_boot
```

参考：
* https://stackoverflow.com/questions/60407082/rabbit-mq-error-while-waiting-for-mnesia-tables
* https://kb.vmware.com/s/article/78027







```bash
$ kubectl patch pv pvc-67b07d36-c2fe-41e2-83d0-6d0301d52f60 --type='json' -p='[{"op": "replace", "path": "/spec/persistentVolumeReclaimPolicy", "value":Retain}]'





$ kubectl -n hz-rabbitmq patch RabbitmqCluster rabbitmqcluster-sample --type='json' -p='[{"op": "replace", "path": "/spec/override/statefulSet/podManagementPolicy", "value":OrderedReady}]'

$ kubectl -n hz-rabbitmq patch RabbitmqCluster rabbitmqcluster-sample --type='json' -p='[{"op": "replace", "path": "/spec/override/statefulSet/podManagementPolicy", "value":Parallel}]'


$ kubectl patch pv pvc-67b07d36-c2fe-41e2-83d0-6d0301d52f60 --type='json' -p='[{"op": "remove", "path": "/spec/claimRef"}]'
 
$ kubectl patch pv pvc-3051eb33-a8ff-4c06-bf22-17b1c6feee7f --type='json' -p='[{"op": "remove", "path": "/spec/claimRef"}]'
 
$ kubectl patch pv pvc-4cf74a16-d474-4678-9880-018fab2ef425 --type='json' -p='[{"op": "remove", "path": "/spec/claimRef"}]'




[root@dataservice-master ~]# kubectl -n hz-rabbitmq get RabbitmqCluster -o yaml
apiVersion: v1
items:
- apiVersion: rabbitmq.com/v1beta1
  kind: RabbitmqCluster
  metadata:
    annotations:
      cpaas.io/creator: admin@cpaas.io
      cpaas.io/operator: admin@cpaas.io
      cpaas.io/updated-at: "2021-10-25T12:15:32Z"
    creationTimestamp: "2021-10-25T12:08:57Z"
    finalizers:
    - deletion.finalizers.rabbitmqclusters.rabbitmq.com
    generation: 2
    name: rabbitmqcluster-sample
    namespace: hz-rabbitmq
    resourceVersion: "18882780"
    uid: eade9c3e-bed9-4834-b37e-40b4e8cefaea
  spec:
    image: rabbitmq:3.8.16-management
    override:
      statefulSet: {}
    persistence:
      storage: 3Gi
      storageClassName: local-path
    rabbitmq: {}
    replicas: 3
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: "1"
        memory: 2Gi
    service:
      type: NodePort
    terminationGracePeriodSeconds: 60
    tls: {}
  status:
    binding:
      name: rabbitmqcluster-sample-default-user
    conditions:
    - lastTransitionTime: "2021-10-25T12:12:13Z"
      reason: AllPodsAreReady
      status: "True"
      type: AllReplicasReady
    - lastTransitionTime: "2021-10-25T12:12:12Z"
      reason: AtLeastOneEndpointAvailable
      status: "True"
      type: ClusterAvailable
    - lastTransitionTime: "2021-10-25T12:08:58Z"
      reason: NoWarnings
      status: "True"
      type: NoWarnings
    - lastTransitionTime: "2021-10-25T12:10:14Z"
      message: Finish reconciling
      reason: Success
      status: "True"
      type: ReconcileSuccess
    defaultUser:
      secretReference:
        keys:
          password: password
          username: username
        name: rabbitmqcluster-sample-default-user
        namespace: hz-rabbitmq
      serviceReference:
        name: rabbitmqcluster-sample
        namespace: hz-rabbitmq
    observedGeneration: 2
kind: List
metadata:
  resourceVersion: ""
  selfLink: ""
[root@dataservice-master ~]#






```



