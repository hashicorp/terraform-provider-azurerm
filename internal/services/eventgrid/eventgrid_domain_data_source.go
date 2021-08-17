package eventgrid

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"input_schema": {
				Type: pluginsdk.TypeString,
			},

			//lintignore:XS003
			"input_mapping_fields": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type: pluginsdk.TypeString,
						},
						"topic": {
							Type: pluginsdk.TypeString,
						},
						"event_time": {
							Type: pluginsdk.TypeString,
						},
						"event_type": {
							Type: pluginsdk.TypeString,
						},
						"subject": {
							Type: pluginsdk.TypeString,
						},
						"data_version": {
							Type: pluginsdk.TypeString,
						},
					},
				},
			},

			//lintignore:XS003
			"input_mapping_default_values": {
				Type:     pluginsdk.TypeList,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"event_type": {
							Type: pluginsdk.TypeString,
						},
						"subject": {
							Type: pluginsdk.TypeString,
						},
						"data_version": {
							Type: pluginsdk.TypeString,
						},
					},
				},
			},

			"public_network_access_enabled": eventSubscriptionPublicNetworkAccessEnabled(),

			"inbound_ip_rule": eventSubscriptionInboundIPRule(),

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

			"tags": tags.Schema(),
		},
	}
}

func dataSourceEventGridDomainRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).EventGrid.DomainsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: EventGrid Domain %s (Resource Group %s) was not found: %+v", name, resourceGroup, err)
		}

		return fmt.Errorf("making Read request on EventGrid Domain '%s': %+v", name, err)
	}

	d.SetId(*resp.ID)
	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.DomainProperties; props != nil {
		d.Set("endpoint", props.Endpoint)

		d.Set("input_schema", string(props.InputSchema))

		inputMappingFields, err := flattenAzureRmEventgridDomainInputMapping(props.InputSchemaMapping)
		if err != nil {
			return fmt.Errorf("flattening `input_schema_mapping_fields` for EventGrid Domain %q (Resource Group %q): %s", name, resourceGroup, err)
		}
		if err := d.Set("input_mapping_fields", inputMappingFields); err != nil {
			return fmt.Errorf("setting `input_schema_mapping_fields` for EventGrid Domain %q (Resource Group %q): %s", name, resourceGroup, err)
		}

		inputMappingDefaultValues, err := flattenAzureRmEventgridDomainInputMappingDefaultValues(props.InputSchemaMapping)
		if err != nil {
			return fmt.Errorf("flattening `input_schema_mapping_default_values` for EventGrid Domain %q (Resource Group %q): %s", name, resourceGroup, err)
		}
		if err := d.Set("input_mapping_default_values", inputMappingDefaultValues); err != nil {
			return fmt.Errorf("setting `input_schema_mapping_fields` for EventGrid Domain %q (Resource Group %q): %s", name, resourceGroup, err)
		}

		publicNetworkAccessEnabled := flattenPublicNetworkAccess(props.PublicNetworkAccess)
		if err := d.Set("public_network_access_enabled", publicNetworkAccessEnabled); err != nil {
			return fmt.Errorf("setting `public_network_access_enabled` in EventGrid Domain %q (Resource Group %q): %+v", name, resourceGroup, err)
		}

		inboundIPRules := flattenInboundIPRules(props.InboundIPRules)
		if err := d.Set("inbound_ip_rule", inboundIPRules); err != nil {
			return fmt.Errorf("setting `inbound_ip_rule` in EventGrid Domain %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	keys, err := client.ListSharedAccessKeys(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving Shared Access Keys for EventGrid Domain %q: %+v", name, err)
	}
	d.Set("primary_access_key", keys.Key1)
	d.Set("secondary_access_key", keys.Key2)

	return tags.FlattenAndSet(d, resp.Tags)
}
