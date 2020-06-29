package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateEndpointCreateUpdate,
		Read:   resourceArmPrivateEndpointRead,
		Update: resourceArmPrivateEndpointCreateUpdate,
		Delete: resourceArmPrivateEndpointDelete,
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
				ValidateFunc: ValidatePrivateLinkName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"private_dns_zone_group": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: ValidatePrivateLinkName,
						},
						"private_dns_zone_ids": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: parse.ValidatePrivateDnsZoneResourceID,
							},
						},
					},
				},
			},

			"private_service_connection": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: ValidatePrivateLinkName,
						},
						"is_manual_connection": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"private_connection_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"subresource_names": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: ValidatePrivateLinkSubResourceName,
							},
						},
						"request_message": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 140),
						},
						"private_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"custom_dns_configs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fqdn": {
							Type:     schema.TypeString,
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

			"private_dns_zone_configs": {
				Type:     schema.TypeList,
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

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpointClient
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if err := ValidatePrivateEndpointSettings(d); err != nil {
		return fmt.Errorf("validating the configuration for the Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_endpoint", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]interface{})
	privateServiceConnections := d.Get("private_service_connection").([]interface{})
	subnetId := d.Get("subnet_id").(string)

	parameters := network.PrivateEndpoint{
		Location: utils.String(location),
		PrivateEndpointProperties: &network.PrivateEndpointProperties{
			PrivateLinkServiceConnections:       expandArmPrivateLinkEndpointServiceConnection(privateServiceConnections, false),
			ManualPrivateLinkServiceConnections: expandArmPrivateLinkEndpointServiceConnection(privateServiceConnections, true),
			Subnet: &network.Subnet{
				ID: utils.String(subnetId),
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		if azure.StringContains(err.Error(), "is missing required parameter 'group Id'") {
			return fmt.Errorf("creating Private Endpoint %q (Resource Group %q) due to missing 'group Id', ensure that the 'subresource_names' type is populated: %+v", name, resourceGroup, err)
		} else {
			return fmt.Errorf("creating Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("retrieving Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	// now create the dns zone group
	// first I have to see if the dns zone group exists, if it does I need to delete it an re-create it because you can only have one per private endpoint
	if d.HasChange("private_dns_zone_group") || d.IsNewResource() {
		oldRaw, newRaw := d.GetChange("private_dns_zone_group")
		oldPrivateDnsZoneGroup := make(map[string]interface{})
		if oldRaw != nil {
			for _, v := range oldRaw.([]interface{}) {
				oldPrivateDnsZoneGroup = v.(map[string]interface{})
			}
		}

		newPrivateDnsZoneGroup := make(map[string]interface{})
		if newRaw != nil {
			for _, v := range newRaw.([]interface{}) {
				newPrivateDnsZoneGroup = v.(map[string]interface{})
			}
		}

		if len(newPrivateDnsZoneGroup) == 0 && len(oldPrivateDnsZoneGroup) != 0 {
			if err := resourceArmPrivateDnsZoneGroupDelete(d, meta, oldPrivateDnsZoneGroup["id"].(string)); err != nil {
				return err
			}
		} else if len(newPrivateDnsZoneGroup) != 0 && len(oldPrivateDnsZoneGroup) != 0 {
			if oldPrivateDnsZoneGroup["name"].(string) != newPrivateDnsZoneGroup["name"].(string) {
				if err := resourceArmPrivateDnsZoneGroupDelete(d, meta, oldPrivateDnsZoneGroup["id"].(string)); err != nil {
					return err
				}
			}
		}
	}

	for _, v := range privateDnsZoneGroup {
		item := v.(map[string]interface{})
		dnsGroupName := item["name"].(string)
		privateDnsZoneIds := item["private_dns_zone_ids"].([]interface{})
		privateDnsZones, err := parse.PrivateDnsZoneResourceIDs(privateDnsZoneIds)
		if err != nil {
			return err
		}

		privateDnsZoneConfigs := make([]network.PrivateDNSZoneConfig, 0)

		for _, item := range *privateDnsZones {
			v := network.PrivateDNSZoneConfig{
				Name: utils.String(item.Name),
				PrivateDNSZonePropertiesFormat: &network.PrivateDNSZonePropertiesFormat{
					PrivateDNSZoneID: utils.String(item.ID),
				},
			}

			privateDnsZoneConfigs = append(privateDnsZoneConfigs, v)
		}

		parameters := network.PrivateDNSZoneGroup{}
		parameters.Name = utils.String(name)
		parameters.PrivateDNSZoneGroupPropertiesFormat = &network.PrivateDNSZoneGroupPropertiesFormat{
			PrivateDNSZoneConfigs: &privateDnsZoneConfigs,
		}

		future, err := dnsClient.CreateOrUpdate(ctx, resourceGroup, name, dnsGroupName, parameters)
		if err != nil {
			return fmt.Errorf("creating Private DNS Zone Group %q Private Endpoint %q (Resource Group %q): %+v", dnsGroupName, name, resourceGroup, err)
		}
		if err = future.WaitForCompletionRef(ctx, dnsClient.Client); err != nil {
			return fmt.Errorf("waiting for creation of Private DNS Zone Group %q Private Endpoint %q (Resource Group %q): %+v", dnsGroupName, name, resourceGroup, err)
		}

		resp, err := dnsClient.Get(ctx, resourceGroup, name, dnsGroupName)
		if err != nil {
			return fmt.Errorf("retrieving Private DNS Zone Group %q (Resource Group %q): %+v", dnsGroupName, resourceGroup, err)
		}
		if resp.ID == nil || *resp.ID == "" {
			return fmt.Errorf("API returns a nil/empty id on Private DNS Zone Group %q (Resource Group %q): %+v", dnsGroupName, resourceGroup, err)
		}
	}

	return resourceArmPrivateEndpointRead(d, meta)
}

func resourceArmPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpointClient
	nicsClient := meta.(*clients.Client).Network.InterfacesClient
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateEndpoints"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.PrivateEndpointProperties; props != nil {
		privateIpAddress := ""

		if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
			nic := (*nics)[0]
			if nic.ID != nil && *nic.ID != "" {
				privateIpAddress = getPrivateIpAddress(ctx, nicsClient, *nic.ID)
			}
		}

		flattenedConnection := flattenArmPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections)
		for _, item := range flattenedConnection {
			v := item.(map[string]interface{})
			v["private_ip_address"] = privateIpAddress
		}
		if err := d.Set("private_service_connection", flattenedConnection); err != nil {
			return fmt.Errorf("setting `private_service_connection`: %+v", err)
		}

		subnetId := ""
		if subnet := props.Subnet; subnet != nil {
			subnetId = *subnet.ID
		}
		d.Set("subnet_id", subnetId)
		d.Set("custom_dns_configs", flattenArmCustomDnsConfigs(props.CustomDNSConfigs))
	}

	// DNS Zone Read Here...
	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]interface{})
	if len(privateDnsZoneGroup) > 0 {
		for _, v := range privateDnsZoneGroup {
			dnsZoneGroup := v.(map[string]interface{})

			dnsResp, err := dnsClient.Get(ctx, resourceGroup, name, dnsZoneGroup["name"].(string))
			if err != nil {
				return fmt.Errorf("reading Private DNS Zone Group %q (Resource Group %q): %+v", dnsZoneGroup["name"].(string), resourceGroup, err)
			}

			if err := d.Set("private_dns_zone_group", flattenArmPrivateDnsZoneGroup(dnsResp)); err != nil {
				return err
			}

			// now split out the private dns zone configs into there own block
			if props := dnsResp.PrivateDNSZoneGroupPropertiesFormat; props != nil {
				if err := d.Set("private_dns_zone_configs", flattenArmPrivateDnsZoneConfigs(props.PrivateDNSZoneConfigs, *dnsResp.ID)); err != nil {
					return fmt.Errorf("setting private_dns_zone_configs : %+v", err)
				}
			}
		}
	} else {
		// remove associated configs, if any
		d.Set("private_dns_zone_configs", make([]interface{}, 0))
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPrivateEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpointClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// delete private dns zone first if it exists
	dnsRaw := d.Get("private_dns_zone_group")
	privateDnsZoneGroup := make(map[string]interface{})
	if dnsRaw != nil {
		for _, v := range dnsRaw.([]interface{}) {
			privateDnsZoneGroup = v.(map[string]interface{})
		}
	}

	if len(privateDnsZoneGroup) != 0 {
		if err := resourceArmPrivateDnsZoneGroupDelete(d, meta, privateDnsZoneGroup["id"].(string)); err != nil {
			return err
		}
	}

	privateEndpoint, err := parse.PrivateEndpointResourceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, privateEndpoint.ResourceGroup, privateEndpoint.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Private Endpoint %q (Resource Group %q): %+v", privateEndpoint.Name, privateEndpoint.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Private Endpoint %q (Resource Group %q): %+v", privateEndpoint.Name, privateEndpoint.ResourceGroup, err)
		}
	}

	return nil
}

func resourceArmPrivateDnsZoneGroupDelete(d *schema.ResourceData, meta interface{}, oldId string) error {
	client := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if oldId == "" {
		return nil
	}

	privateEndpoint, err := parse.PrivateEndpointResourceID(d.Id())
	if err != nil {
		return err
	}

	privateDnsZoneGroupId, err := parse.PrivateDnsZoneGroupResourceID(oldId)
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, privateEndpoint.ResourceGroup, privateEndpoint.Name, privateDnsZoneGroupId.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Private DNS Zone Group %q (Resource Group %q): %+v", privateDnsZoneGroupId.Name, privateEndpoint.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Private DNS Zone Group %q (Resource Group %q): %+v", privateDnsZoneGroupId.Name, privateEndpoint.ResourceGroup, err)
		}
	}

	return nil
}

func expandArmPrivateLinkEndpointServiceConnection(input []interface{}, parseManual bool) *[]network.PrivateLinkServiceConnection {
	results := make([]network.PrivateLinkServiceConnection, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		privateConnectonResourceId := v["private_connection_resource_id"].(string)
		subresourceNames := v["subresource_names"].([]interface{})
		requestMessage := v["request_message"].(string)
		isManual := v["is_manual_connection"].(bool)
		name := v["name"].(string)

		if isManual == parseManual {
			result := network.PrivateLinkServiceConnection{
				Name: utils.String(name),
				PrivateLinkServiceConnectionProperties: &network.PrivateLinkServiceConnectionProperties{
					GroupIds:             utils.ExpandStringSlice(subresourceNames),
					PrivateLinkServiceID: utils.String(privateConnectonResourceId),
				},
			}

			if requestMessage != "" {
				result.PrivateLinkServiceConnectionProperties.RequestMessage = utils.String(requestMessage)
			}

			results = append(results, result)
		}
	}

	return &results
}

func flattenArmPrivateDnsZoneGroup(input network.PrivateDNSZoneGroup) []interface{} {
	output := make([]interface{}, 0)
	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if name := input.Name; name != nil {
		result["name"] = *name
	}

	if props := input.PrivateDNSZoneGroupPropertiesFormat; props != nil {
		result["private_dns_zone_ids"] = flattenArmPrivateDnsZoneIds(props.PrivateDNSZoneConfigs)
	}
	output = append(output, result)
	return output
}

func flattenArmCustomDnsConfigs(customDnsConfigs *[]network.CustomDNSConfigPropertiesFormat) []interface{} {
	results := make([]interface{}, 0)
	if customDnsConfigs == nil {
		return results
	}

	for _, item := range *customDnsConfigs {
		results = append(results, map[string]interface{}{
			"fqdn":         item.Fqdn,
			"ip_addresses": utils.FlattenStringSlice(item.IPAddresses),
		})
	}

	return results
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

func flattenArmPrivateLinkEndpointServiceConnection(serviceConnections *[]network.PrivateLinkServiceConnection, manualServiceConnections *[]network.PrivateLinkServiceConnection) []interface{} {
	results := make([]interface{}, 0)
	if serviceConnections == nil && manualServiceConnections == nil {
		return results
	}

	if serviceConnections != nil {
		for _, item := range *serviceConnections {
			name := ""
			if item.Name != nil {
				name = *item.Name
			}

			privateConnectionId := ""
			subResourceNames := make([]interface{}, 0)

			if props := item.PrivateLinkServiceConnectionProperties; props != nil {
				if v := props.GroupIds; v != nil {
					subResourceNames = utils.FlattenStringSlice(v)
				}
				if props.PrivateLinkServiceID != nil {
					privateConnectionId = *props.PrivateLinkServiceID
				}
			}
			results = append(results, map[string]interface{}{
				"name":                           name,
				"is_manual_connection":           false,
				"private_connection_resource_id": privateConnectionId,
				"subresource_names":              subResourceNames,
			})
		}
	}

	if manualServiceConnections != nil {
		for _, item := range *manualServiceConnections {
			name := ""
			if item.Name != nil {
				name = *item.Name
			}

			privateConnectionId := ""
			requestMessage := ""
			subResourceNames := make([]interface{}, 0)

			if props := item.PrivateLinkServiceConnectionProperties; props != nil {
				if v := props.GroupIds; v != nil {
					subResourceNames = utils.FlattenStringSlice(v)
				}
				if props.PrivateLinkServiceID != nil {
					privateConnectionId = *props.PrivateLinkServiceID
				}
				if props.RequestMessage != nil {
					requestMessage = *props.RequestMessage
				}
			}

			results = append(results, map[string]interface{}{
				"name":                           name,
				"is_manual_connection":           true,
				"private_connection_resource_id": privateConnectionId,
				"request_message":                requestMessage,
				"subresource_names":              subResourceNames,
			})
		}
	}

	return results
}
