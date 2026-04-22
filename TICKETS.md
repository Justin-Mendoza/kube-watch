# kube-watch — Tickets

---

## KW-001: Local Environment Setup
**Priority:** Critical
**Status:** To Do
**Component:** Infrastructure

**Description:**
Set up the local development environment to get the full deploy loop working end-to-end.

**Acceptance Criteria:**
- [ ] Docker Desktop installed and running
- [ ] minikube installed and cluster started
- [ ] Webapp Docker image builds successfully
- [ ] `terraform apply` provisions namespace, deployment, and service
- [ ] Webapp is accessible in browser via `minikube service`
- [ ] Liveness probe passes on `/health` endpoint

---

## KW-002: Fix CLI Broken Imports
**Priority:** Critical
**Status:** To Do
**Component:** CLI
**Depends On:** KW-001

**Description:**
The CLI has a bad import (`"fmt/log"`) that will prevent compilation. Clean up the import block so the project builds.

**Acceptance Criteria:**
- [ ] Remove invalid `"fmt/log"` import
- [ ] Remove any unused imports
- [ ] `go build` compiles successfully

---

## KW-003: Implement CLI `deploy` Command
**Priority:** High
**Status:** To Do
**Component:** CLI
**Depends On:** KW-002

**Description:**
Implement the `deploy` subcommand that automates the full deployment pipeline: building the Docker image and running Terraform apply.

**Acceptance Criteria:**
- [ ] Runs `docker build -t go-webapp:latest ./webapp`
- [ ] Runs `terraform init` + `terraform apply -auto-approve` in the terraform directory
- [ ] Streams command output to the terminal in real time
- [ ] Exits with a clear error message if any step fails
- [ ] Prints the service URL on success

---

## KW-004: Implement CLI `destroy` Command
**Priority:** High
**Status:** To Do
**Component:** CLI
**Depends On:** KW-002

**Description:**
Implement the `destroy` subcommand to tear down all Kubernetes resources created by `deploy`.

**Acceptance Criteria:**
- [ ] Runs `terraform destroy -auto-approve` in the terraform directory
- [ ] Streams output to the terminal
- [ ] Confirms resources are removed on completion
- [ ] Handles case where nothing is deployed gracefully

---

## KW-005: Implement CLI `monitor` Command
**Priority:** High
**Status:** To Do
**Component:** CLI
**Depends On:** KW-003

**Description:**
Implement the `monitor` subcommand using the Kubernetes Go client (`k8s.io/client-go`) to display live workload health.

**Acceptance Criteria:**
- [ ] Connects to cluster via kubeconfig
- [ ] Lists all pods in `local-go-app` namespace with status
- [ ] Shows restart counts per pod
- [ ] Shows CPU and memory requests vs limits per container
- [ ] Refreshes on an interval (default 5s)
- [ ] Clean terminal output (table format or similar)

---

## KW-006: Modularize Terraform for Multi-Cluster
**Priority:** Medium
**Status:** To Do
**Component:** Terraform
**Depends On:** KW-001

**Description:**
Refactor Terraform configuration into reusable modules so the same infrastructure can target multiple Kubernetes cluster contexts.

**Acceptance Criteria:**
- [ ] Extract deployment, service, and namespace into a module
- [ ] Module accepts `kube_context` as an input variable
- [ ] Can apply the same module against two different clusters (e.g., minikube + kind)
- [ ] Each cluster gets its own state file

---

## KW-007: CLI Multi-Cluster Support
**Priority:** Medium
**Status:** To Do
**Component:** CLI
**Depends On:** KW-005, KW-006

**Description:**
Extend the CLI to support deploying and monitoring across multiple Kubernetes clusters.

**Acceptance Criteria:**
- [ ] Add `--context` flag to target a specific cluster
- [ ] Add `--all` flag to run against all configured clusters
- [ ] `monitor` shows aggregated status across clusters
- [ ] Config file or flag to define list of cluster contexts

---

## KW-008: Multi-Cloud Terraform Modules
**Priority:** Medium
**Status:** To Do
**Component:** Terraform
**Depends On:** KW-006

**Description:**
Add Terraform modules for provisioning managed Kubernetes clusters on AWS EKS, Azure AKS, and GCP GKE behind a common interface.

**Acceptance Criteria:**
- [ ] `modules/eks/` provisions an EKS cluster
- [ ] `modules/aks/` provisions an AKS cluster
- [ ] `modules/gke/` provisions a GKE cluster
- [ ] Common variables interface across all three (cluster name, node count, machine type)
- [ ] Outputs cluster kubeconfig for each provider

---

## KW-009: Capacity and Resource Reporting
**Priority:** Medium
**Status:** To Do
**Component:** CLI
**Depends On:** KW-005

**Description:**
Add resource utilization reporting to the CLI — query node and pod level resource usage to assist with capacity planning.

**Acceptance Criteria:**
- [ ] New `capacity` subcommand
- [ ] Shows per-node: allocatable vs requested CPU and memory
- [ ] Shows per-namespace: total resource requests vs limits
- [ ] Flags over-provisioned workloads (requests < 30% of limits)
- [ ] Flags under-provisioned workloads (requests > 90% of allocatable)
- [ ] Summary output with cluster-wide utilization percentage

---

## KW-010: Prometheus Integration
**Priority:** Low
**Status:** To Do
**Component:** Observability
**Depends On:** KW-005

**Description:**
Deploy Prometheus to the cluster and integrate real metrics into the CLI and webapp dashboard.

**Acceptance Criteria:**
- [ ] Prometheus deployed via Helm chart in Terraform
- [ ] Webapp exposes `/metrics` endpoint with basic Go runtime metrics
- [ ] `monitor` command can pull real CPU/memory metrics from Prometheus
- [ ] Webapp dashboard displays real metric values instead of hardcoded counters

---

## KW-011: Workload Placement Policies
**Priority:** Low
**Status:** To Do
**Component:** Platform
**Depends On:** KW-007, KW-009

**Description:**
Implement workload placement logic that decides which cluster a workload should be deployed to based on available capacity.

**Acceptance Criteria:**
- [ ] CLI evaluates cluster capacity before deploying
- [ ] Selects cluster with most available resources by default
- [ ] Supports affinity rules (e.g., prefer a specific cloud provider)
- [ ] Dry-run mode that shows placement decision without deploying

---

## KW-012: Namespace Multi-Tenancy and Quotas
**Priority:** Low
**Status:** To Do
**Component:** Platform
**Depends On:** KW-006

**Description:**
Add support for namespace-level resource quotas and multi-tenancy to simulate platform support for multiple teams.

**Acceptance Criteria:**
- [ ] Terraform creates namespaces with ResourceQuota objects
- [ ] Configurable CPU and memory limits per namespace
- [ ] CLI `capacity` command shows quota usage per namespace
- [ ] Warn when a namespace is approaching its quota limit
