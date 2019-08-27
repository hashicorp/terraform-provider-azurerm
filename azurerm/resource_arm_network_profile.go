package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/locks"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

const azureNetworkProfileResourceName = "azurerm_network_profile"

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"container_network_interface": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"ip_configuration": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validate.NoEmptyStrings,
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

			"container_network_interface_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmNetworkProfileCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ProfileClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for Network Profile creation")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Network Profile %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_network_profile", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	cniConfigs, err := expandNetworkProfileContainerNetworkInterface(d)
	if err != nil {
		return fmt.Errorf("Error building list of Container Network Interface Configurations: %+v", err)
	}

	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(d)
	if err != nil {
		return fmt.Errorf("Error extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByName(name, azureNetworkProfileResourceName)
	defer locks.UnlockByName(name, azureNetworkProfileResourceName)

	locks.MultipleByName(subnetsToLock, subnetResourceName)
	defer locks.UnlockMultipleByName(subnetsToLock, subnetResourceName)

	locks.MultipleByName(vnetsToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, virtualNetworkResourceName)

	parameters := network.Profile{
		Location: &location,
		Tags:     tags.Expand(t),
		ProfilePropertiesFormat: &network.ProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: cniConfigs,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters); err != nil {
		return fmt.Errorf("Error creating/updating Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	profile, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if profile.ID == nil {
		return fmt.Errorf("Cannot read Network Profile %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*profile.ID)

	return resourceArmNetworkProfileRead(d, meta)
}

func resourceArmNetworkProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ProfileClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkProfiles"]

	profile, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(profile.Response) {
			log.Printf("[DEBUG] Network Profile %q was not found in Resource Group %q - removing from state!", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", profile.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := profile.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := profile.ProfilePropertiesFormat; props != nil {
		cniConfigs := flattenNetworkProfileContainerNetworkInterface(props.ContainerNetworkInterfaceConfigurations)
		if err := d.Set("container_network_interface", cniConfigs); err != nil {
			return fmt.Errorf("Error setting `container_network_interface`: %+v", err)
		}

		cniIDs := flattenNetworkProfileContainerNetworkInterfaceIDs(props.ContainerNetworkInterfaces)
		if err := d.Set("container_network_interface_ids", cniIDs); err != nil {
			return fmt.Errorf("Error setting `container_network_interface_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, profile.Tags)
}

func resourceArmNetworkProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.ProfileClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["networkProfiles"]

	read, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] Network Profile %q was not found in Resource Group %q - assuming removed!", name, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(d)
	if err != nil {
		return fmt.Errorf("Error extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByName(name, azureNetworkProfileResourceName)
	defer locks.UnlockByName(name, azureNetworkProfileResourceName)

	locks.MultipleByName(subnetsToLock, subnetResourceName)
	defer locks.UnlockMultipleByName(subnetsToLock, subnetResourceName)

	locks.MultipleByName(vnetsToLock, virtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, virtualNetworkResourceName)

	if _, err = client.Delete(ctx, resourceGroup, name); err != nil {
		return fmt.Errorf("Error deleting Network Profile %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return err
}

func expandNetworkProfileContainerNetworkInterface(d *schema.ResourceData) (*[]network.ContainerNetworkInterfaceConfiguration, error) {
	cniConfigs := d.Get("container_network_interface").([]interface{})
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

func expandNetworkProfileVirtualNetworkSubnetNames(d *schema.ResourceData) (*[]string, *[]string, error) {
	cniConfigs := d.Get("container_network_interface").([]interface{})
	subnetNames := make([]string, 0)
	vnetNames := make([]string, 0)

	for _, cniConfig := range cniConfigs {
		nciData := cniConfig.(map[string]interface{})
		ipConfigs := nciData["ip_configuration"].([]interface{})

		for _, ipConfig := range ipConfigs {
			ipData := ipConfig.(map[string]interface{})
			subnetID := ipData["subnet_id"].(string)

			subnetResourceID, err := azure.ParseAzureResourceID(subnetID)
			if err != nil {
				return nil, nil, err
			}

			subnetName := subnetResourceID.Path["subnets"]
			vnetName := subnetResourceID.Path["virtualNetworks"]

			if !sliceContainsValue(subnetNames, subnetName) {
				subnetNames = append(subnetNames, subnetName)
			}

			if !sliceContainsValue(vnetNames, vnetName) {
				vnetNames = append(vnetNames, vnetName)
			}
		}
	}

	return &subnetNames, &vnetNames, nil
}

func flattenNetworkProfileContainerNetworkInterface(input *[]network.ContainerNetworkInterfaceConfiguration) []interface{} {
	retCNIConfigs := make([]interface{}, 0)
	if input == nil {
		return retCNIConfigs
	}

	// if-continue is used to simplify the deeply nested if-else statement.
	for _, cniConfig := range *input {
		retCNIConfig := make(map[string]interface{})

		if cniConfig.Name != nil {
			retCNIConfig["name"] = *cniConfig.Name
		}

		retIPConfigs := make([]interface{}, 0)
		if cniProps := cniConfig.ContainerNetworkInterfaceConfigurationPropertiesFormat; cniProps != nil && cniProps.IPConfigurations != nil {
			for _, ipConfig := range *cniProps.IPConfigurations {
				retIPConfig := make(map[string]interface{})

				if ipConfig.Name != nil {
					retIPConfig["name"] = *ipConfig.Name
				}

				if ipProps := ipConfig.IPConfigurationProfilePropertiesFormat; ipProps != nil && ipProps.Subnet != nil && ipProps.Subnet.ID != nil {
					retIPConfig["subnet_id"] = *ipProps.Subnet.ID
				}

				retIPConfigs = append(retIPConfigs, retIPConfig)
			}
		}
		retCNIConfig["ip_configuration"] = retIPConfigs

		retCNIConfigs = append(retCNIConfigs, retCNIConfig)
	}

	return retCNIConfigs
}

func flattenNetworkProfileContainerNetworkInterfaceIDs(input *[]network.ContainerNetworkInterface) []string {
	retCNIs := make([]string, 0)
	if input == nil {
		return retCNIs
	}

	for _, retCNI := range *input {
		if retCNI.ID != nil {
			retCNIs = append(retCNIs, *retCNI.ID)
		}
	}

	return retCNIs
}
