// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhostgroups"
	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2021-11-01/dedicatedhosts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/validate"
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
			_, err := dedicatedhosts.ParseHostID(id)
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

			"dedicated_host_group_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: dedicatedhostgroups.ValidateHostGroupID,
			},

			"location": commonschema.Location(),

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
					"LSv3-Type1",
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
					// TODO: remove `None` in 4.0 in favour of this field being set to an empty string (since it's optional)
					string(dedicatedhosts.DedicatedHostLicenseTypesNone),
					string(dedicatedhosts.DedicatedHostLicenseTypesWindowsServerHybrid),
					string(dedicatedhosts.DedicatedHostLicenseTypesWindowsServerPerpetual),
				}, false),
				Default: string(dedicatedhosts.DedicatedHostLicenseTypesNone),
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceDedicatedHostCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	hostGroupId, err := dedicatedhostgroups.ParseHostGroupID(d.Get("dedicated_host_group_id").(string))
	if err != nil {
		return err
	}

	id := dedicatedhosts.NewHostID(hostGroupId.SubscriptionId, hostGroupId.ResourceGroupName, hostGroupId.HostGroupName, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id, dedicatedhosts.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_dedicated_host", id.ID())
		}
	}

	licenseType := dedicatedhosts.DedicatedHostLicenseTypes(d.Get("license_type").(string))
	payload := dedicatedhosts.DedicatedHost{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &dedicatedhosts.DedicatedHostProperties{
			AutoReplaceOnFailure: utils.Bool(d.Get("auto_replace_on_failure").(bool)),
			LicenseType:          &licenseType,
			PlatformFaultDomain:  utils.Int64(int64(d.Get("platform_fault_domain").(int))),
		},
		Sku: dedicatedhosts.Sku{
			Name: utils.String(d.Get("sku_name").(string)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDedicatedHostRead(d, meta)
}

func resourceDedicatedHostRead(d *pluginsdk.ResourceData, meta interface{}) error {
	hostsClient := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhosts.ParseHostID(d.Id())
	if err != nil {
		return err
	}

	resp, err := hostsClient.Get(ctx, *id, dedicatedhosts.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.HostName)
	d.Set("dedicated_host_group_id", dedicatedhostgroups.NewHostGroupID(id.SubscriptionId, id.ResourceGroupName, id.HostGroupName).ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))
		d.Set("sku_name", model.Sku.Name)
		if props := model.Properties; props != nil {
			d.Set("auto_replace_on_failure", props.AutoReplaceOnFailure)
			d.Set("license_type", string(pointer.From(props.LicenseType)))

			platformFaultDomain := 0
			if props.PlatformFaultDomain != nil {
				platformFaultDomain = int(*props.PlatformFaultDomain)
			}
			d.Set("platform_fault_domain", platformFaultDomain)
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceDedicatedHostUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhosts.ParseHostID(d.Id())
	if err != nil {
		return err
	}

	payload := dedicatedhosts.DedicatedHostUpdate{}

	if d.HasChanges("auto_replace_on_failure", "license_type") {
		payload.Properties = &dedicatedhosts.DedicatedHostProperties{}
		if d.HasChange("auto_replace_on_failure") {
			payload.Properties.AutoReplaceOnFailure = utils.Bool(d.Get("auto_replace_on_failure").(bool))
		}
		if d.HasChange("license_type") {
			licenseType := dedicatedhosts.DedicatedHostLicenseTypes(d.Get("license_type").(string))
			payload.Properties.LicenseType = &licenseType
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.UpdateThenPoll(ctx, *id, payload); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceDedicatedHostRead(d, meta)
}

func resourceDedicatedHostDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Compute.DedicatedHostsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := dedicatedhosts.ParseHostID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
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

func dedicatedHostDeletedRefreshFunc(ctx context.Context, client *dedicatedhosts.DedicatedHostsClient, id dedicatedhosts.HostId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id, dedicatedhosts.DefaultGetOperationOptions())
		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return "NotFound", "NotFound", nil
			}

			return nil, "", fmt.Errorf("checking if %s has been deleted: %+v", id, err)
		}

		return res, "Exists", nil
	}
}
