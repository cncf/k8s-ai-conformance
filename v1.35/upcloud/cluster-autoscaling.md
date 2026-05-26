## Prerequisites
* Have UKS cluster running
* Have a node group that has Nvidia L40S accelerator
* Have nvidia device plugin installed. guide: https://upcloud.com/docs/guides/gpu-workloads-managed-kubernetes/ 

## Install autoscaler
Use this guide to install autoscaler to cluster https://upcloud.com/docs/guides/cluster-autoscaler/ 

## Deploy a workload that needs two gpu's
```bash
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 2
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        resources:
          limits:
            nvidia.com/gpu: 1
EOF
```
## Autoscaler logs:
You can check autoscaler pod's logs to confirm that node group is being scaled up
```
orchestrator.go:108] Upcoming 0 nodes
Pod default/nginx-deployment-5bfb4984cc-fdrzr can't be scheduled on 0d530c5a-b9ca-42ba-b5fb-6078e309667c/small, predicate checking error: Insufficient nvidia.com/gpu; predicateName=NodeResourcesFit; reasons: Insufficient nvidia.com/gpu; debugInfo=
orchestrator.go:181] Best option to resize: 0d530c5a-b9ca-42ba-b5fb-6078e309667c/gpu1
orchestrator.go:150] No pod can fit to 0d530c5a-b9ca-42ba-b5fb-6078e309667c/small
orchestrator.go:185] Estimated 1 nodes needed in 0d530c5a-b9ca-42ba-b5fb-6078e309667c/gpu1
orchestrator.go:291] Final scale-up plan: [{0d530c5a-b9ca-42ba-b5fb-6078e309667c/gpu1 1->2 (max: 120)}]
executor.go:147] Scale-up: setting group 0d530c5a-b9ca-42ba-b5fb-6078e309667c/gpu1 size to 2
upcloud_node_group.go:115] scaling node group 0d530c5a-b9ca-42ba-b5fb-6078e309667c/gpu1 from 1 to 2
```
### Confirm that a new node was scheduled

```bash
kubectl get node
NAME                STATUS   ROLES    AGE     VERSION
gpu1-jrztm-bwcg6    Ready    <none>   7m19s   v1.35.3
gpu1-jrztm-p48ww    Ready    <none>   36m     v1.35.3
small-bzkkz-zk9cn   Ready    <none>   101m    v1.35.3
```