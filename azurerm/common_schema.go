package azurerm

// NOTE: methods in this file will be moved to `./helpers/azurerm` in time
// new methods should be added to this directory instead.
// This file exists to be able to consolidate files in the root

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

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
