// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSecurityCenterAutoProvisioning() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSecurityCenterAutoProvisioningCreateUpdate,
		Read:   resourceSecurityCenterAutoProvisioningRead,
		Update: resourceSecurityCenterAutoProvisioningCreateUpdate,
		Delete: resourceSecurityCenterAutoProvisioningDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AutoProvisioningSettingID(id)
			return err
		}),

		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AutoProvisioningV0ToV1{},
		}),
		SchemaVersion: 1,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"auto_provision": {
				Type:     pluginsdk.TypeString,
				Required: true,
				// NOTE: the API seems case insensitive to this string value, 'ON', 'On', 'on' all work
				ValidateFunc: validation.StringInSlice([]string{
					string(security.AutoProvisionOn),
					string(security.AutoProvisionOff),
				}, false),
			},
		},
	}
}

func resourceSecurityCenterAutoProvisioningCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutoProvisioningClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// No need for import check as there's always single resource called 'default'
	// - it cannot be deleted, all this does is set a string property to: "on" or "off"

	// Build settings struct with auto_provision value
	settings := security.AutoProvisioningSetting{
		AutoProvisioningSettingProperties: &security.AutoProvisioningSettingProperties{
			AutoProvision: security.AutoProvision(d.Get("auto_provision").(string)),
		},
	}

	// NOTE: 'default' is the only valid name currently supported by the API
	// No other names can be created and the 'default' resource can not be destroyed
	id := parse.NewAutoProvisioningSettingID(subscriptionId, "default")

	// There is no update function or operation in the API, only create
	if _, err := client.Create(ctx, id.Name, settings); err != nil {
		return fmt.Errorf("updating auto-provisioning setting for %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceSecurityCenterAutoProvisioningRead(d, meta)
}

func resourceSecurityCenterAutoProvisioningRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AutoProvisioningClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutoProvisioningSettingID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state", err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving auto-provisioning setting %s: %+v", id, err)
	}

	if properties := resp.AutoProvisioningSettingProperties; properties != nil {
		d.Set("auto_provision", properties.AutoProvision)
	}

	return nil
}

func resourceSecurityCenterAutoProvisioningDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	// The API has no delete operation
	// Instead we reset back to 'Off' which is the default

	client := meta.(*clients.Client).SecurityCenter.AutoProvisioningClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutoProvisioningSettingID(d.Id())
	if err != nil {
		return err
	}

	settings := security.AutoProvisioningSetting{
		AutoProvisioningSettingProperties: &security.AutoProvisioningSettingProperties{
			AutoProvision: security.AutoProvisionOff,
		},
	}

	// There is no update function or operation in the API, only create
	if _, err := client.Create(ctx, id.Name, settings); err != nil {
		return fmt.Errorf("resetting Security Center auto provisioning to 'Off': %+v", err)
	}

	return nil
}
