// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAdvancedThreatProtection() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAdvancedThreatProtectionCreateUpdate,
		Read:   resourceAdvancedThreatProtectionRead,
		Update: resourceAdvancedThreatProtectionCreateUpdate,
		Delete: resourceAdvancedThreatProtectionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AdvancedThreatProtectionID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.AdvancedThreatProtectionV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"target_resource_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},
		},
	}
}

func resourceAdvancedThreatProtectionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewAdvancedThreatProtectionId(d.Get("target_resource_id").(string))
	if d.IsNewResource() {
		server, err := client.Get(ctx, id.TargetResourceID)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("checking for presence of existing Advanced Threat Protection for %q: %+v", id.TargetResourceID, err)
			}
		}

		if server.ID != nil && *server.ID != "" && server.IsEnabled != nil && *server.IsEnabled {
			return tf.ImportAsExistsError("azurerm_advanced_threat_protection", id.ID())
		}
	}

	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: utils.Bool(d.Get("enabled").(bool)),
		},
	}

	if _, err := client.Create(ctx, id.TargetResourceID, setting); err != nil {
		return fmt.Errorf("updating Advanced Threat protection for %q: %+v", id.TargetResourceID, err)
	}

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("internal-error: context had no deadline")
	}

	// the API appears to be eventually consistent, tracked on https://github.com/Azure/azure-rest-api-specs/issues/25232
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"diff"},
		Target:  []string{"consistent"},
		Refresh: func() (result interface{}, state string, err error) {
			resp, err := client.Get(ctx, id.TargetResourceID)
			if err != nil {
				return resp, "error", err
			}
			if atpp := resp.AdvancedThreatProtectionProperties; atpp != nil {
				respEnabled := atpp.IsEnabled != nil && *atpp.IsEnabled
				if respEnabled == d.Get("enabled").(bool) {
					return resp, "consistent", nil
				} else {
					return resp, "diff", nil
				}
			}
			return resp, "error", fmt.Errorf("Properties was nil")
		},
		MinTimeout:                1 * time.Minute,
		ContinuousTargetOccurence: 3,
		Timeout:                   time.Until(deadline),
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for provisioning state of advanced threat protection: %+v", err)
	}

	d.SetId(id.ID())
	return resourceAdvancedThreatProtectionRead(d, meta)
}

func resourceAdvancedThreatProtectionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AdvancedThreatProtectionID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.TargetResourceID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Advanced Threat Protection was not found for %q: %+v", id.TargetResourceID, err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Advanced Threat Protection status for %q: %+v", id.TargetResourceID, err)
	}

	d.Set("target_resource_id", id.TargetResourceID)
	if atpp := resp.AdvancedThreatProtectionProperties; atpp != nil {
		d.Set("enabled", resp.IsEnabled)
	}

	return nil
}

func resourceAdvancedThreatProtectionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.AdvancedThreatProtectionClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AdvancedThreatProtectionID(d.Id())
	if err != nil {
		return err
	}

	// there is no delete.. so lets just do best effort and set it to false?
	setting := security.AdvancedThreatProtectionSetting{
		AdvancedThreatProtectionProperties: &security.AdvancedThreatProtectionProperties{
			IsEnabled: utils.Bool(false),
		},
	}

	if _, err := client.Create(ctx, id.TargetResourceID, setting); err != nil {
		return fmt.Errorf("removing Advanced Threat Protection for %q: %+v", id.TargetResourceID, err)
	}

	return nil
}
