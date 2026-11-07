// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package container_app_environment_managed_certificate

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/containerapps/2025-07-01/managedenvironments"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.CertificateName,
			Description:  "The name of the Container Apps Managed Certificate.",
		},

		"container_app_environment_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managedenvironments.ValidateManagedEnvironmentID,
			Description:  "The Container App Managed Environment ID to configure this Managed Certificate on.",
		},

		"subject_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			Description:  "The Subject Name of the Certificate. Must be a valid domain name.",
		},

		"domain_control_validation": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
			Default:  string(managedenvironments.ManagedCertificateDomainControlValidationHTTP),
			ValidateFunc: validation.StringInSlice(
				[]string{
					string(managedenvironments.ManagedCertificateDomainControlValidationCNAME),
					string(managedenvironments.ManagedCertificateDomainControlValidationHTTP),
				},
				false,
			),
			Description: "The domain control validation type for the managed certificate. Possible values are `CNAME` and `HTTP`. Defaults to `HTTP`.",
		},

		"tags": commonschema.Tags(),
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"validation_token": {
			Type:        pluginsdk.TypeString,
			Computed:    true,
			Description: "The validation token for the managed certificate.",
		},
	}
}
