package healthcare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
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
			_, err := parse.MedTechServiceFhirDestinationID(id)
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
				ValidateFunc: validate.MedTechServiceID,
			},

			"location": commonschema.Location(),

			"destination_fhir_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.FhirServiceID,
			},

			"destination_identity_resolution_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(healthcareapis.IotIdentityResolutionTypeCreate),
					string(healthcareapis.IotIdentityResolutionTypeLookup),
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
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceFhirDestinationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Med Tech Service Fhir Destination creation.")

	medTechService, err := parse.MedTechServiceID(d.Get("medtech_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Med Tech Service error: %+v", err)
	}
	id := parse.NewMedTechServiceFhirDestinationID(medTechService.SubscriptionId, medTechService.ResourceGroup, medTechService.WorkspaceName, medTechService.IotConnectorName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_healthcare_medtech_service_fhir_destination", id.ID())
		}
	}

	fhirServiceId, err := parse.FhirServiceID(d.Get("destination_fhir_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing fhir destination id err: %+v", err)
	}

	iotFhirServiceParameters := healthcareapis.IotFhirDestination{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		IotFhirDestinationProperties: &healthcareapis.IotFhirDestinationProperties{
			FhirServiceResourceID:          utils.String(fhirServiceId.ID()),
			ResourceIdentityResolutionType: healthcareapis.IotIdentityResolutionType(d.Get("destination_identity_resolution_type").(string)),
		},
	}

	fhirMap := healthcareapis.IotMappingProperties{}
	fhirMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("destination_fhir_mapping_json").(string))
	if err := json.Unmarshal([]byte(fhirMappingJson), &fhirMap); err != nil {
		return err
	}
	iotFhirServiceParameters.IotFhirDestinationProperties.FhirMapping = &fhirMap

	medTechServiceDesFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName, iotFhirServiceParameters)
	if err != nil {
		return fmt.Errorf("updating fhir service %s for the Med Tech Service err: %+v", id, err)
	}
	if err = medTechServiceDesFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisMedTechServiceFhirDestinationRead(d, meta)
}

func resourceHealthcareApisMedTechServiceFhirDestinationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceFhirDestinationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MedTechServiceFhirDestinationID(d.Id())
	if err != nil {
		return err
	}

	medTechServiceId := parse.NewMedTechServiceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName)
	d.Set("medtech_service_id", medTechServiceId.ID())

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Healthcare Apis Med Tech Service Fhir Destination %s was not found", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FhirDestinationName)

	if resp.Location != nil {
		d.Set("location", location.NormalizeNilable(resp.Location))
	}

	if props := resp.IotFhirDestinationProperties; props != nil {
		if props.FhirServiceResourceID != nil {
			d.Set("destination_fhir_service_id", props.FhirServiceResourceID)
		}

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
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceFhirDestinationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	medTechService, err := parse.MedTechServiceID(d.Get("medtech_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing Med Tech Service error: %+v", err)
	}
	id := parse.NewMedTechServiceFhirDestinationID(medTechService.SubscriptionId, medTechService.ResourceGroup, medTechService.WorkspaceName, medTechService.IotConnectorName, d.Get("name").(string))

	fhirServiceId, err := parse.FhirServiceID(d.Get("destination_fhir_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing fhir destination id err: %+v", err)
	}

	medTechFhirServiceParameters := healthcareapis.IotFhirDestination{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		IotFhirDestinationProperties: &healthcareapis.IotFhirDestinationProperties{
			FhirServiceResourceID:          utils.String(fhirServiceId.ID()),
			ResourceIdentityResolutionType: healthcareapis.IotIdentityResolutionType(d.Get("destination_identity_resolution_type").(string)),
		},
	}

	fhirMap := healthcareapis.IotMappingProperties{}
	fhirMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("destination_fhir_mapping_json").(string))
	if err := json.Unmarshal([]byte(fhirMappingJson), &fhirMap); err != nil {
		return err
	}
	medTechFhirServiceParameters.IotFhirDestinationProperties.FhirMapping = &fhirMap

	medTechServiceDesFuture, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName, medTechFhirServiceParameters)
	if err != nil {
		return fmt.Errorf("updating fhir service %s for the Med Tech Service err: %+v", id, err)
	}
	if err = medTechServiceDesFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceHealthcareApisMedTechServiceFhirDestinationRead(d, meta)
}

func resourceHealthcareApisMedTechServiceFhirDestinationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceFhirDestinationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MedTechServiceFhirDestinationID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
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

func healthcareApiMedTechServiceFhirDestinationStateCodeRefreshFunc(ctx context.Context, client *healthcareapis.IotConnectorFhirDestinationClient, id parse.MedTechServiceFhirDestinationId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName, id.FhirDestinationName)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, "Pending", nil
	}
}
