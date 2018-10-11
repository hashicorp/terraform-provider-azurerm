package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/devtestlabs/mgmt/2016-05-15/dtl"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDevTestVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDevTestVirtualMachineCreateUpdate,
		Read:   resourceArmDevTestVirtualMachineRead,
		Update: resourceArmDevTestVirtualMachineCreateUpdate,
		Delete: resourceArmDevTestVirtualMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				// The virtual machine name must be between 1 and 62 characters and cannot contain any spaces or special characters. The name may contain letters, numbers, or '-'. However, it must begin and end with a letter or number, and cannot be all numbers.
				// The name must be between 1 and 15 characters, cannot be entirely numeric, and cannot contain most special characters
			},

			"lab_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DevTestLabName(),
			},

			// There's a bug in the Azure API where this is returned in lower-case
			// BUG: https://github.com/Azure/azure-rest-api-specs/issues/3964
			"resource_group_name": resourceGroupNameDiffSuppressSchema(),

			"location": locationSchema(),

			"os_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Linux",
					"Windows",
				}, false),
			},

			"size": {
				Type:     schema.TypeString,
				Required: true,
				// since this isn't returned from the API
				ForceNew: true,
			},

			"username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"Standard",
					"Premium",
				}, false),
			},

			"lab_subnet_name": {
				Type:     schema.TypeString,
				Required: true,
				// since this isn't returned from the API
				ForceNew: true,
			},

			"lab_virtual_network_id": {
				Type:     schema.TypeString,
				Required: true,
				// since this isn't returned from the API
				ForceNew: true,
			},

			"allow_claim": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"disallow_public_ip_address": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"password": {
				Type:     schema.TypeString,
				Optional: true,
				// since this isn't returned from the API
				ForceNew:  true,
				Sensitive: true,
			},

			"ssh_key": {
				Type:     schema.TypeString,
				Optional: true,
				// since this isn't returned from the API
				ForceNew: true,
			},

			"gallery_image_reference": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"offer": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"publisher": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"sku": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			"inbound_nat_rule": azure.SchemaDevTestVirtualMachineInboundNatRule(),

			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": tagsSchema(),

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"unique_identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArmDevTestVirtualMachineCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestVirtualMachinesClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for DevTest Lab Virtual Machine creation")

	name := d.Get("name").(string)
	labName := d.Get("lab_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	tags := d.Get("tags").(map[string]interface{})

	allowClaim := d.Get("allow_claim").(bool)
	disallowPublicIPAddress := d.Get("disallow_public_ip_address").(bool)
	labSubnetName := d.Get("lab_subnet_name").(string)
	labVirtualNetworkId := d.Get("lab_virtual_network_id").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	osType := d.Get("os_type").(string)
	notes := d.Get("notes").(string)
	password := d.Get("password").(string)
	sshKey := d.Get("ssh_key").(string)
	size := d.Get("size").(string)
	storageType := d.Get("storage_type").(string)
	username := d.Get("username").(string)

	galleryImageReferenceRaw := d.Get("gallery_image_reference").([]interface{})
	galleryImageReference := expandDevTestLabVirtualMachineGalleryImageReference(galleryImageReferenceRaw, osType)

	natRulesRaw := d.Get("inbound_nat_rule").(*schema.Set)
	natRules := azure.ExpandDevTestLabVirtualMachineNatRules(natRulesRaw)

	if len(natRules) > 0 && disallowPublicIPAddress {
		return fmt.Errorf("If `inbound_nat_rule` is specified then `disallow_public_ip_address` must be set to true.")
	}

	nic := dtl.NetworkInterfaceProperties{}
	if disallowPublicIPAddress {
		nic.SharedPublicIPAddressConfiguration = &dtl.SharedPublicIPAddressConfiguration{
			InboundNatRules: &natRules,
		}
	}

	authenticateViaSsh := sshKey != ""
	parameters := dtl.LabVirtualMachine{
		Location: utils.String(location),
		LabVirtualMachineProperties: &dtl.LabVirtualMachineProperties{
			AllowClaim:                 utils.Bool(allowClaim),
			IsAuthenticationWithSSHKey: utils.Bool(authenticateViaSsh),
			DisallowPublicIPAddress:    utils.Bool(disallowPublicIPAddress),
			GalleryImageReference:      galleryImageReference,
			LabSubnetName:              utils.String(labSubnetName),
			LabVirtualNetworkID:        utils.String(labVirtualNetworkId),
			NetworkInterface:           &nic,
			OsType:                     utils.String(osType),
			Notes:                      utils.String(notes),
			Password:                   utils.String(password),
			Size:                       utils.String(size),
			SSHKey:                     utils.String(sshKey),
			StorageType:                utils.String(storageType),
			UserName:                   utils.String(username),
		},
		Tags: expandTags(tags),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, labName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating/updating DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation/update of DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, labName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read DevTest Virtual Machine %q (Lab %q / Resource Group %q) ID", name, labName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmDevTestVirtualMachineRead(d, meta)
}

func resourceArmDevTestVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestVirtualMachinesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	labName := id.Path["labs"]
	name := id.Path["virtualmachines"]

	read, err := client.Get(ctx, resourceGroup, labName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			log.Printf("[DEBUG] DevTest Virtual Machine %q was not found in Lab %q / Resource Group %q - removing from state!", name, labName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	d.Set("name", read.Name)
	d.Set("lab_name", labName)
	d.Set("resource_group_name", resourceGroup)
	if location := read.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := read.LabVirtualMachineProperties; props != nil {
		d.Set("allow_claim", props.AllowClaim)
		d.Set("disallow_public_ip_address", props.DisallowPublicIPAddress)
		d.Set("os_type", props.OsType)
		d.Set("notes", props.Notes)
		d.Set("size", props.Size)
		d.Set("storage_type", props.StorageType)
		d.Set("username", props.UserName)

		flattenedImage := flattenDevTestVirtualMachineGalleryImage(props.GalleryImageReference)
		if err := d.Set("gallery_image_reference", flattenedImage); err != nil {
			return fmt.Errorf("Error flattening `gallery_image_reference`: %+v", err)
		}

		// Computed fields
		d.Set("fqdn", props.Fqdn)
		d.Set("unique_identifier", props.UniqueIdentifier)
	}

	flattenAndSetTags(d, read.Tags)

	return nil
}

func resourceArmDevTestVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).devTestVirtualMachinesClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	labName := id.Path["labs"]
	name := id.Path["virtualmachines"]

	read, err := client.Get(ctx, resourceGroup, labName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(read.Response) {
			// deleted outside of TF
			log.Printf("[DEBUG] DevTest Virtual Machine %q was not found in Lab %q / Resource Group %q - assuming removed!", name, labName, resourceGroup)
			return nil
		}

		return fmt.Errorf("Error retrieving DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	future, err := client.Delete(ctx, resourceGroup, labName, name)
	if err != nil {
		return fmt.Errorf("Error deleting DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for the deletion of DevTest Virtual Machine %q (Lab %q / Resource Group %q): %+v", name, labName, resourceGroup, err)
	}

	return err
}

func expandDevTestLabVirtualMachineGalleryImageReference(input []interface{}, osType string) *dtl.GalleryImageReference {
	if len(input) == 0 {
		return nil
	}

	v := input[0].(map[string]interface{})
	offer := v["offer"].(string)
	publisher := v["publisher"].(string)
	sku := v["sku"].(string)
	version := v["version"].(string)

	return &dtl.GalleryImageReference{
		Offer:     utils.String(offer),
		OsType:    utils.String(osType),
		Publisher: utils.String(publisher),
		Sku:       utils.String(sku),
		Version:   utils.String(version),
	}
}

func flattenDevTestVirtualMachineGalleryImage(input *dtl.GalleryImageReference) []interface{} {
	results := make([]interface{}, 0)

	if input != nil {
		output := make(map[string]interface{}, 0)

		if input.Offer != nil {
			output["offer"] = *input.Offer
		}

		if input.Publisher != nil {
			output["publisher"] = *input.Publisher
		}

		if input.Sku != nil {
			output["sku"] = *input.Sku
		}

		if input.Version != nil {
			output["version"] = *input.Version
		}

		results = append(results, output)
	}

	return results
}
