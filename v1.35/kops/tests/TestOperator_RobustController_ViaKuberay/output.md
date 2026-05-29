

# Robust Controller (with KubeRay)


## Verify KubeRay with a sample RayJob


Creating test namespace "testoperator-robustcontroller-viakuberay-1774009591"


Applying manifest "testdata/rayjob-sample.yaml" to namespace "testoperator-robustcontroller-viakuberay-1774009591"
```bash
> kubectl apply -n testoperator-robustcontroller-viakuberay-1774009591 -f testdata/rayjob-sample.yaml

```

```

rayjob.ray.io/rayjob-sample created
configmap/ray-job-code-sample created

```

```bash
> kubectl wait -n testoperator-robustcontroller-viakuberay-1774009591 --for='jsonpath={.status.jobDeploymentStatus}=Complete' rayjob/rayjob-sample --timeout=300s

```

```

rayjob.ray.io/rayjob-sample condition met

```

```bash
> kubectl logs -n testoperator-robustcontroller-viakuberay-1774009591 -l=job-name=rayjob-sample

```

```

2026-03-20 05:27:35,047	INFO worker.py:1694 -- Connecting to existing Ray cluster at address: 100.96.4.190:6379...
2026-03-20 05:27:35,057	INFO worker.py:1879 -- Connected to Ray cluster. View the dashboard at [1m[32m100.96.4.190:8265 [39m[22m
test_counter got 1
test_counter got 2
test_counter got 3
test_counter got 4
test_counter got 5
2026-03-20 05:27:43,618	SUCC cli.py:65 -- [32m-----------------------------------[39m
2026-03-20 05:27:43,618	SUCC cli.py:66 -- [32mJob 'rayjob-sample-tbqxx' succeeded[39m
2026-03-20 05:27:43,618	SUCC cli.py:67 -- [32m-----------------------------------[39m

```

&check; Found succeeded message in logs, indicating the RayJob completed successfully.



Deleting test namespace "testoperator-robustcontroller-viakuberay-1774009591"


Namespace deletion took 18s
