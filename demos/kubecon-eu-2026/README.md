## AI Conformance Demo

This demo shows the AI Conformance requirements related to Gateway API, Gateway API Inference Extension (GAIE), and disaggregated inference.

The demo uses an AI conformant GKE cluster with 12 TPU7x and an llm-d P/D deployment (1 prefill and 2 decode workers) as an example, hosting a 480B sparse MoE model (35B active parameters) in FP8-Dynamic. It can be adapted to run on a different platform + accelerators + disaggregated inference solution combination.

```bash
# Download demo-magic.sh
curl -O https://raw.githubusercontent.com/paxtonhare/demo-magic/refs/heads/master/demo-magic.sh
# Run the demo on an AI conformant cluster with accelerators, GAIE, and disaggrated inference (adjust namespace and pod label selector if not using llm-d)
./run-demo-eu.sh
```

See it in action:

[<img src="https://img.youtube.com/vi/RRhPMuKPZG8/maxresdefault.jpg" width="560" height="315"
/>](https://www.youtube.com/embed/RRhPMuKPZG8?si=6niy8XUIE1aFvDjq)

Watch the KubeCon EU keynote recording:

[<img src="https://img.youtube.com/vi/GkeNG3wxUfY/maxresdefault.jpg" width="560" height="315"
/>](https://youtu.be/GkeNG3wxUfY?si=7nbRcpTpfhVTncd7&t=1913)