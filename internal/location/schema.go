package location

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func Schema() *pluginsdk.Schema {
	return commonschema.Location()
}

func StateFunc(input interface{}) string {
	return location.StateFunc(input)
}
