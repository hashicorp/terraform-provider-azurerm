# AKS Egress Gateway Example

This example demonstrates how to configure Azure Kubernetes Service (AKS) with the Egress Gateway feature.

## What this example does

This configuration creates:

1. **AKS Cluster with Static Egress Gateway Profile**: Enables the static egress gateway feature at the cluster level
2. **Gateway Node Pool**: A dedicated node pool with mode set to "Gateway" for hosting egress gateways
3. **User Node Pool**: An additional node pool for regular application workloads

## Egress Gateway Feature Overview

The AKS Egress Gateway feature provides:

- **Static IP for egress traffic**: All outbound traffic from the cluster uses a consistent, predictable IP address
- **Dedicated gateway nodes**: Specialized nodes handle egress traffic routing
- **Improved security**: Better control over outbound traffic and firewall rules

## Configuration Details

### Static Egress Gateway Profile

```hcl
network_profile {
  network_plugin = "azure"
  static_egress_gateway_profile {
    enabled = true
  }
}
```

This enables the static egress gateway feature at the cluster level.

### Gateway Node Pool

```hcl
resource "azurerm_kubernetes_cluster_node_pool" "egress_gateway" {
  mode = "Gateway"

  gateway_profile {
    public_ip_prefix_size = 30  # Range: 28-31, default: 31
  }
  # ... other configuration
}
```

Node pools with `mode = "Gateway"` are specifically configured to host the egress gateway infrastructure. The `gateway_profile` block allows you to configure the public IP prefix size (28-31) that determines the number of IP addresses allocated to the gateway nodes.

## Prerequisites

- Azure subscription with appropriate permissions
- Terraform >= 1.0
- Azure CLI (optional, for post-deployment verification)

## Usage

1. Clone this repository and navigate to this example directory
2. Run `terraform init` to initialize the configuration
3. Run `terraform plan` to review the planned changes
4. Run `terraform apply` to create the resources

## Verification

After deployment, you can verify the egress gateway configuration:

```bash
# Get cluster credentials
az aks get-credentials --resource-group example-k8s-egress-gateway --name example-aks

# Check node pools
kubectl get nodes --show-labels

# Verify egress gateway pods
kubectl get pods -n kube-system | grep egress-gateway
```

## Cleanup

To remove all resources:

```bash
terraform destroy
```

## Learn More

- [Azure AKS Egress Gateway Documentation](https://learn.microsoft.com/en-us/azure/aks/configure-static-egress-gateway)
- [Azure REST API Documentation](https://learn.microsoft.com/en-us/rest/api/aks/managed-clusters/create-or-update)