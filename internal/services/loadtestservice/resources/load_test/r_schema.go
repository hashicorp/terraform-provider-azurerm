// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package load_test

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2022-12-01/loadtests"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"name": {
			Description: "The name of this Load Test.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"resource_group_name": commonschema.ResourceGroupName(),
		"description": {
			Description: "Description of the resource.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},
		"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),
		"encryption": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key_url": {
						Description:  "The URI specifying the Key Vault and key to be used to encrypt data in this resource. The URI should include the key version.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
					"identity": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"type": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: validation.StringInSlice(loadtests.PossibleValuesForType(), false),
								},
								"identity_id": {
									Type:         pluginsdk.TypeString,
									Required:     true,
									ForceNew:     true,
									ValidateFunc: commonids.ValidateUserAssignedIdentityID,
								},
							},
						},
					},
				},
			},
		},
		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"data_plane_uri": {
			Description: "Resource data plane URI.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func (r Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LoadTestResourceSchema
			if err := metadata.DecodeDiff(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// If these values are not yet known, we want to avoid returning an error and will instead handle it in the Create/Update methods
			if metadata.IsKnownAt("encryption.0.identity.0.identity_id") && metadata.IsWhollyKnownAt("identity.0.identity_ids") {
				return r.EnsureEncryptionIdentityIDExistsInIdentity(model)
			}

			return nil
		},
	}
}
