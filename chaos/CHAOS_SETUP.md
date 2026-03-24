# Chaos Engineering Setup

This document describes the chaos engineering experiments and tooling for the Phenotype platform.

## Overview

We use LitmusChaos for Kubernetes-based chaos experiments. The experiments are defined as Kubernetes custom resources.

## Prerequisites

- Kubernetes cluster 1.19+
- LitmusChaos operator installed
- Service account with appropriate RBAC

## Installation

```bash
# Install LitmusChaos
kubectl apply -f https://litmuschaos.github.io/litmus-operator/v2.14.0/litmus-operator-v2.14.0.yaml

# Verify installation
kubectl get pods -n litmus
```

## Experiments

### 1. Pod Failure

Simulates a pod being killed unexpectedly.

```yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: pod-failure-chaos
  namespace: default
spec:
  appinfo:
    appns: default
    applabel: "app=phenotype-api"
  chaosServiceAccount: litmus-admin
  experiments:
    - name: pod-delete
      spec:
        components:
          env:
            - name: TOTAL_CHAOS_DURATION
              value: '30'
            - name: CHAOS_INTERVAL
              value: '10'
            - name: FORCE
              value: 'false'
```

### 2. Network Latency

Injects network latency between services.

```yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: network-latency-chaos
  namespace: default
spec:
  appinfo:
    appns: default
    applabel: "app=phenotype-api"
  chaosServiceAccount: litmus-admin
  experiments:
    - name: container-kill
      spec:
        components:
          env:
            - name: TOTAL_CHAOS_DURATION
              value: '60'
            - name: NETWORK_LATENCY
              value: '1000'
            - name: JITTER
              value: '100'
```

### 3. CPU Stress

Generates CPU load on target pods.

```yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: cpu-stress-chaos
  namespace: default
spec:
  appinfo:
    appns: default
    applabel: "app=phenotype-api"
  chaosServiceAccount: litmus-admin
  experiments:
    - name: pod-cpu-hog
      spec:
        components:
          env:
            - name: TOTAL_CHAOS_DURATION
              value: '30'
            - name: CPU_CORE
              value: '1'
            - name: CPU_LOAD
              value: '50'
```

### 4. Memory Stress

Consumes memory on target pods.

```yaml
apiVersion: litmuschaos.io/v1alpha1
kind: ChaosEngine
metadata:
  name: memory-stress-chaos
  namespace: default
spec:
  appinfo:
    appns: default
    applabel: "app=phenotype-api"
  chaosServiceAccount: litmus-admin
  experiments:
    - name: pod-memory-hog
      spec:
        components:
          env:
            - name: TOTAL_CHAOS_DURATION
              value: '30'
            - name: MEMORY_CONSUMPTION
              value: '500'
```

## Running Experiments

```bash
# Apply chaos experiment
kubectl apply -f chaos/pod-failure.yaml

# Monitor experiment
kubectl describe chaosengine pod-failure-chaos

# Check chaos results
kubectl get chaosexperiments
kubectl get chaosresults
```

## Guardrails

- **Run during maintenance windows**: Chaos experiments should be scheduled during low-traffic periods
- **Set limits**: Define maximum blast radius (e.g., max 10% of pods)
- **Monitor**: Ensure alerting is active during experiments
- **Rollback**: Have automatic rollback mechanisms in place
- **Communication**: Notify team before running experiments

## Automation

Run chaos experiments as part of CI/CD:

```bash
#!/bin/bash
# chaos-runner.sh

# Run smoke tests first
./scripts/smoke-tests.sh

# Run chaos experiment
kubectl apply -f chaos/pod-failure.yaml

# Wait for experiment completion
sleep 60

# Verify application health
curl -f http://health endpoint || exit 1

# Run regression tests
./scripts/regression-tests.sh
```
