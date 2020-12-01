package network

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-05-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	privateDnsParse "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/parse"
	privateDnsValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/privatedns/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateEndpoint() *schema.Resource {
	return &schema.Resource{
		// TODO: add a state migration to ensure the ID's stable
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
								ValidateFunc: privateDnsValidate.PrivateDnsZoneID,
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
	// TODO: split this into a Create and an Update

	client := meta.(*clients.Client).Network.PrivateEndpointClient
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPrivateEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if err := ValidatePrivateEndpointSettings(d); err != nil {
		return fmt.Errorf("validating the configuration for the Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
			}
		}

		if existing.PrivateEndpointProperties != nil {
			return tf.ImportAsExistsError("azurerm_private_endpoint", id.ID(""))
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

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		if strings.EqualFold(err.Error(), "is missing required parameter 'group Id'") {
			return fmt.Errorf("creating Private Endpoint %q (Resource Group %q) due to missing 'group Id', ensure that the 'subresource_names' type is populated: %+v", id.Name, id.ResourceGroup, err)
		} else {
			return fmt.Errorf("creating Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.SetId(id.ID(""))

	// now create the dns zone group
	// first I have to see if the dns zone group exists, if it does I need to delete it an re-create it because you can only have one per private endpoint
	if d.HasChange("private_dns_zone_group") || d.IsNewResource() {
		// TODO: shouldn't this be pulling the list from Azure?!
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

		needToRemove := len(newPrivateDnsZoneGroup) == 0 && len(oldPrivateDnsZoneGroup) != 0
		nameHasChanged := (len(newPrivateDnsZoneGroup) != 0 && len(oldPrivateDnsZoneGroup) != 0) && oldPrivateDnsZoneGroup["name"].(string) != newPrivateDnsZoneGroup["name"].(string)
		if needToRemove || nameHasChanged {
			log.Printf("[DEBUG] Deleting the Existing Private DNS Zone Group associated with Private Endpoint %q / Resource Group %q..", id.Name, id.ResourceGroup)
			if err := deletePrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, id); err != nil {
				return err
			}
			log.Printf("[DEBUG] Deleted the Existing Private DNS Zone Group associated with Private Endpoint %q / Resource Group %q.", id.Name, id.ResourceGroup)
		}
	}

	if len(privateDnsZoneGroup) > 0 {
		log.Printf("[DEBUG] Creating Private DNS Zone Group associated with Private Endpoint %q / Resource Group %q..", id.Name, id.ResourceGroup)
		if err := createPrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, id, privateDnsZoneGroup); err != nil {
			return err
		}
		log.Printf("[DEBUG] Created the Existing Private DNS Zone Group associated with Private Endpoint %q / Resource Group %q.", id.Name, id.ResourceGroup)
	}

	return resourceArmPrivateEndpointRead(d, meta)
}

func resourceArmPrivateEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpointClient
	nicsClient := meta.(*clients.Client).Network.InterfacesClient
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.PrivateEndpointProperties; props != nil {
		if err := d.Set("custom_dns_configs", flattenArmCustomDnsConfigs(props.CustomDNSConfigs)); err != nil {
			return fmt.Errorf("setting `custom_dns_configs`: %+v", err)
		}

		privateIpAddress := ""
		if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
			nic := (*nics)[0]
			if nic.ID != nil && *nic.ID != "" {
				privateIpAddress = getPrivateIpAddress(ctx, nicsClient, *nic.ID)
			}
		}
		// TODO: why not pass in the Private IP Address here?!
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
	}

	// DNS Zone Read Here...
	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]interface{})
	// TODO: switch this to looking them up..
	if len(privateDnsZoneGroup) > 0 {
		for _, v := range privateDnsZoneGroup {
			dnsZoneGroup := v.(map[string]interface{})

			dnsResp, err := dnsClient.Get(ctx, id.ResourceGroup, id.Name, dnsZoneGroup["name"].(string))
			if err != nil {
				return fmt.Errorf("reading Private DNS Zone Group %q (Resource Group %q): %+v", dnsZoneGroup["name"].(string), id.ResourceGroup, err)
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
	dnsZoneGroupsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroupClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting the Private DNS Zone Group associated with Private Endpoint %q / Resource Group %q..", id.Name, id.ResourceGroup)
	if err := deletePrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsZoneGroupsClient, *id); err != nil {
		return err
	}
	log.Printf("[DEBUG] Deleted the Private DNS Zone Group associated with Private Endpoint %q / Resource Group %q.", id.Name, id.ResourceGroup)

	log.Printf("[DEBUG] Deleting the Private Endpoint %q / Resource Group %q..", id.Name, id.ResourceGroup)
	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}
	log.Printf("[DEBUG] Deleted the Private Endpoint %q / Resource Group %q.", id.Name, id.ResourceGroup)

	return nil
}

func createPrivateDnsZoneGroupForPrivateEndpoint(ctx context.Context, client *network.PrivateDNSZoneGroupsClient, id parse.PrivateEndpointId, inputRaw []interface{}) error {
	if len(inputRaw) != 1 {
		return fmt.Errorf("expected a single Private DNS Zone Groups but got %d", len(inputRaw))
	}
	item := inputRaw[0].(map[string]interface{})

	dnsGroupName := item["name"].(string)
	privateDnsZoneIdsRaw := item["private_dns_zone_ids"].([]interface{})
	privateDnsZoneConfigs := make([]network.PrivateDNSZoneConfig, 0)
	for _, item := range privateDnsZoneIdsRaw {
		v := item.(string)

		privateDnsZone, err := privateDnsParse.PrivateDnsZoneID(v)
		if err != nil {
			return err
		}

		privateDnsZoneConfigs = append(privateDnsZoneConfigs, network.PrivateDNSZoneConfig{
			Name: utils.String(privateDnsZone.Name),
			PrivateDNSZonePropertiesFormat: &network.PrivateDNSZonePropertiesFormat{
				PrivateDNSZoneID: utils.String(privateDnsZone.ID("")),
			},
		})
	}

	parameters := network.PrivateDNSZoneGroup{
		Name: utils.String(id.Name),
		PrivateDNSZoneGroupPropertiesFormat: &network.PrivateDNSZoneGroupPropertiesFormat{
			PrivateDNSZoneConfigs: &privateDnsZoneConfigs,
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, dnsGroupName, parameters)
	if err != nil {
		return fmt.Errorf("creating Private DNS Zone Group %q for Private Endpoint %q (Resource Group %q): %+v", dnsGroupName, id.Name, id.ResourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Private DNS Zone Group %q for Private Endpoint %q (Resource Group %q): %+v", dnsGroupName, id.Name, id.ResourceGroup, err)
	}

	return nil
}

func deletePrivateDnsZoneGroupForPrivateEndpoint(ctx context.Context, client *network.PrivateDNSZoneGroupsClient, id parse.PrivateEndpointId) error {
	// lookup and delete the (should be, Single) Private DNS Zone Group associated with this Private Endpoint
	dnsZones, err := client.ListComplete(ctx, id.Name, id.ResourceGroup) // looks odd.. matches the SDK method
	if err != nil {
		if !utils.ResponseWasNotFound(dnsZones.Response().Response) {
			return fmt.Errorf("retrieving Private DNS Zone Groups for Private Endpoint %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}
	for dnsZones.NotDone() {
		privateDnsZoneGroup := dnsZones.Value()
		if privateDnsZoneGroup.ID != nil {
			groupId, err := parse.PrivateDnsZoneGroupID(*privateDnsZoneGroup.ID)
			if err != nil {
				return err
			}

			log.Printf("[DEBUG] Deleting Private DNS Zone Group %q (Private Endpoint %q / Resource Group %q)..", groupId.Name, groupId.PrivateEndpointName, groupId.ResourceGroup)
			future, err := client.Delete(ctx, groupId.ResourceGroup, groupId.PrivateEndpointName, groupId.Name)
			if err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("deleting Private DNS Zone Group %q (Private Endpoint %q / Resource Group %q): %+v", groupId.Name, groupId.PrivateEndpointName, groupId.ResourceGroup, err)
				}
			}

			if !response.WasNotFound(future.Response()) {
				log.Printf("[DEBUG] Waiting for deletion of Private DNS Zone Group %q (Private Endpoint %q / Resource Group %q)..", groupId.Name, groupId.PrivateEndpointName, groupId.ResourceGroup)
				if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
					if !response.WasNotFound(future.Response()) {
						return fmt.Errorf("waiting for deletion of Private DNS Zone Group %q (Private Endpoint %q / Resource Group %q): %+v", groupId.Name, groupId.PrivateEndpointName, groupId.ResourceGroup, err)
					}
				}
				log.Printf("[DEBUG] Deleted Private DNS Zone Group %q (Private Endpoint %q / Resource Group %q).", groupId.Name, groupId.PrivateEndpointName, groupId.ResourceGroup)
			}
		}

		if err := dnsZones.NextWithContext(ctx); err != nil {
			return err
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
			// I have to consturct this because the SDK does not expose it in its PrivateDNSZoneConfig struct
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
