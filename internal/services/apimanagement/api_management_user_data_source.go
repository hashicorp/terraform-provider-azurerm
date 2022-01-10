package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceApiManagementUser() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementUserRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"user_id": schemaz.SchemaApiManagementUserDataSourceName(),

			"api_management_name": schemaz.SchemaApiManagementDataSourceName(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"first_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"email": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"last_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"note": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementUserRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.UsersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewUserID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("user_id").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making read request %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if props := resp.UserContractProperties; props != nil {
		d.Set("first_name", props.FirstName)
		d.Set("last_name", props.LastName)
		d.Set("email", props.Email)
		d.Set("note", props.Note)
		d.Set("state", string(props.State))
	}

	return nil
}
