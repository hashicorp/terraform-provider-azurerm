package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2021-11-01/compute"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceDedicatedHost() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceDedicatedHostCreate,
		Read:   resourceDedicatedHostRead,
		Update: resourceDedicatedHostUpdate,
		Delete: resourceDedicatedHostDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.DedicatedHostID(id)
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
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DedicatedHostName(),
			},

			"location": azure.SchemaLocation(),

			"dedicated_host_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.DedicatedHostGroupID,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				ForceNew: true,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DADSv5-Type1",
					"DASv4-Type1",
					"DASv4-Type2",
					"DASv5-Type1",
					"DCSv2-Type1",
					"DDSv4-Type1",
					"DDSv4-Type2",
					"DDSv5-Type1",
					"DSv3-Type1",
					"DSv3-Type2",
					"DSv3-Type3",
					"DSv3-Type4",
					"DSv4-Type1",
					"DSv4-Type2",
					"DSv5-Type1",
					"EADSv5-Type1",
					"EASv4-Type1",
					"EASv4-Type2",
					"EASv5-Type1",
					"EDSv4-Type1",
					"EDSv4-Type2",
					"EDSv5-Type1",
					"ESv3-Type1",
					"ESv3-Type2",
					"ESv3-Type3",
					"ESv3-Type4",
					"ESv4-Type1",
					"ESv4-Type2",
					"ESv5-Type1",
					"FSv2-Type2",
					"FSv2-Type3",
					"FSv2-Type4",
					"FXmds-Type1",
					"LSv2-Type1",
					"MDMSv2MedMem-Type1",
					"MDSv2MedMem-Type1",
					"MMSv2MedMem-Type1",
					"MS-Type1",
					"MSm-Type1",
					"MSmv2-Type1",
					"MSv2-Type1",
					"MSv2MedMem-Type1",
					"NVASv4-Type1",
					"NVSv3-Type1",
				}, false),
			},

			"platform_fault_domain": {
				Type:     pluginsdk.TypeInt,
				ForceNew: true,
				Required: true,
			},

			"auto_replace_on_failure": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"license_type": {
				Type:     pluginsdk.TypeString,
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

func resourceDedicatedHostCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	hostGroupId, err := parse.DedicatedHostGroupID(d.Get("dedicated_host_group_id").(string))
	if err != nil {
		return err
	}

	id := parse.NewDedicatedHostID(hostGroupId.SubscriptionId, hostGroupId.ResourceGroup, hostGroupId.HostGroupName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_dedicated_host", id.ID())
		}
	}

	parameters := compute.DedicatedHost{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		DedicatedHostProperties: &compute.DedicatedHostProperties{
			AutoReplaceOnFailure: utils.Bool(d.Get("auto_replace_on_failure").(bool)),
			LicenseType:          compute.DedicatedHostLicenseTypes(d.Get("license_type").(string)),
			PlatformFaultDomain:  utils.Int32(int32(d.Get("platform_fault_domain").(int))),
		},
		Sku: &compute.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDedicatedHostRead(d, meta)
}

func resourceDedicatedHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	groupsClient := meta.(*clients.Client).Compute.DedicatedHostGroupsClient
	hostsClient := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	hostGroupId := parse.NewDedicatedHostGroupID(id.SubscriptionId, id.ResourceGroup, id.HostGroupName)
	group, err := groupsClient.Get(ctx, hostGroupId.ResourceGroup, hostGroupId.HostGroupName, "")
	if err != nil {
		if utils.ResponseWasNotFound(group.Response) {
			log.Printf("[INFO] Parent %s does not exist - removing from state", hostGroupId)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", hostGroupId, err)
	}

	resp, err := hostsClient.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Dedicated Host %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Dedicated Host %q (Host Group Name %q / Resource Group %q): %+v", id.HostName, id.HostGroupName, id.ResourceGroup, err)
	}

	d.Set("name", id.HostName)
	d.Set("dedicated_host_group_id", hostGroupId.ID())

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	d.Set("sku_name", resp.Sku.Name)
	if props := resp.DedicatedHostProperties; props != nil {
		d.Set("auto_replace_on_failure", props.AutoReplaceOnFailure)
		d.Set("license_type", props.LicenseType)

		platformFaultDomain := 0
		if props.PlatformFaultDomain != nil {
			platformFaultDomain = int(*props.PlatformFaultDomain)
		}
		d.Set("platform_fault_domain", platformFaultDomain)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceDedicatedHostUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	parameters := compute.DedicatedHostUpdate{
		DedicatedHostProperties: &compute.DedicatedHostProperties{
			AutoReplaceOnFailure: utils.Bool(d.Get("auto_replace_on_failure").(bool)),
			LicenseType:          compute.DedicatedHostLicenseTypes(d.Get("license_type").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", *id, err)
	}

	return resourceDedicatedHostRead(d, meta)
}

func resourceDedicatedHostDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DedicatedHostID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.HostGroupName, id.HostName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
		}
	}

	// API has bug, which appears to be eventually consistent. Tracked by this issue: https://github.com/Azure/azure-rest-api-specs/issues/8137
	log.Printf("[DEBUG] Waiting for %s to be fully deleted..", *id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Exists"},
		Target:                    []string{"NotFound"},
		Refresh:                   dedicatedHostDeletedRefreshFunc(ctx, client, *id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 20,
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be fully deleted: %+v", *id, err)
	}

	return nil
}

func dedicatedHostDeletedRefreshFunc(ctx context.Context, client *compute.DedicatedHostsClient, id parse.DedicatedHostId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.HostGroupName, id.HostName, "")
		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
