

# Gang Scheduling (via Kueue)


## Simple gang scheduling test using Kueue


Creating a Kueue Job that requires gang scheduling


Creating test namespace "gangscheduling-kueue-1774011087"


Applying manifest "testdata/gangscheduling-kueue.yaml" to namespace "gangscheduling-kueue-1774011087"
```bash
> kubectl apply -n gangscheduling-kueue-1774011087 -f testdata/gangscheduling-kueue.yaml

```

```

resourceflavor.kueue.x-k8s.io/default-flavor created
clusterqueue.kueue.x-k8s.io/cluster-queue created
localqueue.kueue.x-k8s.io/team1 created
job.batch/gangscheduling-kueue created


Warning: This version is deprecated. Use v1beta2 instead.

```



Waiting for Job to complete
```bash
> kubectl wait --namespace gangscheduling-kueue-1774011087 --for=condition=complete job/gangscheduling-kueue --timeout=300s

```

```

job.batch/gangscheduling-kueue condition met

```

&check; Gang scheduling via Kueue test completed successfully.



Deleting test namespace "gangscheduling-kueue-1774011087"


Namespace deletion took 6s
