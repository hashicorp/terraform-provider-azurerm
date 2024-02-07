// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package blueprints

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/blueprint"
	"github.com/hashicorp/go-azure-sdk/resource-manager/blueprints/2018-11-01-preview/publishedblueprint"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/blueprints/validate"
	mgValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func dataSourceBlueprintDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceBlueprintDefinitionRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.DefinitionName,
			},

			"scope_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.Any(
					azure.ValidateResourceID,
					mgValidate.ManagementGroupID,
				),
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

			"versions": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func dataSourceBlueprintDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Blueprints.BlueprintsClient
	publishedClient := meta.(*clients.Client).Blueprints.PublishedBlueprintsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := blueprint.NewScopedBlueprintID(d.Get("scope_id").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("the Blueprint Definition %q not found in Scope (%q): %+v", id.BlueprintName, id.ResourceScope, err)
		}

		return fmt.Errorf("read failed for Blueprint Definition (%q) in Sccope (%q): %+v", id.BlueprintName, id.ResourceScope, err)
	}

	d.SetId(id.ID())

	if m := resp.Model; m != nil {
		p := m.Properties

		d.Set("description", pointer.From(p.Description))
		d.Set("display_name", pointer.From(p.DisplayName))
		d.Set("last_modified", p.Status.LastModified)
		d.Set("time_created", p.Status.TimeCreated)
		d.Set("target_scope", p.TargetScope)

		publishedId := publishedblueprint.NewScopedBlueprintID(id.ResourceScope, id.BlueprintName)

		versionList := make([]string, 0)
		resp, err := publishedClient.List(ctx, publishedId)
		if err != nil {
			return fmt.Errorf("listing blue print versions for %s error: %+v", publishedId.String(), err)
		}

		if m := resp.Model; m != nil {
			for _, v := range *resp.Model {
				if v.Properties.BlueprintName != nil {
					versionList = append(versionList, *v.Properties.BlueprintName)
				}
			}
		}

		d.Set("versions", versionList)
	}
	return nil
}
