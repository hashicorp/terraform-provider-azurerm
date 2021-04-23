package springcloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
)

const (
	springCloudAppAssociationTypeCosmosDb = "Microsoft.DocumentDB"
	springCloudAppAssociationTypeMysql    = "Microsoft.DBforMySQL"
	springCloudAppAssociationTypeRedis    = "Microsoft.Cache"
)

func importSpringCloudAppAssociation(resourceType string) func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
	return func(d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.SpringCloudAppAssociationID(d.Id())
		if err != nil {
			return []*schema.ResourceData{}, err
		}

		client := meta.(*clients.Client).AppPlatform.BindingsClient
		ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
		defer cancel()

		resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
		if err != nil {
			return []*schema.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Properties == nil || resp.Properties.ResourceType == nil {
			return []*schema.ResourceData{}, fmt.Errorf("retrieving %s: `properties` or `properties.resourceType` was nil", id)
		}

		if *resp.Properties.ResourceType != resourceType {
			return []*schema.ResourceData{}, fmt.Errorf(`spring Cloud App Association "type" mismatch, expected "%s", got "%s"`, resourceType, *resp.Properties.ResourceType)
		}

		return []*schema.ResourceData{d}, nil
	}
}
