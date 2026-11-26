// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service_fhir_destination

import (
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func healthcareMedtechServiceFhirDestinationSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Healthcare Med Tech Service Fhir Destination.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.MedTechServiceName(),
		},

		"medtech_service_id": {
			Description:  "The name of the Healthcare Med Tech Service where the Healthcare Med Tech Service Fhir Destination should exist.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: iotconnectors.ValidateIotConnectorID,
		},

		"location": commonschema.Location(),

		"destination_fhir_service_id": {
			Description:  "The destination fhir service ID of the Med Tech Service Fhir Destination.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: fhirservices.ValidateFhirServiceID,
		},

		"destination_identity_resolution_type": {
			Description:  "The destination identity resolution type where the Healthcare Med Tech Service Fhir Destination should be created.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice(iotconnectors.PossibleValuesForIotIdentityResolutionType(), false),
		},

		"destination_fhir_mapping_json": {
			Description:      "The destination Fhir mappings of the Med Tech Service Fhir Destination.",
			Type:             pluginsdk.TypeString,
			Required:         true,
			StateFunc:        pluginsdk.NormalizeJson,
			DiffSuppressFunc: helpers.SuppressJsonOrderingDifference,
		},
	}
}
