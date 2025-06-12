// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	mariadbServers "github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/privatednszonegroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2024-05-01/privateendpoints"
	postgresqlServers "github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01/redis"
	"github.com/hashicorp/go-azure-sdk/resource-manager/signalr/2024-03-01/signalr"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	cosmosParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePrivateEndpoint() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateEndpointCreate,
		Read:   resourcePrivateEndpointRead,
		Update: resourcePrivateEndpointUpdate,
		Delete: resourcePrivateEndpointDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privateendpoints.ParsePrivateEndpointID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.PrivateLinkName,
			},

			"location": commonschema.Location(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"subnet_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"network_interface": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"custom_network_interface_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"private_dns_zone_group": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.PrivateLinkName,
						},
						"private_dns_zone_ids": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
							},
						},
					},
				},
			},

			"private_service_connection": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.PrivateLinkName,
						},
						"is_manual_connection": {
							Type:     pluginsdk.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"private_connection_resource_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
							ExactlyOneOf: []string{"private_service_connection.0.private_connection_resource_alias", "private_service_connection.0.private_connection_resource_id"},
						},
						"private_connection_resource_alias": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validate.PrivateConnectionResourceAlias,
							ExactlyOneOf: []string{"private_service_connection.0.private_connection_resource_alias", "private_service_connection.0.private_connection_resource_id"},
						},
						"subresource_names": {
							Type:     pluginsdk.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.PrivateLinkSubResourceName,
							},
						},
						"request_message": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 140),
						},
						"private_ip_address": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ip_configuration": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validate.PrivateLinkName,
						},
						"private_ip_address": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"subresource_name": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						// lintignore:S013
						"member_name": {
							Type: pluginsdk.TypeString,
							// NOTE: O+C This value should remain optional computed as there are certain cases where Azure will error if you pass in a member id when it isn't expecting one.
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},
				},
			},

			"custom_dns_configs": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"fqdn": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"ip_addresses": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"private_dns_zone_configs": {
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
						"private_dns_zone_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"record_sets": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"type": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"fqdn": {
										Type:     pluginsdk.TypeString,
										Computed: true,
									},
									"ttl": {
										Type:     pluginsdk.TypeInt,
										Computed: true,
									},
									"ip_addresses": {
										Type:     pluginsdk.TypeList,
										Computed: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
										},
									},
								},
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourcePrivateEndpointCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privateendpoints.NewPrivateEndpointID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if err := validatePrivateEndpointSettings(d); err != nil {
		return fmt.Errorf("validating the configuration for %s: %+v", id, err)
	}

	existing, err := client.Get(ctx, id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if existing.Model != nil {
		return tf.ImportAsExistsError("azurerm_private_endpoint", id.ID())
	}

	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]interface{})

	parameters := privateendpoints.PrivateEndpoint{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &privateendpoints.PrivateEndpointProperties{
			PrivateLinkServiceConnections:       expandPrivateLinkEndpointServiceConnection(d.Get("private_service_connection").([]interface{}), false),
			ManualPrivateLinkServiceConnections: expandPrivateLinkEndpointServiceConnection(d.Get("private_service_connection").([]interface{}), true),
			Subnet: &privateendpoints.Subnet{
				Id: pointer.To(d.Get("subnet_id").(string)),
			},
			IPConfigurations:           expandPrivateEndpointIPConfigurations(d.Get("ip_configuration").([]interface{})),
			CustomNetworkInterfaceName: pointer.To(d.Get("custom_network_interface_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	err = validatePrivateLinkServiceId(*parameters.Properties.PrivateLinkServiceConnections)
	if err != nil {
		return err
	}
	err = validatePrivateLinkServiceId(*parameters.Properties.ManualPrivateLinkServiceConnections)
	if err != nil {
		return err
	}

	cosmosDbResIds := getCosmosDbResIdInPrivateServiceConnections(parameters.Properties)
	for _, cosmosDbResId := range cosmosDbResIds {
		log.Printf("[DEBUG] Add Lock For Private Endpoint %q, lock name: %q", id.PrivateEndpointName, cosmosDbResId)
		locks.ByName(cosmosDbResId, "azurerm_private_endpoint")
		//goland:noinspection GoDeferInLoop
		defer locks.UnlockByName(cosmosDbResId, "azurerm_private_endpoint")
	}

	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if err = client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
			switch {
			case strings.EqualFold(err.Error(), "is missing required parameter 'group Id'"):
				{
					return &pluginsdk.RetryError{
						Err:       fmt.Errorf("creating %s due to missing 'group Id', ensure that the 'subresource_names' type is populated: %+v", id, err),
						Retryable: false,
					}
				}
			case strings.Contains(err.Error(), "PrivateLinkServiceId Invalid private link service id"):
				{
					return &pluginsdk.RetryError{
						Err:       fmt.Errorf("creating Private Endpoint %s: %+v", id, err),
						Retryable: true,
					}
				}
			default:
				return &pluginsdk.RetryError{
					Err:       fmt.Errorf("creating %s: %+v", id, err),
					Retryable: false,
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	d.SetId(id.ID())

	// 1 Private Endpoint can have 1 Private DNS Zone Group
	// since this is a new resource, there shouldn't be an existing one - so there's no need to delete it
	if len(privateDnsZoneGroup) > 0 {
		log.Printf("[DEBUG] Creating Private DNS Zone Group associated with %s..", id)
		if err := createPrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, id, privateDnsZoneGroup); err != nil {
			return err
		}
		log.Printf("[DEBUG] Created the Existing Private DNS Zone Group associated with %s", id)
	}

	return resourcePrivateEndpointRead(d, meta)
}

func validatePrivateLinkServiceId(endpoints []privateendpoints.PrivateLinkServiceConnection) error {
	for _, connection := range endpoints {
		if connection.Properties == nil || connection.Properties.PrivateLinkServiceId == nil {
			return fmt.Errorf("properties/id was nil for %+v", connection)
		}
		_, errors := azure.ValidateResourceID(*connection.Properties.PrivateLinkServiceId, "PrivateLinkServiceID")
		if len(errors) == 0 {
			continue
		}
		_, errors = validate.PrivateConnectionResourceAlias(*connection.Properties.PrivateLinkServiceId, "PrivateLinkServiceID")
		if len(errors) != 0 {
			return fmt.Errorf("PrivateLinkServiceId Invalid: %q", *connection.Properties.PrivateLinkServiceId)
		}
	}
	return nil
}

func getCosmosDbResIdInPrivateServiceConnections(p *privateendpoints.PrivateEndpointProperties) []string {
	var ids []string
	exists := make(map[string]struct{})

	for _, l := range *p.PrivateLinkServiceConnections {
		if l.Properties.PrivateLinkServiceId == nil {
			continue
		}
		id := *l.Properties.PrivateLinkServiceId
		if _, err := cosmosParse.DatabaseAccountID(id); err == nil {
			_, ok := exists[id]
			if !ok {
				ids = append(ids, id)
				exists[id] = struct{}{}
			}
		}
	}
	for _, l := range *p.ManualPrivateLinkServiceConnections {
		if l.Properties.PrivateLinkServiceId == nil {
			continue
		}
		id := *l.Properties.PrivateLinkServiceId
		if _, err := cosmosParse.DatabaseAccountID(id); err == nil {
			_, ok := exists[id]
			if !ok {
				ids = append(ids, id)
				exists[id] = struct{}{}
			}
		}
	}
	// Sort ids, force adding lock in consistent order to avoid potential deadlock
	sort.Strings(ids)
	return ids
}

func resourcePrivateEndpointUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateendpoints.ParsePrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	if err := validatePrivateEndpointSettings(d); err != nil {
		return fmt.Errorf("validating the configuration for %s: %+v", id, err)
	}

	// Ensure we don't overwrite the existing ApplicationSecurityGroups
	existing, err := client.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model.Properties` was nil", *id)
	}

	applicationSecurityGroupAssociation := existing.Model.Properties.ApplicationSecurityGroups
	location := azure.NormalizeLocation(d.Get("location").(string))
	privateDnsZoneGroup := d.Get("private_dns_zone_group").([]interface{})
	privateServiceConnections := d.Get("private_service_connection").([]interface{})
	ipConfigurations := d.Get("ip_configuration").([]interface{})
	subnetId := d.Get("subnet_id").(string)
	customNicName := d.Get("custom_network_interface_name").(string)

	// TODO: in future it'd be nice to support conditional updates here, but one problem at a time
	parameters := privateendpoints.PrivateEndpoint{
		Location: pointer.To(location),
		Properties: &privateendpoints.PrivateEndpointProperties{
			ApplicationSecurityGroups:           applicationSecurityGroupAssociation,
			PrivateLinkServiceConnections:       expandPrivateLinkEndpointServiceConnection(privateServiceConnections, false),
			ManualPrivateLinkServiceConnections: expandPrivateLinkEndpointServiceConnection(privateServiceConnections, true),
			Subnet: &privateendpoints.Subnet{
				Id: pointer.To(subnetId),
			},
			IPConfigurations:           expandPrivateEndpointIPConfigurations(ipConfigurations),
			CustomNetworkInterfaceName: pointer.To(customNicName),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	err = validatePrivateLinkServiceId(*parameters.Properties.PrivateLinkServiceConnections)
	if err != nil {
		return err
	}
	err = validatePrivateLinkServiceId(*parameters.Properties.ManualPrivateLinkServiceConnections)
	if err != nil {
		return err
	}

	err = pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if err = client.CreateOrUpdateThenPoll(ctx, *id, parameters); err != nil {
			switch {
			case strings.EqualFold(err.Error(), "is missing required parameter 'group Id'"):
				{
					return &pluginsdk.RetryError{
						Err:       fmt.Errorf("updating %s due to missing 'group Id', ensure that the 'subresource_names' type is populated: %+v", id, err),
						Retryable: false,
					}
				}
			case strings.Contains(err.Error(), "PrivateLinkServiceId Invalid private link service id"):
				{
					return &pluginsdk.RetryError{
						Err:       fmt.Errorf("creating Private Endpoint %s: %+v", id, err),
						Retryable: true,
					}
				}
			default:
				return &pluginsdk.RetryError{
					Err: fmt.Errorf("updating %s: %+v", id, err),
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// 1 Private Endpoint can have 1 Private DNS Zone Group - so to update we need to Delete & Recreate
	if d.HasChange("private_dns_zone_group") {
		existingDnsZoneGroups, err := retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx, dnsClient, *id)
		if err != nil {
			return err
		}

		newDnsZoneGroups := d.Get("private_dns_zone_group").([]interface{})
		newDnsZoneName := ""
		idHasBeenChanged := false
		if len(newDnsZoneGroups) > 0 {
			groupRaw := newDnsZoneGroups[0].(map[string]interface{})
			newDnsZoneName = groupRaw["name"].(string)

			// it is possible to add or remove a private_dns_zone_id, but if an id is added at the same time as one as been removed and the name has not been changed
			// an existing entry is updated, which is not allowed, so we need to delete the existing private dns zone groups
			if d.HasChange("private_dns_zone_group.0.private_dns_zone_ids") {
				o, n := d.GetChange("private_dns_zone_group.0.private_dns_zone_ids")
				if len(o.([]interface{})) == len(n.([]interface{})) {
					idHasBeenChanged = true
				}
			}
		}

		needToRemove := newDnsZoneName == ""
		nameHasChanged := false
		if existingDnsZoneGroups != nil && newDnsZoneName != "" {
			needToRemove = len(*existingDnsZoneGroups) > 0 && len(newDnsZoneGroups) == 0

			// there should only be a single one, but there's no harm checking all returned
			for _, existing := range *existingDnsZoneGroups {
				if existing.PrivateDnsZoneGroupName != newDnsZoneName {
					nameHasChanged = true
					break
				}
			}
		}

		if needToRemove || nameHasChanged || idHasBeenChanged {
			log.Printf("[DEBUG] Deleting the Existing Private DNS Zone Group associated with %s..", id)
			if err := deletePrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, *id); err != nil {
				return err
			}
			log.Printf("[DEBUG] Deleted the Existing Private DNS Zone Group associated with %s.", id)
		}

		if len(privateDnsZoneGroup) > 0 {
			log.Printf("[DEBUG] Creating Private DNS Zone Group associated with %s..", id)
			if err := createPrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsClient, *id, privateDnsZoneGroup); err != nil {
				return err
			}
			log.Printf("[DEBUG] Created the Existing Private DNS Zone Group associated with %s", id)
		}
	}

	return resourcePrivateEndpointRead(d, meta)
}

func resourcePrivateEndpointRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	nicsClient := meta.(*clients.Client).Network.NetworkInterfaces
	dnsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateendpoints.ParsePrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Private Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", id, err)
	}

	privateDnsZoneIds, err := retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx, dnsClient, *id)
	if err != nil {
		return err
	}

	d.Set("name", id.PrivateEndpointName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			if err := d.Set("custom_dns_configs", flattenCustomDnsConfigs(props.CustomDnsConfigs)); err != nil {
				return fmt.Errorf("setting `custom_dns_configs`: %+v", err)
			}

			networkInterfaceId := ""
			privateIpAddress := ""
			if nics := props.NetworkInterfaces; nics != nil && len(*nics) > 0 {
				nic := (*nics)[0]
				if nic.Id != nil && *nic.Id != "" {
					networkInterfaceId = *nic.Id
					privateIpAddress = getPrivateIpAddress(ctx, nicsClient, networkInterfaceId)
				}
			}

			networkInterface := flattenNetworkInterface(networkInterfaceId)
			if err := d.Set("network_interface", networkInterface); err != nil {
				return fmt.Errorf("setting `network_interface`: %+v", err)
			}

			flattenedConnection := flattenPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections, privateIpAddress)
			if err := d.Set("private_service_connection", flattenedConnection); err != nil {
				return fmt.Errorf("setting `private_service_connection`: %+v", err)
			}

			flattenedipconfiguration := flattenPrivateEndpointIPConfigurations(props.IPConfigurations)
			if err := d.Set("ip_configuration", flattenedipconfiguration); err != nil {
				return fmt.Errorf("setting `ip_configuration`: %+v", err)
			}

			subnetId := ""
			if props.Subnet != nil && props.Subnet.Id != nil {
				subnetId = *props.Subnet.Id
			}
			d.Set("subnet_id", subnetId)
			customNicName := ""
			if props.CustomNetworkInterfaceName != nil {
				customNicName = *props.CustomNetworkInterfaceName
			}
			d.Set("custom_network_interface_name", customNicName)
		}

		privateDnsZoneConfigs := make([]interface{}, 0)
		privateDnsZoneGroups := make([]interface{}, 0)
		if privateDnsZoneIds != nil {
			for _, dnsZoneId := range *privateDnsZoneIds {
				flattened, err := retrieveAndFlattenPrivateDnsZone(ctx, dnsClient, dnsZoneId)
				if err != nil {
					return fmt.Errorf("reading %s for %s: %+v", dnsZoneId, id, err)
				}

				// an exceptional case but no harm in handling
				if flattened == nil {
					continue
				}

				privateDnsZoneConfigs = append(privateDnsZoneConfigs, flattened.DnsZoneConfig...)
				privateDnsZoneGroups = append(privateDnsZoneGroups, flattened.DnsZoneGroup)
			}
		}
		if err = d.Set("private_dns_zone_configs", privateDnsZoneConfigs); err != nil {
			return fmt.Errorf("setting `private_dns_zone_configs`: %+v", err)
		}
		if err = d.Set("private_dns_zone_group", privateDnsZoneGroups); err != nil {
			return fmt.Errorf("setting `private_dns_zone_group`: %+v", err)
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourcePrivateEndpointDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateEndpoints
	dnsZoneGroupsClient := meta.(*clients.Client).Network.PrivateDnsZoneGroups
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privateendpoints.ParsePrivateEndpointID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting the Private DNS Zone Group associated with %s", id)
	if err := deletePrivateDnsZoneGroupForPrivateEndpoint(ctx, dnsZoneGroupsClient, *id); err != nil {
		return err
	}
	log.Printf("[DEBUG] Deleted the Private DNS Zone Group associated with %s.", id)

	existing, err := client.Get(ctx, *id, privateendpoints.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil {
		return fmt.Errorf("retrieving existing %s: `model` was nil", *id)
	}
	subnetId := ""
	if model := existing.Model; model != nil {
		if props := model.Properties; props != nil {
			if subnet := props.Subnet; subnet != nil && subnet.Id != nil {
				subnetId = *subnet.Id
			}
		}
	}
	if subnetId == "" {
		// this also captures `model.Properties` being nil below, since otherwise we wouldn't get the Subnet
		return fmt.Errorf("retrieving existing %s: `model.Properties.Subnet.Id` was nil", *id)
	}

	cosmosDbResIds := getCosmosDbResIdInPrivateServiceConnections(existing.Model.Properties)
	for _, cosmosDbResId := range cosmosDbResIds {
		locks.ByName(cosmosDbResId, "azurerm_private_endpoint")
		//goland:noinspection GoDeferInLoop
		defer locks.UnlockByName(cosmosDbResId, "azurerm_private_endpoint")
	}

	log.Printf("[DEBUG] Deleting %s", id)
	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	log.Printf("[DEBUG] Deleted %s", id)

	return nil
}

func expandPrivateLinkEndpointServiceConnection(input []interface{}, parseManual bool) *[]privateendpoints.PrivateLinkServiceConnection {
	results := make([]privateendpoints.PrivateLinkServiceConnection, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		privateConnectionResourceId := v["private_connection_resource_id"].(string)
		if privateConnectionResourceId == "" {
			privateConnectionResourceId = v["private_connection_resource_alias"].(string)
		}
		subresourceNames := v["subresource_names"].([]interface{})
		requestMessage := v["request_message"].(string)
		isManual := v["is_manual_connection"].(bool)
		name := v["name"].(string)

		if isManual == parseManual {
			result := privateendpoints.PrivateLinkServiceConnection{
				Name: pointer.To(name),
				Properties: &privateendpoints.PrivateLinkServiceConnectionProperties{
					GroupIds:             utils.ExpandStringSlice(subresourceNames),
					PrivateLinkServiceId: pointer.To(privateConnectionResourceId),
				},
			}

			if requestMessage != "" {
				result.Properties.RequestMessage = pointer.To(requestMessage)
			}

			results = append(results, result)
		}
	}

	return &results
}

func expandPrivateEndpointIPConfigurations(input []interface{}) *[]privateendpoints.PrivateEndpointIPConfiguration {
	results := make([]privateendpoints.PrivateEndpointIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		privateIPAddress := v["private_ip_address"].(string)
		subResourceName := v["subresource_name"].(string)
		memberName := v["member_name"].(string)
		if memberName == "" {
			memberName = subResourceName
		}
		name := v["name"].(string)
		result := privateendpoints.PrivateEndpointIPConfiguration{
			Name: pointer.To(name),
			Properties: &privateendpoints.PrivateEndpointIPConfigurationProperties{
				PrivateIPAddress: pointer.To(privateIPAddress),
				GroupId:          pointer.To(subResourceName),
				MemberName:       pointer.To(memberName),
			},
		}
		results = append(results, result)
	}

	return &results
}

func flattenPrivateEndpointIPConfigurations(ipConfigurations *[]privateendpoints.PrivateEndpointIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if ipConfigurations == nil {
		return results
	}

	for _, item := range *ipConfigurations {
		if props := item.Properties; props != nil {
			results = append(results, map[string]interface{}{
				"name":               item.Name,
				"private_ip_address": props.PrivateIPAddress,
				"subresource_name":   props.GroupId,
				"member_name":        props.MemberName,
			})
		}
	}

	return results
}

func flattenCustomDnsConfigs(customDnsConfigs *[]privateendpoints.CustomDnsConfigPropertiesFormat) []interface{} {
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

func flattenPrivateLinkEndpointServiceConnection(serviceConnections *[]privateendpoints.PrivateLinkServiceConnection, manualServiceConnections *[]privateendpoints.PrivateLinkServiceConnection, privateIPAddress string) []interface{} {
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

			if props := item.Properties; props != nil {
				if v := props.GroupIds; v != nil {
					subResourceNames = utils.FlattenStringSlice(v)
				}
				if props.PrivateLinkServiceId != nil {
					privateConnectionId = *props.PrivateLinkServiceId
				}
			}
			attrs := map[string]interface{}{
				"name":                 name,
				"is_manual_connection": false,
				"private_ip_address":   privateIPAddress,
				"subresource_names":    subResourceNames,
			}
			if strings.HasSuffix(privateConnectionId, ".azure.privatelinkservice") {
				attrs["private_connection_resource_alias"] = privateConnectionId
			} else {
				privateConnectionId = normalizePrivateConnectionId(privateConnectionId)
				attrs["private_connection_resource_id"] = privateConnectionId
			}

			results = append(results, attrs)
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

			if props := item.Properties; props != nil {
				if v := props.GroupIds; v != nil {
					subResourceNames = utils.FlattenStringSlice(v)
				}
				if props.PrivateLinkServiceId != nil {
					privateConnectionId = *props.PrivateLinkServiceId
				}
				if props.RequestMessage != nil {
					requestMessage = *props.RequestMessage
				}
			}

			attrs := map[string]interface{}{
				"name":                 name,
				"is_manual_connection": true,
				"private_ip_address":   privateIPAddress,
				"request_message":      requestMessage,
				"subresource_names":    subResourceNames,
			}
			if strings.HasSuffix(privateConnectionId, ".azure.privatelinkservice") {
				attrs["private_connection_resource_alias"] = privateConnectionId
			} else {
				privateConnectionId = normalizePrivateConnectionId(privateConnectionId)
				attrs["private_connection_resource_id"] = privateConnectionId
			}

			results = append(results, attrs)
		}
	}

	return results
}

func createPrivateDnsZoneGroupForPrivateEndpoint(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privateendpoints.PrivateEndpointId, inputRaw []interface{}) error {
	if len(inputRaw) != 1 {
		return fmt.Errorf("expected a single Private DNS Zone Groups but got %d", len(inputRaw))
	}
	item := inputRaw[0].(map[string]interface{})

	dnsZoneGroupId := privatednszonegroups.NewPrivateDnsZoneGroupID(id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName, item["name"].(string))
	privateDnsZoneIdsRaw := item["private_dns_zone_ids"].([]interface{})
	privateDnsZoneConfigs := make([]privatednszonegroups.PrivateDnsZoneConfig, 0)
	for _, item := range privateDnsZoneIdsRaw {
		v := item.(string)

		privateDnsZone, err := privatezones.ParsePrivateDnsZoneID(v)
		if err != nil {
			return err
		}

		privateDnsZoneConfigs = append(privateDnsZoneConfigs, privatednszonegroups.PrivateDnsZoneConfig{
			Name: pointer.To(privateDnsZone.PrivateDnsZoneName),
			Properties: &privatednszonegroups.PrivateDnsZonePropertiesFormat{
				PrivateDnsZoneId: pointer.To(privateDnsZone.ID()),
			},
		})
	}

	parameters := privatednszonegroups.PrivateDnsZoneGroup{
		Name: pointer.To(id.PrivateEndpointName),
		Properties: &privatednszonegroups.PrivateDnsZoneGroupPropertiesFormat{
			PrivateDnsZoneConfigs: &privateDnsZoneConfigs,
		},
	}
	if err := client.CreateOrUpdateThenPoll(ctx, dnsZoneGroupId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	return nil
}

func deletePrivateDnsZoneGroupForPrivateEndpoint(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privateendpoints.PrivateEndpointId) error {
	// lookup and delete the (should be, Single) Private DNS Zone Group associated with this Private Endpoint
	privateDnsZoneIds, err := retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx, client, id)
	if err != nil {
		return err
	}

	for _, privateDnsZoneId := range *privateDnsZoneIds {
		log.Printf("[DEBUG] Deleting %s..", privateDnsZoneId)
		if err := client.DeleteThenPoll(ctx, privateDnsZoneId); err != nil {
			return fmt.Errorf("deleting %s: %+v", privateDnsZoneId, err)
		}
	}

	return nil
}

func retrievePrivateDnsZoneGroupsForPrivateEndpoint(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privateendpoints.PrivateEndpointId) (*[]privatednszonegroups.PrivateDnsZoneGroupId, error) {
	output := make([]privatednszonegroups.PrivateDnsZoneGroupId, 0)

	privateEndpointId := privatednszonegroups.NewPrivateEndpointID(id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName)
	dnsZones, err := client.ListComplete(ctx, privateEndpointId)
	if err != nil {
		if response.WasNotFound(dnsZones.LatestHttpResponse) {
			return &output, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	for _, zone := range dnsZones.Items {
		if zone.Id != nil {
			groupId, err := privatednszonegroups.ParsePrivateDnsZoneGroupID(*zone.Id)
			if err != nil {
				return nil, err
			}

			output = append(output, *groupId)
		}
	}

	return &output, nil
}

type flattenedPrivateDnsZoneGroup struct {
	DnsZoneConfig []interface{}
	DnsZoneGroup  map[string]interface{}
}

func retrieveAndFlattenPrivateDnsZone(ctx context.Context, client *privatednszonegroups.PrivateDnsZoneGroupsClient, id privatednszonegroups.PrivateDnsZoneGroupId) (*flattenedPrivateDnsZoneGroup, error) {
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}

	privateDnsZoneIds := make([]string, 0)
	dnsZoneConfigs := make([]interface{}, 0)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil && props.PrivateDnsZoneConfigs != nil {
			for _, config := range *props.PrivateDnsZoneConfigs {
				if config.Name == nil {
					// necessary to build up the ID
					continue
				}
				if config.Properties == nil || config.Properties.PrivateDnsZoneId == nil {
					// necessary for a bunch of other things
					continue
				}
				configProps := *config.Properties
				name := *config.Name
				privateDnsZoneId := *configProps.PrivateDnsZoneId

				privateDnsZoneIds = append(privateDnsZoneIds, privateDnsZoneId)

				recordSets := flattenPrivateDnsZoneGroupRecordSets(configProps.RecordSets)
				dnsZoneConfigs = append(dnsZoneConfigs, map[string]interface{}{
					"id":                  parse.NewPrivateDnsZoneConfigID(id.SubscriptionId, id.ResourceGroupName, id.PrivateEndpointName, id.PrivateDnsZoneGroupName, name).ID(),
					"name":                name,
					"private_dns_zone_id": privateDnsZoneId,
					"record_sets":         recordSets,
				})
			}
		}
	}

	return &flattenedPrivateDnsZoneGroup{
		DnsZoneConfig: dnsZoneConfigs,
		DnsZoneGroup: map[string]interface{}{
			"id":                   id.ID(),
			"name":                 id.PrivateDnsZoneGroupName,
			"private_dns_zone_ids": privateDnsZoneIds,
		},
	}, nil
}

func flattenPrivateDnsZoneGroupRecordSets(input *[]privatednszonegroups.RecordSet) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	for _, v := range *input {
		fqdn := ""
		if v.Fqdn != nil {
			fqdn = *v.Fqdn
		}

		name := ""
		if v.RecordSetName != nil {
			name = *v.RecordSetName
		}

		recordType := ""
		if v.RecordType != nil {
			recordType = *v.RecordType
		}

		ttl := 0
		if v.Ttl != nil {
			ttl = int(*v.Ttl)
		}

		ipAddresses := make([]string, 0)
		if v.IPAddresses != nil {
			ipAddresses = *v.IPAddresses
		}

		output = append(output, map[string]interface{}{
			"fqdn":         fqdn,
			"ip_addresses": ipAddresses,
			"name":         name,
			"ttl":          ttl,
			"type":         recordType,
		})
	}

	return output
}

func validatePrivateEndpointSettings(d *pluginsdk.ResourceData) error {
	privateServiceConnections := d.Get("private_service_connection").([]interface{})

	for _, psc := range privateServiceConnections {
		privateServiceConnection := psc.(map[string]interface{})
		name := privateServiceConnection["name"].(string)

		// If this is not a manual connection and the message is set return an error since this does not make sense.
		if !privateServiceConnection["is_manual_connection"].(bool) && privateServiceConnection["request_message"].(string) != "" {
			return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute cannot be set if the "is_manual_connection" attribute is "false"`, name)
		}

		// If this is a manual connection and the message isn't set return an error.
		if privateServiceConnection["is_manual_connection"].(bool) && strings.TrimSpace(privateServiceConnection["request_message"].(string)) == "" {
			return fmt.Errorf(`"private_service_connection":%q is invalid, the "request_message" attribute must not be empty`, name)
		}
	}

	return nil
}

// normalize the PrivateConnectionId due to the casing change at service side
func normalizePrivateConnectionId(privateConnectionId string) string {
	// intentionally including the extra segment to handle Redis vs Redis Enterprise (which is within the same RP)
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.cache/redis/") {
		if cacheId, err := redis.ParseRediIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = cacheId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.dbforpostgresql") {
		if serverId, err := postgresqlServers.ParseServerIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serverId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.dbformysql") {
		if serverId, err := servers.ParseServerIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serverId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.dbformariadb") {
		if serverId, err := mariadbServers.ParseServerIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serverId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.kusto") {
		if clusterId, err := commonids.ParseKustoClusterIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = clusterId.ID()
		}
	}
	if strings.Contains(strings.ToLower(privateConnectionId), "microsoft.signalrservice") {
		if serviceId, err := signalr.ParseSignalRIDInsensitively(privateConnectionId); err == nil {
			privateConnectionId = serviceId.ID()
		}
	}
	return privateConnectionId
}
