package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/hashicorp/go-azure-helpers/response"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-07-01/compute"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmDedicatedHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmDedicatedHostCreate,
		Read:   resourceArmDedicatedHostRead,
		Update: resourceArmDedicatedHostUpdate,
		Delete: resourceArmDedicatedHostDelete,

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
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDedicatedHostName(),
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"host_group_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateDedicatedHostGroupName(),
			},

			// In order to follow convention for both protal and azure cli, `sku` here means
			// sku name only.
			"sku": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DSv3-Type1",
					"ESv3-Type1",
					"FSv2-Type2",
				}, false),
			},

			"platform_fault_domain": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},

			"auto_replace_on_failure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"license_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.DedicatedHostLicenseTypesNone),
					string(compute.DedicatedHostLicenseTypesWindowsServerHybrid),
					string(compute.DedicatedHostLicenseTypesWindowsServerPerpetual),
				}, false),
				Default: string(compute.DedicatedHostLicenseTypesNone),
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmDedicatedHostCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	hostGroupName := d.Get("host_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroupName, hostGroupName, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for present of existing Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_dedicated_host", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})
	parameters := compute.DedicatedHost{
		Location: utils.String(location),
		DedicatedHostProperties: &compute.DedicatedHostProperties{
			AutoReplaceOnFailure: utils.Bool(d.Get("auto_replace_on_failure").(bool)),
			LicenseType:          compute.DedicatedHostLicenseTypes(d.Get("license_type").(string)),
		},
		Sku: &compute.Sku{
			Name: utils.String(d.Get("sku").(string)),
		},
		Tags: tags.Expand(t),
	}
	if platformFaultDomain, ok := d.GetOk("platform_fault_domain"); ok {
		parameters.DedicatedHostProperties.PlatformFaultDomain = utils.Int32(int32(platformFaultDomain.(int)))
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroupName, hostGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}

	resp, err := client.Get(ctx, resourceGroupName, hostGroupName, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Dedicated Host %q (Host Group Name %q / Resource Group %q) ID", name, hostGroupName, resourceGroupName)
	}
	d.SetId(*resp.ID)

	return resourceArmDedicatedHostRead(d, meta)
}

func resourceArmDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	hostGroupName := id.Path["hostGroups"]
	name := id.Path["hosts"]

	resp, err := client.Get(ctx, resourceGroupName, hostGroupName, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Dedicated Host %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroupName)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("host_group_name", hostGroupName)
	d.Set("sku", resp.Sku.Name)
	if props := resp.DedicatedHostProperties; props != nil {
		d.Set("platform_fault_domain", props.PlatformFaultDomain)
		d.Set("auto_replace_on_failure", props.AutoReplaceOnFailure)
		d.Set("license_type", props.LicenseType)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroupName := d.Get("resource_group_name").(string)
	hostGroupName := d.Get("host_group_name").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := compute.DedicatedHostUpdate{
		DedicatedHostProperties: &compute.DedicatedHostProperties{
			AutoReplaceOnFailure: utils.Bool(d.Get("auto_replace_on_failure").(bool)),
			LicenseType:          compute.DedicatedHostLicenseTypes(d.Get("license_type").(string)),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.Update(ctx, resourceGroupName, hostGroupName, name, parameters)
	if err != nil {
		return fmt.Errorf("Error updating Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}

	return resourceArmDedicatedHostRead(d, meta)
}

func resourceArmDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroupName := id.ResourceGroup
	hostGroupName := id.Path["hostGroups"]
	name := id.Path["hosts"]

	future, err := client.Delete(ctx, resourceGroupName, hostGroupName, name)
	if err != nil {
		return fmt.Errorf("Error deleting Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", name, hostGroupName, resourceGroupName, err)
		}
	}

	// API has bug, which appears to be eventually consistent. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/8137
	log.Printf("[DEBUG] Waiting for Dedicated Host %q (Host Group Name %q / Resource Group %q) to disappear", name, hostGroupName, resourceGroupName)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   dedicatedHostDeletedRefreshFunc(ctx, client, id),
		MinTimeout:                5 * time.Second,
		ContinuousTargetOccurence: 5,
	}

	if features.SupportsCustomTimeouts() {
		stateConf.Timeout = d.Timeout(schema.TimeoutDelete)
	} else {
		stateConf.Timeout = 10 * time.Minute
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Dedicated Host %q (Host Group Name %q / Resource Group %q) to become available: %+v", name, hostGroupName, resourceGroupName, err)
	}

	return nil
}

func dedicatedHostDeletedRefreshFunc(ctx context.Context, client *compute.DedicatedHostsClient, id *azure.ResourceID) resource.StateRefreshFunc {
	resourceGroupName := id.ResourceGroup
	hostGroupName := id.Path["hostGroups"]
	name := id.Path["hosts"]
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, hostGroupName, name, "")
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}
			return nil, "", fmt.Errorf("Error issuing read request in  dedicatedHostDeletedRefreshFunc: %+v", err)
		}

		return res, "Exists", nil
	}
}
