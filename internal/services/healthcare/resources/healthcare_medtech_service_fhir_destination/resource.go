// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service_fhir_destination

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_medtech_service_fhir_destination/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisMedTechServiceFhirDestinationCreate,
		Read:   resourceHealthcareApisMedTechServiceFhirDestinationRead,
		Update: resourceHealthcareApisMedTechServiceFhirDestinationUpdate,
		Delete: resourceHealthcareApisMedTechServiceFhirDestinationDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := iotconnectors.ParseFhirDestinationID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareMedTechServiceFhirDestinationV0ToV1{},
		}),

		Schema: healthcareMedtechServiceFhirDestinationSchema(),
	}
}
