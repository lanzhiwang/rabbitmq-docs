



```bash

sample-lanzhiwang

apiVersion: rabbitmq.com/v1beta1
kind: RabbitmqCluster
metadata:  
  name: sample-lanzhiwang
  namespace: operators
spec:
  image: rabbitmq:3.8.12-management
  override: {}
  persistence:
    storage: 1Gi
  rabbitmq: {}
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
  tls: {}



[root@dataservice-master huzhi]# kubectl -n operators get sts sample-lanzhiwang-server -o yaml
apiVersion: apps/v1
kind: StatefulSet
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: sample-lanzhiwang
    app.kubernetes.io/part-of: rabbitmq
  name: sample-lanzhiwang-server
  namespace: operators
spec:
  podManagementPolicy: OrderedReady
  replicas: 3
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app.kubernetes.io/name: sample-lanzhiwang
  serviceName: sample-lanzhiwang-nodes
  template:
    metadata:
      annotations:
        prometheus.io/port: "15692"
        prometheus.io/scrape: "true"
      creationTimestamp: null
      labels:
        app.kubernetes.io/component: rabbitmq
        app.kubernetes.io/name: sample-lanzhiwang
        app.kubernetes.io/part-of: rabbitmq
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
          value: sample-lanzhiwang-nodes
        - name: RABBITMQ_USE_LONGNAME
          value: "true"
        - name: RABBITMQ_NODENAME
          value: rabbit@$(MY_POD_NAME).$(K8S_SERVICE_NAME).$(MY_POD_NAMESPACE)
        - name: K8S_HOSTNAME_SUFFIX
          value: .$(K8S_SERVICE_NAME).$(MY_POD_NAMESPACE)
        image: rabbitmq:3.8.12-management
        imagePullPolicy: IfNotPresent
        lifecycle:
          preStop:
            exec:
              command:
              - /bin/bash
              - -c
              - if [ ! -z "$(cat /etc/pod-info/skipPreStopChecks)" ]; then exit 0;
                fi; rabbitmq-upgrade await_online_quorum_plus_one -t 604800; rabbitmq-upgrade
                await_online_synchronized_mirror -t 604800; rabbitmq-upgrade drain
                -t 604800
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
            memory: 2000Mi
          requests:
            cpu: "1"
            memory: 2000Mi
        securityContext:
          privileged: true
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
      dnsPolicy: ClusterFirst
      initContainers:
      - command:
        - sh
        - -c
        - cp /tmp/erlang-cookie-secret/.erlang.cookie /var/lib/rabbitmq/.erlang.cookie
          && chown 999:999 /var/lib/rabbitmq/.erlang.cookie && chmod 600 /var/lib/rabbitmq/.erlang.cookie
          ; cp /tmp/rabbitmq-plugins/enabled_plugins /operator/enabled_plugins &&
          chown 999:999 /operator/enabled_plugins ; chown 999:999 /var/lib/rabbitmq/mnesia/
          ; echo '[default]' > /var/lib/rabbitmq/.rabbitmqadmin.conf && sed -e 's/default_user/username/'
          -e 's/default_pass/password/' /tmp/default_user.conf >> /var/lib/rabbitmq/.rabbitmqadmin.conf
          && chown 999:999 /var/lib/rabbitmq/.rabbitmqadmin.conf && chmod 600 /var/lib/rabbitmq/.rabbitmqadmin.conf
        image: rabbitmq:3.8.12-management
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
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext:
        fsGroup: 999
        runAsGroup: 999
        runAsUser: 999
      serviceAccount: sample-lanzhiwang-server
      serviceAccountName: sample-lanzhiwang-server
      terminationGracePeriodSeconds: 604800
      topologySpreadConstraints:
      - labelSelector:
          matchLabels:
            app.kubernetes.io/name: sample-lanzhiwang
        maxSkew: 1
        topologyKey: topology.kubernetes.io/zone
        whenUnsatisfiable: ScheduleAnyway
      volumes:
      - configMap:
          defaultMode: 420
          name: sample-lanzhiwang-plugins-conf
        name: plugins-conf
      - name: rabbitmq-confd
        projected:
          defaultMode: 420
          sources:
          - secret:
              items:
              - key: default_user.conf
                path: default_user.conf
              name: sample-lanzhiwang-default-user
          - configMap:
              items:
              - key: operatorDefaults.conf
                path: operatorDefaults.conf
              - key: userDefinedConfiguration.conf
                path: userDefinedConfiguration.conf
              name: sample-lanzhiwang-server-conf
      - emptyDir: {}
        name: rabbitmq-erlang-cookie
      - name: erlang-cookie-secret
        secret:
          defaultMode: 420
          secretName: sample-lanzhiwang-erlang-cookie
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
  updateStrategy:
    rollingUpdate:
      partition: 0
    type: RollingUpdate
  volumeClaimTemplates:
  - apiVersion: v1
    kind: PersistentVolumeClaim
    metadata:
      creationTimestamp: null
      labels:
        app.kubernetes.io/component: rabbitmq
        app.kubernetes.io/name: sample-lanzhiwang
        app.kubernetes.io/part-of: rabbitmq
      name: persistence
      namespace: operators
      ownerReferences:
      - apiVersion: rabbitmq.com/v1beta1
        blockOwnerDeletion: false
        controller: true
        kind: RabbitmqCluster
        name: sample-lanzhiwang
        uid: 6ccc9e27-e476-4bf8-881c-4f85e1837d12
    spec:
      accessModes:
      - ReadWriteOnce
      resources:
        requests:
          storage: 1Gi
      volumeMode: Filesystem
    status:
      phase: Pending



configMap: sample-lanzhiwang-plugins-conf
configMap: sample-lanzhiwang-server-conf

secret: sample-lanzhiwang-default-user
secret: sample-lanzhiwang-erlang-cookie


[root@dataservice-master huzhi]# kubectl -n operators get cm sample-lanzhiwang-plugins-conf -o yaml
apiVersion: v1
data:
  enabled_plugins: '[rabbitmq_peer_discovery_k8s,rabbitmq_prometheus,rabbitmq_management].'
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: sample-lanzhiwang
    app.kubernetes.io/part-of: rabbitmq
  name: sample-lanzhiwang-plugins-conf
  namespace: operators



[root@dataservice-master huzhi]# kubectl -n operators get cm sample-lanzhiwang-server-conf -o yaml
apiVersion: v1
data:
  operatorDefaults.conf: |
    cluster_formation.peer_discovery_backend = rabbit_peer_discovery_k8s
    cluster_formation.k8s.host               = kubernetes.default
    cluster_formation.k8s.address_type       = hostname
    cluster_partition_handling               = pause_minority
    queue_master_locator                     = min-masters
    disk_free_limit.absolute                 = 2GB
    cluster_name                             = sample-lanzhiwang
  userDefinedConfiguration.conf: |
    total_memory_available_override_value = 1677721600
kind: ConfigMap
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: sample-lanzhiwang
    app.kubernetes.io/part-of: rabbitmq
  name: sample-lanzhiwang-server-conf
  namespace: operators



[root@dataservice-master huzhi]# kubectl -n operators get secret sample-lanzhiwang-default-user -o yaml
apiVersion: v1
data:
  default_user.conf: ZGVmYXVsdF91c2VyID0gZHhfTEVIaFc3ZHUyOXNkbXBGdlByRWZha2E2UmJmblMKZGVmYXVsdF9wYXNzID0gM1lnd3pKdTB5OGFXVlhsSVUxd0I5TGthTnBVNEg3c2gK
  password: M1lnd3pKdTB5OGFXVlhsSVUxd0I5TGthTnBVNEg3c2g=
  username: ZHhfTEVIaFc3ZHUyOXNkbXBGdlByRWZha2E2UmJmblM=
kind: Secret
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: sample-lanzhiwang
    app.kubernetes.io/part-of: rabbitmq
  name: sample-lanzhiwang-default-user
  namespace: operators
type: Opaque



[root@dataservice-master huzhi]# kubectl -n operators get secret sample-lanzhiwang-erlang-cookie -o yaml
apiVersion: v1
data:
  .erlang.cookie: U0ZkVEljYXp5Q2lrMWk5am9PZnVwMmNBTm15X19udVk=
kind: Secret
metadata:
  labels:
    app.kubernetes.io/component: rabbitmq
    app.kubernetes.io/name: sample-lanzhiwang
    app.kubernetes.io/part-of: rabbitmq
  name: sample-lanzhiwang-erlang-cookie
  namespace: operators
type: Opaque




[root@dataservice-master huzhi]# kubectl -n operators get secret  sample-lanzhiwang-server-token-mq6wt -o yaml
apiVersion: v1
data:
  ca.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUM1ekNDQWMrZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRJeE1EY3hOakEwTVRnMU1Gb1hEVE14TURjeE5EQTBNVGcxTUZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTDBhCmF1MkdCNnFtcjU0ZlUyaG93eUZ3enNoankvbnNPMTViSnZhYVk2dHk2NUpRNGM4V2VyQVk4Tnl6cDZab0hreE4KbnpPcWFwSkZMUUNtak9KN3hIdHduOE9JYTdkWllNWXJZK0xwbERFbXIyV0Zyc0xQbzlQbEJ0MGlNVjl2dHUyZQpFSVk4NkR0SU1sTSsxTEJzRXMyeHJYV1Ywb3RKdWFMOGFNK0NuK2pKQWdOWVNELy9ndE5BNU44MG5Dc20ycmFRCk90VVdHU3g2Y2kxdUxHMDZMTEw2MVBDUXkzMmlrYXFaalp2dWd6L29jVTc4RlZLdTFZbnhrcEtYd04xcHdITVgKWHlzdC9wSVM0clV3TVFLUCtsME4ySlBFWERCRnoyY3gwMnpxQWFXaFFIM0NacUdLbEtucVVRQUl5YkszenZObQpvamxUZU9BKzUzS2t1eVJFWTlVQ0F3RUFBYU5DTUVBd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0hRWURWUjBPQkJZRUZCa1czVkltQnhLM0NaaG42RmdiQUVLNHBIb01NQTBHQ1NxR1NJYjMKRFFFQkN3VUFBNElCQVFDRmZHSlVTTGhnK0hCcmtydi91Y1NhZWhTYmZYb0Z1eW5pWFRqUnhPYyt1WHNVNUI2dQo3bHdoMUE4YTd4eWY3NmYzUDFMeUlJaGI0bWVSSU13NlpXY096UU1MQkdZUjJ6eXZnSWZjQUxTd1JESzQvcVZICmNPWUozRFRIQ05RdE4wZi8xUUEvVFlBRWZxMHRZOWx4aVoxRUlqMDM5R0xYcW1JYlZ2NGtSYVQvcFRLeUJwQjIKK0U5MFh6d0xwbVFGZkFxZmljVXVRcG82L0E5MDdpTFdBRUdOKzVmbXcvV2hkN1VSNzIzelhiMEJtSUdqNE5WdgpFMmM2M2tQRjJ6andsell5WDRTZGtKSGQ0N3JIbDJBeTdqVEpDbktyc3AyYkswWlV1bFRZNHBtWk5mc29oZHJzCnN5UzRhclorTDFjb0VwV2tQR3JtNE8rdTJ4K3o5U1RXZTk5awotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==
  namespace: b3BlcmF0b3Jz
  token: ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklrdDJZalZCVlVKNlF5MTVUbGd5T0dkMlNXTXRVRkZxWlZreE1XTTJhbVZPYlhCelVGWTVaV2wzY2pBaWZRLmV5SnBjM01pT2lKcmRXSmxjbTVsZEdWekwzTmxjblpwWTJWaFkyTnZkVzUwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXVZVzFsYzNCaFkyVWlPaUp2Y0dWeVlYUnZjbk1pTENKcmRXSmxjbTVsZEdWekxtbHZMM05sY25acFkyVmhZMk52ZFc1MEwzTmxZM0psZEM1dVlXMWxJam9pYzJGdGNHeGxMV3hoYm5wb2FYZGhibWN0YzJWeWRtVnlMWFJ2YTJWdUxXMXhObmQwSWl3aWEzVmlaWEp1WlhSbGN5NXBieTl6WlhKMmFXTmxZV05qYjNWdWRDOXpaWEoyYVdObExXRmpZMjkxYm5RdWJtRnRaU0k2SW5OaGJYQnNaUzFzWVc1NmFHbDNZVzVuTFhObGNuWmxjaUlzSW10MVltVnlibVYwWlhNdWFXOHZjMlZ5ZG1salpXRmpZMjkxYm5RdmMyVnlkbWxqWlMxaFkyTnZkVzUwTG5WcFpDSTZJakF5TkRnMlpETmxMVFJoT0RRdE5Ea3pNQzA0T0RnMExXWXlNR1l6TWpSa1l6ZzFOU0lzSW5OMVlpSTZJbk41YzNSbGJUcHpaWEoyYVdObFlXTmpiM1Z1ZERwdmNHVnlZWFJ2Y25NNmMyRnRjR3hsTFd4aGJucG9hWGRoYm1jdGMyVnlkbVZ5SW4wLlFZNEM2dDJrSDlfOS1WdWVNRmNSaFV2MDFjbURMUTF2aVlLNFprZDlLcnBPaFRHZXNJOGZvcFpwY3hFMlA3aTFYVnJaQXlDRnJBTWJjZXBGZlpXX1NOXzBvWGxFUFhTeU5SVktMN292ZGlYdXg2U2VyNjd6QlhmN09RMlQtUmg4TjdITVZXN0tlclo3WmY3b1hrbmloVF9iam9ZNmp5ekdUMC1fRDdfajk5UUNSR01jNU1XRUNZY3hQdjVadGl3TEdueGt5cm0xbFF2bnVYclJJV2NUTmszOG5pQnpKV19BM2xjaU9INEM0SllMZWNFbzR3dXhfSmtKT21QcGR2UjhlOW5TNGFVWWdxQV9la3FiMENBS29NR1R4clNTaGxPbXZibTV1eTVkRUgzLThRcGRHS1ROVndIWGJydUlTRFJZTnRvb0JDYnprYWdwbkZJdEJCUlFXZw==
kind: Secret
metadata:
  annotations:
    kubernetes.io/service-account.name: sample-lanzhiwang-server
    kubernetes.io/service-account.uid: 02486d3e-4a84-4930-8884-f20f324dc855
  name: sample-lanzhiwang-server-token-mq6wt
  namespace: operators
type: kubernetes.io/service-account-token






[root@dataservice-master huzhi]# kubectl -n operators get pvc
NAME                                     STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistence-sample-lanzhiwang-server-0   Bound    pvc-d37e1542-21bc-4dc1-be42-d04070534b2a   1Gi        RWO            topolvm        16m
persistence-sample-lanzhiwang-server-1   Bound    pvc-e3318cb2-d7d3-4832-936d-9edd3e21429a   1Gi        RWO            topolvm        15m
persistence-sample-lanzhiwang-server-2   Bound    pvc-71dfa785-b10b-407c-8d09-33deeadfed2f   1Gi        RWO            topolvm        15m


pvc-d37e1542-21bc-4dc1-be42-d04070534b2a pvc-e3318cb2-d7d3-4832-936d-9edd3e21429a pvc-71dfa785-b10b-407c-8d09-33deeadfed2f





[root@dataservice-master huzhi]# kubectl -n operators get pvc persistence-sample-lanzhiwang-server-0 -o yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: persistence-sample-lanzhiwang-server-0
  namespace: operators
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: topolvm
  volumeMode: Filesystem
  volumeName: pvc-d37e1542-21bc-4dc1-be42-d04070534b2a


[root@dataservice-master huzhi]# kubectl -n operators get pvc persistence-sample-lanzhiwang-server-1 -o yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: persistence-sample-lanzhiwang-server-1
  namespace: operators
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: topolvm
  volumeMode: Filesystem
  volumeName: pvc-e3318cb2-d7d3-4832-936d-9edd3e21429a


[root@dataservice-master huzhi]# kubectl -n operators get pvc persistence-sample-lanzhiwang-server-2 -o yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: persistence-sample-lanzhiwang-server-2
  namespace: operators
spec:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: topolvm
  volumeMode: Filesystem
  volumeName: pvc-71dfa785-b10b-407c-8d09-33deeadfed2f







[root@dataservice-master huzhi]# kubectl -n operators get services sample-lanzhiwang-test -o yaml
apiVersion: v1
kind: Service
metadata:
  name: sample-lanzhiwang
  namespace: operators
spec:
  externalTrafficPolicy: Cluster
  ports:
  - name: amqp
    nodePort: 30736
    port: 5672
    protocol: TCP
    targetPort: 5672
  - name: management
    nodePort: 30364
    port: 15672
    protocol: TCP
    targetPort: 15672
  selector:
    app.kubernetes.io/name: sample-lanzhiwang
  sessionAffinity: None
  type: NodePort




[root@dataservice-master huzhi]# kubectl -n operators get services sample-lanzhiwang-test-nodes -o yaml
apiVersion: v1
kind: Service
metadata:
  name: sample-lanzhiwang-nodes
  namespace: operators
spec:
  clusterIP: None
  ports:
  - name: epmd
    port: 4369
    protocol: TCP
    targetPort: 4369
  - name: cluster-rpc
    port: 25672
    protocol: TCP
    targetPort: 25672
  publishNotReadyAddresses: true
  selector:
    app.kubernetes.io/name: sample-lanzhiwang
  sessionAffinity: None
  type: ClusterIP










```







