package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-08-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

var subnetResourceName = "azurerm_subnet"

func resourceArmSubnet() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmSubnetCreate,
		Read:   resourceArmSubnetRead,
		Update: resourceArmSubnetCreate,
		Delete: resourceArmSubnetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": resourceGroupNameSchema(),

			"virtual_network_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"address_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},

			"network_security_group_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use the `azurerm_subnet_network_security_group_association` resource instead.",
			},

			"route_table_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "Use the `azurerm_subnet_route_table_association` resource instead.",
			},

			"ip_configurations": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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
									"service_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"actions": {
										Type:     schema.TypeList,
										Optional: true,
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

func resourceArmSubnetCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Azure ARM Subnet creation.")

	name := d.Get("name").(string)
	vnetName := d.Get("virtual_network_name").(string)
	resGroup := d.Get("resource_group_name").(string)
	addressPrefix := d.Get("address_prefix").(string)

	azureRMLockByName(vnetName, virtualNetworkResourceName)
	defer azureRMUnlockByName(vnetName, virtualNetworkResourceName)

	properties := network.SubnetPropertiesFormat{
		AddressPrefix: &addressPrefix,
	}

	if v, ok := d.GetOk("network_security_group_id"); ok {
		nsgId := v.(string)
		properties.NetworkSecurityGroup = &network.SecurityGroup{
			ID: &nsgId,
		}

		networkSecurityGroupName, err := parseNetworkSecurityGroupName(nsgId)
		if err != nil {
			return err
		}

		azureRMLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureRMUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	} else {
		properties.NetworkSecurityGroup = nil
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		rtId := v.(string)
		properties.RouteTable = &network.RouteTable{
			ID: &rtId,
		}

		routeTableName, err := parseRouteTableName(rtId)
		if err != nil {
			return err
		}

		azureRMLockByName(routeTableName, routeTableResourceName)
		defer azureRMUnlockByName(routeTableName, routeTableResourceName)
	} else {
		properties.RouteTable = nil
	}

	serviceEndpoints := expandSubnetServiceEndpoints(d)
	properties.ServiceEndpoints = &serviceEndpoints

	delegations := expandSubnetDelegations(d)
	properties.Delegations = &delegations

	subnet := network.Subnet{
		Name:                   &name,
		SubnetPropertiesFormat: &properties,
	}

	future, err := client.CreateOrUpdate(ctx, resGroup, vnetName, name, subnet)
	if err != nil {
		return fmt.Errorf("Error Creating/Updating Subnet %q (VN %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion of Subnet %q (VN %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	read, err := client.Get(ctx, resGroup, vnetName, name, "")
	if err != nil {
		return err
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read ID of Subnet %q (VN %q / Resource Group %q)", vnetName, name, resGroup)
	}

	d.SetId(*read.ID)

	return resourceArmSubnetRead(d, meta)
}

func resourceArmSubnetRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	vnetName := id.Path["virtualNetworks"]
	name := id.Path["subnets"]

	resp, err := client.Get(ctx, resGroup, vnetName, name, "")

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error making Read request on Azure Subnet %q: %+v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("virtual_network_name", vnetName)

	if props := resp.SubnetPropertiesFormat; props != nil {
		d.Set("address_prefix", props.AddressPrefix)

		var securityGroupId *string
		if props.NetworkSecurityGroup != nil {
			securityGroupId = props.NetworkSecurityGroup.ID
		}
		d.Set("network_security_group_id", securityGroupId)

		var routeTableId string
		if props.RouteTable != nil && props.RouteTable.ID != nil {
			routeTableId = *props.RouteTable.ID
		}
		d.Set("route_table_id", routeTableId)

		ips := flattenSubnetIPConfigurations(props.IPConfigurations)
		if err := d.Set("ip_configurations", ips); err != nil {
			return err
		}

		serviceEndpoints := flattenSubnetServiceEndpoints(props.ServiceEndpoints)
		if err := d.Set("service_endpoints", serviceEndpoints); err != nil {
			return err
		}
	}

	return nil
}

func resourceArmSubnetDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).subnetClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resGroup := id.ResourceGroup
	name := id.Path["subnets"]
	vnetName := id.Path["virtualNetworks"]

	if v, ok := d.GetOk("network_security_group_id"); ok {
		networkSecurityGroupId := v.(string)
		networkSecurityGroupName, err2 := parseNetworkSecurityGroupName(networkSecurityGroupId)
		if err2 != nil {
			return err2
		}

		azureRMLockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
		defer azureRMUnlockByName(networkSecurityGroupName, networkSecurityGroupResourceName)
	}

	if v, ok := d.GetOk("route_table_id"); ok {
		rtId := v.(string)
		routeTableName, err2 := parseRouteTableName(rtId)
		if err2 != nil {
			return err2
		}

		azureRMLockByName(routeTableName, routeTableResourceName)
		defer azureRMUnlockByName(routeTableName, routeTableResourceName)
	}

	azureRMLockByName(vnetName, virtualNetworkResourceName)
	defer azureRMUnlockByName(vnetName, virtualNetworkResourceName)

	azureRMLockByName(name, subnetResourceName)
	defer azureRMUnlockByName(name, subnetResourceName)

	future, err := client.Delete(ctx, resGroup, vnetName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Subnet %q (VN %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for completion for Subnet %q (VN %q / Resource Group %q): %+v", name, vnetName, resGroup, err)
	}

	return nil
}

func expandSubnetServiceEndpoints(d *schema.ResourceData) []network.ServiceEndpointPropertiesFormat {
	serviceEndpoints := d.Get("service_endpoints").([]interface{})
	enpoints := make([]network.ServiceEndpointPropertiesFormat, 0)

	for _, serviceEndpointsRaw := range serviceEndpoints {
		data := serviceEndpointsRaw.(string)

		endpoint := network.ServiceEndpointPropertiesFormat{
			Service: &data,
		}

		enpoints = append(enpoints, endpoint)
	}

	return enpoints
}

func flattenSubnetServiceEndpoints(serviceEndpoints *[]network.ServiceEndpointPropertiesFormat) []string {
	endpoints := make([]string, 0)

	if serviceEndpoints != nil {
		for _, endpoint := range *serviceEndpoints {
			endpoints = append(endpoints, *endpoint.Service)
		}
	}

	return endpoints
}

func flattenSubnetIPConfigurations(ipConfigurations *[]network.IPConfiguration) []string {
	ips := make([]string, 0)

	if ipConfigurations != nil {
		for _, ip := range *ipConfigurations {
			ips = append(ips, *ip.ID)
		}
	}

	return ips
}

func expandSubnetDelegations(d *schema.ResourceData) []network.Delegation {
	delegations := d.Get("delegation").([]interface{})
	retDelegations := make([]network.Delegation, 0)

	for _, deleValue := range delegations {
		deleData := deleValue.(map[string]interface{})
		deleName := deleData["name"].(string)
		srvDelegations := deleData["service_delegation"].([]interface{})
		srvDelegation := srvDelegations[0].(map[string]interface{})
		srvName := srvDelegation["service_name"].(string)
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

	return retDelegations
}
