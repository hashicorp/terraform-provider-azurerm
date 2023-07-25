// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package eventgrid

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15/domains"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceEventGridDomain() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceEventGridDomainRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.LocationComputed(),

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"input_schema": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			//lintignore:XS003
			"input_mapping_fields": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"topic": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"event_time": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"event_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subject": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"data_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			//lintignore:XS003
			"input_mapping_default_values": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"event_type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"subject": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"data_version": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"inbound_ip_rule": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ip_mask": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"action": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"tags": commonschema.TagsDataSource(),
		},
	}
}

func dataSourceEventGridDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.Domains
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := domains.NewDomainID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	keys, err := client.ListSharedAccessKeys(ctx, id)
	if err != nil {
		if !response.WasForbidden(resp.HttpResponse) {
			return fmt.Errorf("retrieving Shared Access Keys for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	d.Set("name", id.DomainName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("endpoint", props.Endpoint)

			inputSchema := ""
			if props.InputSchema != nil {
				inputSchema = string(*props.InputSchema)
			}
			d.Set("input_schema", inputSchema)

			inputMappingFields := flattenDomainInputMapping(props.InputSchemaMapping)
			if err := d.Set("input_mapping_fields", inputMappingFields); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields` for %s: %+v", id, err)
			}

			inputMappingDefaultValues := flattenDomainInputMappingDefaultValues(props.InputSchemaMapping)
			if err := d.Set("input_mapping_default_values", inputMappingDefaultValues); err != nil {
				return fmt.Errorf("setting `input_schema_mapping_fields` for %s: %+v", id, err)
			}

			publicNetworkAccessEnabled := true
			if props.PublicNetworkAccess != nil && *props.PublicNetworkAccess == domains.PublicNetworkAccessDisabled {
				publicNetworkAccessEnabled = false
			}
			d.Set("public_network_access_enabled", publicNetworkAccessEnabled)

			inboundIPRules := flattenDomainInboundIPRules(props.InboundIPRules)
			if err := d.Set("inbound_ip_rule", inboundIPRules); err != nil {
				return fmt.Errorf("setting `inbound_ip_rule` in %s: %+v", id, err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return fmt.Errorf("setting `tags`: %+v", err)
		}
	}

	if model := keys.Model; model != nil {
		d.Set("primary_access_key", model.Key1)
		d.Set("secondary_access_key", model.Key2)
	}

	return nil
}
