package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmApiManagementUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmApiManagementUserRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"user_id": azure.SchemaApiManagementUserDataSourceName(),

			"api_management_name": azure.SchemaApiManagementDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"first_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"note": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmApiManagementUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	userId := d.Get("user_id").(string)

	resp, err := client.Get(ctx, resourceGroup, serviceName, userId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("User %q was not found in API Management Service %q / Resource Group %q", userId, serviceName, resourceGroup)
		}

		return fmt.Errorf("Error making Read request on User %q (API Management Service %q / Resource Group %q): %+v", userId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	if props := resp.UserContractProperties; props != nil {
		d.Set("first_name", props.FirstName)
		d.Set("last_name", props.LastName)
		d.Set("email", props.Email)
		d.Set("note", props.Note)
		d.Set("state", string(props.State))
	}

	return nil
}
