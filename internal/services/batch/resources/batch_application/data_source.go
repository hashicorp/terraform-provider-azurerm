// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package batch_application

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/batch/2024-07-01/application"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	batchClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Use this data source to access information about an existing Batch Application instance.",
		Read:        dataSourceBatchApplicationRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The name of the Application.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ApplicationName,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"account_name": {
				Description:  "The name of the Batch Account.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.AccountName,
			},

			"allow_updates": {
				Description: "Whether packages within the application may be overwritten using the same version string.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"default_version": {
				Description: "The package to use if a client requests the application but does not specify a version.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"display_name": {
				Description: "The display name for the application.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceBatchApplicationRead(d *pluginsdk.ResourceData, meta any) error {
	client := batchClient.From(meta).ApplicationClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	subscriptionId := account.From(meta).SubscriptionId
	id := application.NewApplicationID(subscriptionId, d.Get("resource_group_name").(string), d.Get("account_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	d.SetId(id.ID())

	d.Set("name", id.ApplicationName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("account_name", id.BatchAccountName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("allow_updates", props.AllowUpdates)
			d.Set("default_version", props.DefaultVersion)
			d.Set("display_name", props.DisplayName)
		}
	}

	return nil
}
