# kube-watch

A Kubernetes platform toolkit for deploying, managing, and monitoring containerized workloads across clusters. Built with Go and Terraform.

## Overview

kube-watch provides three core components:

- **CLI** (`cli/`) — A Go command-line tool for managing the full lifecycle of Kubernetes workloads: deploying, tearing down, and monitoring applications across clusters.
- **Webapp** (`webapp/`) — A lightweight Go HTTP service used as the reference workload. Containerized with a multi-stage Docker build (Go builder to scratch) for minimal image size. Includes a health check endpoint for liveness probes.
- **Terraform** (`terraform/`) — Infrastructure-as-code for provisioning Kubernetes resources: namespaces, deployments with resource limits, services with NodePort exposure, and liveness probes.

## Project Structure

```
kube-watch/
├── cli/                  # Go CLI tool (deploy, destroy, monitor)
│   ├── main.go
│   └── go.mod
├── terraform/            # Kubernetes infrastructure-as-code
│   ├── main.tf           # Namespace, Deployment, Service definitions
│   ├── variables.tf      # Configurable inputs (image, replicas, ports)
│   ├── outputs.tf        # Post-apply output values and debug commands
│   └── terraform.tfvars  # Variable overrides
├── webapp/               # Go HTTP service (reference workload)
│   ├── main.go
│   ├── Dockerfile
│   └── go.mod
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker
- Terraform >= 1.0
- A local Kubernetes cluster (minikube or Docker Desktop)
- kubectl configured with cluster access

### Build and Deploy

1. Build the Docker image locally:
   ```sh
   cd webapp
   docker build -t go-webapp:latest .
   ```

2. Deploy to your local cluster with Terraform:
   ```sh
   cd terraform
   terraform init
   terraform apply
   ```

3. Access the service:
   ```sh
   curl http://localhost:30080
   ```

### CLI Usage

```sh
cd cli
go run main.go [deploy|destroy|monitor]
```

- `deploy` — Provisions infrastructure and deploys the webapp
- `destroy` — Tears down all deployed resources
- `monitor` — Watches cluster and workload health

## Current State

- Webapp is functional with health check endpoint and multi-stage Dockerfile
- Terraform provisions a namespace, deployment (2 replicas), and NodePort service on a local Kubernetes cluster
- Liveness probes are configured against the health endpoint
- Resource requests and limits are set per container
- CLI scaffolding is in place with command routing

## Roadmap

### CLI Completion
- Implement `deploy` command wrapping Terraform apply and image builds
- Implement `destroy` command for full teardown of resources
- Implement `monitor` command with live pod status, restart counts, and resource consumption
- Add structured logging and error handling

### Multi-Cluster Support
- Extend CLI and Terraform to target multiple Kubernetes contexts
- Support deploying and monitoring workloads across clusters simultaneously
- Add cluster-aware configuration management

### Multi-Cloud Infrastructure
- Modularize Terraform to support AWS EKS, Azure AKS, and GCP GKE
- Abstract provider-specific configuration behind a common interface
- Support provisioning clusters and deploying workloads across cloud providers

### Capacity and Resource Management
- Query and report node-level resource usage across clusters
- Track pod resource requests vs actual utilization per namespace
- Flag over-provisioned or under-provisioned workloads
- Surface capacity trends to assist with planning

### Observability
- Integrate Prometheus metrics collection
- Surface pod health, restart history, and resource metrics through the CLI
- Add alerting thresholds for workload anomalies
- Dashboard output for fleet-wide status

### Platform Tooling
- Namespace-level multi-tenancy and resource quota management
- Workload placement policies across clusters
- Self-service deployment workflows for application teams

## Technologies

- **Languages:** Go
- **Infrastructure:** Terraform, Docker
- **Platforms:** Kubernetes, AWS, Azure, GCP (planned)
