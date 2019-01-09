# Azure Container Registry Sample

Sample to deploy an Azure Container Registry

## Creates

1. A Resource Group
2. An [Azure Container Registry](https://azure.microsoft.com/en-us/services/container-registry/)

## Usage

- Provide values to all variables.
- Create with `terraform apply`
- Destroy all with `terraform destroy --force`

## Geo-Replication

Azure Container Registry supports [geo-replication](https://docs.microsoft.com/en-us/azure/container-registry/container-registry-geo-replication) to help you with multi-region deployment. It is only supported by `Premium` SKU.

To enable geo-replication with Terraform, just fill the `georeplication_locations` set with the list of Azure locations where you want the registry to be geo-replicated.