package healthcare

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	fhirService "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/fhirservices"
	iotConnector "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/iotconnectors"
	workspace "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/sdk/2021-06-01-preview/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisIotConnector() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisIotConnectorCreateUpdate,
		Read:   resoruceHealthcareApisIotConnectorRead,
		Update: resourceHealthcareApisIotConnectorCreateUpdate,
		Delete: resourceHealthcareApisIotConnectorDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := iotConnector.ParseIotConnectorID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				//todo: check the name validation function
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: workspace.ValidateWorkspaceID,
			},

			"location": azure.SchemaLocation(),

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(fhirService.ManagedServiceIdentityTypeSystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"eventhub_namespace_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubNamespaceName(),
			},

			"eventhub_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"eventhub_consumer_group_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhubValidate.ValidateEventHubConsumerName(),
			},

			"destination_fhir_service_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: fhirService.ValidateFhirServiceID,
			},

			//todo use the fhir service name validation
			"destination_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"destination_identity_resolution_type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(iotConnector.IotIdentityResolutionTypeCreate),
					string(iotConnector.IotIdentityResolutionTypeLookup),
				}, false),
			},

			"destination_fhir_mapping": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"device_mapping": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisIotConnectorCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare Iot connector creation.")

	workspace, err := workspace.ParseWorkspaceIDInsensitively(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare api workspace error: %+v", err)
	}
	iotConnectorId := iotConnector.NewIotConnectorID(workspace.SubscriptionId, workspace.ResourceGroupName, workspace.WorkspaceName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, iotConnectorId)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presense of existing %s: %+v", iotConnectorId, err)
			}
		}
		if existing.Model != nil {
			return tf.ImportAsExistsError("azurerm_healthcareapis_iot_connector", iotConnectorId.ID())
		}
	}

	namespaceName := d.Get("eventhub_namespace_name").(string) + ".servicebus.windows.net"

	parameters := iotConnector.IotConnector{
		Name:     utils.String(iotConnectorId.IotConnectorName),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Identity: expandIotConnectorIdentity(d.Get("identity").([]interface{})),
		Properties: &iotConnector.IotConnectorProperties{
			IngestionEndpointConfiguration: &iotConnector.IotEventHubIngestionEndpointConfiguration{
				EventHubName:                    utils.String(d.Get("eventhub_name").(string)),
				ConsumerGroup:                   utils.String(d.Get("eventhub_consumer_group_name").(string)),
				FullyQualifiedEventHubNamespace: &namespaceName,
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	deviceContentMap := iotConnector.IotMappingProperties{}
	deviceMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("device_mapping").(string))
	if err := json.Unmarshal([]byte(deviceMappingJson), &deviceContentMap); err != nil {
		return err //todo fmt
	}
	parameters.Properties.DeviceMapping = &deviceContentMap

	if err := client.CreateOrUpdateThenPoll(ctx, iotConnectorId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", iotConnectorId, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 12,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"Creating", "Updating"},
		Target:                    []string{"Succeeded"},
		Refresh:                   iotConnectorStateRefreshFunc(ctx, client, iotConnectorId),
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Iot Connetcor %s to settle down: %+v", iotConnectorId, err)
	}

	destinationName := d.Get("destination_name").(string)
	fhirServiceId, err := fhirService.ParseFhirServiceID(d.Get("destination_fhir_service_id").(string))
	if err != nil {
		return fmt.Errorf("parsing fhir destination id err: %+v", err)
	}

	iotFhirServiceId := iotConnector.NewFhirDestinationID(iotConnectorId.SubscriptionId, iotConnectorId.ResourceGroupName, iotConnectorId.WorkspaceName, iotConnectorId.IotConnectorName, destinationName)
	iotFhirServiceParameters := iotConnector.IotFhirDestination{
		Name:     &destinationName,
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Properties: iotConnector.IotFhirDestinationProperties{
			FhirServiceResourceId:          fhirServiceId.ID(),
			ResourceIdentityResolutionType: iotConnector.IotIdentityResolutionType(d.Get("destination_identity_resolution_type").(string)),
		},
	}

	fhirMap := iotConnector.IotMappingProperties{}
	fhirMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("destination_fhir_mapping").(string))
	if err := json.Unmarshal([]byte(fhirMappingJson), &fhirMap); err != nil {
		return err
	}
	iotFhirServiceParameters.Properties.FhirMapping = fhirMap

	if err := client.IotConnectorFhirDestinationCreateOrUpdateThenPoll(ctx, iotFhirServiceId, iotFhirServiceParameters); err != nil {
		return fmt.Errorf("updating fhir service %s for the iot connector err: %+v", iotFhirServiceId, err)
	}

	d.SetId(iotConnectorId.ID())

	return resourceHealthcareServiceRead(d, meta)
}

func resoruceHealthcareApisIotConnectorRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotConnector.ParseIotConnectorID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] Healthcare Apis Iot Connector %s was not found", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IotConnectorName)
	workspaceId := workspace.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName)
	d.Set("workspace_id", workspaceId.ID())

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		d.Set("identity", flattenIotConnectorIdentity(model.Identity))

		if props := model.Properties; props != nil {
			if props.IngestionEndpointConfiguration.EventHubName != nil {
				d.Set("eventhub_name", props.IngestionEndpointConfiguration.EventHubName)
			}

			if props.IngestionEndpointConfiguration.ConsumerGroup != nil {
				d.Set("eventhub_consumer_group_name", props.IngestionEndpointConfiguration.ConsumerGroup)
			}

			if props.DeviceMapping != nil {
				deviceMapData, err := json.Marshal(props.DeviceMapping)
				if err != nil {
					return err
				}

				var m map[string]*json.RawMessage
				if err = json.Unmarshal(deviceMapData, &m); err != nil {
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
				d.Set("device_mapping", mapContent)
			}

			if props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace != nil {
				d.Set("eventhub_namespace_name", strings.TrimSuffix(*props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace, ".servicebus.windows.net"))
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	iotFhirService, err := client.FhirDestinationsListByIotConnector(ctx, *id)
	if err != nil {
		log.Printf("retrieving Fhir Destination for iot connector %s error", id)
		d.Set("destination_fhir_service_id", "")
	}

	//todo: deal with next page
	//{
	//    "error": {
	//        "code": "BadRequest",
	//        "message": "Only 1 resource type 'fhirdestinations' can be provisioned in iotconnectors."
	//    }
	//}
	if props := iotFhirService.Model; props != nil && len(*props) > 0 {
		for _, item := range *props {
			if item.Name != nil {
				d.Set("destination_name", item.Name)
			}
			if item.Properties.FhirServiceResourceId != "" {
				d.Set("destination_fhir_service_id", item.Properties.FhirServiceResourceId)
			}
			if item.Properties.FhirMapping.Content != nil {
				fhirMapData, err := json.Marshal(item.Properties.FhirMapping)
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
				d.Set("destination_fhir_mapping", mapContent)
			}
		}
	}

	return nil
}
func resourceHealthcareApisIotConnectorDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceIotConnectorClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := iotConnector.ParseIotConnectorID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(future.HttpResponse) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}
	return waitForHealthcareApisIotConnetorToBeDeleted(ctx, client, *id)
}

func waitForHealthcareApisIotConnetorToBeDeleted(ctx context.Context, client *iotConnector.IotConnectorsClient, id iotConnector.IotConnectorId) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context has no deadline")
	}

	log.Printf("[DEBUG] Waiting for %s to be deleted..", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{"200"},
		Target:  []string{"404"},
		Refresh: healthcareApiIotConnectorStateCodeRefreshFunc(ctx, client, id),
		Timeout: time.Until(deadline),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func healthcareApiIotConnectorStateCodeRefreshFunc(ctx context.Context, client *iotConnector.IotConnectorsClient, id iotConnector.IotConnectorId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id)
		if res.HttpResponse != nil {
			log.Printf("Retrieving %s returned Status %d", id, res.HttpResponse.StatusCode)
		}

		if err != nil {
			if response.WasNotFound(res.HttpResponse) {
				return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
			}
			return nil, "", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, strconv.Itoa(res.HttpResponse.StatusCode), nil
	}
}
func expandIotConnectorIdentity(input []interface{}) *iotConnector.ServiceManagedIdentityIdentity {
	//todo: is there any other way to set the address?
	typeNone := iotConnector.ManagedServiceIdentityTypeNone
	if len(input) == 0 {
		return &iotConnector.ServiceManagedIdentityIdentity{
			Type: &typeNone,
		}
	}

	identity := input[0].(map[string]interface{})
	inputType := iotConnector.ManagedServiceIdentityType(identity["type"].(string))
	return &iotConnector.ServiceManagedIdentityIdentity{
		Type: &inputType,
	}
}

func flattenIotConnectorIdentity(identity *iotConnector.ServiceManagedIdentityIdentity) []interface{} {
	if identity == nil || *identity.Type == iotConnector.ManagedServiceIdentityTypeNone {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(*identity.Type)

	//todo:check if there is any tenantID and principalID will be added in the stable api
	//if identity.PrincipalID != nil {
	//	result["principal_id"] = identity.PrincipalID.String()
	//}
	//
	//if identity.TenantID != nil {
	//	result["tenant_id"] = identity.TenantID.String()
	//}
	return []interface{}{result}
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}

func iotConnectorStateRefreshFunc(ctx context.Context, client *iotConnector.IotConnectorsClient, iotConnectorId iotConnector.IotConnectorId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, iotConnectorId)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil, "", fmt.Errorf("unable to retrieve iot connector %q: %+v", iotConnectorId, err)
			}
		}

		return resp, string(*resp.Model.Properties.ProvisioningState), nil
	}
}
