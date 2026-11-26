// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package healthcare_medtech_service_fhir_destination

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceHealthcareApisMedTechServiceFhirDestinationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotconnectors.ParseFhirDestinationID(d.Id())
	if err != nil {
		return err
	}

	d.Set("medtech_service_id", iotconnectors.NewIotConnectorID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.IotConnectorName).ID())

	resp, err := client.IotConnectorFhirDestinationGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] Healthcare Apis Med Tech Service Fhir Destination %s was not found", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FhirDestinationName)

	if m := resp.Model; m != nil {
		d.Set("location", location.NormalizeNilable(m.Location))

		props := m.Properties
		d.Set("destination_fhir_service_id", props.FhirServiceResourceId)

		if props.FhirMapping.Content != nil {
			fhirMapData, err := json.Marshal(props.FhirMapping)
			if err != nil {
				return err
			}

			var m map[string]*json.RawMessage
			if err = json.Unmarshal(fhirMapData, &m); err != nil {
				return err
			}
			mapContent := ""
			if v, ok := m["content"]; ok {
				contents, err := json.Marshal(v)
				if err != nil {
					return err
				}
				mapContent = string(contents)
			}
			d.Set("destination_fhir_mapping_json", mapContent)
		}
		d.Set("destination_identity_resolution_type", props.ResourceIdentityResolutionType)
	}

	return nil
}
