package network

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-03-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var SubnetResourceName = "azurerm_subnet"

func resourceArmSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSubnetCreate,
		Read:   resourceArmSubnetRead,
		Update: resourceArmSubnetUpdate,
		Delete: resourceArmSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_network_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				// TODO Remove this in the next major version release
				Deprecated:   "Use the `address_prefixes` property instead.",
				ExactlyOneOf: []string{"address_prefix", "address_prefixes"},
			},

			"address_prefixes": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				ExactlyOneOf: []string{"address_prefix", "address_prefixes"},
			},

			"service_endpoints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"delegation": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"service_delegation": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"Microsoft.BareMetal/AzureVMware",
											"Microsoft.BareMetal/CrayServers",
											"Microsoft.Batch/batchAccounts",
											"Microsoft.ContainerInstance/containerGroups",
											"Microsoft.Databricks/workspaces",
											"Microsoft.DBforPostgreSQL/serversv2",
											"Microsoft.HardwareSecurityModules/dedicatedHSMs",
											"Microsoft.Logic/integrationServiceEnvironments",
											"Microsoft.Netapp/volumes",
											"Microsoft.ServiceFabricMesh/networks",
											"Microsoft.Sql/managedInstances",
											"Microsoft.Sql/servers",
											"Microsoft.StreamAnalytics/streamingJobs",
											"Microsoft.Web/hostingEnvironments",
											"Microsoft.Web/serverFarms",
										}, false),
									},

									"actions": {
										Type:       schema.TypeList,
										Optional:   true,
										ConfigMode: schema.SchemaConfigModeAttr,
										Elem: &schema.Schema{
											Type: schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{
												"Microsoft.Network/networkinterfaces/*",
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

			"enforce_private_link_endpoint_network_policies": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enforce_private_link_service_network_policies": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

// TODO: refactor the create/flatten functions
func resourceArmSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for Azure ARM Subnet creation.")

	name := d.Get("name").(string)
	vnetName := d.Get("virtual_network_name").(string)
	resGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resGroup, vnetName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Subnet %q (Virtual Network %q / Resource Group %q): %s", name, vnetName, resGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_subnet", *existing.ID)
		}
	}

	locks.ByName(vnetName, VirtualNetworkResourceName)
	defer locks.UnlockByName(vnetName, VirtualNetworkResourceName)

	properties := network.SubnetPropertiesFormat{}
	if value, ok := d.GetOk("address_prefixes"); ok {
		var addressPrefixes []string
		for _, item := range value.([]interface{}) {
			addressPrefixes = append(addressPrefixes, item.(string))
		}
		properties.AddressPrefixes = &addressPrefixes
	}
	if value, ok := d.GetOk("address_prefix"); ok {
		addressPrefix := value.(string)
		properties.AddressPrefix = &addressPrefix
	}
	if properties.AddressPrefixes != nil && len(*properties.AddressPrefixes) == 1 {
		properties.AddressPrefix = &(*properties.AddressPrefixes)[0]
		properties.AddressPrefixes = nil
	}

	// To enable private endpoints you must disable the network policies for the subnet because
	// Network policies like network security groups are not supported by private endpoints.
	privateEndpointNetworkPolicies := d.Get("enforce_private_link_endpoint_network_policies").(bool)
	privateLinkServiceNetworkPolicies := d.Get("enforce_private_link_service_network_policies").(bool)
	properties.PrivateEndpointNetworkPolicies = expandSubnetPrivateLinkNetworkPolicy(privateEndpointNetworkPolicies)
	properties.PrivateLinkServiceNetworkPolicies = expandSubnetPrivateLinkNetworkPolicy(privateLinkServiceNetworkPolicies)

	serviceEndpointsRaw := d.Get("service_endpoints").([]interface{})
	properties.ServiceEndpoints = expandSubnetServiceEndpoints(serviceEndpointsRaw)

	delegationsRaw := d.Get("delegation").([]interface{})
	properties.Delegations = expandSubnetDelegation(delegationsRaw)

	subnet := network.Subnet{
		Name:                   &name,
		SubnetPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, vnetName, name, subnet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, vnetName, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Subnet %q (Virtual Network %q / Resource Group %q)", vnetName, name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSubnetRead(d, meta)
}

func resourceArmSubnetUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	networkName := id.Path["virtualNetworks"]
	name := id.Path["subnets"]

	existing, err := client.Get(ctx, resourceGroup, networkName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving existing Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, networkName, resourceGroup, err)
	}

	if existing.SubnetPropertiesFormat == nil {
		return fmt.Errorf("Error retrieving existing Subnet %q (Virtual Network %q / Resource Group %q): `properties` was nil", name, networkName, resourceGroup)
	}

	// TODO: locking on the NSG/Route Table if applicable

	props := *existing.SubnetPropertiesFormat

	if d.HasChange("address_prefix") {
		props.AddressPrefix = utils.String(d.Get("address_prefix").(string))
	}

	if d.HasChange("address_prefixes") {
		addressPrefixesRaw := d.Get("address_prefixes").([]interface{})
		props.AddressPrefixes = utils.ExpandStringSlice(addressPrefixesRaw)
		if props.AddressPrefixes != nil && len(*props.AddressPrefixes) == 1 {
			props.AddressPrefix = &(*props.AddressPrefixes)[0]
			props.AddressPrefixes = nil
		}
	}

	if d.HasChange("delegation") {
		delegationsRaw := d.Get("delegation").([]interface{})
		props.Delegations = expandSubnetDelegation(delegationsRaw)
	}

	if d.HasChange("enforce_private_link_endpoint_network_policies") {
		v := d.Get("enforce_private_link_endpoint_network_policies").(bool)
		props.PrivateEndpointNetworkPolicies = expandSubnetPrivateLinkNetworkPolicy(v)
	}

	if d.HasChange("enforce_private_link_service_network_policies") {
		v := d.Get("enforce_private_link_service_network_policies").(bool)
		props.PrivateLinkServiceNetworkPolicies = expandSubnetPrivateLinkNetworkPolicy(v)
	}

	if d.HasChange("service_endpoints") {
		serviceEndpointsRaw := d.Get("service_endpoints").([]interface{})
		props.ServiceEndpoints = expandSubnetServiceEndpoints(serviceEndpointsRaw)
	}

	subnet := network.Subnet{
		Name:                   utils.String(name),
		SubnetPropertiesFormat: &props,
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, networkName, name, subnet)
	if err != nil {
		return fmt.Errorf("Error updating Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, networkName, resourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, networkName, resourceGroup, err)
	}

	return resourceArmSubnetRead(d, meta)
}

func resourceArmSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	networkName := id.Path["virtualNetworks"]
	name := id.Path["subnets"]

	resp, err := client.Get(ctx, resourceGroup, networkName, name, "")

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, networkName, resourceGroup, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("virtual_network_name", networkName)

	if props := resp.SubnetPropertiesFormat; props != nil {
		d.Set("address_prefix", props.AddressPrefix)
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
			return fmt.Errorf("Error flattening `delegation`: %+v", err)
		}

		d.Set("enforce_private_link_endpoint_network_policies", flattenSubnetPrivateLinkNetworkPolicy(props.PrivateEndpointNetworkPolicies))
		d.Set("enforce_private_link_service_network_policies", flattenSubnetPrivateLinkNetworkPolicy(props.PrivateLinkServiceNetworkPolicies))

		serviceEndpoints := flattenSubnetServiceEndpoints(props.ServiceEndpoints)
		if err := d.Set("service_endpoints", serviceEndpoints); err != nil {
			return fmt.Errorf("Error setting `service_endpoints`: %+v", err)
		}
	}

	return nil
}

func resourceArmSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.SubnetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["subnets"]
	networkName := id.Path["virtualNetworks"]

	locks.ByName(networkName, VirtualNetworkResourceName)
	defer locks.UnlockByName(networkName, VirtualNetworkResourceName)

	locks.ByName(name, SubnetResourceName)
	defer locks.UnlockByName(name, SubnetResourceName)

	future, err := client.Delete(ctx, resourceGroup, networkName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, networkName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of Subnet %q (Virtual Network %q / Resource Group %q): %+v", name, networkName, resourceGroup, err)
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

func flattenSubnetServiceEndpoints(serviceEndpoints *[]network.ServiceEndpointPropertiesFormat) []string {
	endpoints := make([]string, 0)

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

	for _, dele := range *delegations {
		retDele := make(map[string]interface{})
		if v := dele.Name; v != nil {
			retDele["name"] = *v
		}

		svcDeles := make([]interface{}, 0)
		svcDele := make(map[string]interface{})
		if props := dele.ServiceDelegationPropertiesFormat; props != nil {
			if v := props.ServiceName; v != nil {
				svcDele["name"] = *v
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

// TODO: confirm this logic below

func expandSubnetPrivateLinkNetworkPolicy(enabled bool) *string {
	// This is strange logic, but to get the schema to make sense for the end user
	// I exposed it with the same name that the Azure CLI does to be consistent
	// between the tool sets, which means true == Disabled.
	if enabled {
		return utils.String("Disabled")
	}

	return utils.String("Enabled")
}

func flattenSubnetPrivateLinkNetworkPolicy(input *string) bool {
	// This is strange logic, but to get the schema to make sense for the end user
	// I exposed it with the same name that the Azure CLI does to be consistent
	// between the tool sets, which means true == Disabled.
	if input == nil {
		return false
	}

	return strings.EqualFold(*input, "Disabled")
}
