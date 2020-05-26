package msi

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmUserAssignedIdentity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmUserAssignedIdentityRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 128),
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"location": azure.SchemaLocationForDataSource(),

			"principal_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"client_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),
		},
	}
}

func dataSourceArmUserAssignedIdentityRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSI.UserAssignedIdentitiesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("User Assigned Identity %q was not found in Resource Group %q", name, resourceGroup)
		}
		return fmt.Errorf("Error making Read request on User Assigned Identity %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.IdentityProperties; props != nil {
		if principalId := props.PrincipalID; principalId != nil {
			d.Set("principal_id", principalId.String())
		}

		if clientId := props.ClientID; clientId != nil {
			d.Set("client_id", clientId.String())
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
