package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-08-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
)

func resourceArmNetworkProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmNetworkProfileCreateUpdate,
		Read:   resourceArmNetworkProfileRead,
		Update: resourceArmNetworkProfileCreateUpdate,
		Delete: resourceArmNetworkProfileDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"container_network_interface_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"ip_configuration": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validation.NoZeroValues,
									},
									"subnet_id": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: azure.ValidateResourceID,
									},
								},
							},
						},
					},
				},
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmNetworkProfileCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netProfileClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM Network Profile creation")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	cniConfigs, err := expandArmContainerNetworkInterfaceConfigurations(d)
	if err != nil {
		return fmt.Errorf("Error building list of Azure Container Network Interface Configurations: %+v", err)
	}

	parameters := network.Profile{
		Location: &location,
		Tags:     expandTags(tags),
		ProfilePropertiesFormat: &network.ProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: cniConfigs,
		},
	}

	_, err = client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating Azure Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	profile, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Azure Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if profile.ID == nil {
		return fmt.Errorf("Cannot read Azure Network Profile %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*profile.ID)

	return resourceArmNetworkProfileRead(d, meta)
}

func resourceArmNetworkProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netProfileClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkProfiles"]

	profile, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error making read request on Azure Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", profile.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := profile.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := profile.ProfilePropertiesFormat; props != nil {
		cniConfigs := flattenArmContainerNetworkInterfaceConfigurations(props.ContainerNetworkInterfaceConfigurations)
		if err := d.Set("container_network_interface_configuration", cniConfigs); err != nil {
			return fmt.Errorf("Error setting `container_network_interface_configuration`: %+v", err)
		}
	}

	flattenAndSetTags(d, profile.Tags)

	return nil
}

func resourceArmNetworkProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).netProfileClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["azureFirewalls"]

	_, err = client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting Azure Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandArmContainerNetworkInterfaceConfigurations(d *schema.ResourceData) (*[]network.ContainerNetworkInterfaceConfiguration, error) {
	cniConfigs := d.Get("container_network_interface_configuration").([]interface{})
	retCNIConfigs := make([]network.ContainerNetworkInterfaceConfiguration, 0)

	for _, cniConfig := range cniConfigs {
		nciData := cniConfig.(map[string]interface{})
		nciName := nciData["name"].(string)
		ipConfigs := nciData["ip_configuration"].([]interface{})

		retIPConfigs := make([]network.IPConfigurationProfile, 0)
		for _, ipConfig := range ipConfigs {
			ipData := ipConfig.(map[string]interface{})
			ipName := ipData["name"].(string)
			subNetID := ipData["subnet_id"].(string)

			retIPConfig := network.IPConfigurationProfile{
				Name: &ipName,
				IPConfigurationProfilePropertiesFormat: &network.IPConfigurationProfilePropertiesFormat{
					Subnet: &network.Subnet{
						ID: &subNetID,
					},
				},
			}

			retIPConfigs = append(retIPConfigs, retIPConfig)
		}

		retCNIConfig := network.ContainerNetworkInterfaceConfiguration{
			Name: &nciName,
			ContainerNetworkInterfaceConfigurationPropertiesFormat: &network.ContainerNetworkInterfaceConfigurationPropertiesFormat{
				IPConfigurations: &retIPConfigs,
			},
		}

		retCNIConfigs = append(retCNIConfigs, retCNIConfig)
	}

	return &retCNIConfigs, nil
}

func flattenArmContainerNetworkInterfaceConfigurations(input *[]network.ContainerNetworkInterfaceConfiguration) []interface{} {
	retCNIConfigs := make([]interface{}, 0)
	if input == nil {
		return retCNIConfigs
	}

	for _, cniConfig := range *input {
		retCNIConfig := make(map[string]interface{})
		cniProps := cniConfig.ContainerNetworkInterfaceConfigurationPropertiesFormat
		if cniProps == nil || cniProps.IPConfigurations == nil {
			continue
		}

		if cniConfig.Name != nil {
			retCNIConfig["name"] = *cniConfig.Name
		}

		retIPConfigs := make([]interface{}, 0)
		for _, ipConfig := range *cniProps.IPConfigurations {
			retIPConfig := make(map[string]interface{})
			ipProps := ipConfig.IPConfigurationProfilePropertiesFormat
			if ipProps == nil || ipProps.Subnet == nil {
				continue
			}

			if ipConfig.Name != nil {
				retIPConfig["name"] = *ipConfig.Name
			}

			retIPConfig["subnet_id"] = *ipProps.Subnet.ID

			retIPConfigs = append(retIPConfigs, retIPConfig)
		}

		retCNIConfig["ip_configuration"] = retIPConfigs

		retCNIConfigs = append(retCNIConfigs, retCNIConfig)
	}

	return retCNIConfigs
}
