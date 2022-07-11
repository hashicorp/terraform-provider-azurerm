package devtestlabs

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2018-09-15/dtl"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceArmDevTestVirtualNetwork() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmDevTestVirtualNetworkCreate,
		Read:   resourceArmDevTestVirtualNetworkRead,
		Update: resourceArmDevTestVirtualNetworkUpdate,
		Delete: resourceArmDevTestVirtualNetworkDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DevTestVirtualNetworkID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.DevTestVirtualNetworkUpgradeV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateDevTestVirtualNetworkName(),
			},

			"lab_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"subnet": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				// whilst the API accepts multiple, in practice only one is usable
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"use_in_virtual_machine_creation": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(dtl.Allow),
							ValidateFunc: validate.DevTestVirtualNetworkUsagePermissionType(),
						},

						"use_public_ip_address": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Default:      string(dtl.Allow),
							ValidateFunc: validate.DevTestVirtualNetworkUsagePermissionType(),
						},
					},
				},
			},

			"tags": tags.Schema(),

			"unique_identifier": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDevTestVirtualNetworkCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Virtual Network creation")

	id := parse.NewDevTestVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.LabName, id.VirtualNetworkName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_dev_test_virtual_network", id.ID())
		}
	}

	description := d.Get("description").(string)
	t := d.Get("tags").(map[string]interface{})

	subnetsRaw := d.Get("subnet").([]interface{})
	subnets := expandDevTestVirtualNetworkSubnets(subnetsRaw, subscriptionId, id.ResourceGroup, id.VirtualNetworkName)

	parameters := dtl.VirtualNetwork{
		Tags: tags.Expand(t),
		VirtualNetworkProperties: &dtl.VirtualNetworkProperties{
			Description:     utils.String(description),
			SubnetOverrides: subnets,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LabName, id.VirtualNetworkName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestVirtualNetworkUpdate(d, meta)
}

func resourceArmDevTestVirtualNetworkRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DevTestVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LabName, id.VirtualNetworkName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("name", id.VirtualNetworkName)
	d.Set("lab_name", id.LabName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := read.VirtualNetworkProperties; props != nil {
		d.Set("description", props.Description)

		flattenedSubnets := flattenDevTestVirtualNetworkSubnets(props.SubnetOverrides)
		if err := d.Set("subnet", flattenedSubnets); err != nil {
			return fmt.Errorf("setting `subnet`: %+v", err)
		}

		// Computed fields
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	return tags.FlattenAndSet(d, read.Tags)
}

func resourceArmDevTestVirtualNetworkUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for DevTest Virtual Network creation")

	id := parse.NewDevTestVirtualNetworkID(subscriptionId, d.Get("resource_group_name").(string), d.Get("lab_name").(string), d.Get("name").(string))

	description := d.Get("description").(string)
	t := d.Get("tags").(map[string]interface{})

	subnetsRaw := d.Get("subnet").([]interface{})
	subnets := expandDevTestVirtualNetworkSubnets(subnetsRaw, subscriptionId, id.ResourceGroup, id.VirtualNetworkName)

	parameters := dtl.VirtualNetwork{
		Tags: tags.Expand(t),
		VirtualNetworkProperties: &dtl.VirtualNetworkProperties{
			Description:     utils.String(description),
			SubnetOverrides: subnets,
		},
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.LabName, id.VirtualNetworkName, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceArmDevTestVirtualNetworkRead(d, meta)
}

func resourceArmDevTestVirtualNetworkDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DevTestLabs.VirtualNetworksClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DevTestVirtualNetworkID(d.Id())
	if err != nil {
		return err
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.LabName, id.VirtualNetworkName, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] %s was not found - assuming removed!", *id)
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.LabName, id.VirtualNetworkName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}

func ValidateDevTestVirtualNetworkName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile("^[A-Za-z0-9_-]+$"),
		"Virtual Network Name can only include alphanumeric characters, underscores, hyphens.")
}

func expandDevTestVirtualNetworkSubnets(input []interface{}, subscriptionId, resourceGroupName, virtualNetworkName string) *[]dtl.SubnetOverride {
	results := make([]dtl.SubnetOverride, 0)
	// default found from the Portal
	name := fmt.Sprintf("%sSubnet", virtualNetworkName)
	idFmt := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s"
	subnetId := fmt.Sprintf(idFmt, subscriptionId, resourceGroupName, virtualNetworkName, name)
	if len(input) == 0 {
		result := dtl.SubnetOverride{
			ResourceID:                   utils.String(subnetId),
			LabSubnetName:                utils.String(name),
			UsePublicIPAddressPermission: dtl.Allow,
			UseInVMCreationPermission:    dtl.Allow,
		}
		results = append(results, result)
		return &results
	}

	for _, val := range input {
		v := val.(map[string]interface{})
		usePublicIPAddress := v["use_public_ip_address"].(string)
		useInVirtualMachineCreation := v["use_in_virtual_machine_creation"].(string)

		subnet := dtl.SubnetOverride{
			ResourceID:                   utils.String(subnetId),
			LabSubnetName:                utils.String(name),
			UsePublicIPAddressPermission: dtl.UsagePermissionType(usePublicIPAddress),
			UseInVMCreationPermission:    dtl.UsagePermissionType(useInVirtualMachineCreation),
		}
		results = append(results, subnet)
	}

	return &results
}

func flattenDevTestVirtualNetworkSubnets(input *[]dtl.SubnetOverride) []interface{} {
	outputs := make([]interface{}, 0)
	if input == nil {
		return outputs
	}

	for _, v := range *input {
		output := make(map[string]interface{})
		if v.LabSubnetName != nil {
			output["name"] = *v.LabSubnetName
		}
		output["use_public_ip_address"] = string(v.UsePublicIPAddressPermission)
		output["use_in_virtual_machine_creation"] = string(v.UseInVMCreationPermission)

		outputs = append(outputs, output)
	}

	return outputs
}
