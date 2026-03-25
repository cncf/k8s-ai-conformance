## AI Conformance Demo

This demo shows the AI Conformance requirements related to Gateway API, Gateway API Inference Extension (GAIE), and disaggregated inference.

The demo uses an AI conformant GKE cluster with TPU 7x and an llm-d P/D deployment as an example. It can be adapted to run on a different platform + accelerators + disaggregated inference solution combination.

```bash
# Download demo-magic.sh
curl -O https://raw.githubusercontent.com/paxtonhare/demo-magic/refs/heads/master/demo-magic.sh
# Run the demo on an AI conformant cluster with accelerators, GAIE, and disaggrated inference (adjust namespace and pod label selector if not using llm-d)
./run-demo-eu.sh
```

See it in action:

[<img src="https://img.youtube.com/vi/RRhPMuKPZG8/maxresdefault.jpg" width="560" height="315"
/>](https://www.youtube.com/embed/RRhPMuKPZG8?si=6niy8XUIE1aFvDjq)
