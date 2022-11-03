## When to choose for an inline property (block) vs a new resource
Sometimes it might be a bit unclear whether it is necessary to create a new resource vs add a new property or property block to an existing resource.

To get a bit of insight in how the decission will be made, this are some rules of thumb to decide. In case it is unclear, please contact one of the HashiCorp Maintainers.

### Inline
Most extra functionality will end up inline an existing resource. Properties to enable or disable functionality with maybe some extra configuration on how it works on this specific resource like schedules, relations with clearly separate resources (i.e. which VNet a K8s cluster will land in) are most of the time inlined.

#### Category 1: properties
- When it is 'just a property' of this resource

#### Category 2: child resources which cannot be sparated
- It has a strict `1:1` relation with it's parent resource
- It cannot be deleted, only returned to a default state
- It doesn't have it's own unique Resource ID or name (i.e. `<parentId>/default` or `<parentId>/keyrotationpolicy`, not something like `<parentId>/subResource/MySubResourceName`)
- It does not have it's own API endpoint but uses the parent resource endpoint
- It does not cross security/team boundaries in the most common client situation
- It does not contain backwards compatibility issues

#### Category 3: relations between resources
- Resources are really separate and have `1:many` and `many:1` relations with a (i.e. a relation property with a resource within a completely different Resource Provider `/subscriptions/<subscriptionId>/resourceGroups/<resourceGroup>/providers/Microsoft.Network/networkSecurityGroups/example-nsg` and `/subscriptions/<subscriptionId>/resourceGroups/<resourceGroup>/providers/Microsoft.Storage/storageAccounts/example-storage`: which `azurerm_storage_account` is used to store the logs from a resource).
- Connecting these resources does not require extensive permissions like `Owner` on the other resource.

### Separate resource
- It is a new resource with it's own lifecycle
- It needs more permissions on the already existing resource than the current parent resource requires (i.e. Key Vault Key Rotation Policies require more permissions than the Key Vault Key it really belongs to)
- It crosses team/security boundaries (i.e. Infra team vs Application team)

### Both inline and separate
- It depends on the use-case which resource comes in handy, i.e. VNet with inline Subnets vs separate, NSG with rules inlined and separated