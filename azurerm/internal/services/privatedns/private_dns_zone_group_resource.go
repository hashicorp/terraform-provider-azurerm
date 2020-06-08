package privatedns

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateDnsZoneGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateDnsZoneGroupCreateUpdate,
		Read:   resourceArmPrivateDnsZoneGroupRead,
		Update: resourceArmPrivateDnsZoneGroupCreateUpdate,
		Delete: resourceArmPrivateDnsZoneGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidatePrivateDnsZoneGroupName,
			},

			"private_endpoint_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"private_dns_zone_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"private_dns_zone_configs": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_dns_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"record_sets": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"fqdn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ttl": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"ip_addresses": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceArmPrivateDnsZoneGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	privateEndpointId := d.Get("private_endpoint_id").(string)
	privateDnsZoneIds := d.Get("private_dns_zone_ids").([]interface{})

	privateEndpoint, err := parse.PrivateEndpointResourceID(privateEndpointId)
	if err != nil {
		return err
	}
	privateDnsZones, err := parse.PrivateDnsZoneResourceIDs(privateDnsZoneIds)
	if err != nil {
		return err
	}

	name := d.Get("name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, privateEndpoint.ResourceGroup, privateEndpoint.Name, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Private Endpoint DNS Zone Group %q (Resource Group %q): %+v", name, privateEndpoint.ResourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_dns_zone_group", *existing.ID)
		}
	}

	parameters := network.PrivateDNSZoneGroup{}

	parameters.Name = utils.String(name)
	privateDnsZoneConfigs := make([]network.PrivateDNSZoneConfig, 0)

	for _, item := range privateDnsZones {
		v := network.PrivateDNSZoneConfig{
			Name: utils.String(item.Name),
			PrivateDNSZonePropertiesFormat: &network.PrivateDNSZonePropertiesFormat{
				PrivateDNSZoneID: utils.String(item.ID),
			},
		}

		privateDnsZoneConfigs = append(privateDnsZoneConfigs, v)
	}

	parameters.PrivateDNSZoneGroupPropertiesFormat = &network.PrivateDNSZoneGroupPropertiesFormat{
		PrivateDNSZoneConfigs: &privateDnsZoneConfigs,
	}

	future, err := client.CreateOrUpdate(ctx, privateEndpoint.ResourceGroup, privateEndpoint.Name, name, parameters)
	if err != nil {
		return fmt.Errorf("creating Private Endpoint DNS Zone Group %q (Resource Group %q): %+v", name, privateEndpoint.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Private Endpoint DNS Zone Group %q (Resource Group %q): %+v", name, privateEndpoint.ResourceGroup, err)
	}

	resp, err := client.Get(ctx, privateEndpoint.ResourceGroup, privateEndpoint.Name, name)
	if err != nil {
		return fmt.Errorf("retrieving Private DNS Zone Group %q (Resource Group %q): %+v", name, privateEndpoint.ResourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private DNS Zone Group %q (Resource Group %q): %+v", name, privateEndpoint.ResourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmPrivateDnsZoneGroupRead(d, meta)
}

func resourceArmPrivateDnsZoneGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	name := id.Path["privateDnsZoneGroups"]
	privateEndpointName := id.Path["privateEndpoints"]
	resourceGroup := id.ResourceGroup

	resp, err := client.Get(ctx, resourceGroup, privateEndpointName, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private DNS Zone Group %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Private DNS Zone Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("private_endpoint_id", d.Get("private_endpoint_id").(string))

	if props := resp.PrivateDNSZoneGroupPropertiesFormat; props != nil {
		d.Set("private_dns_zone_ids", flattenArmPrivateDnsZoneIds(props.PrivateDNSZoneConfigs))
		if err := d.Set("private_dns_zone_configs", flattenArmPrivateDnsZoneConfigs(props.PrivateDNSZoneConfigs, d.Id())); err != nil {
			return fmt.Errorf("setting private_dns_zone_configs : %+v", err)
		}
	}

	return nil
}

func flattenArmPrivateDnsZoneIds(input *[]network.PrivateDNSZoneConfig) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if props := item.PrivateDNSZonePropertiesFormat; props != nil {
			if privateDnsZoneId := props.PrivateDNSZoneID; privateDnsZoneId != nil {
				results = append(results, *props.PrivateDNSZoneID)
			}
		}
	}

	return results
}

func resourceArmPrivateDnsZoneGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateDnsZoneGroups"]
	privateEndpointName := id.Path["privateEndpoints"]

	future, err := client.Delete(ctx, resourceGroup, privateEndpointName, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Private DNS Zone Group %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Private DNS Zone Group %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmPrivateDnsZoneGroup(input []interface{}) network.PrivateDNSZoneGroup {
	result := network.PrivateDNSZoneGroup{}
	if len(input) == 0 {
		return result
	}

	dnsZoneConfigs := make([]network.PrivateDNSZoneConfig, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		name := v["name"].(string)
		zoneConfigs := v["zone_config"].([]interface{})

		result.Name = utils.String(name)

		for _, zoneConfig := range zoneConfigs {
			z := zoneConfig.(map[string]interface{})
			zoneName := z["name"].(string)
			zoneId := z["private_dns_zone_id"].(string)

			config := network.PrivateDNSZoneConfig{
				Name: utils.String(zoneName),
				PrivateDNSZonePropertiesFormat: &network.PrivateDNSZonePropertiesFormat{
					PrivateDNSZoneID: utils.String(zoneId),
				},
			}

			dnsZoneConfigs = append(dnsZoneConfigs, config)
		}
	}

	result.PrivateDNSZoneGroupPropertiesFormat = &network.PrivateDNSZoneGroupPropertiesFormat{
		PrivateDNSZoneConfigs: &dnsZoneConfigs,
	}

	return result
}

func flattenArmPrivateDnsZoneConfigs(input *[]network.PrivateDNSZoneConfig, zoneGroupId string) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		result := make(map[string]interface{})

		if name := v.Name; name != nil {
			result["name"] = *name
			// I have to consturct this because the SDK does not expose it in it's PrivateDNSZoneConfig struct
			result["id"] = fmt.Sprintf("%s/privateDnsZoneConfigs/%s", zoneGroupId, *name)
		}

		if props := v.PrivateDNSZonePropertiesFormat; props != nil {
			if zoneId := props.PrivateDNSZoneID; zoneId != nil {
				result["private_dns_zone_id"] = *zoneId
			}

			if recordSets := props.RecordSets; recordSets != nil {
				result["record_sets"] = flattenArmPrivateDnsZoneRecordSets(recordSets)
			}
		}

		output = append(output, result)
	}

	log.Printf("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n*********************")
	log.Printf("\nconfigs = %+v", output)
	log.Printf("*********************\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n")

	return output
}

func flattenArmPrivateDnsZoneRecordSets(input *[]network.RecordSet) []map[string]interface{} {
	output := make([]map[string]interface{}, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		result := make(map[string]interface{})

		if recordName := v.RecordSetName; recordName != nil {
			result["name"] = *recordName
		}

		if recordType := v.RecordType; recordType != nil {
			result["type"] = *recordType
		}

		if fqdn := v.Fqdn; fqdn != nil {
			result["fqdn"] = *fqdn
		}

		if ttl := v.TTL; ttl != nil {
			result["ttl"] = int(*ttl)
		}

		if ipAddresses := v.IPAddresses; ipAddresses != nil {
			result["ip_addresses"] = *ipAddresses
		}

		output = append(output, result)
	}

	return output
}
