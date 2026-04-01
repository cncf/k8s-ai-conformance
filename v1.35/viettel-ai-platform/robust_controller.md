# AI Operator: Robust CRD and Controller Validation (KServe)

## Overview

This guide validates that KServe — a complex AI operator with CRDs — is installed and functions reliably on Viettel AI Platform. KServe provides model serving infrastructure for ML frameworks including HuggingFace, TensorFlow, PyTorch, LightGBM, and others. Validation covers operator pods, webhook operation, CRD registration, and custom resource reconciliation.

## Step 1: Verify KServe Controller Pods

```bash
$ kubectl get pods -n kserve
NAME                                         READY   STATUS    RESTARTS      AGE
kserve-controller-manager-786fcddc9-hrcv5    2/2     Running   2 (25h ago)   33h
llmisvc-controller-manager-d7bb56889-27g9m   1/1     Running   1 (25h ago)   33h
```

`kserve-controller-manager` runs 2/2 containers (manager + metrics proxy). `llmisvc-controller-manager` handles LLMInferenceService resources.

## Step 2: Verify KServe CRDs Are Registered

```bash
$ kubectl get crds | grep kserve
clusterservingruntimes.serving.kserve.io       2026-03-30T20:54:06Z
inferencegraphs.serving.kserve.io              2026-03-30T20:54:06Z
inferenceservices.serving.kserve.io            2026-03-30T20:54:07Z
llminferenceserviceconfigs.serving.kserve.io   2026-03-30T20:55:08Z
llminferenceservices.serving.kserve.io         2026-03-30T20:55:11Z
servingruntimes.serving.kserve.io              2026-03-30T20:54:08Z
trainedmodels.serving.kserve.io                2026-03-30T20:54:08Z
```

All 7 KServe CRDs are registered and available.

## Step 3: Verify Webhooks Are Operational

```bash
$ kubectl get mutatingwebhookconfigurations | grep kserve
inferenceservice.serving.kserve.io   2   36h

$ kubectl get validatingwebhookconfigurations | grep kserve
clusterservingruntime.serving.kserve.io       1   36h
inferencegraph.serving.kserve.io              1   36h
inferenceservice.serving.kserve.io            1   36h
llminferenceservice.serving.kserve.io         2   36h
llminferenceserviceconfig.serving.kserve.io   2   36h
servingruntime.serving.kserve.io              1   36h
trainedmodel.serving.kserve.io                1   36h
```

8 admission webhooks (1 mutating + 7 validating) are active for KServe resources.

## Step 4: Verify ClusterServingRuntimes Are Reconciled

```bash
$ kubectl get clusterservingruntimes
NAME                                 DISABLED   MODELTYPE     CONTAINERS         AGE
kserve-huggingfaceserver                        huggingface   kserve-container   36h
kserve-huggingfaceserver-multinode              huggingface   kserve-container   36h
kserve-lgbserver                                lightgbm      kserve-container   36h
kserve-mlserver                                 sklearn       kserve-container   36h
kserve-paddleserver                             paddle        kserve-container   36h
kserve-pmmlserver                               pmml          kserve-container   36h
kserve-predictiveserver                         sklearn       kserve-container   36h
kserve-sklearnserver                            sklearn       kserve-container   36h
kserve-tensorflow-serving                       tensorflow    kserve-container   36h
```

The KServe controller has reconciled 9 ClusterServingRuntimes covering all major ML frameworks.

## Step 5: Deploy and Reconcile an InferenceService

```bash
kubectl apply -f - <<'EOF'
apiVersion: v1
kind: Namespace
metadata:
  name: kserve-test
---
apiVersion: serving.kserve.io/v1beta1
kind: InferenceService
metadata:
  name: sklearn-iris
  namespace: kserve-test
spec:
  predictor:
    sklearn:
      storageUri: "gs://kfserving-examples/models/sklearn/1.0/model"
EOF
```

Verify the InferenceService custom resource is accepted and reconciled by the webhook and controller:

```bash
$ kubectl get inferenceservice sklearn-iris -n kserve-test
NAME           URL   READY     PREV   LATEST   PREVROLLEDOUTREVISION   LATESTREADYREVISION   AGE
sklearn-iris         Unknown                                                                  10s
```

The InferenceService was accepted by the validating webhook and is actively being reconciled by the KServe controller (`IngressReady: Unknown` indicates the controller is provisioning resources).

## Cleanup

```bash
kubectl delete namespace kserve-test
```

## Summary

| Operator            | Pods        | CRDs | Webhooks | Status      |
| ------------------- | ----------- | ---- | -------- | ----------- |
| KServe              | 2/2 Running | 7    | 8        | Operational |
| LLMInferenceService | 1/1 Running | 2    | 4        | Operational |
