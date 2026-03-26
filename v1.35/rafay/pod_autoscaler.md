# Horizontal Pod Autoscaling

This document demonstrates Horizontal Pod Autoscaling (HPA) on a Rafay MKS cluster. An HPA resource is configured to target 20% average CPU utilization on a test webserver deployment, scaling between 1 and 10 replicas. Load generators are deployed to drive CPU usage above the threshold, triggering scale-up. Once the load is removed, the HPA scales the deployment back down to 1 replica.

## Prerequisites

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform. The cluster must have a working metrics-server so that the HPA can read CPU utilization. A Rafay v1.35 MKS cluster with the default blueprint includes metrics-server preinstalled.



## Steps

### Step 1: Deploy the test webserver Service and Deployment

Create a ClusterIP Service and a single-replica Deployment running a Python HTTP server with CPU requests (100m) and limits (500m) in the `test` namespace.

```
cat <<EOF | kubectl apply -f -
apiVersion: v1
kind: Service
metadata:
  name: test-webserver
  namespace: test
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
  namespace: test
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

### Step 2: Create the HorizontalPodAutoscaler

Create an HPA targeting the test-webserver Deployment with 20% average CPU utilization threshold, scaling between 1 and 10 replicas, and a 30-second scale-down stabilization window.

```
cat <<EOF | kubectl apply -f -
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: hpa-webserver
  namespace: test
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

### Step 3: Deploy load generators

Deploy 3 load-generator pods that continuously send HTTP requests to the test-webserver to drive up CPU utilization.

```
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: load-generator
  namespace: test
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



### Step 4: Verify load generators and webserver pods are running

Confirm all load-generator and test-webserver pods are running.

```
$ kubectl get po -n test

NAME                              READY   STATUS    RESTARTS   AGE
load-generator-5b5764d8fd-8fjjz   1/1     Running   0          14s
load-generator-5b5764d8fd-l4nch   1/1     Running   0          17s
load-generator-5b5764d8fd-lvkrb   1/1     Running   0          16s
test-webserver-d54977bb6-988dt    1/1     Running   0          4m20s
```

### Step 5: Observe CPU usage under load

Verify that the load generators are consuming CPU and driving up the test-webserver's CPU usage.

```
$ kubectl top pods -n test

NAME                              CPU(cores)   MEMORY(bytes)
load-generator-5b5764d8fd-8fjjz   261m         14Mi
load-generator-5b5764d8fd-l4nch   241m         14Mi
load-generator-5b5764d8fd-lvkrb   251m         14Mi
test-webserver-d54977bb6-988dt    45m          13Mi
```

### Step 6: Confirm HPA has scaled up replicas

Check the HPA status. CPU utilization has risen to 210% (well above the 20% target), causing the HPA to scale the deployment to 10 replicas (the configured maximum).

```
$ kubectl get hpa -A

NAMESPACE   NAME            REFERENCE                   TARGETS         MINPODS   MAXPODS   REPLICAS   AGE
test        hpa-webserver   Deployment/test-webserver   cpu: 210%/20%   1         10        10          6m59s
```

### Step 7: Verify all scaled-up pods are running

Confirm that the HPA has scaled the test-webserver deployment up to 10 replicas to handle the load.

```
$ kubectl get po -n test

NAME                              READY   STATUS    RESTARTS   AGE
load-generator-5b5764d8fd-8fjjz   1/1     Running   0          3m18s
load-generator-5b5764d8fd-l4nch   1/1     Running   0          3m21s
load-generator-5b5764d8fd-lvkrb   1/1     Running   0          3m20s
test-webserver-d54977bb6-4vj76    1/1     Running   0          83s
test-webserver-d54977bb6-6x4sr    1/1     Running   0          2m23s
test-webserver-d54977bb6-988dt    1/1     Running   0          7m24s
test-webserver-d54977bb6-bdgzv    1/1     Running   0          23s
test-webserver-d54977bb6-c85cc    1/1     Running   0          23s
test-webserver-d54977bb6-gcgpz    1/1     Running   0          83s
test-webserver-d54977bb6-lsl95    1/1     Running   0          83s
test-webserver-d54977bb6-n4qt5    1/1     Running   0          23s
test-webserver-d54977bb6-sf9f2    1/1     Running   0          2m23s
test-webserver-d54977bb6-vlgmx    1/1     Running   0          83s
```

### Step 8: Remove load and verify scale-down

Delete the load-generator deployment to remove CPU pressure.

```
$ kubectl delete deploy -n test load-generator
```

### Step 9: Confirm HPA has scaled back down

After the stabilization window, the HPA scales the deployment back to 1 replica as CPU utilization drops to 1%.

```
$ kubectl get hpa -A

NAMESPACE   NAME            REFERENCE                   TARGETS       MINPODS   MAXPODS   REPLICAS   AGE
test        hpa-webserver   Deployment/test-webserver   cpu: 1%/20%   1         10        1          11m
```

### Step 10: Verify single replica is running

Confirm only the original pod remains, proving the scale-down completed successfully.

```
$ kubectl get po -n test

NAME                             READY   STATUS    RESTARTS   AGE
test-webserver-d54977bb6-988dt   1/1     Running   0          12m
```

## Summary

The Horizontal Pod Autoscaler successfully:

- Detected increased CPU utilization (210%) caused by load-generator pods
- Scaled the test-webserver deployment from 1 replica up to 10 replicas (the configured maximum)
- Scaled the deployment back down to 1 replica after the load was removed and CPU dropped to 1%

This confirms that CPU-based Horizontal Pod Autoscaling is functioning correctly on the Rafay MKS cluster.
