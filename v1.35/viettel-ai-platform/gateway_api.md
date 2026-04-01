# Gateway API Validation

## Overview

Viettel AI Platform supports the Kubernetes Gateway API (`gateway.networking.k8s.io/v1`) with Envoy Gateway as the implementation and Envoy AI Gateway for AI-specific traffic management (LLM routing, model-backend selection, token-based rate limiting). This guide demonstrates installation and validation of the Gateway API on the platform.

## Step 1: Verify Gateway API Resources

```bash
$ kubectl api-resources --api-group=gateway.networking.k8s.io
NAME                 SHORTNAMES   APIVERSION                           NAMESPACED   KIND
backendtlspolicies   btlspolicy   gateway.networking.k8s.io/v1         true         BackendTLSPolicy
gatewayclasses       gc           gateway.networking.k8s.io/v1         false        GatewayClass
gateways             gtw          gateway.networking.k8s.io/v1         true         Gateway
grpcroutes                        gateway.networking.k8s.io/v1         true         GRPCRoute
httproutes                        gateway.networking.k8s.io/v1         true         HTTPRoute
listenersets         lset         gateway.networking.k8s.io/v1         true         ListenerSet
referencegrants      refgrant     gateway.networking.k8s.io/v1         true         ReferenceGrant
tcproutes                         gateway.networking.k8s.io/v1alpha2   true         TCPRoute
tlsroutes                         gateway.networking.k8s.io/v1         true         TLSRoute
udproutes                         gateway.networking.k8s.io/v1alpha2   true         UDPRoute
```

## Step 2: Verify Envoy Gateway and AI Gateway Controllers

```bash
$ kubectl get pods -n envoy-gateway-system
NAME                             READY   STATUS    RESTARTS      AGE
envoy-gateway-5d54cdccd6-f55xq   1/1     Running   1 (28h ago)   36h

$ kubectl get pods -n envoy-ai-gateway-system
NAME                                     READY   STATUS    RESTARTS      AGE
ai-gateway-controller-6495f95969-nsdth   1/1     Running   1 (27h ago)   36h
```

## Step 3: Create GatewayClass, Gateway, Backend, and HTTPRoute

```bash
kubectl apply -f - <<'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: GatewayClass
metadata:
  name: eg
spec:
  controllerName: gateway.envoyproxy.io/gatewayclass-controller
---
apiVersion: v1
kind: Namespace
metadata:
  name: gw-test
---
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: test-gateway
  namespace: gw-test
spec:
  gatewayClassName: eg
  listeners:
  - name: http
    protocol: HTTP
    port: 80
    allowedRoutes:
      namespaces:
        from: Same
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  namespace: gw-test
spec:
  replicas: 1
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: python:3.11-slim
        command: ["python3", "-c"]
        args:
        - |
          from http.server import HTTPServer, BaseHTTPRequestHandler
          class H(BaseHTTPRequestHandler):
            def do_GET(self):
              self.send_response(200)
              self.end_headers()
              self.wfile.write(b'Hello from Viettel AI Platform')
            def log_message(self, *a): pass
          HTTPServer(('0.0.0.0', 8080), H).serve_forever()
        ports:
        - containerPort: 8080
---
apiVersion: v1
kind: Service
metadata:
  name: backend
  namespace: gw-test
spec:
  selector:
    app: backend
  ports:
  - port: 8080
    targetPort: 8080
---
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: backend-route
  namespace: gw-test
spec:
  parentRefs:
  - name: test-gateway
  hostnames:
  - test.viettelai.local
  rules:
  - matches:
    - path:
        type: PathPrefix
        value: /
    backendRefs:
    - name: backend
      port: 8080
EOF
```

## Step 4: Verify Gateway is Programmed

```bash
$ kubectl get gateway test-gateway -n gw-test
NAME           CLASS   ADDRESS        PROGRAMMED   AGE
test-gateway   eg      10.24.10.206   True         3m48s

$ kubectl get httproute -n gw-test
NAME            HOSTNAMES                  AGE
backend-route   ["test.viettelai.local"]   3m53s
```

The Gateway has been assigned address `10.24.10.206` and is `Programmed: True`.

## Step 5: Test L7 Traffic Routing Through the Gateway

```bash
$ curl -v -H 'Host: test.viettelai.local' http://10.24.10.206:80/
> GET / HTTP/1.1
> Host: test.viettelai.local
> User-Agent: curl/8.7.1
> Accept: */*
>
< HTTP/1.1 200 OK
< server: BaseHTTP/0.6 Python/3.11.15
< date: Wed, 01 Apr 2026 09:34:23 GMT
< transfer-encoding: chunked
<
Hello from Viettel AI Platform
```

**The above test confirms that an HTTP request with `Host: test.viettelai.local` forwarded to the Envoy Gateway is correctly routed to the backend service (L7 routing).**

## Cleanup

```bash
kubectl delete namespace gw-test
kubectl delete gatewayclass eg
```
