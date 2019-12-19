package azurerm

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmPrivateLinkServiceEndpointConnections() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmPrivateLinkServiceEndpointConnectionsRead,
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"private_endpoint_connections": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connection_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_endpoint_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_endpoint_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceArmPrivateLinkServiceEndpointConnectionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	serviceId := d.Get("service_id").(string)

	id, err := azure.ParseAzureResourceID(serviceId)
	if err != nil {
		return fmt.Errorf("Error parsing %q: %s", serviceId, err)
	}

	name := id.Path["privateLinkServices"]
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Error: Private Link Service %q (Resource Group %q) was not found", name, resourceGroup)
		}
		return fmt.Errorf("Error reading Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Service Endpoint Connection Status %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("service_id", serviceId)
	d.Set("service_name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if props := resp.PrivateLinkServiceProperties; props != nil {
		if err := d.Set("private_endpoint_connections", dataSourceflattenArmPrivateLinkServicePrivateEndpointConnections(props.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("Error setting `private_endpoint_connections`: %+v", err)
		}
	}

	d.SetId(fmt.Sprintf("%s/privateLinkServiceEndpointConnections/%s", *resp.ID, name))

	return nil
}

func dataSourceflattenArmPrivateLinkServicePrivateEndpointConnections(input *[]network.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})
		if id := item.ID; id != nil {
			v["connection_id"] = *id
		}
		if name := item.Name; name != nil {
			v["connection_name"] = *name
		}

		if props := item.PrivateEndpointConnectionProperties; props != nil {
			if p := props.PrivateEndpoint; p != nil {
				if id := p.ID; id != nil {
					v["private_endpoint_id"] = *id

					id, _ := azure.ParseAzureResourceID(*id)
					name := id.Path["privateEndpoints"]
					if name != "" {
						v["private_endpoint_name"] = name
					}
				}
			}

			if s := props.PrivateLinkServiceConnectionState; s != nil {
				if a := s.ActionsRequired; a != nil {
					v["action_required"] = *a
				} else {
					v["action_required"] = "none"
				}
				if d := s.Description; d != nil {
					v["description"] = *d
				}
				if t := s.Status; t != nil {
					v["status"] = *t
				}
			}
		}

		results = append(results, v)
	}

	return results
}
