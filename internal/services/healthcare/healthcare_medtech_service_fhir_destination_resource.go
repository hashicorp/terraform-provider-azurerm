// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package healthcare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/fhirservices"
	"github.com/hashicorp/go-azure-sdk/resource-manager/healthcareapis/2022-12-01/iotconnectors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisMedTechServiceFhirDestination() *pluginsdk.Resource {
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

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MedTechServiceName(),
			},

			"medtech_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iotconnectors.ValidateIotConnectorID,
			},

			"location": commonschema.Location(),

			"destination_fhir_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: fhirservices.ValidateFhirServiceID,
			},

			"destination_identity_resolution_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(iotconnectors.IotIdentityResolutionTypeCreate),
					string(iotconnectors.IotIdentityResolutionTypeLookup),
				}, false),
			},

			"destination_fhir_mapping_json": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},
		},
	}
}

func resourceHealthcareApisMedTechServiceFhirDestinationCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Med Tech Service Fhir Destination creation.")

	medTechService, err := iotconnectors.ParseIotConnectorID(d.Get("medtech_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Med Tech Service error: %+v", err)
	}
	id := iotconnectors.NewFhirDestinationID(medTechService.SubscriptionId, medTechService.ResourceGroupName, medTechService.WorkspaceName, medTechService.IotConnectorName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.IotConnectorFhirDestinationGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_healthcare_medtech_service_fhir_destination", id.ID())
		}
	}

	iotFhirServiceParameters := iotconnectors.IotFhirDestination{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
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
	iotFhirServiceParameters.Properties.FhirMapping = fhirMap

	err = client.IotConnectorFhirDestinationCreateOrUpdateThenPoll(ctx, id, iotFhirServiceParameters)
	if err != nil {
		return fmt.Errorf("updating fhir service %s for the Med Tech Service err: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisMedTechServiceFhirDestinationRead(d, meta)
}

func resourceHealthcareApisMedTechServiceFhirDestinationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceHealthcareApisMedTechServiceFhirDestinationUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	medTechService, err := iotconnectors.ParseIotConnectorID(d.Get("medtech_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Med Tech Service error: %+v", err)
	}
	id := iotconnectors.NewFhirDestinationID(medTechService.SubscriptionId, medTechService.ResourceGroupName, medTechService.WorkspaceName, medTechService.IotConnectorName, d.Get("name").(string))

	medTechFhirServiceParameters := iotconnectors.IotFhirDestination{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
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

	err = client.IotConnectorFhirDestinationCreateOrUpdateThenPoll(ctx, id, medTechFhirServiceParameters)
	if err != nil {
		return fmt.Errorf("updating fhir service %s for the Med Tech Service err: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisMedTechServiceFhirDestinationRead(d, meta)
}

func resourceHealthcareApisMedTechServiceFhirDestinationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotconnectors.ParseFhirDestinationID(d.Id())
	if err != nil {
		return err
	}

	err = client.IotConnectorFhirDestinationDeleteThenPoll(ctx, *id)
	if err != nil {

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Pending"},
		Target:                    []string{"Deleted"},
		Refresh:                   healthcareApiMedTechServiceFhirDestinationStateCodeRefreshFunc(ctx, client, *id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 3,
		PollInterval:              10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func healthcareApiMedTechServiceFhirDestinationStateCodeRefreshFunc(ctx context.Context, client *iotconnectors.IotConnectorsClient, id iotconnectors.FhirDestinationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.IotConnectorFhirDestinationGet(ctx, id)

		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return resp, "Pending", nil
	}
}
