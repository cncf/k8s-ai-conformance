

# DRA API Availability


## Checking for DRA API version v1
```bash
> kubectl api-versions | grep resource.k8s.io/v1

```

```

resource.k8s.io/v1

```

&check; DRA API version resource.k8s.io/v1 is available.



## Checking for deviceclasses
```bash
> kubectl api-resources --api-group=resource.k8s.io | grep deviceclasses

```

```

deviceclasses                         resource.k8s.io/v1   false        DeviceClass

```

&check; DRA API resource deviceclasses is available.



## Checking for resourceclaims
```bash
> kubectl api-resources --api-group=resource.k8s.io | grep resourceclaims

```

```

resourceclaims                        resource.k8s.io/v1   true         ResourceClaim

```

&check; DRA API resource resourceclaims is available.



## Checking for resourceclaimtemplates
```bash
> kubectl api-resources --api-group=resource.k8s.io | grep resourceclaimtemplates

```

```

resourceclaimtemplates                resource.k8s.io/v1   true         ResourceClaimTemplate

```

&check; DRA API resource resourceclaimtemplates is available.



## Checking for resourceslices
```bash
> kubectl api-resources --api-group=resource.k8s.io | grep resourceslices

```

```

resourceslices                        resource.k8s.io/v1   false        ResourceSlice

```

&check; DRA API resource resourceslices is available.

