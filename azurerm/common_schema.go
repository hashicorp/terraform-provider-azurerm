package azurerm

// NOTE: methods in this file will be moved to `./helpers/azurerm` in time
// new methods should be added to this directory instead.
// This file exists to be able to consolidate files in the root

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func locationSchema() *schema.Schema {
	return azure.SchemaLocation()
}

func locationForDataSourceSchema() *schema.Schema {
	return azure.SchemaLocationForDataSource()
}

func deprecatedLocationSchema() *schema.Schema {
	return azure.SchemaLocationDeprecated()
}

func azureRMNormalizeLocation(location interface{}) string {
	return azure.NormalizeLocation(location.(string))
}

func azureRMSuppressLocationDiff(k, old, new string, d *schema.ResourceData) bool {
	return azure.SuppressLocationDiff(k, old, new, d)
}

func resourceGroupNameSchema() *schema.Schema {
	return azure.SchemaResourceGroupName()
}

func resourceGroupNameDiffSuppressSchema() *schema.Schema {
	return azure.SchemaResourceGroupNameDiffSuppress()
}

func resourceGroupNameForDataSourceSchema() *schema.Schema {
	return azure.SchemaResourceGroupNameForDataSource()
}

func zonesSchema() *schema.Schema {
	return azure.SchemaZones()
}

func singleZonesSchema() *schema.Schema {
	return azure.SchemaSingleZone()
}

func zonesSchemaComputed() *schema.Schema {
	return azure.SchemaZonesComputed()
}

func expandZones(v []interface{}) *[]string {
	return azure.ExpandZones(v)
}

func azureRMHashLocation(location interface{}) int {
	return azure.HashAzureLocation(location)
}
