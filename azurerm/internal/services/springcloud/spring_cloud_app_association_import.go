package springcloud

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/springcloud/parse"
)

const (
	springCloudAppAssociationTypeCosmosDb = "Microsoft.DocumentDB"
	springCloudAppAssociationTypeMysql    = "Microsoft.DBforMySQL"
	springCloudAppAssociationTypeRedis    = "Microsoft.Cache"
)

func importSpringCloudAppAssociation(resourceType string) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) (data []*schema.ResourceData, err error) {
		id, err := parse.SpringCloudAppAssociationID(d.Id())
		if err != nil {
			return []*pluginsdk.ResourceData{}, err
		}

		client := meta.(*clients.Client).AppPlatform.BindingsClient
		resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.AppName, id.BindingName)
		if err != nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Properties == nil || resp.Properties.ResourceType == nil {
			return []*pluginsdk.ResourceData{}, fmt.Errorf("retrieving %s: `properties` or `properties.resourceType` was nil", id)
		}

		if *resp.Properties.ResourceType != resourceType {
			return []*pluginsdk.ResourceData{}, fmt.Errorf(`spring Cloud App Association "type" mismatch, expected "%s", got "%s"`, resourceType, *resp.Properties.ResourceType)
		}

		return []*pluginsdk.ResourceData{d}, nil
	}
}
