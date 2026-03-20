

# Observability: AI Service Metrics


Found Prometheus service: monitoring/kube-prometheus-stack-prometheus


Found Prometheus endpoint: http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090


## Deploy a workload exposing Prometheus metrics and verify collection


Creating test namespace "ai-service-metrics-1774009432"


Applying manifest "testdata/ai-service-metrics.yaml" to namespace "ai-service-metrics-1774009432"
```bash
> kubectl apply -n ai-service-metrics-1774009432 -f testdata/ai-service-metrics.yaml

```

```

configmap/metrics-fake-workload created
configmap/metrics-fake-workload-nginx-conf created
deployment.apps/metrics-fake-workload created
service/metrics-fake-workload created
podmonitor.monitoring.coreos.com/metrics-fake-workload created

```



Waiting for Deployment metrics-fake-workload to be ready
```bash
> kubectl wait -n ai-service-metrics-1774009432 --for=condition=Available deployment.apps/metrics-fake-workload --timeout=300s

```

```

deployment.apps/metrics-fake-workload condition met

```



### Verify workload exposes metrics
```bash
> kubectl run direct-scrape -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS http://metrics-fake-workload.ai-service-metrics-1774009432.svc.cluster.local:8080/metrics

```

```

pod/direct-scrape created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/direct-scrape --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/direct-scrape condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 direct-scrape

```

```

# HELP aiconformancetest_requests_total A fake counter metric for AI conformance validation.
# TYPE aiconformancetest_requests_total counter
aiconformancetest_requests_total{method="GET",endpoint="/predict"} 1027
aiconformancetest_requests_total{method="POST",endpoint="/predict"} 2563
# HELP aiconformancetest_request_duration_seconds A fake histogram metric for AI conformance validation.
# TYPE aiconformancetest_request_duration_seconds histogram
aiconformancetest_request_duration_seconds_bucket{le="0.01"} 500
aiconformancetest_request_duration_seconds_bucket{le="0.1"} 1500
aiconformancetest_request_duration_seconds_bucket{le="1"} 3000
aiconformancetest_request_duration_seconds_bucket{le="+Inf"} 3590
aiconformancetest_request_duration_seconds_sum 245.7
aiconformancetest_request_duration_seconds_count 3590

```



Confirmed workload exposes metric aiconformancetest_requests_total


### Verify Prometheus collects the metric
```bash
> kubectl run prom-query-1 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-1 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-1 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-1 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-1

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 1 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 1: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-2 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-2 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-2 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-2 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-2

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 2 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 2: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-3 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-3 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-3 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-3 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-3

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 3 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 3: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-4 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-4 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-4 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-4 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-4

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 4 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 4: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-5 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-5 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-5 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-5 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-5

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 5 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 5: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-6 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-6 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-6 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-6 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-6

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 6 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 6: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-7 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-7 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-7 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-7 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-7

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 7 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 7: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-8 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-8 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-8 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-8 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-8

```

```

{"status":"success","data":{"resultType":"vector","result":[]}}

```



Prometheus query attempt 8 response: {"status":"success","data":{"resultType":"vector","result":[]}}


Attempt 8: metric not yet collected by Prometheus, retrying in 15s...
```bash
> kubectl run prom-query-9 -n ai-service-metrics-1774009432 --image=registry.k8s.io/e2e-test-images/agnhost:2.39 --restart=Never --command -- curl -sS -X POST --data-urlencode 'query=aiconformancetest_requests_total{namespace="ai-service-metrics-1774009432"}' 'http://kube-prometheus-stack-prometheus.monitoring.svc.cluster.local:9090/api/v1/query'

```

```

pod/prom-query-9 created

```

```bash
> kubectl wait -n ai-service-metrics-1774009432 pod/prom-query-9 --for=jsonpath='{.status.phase}'=Succeeded --timeout=120s

```

```

pod/prom-query-9 condition met

```

```bash
> kubectl logs -n ai-service-metrics-1774009432 prom-query-9

```

```

{"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"aiconformancetest_requests_total","container":"webserver","endpoint":"metrics","exported_endpoint":"/predict","instance":"100.96.4.87:8080","job":"ai-service-metrics-1774009432/metrics-fake-workload","method":"GET","namespace":"ai-service-metrics-1774009432","pod":"metrics-fake-workload-85c55c7fd4-q8ffg"},"value":[1774009575.241,"1027"]},{"metric":{"__name__":"aiconformancetest_requests_total","container":"webserver","endpoint":"metrics","exported_endpoint":"/predict","instance":"100.96.4.87:8080","job":"ai-service-metrics-1774009432/metrics-fake-workload","method":"POST","namespace":"ai-service-metrics-1774009432","pod":"metrics-fake-workload-85c55c7fd4-q8ffg"},"value":[1774009575.241,"2563"]}]}}

```



Prometheus query attempt 9 response: {"status":"success","data":{"resultType":"vector","result":[{"metric":{"__name__":"aiconformancetest_requests_total","container":"webserver","endpoint":"metrics","exported_endpoint":"/predict","instance":"100.96.4.87:8080","job":"ai-service-metrics-1774009432/metrics-fake-workload","method":"GET","namespace":"ai-service-metrics-1774009432","pod":"metrics-fake-workload-85c55c7fd4-q8ffg"},"value":[1774009575.241,"1027"]},{"metric":{"__name__":"aiconformancetest_requests_total","container":"webserver","endpoint":"metrics","exported_endpoint":"/predict","instance":"100.96.4.87:8080","job":"ai-service-metrics-1774009432/metrics-fake-workload","method":"POST","namespace":"ai-service-metrics-1774009432","pod":"metrics-fake-workload-85c55c7fd4-q8ffg"},"value":[1774009575.241,"2563"]}]}}


Prometheus collected metric aiconformancetest_requests_total on attempt 9: {Metric:map[__name__:aiconformancetest_requests_total container:webserver endpoint:metrics exported_endpoint:/predict instance:100.96.4.87:8080 job:ai-service-metrics-1774009432/metrics-fake-workload method:GET namespace:ai-service-metrics-1774009432 pod:metrics-fake-workload-85c55c7fd4-q8ffg] Value:[1.774009575241e+09 1027]}


Prometheus collected metric aiconformancetest_requests_total on attempt 9: {Metric:map[__name__:aiconformancetest_requests_total container:webserver endpoint:metrics exported_endpoint:/predict instance:100.96.4.87:8080 job:ai-service-metrics-1774009432/metrics-fake-workload method:POST namespace:ai-service-metrics-1774009432 pod:metrics-fake-workload-85c55c7fd4-q8ffg] Value:[1.774009575241e+09 2563]}


Successfully verified: monitoring system discovered and collected metrics from a workload exposing Prometheus exposition format


Deleting test namespace "ai-service-metrics-1774009432"


Namespace deletion took 12s
