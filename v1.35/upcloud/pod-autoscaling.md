## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator

## Enable time-slicing
Follow guide at [util-time-slicing.md](util-time-slicing.md)

## Install metrics server
```bash
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
kubectl patch deployment metrics-server -n kube-system --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

## Deploy app that will be scaled
```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
spec:
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      containers:
      - name: app
        image: k8s.gcr.io/hpa-example
        resources:
          requests:
            cpu: 100m
            nvidia.com/gpu.shared: 2
          limits:
            nvidia.com/gpu.shared: 2
---
apiVersion: v1
kind: Service
metadata:
  name: example
spec:
  ports:
  - port: 80
  selector:
    app: example
EOF
```

## Define HPA
```bash
kubectl apply -f - <<EOF
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: example
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: example
  minReplicas: 1
  maxReplicas: 5
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 60
EOF
```


## Results

```bash
 kubectl describe hpa example
Name:                                                  example
Namespace:                                             default
Labels:                                                <none>
Annotations:                                           <none>
CreationTimestamp:                                     Tue, 28 Apr 2026 16:23:34 +0300
Reference:                                             Deployment/example
Metrics:                                               ( current / target )
  resource cpu on pods  (as a percentage of request):  956% (956m) / 60%
Min replicas:                                          1
Max replicas:                                          5
Deployment pods:                                       1 current / 4 desired
Conditions:
  Type            Status  Reason            Message
  ----            ------  ------            -------
  AbleToScale     True    SucceededRescale  the HPA controller was able to update the target scale to 4
  ScalingActive   True    ValidMetricFound  the HPA was able to successfully calculate a replica count from cpu resource utilization (percentage of request)
  ScalingLimited  True    ScaleUpLimit      the desired replica count is increasing faster than the maximum scale rate
Events:
  Type    Reason             Age   From                       Message
  ----    ------             ----  ----                       -------
  Normal  SuccessfulRescale  6s    horizontal-pod-autoscaler  New size: 4; reason: cpu resource utilization (percentage of request) above target
```

```bash
kubectl describe pod -l app=example | grep -A 5 "Requests:"
    Requests:
      cpu:                    100m
      nvidia.com/gpu.shared:  2
...
```
