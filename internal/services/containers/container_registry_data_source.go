package containers

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceContainerRegistry() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceContainerRegistryRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: dataSourceContainerRegistrySchema(),
	}
}

func dataSourceContainerRegistryRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Containers.RegistriesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewRegistryID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making Read request on %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("admin_enabled", resp.AdminUserEnabled)
	d.Set("login_server", resp.LoginServer)

	if sku := resp.Sku; sku != nil {
		d.Set("sku", string(sku.Tier))
	}

	// Deprecated as it is not returned by the API now.
	d.Set("storage_account_id", "")

	if *resp.AdminUserEnabled {
		credsResp, err := client.ListCredentials(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("making Read request on %s: %s", id, err)
		}

		d.Set("admin_username", credsResp.Username)
		for _, v := range *credsResp.Passwords {
			d.Set("admin_password", v.Value)
			break
		}
	} else {
		d.Set("admin_username", "")
		d.Set("admin_password", "")
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func dataSourceContainerRegistrySchema() map[string]*pluginsdk.Schema {
	out := map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validate.ContainerRegistryName,
		},

		"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

		"location": azure.SchemaLocationForDataSource(),

		"admin_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"admin_password": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"admin_username": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"login_server": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"sku": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"tags": tags.SchemaDataSource(),
	}
	if !features.ThreePointOhBeta() {
		out["storage_account_id"] = &pluginsdk.Schema{
			Type:       pluginsdk.TypeString,
			Computed:   true,
			Deprecated: "this attribute is no longer recognized by the API and is not functional anymore, thus this property will be removed in v3.0",
		}
	}

	return out
}
