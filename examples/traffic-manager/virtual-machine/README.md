## Example: Traffic Manager fronting Virtual Machines

This example provisions a Traffic Manager Profile which fronts a couple of Virtual Machines.

This Terraform template was based on [this](https://github.com/Azure/azure-quickstart-templates/tree/master/201-traffic-manager-vm) Azure Quickstart Template. Changes to the ARM template that may have occurred since the creation of this example may not be reflected here.

This template shows how to create an Azure Traffic Manager profile to load-balance across a couple of Azure virtual machines. Each endpoint has an equal weight but different weights can be specified to distribute load non-uniformly.

See also:

- [Traffic Manager routing methods for details of the different routing methods available.](https://azure.microsoft.com/en-us/documentation/articles/traffic-manager-routing-methods/)
- [Create or update a Traffic Manager profile](https://msdn.microsoft.com/en-us/library/azure/mt163581.aspx) for details of the JSON elements relating to a Traffic Manager profile.
