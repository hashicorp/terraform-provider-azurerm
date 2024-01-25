// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package monitor

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/insights/2021-07-01-preview/privatelinkscopesapis"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMonitorPrivateLinkScope() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMonitorPrivateLinkScopeCreateUpdate,
		Read:   resourceMonitorPrivateLinkScopeRead,
		Update: resourceMonitorPrivateLinkScopeCreateUpdate,
		Delete: resourceMonitorPrivateLinkScopeDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privatelinkscopesapis.ParsePrivateLinkScopeID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateLinkScopeName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"access_mode": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ingestion": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(privatelinkscopesapis.AccessModePrivateOnly),
							ValidateFunc: validation.StringInSlice([]string{
								string(privatelinkscopesapis.AccessModeOpen),
								string(privatelinkscopesapis.AccessModePrivateOnly),
							}, false),
						},

						"query": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  string(privatelinkscopesapis.AccessModePrivateOnly),
							ValidateFunc: validation.StringInSlice([]string{
								string(privatelinkscopesapis.AccessModeOpen),
								string(privatelinkscopesapis.AccessModePrivateOnly),
							}, false),
						},

						"exclusions": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"connection_name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},

									"ingestion": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(privatelinkscopesapis.AccessModePrivateOnly),
										ValidateFunc: validation.StringInSlice([]string{
											string(privatelinkscopesapis.AccessModeOpen),
											string(privatelinkscopesapis.AccessModePrivateOnly),
										}, false),
									},

									"query": {
										Type:     pluginsdk.TypeString,
										Optional: true,
										Default:  string(privatelinkscopesapis.AccessModePrivateOnly),
										ValidateFunc: validation.StringInSlice([]string{
											string(privatelinkscopesapis.AccessModeOpen),
											string(privatelinkscopesapis.AccessModePrivateOnly),
										}, false),
									},
								},
							},
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMonitorPrivateLinkScopeCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := privatelinkscopesapis.NewPrivateLinkScopeID(subscriptionId, resourceGroup, name)

	if d.IsNewResource() {
		existing, err := client.PrivateLinkScopesGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_monitor_private_link_scope", id.ID())
		}
	}

	parameters := privatelinkscopesapis.AzureMonitorPrivateLinkScope{
		Location: "Global",
		Tags:     utils.ExpandPtrMapStringString(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("access_mode"); ok {
		parameters.Properties.AccessModeSettings = expandMonitorPrivateLinkScopeAccessMode(v.(*pluginsdk.ResourceData))
	}
	if _, err := client.PrivateLinkScopesCreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceMonitorPrivateLinkScopeRead(d, meta)
}

func resourceMonitorPrivateLinkScopeRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopesapis.ParsePrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.PrivateLinkScopesGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s does not exist - removing from state!", *id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrivateLinkScopeName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if err = d.Set("tags", utils.FlattenPtrMapStringString(model.Tags)); err != nil {
			return err
		}

		if err = d.Set("access_mode", flattenMonitorPrivateLinkScopeAccessMode(model.Properties.AccessModeSettings)); err != nil {
			return err
		}

	}

	return nil
}

func resourceMonitorPrivateLinkScopeDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Monitor.PrivateLinkScopesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkscopesapis.ParsePrivateLinkScopeID(d.Id())
	if err != nil {
		return err
	}

	err = client.PrivateLinkScopesDeleteThenPoll(ctx, *id)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandMonitorPrivateLinkScopeAccessMode(d *pluginsdk.ResourceData) privatelinkscopesapis.AccessModeSettings {
	input, ok := d.Get("job_storage_account").([]interface{})
	if input == nil || !ok {
		return privatelinkscopesapis.AccessModeSettings{}
	}

	v := input[0].(map[string]interface{})

	action := privatelinkscopesapis.AccessModeSettings{
		Exclusions:          expandMonitorPrivateLinkScopeExclusions(v["exclusions"].([]interface{})),
		QueryAccessMode:     privatelinkscopesapis.AccessMode(v["query"].(string)),
		IngestionAccessMode: privatelinkscopesapis.AccessMode(v["ingestion"].(string)),
	}

	return action
}

func expandMonitorPrivateLinkScopeExclusions(i []interface{}) *[]privatelinkscopesapis.AccessModeSettingsExclusion {
	if i == nil {
		return nil
	}

	exclusions := make([]privatelinkscopesapis.AccessModeSettingsExclusion, len(i))

	for i, v := range i {
		exclusion := v.(map[string]interface{})

		exclusions[i] = privatelinkscopesapis.AccessModeSettingsExclusion{
			PrivateEndpointConnectionName: pointer.To(exclusion["connection_name"].(string)),
			QueryAccessMode:               pointer.To(privatelinkscopesapis.AccessMode(exclusion["query"].(string))),
			IngestionAccessMode:           pointer.To(privatelinkscopesapis.AccessMode(exclusion["ingestion"].(string))),
		}
	}

	return &exclusions
}

func flattenMonitorPrivateLinkScopeAccessMode(accessModeSettings privatelinkscopesapis.AccessModeSettings) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"exclusions": flattenMonitorPrivateLinkScopeExclusions(accessModeSettings.Exclusions),
			"query":      accessModeSettings.QueryAccessMode,
			"ingestion":  accessModeSettings.IngestionAccessMode,
		},
	}
}

func flattenMonitorPrivateLinkScopeExclusions(exclusions *[]privatelinkscopesapis.AccessModeSettingsExclusion) []map[string]interface{} {
	if exclusions == nil {
		return nil
	}

	exclusion := make([]map[string]interface{}, len(*exclusions))

	for i, v := range *exclusions {
		exclusion[i] = map[string]interface{}{
			"connection_name": v.PrivateEndpointConnectionName,
			"query":           v.QueryAccessMode,
			"ingestion":       v.IngestionAccessMode,
		}
	}

	return exclusion
}
