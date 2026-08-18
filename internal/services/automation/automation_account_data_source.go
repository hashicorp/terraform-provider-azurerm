// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package automation

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2019-06-01/agentregistrationinformation"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"primary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secondary_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityComputed(),

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"location": commonschema.LocationComputed(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": commonschema.TagsDataSource(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"encryption": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"user_assigned_identity_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"key_vault_key_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"local_authentication_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"dsc_server_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"dsc_primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"dsc_secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"private_endpoint_connection": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"hybrid_service_url": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	iclient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	client := meta.(*clients.Client).Automation.AutomationAccount
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId

	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := automationaccount.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retreiving %s: %+v", id, err)
	}
	d.SetId(id.ID())

	infoId := agentregistrationinformation.NewAutomationAccountID(id.SubscriptionId, id.ResourceGroupName, id.AutomationAccountName)
	infoResp, err := iclient.Get(ctx, infoId)
	if err != nil {
		return fmt.Errorf("retreiving Agent Registration Information for %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		flattenedIdentity, err := identity.FlattenSystemAndUserAssignedMap(model.Identity)
		if err != nil {
			return fmt.Errorf("flattening `identity`: %+v", err)
		}
		if err := d.Set("identity", flattenedIdentity); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("private_endpoint_connection", flattenPrivateEndpointConnectionsDataSource(props.PrivateEndpointConnections))
			d.Set("hybrid_service_url", props.AutomationHybridServiceURL)
			d.Set("local_authentication_enabled", !pointer.From(props.DisableLocalAuth))
			d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess))
			d.Set("sku_name", pointer.From(props.Sku))

			if err := d.Set("encryption", flattenEncryption(props.Encryption)); err != nil {
				return fmt.Errorf("setting `encryption`: %+v", err)
			}

		}
		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	endpoint := ""
	primaryKey := ""
	secondaryKey := ""
	if model := infoResp.Model; model != nil {
		endpoint = pointer.From(model.Endpoint)
		if keys := model.Keys; keys != nil {
			primaryKey = pointer.From(keys.Primary)
			secondaryKey = pointer.From(keys.Secondary)
		}
	}
	d.Set("endpoint", endpoint)
	d.Set("primary_key", primaryKey)
	d.Set("secondary_key", secondaryKey)

	d.Set("dsc_server_endpoint", endpoint)
	d.Set("dsc_primary_access_key", primaryKey)
	d.Set("dsc_secondary_access_key", secondaryKey)

	return nil
}

func flattenPrivateEndpointConnectionsDataSource(input *[]automationaccount.PrivateEndpointConnection) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, item := range *input {
		output = append(output, map[string]interface{}{
			"id":   pointer.From(item.Id),
			"name": pointer.From(item.Name),
		})
	}
	return output
}
