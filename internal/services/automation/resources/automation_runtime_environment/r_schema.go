// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package automation_runtime_environment

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2024-10-23/runtimeenvironment"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (m Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			rd := metadata.ResourceDiff

			if rdp, ok := rd.GetOk("runtime_default_packages"); ok && rdp != nil {
				if rl, ok := rd.GetOk("runtime_language"); ok {
					if rl.(string) == "Python" {
						return errors.New("runtime_default_packages cannot be set for runtime_language `Python`")
					}
				}
			}

			if rd.HasChange("runtime_default_packages") {
				old, new := rd.GetChange("runtime_default_packages")
				oldMap, newMap := old.(map[string]any), new.(map[string]any)

				// Azure API limitation: Runtime environment packages cannot be removed once added
				if len(oldMap) > len(newMap) {
					if err := rd.ForceNew("runtime_default_packages"); err != nil {
						return err
					}
				}

				// Azure API limitation: Package names cannot change, only versions can be updated
				for k := range oldMap {
					if _, ok := newMap[k]; !ok {
						if err := rd.ForceNew("runtime_default_packages"); err != nil {
							return err
						}
					}
				}
			}

			return nil
		},
		Timeout: 5 * time.Minute,
	}
}

func (m Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"location": commonschema.Location(),

		"automation_account_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: runtimeenvironment.ValidateAutomationAccountID,
		},

		"runtime_language": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				"Python",
				"PowerShell",
			}, false),
		},

		"runtime_version": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			ForceNew:     true,
		},

		"runtime_default_packages": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem:     &pluginsdk.Schema{Type: pluginsdk.TypeString},
		},

		"tags": {
			Type:         pluginsdk.TypeMap,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: tags.Validate,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (m Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
