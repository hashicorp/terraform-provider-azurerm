package healthcare

import (
	"context"
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
	consumerGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/sdk/2017-04-01/consumergroups"
	eventhub "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/sdk/2017-04-01/eventhubs"
	namespaceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/sdk/2021-01-01-preview/namespaces"
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
				ValidateFunc: eventhubValidate.ValidateEventHubName(),
			},

			"eventhub_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: eventhub.ValidateEventhubID,
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

			//todo check if there is any enum of "CollectionContent" in the stable api
			"device_mapping": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"template_type": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"CollectionContent"}, false),
						},

						"template": {
							Type:     pluginsdk.TypeList,
							Optional: true,
						},
					},
				},
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

	eventhubId, err := eventhub.ParseEventhubID(d.Get("eventhub_id").(string))
	if err != nil {
		return fmt.Errorf("parsing eventhub id error: %+v", err)
	}

	consumerGroupId, err := consumerGroup.ParseConsumerGroupID(d.Get("eventhub_consumer_group_name").(string))
	if err != nil {
		return fmt.Errorf("parsing eventhub consumer group id error: %+v", err)
	}

	parameters := iotConnector.IotConnector{
		Name:     utils.String(iotConnectorId.IotConnectorName),
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Identity: expandIotConnectorIdentity(d.Get("identity").([]interface{})),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &iotConnector.IotConnectorProperties{
			IngestionEndpointConfiguration: &iotConnector.IotEventHubIngestionEndpointConfiguration{
				EventHubName:                    &eventhubId.EventHubName,
				ConsumerGroup:                   &consumerGroupId.ConsumerGroupName,
				FullyQualifiedEventHubNamespace: utils.String(d.Get("eventhub_namespace_name").(string)),
			},
			DeviceMapping: expandIotConnectorDeviceMapping(d.Get("device_mapping").([]interface{})),
		},
	}

	if err := client.CreateOrUpdateThenPoll(ctx, iotConnectorId, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", iotConnectorId, err)
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
			if props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace != nil {
				//todo: to deal with the port number, validate if this works
				namespaceWithSuffix := strings.TrimPrefix(*props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace, "https://")
				namespace:= namespaceWithSuffix[:strings.IndexAny(namespaceWithSuffix, ".")]

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

//
func expandIotConnectorDeviceMapping(input []interface{}) *iotConnector.IotMappingProperties {
	rawData := input[0].(map[string]interface{})

	content := make(map[string]interface{})
	content["templateType"] = rawData["template_type"]
	content["template"] = rawData["template"]

	return &iotConnector.IotMappingProperties{
		Content: &[]interface{}{content}[0],
	}
}

//func flattenIotConnectorDeviceMapping(input *iotConnector.IotMappingProperties) []interface{} {
//	if input == nil {
//		return []interface{}{}
//	}
//
//	result := make(map[string]interface{})
//	if data := input.Content; data != nil {
//
//	}
//
//}
