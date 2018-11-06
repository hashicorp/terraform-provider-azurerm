package azurerm

// NOTE: methods in this file will be moved to `./helpers/azurerm` in time
// new methods should be added to this directory instead.
// This file exists to be able to consolidate files in the root

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceGroupNameSchema() *schema.Schema {
	return azure.SchemaResourceGroupName()
}

func resourceGroupNameDiffSuppressSchema() *schema.Schema {
	return azure.SchemaResourceGroupNameDiffSuppress()
}

func resourceGroupNameForDataSourceSchema() *schema.Schema {
	return azure.SchemaResourceGroupNameForDataSource()
}
