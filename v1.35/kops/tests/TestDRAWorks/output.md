

# DRA functional tests


## Listing device classes


* compute-domain-daemon.nvidia.com


* compute-domain-default-channel.nvidia.com


* gpu.nvidia.com


* mig.nvidia.com


* vfio.gpu.nvidia.com


## Listing resource slices


* i-030791dc5a37cd257-compute-domain.nvidia.com-ks4gr


* i-030791dc5a37cd257-gpu.nvidia.com-gld2v


* i-034383b8fabf72012-compute-domain.nvidia.com-j9pjr


* i-034383b8fabf72012-gpu.nvidia.com-ns577


## Run cuda-smoketest


Creating test namespace "testdraworks-1774009307"


Applying manifest "testdata/cuda-smoketest.yaml" to namespace "testdraworks-1774009307"
```bash
> kubectl apply -n testdraworks-1774009307 -f testdata/cuda-smoketest.yaml

```

```

resourceclaim.resource.k8s.io/cuda-smoketest created
job.batch/cuda-smoketest created

```

```bash
> kubectl wait --for=condition=complete --namespace testdraworks-1774009307 job/cuda-smoketest --timeout=5m

```

```

job.batch/cuda-smoketest condition met

```

```bash
> kubectl logs --namespace testdraworks-1774009307 job/cuda-smoketest

```

```

[Vector addition of 50000 elements]
Copy input data from the host memory to the CUDA device
CUDA kernel launch with 196 blocks of 256 threads
Copy output data from the CUDA device to the host memory
Test PASSED
Done

```



Deleting test namespace "testdraworks-1774009307"


Namespace deletion took 6s
