// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_desktop_application

import (
	"context"
	"errors"
	"regexp"

	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/application"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2025-10-10/applicationgroup"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func virtualDesktopApplicationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description: "The name of the Virtual Desktop Application.",
			Type:        pluginsdk.TypeString,
			Required:    true,
			ForceNew:    true,
			ValidateFunc: validation.All(
				// NOTE: name string length is incorrect in the API specs, see: https://github.com/Azure/azure-rest-api-specs/issues/46813
				validation.StringLenBetween(1, 260),
				validation.StringMatch(
					regexp.MustCompile(`^[A-Za-z0-9@._ -]*$`),
					"Virtual desktop application name must be 1 - 260 characters long and may only contain letters, numbers, spaces, periods, underscores, hyphens, and @.",
				),
			),
		},

		"application_group_id": {
			Description:  "Resource ID for a Virtual Desktop Application Group to associate with the Virtual Desktop Application.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: applicationgroup.ValidateApplicationGroupID,
		},

		"command_line_argument_policy": {
			Description:  "Whether this published application can be launched with command line arguments provided by the client, command line arguments specified at publish time, or no command line arguments at all.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(application.PossibleValuesForCommandLineSetting(), false),
		},

		"application_type": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			Default:      string(application.RemoteApplicationTypeInBuilt),
			ValidateFunc: validation.StringInSlice(application.PossibleValuesForRemoteApplicationType(), false),
		},

		"command_line_arguments": {
			Description: "Command Line Arguments for Virtual Desktop Application.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
		},

		"description": {
			Description:  "Option to set a description for the Virtual Desktop Application.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringLenBetween(1, 512),
		},

		"friendly_name": {
			Description: "Option to set a friendly name for the Virtual Desktop Application.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C The API will use the value in `name` as the default
			Computed:     true,
			ValidateFunc: validation.StringLenBetween(1, 64),
		},

		"icon_index": {
			Description: "The index of the icon you wish to use.",
			Type:        pluginsdk.TypeInt,
			Optional:    true,
		},

		"icon_path": {
			Description: "The path for an icon which will be used for this Virtual Desktop Application.",
			Type:        pluginsdk.TypeString,
			Optional:    true,
			// NOTE: O+C The API will use the value in `path` as the default
			Computed: true,
		},

		"msix_package_application_id": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"msix_package_family_name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"path": {
			Description:  "The file path location of the app on the Virtual Desktop OS.",
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"show_in_portal": {
			Description: "Whether to show the RemoteApp program in the RD Web Access server.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
		},
	}
}

func virtualDesktopApplicationCustomizeDiff(ctx context.Context, d *pluginsdk.ResourceDiff, _ any) error {
	_, pathOk := d.GetOk("path")
	applicationType, applicationTypeOk := d.GetOk("application_type")
	msixPackageApplicationId := d.GetRawConfig().AsValueMap()["msix_package_application_id"]
	msixPackageFamilyName := d.GetRawConfig().AsValueMap()["msix_package_family_name"]

	if applicationTypeOk && applicationType == string(application.RemoteApplicationTypeMsixApplication) {
		if pathOk {
			return errors.New("`path` cannot be set when `application_type` is `MsixApplication`")
		}

		if msixPackageApplicationId.IsNull() || msixPackageFamilyName.IsNull() {
			return errors.New("`msix_package_application_id` and `msix_package_family_name` must be set when `application_type` is `MsixApplication`")
		}
	} else {
		if !pathOk {
			return errors.New("`path` must be set when `application_type` is not `MsixApplication`")
		}

		if !msixPackageApplicationId.IsNull() || !msixPackageFamilyName.IsNull() {
			return errors.New("`application_type` must be set as `MsixApplication` when `msix_package_application_id` or `msix_package_family_name` is set")
		}
	}

	return nil
}
