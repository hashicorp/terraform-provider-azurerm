# AKS Node Pool Properties Implementation Plan

## Overview
This document outlines the plan to implement missing AKS node pool properties across separate git branches.

## Properties to Implement

### 1. Gateway Profile (`gateway_profile`)
**Branch:** `feature/aks-node-pool-gateway-profile`
**Status:** In Progress
**Files to modify:**
- ✅ Schema: `internal/services/containers/kubernetes_cluster_node_pool_resource.go`
- ✅ Expand/Flatten functions: Added
- ⏳ Create function: Need to add read/update support
- ⏳ Tests: `internal/services/containers/kubernetes_cluster_node_pool_resource_test.go`
- ⏳ Documentation: `website/docs/r/kubernetes_cluster_node_pool.html.markdown`

### 2. Pod IP Allocation Mode (`pod_ip_allocation_mode`)
**Branch:** `feature/aks-node-pool-pod-ip-allocation-mode`
**Status:** Not Started
**Files to modify:**
- Schema: `internal/services/containers/kubernetes_cluster_node_pool_resource.go`
- Create/Read/Update functions
- Tests: `internal/services/containers/kubernetes_cluster_node_pool_resource_test.go`
- Documentation: `website/docs/r/kubernetes_cluster_node_pool.html.markdown`

### 3. Security Profile (`security_profile`)
**Branch:** `feature/aks-node-pool-security-profile`
**Status:** Not Started
**Files to modify:**
- Schema: `internal/services/containers/kubernetes_cluster_node_pool_resource.go`
- Expand/Flatten functions
- Create/Read/Update functions
- Tests: `internal/services/containers/kubernetes_cluster_node_pool_resource_test.go`
- Documentation: `website/docs/r/kubernetes_cluster_node_pool.html.markdown`

### 4. Virtual Machine Nodes Status (`virtual_machine_nodes_status`)
**Branch:** `feature/aks-node-pool-vm-nodes-status`
**Status:** Not Started
**Files to modify:**
- Schema: `internal/services/containers/kubernetes_cluster_node_pool_resource.go` (Computed only)
- Flatten function only (no expand needed - computed property)
- Read function
- Tests: `internal/services/containers/kubernetes_cluster_node_pool_resource_test.go`
- Documentation: `website/docs/r/kubernetes_cluster_node_pool.html.markdown`

### 5. Virtual Machine Profile (`virtual_machine_profile`)
**Branch:** `feature/aks-node-pool-vm-profile`
**Status:** Not Started
**Files to modify:**
- Schema: `internal/services/containers/kubernetes_cluster_node_pool_resource.go`
- Expand/Flatten functions
- Create/Read/Update functions
- Tests: `internal/services/containers/kubernetes_cluster_node_pool_resource_test.go`
- Documentation: `website/docs/r/kubernetes_cluster_node_pool.html.markdown`

## Implementation Checklist for Each Property

For each property, complete the following steps:

1. **Code Implementation**
   - [ ] Add schema definition
   - [ ] Add expand function (if applicable)
   - [ ] Add flatten function
   - [ ] Update create function
   - [ ] Update read function
   - [ ] Update update function (if applicable)
   - [ ] Run `make fmt`
   - [ ] Run `go build ./internal/services/containers/`

2. **Testing**
   - [ ] Add basic test case
   - [ ] Add update test case (if applicable)
   - [ ] Run `make test`
   - [ ] Run acceptance tests (optional due to quota issues)

3. **Documentation**
   - [ ] Add property to arguments reference
   - [ ] Add property to attributes reference (if computed)
   - [ ] Add example usage
   - [ ] Document any special behavior or requirements

4. **Commit and Branch**
   - [ ] Commit all changes
   - [ ] Push branch to remote
   - [ ] Create PR with proper description

## Key Files

- **Resource Implementation:** `internal/services/containers/kubernetes_cluster_node_pool_resource.go`
- **Tests:** `internal/services/containers/kubernetes_cluster_node_pool_resource_test.go`
- **Documentation:** `website/docs/r/kubernetes_cluster_node_pool.html.markdown`
- **API Models:** `vendor/github.com/hashicorp/go-azure-sdk/resource-manager/containerservice/2025-07-01/agentpools/`

## Next Steps

1. Complete the gateway_profile implementation (current branch)
2. Create remaining branches
3. Implement each property individually
4. Test and document each property
5. Create PRs for review

