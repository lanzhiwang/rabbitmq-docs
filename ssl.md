```bash
$ git clone https://github.com/michaelklishin/tls-gen.git

$ cd tls-gen/basic

# make SERVER_ALT_NAME=<rabbitmq集群名称>.<命名空间>.svc.cluster.local
$ make SERVER_ALT_NAME=rabbitmqcluster-sample.hz-rabbitmq.svc.cluster.local

$ tree -a result/
result/
├── ca_certificate.pem
├── ca_key.pem
├── client_certificate.pem
├── client_key.p12
├── client_key.pem
├── server_certificate.pem
├── server_key.p12
└── server_key.pem

$ kubectl -n hz-rabbitmq create secret tls rabbitmq-server-certs --cert=./result/server_certificate.pem --key=./result/server_key.pem

$ kubectl -n hz-rabbitmq create secret generic rabbitmq-ca-cert --from-file=ca.crt=./result/ca_certificate.pem

$ kubectl -n hz-rabbitmq get secret rabbitmq-server-certs -o yaml
apiVersion: v1
data:
  tls.crt: LS0tS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo=
  tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLSK
kind: Secret

$ kubectl -n hz-rabbitmq get secret rabbitmq-ca-cert -o yaml
apiVersion: v1
data:
  ca.crt: LS0tLS1VElGSUNBVEUtLS0tLQo=
kind: Secret


apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  name: rabbitmqcluster-sample
  namespace: hz-rabbitmq
spec:
  replicas: 3
  service:
    type: NodePort
  persistence:
    storageClassName: local-path
    storage: 3Gi
  resources:
    requests:
      cpu: 1000m
      memory: 2Gi
    limits:
      cpu: 1000m
      memory: 2Gi
  tls:
    secretName: rabbitmq-server-certs
    caSecretName: rabbitmq-ca-cert
    disableNonTLSListeners: false
  terminationGracePeriodSeconds: 60
  rabbitmq:
    additionalConfig: |
      ssl_options.fail_if_no_peer_cert = true

config file(s) :
/etc/rabbitmq/rabbitmq.conf
/etc/rabbitmq/conf.d/10-operatorDefaults.conf
/etc/rabbitmq/conf.d/11-default_user.conf
/etc/rabbitmq/conf.d/90-userDefinedConfiguration.conf

$ kubectl -n hz-rabbitmq get pods rabbitmqcluster-sample-server-0 -o yaml
apiVersion: v1
kind: Pod
metadata:
  annotations:
    kubernetes.io/psp: 20-user-restricted
    prometheus.io/port: "15691"
    prometheus.io/scrape: "true"
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: rabbitmqcluster-sample
    app.kubernetes.io/part-of: rabbitmq
  name: rabbitmqcluster-sample-server-0
  namespace: hz-rabbitmq
spec:
  automountServiceAccountToken: true
  containers:
  - env:
    - name: MY_POD_NAME
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.name
    - name: MY_POD_NAMESPACE
      valueFrom:
        fieldRef:
          apiVersion: v1
          fieldPath: metadata.namespace
    - name: RABBITMQ_ENABLED_PLUGINS_FILE
      value: /operator/enabled_plugins
    - name: K8S_SERVICE_NAME
      value: rabbitmqcluster-sample-nodes
    - name: RABBITMQ_USE_LONGNAME
      value: "true"
    - name: RABBITMQ_NODENAME
      value: rabbit@$(MY_POD_NAME).$(K8S_SERVICE_NAME).$(MY_POD_NAMESPACE)
    - name: K8S_HOSTNAME_SUFFIX
      value: .$(K8S_SERVICE_NAME).$(MY_POD_NAMESPACE)
    image: 192.168.131.207:60080/3rdparty/rabbitmq:3.8.16-management
    imagePullPolicy: IfNotPresent
    lifecycle:
      preStop:
        exec:
          command:
          - /bin/bash
          - -c
          - if [ ! -z "$(cat /etc/pod-info/skipPreStopChecks)" ]; then exit 0; fi;
            rabbitmq-upgrade await_online_quorum_plus_one -t 60; rabbitmq-upgrade
            await_online_synchronized_mirror -t 60; rabbitmq-upgrade drain -t 60
    name: rabbitmq
    ports:
    - containerPort: 4369
      name: epmd
      protocol: TCP
    - containerPort: 5672
      name: amqp
      protocol: TCP
    - containerPort: 15672
      name: management
      protocol: TCP
    - containerPort: 15692
      name: prometheus
      protocol: TCP
    - containerPort: 5671
      name: amqps
      protocol: TCP
    - containerPort: 15671
      name: management-tls
      protocol: TCP
    - containerPort: 15691
      name: prometheus-tls
      protocol: TCP
    readinessProbe:
      failureThreshold: 3
      initialDelaySeconds: 10
      periodSeconds: 10
      successThreshold: 1
      tcpSocket:
        port: amqp
      timeoutSeconds: 5
    resources:
      limits:
        cpu: "1"
        memory: 2Gi
      requests:
        cpu: 25m
        memory: 51Mi
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /var/lib/rabbitmq/
      name: rabbitmq-erlang-cookie
    - mountPath: /var/lib/rabbitmq/mnesia/
      name: persistence
    - mountPath: /operator
      name: rabbitmq-plugins
    - mountPath: /etc/rabbitmq/conf.d/10-operatorDefaults.conf
      name: rabbitmq-confd
      subPath: operatorDefaults.conf
    - mountPath: /etc/rabbitmq/conf.d/11-default_user.conf
      name: rabbitmq-confd
      subPath: default_user.conf
    - mountPath: /etc/rabbitmq/conf.d/90-userDefinedConfiguration.conf
      name: rabbitmq-confd
      subPath: userDefinedConfiguration.conf
    - mountPath: /etc/pod-info/
      name: pod-info
    - mountPath: /etc/rabbitmq-tls/
      name: rabbitmq-tls
      readOnly: true
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-ztztg
      readOnly: true
  dnsPolicy: ClusterFirst
  enableServiceLinks: true
  hostname: rabbitmqcluster-sample-server-0
  initContainers:
  - command:
    - sh
    - -c
    - cp /tmp/erlang-cookie-secret/.erlang.cookie /var/lib/rabbitmq/.erlang.cookie
      && chown 999:999 /var/lib/rabbitmq/.erlang.cookie && chmod 600 /var/lib/rabbitmq/.erlang.cookie
      ; cp /tmp/rabbitmq-plugins/enabled_plugins /operator/enabled_plugins && chown
      999:999 /operator/enabled_plugins ; chown 999:999 /var/lib/rabbitmq/mnesia/
      ; echo '[default]' > /var/lib/rabbitmq/.rabbitmqadmin.conf && sed -e 's/default_user/username/'
      -e 's/default_pass/password/' /tmp/default_user.conf >> /var/lib/rabbitmq/.rabbitmqadmin.conf
      && chown 999:999 /var/lib/rabbitmq/.rabbitmqadmin.conf && chmod 600 /var/lib/rabbitmq/.rabbitmqadmin.conf
    image: 192.168.131.207:60080/3rdparty/rabbitmq:3.8.16-management
    imagePullPolicy: IfNotPresent
    name: setup-container
    resources:
      limits:
        cpu: 100m
        memory: 500Mi
      requests:
        cpu: 100m
        memory: 500Mi
    securityContext:
      capabilities:
        drop:
        - FSETID
        - KILL
        - SETGID
        - SETUID
        - SETPCAP
        - NET_BIND_SERVICE
        - NET_RAW
        - SYS_CHROOT
        - MKNOD
        - AUDIT_WRITE
        - SETFCAP
      runAsUser: 0
    terminationMessagePath: /dev/termination-log
    terminationMessagePolicy: File
    volumeMounts:
    - mountPath: /tmp/rabbitmq-plugins/
      name: plugins-conf
    - mountPath: /var/lib/rabbitmq/
      name: rabbitmq-erlang-cookie
    - mountPath: /tmp/erlang-cookie-secret/
      name: erlang-cookie-secret
    - mountPath: /operator
      name: rabbitmq-plugins
    - mountPath: /var/lib/rabbitmq/mnesia/
      name: persistence
    - mountPath: /tmp/default_user.conf
      name: rabbitmq-confd
      subPath: default_user.conf
    - mountPath: /var/run/secrets/kubernetes.io/serviceaccount
      name: kube-api-access-ztztg
      readOnly: true
  nodeName: 192.168.131.208
  preemptionPolicy: PreemptLowerPriority
  priority: 0
  restartPolicy: Always
  schedulerName: default-scheduler
  securityContext:
    fsGroup: 999
    runAsGroup: 999
    runAsUser: 999
  serviceAccount: rabbitmqcluster-sample-server
  serviceAccountName: rabbitmqcluster-sample-server
  subdomain: rabbitmqcluster-sample-nodes
  terminationGracePeriodSeconds: 60
  tolerations:
  - effect: NoExecute
    key: node.kubernetes.io/not-ready
    operator: Exists
    tolerationSeconds: 30
  - effect: NoExecute
    key: node.kubernetes.io/unreachable
    operator: Exists
    tolerationSeconds: 30
  topologySpreadConstraints:
  - labelSelector:
      matchLabels:
        app.kubernetes.io/name: rabbitmqcluster-sample
    maxSkew: 1
    topologyKey: topology.kubernetes.io/zone
    whenUnsatisfiable: ScheduleAnyway
  volumes:
  - name: persistence
    persistentVolumeClaim:
      claimName: persistence-rabbitmqcluster-sample-server-0
  - configMap:
      defaultMode: 420
      name: rabbitmqcluster-sample-plugins-conf
    name: plugins-conf
  - name: rabbitmq-confd
    projected:
      defaultMode: 420
      sources:
      - secret:
          items:
          - key: default_user.conf
            path: default_user.conf
          name: rabbitmqcluster-sample-default-user
      - configMap:
          items:
          - key: operatorDefaults.conf
            path: operatorDefaults.conf
          - key: userDefinedConfiguration.conf
            path: userDefinedConfiguration.conf
          name: rabbitmqcluster-sample-server-conf
  - emptyDir: {}
    name: rabbitmq-erlang-cookie
  - name: erlang-cookie-secret
    secret:
      defaultMode: 420
      secretName: rabbitmqcluster-sample-erlang-cookie
  - emptyDir: {}
    name: rabbitmq-plugins
  - downwardAPI:
      defaultMode: 420
      items:
      - fieldRef:
          apiVersion: v1
          fieldPath: metadata.labels['skipPreStopChecks']
        path: skipPreStopChecks
    name: pod-info
  - name: rabbitmq-tls
    projected:
      defaultMode: 400
      sources:
      - secret:
          name: rabbitmq-server-certs
          optional: true
      - secret:
          name: rabbitmq-ca-cert
          optional: true
  - name: kube-api-access-ztztg
    projected:
      defaultMode: 420
      sources:
      - serviceAccountToken:
          expirationSeconds: 3607
          path: token
      - configMap:
          items:
          - key: ca.crt
            path: ca.crt
          name: kube-root-ca.crt
      - downwardAPI:
          items:
          - fieldRef:
              apiVersion: v1
              fieldPath: metadata.namespace
            path: namespace

$ kubectl -n hz-rabbitmq get cm rabbitmqcluster-sample-plugins-conf -o yaml
apiVersion: v1
data:
  enabled_plugins: '[rabbitmq_peer_discovery_k8s,rabbitmq_prometheus,rabbitmq_management].'
kind: ConfigMap

$ kubectl -n hz-rabbitmq get secret rabbitmqcluster-sample-default-user -o yaml
apiVersion: v1
data:
  default_user.conf: ZGVmYXVsdF91c2VyID0gSDJ5R0VUTXFrRnFwMVEwaHl2Ym9SMm41SDZpV3JDWWIKZGVmYXVsdF9wYXNzID0gR0U3S05KeTlfN0t6SXZPSERCM2NPWXY5bGpjSFZfdi0K
  password: R0U3S05KeTlfN0t6SXZPSERCM2NPWXY5bGpjSFZfdi0=
  provider: cmFiYml0bXE=
  type: cmFiYml0bXE=
  username: SDJ5R0VUTXFrRnFwMVEwaHl2Ym9SMm41SDZpV3JDWWI=
kind: Secret

$ kubectl -n hz-rabbitmq get cm rabbitmqcluster-sample-server-conf -o yaml
apiVersion: v1
data:
  operatorDefaults.conf: |
    cluster_formation.peer_discovery_backend             = rabbit_peer_discovery_k8s
    cluster_formation.k8s.host                           = kubernetes.default
    cluster_formation.k8s.address_type                   = hostname
    cluster_partition_handling                           = pause_minority
    queue_master_locator                                 = min-masters
    disk_free_limit.absolute                             = 2GB
    cluster_formation.randomized_startup_delay_range.min = 0
    cluster_formation.randomized_startup_delay_range.max = 60
    cluster_name                                         = rabbitmqcluster-sample
  userDefinedConfiguration.conf: |
    listeners.ssl.default                 = 5671
    ssl_options.cacertfile                = /etc/rabbitmq-tls/ca.crt
    ssl_options.certfile                  = /etc/rabbitmq-tls/tls.crt
    ssl_options.keyfile                   = /etc/rabbitmq-tls/tls.key
    ssl_options.verify                    = verify_peer
    ssl_options.fail_if_no_peer_cert      = true

    management.tcp.port                   = 15672
    management.ssl.port                   = 15671
    management.ssl.cacertfile             = /etc/rabbitmq-tls/ca.crt
    management.ssl.certfile               = /etc/rabbitmq-tls/tls.crt
    management.ssl.keyfile                = /etc/rabbitmq-tls/tls.key

    prometheus.tcp.port                   = 15692
    prometheus.ssl.port                   = 15691
    prometheus.ssl.cacertfile             = /etc/rabbitmq-tls/ca.crt
    prometheus.ssl.certfile               = /etc/rabbitmq-tls/tls.crt
    prometheus.ssl.keyfile                = /etc/rabbitmq-tls/tls.key
    total_memory_available_override_value = 1717986919
kind: ConfigMap

$ openssl s_client -connect 192.168.131.213:30570 -ssl3

$ openssl s_client -connect 192.168.131.213:30570 -tls1_2

$ keytool -importcert -alias server1 -file result/server_certificate.pem -keystore result/rabbitstore
Enter keystore password: jixjo567233
Re-enter new password: jixjo567233
Trust this certificate? [no]:  y
Certificate was added to keystore

$ keytool -importcert -alias client1 -file result/client_certificate.pem -keystore result/rabbitclient
Enter keystore password: jixjo567233
Re-enter new password: jixjo567233
Trust this certificate? [no]:  y
Certificate was added to keystore

$ ./runjava com.rabbitmq.perf.PerfTest --help

$ JAVA_OPTS="
-Djavax.net.ssl.trustStore=/root/huzhi/tls-gen/basic/result/rabbitstore
-Djavax.net.ssl.trustStorePassword=jixjo567233
-Djavax.net.ssl.keyStore=/root/huzhi/tls-gen/basic/result/client_key.p12
-Djavax.net.ssl.keyStoreType=PKCS12"
./runjava com.rabbitmq.perf.PerfTest -h amqps://H2yGETMqkFqp1Q0hyvboR2n5H6iWrCYb:GE7KNJy9_7KzIvOHDB3cOYv9ljcHV_v-@192.168.131.213:30570 -x 1 -y 2 -u "throughput-test-1" -a --id "test 1"















$ kubectl get crd | grep cert-manager
certificaterequests.cert-manager.io                              2021-10-08T07:53:03Z
certificates.cert-manager.io                                     2021-10-08T07:53:03Z
challenges.acme.cert-manager.io                                  2021-10-08T07:53:03Z
clusterissuers.cert-manager.io                                   2021-10-08T07:53:03Z
issuers.cert-manager.io                                          2021-10-08T07:53:04Z
orders.acme.cert-manager.io                                      2021-10-08T07:53:04Z



其中 Issuer 代表的是证书颁发者，可以定义各种提供者的证书颁发者，当前支持基于 Letsencrypt、vault 和 CA 的证书颁发者，还可以定义不同环境下的证书颁发者。

$ kubectl get issuers.cert-manager.io -A
NAMESPACE      NAME       READY   AGE
cpaas-system   cpaas-ca   True    24d

$ kubectl get issuers.cert-manager.io -n cpaas-system cpaas-ca -o yaml
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  annotations:
    helm.sh/chart-version: v3.6.24
    helm.sh/original-name: cpaas-ca
  creationTimestamp: "2021-10-08T07:53:12Z"
  generation: 1
  labels:
    helm.sh/chart-name: cert-manager
    helm.sh/release-name: cert-manager
    helm.sh/release-namespace: cpaas-system
  name: cpaas-ca
  namespace: cpaas-system
  resourceVersion: "3154"
  uid: 01783fe0-abfc-4b10-8b76-89d0deb55d53
spec:
  selfSigned: {}
status:
  conditions:
  - lastTransitionTime: "2021-10-08T07:53:28Z"
    observedGeneration: 1
    reason: IsReady
    status: "True"
    type: Ready
$




而 Certificates 代表的是生成证书的请求，一般其中存入生成证书的元信息，如域名等等。

一旦在 k8s 中定义了上述两类资源，部署的 cert-manager 则会根据 Issuer 和 Certificate 生成 TLS 证书，并将证书保存进 k8s 的 Secret 资源中，然后在 Ingress 资源中就可以引用到这些生成的 Secret 资源。对于已经生成的证书，还是定期检查证书的有效期，如即将超过有效期，还会自动续期。


$ kubectl get certificates.cert-manager.io -n chenjin-service redis-xxli-cert -o yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  creationTimestamp: "2021-10-18T07:45:52Z"
  generation: 1
  labels:
    managed-by: redis-cluster-operator
    redis.kun/name: redis-xxli
  name: redis-xxli-cert
  namespace: chenjin-service
  ownerReferences:
  - apiVersion: redis.kun/v1alpha1
    kind: DistributedRedisCluster
    name: redis-xxli
    uid: 2dc7bc2a-efc7-4acc-8e74-d3ee70033748
  resourceVersion: "10475129"
  uid: e1d85f82-5144-42dc-a7da-161ca8d9b263
spec:
  dnsNames:
  - redis-xxli.chenjin-service.svc
  - redis-xxli-proxy.chenjin-service.svc
  duration: 87600h0m0s
  issuerRef:
    kind: ClusterIssuer
    name: cpaas-ca
  secretName: redis-xxli-tls
status:
  conditions:
  - lastTransitionTime: "2021-10-18T07:45:53Z"
    message: Certificate is up to date and has not expired
    observedGeneration: 1
    reason: Ready
    status: "True"
    type: Ready
  notAfter: "2031-10-16T07:45:53Z"
  notBefore: "2021-10-18T07:45:53Z"
  renewalTime: "2028-06-16T15:45:53Z"
  revision: 1


$ kubectl get secret -n chenjin-service redis-xxli-tls -o yaml
apiVersion: v1
data:
  ca.crt: LS0tLS1
  tls.crt: LS0t
  tls.key: LS0tLS1
kind: Secret
metadata:
  annotations:
    cert-manager.io/alt-names: redis-xxli.chenjin-service.svc,redis-xxli-proxy.chenjin-service.svc
    cert-manager.io/certificate-name: redis-xxli-cert
    cert-manager.io/common-name: ""
    cert-manager.io/ip-sans: ""
    cert-manager.io/issuer-group: ""
    cert-manager.io/issuer-kind: ClusterIssuer
    cert-manager.io/issuer-name: cpaas-ca
    cert-manager.io/uri-sans: ""
  creationTimestamp: "2021-10-18T07:45:53Z"
  name: redis-xxli-tls
  namespace: chenjin-service
  resourceVersion: "10475123"
  uid: 5d3d9b2d-ca9b-429f-9bf7-1a9a4e517dfe
type: kubernetes.io/tls




$ kubectl get certificates.cert-manager.io -n operators redis-cluster-demo-cert -o yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  creationTimestamp: "2021-10-14T09:17:15Z"
  generation: 1
  labels:
    managed-by: redis-cluster-operator
    redis.kun/name: redis-cluster-demo
  name: redis-cluster-demo-cert
  namespace: operators
  ownerReferences:
  - apiVersion: redis.kun/v1alpha1
    kind: DistributedRedisCluster
    name: redis-cluster-demo
    uid: 9f59a480-f502-4d44-8339-ec50509c403b
  resourceVersion: "6222868"
  uid: a8fd99c0-a3f2-4454-a133-14ce5c43e7ba
spec:
  dnsNames:
  - redis-cluster-demo.operators.svc
  - redis-cluster-demo-proxy.operators.svc
  duration: 87600h0m0s
  issuerRef:
    kind: ClusterIssuer
    name: cpaas-ca
  secretName: redis-cluster-demo-tls
status:
  conditions:
  - lastTransitionTime: "2021-10-14T09:17:15Z"
    message: Certificate is up to date and has not expired
    observedGeneration: 1
    reason: Ready
    status: "True"
    type: Ready
  notAfter: "2031-10-12T09:17:15Z"
  notBefore: "2021-10-14T09:17:15Z"
  renewalTime: "2028-06-12T17:17:15Z"
  revision: 1


[root@dataservice-master ~]# kubectl get secret -n operators redis-cluster-demo-tls -o yaml
apiVersion: v1
data:
  ca.crt: LS0tL
  tls.crt: LS0tLS1C
  tls.key: LS0tL
kind: Secret
metadata:
  annotations:
    cert-manager.io/alt-names: redis-cluster-demo.operators.svc,redis-cluster-demo-proxy.operators.svc
    cert-manager.io/certificate-name: redis-cluster-demo-cert
    cert-manager.io/common-name: ""
    cert-manager.io/ip-sans: ""
    cert-manager.io/issuer-group: ""
    cert-manager.io/issuer-kind: ClusterIssuer
    cert-manager.io/issuer-name: cpaas-ca
    cert-manager.io/uri-sans: ""
  creationTimestamp: "2021-10-14T09:17:15Z"
  name: redis-cluster-demo-tls
  namespace: operators
  resourceVersion: "6222862"
  uid: 85bc4211-4f58-4d20-8dd6-8d440dddd317
type: kubernetes.io/tls
[root@dataservice-master ~]#
















apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: rabbitmqcluster-sample
    app.kubernetes.io/part-of: rabbitmq
  name: rabbitmq-server-certs
  namespace: hz-rabbitmq
  ownerReferences:
  - apiVersion: rabbitmq.com/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: RabbitmqCluster
    name: rabbitmqcluster-sample
    uid: 70aac8c4-1ee4-4465-bc19-0c8c15db9c6b
spec:
  dnsNames:
  - rabbitmqcluster-sample.hz-rabbitmq.svc.cluster.local
  duration: 87600h0m0s
  issuerRef:
    kind: ClusterIssuer
    name: cpaas-ca
  secretName: rabbitmq-server-certs

---

apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: rabbitmqcluster-sample
    app.kubernetes.io/part-of: rabbitmq
  name: rabbitmq-ca-cert
  namespace: hz-rabbitmq
  ownerReferences:
  - apiVersion: rabbitmq.com/v1beta1
    blockOwnerDeletion: true
    controller: true
    kind: RabbitmqCluster
    name: rabbitmqcluster-sample
    uid: 70aac8c4-1ee4-4465-bc19-0c8c15db9c6b
spec:
  dnsNames:
  - rabbitmqcluster-sample.hz-rabbitmq.svc.cluster.local
  duration: 87600h0m0s
  issuerRef:
    kind: ClusterIssuer
    name: cpaas-ca
  secretName: rabbitmq-ca-cert



















```

