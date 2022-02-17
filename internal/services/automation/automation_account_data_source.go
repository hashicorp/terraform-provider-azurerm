package automation

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceAutomationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceAutomationAccountRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"primary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"secondary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	iclient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAutomationAccountID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retreiving %s: %+v", id, err)
	}
	d.SetId(id.ID())

	iresp, err := iclient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(iresp.Response) {
			return fmt.Errorf("%q Account Registration Information was not found", id)
		}
		return fmt.Errorf("retreiving Automation Account Registration Information %s: %+v", id, err)
	}
	if iresp.Keys != nil {
		d.Set("primary_key", iresp.Keys.Primary)
		d.Set("secondary_key", iresp.Keys.Secondary)
	}
	d.Set("endpoint", iresp.Endpoint)
	return nil
}
