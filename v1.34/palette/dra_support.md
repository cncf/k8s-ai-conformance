# Spectro Cloud Palette DRA Support Validation

This document provides evidence for Dynamic Resource Allocation (DRA) support on a live PXK 1.34 cluster deployed and managed with Spectro Cloud Palette.

## Cluster Version

    kubectl version

Output:

    Client Version: v1.35.2
    Kustomize Version: v5.7.1
    Server Version: v1.34.2

## DRA API Availability

    kubectl api-resources | grep resource.k8s.io

Output:

    NAME                                SHORTNAMES                          APIVERSION                                NAMESPACED   KIND
    deviceclasses                                                           resource.k8s.io/v1                        false        DeviceClass
    resourceclaims                                                          resource.k8s.io/v1                        true         ResourceClaim
    resourceclaimtemplates                                                  resource.k8s.io/v1                        true         ResourceClaimTemplate
    resourceslices                                                          resource.k8s.io/v1                        false        ResourceSlice

## DRA API Discovery

    kubectl get --raw /apis/resource.k8s.io/v1

Output (truncated):

    {
      "kind": "APIResourceList",
      "apiVersion": "v1",
      "groupVersion": "resource.k8s.io/v1"
    }

## Platform Support for DRA

Spectro Cloud Palette manages upstream-conformant Kubernetes clusters (PXK 1.34) and does not restrict the use of Kubernetes Dynamic Resource Allocation.

Clusters expose the `resource.k8s.io/v1` API group, enabling the installation and use of DRA-compatible resource drivers (e.g., GPU or custom device drivers).

Palette supports deployment of additional components (operators, device plugins, and drivers) via its pack-based architecture, allowing users to install and operate DRA drivers as needed.

## Conclusion

A live PXK 1.34 cluster managed with Spectro Cloud Palette exposes the Kubernetes DRA v1 API group (`resource.k8s.io/v1`) and the expected DRA resource types, including DeviceClass, ResourceClaim, ResourceClaimTemplate, and ResourceSlice.

This demonstrates that Spectro Cloud Palette supports Kubernetes Dynamic Resource Allocation (DRA) and allows DRA-compatible resource drivers to be installed and used without restriction.