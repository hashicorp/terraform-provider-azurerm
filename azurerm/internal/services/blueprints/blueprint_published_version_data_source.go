package blueprints

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	mgValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceBlueprintPublishedVersion() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBlueprintPublishedVersionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"blueprint_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"scope_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					azure.ValidateResourceID,
					mgValidate.ManagementGroupID,
				),
			},

			"version": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// Computed
			"description": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"last_modified": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"target_scope": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"time_created": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBlueprintPublishedVersionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.PublishedBlueprintsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	scope := d.Get("scope_id").(string)
	blueprintName := d.Get("blueprint_name").(string)
	versionID := d.Get("version").(string)

	resp, err := client.Get(ctx, scope, blueprintName, versionID)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("Published Blueprint Version %q not found: %+v", versionID, err)
		}

		return fmt.Errorf("Read failed for Published Blueprint (%q) Version (%q): %+v", blueprintName, versionID, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("Failed to retrieve ID for Blueprint %q Version %q", blueprintName, versionID)
	} else {
		d.SetId(*resp.ID)
	}

	if resp.Type != nil {
		d.Set("type", resp.Type)
	}

	if resp.Status != nil {
		if resp.Status.TimeCreated != nil {
			d.Set("time_created", resp.Status.TimeCreated.String())
		}

		if resp.Status.LastModified != nil {
			d.Set("last_modified", resp.Status.LastModified.String())
		}
	}

	d.Set("target_scope", resp.TargetScope)

	if resp.DisplayName != nil {
		d.Set("display_name", resp.DisplayName)
	}

	if resp.Description != nil {
		d.Set("description", resp.Description)
	}

	return nil
}
