// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package geo_catalog

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9-]{3,24}$"),
				"GeoCatalog name must consist only of letters, numbers, hyphens (-), and have length between 3 and 24 characters (inclusive)",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if rawIdentityType, ok := metadata.ResourceDiff.GetOk("identity.0.type"); ok {
				identityType := identity.Type(rawIdentityType.(string))
				switch identityType {
				case identity.TypeSystemAssigned, identity.TypeSystemAssignedUserAssigned:
					return fmt.Errorf("`%s` and `%s` resource identities are not supported", identity.TypeSystemAssigned, identity.TypeSystemAssignedUserAssigned)
				case identity.TypeUserAssigned:
					identityIds := metadata.ResourceDiff.GetRawConfig().AsValueMap()["identity"].AsValueSlice()[0].AsValueMap()["identity_ids"]
					if identityIds.IsNull() || identityIds.LengthInt() == 0 {
						return fmt.Errorf("`identity_ids` must be specified when `type` is set to `%s`", identity.TypeUserAssigned)
					}
				}
			}

			return nil
		},
	}
}
