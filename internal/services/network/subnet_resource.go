// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicies"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/ipampools"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

//go:generate go run ../../tools/generator-tests resourceidentity -resource-name subnet -service-package-name network -properties "name,resource_group_name,virtual_network_name" -known-values "subscription_id:data.Subscriptions.Primary"

var SubnetResourceName = "azurerm_subnet"

var subnetDelegationServiceNames = []string{
	"GitHub.Network/networkSettings",
	"Informatica.DataManagement/organizations",
	"Microsoft.ApiManagement/service",
	"Microsoft.Apollo/npu",
	"Microsoft.App/environments",
	"Microsoft.App/testClients",
	"Microsoft.AVS/PrivateClouds",
	"Microsoft.AzureCosmosDB/clusters",
	"Microsoft.BareMetal/AzureHostedService",
	"Microsoft.BareMetal/AzureHPC",
	"Microsoft.BareMetal/AzurePaymentHSM",
	"Microsoft.BareMetal/AzureVMware",
	"Microsoft.BareMetal/CrayServers",
	"Microsoft.BareMetal/MonitoringServers",
	"Microsoft.Batch/batchAccounts",
	"Microsoft.CloudTest/hostedpools",
	"Microsoft.CloudTest/images",
	"Microsoft.CloudTest/pools",
	"Microsoft.Codespaces/plans",
	"Microsoft.ContainerInstance/containerGroups",
	"Microsoft.ContainerService/managedClusters",
	"Microsoft.ContainerService/TestClients",
	"Microsoft.Databricks/workspaces",
	"Microsoft.DBforMySQL/flexibleServers",
	"Microsoft.DBforMySQL/servers",
	"Microsoft.DBforMySQL/serversv2",
	"Microsoft.DBforPostgreSQL/flexibleServers",
	"Microsoft.DBforPostgreSQL/serversv2",
	"Microsoft.DBforPostgreSQL/singleServers",
	"Microsoft.DelegatedNetwork/controller",
	"Microsoft.DevCenter/networkConnection",
	"Microsoft.DevOpsInfrastructure/pools",
	"Microsoft.DocumentDB/cassandraClusters",
	"Microsoft.Fidalgo/networkSettings",
	"Microsoft.HardwareSecurityModules/dedicatedHSMs",
	"Microsoft.Kusto/clusters",
	"Microsoft.LabServices/labplans",
	"Microsoft.Logic/integrationServiceEnvironments",
	"Microsoft.MachineLearningServices/workspaces",
	"Microsoft.Netapp/volumes",
	"Microsoft.Network/applicationGateways",
	"Microsoft.Network/dnsResolvers",
	"Microsoft.Network/managedResolvers",
	"Microsoft.Network/fpgaNetworkInterfaces",
	"Microsoft.Network/networkWatchers.",
	"Microsoft.Network/virtualNetworkGateways",
	"Microsoft.Orbital/orbitalGateways",
	"Microsoft.PowerAutomate/hostedRpa",
	"Microsoft.PowerPlatform/enterprisePolicies",
	"Microsoft.PowerPlatform/vnetaccesslinks",
	"Microsoft.ServiceFabricMesh/networks",
	"Microsoft.ServiceNetworking/trafficControllers",
	"Microsoft.Singularity/accounts/networks",
	"Microsoft.Singularity/accounts/npu",
	"Microsoft.Sql/managedInstances",
	"Microsoft.Sql/managedInstancesOnebox",
	"Microsoft.Sql/managedInstancesStage",
	"Microsoft.Sql/managedInstancesTest",
	"Microsoft.Sql/servers",
	"Microsoft.StoragePool/diskPools",
	"Microsoft.StreamAnalytics/streamingJobs",
	"Microsoft.Synapse/workspaces",
	"Microsoft.Web/hostingEnvironments",
	"Microsoft.Web/serverFarms",
	"NGINX.NGINXPLUS/nginxDeployments",
	"Oracle.Database/networkAttachments",
	"PaloAltoNetworks.Cloudngfw/firewalls",
	"Qumulo.Storage/fileSystems",
}

func resourceSubnet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:   resourceSubnetCreate,
		Read:     resourceSubnetRead,
		Update:   resourceSubnetUpdate,
		Delete:   resourceSubnetDelete,
		Importer: pluginsdk.ImporterValidatingIdentity(&commonids.SubnetId{}),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
				// Validate `sharing_scope` cannot be set when `default_outbound_access_enabled` is true.
				if diff.Get("sharing_scope").(string) != "" && diff.Get("default_outbound_access_enabled").(bool) {
					return fmt.Errorf("`sharing_scope` cannot be set if `default_outbound_access_enabled` is set to `true`")
				}
				return nil
			}),
		),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&commonids.SubnetId{}),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"virtual_network_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_prefixes": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MinItems:     1,
				ExactlyOneOf: []string{"address_prefixes", "ip_address_pool"},
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					// If `ip_address_pool` is used instead of `address_prefixes` there is a perpetual diff
					// due to the API returning a CIDR range provisioned by the IP Address Management Pool.
					// Note: using `GetRawConfig` to avoid suppressing a diff if a user updates from `ip_address_pool` to `address_prefixes`.
					rawIpAddressPool := d.GetRawConfig().AsValueMap()["ip_address_pool"]
					if !rawIpAddressPool.IsNull() && len(rawIpAddressPool.AsValueSlice()) > 0 {
						return true
					}

					return false
				},
			},

			"service_endpoints": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
				Set:      pluginsdk.HashString,
			},

			"service_endpoint_policy_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: serviceendpointpolicies.ValidateServiceEndpointPolicyID,
				},
			},

			"sharing_scope": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// todo "Tenant" is only supported until https://github.com/Azure/azure-rest-api-specs/issues/36446 is addressed
				ValidateFunc: validation.StringInSlice([]string{string(subnets.SharingScopeTenant)}, false),
			},

			"delegation": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
						"service_delegation": {
							Type:     pluginsdk.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringInSlice(subnetDelegationServiceNames, false),
									},

									"actions": {
										Type:     pluginsdk.TypeSet,
										Optional: true,
										Elem: &pluginsdk.Schema{
											Type: pluginsdk.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Microsoft.Network/networkinterfaces/*",
												"Microsoft.Network/publicIPAddresses/join/action",
												"Microsoft.Network/publicIPAddresses/read",
												"Microsoft.Network/virtualNetworks/read",
												"Microsoft.Network/virtualNetworks/subnets/action",
												"Microsoft.Network/virtualNetworks/subnets/join/action",
												"Microsoft.Network/virtualNetworks/subnets/prepareNetworkPolicies/action",
												"Microsoft.Network/virtualNetworks/subnets/unprepareNetworkPolicies/action",
											}, false),
										},
									},
								},
							},
						},
					},
				},
			},

			"ip_address_pool": {
				Type:         pluginsdk.TypeList,
				Optional:     true,
				MaxItems:     1,
				ExactlyOneOf: []string{"address_prefixes", "ip_address_pool"},
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: ipampools.ValidateIPamPoolID,
						},

						"number_of_ip_addresses": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile(`^[1-9]\d*$`),
								"`number_of_ip_addresses` must be a string that represents a positive number",
							),
						},

						"allocated_ip_address_prefixes": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"default_outbound_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				Optional: true,
			},

			"private_endpoint_network_policies": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled),
				ValidateFunc: validation.StringInSlice(subnets.PossibleValuesForVirtualNetworkPrivateEndpointNetworkPolicies(), false),
			},

			"private_link_service_network_policies_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

// TODO: refactor the create/flatten functions
func resourceSubnetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Subnet creation.")

	id := commonids.NewSubnetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_subnet", id.ID())
	}

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	properties := subnets.SubnetPropertiesFormat{}
	if value, ok := d.GetOk("address_prefixes"); ok {
		var addressPrefixes []string
		for _, item := range value.([]interface{}) {
			addressPrefixes = append(addressPrefixes, item.(string))
		}
		properties.AddressPrefixes = &addressPrefixes
	}
	if properties.AddressPrefixes != nil && len(*properties.AddressPrefixes) == 1 {
		properties.AddressPrefix = &(*properties.AddressPrefixes)[0]
		properties.AddressPrefixes = nil
	}

	properties.IPamPoolPrefixAllocations = expandSubnetIPAddressPool(d.Get("ip_address_pool").([]interface{}))

	// To enable private endpoints you must disable the network policies for the subnet because
	// Network policies like network security groups are not supported by private endpoints.
	var privateEndpointNetworkPolicies subnets.VirtualNetworkPrivateEndpointNetworkPolicies
	var privateLinkServiceNetworkPolicies subnets.VirtualNetworkPrivateLinkServiceNetworkPolicies

	privateEndpointNetworkPoliciesRaw := d.Get("private_endpoint_network_policies").(string)
	privateLinkServiceNetworkPoliciesRaw := d.Get("private_link_service_network_policies_enabled").(bool)

	privateEndpointNetworkPolicies = subnets.VirtualNetworkPrivateEndpointNetworkPolicies(privateEndpointNetworkPoliciesRaw)
	privateLinkServiceNetworkPolicies = subnets.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandSubnetNetworkPolicy(privateLinkServiceNetworkPoliciesRaw))

	properties.PrivateEndpointNetworkPolicies = pointer.To(privateEndpointNetworkPolicies)
	properties.PrivateLinkServiceNetworkPolicies = pointer.To(privateLinkServiceNetworkPolicies)

	serviceEndpointPoliciesRaw := d.Get("service_endpoint_policy_ids").(*pluginsdk.Set).List()
	properties.ServiceEndpointPolicies = expandSubnetServiceEndpointPolicies(serviceEndpointPoliciesRaw)

	serviceEndpointsRaw := d.Get("service_endpoints").(*pluginsdk.Set).List()
	properties.ServiceEndpoints = expandSubnetServiceEndpoints(serviceEndpointsRaw)

	properties.SharingScope = pointer.ToEnum[subnets.SharingScope](d.Get("sharing_scope").(string))

	properties.DefaultOutboundAccess = pointer.To(d.Get("default_outbound_access_enabled").(bool))

	delegationsRaw := d.Get("delegation").([]interface{})
	properties.Delegations = expandSubnetDelegation(delegationsRaw)

	subnet := subnets.Subnet{
		Name:       utils.String(id.SubnetName),
		Properties: &properties,
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, subnet); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceSubnetRead(d, meta)
}

func resourceSubnetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	vnetClient := meta.(*clients.Client).Network.VirtualNetworks
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	locks.ByName(id.SubnetName, SubnetResourceName)
	defer locks.UnlockByName(id.SubnetName, SubnetResourceName)

	existing, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", *id)
	}

	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// TODO: locking on the NSG/Route Table if applicable

	props := *existing.Model.Properties

	if d.HasChange("address_prefixes") {
		addressPrefixesRaw := d.Get("address_prefixes").([]interface{})
		switch len(addressPrefixesRaw) {
		case 0:
			// this is the case IPAddressPool is used, so we shall clear the `AddressPrefix` and `AddressPrefixes`.
			props.AddressPrefix = nil
			props.AddressPrefixes = nil
		case 1:
			// N->1: we shall insist on using the `AddressPrefix` and clear the `AddressPrefixes`.
			props.AddressPrefix = utils.String(addressPrefixesRaw[0].(string))
			props.AddressPrefixes = nil
		default:
			// 1->N: we shall insist on using the `AddressPrefixes` and clear the `AddressPrefix`. If both are set, service be confused and (currently) will only
			// return the `AddressPrefix` in response.
			props.AddressPrefixes = utils.ExpandStringSlice(addressPrefixesRaw)
			props.AddressPrefix = nil
		}
	}

	if d.HasChange("ip_address_pool") {
		if v := d.Get("ip_address_pool").([]interface{}); len(v) > 0 {
			expandedIPAddressPool := expandSubnetIPAddressPool(d.Get("ip_address_pool").([]interface{}))

			if props.IPamPoolPrefixAllocations != nil {
				for _, existingAllocation := range *props.IPamPoolPrefixAllocations {
					for _, expandedAllocation := range *expandedIPAddressPool {
						if existingAllocation.Pool != nil && expandedAllocation.Pool != nil && strings.EqualFold(pointer.From(existingAllocation.Pool.Id), pointer.From(expandedAllocation.Pool.Id)) &&
							existingAllocation.NumberOfIPAddresses != nil && expandedAllocation.NumberOfIPAddresses != nil && *existingAllocation.NumberOfIPAddresses > *expandedAllocation.NumberOfIPAddresses {
							return fmt.Errorf("`number_of_ip_addresses` cannot be decreased from %v to %v on pool: %v", *existingAllocation.NumberOfIPAddresses, *expandedAllocation.NumberOfIPAddresses, *expandedAllocation.Pool.Id)
						}
					}
				}
			}

			props.IPamPoolPrefixAllocations = expandedIPAddressPool

			// Set nil for AddressPrefixes when changing from `AddressPrefixes` to `IPAddressPool` but the change for `AddressPrefix` is not detected due to the diffSuppressFunc.
			props.AddressPrefix = nil
			props.AddressPrefixes = nil
		} else {
			props.IPamPoolPrefixAllocations = nil
		}
	}

	if d.HasChange("default_outbound_access_enabled") {
		props.DefaultOutboundAccess = pointer.To(d.Get("default_outbound_access_enabled").(bool))
	}

	if d.HasChange("delegation") {
		delegationsRaw := d.Get("delegation").([]interface{})
		props.Delegations = expandSubnetDelegation(delegationsRaw)
	}

	if d.HasChange("private_endpoint_network_policies") {
		v := d.Get("private_endpoint_network_policies").(string)
		props.PrivateEndpointNetworkPolicies = pointer.To(subnets.VirtualNetworkPrivateEndpointNetworkPolicies(v))
	}

	if d.HasChange("private_link_service_network_policies_enabled") {
		v := d.Get("private_link_service_network_policies_enabled").(bool)
		props.PrivateLinkServiceNetworkPolicies = pointer.To(subnets.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandSubnetNetworkPolicy(v)))
	}

	if d.HasChange("sharing_scope") {
		props.SharingScope = pointer.ToEnum[subnets.SharingScope](d.Get("sharing_scope").(string))
	}

	if d.HasChange("service_endpoints") {
		serviceEndpointsRaw := d.Get("service_endpoints").(*pluginsdk.Set).List()
		props.ServiceEndpoints = expandSubnetServiceEndpoints(serviceEndpointsRaw)
	}

	if d.HasChange("service_endpoint_policy_ids") {
		serviceEndpointPoliciesRaw := d.Get("service_endpoint_policy_ids").(*pluginsdk.Set).List()
		props.ServiceEndpointPolicies = expandSubnetServiceEndpointPolicies(serviceEndpointPoliciesRaw)
	}

	subnet := subnets.Subnet{
		Name:       utils.String(id.SubnetName),
		Properties: &props,
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, subnet); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(subnets.ProvisioningStateUpdating)},
		Target:     []string{string(subnets.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}

	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	return resourceSubnetRead(d, meta)
}

func resourceSubnetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, subnets.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SubnetName)
	d.Set("virtual_network_name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			if props.AddressPrefixes == nil {
				if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
					d.Set("address_prefixes", []string{*props.AddressPrefix})
				} else {
					d.Set("address_prefixes", []string{})
				}
			} else {
				d.Set("address_prefixes", props.AddressPrefixes)
			}

			defaultOutboundAccessEnabled := true
			if props.DefaultOutboundAccess != nil {
				defaultOutboundAccessEnabled = *props.DefaultOutboundAccess
			}
			d.Set("default_outbound_access_enabled", defaultOutboundAccessEnabled)

			delegation := flattenSubnetDelegation(props.Delegations)
			if err := d.Set("delegation", delegation); err != nil {
				return fmt.Errorf("flattening `delegation`: %+v", err)
			}

			if err := d.Set("ip_address_pool", flattenSubnetIPAddressPool(props.IPamPoolPrefixAllocations)); err != nil {
				return fmt.Errorf("setting `ip_address_pool`: %+v", err)
			}

			d.Set("private_endpoint_network_policies", string(pointer.From(props.PrivateEndpointNetworkPolicies)))
			d.Set("private_link_service_network_policies_enabled", flattenSubnetNetworkPolicy(string(pointer.From(props.PrivateLinkServiceNetworkPolicies))))
			d.Set("sharing_scope", pointer.FromEnum(props.SharingScope))

			serviceEndpoints := flattenSubnetServiceEndpoints(props.ServiceEndpoints)
			if err := d.Set("service_endpoints", serviceEndpoints); err != nil {
				return fmt.Errorf("setting `service_endpoints`: %+v", err)
			}

			serviceEndpointPolicies := flattenSubnetServiceEndpointPolicies(props.ServiceEndpointPolicies)
			if err := d.Set("service_endpoint_policy_ids", serviceEndpointPolicies); err != nil {
				return fmt.Errorf("setting `service_endpoint_policy_ids`: %+v", err)
			}
		}
	}

	if err := pluginsdk.SetResourceIdentityData(d, id); err != nil {
		return err
	}

	return nil
}

func resourceSubnetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.Subnets
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	locks.ByName(id.SubnetName, SubnetResourceName)
	defer locks.UnlockByName(id.SubnetName, SubnetResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandSubnetServiceEndpoints(input []interface{}) *[]subnets.ServiceEndpointPropertiesFormat {
	endpoints := make([]subnets.ServiceEndpointPropertiesFormat, 0)

	for _, svcEndpointRaw := range input {
		if svc, ok := svcEndpointRaw.(string); ok {
			endpoint := subnets.ServiceEndpointPropertiesFormat{
				Service: &svc,
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return &endpoints
}

func flattenSubnetServiceEndpoints(serviceEndpoints *[]subnets.ServiceEndpointPropertiesFormat) []interface{} {
	endpoints := make([]interface{}, 0)

	if serviceEndpoints == nil {
		return endpoints
	}

	for _, endpoint := range *serviceEndpoints {
		if endpoint.Service != nil {
			endpoints = append(endpoints, *endpoint.Service)
		}
	}

	return endpoints
}

func expandSubnetDelegation(input []interface{}) *[]subnets.Delegation {
	retDelegations := make([]subnets.Delegation, 0)

	for _, deleValue := range input {
		deleData := deleValue.(map[string]interface{})
		deleName := deleData["name"].(string)
		srvDelegations := deleData["service_delegation"].([]interface{})
		srvDelegation := srvDelegations[0].(map[string]interface{})
		srvName := srvDelegation["name"].(string)
		srvActions := srvDelegation["actions"].(*pluginsdk.Set).List()

		retSrvActions := make([]string, 0)
		for _, srvAction := range srvActions {
			srvActionData := srvAction.(string)
			retSrvActions = append(retSrvActions, srvActionData)
		}

		retDelegation := subnets.Delegation{
			Name: &deleName,
			Properties: &subnets.ServiceDelegationPropertiesFormat{
				ServiceName: &srvName,
				Actions:     &retSrvActions,
			},
		}

		retDelegations = append(retDelegations, retDelegation)
	}

	return &retDelegations
}

func flattenSubnetDelegation(delegations *[]subnets.Delegation) []interface{} {
	if delegations == nil {
		return []interface{}{}
	}

	retDeles := make([]interface{}, 0)

	normalizeServiceName := map[string]string{}
	for _, normName := range subnetDelegationServiceNames {
		normalizeServiceName[strings.ToLower(normName)] = normName
	}

	for _, dele := range *delegations {
		retDele := make(map[string]interface{})
		if v := dele.Name; v != nil {
			retDele["name"] = *v
		}

		svcDeles := make([]interface{}, 0)
		svcDele := make(map[string]interface{})
		if props := dele.Properties; props != nil {
			if v := props.ServiceName; v != nil {
				name := *v
				if nv, ok := normalizeServiceName[strings.ToLower(name)]; ok {
					name = nv
				}
				svcDele["name"] = name
			}

			if v := props.Actions; v != nil {
				svcDele["actions"] = *v
			}
		}

		svcDeles = append(svcDeles, svcDele)

		retDele["service_delegation"] = svcDeles

		retDeles = append(retDeles, retDele)
	}

	return retDeles
}

func expandSubnetNetworkPolicy(enabled bool) string {
	if enabled {
		return string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled)
	}

	return string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled)
}

func flattenSubnetNetworkPolicy(input string) bool {
	return strings.EqualFold(input, string(subnets.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled))
}

func expandSubnetServiceEndpointPolicies(input []interface{}) *[]subnets.ServiceEndpointPolicy {
	output := make([]subnets.ServiceEndpointPolicy, 0)
	for _, policy := range input {
		policy := policy.(string)
		output = append(output, subnets.ServiceEndpointPolicy{Id: &policy})
	}
	return &output
}

func flattenSubnetServiceEndpointPolicies(input *[]subnets.ServiceEndpointPolicy) []interface{} {
	if input == nil {
		return nil
	}

	output := make([]interface{}, 0, len(*input))
	for _, policy := range *input {
		id := ""
		if policy.Id != nil {
			id = *policy.Id
		}
		output = append(output, id)
	}
	return output
}

func expandSubnetIPAddressPool(input []interface{}) *[]subnets.IPamPoolPrefixAllocation {
	if len(input) == 0 {
		return nil
	}

	outputs := make([]subnets.IPamPoolPrefixAllocation, 0)
	for _, v := range input {
		ipPoolRaw := v.(map[string]interface{})
		output := subnets.IPamPoolPrefixAllocation{}

		if v, ok := ipPoolRaw["number_of_ip_addresses"]; ok {
			output.NumberOfIPAddresses = pointer.To(v.(string))
		}

		if v, ok := ipPoolRaw["id"]; ok {
			output.Pool = &subnets.IPamPoolPrefixAllocationPool{
				Id: pointer.To(v.(string)),
			}
		}

		outputs = append(outputs, output)
	}

	return &outputs
}

func flattenSubnetIPAddressPool(input *[]subnets.IPamPoolPrefixAllocation) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	outputs := make([]interface{}, 0)
	for _, v := range *input {
		output := map[string]interface{}{
			"number_of_ip_addresses":        pointer.From(v.NumberOfIPAddresses),
			"allocated_ip_address_prefixes": pointer.From(v.AllocatedAddressPrefixes),
		}
		if v.Pool != nil {
			output["id"] = pointer.From(v.Pool.Id)
		}
		outputs = append(outputs, output)
	}

	return outputs
}

func SubnetProvisioningStateRefreshFunc(ctx context.Context, client *subnets.SubnetsClient, id commonids.SubnetId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id, subnets.DefaultGetOperationOptions())
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id.String(), err)
		}

		if res.Model != nil && res.Model.Properties != nil && res.Model.Properties.ProvisioningState != nil {
			return res, string(*res.Model.Properties.ProvisioningState), nil
		}
		return nil, "", fmt.Errorf("unable to read provisioning state")
	}
}
