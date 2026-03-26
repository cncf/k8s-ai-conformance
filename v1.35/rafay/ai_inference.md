# Evidence: AI Inference Traffic Management

---

## Overview

This test validates that the Kubernetes platform supports Gateway API for managing AI inference traffic. It demonstrates weighted routing between two model versions using HTTPRoute, enabling canary deployments and gradual traffic shifting for AI/ML workloads.

---

## Prerequisite

Before running this test, you must provision a Kubernetes 1.35 MKS (Managed Kubernetes Service) cluster through Rafay Platform. A Gateway API compatible ingress controller (such as Traefik) must be installed.

---

## Step 1: Install Gateway API CRDs and Traefik RBAC

Install the standard Gateway API CRDs and Traefik's Gateway API RBAC configuration:

```
$ kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml

$ kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v3.5/docs/content/reference/dynamic-configuration/kubernetes-gateway-rbac.yml
```

---

## Step 2: Verify Traefik Gateway Controller

Verify that the Traefik controller pod is running:

```
$ kubectl get po -n traefik

NAME                       READY   STATUS    RESTARTS   AGE
traefik-65cbd7986f-dlq2h   1/1     Running   0          39s
```

> **Note:** ✅ Traefik Gateway controller pod is running successfully.

---

## Step 3: Create Gateway Resource

Create a Gateway resource that listens on port 8000 for HTTP traffic:

```
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: traefik-gateway
  namespace: test
spec:
  gatewayClassName: traefik
  listeners:
    - name: web
      protocol: HTTP
      port: 8000
      allowedRoutes:
        namespaces:
          from: Same
```

---

## Step 4: Deploy Inference Model Versions

Deploy two versions of the inference service (v1 and v2) along with their corresponding Services. These simulate different model versions for A/B testing:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inference-v1
  namespace: testing
spec:
  replicas: 1
  selector:
    matchLabels:
      app: inference
      version: v1
  template:
    metadata:
      labels:
        app: inference
        version: v1
    spec:
      containers:
        - name: server
          image: hashicorp/http-echo:latest
          args: ["-listen=:80", "-text=model-v1"]
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: inference-v2
  namespace: testing
spec:
  replicas: 1
  selector:
    matchLabels:
      app: inference
      version: v2
  template:
    metadata:
      labels:
        app: inference
        version: v2
    spec:
      containers:
        - name: server
          image: hashicorp/http-echo:latest
          args: ["-listen=:80", "-text=model-v2"]
---
apiVersion: v1
kind: Service
metadata:
  name: inference-v1-svc
  namespace: testing
spec:
  selector:
    app: inference
    version: v1
  ports:
    - port: 80
---
apiVersion: v1
kind: Service
metadata:
  name: inference-v2-svc
  namespace: testing
spec:
  selector:
    app: inference
    version: v2
  ports:
    - port: 80
```

---

## Step 5: Verify Inference Pods

Verify that both inference model version pods are running:

```
$ kubectl get po -n test

NAME                            READY   STATUS    RESTARTS   AGE
inference-v1-64f9cd688f-jjq75   1/1     Running   0          6s
inference-v2-647b44d5ff-tvfmk   1/1     Running   0          6s
```

> **Note:** ✅ Both inference model version pods (v1 and v2) are running.

---

## Step 6: Create HTTPRoute with Weighted Routing

Create an HTTPRoute that distributes traffic between the two model versions with an 80/20 weight split (80% to v1, 20% to v2):

```
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: inference-weighted
  namespace: testing
spec:
  parentRefs:
    - name: traefik-gateway
      sectionName: web
  hostnames:
    - inference.localhost
  rules:
    - backendRefs:
        - name: inference-v1-svc
          port: 80
          weight: 80
        - name: inference-v2-svc
          port: 80
          weight: 20
```

---

## Step 7: Verify HTTPRoute Creation

Verify that the HTTPRoute was created successfully:

```
$ kubectl get httproutes -A

NAMESPACE   NAME                 HOSTNAMES                 AGE
test        inference-weighted   ["inference.localhost"]   32s
```

> **Note:** ✅ HTTPRoute with weighted routing configuration created successfully.

---

## Step 8: Test Weighted Traffic Distribution

Test the weighted routing by sending multiple requests and observing the distribution between model versions:

```
$ kubectl -n traefik port-forward svc/traefik 18080:80 &
```

```
$ for i in {1..10}; do curl -s -H 'Host: inference.localhost' http://127.0.0.1:18080; echo; done

Handling connection for 18080
model-v1

Handling connection for 18080
model-v1

Handling connection for 18080
model-v1

Handling connection for 18080
model-v2

Handling connection for 18080
model-v1

Handling connection for 18080
model-v1

Handling connection for 18080
model-v1

Handling connection for 18080
model-v1

Handling connection for 18080
model-v1

Handling connection for 18080
model-v2
```

> **Note:** ✅ The weighted routing is working correctly. Out of 10 requests, approximately 80% (8 requests) went to model-v1 and 20% (2 requests) went to model-v2, matching the configured weight distribution.

