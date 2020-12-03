package example

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
)

type ExampleResource struct {
}

func (r ExampleResource) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"number": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
}

func (r ExampleResource) Attributes() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"output": {
			Type: schema.TypeString,
		},
	}
}

func (r ExampleResource) ResourceType() string {
	return "azurerm_example"
}

// NOTE: i guess we could return schema object to ensure everything is mapped and valid idk
type ExampleObj struct {
	Name   string `hcl:"name"`
	Number int    `hcl:"number"`
}

func (r ExampleResource) Create() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			//metadata.ResourceData
			//metadata.Logger.WarnF("OHHAI %d", 3)
			//metadata.Client.Account.SubscriptionId
			metadata.Logger.Info("HEYO")

			var obj ExampleObj
			if err := metadata.Decode(&obj); err != nil {
				return err
			}

			id := parse.SubnetId{
				ResourceGroup:      "production-resources",
				VirtualNetworkName: "production-network",
				Name:               obj.Name,
			}

			metadata.Logger.InfoF("Name is %s", obj.Name)
			metadata.Logger.InfoF("Number is %d", obj.Number)

			metadata.SetID(id)
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ExampleResource) Read() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			return metadata.Encode(&ExampleObj{
				Name:   "updated",
				Number: 123,
			})
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ExampleResource) Update() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ExampleResource) Delete() ResourceFunc {
	return ResourceFunc{
		Func: func(ctx context.Context, metadata ResourceMetaData) error {
			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ExampleResource) IDValidationFunc() schema.SchemaValidateFunc {
	return validate.SubnetID
}
