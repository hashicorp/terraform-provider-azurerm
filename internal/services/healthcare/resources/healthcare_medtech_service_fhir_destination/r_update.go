// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service_fhir_destination

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisMedTechServiceFhirDestinationUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	medTechService, err := iotconnectors.ParseIotConnectorID(d.Get("medtech_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Med Tech Service error: %+v", err)
	}
	id := iotconnectors.NewFhirDestinationID(medTechService.SubscriptionId, medTechService.ResourceGroupName, medTechService.WorkspaceName, medTechService.IotConnectorName, d.Get("name").(string))

	medTechFhirServiceParameters := iotconnectors.IotFhirDestination{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: iotconnectors.IotFhirDestinationProperties{
			FhirServiceResourceId:          d.Get("destination_fhir_service_id").(string),
			ResourceIdentityResolutionType: iotconnectors.IotIdentityResolutionType(d.Get("destination_identity_resolution_type").(string)),
		},
	}

	fhirMap := iotconnectors.IotMappingProperties{}
	fhirMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("destination_fhir_mapping_json").(string))
	if err := json.Unmarshal([]byte(fhirMappingJson), &fhirMap); err != nil {
		return err
	}
	medTechFhirServiceParameters.Properties.FhirMapping = fhirMap

	if err = client.IotConnectorFhirDestinationCreateOrUpdateThenPoll(ctx, id, medTechFhirServiceParameters); err != nil {
		return fmt.Errorf("updating fhir service %s for the Med Tech Service err: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisMedTechServiceFhirDestinationRead(d, meta)
}
