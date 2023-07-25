// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/network/2022-07-01/network"
)

var SubnetResourceName = "azurerm_subnet"

var subnetDelegationServiceNames = []string{
	"GitHub.Network/networkSettings",
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
	"Microsoft.DocumentDB/cassandraClusters",
	"Microsoft.Fidalgo/networkSettings",
	"Microsoft.HardwareSecurityModules/dedicatedHSMs",
	"Microsoft.Kusto/clusters",
	"Microsoft.LabServices/labplans",
	"Microsoft.Logic/integrationServiceEnvironments",
	"Microsoft.MachineLearningServices/workspaces",
	"Microsoft.Netapp/volumes",
	"Microsoft.Network/dnsResolvers",
	"Microsoft.Network/managedResolvers",
	"Microsoft.Network/fpgaNetworkInterfaces",
	"Microsoft.Network/networkWatchers.",
	"Microsoft.Network/virtualNetworkGateways",
	"Microsoft.Orbital/orbitalGateways",
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
	"PaloAltoNetworks.Cloudngfw/firewalls",
	"Qumulo.Storage/fileSystems",
}

func resourceSubnet() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceSubnetCreate,
		Read:   resourceSubnetRead,
		Update: resourceSubnetUpdate,
		Delete: resourceSubnetDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := commonids.ParseSubnetID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
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
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
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
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validate.SubnetServiceEndpointStoragePolicyID,
				},
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
										Type:       pluginsdk.TypeList,
										Optional:   true,
										ConfigMode: pluginsdk.SchemaConfigModeAttr,
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

			"private_endpoint_network_policies_enabled": {
				Type: pluginsdk.TypeBool,
				Computed: func() bool {
					return !features.FourPointOh()
				}(),
				Optional: true,
				Default: func() interface{} {
					if !features.FourPointOh() {
						return nil
					}
					return !features.FourPointOh()
				}(),
				ConflictsWith: func() []string {
					if !features.FourPointOh() {
						return []string{"enforce_private_link_endpoint_network_policies"}
					}
					return []string{}
				}(),
			},

			"private_link_service_network_policies_enabled": {
				Type: pluginsdk.TypeBool,
				Computed: func() bool {
					return !features.FourPointOh()
				}(),
				Optional: true,
				Default: func() interface{} {
					if !features.FourPointOh() {
						return nil
					}
					return features.FourPointOh()
				}(),
				ConflictsWith: func() []string {
					if !features.FourPointOh() {
						return []string{"enforce_private_link_service_network_policies"}
					}
					return []string{}
				}(),
			},
		},
	}

	if !features.FourPointOhBeta() {
		resource.Schema["enforce_private_link_endpoint_network_policies"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Computed:      true,
			Optional:      true,
			Deprecated:    "`enforce_private_link_endpoint_network_policies` will be removed in favour of the property `private_endpoint_network_policies_enabled` in version 4.0 of the AzureRM Provider",
			ConflictsWith: []string{"private_endpoint_network_policies_enabled"},
		}

		resource.Schema["enforce_private_link_service_network_policies"] = &pluginsdk.Schema{
			Type:          pluginsdk.TypeBool,
			Computed:      true,
			Optional:      true,
			Deprecated:    "`enforce_private_link_service_network_policies` will be removed in favour of the property `private_link_service_network_policies_enabled` in version 4.0 of the AzureRM Provider",
			ConflictsWith: []string{"private_link_service_network_policies_enabled"},
		}
	}

	return resource
}

// TODO: refactor the create/flatten functions
func resourceSubnetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	vnetClient := meta.(*clients.Client).Network.VnetClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Subnet creation.")

	id := commonids.NewSubnetID(subscriptionId, d.Get("resource_group_name").(string), d.Get("virtual_network_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName, "")
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_subnet", id.ID())
	}

	locks.ByName(id.VirtualNetworkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(id.VirtualNetworkName, VirtualNetworkResourceName)

	properties := network.SubnetPropertiesFormat{}
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

	// To enable private endpoints you must disable the network policies for the subnet because
	// Network policies like network security groups are not supported by private endpoints.
	var privateEndpointNetworkPolicies network.VirtualNetworkPrivateEndpointNetworkPolicies
	var privateLinkServiceNetworkPolicies network.VirtualNetworkPrivateLinkServiceNetworkPolicies

	if features.FourPointOhBeta() {
		privateEndpointNetworkPoliciesRaw := d.Get("private_endpoint_network_policies_enabled").(bool)
		privateLinkServiceNetworkPoliciesRaw := d.Get("private_link_service_network_policies_enabled").(bool)

		privateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPolicies(expandSubnetNetworkPolicy(privateEndpointNetworkPoliciesRaw))
		privateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandSubnetNetworkPolicy(privateLinkServiceNetworkPoliciesRaw))
	} else {
		var enforceOk bool
		var enforceServiceOk bool
		var enableOk bool
		var enableServiceOk bool
		var enforcePrivateEndpointNetworkPoliciesRaw bool
		var enforcePrivateLinkServiceNetworkPoliciesRaw bool
		var privateEndpointNetworkPoliciesRaw bool
		var privateLinkServiceNetworkPoliciesRaw bool

		// Set the legacy default value since they are now computed optional
		privateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled
		privateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPoliciesEnabled

		// This is the only way I was able to figure out if the fields are actually in the config or not,
		// which is needed here because these are all now optional computed fields...
		if !pluginsdk.IsExplicitlyNullInConfig(d, "enforce_private_link_endpoint_network_policies") {
			enforceOk = true
			enforcePrivateEndpointNetworkPoliciesRaw = d.Get("enforce_private_link_endpoint_network_policies").(bool)
		}

		if !pluginsdk.IsExplicitlyNullInConfig(d, "enforce_private_link_service_network_policies") {
			enforceServiceOk = true
			enforcePrivateLinkServiceNetworkPoliciesRaw = d.Get("enforce_private_link_service_network_policies").(bool)
		}

		if !pluginsdk.IsExplicitlyNullInConfig(d, "private_endpoint_network_policies_enabled") {
			enableOk = true
			privateEndpointNetworkPoliciesRaw = d.Get("private_endpoint_network_policies_enabled").(bool)
		}

		if !pluginsdk.IsExplicitlyNullInConfig(d, "private_link_service_network_policies_enabled") {
			enableServiceOk = true
			privateLinkServiceNetworkPoliciesRaw = d.Get("private_link_service_network_policies_enabled").(bool)
		}

		// Only one of these values can be set since they conflict with each other
		// if neither of them are set use the default values
		if enforceOk || enableOk {
			if enforceOk {
				privateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPolicies(expandEnforceSubnetNetworkPolicy(enforcePrivateEndpointNetworkPoliciesRaw))
			} else if enableOk {
				privateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPolicies(expandSubnetNetworkPolicy(privateEndpointNetworkPoliciesRaw))
			}
		}

		if enforceServiceOk || enableServiceOk {
			if enforceServiceOk {
				privateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandEnforceSubnetNetworkPolicy(enforcePrivateLinkServiceNetworkPoliciesRaw))
			} else if enableServiceOk {
				privateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandSubnetNetworkPolicy(privateLinkServiceNetworkPoliciesRaw))
			}
		}
	}

	properties.PrivateEndpointNetworkPolicies = privateEndpointNetworkPolicies
	properties.PrivateLinkServiceNetworkPolicies = privateLinkServiceNetworkPolicies

	serviceEndpointPoliciesRaw := d.Get("service_endpoint_policy_ids").(*pluginsdk.Set).List()
	properties.ServiceEndpointPolicies = expandSubnetServiceEndpointPolicies(serviceEndpointPoliciesRaw)

	serviceEndpointsRaw := d.Get("service_endpoints").(*pluginsdk.Set).List()
	properties.ServiceEndpoints = expandSubnetServiceEndpoints(serviceEndpointsRaw)

	delegationsRaw := d.Get("delegation").([]interface{})
	properties.Delegations = expandSubnetDelegation(delegationsRaw)

	subnet := network.Subnet{
		Name:                   utils.String(id.SubnetName),
		SubnetPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName, subnet)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    VirtualNetworkProvisioningStateRefreshFunc(ctx, vnetClient, vnetId),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = vnetStateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of virtual network for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSubnetRead(d, meta)
}

func resourceSubnetUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	vnetClient := meta.(*clients.Client).Network.VnetClient
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

	existing, err := client.Get(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName, "")
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	if existing.SubnetPropertiesFormat == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	// TODO: locking on the NSG/Route Table if applicable

	props := *existing.SubnetPropertiesFormat

	if d.HasChange("address_prefixes") {
		addressPrefixesRaw := d.Get("address_prefixes").([]interface{})
		switch len(addressPrefixesRaw) {
		case 0:
			// Will never happen as the "MinItem: 1" constraint is set on "address_prefixes"
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

	if d.HasChange("delegation") {
		delegationsRaw := d.Get("delegation").([]interface{})
		props.Delegations = expandSubnetDelegation(delegationsRaw)
	}

	if features.FourPointOhBeta() {
		if d.HasChange("private_endpoint_network_policies_enabled") {
			v := d.Get("private_endpoint_network_policies_enabled").(bool)
			props.PrivateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPolicies(expandSubnetNetworkPolicy(v))
		}

		if d.HasChange("private_link_service_network_policies_enabled") {
			v := d.Get("private_link_service_network_policies_enabled").(bool)
			props.PrivateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandSubnetNetworkPolicy(v))
		}
	} else {
		// This is the best case we can do in this state since they are computed optional fields now
		// If you remove the fields from the config they will just persist as they are, if you change
		// one it will update it to the value that was changed and in the read the other value will be
		// updated as well to reflect the new value so it is safe to toggle between which field you want
		// to use to define this behavior...
		var privateEndpointNetworkPolicies network.VirtualNetworkPrivateEndpointNetworkPolicies
		var privateLinkServiceNetworkPolicies network.VirtualNetworkPrivateLinkServiceNetworkPolicies

		if d.HasChange("enforce_private_link_endpoint_network_policies") || d.HasChange("private_endpoint_network_policies_enabled") {
			enforcePrivateEndpointNetworkPoliciesRaw := d.Get("enforce_private_link_endpoint_network_policies").(bool)
			privateEndpointNetworkPoliciesRaw := d.Get("private_endpoint_network_policies_enabled").(bool)

			if d.HasChange("enforce_private_link_endpoint_network_policies") {
				privateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPolicies(expandEnforceSubnetNetworkPolicy(enforcePrivateEndpointNetworkPoliciesRaw))
			} else if d.HasChange("private_endpoint_network_policies_enabled") {
				privateEndpointNetworkPolicies = network.VirtualNetworkPrivateEndpointNetworkPolicies(expandSubnetNetworkPolicy(privateEndpointNetworkPoliciesRaw))
			}

			props.PrivateEndpointNetworkPolicies = privateEndpointNetworkPolicies
		}

		if d.HasChange("enforce_private_link_service_network_policies") || d.HasChange("private_link_service_network_policies_enabled") {
			enforcePrivateLinkServiceNetworkPoliciesRaw := d.Get("enforce_private_link_service_network_policies").(bool)
			privateLinkServiceNetworkPoliciesRaw := d.Get("private_link_service_network_policies_enabled").(bool)

			if d.HasChange("enforce_private_link_service_network_policies") {
				privateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandEnforceSubnetNetworkPolicy(enforcePrivateLinkServiceNetworkPoliciesRaw))
			} else if d.HasChange("private_link_service_network_policies_enabled") {
				privateLinkServiceNetworkPolicies = network.VirtualNetworkPrivateLinkServiceNetworkPolicies(expandSubnetNetworkPolicy(privateLinkServiceNetworkPoliciesRaw))
			}

			props.PrivateLinkServiceNetworkPolicies = privateLinkServiceNetworkPolicies
		}
	}

	if d.HasChange("service_endpoints") {
		serviceEndpointsRaw := d.Get("service_endpoints").(*pluginsdk.Set).List()
		props.ServiceEndpoints = expandSubnetServiceEndpoints(serviceEndpointsRaw)
	}

	if d.HasChange("service_endpoint_policy_ids") {
		serviceEndpointPoliciesRaw := d.Get("service_endpoint_policy_ids").(*pluginsdk.Set).List()
		props.ServiceEndpointPolicies = expandSubnetServiceEndpointPolicies(serviceEndpointPoliciesRaw)
	}

	subnet := network.Subnet{
		Name:                   utils.String(id.SubnetName),
		SubnetPropertiesFormat: &props,
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName, subnet)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	timeout, _ := ctx.Deadline()

	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
		Refresh:    SubnetProvisioningStateRefreshFunc(ctx, client, *id),
		MinTimeout: 1 * time.Minute,
		Timeout:    time.Until(timeout),
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of %s: %+v", id, err)
	}

	vnetId := commonids.NewVirtualNetworkID(id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName)
	vnetStateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{string(network.ProvisioningStateUpdating)},
		Target:     []string{string(network.ProvisioningStateSucceeded)},
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
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := commonids.ParseSubnetID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.SubnetName)
	d.Set("virtual_network_name", id.VirtualNetworkName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if props := resp.SubnetPropertiesFormat; props != nil {
		if props.AddressPrefixes == nil {
			if props.AddressPrefix != nil && len(*props.AddressPrefix) > 0 {
				d.Set("address_prefixes", []string{*props.AddressPrefix})
			} else {
				d.Set("address_prefixes", []string{})
			}
		} else {
			d.Set("address_prefixes", props.AddressPrefixes)
		}

		delegation := flattenSubnetDelegation(props.Delegations)
		if err := d.Set("delegation", delegation); err != nil {
			return fmt.Errorf("flattening `delegation`: %+v", err)
		}

		if !features.FourPointOhBeta() {
			d.Set("enforce_private_link_endpoint_network_policies", flattenEnforceSubnetNetworkPolicy(string(props.PrivateEndpointNetworkPolicies)))
			d.Set("enforce_private_link_service_network_policies", flattenEnforceSubnetNetworkPolicy(string(props.PrivateLinkServiceNetworkPolicies)))
		}

		d.Set("private_endpoint_network_policies_enabled", flattenSubnetNetworkPolicy(string(props.PrivateEndpointNetworkPolicies)))
		d.Set("private_link_service_network_policies_enabled", flattenSubnetNetworkPolicy(string(props.PrivateLinkServiceNetworkPolicies)))

		serviceEndpoints := flattenSubnetServiceEndpoints(props.ServiceEndpoints)
		if err := d.Set("service_endpoints", serviceEndpoints); err != nil {
			return fmt.Errorf("setting `service_endpoints`: %+v", err)
		}

		serviceEndpointPolicies := flattenSubnetServiceEndpointPolicies(props.ServiceEndpointPolicies)
		if err := d.Set("service_endpoint_policy_ids", serviceEndpointPolicies); err != nil {
			return fmt.Errorf("setting `service_endpoint_policy_ids`: %+v", err)
		}
	}

	return nil
}

func resourceSubnetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
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

	future, err := client.Delete(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	return nil
}

func expandSubnetServiceEndpoints(input []interface{}) *[]network.ServiceEndpointPropertiesFormat {
	endpoints := make([]network.ServiceEndpointPropertiesFormat, 0)

	for _, svcEndpointRaw := range input {
		if svc, ok := svcEndpointRaw.(string); ok {
			endpoint := network.ServiceEndpointPropertiesFormat{
				Service: &svc,
			}
			endpoints = append(endpoints, endpoint)
		}
	}

	return &endpoints
}

func flattenSubnetServiceEndpoints(serviceEndpoints *[]network.ServiceEndpointPropertiesFormat) []interface{} {
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

func expandSubnetDelegation(input []interface{}) *[]network.Delegation {
	retDelegations := make([]network.Delegation, 0)

	for _, deleValue := range input {
		deleData := deleValue.(map[string]interface{})
		deleName := deleData["name"].(string)
		srvDelegations := deleData["service_delegation"].([]interface{})
		srvDelegation := srvDelegations[0].(map[string]interface{})
		srvName := srvDelegation["name"].(string)
		srvActions := srvDelegation["actions"].([]interface{})

		retSrvActions := make([]string, 0)
		for _, srvAction := range srvActions {
			srvActionData := srvAction.(string)
			retSrvActions = append(retSrvActions, srvActionData)
		}

		retDelegation := network.Delegation{
			Name: &deleName,
			ServiceDelegationPropertiesFormat: &network.ServiceDelegationPropertiesFormat{
				ServiceName: &srvName,
				Actions:     &retSrvActions,
			},
		}

		retDelegations = append(retDelegations, retDelegation)
	}

	return &retDelegations
}

func flattenSubnetDelegation(delegations *[]network.Delegation) []interface{} {
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
		if props := dele.ServiceDelegationPropertiesFormat; props != nil {
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

// TODO 4.0: Remove expandEnforceSubnetPrivateLinkNetworkPolicy function
func expandEnforceSubnetNetworkPolicy(enabled bool) string {
	// This is strange logic, but to get the schema to make sense for the end user
	// I exposed it with the same name that the Azure CLI does to be consistent
	// between the tool sets, which means true == Disabled.
	if enabled {
		return string(network.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled)
	}

	return string(network.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled)
}

func expandSubnetNetworkPolicy(enabled bool) string {
	if enabled {
		return string(network.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled)
	}

	return string(network.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled)
}

// TODO 4.0: Remove flattenEnforceSubnetPrivateLinkNetworkPolicy function
func flattenEnforceSubnetNetworkPolicy(input string) bool {
	// This is strange logic, but to get the schema to make sense for the end user
	// I exposed it with the same name that the Azure CLI does to be consistent
	// between the tool sets, which means true == Disabled.
	return strings.EqualFold(input, string(network.VirtualNetworkPrivateEndpointNetworkPoliciesDisabled))
}

func flattenSubnetNetworkPolicy(input string) bool {
	return strings.EqualFold(input, string(network.VirtualNetworkPrivateEndpointNetworkPoliciesEnabled))
}

func expandSubnetServiceEndpointPolicies(input []interface{}) *[]network.ServiceEndpointPolicy {
	output := make([]network.ServiceEndpointPolicy, 0)
	for _, policy := range input {
		policy := policy.(string)
		output = append(output, network.ServiceEndpointPolicy{ID: &policy})
	}
	return &output
}

func flattenSubnetServiceEndpointPolicies(input *[]network.ServiceEndpointPolicy) []interface{} {
	if input == nil {
		return nil
	}

	var output []interface{}
	for _, policy := range *input {
		id := ""
		if policy.ID != nil {
			id = *policy.ID
		}
		output = append(output, id)
	}
	return output
}

func SubnetProvisioningStateRefreshFunc(ctx context.Context, client *network.SubnetsClient, id commonids.SubnetId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName, "")
		if err != nil {
			return nil, "", fmt.Errorf("polling for %s: %+v", id.String(), err)
		}

		return res, string(res.ProvisioningState), nil
	}
}
