// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package proximity_placement_group

import (
	"context"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/compute/2022-03-01/proximityplacementgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

//go:generate go run ../../../../tools/generator-tests resourceidentity

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceProximityPlacementGroupCreateUpdate,
		Read:   resourceProximityPlacementGroupRead,
		Update: resourceProximityPlacementGroupCreateUpdate,
		Delete: resourceProximityPlacementGroupDelete,

		Importer: pluginsdk.ImporterValidatingIdentity(&proximityplacementgroups.ProximityPlacementGroupId{}),

		Identity: &schema.ResourceIdentity{
			SchemaFunc: pluginsdk.GenerateIdentitySchema(&proximityplacementgroups.ProximityPlacementGroupId{}),
		},

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: proximityPlacementGroupSchema(),

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("allowed_vm_sizes", func(ctx context.Context, old, new, meta any) bool {
				return len(old.(*pluginsdk.Set).List()) > 0 && len(new.(*pluginsdk.Set).List()) == 0
			}),
		),
	}
}
