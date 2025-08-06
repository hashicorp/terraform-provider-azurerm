# Guide: When to inline new functionality (as either a block or property) versus a new resource

Sometimes when implementing new functionality it can be a bit unclear whether it is necessary to create a new resource versus to add a new property or block to an existing resource.

To get a bit of insight in how a decision can be made, these are some rules of thumb to decide. In case it is unclear, please contact one of the HashiCorp Maintainers.

## Inline
Most additional functionality will end up inline in an existing resource.

APIs to enable or disable functionality, define the resource functionality or configure details on a resource are most of the time inlined. Relations between resources with clearly separate concern (i.e. which VNet a K8s cluster will land in) are most of the time inlined.

A few categories of inlined functionality with possible motivations to inline are summed up below.

### Category 1: properties
- When it is 'just a property' of this resource, like `sku` in the example below.
- It would require a lot of extra work to make these separate resources.

```hcl
resource "azurerm_example_resource" "example" {
  name = "ThePerfectExample"
  sku  = "Gold"
}
```

### Category 2: child resources which cannot be separated
- It has a strict `1:1` relation with its parent resource
- It cannot be deleted, only returned to a default state (i.e. you might find an API, but only to update the resource)
- It doesn't have its own unique Resource ID or name (i.e. `<parentId>/default` or `<parentId>/keyrotationpolicy`, not something like `<parentId>/subResource/MySubResourceName`)
- It does not have its own API endpoint but uses the parent resource endpoint
- It does not cross security/team boundaries in the most common client situation
- It does not contain backwards compatibility issues

### Category 3: relations between resources
- Resources are really separate and have `1:many` and `many:1` relations with a (i.e. a relation property with a resource within a completely different Resource Provider `/subscriptions/<subscriptionId>/resourceGroups/<resourceGroup>/providers/Microsoft.Network/networkSecurityGroups/example-nsg` and `/subscriptions/<subscriptionId>/resourceGroups/<resourceGroup>/providers/Microsoft.Storage/storageAccounts/example-storage`: which `azurerm_storage_account` is used to store the logs from a resource).
- These relations are created with API calls to the original resource provider, not the connected one
- Reading the connections between the resources does not require extra permissions than previously necessary to create the resource (i.e. the Service Principal used to read the connection/relation should be also `Owner`/`Contributor` on the resource group the connection is made to)

## Separate resource

While inlining might make a lot of sense for many APIs, there are also good reasons to separate them out. These arguments may not be conclusive, but can help steer in the right direction.

### Category 1: the obvious new resource
- It is a new resource with its own lifecycle, own API endpoints (at least `Update` and `Delete`), it just feels natural to put it in a separate resource

### Category 2: child resources
- It does have its own unique Resource ID or name (i.e. `<parentId>/subResource/MySubResourceName`)
- It has its own API endpoint for `Create`, `Update` and `Delete` actions
- It needs more permissions on the already existing resource than the current parent resource requires (i.e. Key Vault Key Rotation Policies require more permissions than the Key Vault Key it really belongs to)
- Control Plane vs Data plane: the functionality is acting on the Data Plane instead of Control Plane of the service or vice versa. (i.e. Azure Storage Account management vs the actual Blobs put in there, Azure Key Vault management vs the Keys/Certs/Secrets inside)
- Its functionality and therefore the scope of the resource crosses team/security boundaries (i.e. Infra team vs Application team).

### Category 3: relations between resources (_"It is complicated"_)
- It is a mediator: there is a separate endpoint to create a relation between two existing resources
- It requires more permissions on another resource to create the connection than to create the resource itself (i.e. connecting a NSG resource to a Subnet)

## Both inline and separate
It might be that there are multiple use-cases and scenarios necessary. Sometimes it makes sense to create it inline, sometimes it makes more sense to separate them.

This requires caution from both the implementer and the user. In most cases it should be explained with some notes in the `docs`. Within the inline resource implementation it requires that it doesn't delete or update properties created externally when it is not explicitly configured in the resource. For the user this might have the drawback that the inlined resource is not strict in enforcing the existing inlined properties. Mixed use within the same context might end up in a mess and is advised not to do.

A few examples of resources which are both inlined and separate resources:
- Subnets (part of VNet resource as well)
- NSG rules (part of NSG resource as well)
- Key Vault permissions (part of Key Vault resource as well)