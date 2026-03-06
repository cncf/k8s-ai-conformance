# Gateway API Validation

## Overview

The Kubernetes Gateway API is a role-oriented, extensible API for managing ingress and network traffic. Compared to traditional Ingress, it provides advanced traffic management capabilities such as weighted traffic splitting, header-based routing, and traffic mirroring — making it well suited for inference service deployments where fine-grained control over request routing is required.

Gateway API with Traefik as the gateway controller can be deployed and verified on Ericsson Cloud Container Distribution 2.34.0 with the steps below.

## Step 1: Install the Gateway CRDs and add RBAC for traefik-controller to access the CRDs

```bash
#Create Gateway CRDs
kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml

#Set RBAC for traefik to access the gateway CRDs
kubectl apply -f https://raw.githubusercontent.com/traefik/traefik/v3.5/docs/content/reference/dynamic-configuration/kubernetes-gateway-rbac.yml
```

## Step 2: Install the traefik helm chart

Create values.yaml to set the Gateway provider

```bash
cat << EOF > values.yaml
providers:
  kubernetesGateway:
    enabled: true

nodeSelector:
  node-role.kubernetes.io/control-plane: ""

tolerations:
  - key: "node-role.kubernetes.io/control-plane"
    operator: "Exists"
    effect: "NoSchedule"
EOF
```

Install the helm chart for traefik

```bash
#Install the traefik helm chart

helm repo add traefik https://traefik.github.io/charts
helm repo update
helm upgrade --install traefik traefik/traefik \
  -n traefik --create-namespace \
  -f values.yaml
```


## Step 3: Create a Gateway

Create a testing namespace
```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: testing
EOF
```

Create a Gateway for the test

```bash
kubectl apply -f - <<'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: Gateway
metadata:
  name: traefik-gateway
  namespace: testing
spec:
  gatewayClassName: traefik
  listeners:
    - name: web
      protocol: HTTP
      port: 8000
      allowedRoutes:
        namespaces:
          from: Same
EOF
```

## Step 4: Create a test webserver with a service

Export the image tag for the test image

```bash
export IMAGE_TAG=$(curl -s -q https://registry.eccd.local:5000/v2/ccd-task-exec-job/tags/list | jq '.tags[]' -r | grep -v 'sha256')
```

Deploy the test webserver

```bash
cat <<EOF | envsubst | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-dep
  namespace: testing
spec:
  selector:
    matchLabels:
      app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
        - name: test
          image: registry.eccd.local:5000/ccd-task-exec-job:${IMAGE_TAG}
          command: ["python3", "-c"]
          args:
          - |
            from http.server import HTTPServer, BaseHTTPRequestHandler
            class Handler(BaseHTTPRequestHandler):
              def do_GET(self):
                self.send_response(200)
                self.end_headers()
                self.wfile.write(b'Hello from test app')
            HTTPServer(('0.0.0.0', 80), Handler).serve_forever()
---
apiVersion: v1
kind: Service
metadata:
  name: test-svc
  namespace: testing
spec:
  selector:
    app: test
  ports:
    - port: 80
      targetPort: 80
EOF
```

Ensure that the webserver is running

```bash
kubectl get po,service -n testing
NAME                            READY   STATUS    RESTARTS   AGE
pod/test-dep-6c568bf575-qzh96   1/1     Running   0          81s

NAME               TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)   AGE
service/test-svc   ClusterIP   10.108.205.55   <none>        80/TCP    81s
```

## Step 5: Create an HTTPRoute to map traffic

The HTTPRoute maps all traffic from the gateway's web listener to the backend service(test-svc)

```bash
kubectl apply -f - <<'EOF'
apiVersion: gateway.networking.k8s.io/v1
kind: HTTPRoute
metadata:
  name: test-http
  namespace: testing
spec:
  parentRefs:
    - name: traefik-gateway
      sectionName: web
  hostnames:
    - test.localhost
  rules:
    - matches:
        - path:
            type: PathPrefix
            value: /
      backendRefs:
        - name: test-svc
          port: 80
EOF
```

## Step 6: Test traffic flow through the Gateway

Forward traffic from a local port to the traefik controller service

Terminal 1
```bash
kubectl -n traefik port-forward svc/traefik --address 127.0.0.1 18080:80
Forwarding from 127.0.0.1:18080 -> 8000
```

Send an HTTP request to the localhost port with the hostname "test.localhost" from a different terminal to test L7 routing.

Terminal 2
```bash
curl -H 'Host: test.localhost' http://127.0.0.1:18080
Hello from test app

#Response from test app
```

Terminal 1
```bash
Forwarding from 127.0.0.1:18080 -> 8000
Handling connection for 18080
```

**The above test confirms that an HTTP request with hostname "test.localhost" forwarded to the traefik service is redirected to the test service(L7 routing)**
