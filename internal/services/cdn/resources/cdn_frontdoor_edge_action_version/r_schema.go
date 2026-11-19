// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package cdn_frontdoor_edge_action_version

import (
	"context"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cdn/2026-10-01/edgeactionversions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9]{1,50}$`),
				"`name` must be between 1 and 50 characters in length and may only contain letters and numbers",
			),
		},

		"cdn_frontdoor_edge_action_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: edgeactionversions.ValidateEdgeActionID,
		},

		"deployment_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(edgeactionversions.PossibleValuesForEdgeActionVersionDeploymentType(), false),
		},

		"code_content_base64": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsBase64,
			RequiredWith: []string{"code_file_name"},
		},

		"code_file_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			RequiredWith: []string{"code_content_base64"},
		},

		"tags": commonschema.Tags(),
	}
}

func (Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"default_version_enabled": {
			Type:     pluginsdk.TypeBool,
			Computed: true,
		},

		"validation_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff
			if rd.Id() == "" || !rd.HasChanges("code_content_base64", "code_file_name") {
				return nil
			}

			oldCode, newCode := rd.GetChange("code_content_base64")
			// `validation_status` is only returned by the API once code has been deployed to the version, which also
			// covers code deployed outside of Terraform (e.g. before the version was imported), since the deployed code
			// itself can't be read back in a form which can be compared with the configuration
			hadCode := oldCode.(string) != "" || rd.Get("validation_status").(string) != ""

			// the API can't remove deployed code from a version, and can't redeploy the code of the
			// default version once it has been deployed, so a new version needs to be created instead.
			// a value which isn't known until apply is reported as empty, so it isn't treated as removing the code
			codeRemoved := hadCode && rd.NewValueKnown("code_content_base64") && newCode.(string) == ""
			defaultVersionRedeployed := hadCode && rd.Get("default_version_enabled").(bool)

			if codeRemoved || defaultVersionRedeployed {
				for _, key := range []string{"code_content_base64", "code_file_name"} {
					if rd.HasChange(key) {
						if err := rd.ForceNew(key); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
	}
}
