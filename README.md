# Kubernetes cluster-api-provider-virtlet

## WARNING

This is a research project. Do not try to use it in production.

## Getting Started

### Prerequisites

1. Prepare a Kubernetes cluster with [Virtlet](https://github.com/Mirantis/virtlet) installed.
2. Deploy cluster-api and provider `kubectl apply -f provider-components.yaml`

### Cluster Creation

Create cluster and machines:

```bash

kubectl apply -f hack/examples/cluster.yml
kubectl apply -f hack/examples/master-machine.yml
kubectl apply -f hack/examples/machine.yml
```
