// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func iothubDpsCertificateSchema() map[string]*pluginsdk.Schema {
	s := map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Iot Device Provisioning Service Certificate resource.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubDpsCertificateName,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"iot_dps_name": {
			Description:  "The name of the IoT Device Provisioning Service that this certificate will be attached to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IoTHubName,
		},

		"certificate_content": {
			Description:  "The Base-64 representation of the X509 leaf certificate .cer file or just a .pem file content.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"verified": {
			Description: "Whether the certificate is created in verified state.",
			Type:        pluginsdk.TypeBool,
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
	}

	if !features.SixPointOh() {
		s["is_verified"] = &pluginsdk.Schema{ // azignore:AZS023 - the old name, kept until 6.0 removes it
			Description:   "Whether the certificate is created in verified state.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"verified"},
			Deprecated:    "`is_verified` has been deprecated in favour of `verified` and will be removed in v6.0 of the AzureRM Provider",
		}

		s["verified"] = &pluginsdk.Schema{
			Description:   "Whether the certificate is created in verified state.",
			Type:          pluginsdk.TypeBool,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"is_verified"},
		}
	}

	return s
}
