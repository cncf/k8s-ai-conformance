

# Gateway API support for AI inference


## Verify Weighted Traffic Splitting


Creating test namespace "weighted-traffic-splitting-1774009334"


Applying manifest "testdata/weighted-traffic-splitting.yaml" to namespace "weighted-traffic-splitting-1774009334"
```bash
> kubectl apply -n weighted-traffic-splitting-1774009334 -f testdata/weighted-traffic-splitting.yaml

```

```

gateway.gateway.networking.k8s.io/gateway created
service/app1 created
deployment.apps/app1 created
service/app2 created
deployment.apps/app2 created
httproute.gateway.networking.k8s.io/weighted-traffic-splitting created

```



Waiting for HTTPRoute weighted-traffic-splitting to be accepted
```bash
> kubectl wait -n weighted-traffic-splitting-1774009334 --for=jsonpath='{.status.parents[0].conditions[?(@.type=="Accepted")].status}'=True httproute.gateway.networking.k8s.io/weighted-traffic-splitting --timeout=300s

```

```

httproute.gateway.networking.k8s.io/weighted-traffic-splitting condition met

```



Deleting test namespace "weighted-traffic-splitting-1774009334"


Namespace deletion took 42s


## Verify Header Based Routing


Creating test namespace "header-based-routing-1774009378"


Applying manifest "testdata/header-based-routing.yaml" to namespace "header-based-routing-1774009378"
```bash
> kubectl apply -n header-based-routing-1774009378 -f testdata/header-based-routing.yaml

```

```

gateway.gateway.networking.k8s.io/gateway created
service/app1 created
deployment.apps/app1 created
service/app2 created
deployment.apps/app2 created
httproute.gateway.networking.k8s.io/header-based-routing created

```



Waiting for HTTPRoute header-based-routing to be accepted
```bash
> kubectl wait -n header-based-routing-1774009378 --for=jsonpath='{.status.parents[0].conditions[?(@.type=="Accepted")].status}'=True httproute.gateway.networking.k8s.io/header-based-routing --timeout=300s

```

```

httproute.gateway.networking.k8s.io/header-based-routing condition met

```



Deleting test namespace "header-based-routing-1774009378"


Namespace deletion took 42s
