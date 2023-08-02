// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package blueprints

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/publishedblueprint"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

	id := publishedblueprint.NewScopedVersionID(d.Get("scope_id").(string), d.Get("blueprint_name").(string), d.Get("version").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("Published Blueprint Version %q not found: %+v", id.String(), err)
		}

		return fmt.Errorf("reading Published Blueprint Version (%q): %+v", id.String(), err)
	}

	d.SetId(id.ID())

	if m := resp.Model; m != nil {
		p := m.Properties

		d.Set("type", pointer.From(m.Type))
		d.Set("target_scope", pointer.From(p.TargetScope))
		d.Set("display_name", pointer.From(p.DisplayName))
		d.Set("description", pointer.From(p.Description))

		if s := p.Status; s != nil {
			d.Set("time_created", pointer.From(s.TimeCreated))
			d.Set("last_modified", pointer.From(s.LastModified))
		}
	}
	return nil
}
