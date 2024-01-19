```yaml
apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  labels:
    app: rabbitmq
  annotations:
    some: annotation
  name: rabbitmqcluster-sample
  namespace: hz-rabbitmq
spec:
  replicas: 3
  image: my-private-registry/rabbitmq:my-custom-tag
  imagePullSecrets:
  - name: some-secret
  service:
    type: ClusterIP | NodePort | LoadBalancer
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-internal: 0.0.0.0/0
  persistence:
    storageClassName: fast
    storage: 20Gi
  resources:
    requests:
      cpu: 1000m
      memory: 2Gi
    limits:
      cpu: 1000m
      memory: 2Gi
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: kubernetes.io/hostname
            operator: In
            values:
            - node-1
  tolerations:
    - key: "dedicated"
      operator: "Equal"
      value: "rabbitmq"
      effect: "NoSchedule"
  rabbitmq:
    additionalConfig: |
      channel_max = 1050
    advancedConfig: |
      [
          {ra, [
              {wal_data_dir, '/var/lib/rabbitmq/quorum-wal'}
          ]}
      ].
    envConfig: |
      RABBITMQ_DISTRIBUTION_BUFFER_SIZE=some_value
    additionalPlugins:
      - rabbitmq_top
      - rabbitmq_shovel
  tls:
    secretName: rabbitmq-server-certs
    caSecretName: rabbitmq-ca-cert
    disableNonTLSListeners: true
  skipPostDeploySteps: false
  terminationGracePeriodSeconds: 60
  override:
    service:
      spec:
        ports:
          - name: additional-port # adds an additional port on the service
            protocol: TCP
            port: 12345
    statefulSet:
      spec:
        template:
          spec:
            containers:
              - name: rabbitmq
                ports:
                  - containerPort: 12345 # opens an additional port on the rabbitmq server container
                    name: additional-port
                    protocol: TCP

---
apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:
  labels:
    prometheus.io/port: '15692'
    prometheus.io/scrape: 'true'
  name: upstream
  namespace: hz-rabbitmq
spec:
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        - labelSelector:
            matchLabels:
              app.kubernetes.io/name: upstream
          topologyKey: kubernetes.io/hostname
  image: 192.168.134.214:60080/middleware/test/rabbitmq3812-management:add-network-tool.2203151432
  override: {}
  persistence:
    storage: 3Gi
    storageClassName: local-path
  rabbitmq:
    additionalConfig: |-
      vm_memory_high_watermark.relative = 0.8
      vm_memory_high_watermark_paging_ratio = 0.9
      total_memory_available_override_value = 2GB
      disk_free_limit.relative = 3.0
      log.file.level = debug
      queue_index_embed_msgs_below = 10240
      tcp_listen_options.backlog = 128
      tcp_listen_options.nodelay = true
      tcp_listen_options.linger.on = true
      tcp_listen_options.linger.timeout = 0
      tcp_listen_options.sndbuf = 1966080
      tcp_listen_options.recbuf = 1966080
      channel_max = 160
    additionalPlugins:
      - rabbitmq_shovel
      - rabbitmq_shovel_management
      - rabbitmq_federation
      - rabbitmq_federation_management
      - rabbitmq_top
      - rabbitmq_event_exchange
      - rabbitmq_tracing
    envConfig: |
      RABBITMQ_DISTRIBUTION_BUFFER_SIZE=256000
    advancedConfig: |
      [
          {rabbit, [
            {msg_store_credit_disc_bound, {40000, 8000}},
            {credit_flow_default_credit, {2000, 500}},
            {queue_index_max_journal_entries, 32768}
          ]}
      ].
  replicas: 3
  resources:
    limits:
      cpu: '1'
      memory: 2Gi
    requests:
      cpu: '1'
      memory: 2Gi
  service:
    type: NodePort
  terminationGracePeriodSeconds: 604800
  tls: {}









```



