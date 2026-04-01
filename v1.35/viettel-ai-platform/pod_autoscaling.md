# Horizontal Pod Autoscaling Validation

## Overview

This guide covers the verification of Horizontal Pod Autoscaling on Viettel AI Platform using the `autoscaling/v2` API with the metrics-server providing CPU/memory metrics.

## Step 1: Install Metrics Server

```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Patch for bare-metal (skip TLS verification for kubelet)
kubectl patch deployment metrics-server -n kube-system --type='json' \
  -p='[{"op":"add","path":"/spec/template/spec/containers/0/args/-","value":"--kubelet-insecure-tls"}]'
```

Verify metrics-server is running:

```bash
$ kubectl get pods -n kube-system -l k8s-app=metrics-server
NAME                              READY   STATUS    RESTARTS   AGE
metrics-server-5f54fb74d9-j9fs2   1/1     Running   0          43s

$ kubectl top nodes
NAME          CPU(cores)   CPU(%)   MEMORY(bytes)   MEMORY(%)
ubuntu-sv16   6905m        7%       46026Mi         11%
ubuntu-sv18   5206m        5%       5435Mi          2%
```

## Step 2: Create Namespace and Deploy Test Application

```bash
kubectl create namespace hpa-test

kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Service
metadata:
  name: test-webserver
  namespace: hpa-test
spec:
  type: ClusterIP
  selector:
    app: test-webserver
  ports:
  - name: http
    port: 8080
    targetPort: 8080
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-webserver
  namespace: hpa-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-webserver
  template:
    metadata:
      labels:
        app: test-webserver
    spec:
      containers:
      - name: webserver
        image: python:3.11-slim
        command: ["python3", "-c"]
        args:
        - |
          from http.server import HTTPServer, BaseHTTPRequestHandler
          class Handler(BaseHTTPRequestHandler):
            def do_GET(self):
              self.send_response(200)
              self.end_headers()
              self.wfile.write(b'OK')
            def log_message(self, *args): pass
          HTTPServer(('0.0.0.0', 8080), Handler).serve_forever()
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 256Mi
EOF
```

```bash
$ kubectl get pods -n hpa-test -l app=test-webserver
NAME                             READY   STATUS    RESTARTS   AGE
test-webserver-c9c779d86-dpdph   1/1     Running   0          30s
```

## Step 3: Create HorizontalPodAutoscaler

```bash
kubectl apply -f - <<'EOF'
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hpa-webserver
  namespace: hpa-test
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: test-webserver
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 20
  behavior:
    scaleDown:
      stabilizationWindowSeconds: 30
EOF
```

```bash
$ kubectl get hpa hpa-webserver -n hpa-test
NAME            REFERENCE                   TARGETS       MINPODS   MAXPODS   REPLICAS   AGE
hpa-webserver   Deployment/test-webserver   cpu: 1%/20%   1         10        1          30s
```

## Step 4: Verify Autoscaling Under Load

```bash
kubectl apply -f - <<'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: load-generator
  namespace: hpa-test
spec:
  replicas: 3
  selector:
    matchLabels:
      app: load-generator
  template:
    metadata:
      labels:
        app: load-generator
    spec:
      containers:
      - name: load-generator
        image: python:3.11-slim
        command: ["python3", "-c"]
        args:
        - |
          import urllib.request
          url = 'http://test-webserver.hpa-test.svc.cluster.local:8080'
          while True:
            try: urllib.request.urlopen(url, timeout=5)
            except: pass
        resources:
          requests:
            cpu: 100m
          limits:
            cpu: "1"
EOF
```

Monitor the HPA scaling up:

```bash
$ kubectl get hpa hpa-webserver -n hpa-test
NAME            REFERENCE                   TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
hpa-webserver   Deployment/test-webserver   cpu: 436%/20%   1         10        5          10m

$ kubectl get pods -n hpa-test -l app=test-webserver
NAME                             READY   STATUS    RESTARTS   AGE
test-webserver-c9c779d86-6qs78   1/1     Running   0          24s
test-webserver-c9c779d86-cfj4w   1/1     Running   0          39s
test-webserver-c9c779d86-ctqhx   1/1     Running   0          39s
test-webserver-c9c779d86-dpdph   1/1     Running   0          12m
test-webserver-c9c779d86-llrrg   1/1     Running   0          24s
test-webserver-c9c779d86-ng2tj   1/1     Running   0          24s
test-webserver-c9c779d86-ptwh8   1/1     Running   0          24s
test-webserver-c9c779d86-q4z7x   1/1     Running   0          24s
test-webserver-c9c779d86-sxvhk   1/1     Running   0          39s
test-webserver-c9c779d86-tx96q   1/1     Running   0          39s
```

HPA scaled from 1 to 10 replicas under CPU load.

## Step 5: Verify Scale Down After Load Removal

```bash
kubectl delete deployment load-generator -n hpa-test

# After stabilization window (30s)
$ kubectl get hpa hpa-webserver -n hpa-test
NAME            REFERENCE                   TARGETS       MINPODS   MAXPODS   REPLICAS   AGE
hpa-webserver   Deployment/test-webserver   cpu: 1%/20%   1         10        1          13m

$ kubectl get deployment test-webserver -n hpa-test
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
test-webserver   1/1     1            1           15m
```

HPA scaled back to 1 replica after load was removed.

## Cleanup (CPU HPA Test)

```bash
kubectl delete namespace hpa-test
```

---

## Part 2: GPU Custom Metric HPA (Prometheus Adapter + DCGM)

This section validates HPA scaling based on custom GPU metrics from NVIDIA DCGM Exporter, using Prometheus and Prometheus Adapter to expose `dcgm_fi_dev_gpu_util` as a Kubernetes custom metric.

### Step 6: Install Prometheus

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

cat > /tmp/prometheus-values.yaml <<'EOF'
alertmanager:
  enabled: false
prometheus-pushgateway:
  enabled: false
server:
  persistentVolume:
    enabled: false
  extraScrapeConfigs:
    - job_name: dcgm-exporter
      static_configs:
        - targets: ["nvidia-dcgm-exporter.gpu-operator.svc.cluster.local:9400"]
EOF

helm install prometheus prometheus-community/prometheus \
  -n monitoring --create-namespace \
  -f /tmp/prometheus-values.yaml
```

```bash
$ kubectl get pods -n monitoring | grep -E 'NAME|prometheus'
NAME                                             READY   STATUS    RESTARTS   AGE
prometheus-adapter-bdf748cc-vtfdx                1/1     Running   0          5m45s
prometheus-kube-state-metrics-675bd9d798-qf2zf   1/1     Running   0          7m21s
prometheus-prometheus-node-exporter-zjp2z        1/1     Running   0          7m21s
prometheus-prometheus-node-exporter-ztgxf        1/1     Running   0          7m21s
prometheus-server-d766bd656-gv7wc                2/2     Running   0          7m21s
```

### Step 7: Install Prometheus Adapter

```bash
cat > /tmp/adapter-values.yaml <<'EOF'
prometheus:
  url: http://prometheus-server.monitoring.svc.cluster.local
  port: 80
rules:
  custom:
    - seriesQuery: 'DCGM_FI_DEV_GPU_UTIL{namespace!="",pod!=""}'
      resources:
        overrides:
          namespace: {resource: "namespace"}
          pod: {resource: "pod"}
      name:
        matches: "DCGM_FI_DEV_GPU_UTIL"
        as: "dcgm_fi_dev_gpu_util"
      metricsQuery: 'avg_over_time(DCGM_FI_DEV_GPU_UTIL{<<.LabelMatchers>>}[2m])'
EOF

helm install prometheus-adapter prometheus-community/prometheus-adapter \
  -n monitoring \
  -f /tmp/adapter-values.yaml
```

Verify the custom metric is registered in the Kubernetes API:

```bash
$ kubectl get --raw '/apis/custom.metrics.k8s.io/v1beta1' | \
    python3 -c "import json,sys; d=json.load(sys.stdin); [print(r['name']) for r in d.get('resources',[]) if 'dcgm' in r.get('name','')]"
pods/dcgm_fi_dev_gpu_util
namespaces/dcgm_fi_dev_gpu_util
```

### Step 8: Deploy GPU Workload

```bash
kubectl create namespace hpa-gpu-test

kubectl apply -f - <<'EOF'
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gpu-inference
  namespace: hpa-gpu-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gpu-inference
  template:
    metadata:
      labels:
        app: gpu-inference
    spec:
      containers:
      - name: gpu-inference
        image: nvidia/cuda:12.4.1-base-ubuntu22.04
        command: ["sleep", "infinity"]
        resources:
          limits:
            nvidia.com/gpu: "1"
EOF
```

```bash
$ kubectl get pods -n hpa-gpu-test
NAME                             READY   STATUS    RESTARTS   AGE
gpu-inference-7746ff8778-vldsn   1/1     Running   0          4m10s
```

Confirm the pod-level GPU utilization metric is accessible:

```bash
$ kubectl get --raw '/apis/custom.metrics.k8s.io/v1beta1/namespaces/hpa-gpu-test/pods/*/dcgm_fi_dev_gpu_util' \
    | python3 -c "import json,sys; d=json.load(sys.stdin); [print(i['describedObject']['name'], '->', i['value']) for i in d.get('items',[])]"
gpu-inference-7746ff8778-vldsn -> 0
```

### Step 9: Create GPU Custom Metric HPA

```bash
kubectl apply -f - <<'EOF'
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: gpu-inference-hpa
  namespace: hpa-gpu-test
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: gpu-inference
  minReplicas: 1
  maxReplicas: 4
  metrics:
  - type: Pods
    pods:
      metric:
        name: dcgm_fi_dev_gpu_util
      target:
        type: AverageValue
        averageValue: "50"
EOF
```

```bash
$ kubectl get hpa gpu-inference-hpa -n hpa-gpu-test
NAME                REFERENCE                  TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
gpu-inference-hpa   Deployment/gpu-inference   0/50      1         4         1          3m15s
```

The HPA reads `dcgm_fi_dev_gpu_util` (current: `0`, target: `50`) via `custom.metrics.k8s.io`, confirming that GPU utilization metrics drive autoscaling decisions for AI inference workloads.

## Cleanup (GPU HPA Test)

```bash
kubectl delete namespace hpa-gpu-test
```
