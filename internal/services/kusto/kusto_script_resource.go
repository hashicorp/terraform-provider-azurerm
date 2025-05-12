// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kusto

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kusto/2024-04-13/scripts"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceKustoDatabaseScript() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceKustoDatabaseScriptCreateUpdate,
		Read:   resourceKustoDatabaseScriptRead,
		Update: resourceKustoDatabaseScriptCreateUpdate,
		Delete: resourceKustoDatabaseScriptDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.KustoDatabaseScriptV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := scripts.ParseScriptID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"database_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateKustoDatabaseID,
			},

			"continue_on_errors_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  false,
			},

			"force_an_update_when_value_changed": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"url": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"url", "script_content"},
				RequiredWith: []string{"sas_token"},
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sas_token": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"url"},
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"script_content": {
				Type:         pluginsdk.TypeString,
				ExactlyOneOf: []string{"url", "script_content"},
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceKustoDatabaseScriptCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ScriptsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	databaseId, err := commonids.ParseKustoDatabaseID(d.Get("database_id").(string))
	if err != nil {
		return err
	}
	id := scripts.NewScriptID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.KustoClusterName, databaseId.KustoDatabaseName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %q: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_kusto_script", id.ID())
		}
	}

	clusterId := commonids.NewKustoClusterID(databaseId.SubscriptionId, databaseId.ResourceGroupName, databaseId.KustoClusterName)
	locks.ByID(clusterId.ID())
	defer locks.UnlockByID(clusterId.ID())

	forceUpdateTag := d.Get("force_an_update_when_value_changed").(string)
	if len(forceUpdateTag) == 0 {
		forceUpdateTag, _ = uuid.GenerateUUID()
	}

	parameters := scripts.Script{
		Properties: &scripts.ScriptProperties{
			ContinueOnErrors: utils.Bool(d.Get("continue_on_errors_enabled").(bool)),
			ForceUpdateTag:   utils.String(forceUpdateTag),
		},
	}

	if scriptURL, ok := d.GetOk("url"); ok {
		parameters.Properties.ScriptURL = utils.String(scriptURL.(string))
	}

	if scriptURLSasToken, ok := d.GetOk("sas_token"); ok {
		parameters.Properties.ScriptURLSasToken = utils.String(scriptURLSasToken.(string))
	}

	if scriptContent, ok := d.GetOk("script_content"); ok {
		parameters.Properties.ScriptContent = utils.String(scriptContent.(string))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %q: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceKustoDatabaseScriptRead(d, meta)
}

func resourceKustoDatabaseScriptRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ScriptsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scripts.ParseScriptID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	d.Set("name", id.ScriptName)
	d.Set("database_id", commonids.NewKustoDatabaseID(id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.DatabaseName).ID())

	if resp.Model != nil {
		if props := resp.Model.Properties; props != nil {
			d.Set("continue_on_errors_enabled", props.ContinueOnErrors)
			d.Set("force_an_update_when_value_changed", props.ForceUpdateTag)
			d.Set("url", props.ScriptURL)
		}
	}
	return nil
}

func resourceKustoDatabaseScriptDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Kusto.ScriptsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := scripts.ParseScriptID(d.Id())
	if err != nil {
		return err
	}

	// DELETE operation for script does not support running concurrently at cluster level
	locks.ByName(id.ClusterName, "azurerm_kusto_cluster")
	defer locks.UnlockByName(id.ClusterName, "azurerm_kusto_cluster")

	err = client.DeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %q: %+v", id, err)
	}

	return nil
}
