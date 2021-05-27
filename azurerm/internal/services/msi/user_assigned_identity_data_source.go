package msi

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/msi/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmUserAssignedIdentity() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceArmUserAssignedIdentityRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"principal_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"client_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmUserAssignedIdentityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSI.UserAssignedIdentitiesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewUserAssignedIdentityID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("User Assigned Identity %q was not found in Resource Group %q", id.Name, id.ResourceGroup)
		}
		return fmt.Errorf("retrieving User Assigned Identity %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.SetId(id.ID())

	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.UserAssignedIdentityProperties; props != nil {
		if principalId := props.PrincipalID; principalId != nil {
			d.Set("principal_id", principalId.String())
		}

		if clientId := props.ClientID; clientId != nil {
			d.Set("client_id", clientId.String())
		}

		if tenantId := props.TenantID; tenantId != nil {
			d.Set("tenant_id", tenantId.String())
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
