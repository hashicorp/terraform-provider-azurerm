package disks

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithUpdate = DiskPoolResource{}

type DiskPoolResource struct {
}

func (DiskPoolResource) Arguments() map[string]*schema.Schema {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) Attributes() map[string]*schema.Schema {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) ModelObject() interface{} {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) ResourceType() string {
	return "azurerm_disk_pool"
}

func (DiskPoolResource) Create() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) Read() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) Delete() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	//TODO implement me
	panic("implement me")
}

func (DiskPoolResource) Update() sdk.ResourceFunc {
	//TODO implement me
	panic("implement me")
}
