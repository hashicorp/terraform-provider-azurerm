// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_dicom_service

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2024-03-31/dicomservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/resources/healthcare_dicom_service/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisDicomServiceCreate,
		Read:   resourceHealthcareApisDicomServiceRead,
		Update: resourceHealthcareApisDicomServiceUpdate,
		Delete: resourceHealthcareApisDicomServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareDicomServiceV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := dicomservices.ParseDicomServiceID(id)
			return err
		}),

		Schema: healthcareDicomServiceSchema(),
	}
}
