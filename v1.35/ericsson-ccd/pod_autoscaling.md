# Horizontal Pod Autoscaling Verification Guide

This guide covers the verification of Horizontal Pod Autoscaling on Ericsson Cloud Container Distribution (ECCD).

> **Note:** All the commands in the following steps are to be executed on a control plane node of the Ericsson CCD deployment.

## Step 1: Create Namespace and Deploy Test Application

```bash
kubectl create namespace hpa-test
```

```bash
cat <<EOF | kubectl apply -f -
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
      securityContext:
        runAsNonRoot: true
        runAsUser: 2000
        runAsGroup: 2000
        fsGroup: 2000
      containers:
      - name: webserver
        image: python:3.11-slim
        command: ["python3", "-c"]
        args:
        - |
          import signal, sys
          from http.server import HTTPServer, BaseHTTPRequestHandler
          class Handler(BaseHTTPRequestHandler):
            def do_GET(self):
              self.send_response(200)
              self.end_headers()
              self.wfile.write(b'OK')
          def shutdown(signum, frame):
            raise SystemExit(0)
          signal.signal(signal.SIGTERM, shutdown)
          server = HTTPServer(('0.0.0.0', 8080), Handler)
          server.serve_forever()
        imagePullPolicy: Always
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 256Mi
        securityContext:
          allowPrivilegeEscalation: false
          runAsNonRoot: true
          capabilities:
            drop:
            - ALL
EOF
```

```bash
# Verify pod is running
kubectl get pods -n hpa-test -l app=test-webserver
```

**Expected Output:**
```
NAME                              READY   STATUS    RESTARTS   AGE
test-webserver-6d9f8b7c4d-xk2pq   1/1     Running   0          30s
```

## Step 2: Create HorizontalPodAutoscaler Resource

```bash
cat <<EOF | kubectl apply -f -
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
# Verify HPA is created and shows CPU metrics
kubectl get hpa hpa-webserver -n hpa-test
```

**Expected Output:**
```
NAME            REFERENCE                   TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
hpa-webserver   Deployment/test-webserver   0%/20%    1         10        1          30s
```

## Step 3: Verify Autoscaling Under Load

```bash
cat <<EOF | kubectl apply -f -
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
      restartPolicy: Always
      securityContext:
        runAsNonRoot: true
        runAsUser: 1000
      containers:
      - name: load-generator
        image: python:3.11-slim
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
        command: ["python3", "-c"]
        args:
        - |
          import signal, urllib.request
          def shutdown(signum, frame):
            raise SystemExit(0)
          signal.signal(signal.SIGTERM, shutdown)
          url = 'http://test-webserver.hpa-test.svc.cluster.local:8080'
          while True:
            try:
              urllib.request.urlopen(url, timeout=5)
            except Exception:
              pass
        resources:
          requests:
            cpu: 100m
          limits:
            cpu: 1
EOF
```

```bash
# Verify load generators are running
kubectl get pods -n hpa-test -l app=load-generator
```

**Expected Output:**
```
NAME                               READY   STATUS    RESTARTS   AGE
load-generator-7b9f6d8c5f-4jxpq    1/1     Running   0          15s
load-generator-7b9f6d8c5f-9mrkv    1/1     Running   0          15s
load-generator-7b9f6d8c5f-tz2wn    1/1     Running   0          15s
```

```bash
# Monitor HPA for scale-up
kubectl get hpa hpa-webserver -n hpa-test -w
```

**Expected Output:**
```
NAME            REFERENCE                   TARGETS    MINPODS   MAXPODS   REPLICAS   AGE
hpa-webserver   Deployment/test-webserver   0%/20%     1         10        1          2m
hpa-webserver   Deployment/test-webserver   85%/20%    1         10        1          3m
hpa-webserver   Deployment/test-webserver   85%/20%    1         10        5          3m30s
hpa-webserver   Deployment/test-webserver   32%/20%    1         10        5          4m
```

## Step 4: Verify Scale Down After the Load is Reduced

```bash
# Delete load generators to remove CPU pressure
kubectl delete deployment load-generator -n hpa-test

# Monitor HPA for scale-down (stabilization window: 30s)
kubectl get hpa hpa-webserver -n hpa-test -w
```

**Expected Output:**
```
NAME            REFERENCE                   TARGETS   MINPODS   MAXPODS   REPLICAS   AGE
hpa-webserver   Deployment/test-webserver   32%/20%   1         10        5          8m
hpa-webserver   Deployment/test-webserver   0%/20%    1         10        5          9m
hpa-webserver   Deployment/test-webserver   0%/20%    1         10        1          9m30s
```

```bash
# Confirm deployment scaled back to 1 replica
kubectl get deployment test-webserver -n hpa-test
```

**Expected Output:**
```
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
test-webserver   1/1     1            1           10m
```

## Cleanup

```bash
kubectl delete namespace hpa-test
```
