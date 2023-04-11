package healthcare

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/healthcareapis/mgmt/2021-11-01/healthcareapis" // nolint: staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	eventhubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceHealthcareApisMedTechService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceHealthcareApisMedTechServiceCreate,
		Read:   resourceHealthcareApisMedTechServiceRead,
		Update: resourceHealthcareApisMedTechServiceUpdate,
		Delete: resourceHealthcareApisMedTechServiceDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(90 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(90 * time.Minute),
		},

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.HealthCareIoTConnectorV0ToV1{},
		}),

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.MedTechServiceID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.MedTechServiceName(),
			},

			"workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"location": commonschema.Location(),

			"identity": commonschema.SystemAssignedIdentityOptional(),

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

			"device_mapping_json": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceHealthcareApisMedTechServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	log.Printf("[INFO] preparing arguments for AzureRM Healthcare MedTech Service creation.")

	workspace, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare api workspace error: %+v", err)
	}
	medTechServiceId := parse.NewMedTechServiceID(workspace.SubscriptionId, workspace.ResourceGroup, workspace.Name, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, medTechServiceId.ResourceGroup, medTechServiceId.WorkspaceName, medTechServiceId.IotConnectorName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", medTechServiceId, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_healthcare_medtech_service", medTechServiceId.ID())
		}
	}

	namespaceName := d.Get("eventhub_namespace_name").(string) + ".servicebus.windows.net"
	identity, err := expandMedTechServiceIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := healthcareapis.IotConnector{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Identity: identity,
		IotConnectorProperties: &healthcareapis.IotConnectorProperties{
			IngestionEndpointConfiguration: &healthcareapis.IotEventHubIngestionEndpointConfiguration{
				EventHubName:                    utils.String(d.Get("eventhub_name").(string)),
				ConsumerGroup:                   utils.String(d.Get("eventhub_consumer_group_name").(string)),
				FullyQualifiedEventHubNamespace: &namespaceName,
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	deviceContentMap := healthcareapis.IotMappingProperties{}
	deviceMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("device_mapping_json").(string))
	if err := json.Unmarshal([]byte(deviceMappingJson), &deviceContentMap); err != nil {
		return err
	}
	parameters.IotConnectorProperties.DeviceMapping = &deviceContentMap

	future, err := client.CreateOrUpdate(ctx, medTechServiceId.ResourceGroup, medTechServiceId.WorkspaceName, medTechServiceId.IotConnectorName, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", medTechServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", medTechServiceId, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 12,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"Creating", "Updating"},
		Target:                    []string{"Succeeded"},
		Refresh:                   medTechServiceCreateStateRefreshFunc(ctx, client, medTechServiceId),
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for MedTech Service %s to settle down: %+v", medTechServiceId, err)
	}

	d.SetId(medTechServiceId.ID())

	return resourceHealthcareApisMedTechServiceRead(d, meta)
}

func resourceHealthcareApisMedTechServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceClient
	domainSuffix, ok := meta.(*clients.Client).Account.Environment.ServiceBus.DomainSuffix()
	if !ok {
		return fmt.Errorf("unable to retrieve the Domain Suffix for ServiceBus, this is not configured for this Cloud Environment")
	}
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MedTechServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] Healthcare Apis MedTech Service %s was not found", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IotConnectorName)
	d.Set("workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("identity", flattenMedTechServiceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.IotConnectorProperties; props != nil {
		eventHubConsumerGroupName := ""
		eventHubName := ""
		eventHubNamespaceName := ""
		if config := props.IngestionEndpointConfiguration; config != nil {
			if config.ConsumerGroup != nil {
				eventHubConsumerGroupName = *config.ConsumerGroup
			}

			if config.EventHubName != nil {
				eventHubName = *config.EventHubName
			}

			if props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace != nil {
				suffixToTrim := "." + *domainSuffix
				eventHubNamespaceName = strings.TrimSuffix(*props.IngestionEndpointConfiguration.FullyQualifiedEventHubNamespace, suffixToTrim)
			}
		}
		d.Set("eventhub_consumer_group_name", eventHubConsumerGroupName)
		d.Set("eventhub_name", eventHubName)
		d.Set("eventhub_namespace_name", eventHubNamespaceName)

		mapContent := ""
		if props.DeviceMapping != nil {
			deviceMapData, err := json.Marshal(props.DeviceMapping)
			if err != nil {
				return err
			}

			var m map[string]*json.RawMessage
			if err = json.Unmarshal(deviceMapData, &m); err != nil {
				return err
			}
			if v, ok := m["content"]; ok {
				contents, err := json.Marshal(v)
				if err != nil {
					return err
				}
				mapContent = string(contents)
			}
		}
		d.Set("device_mapping_json", mapContent)
	}

	if err := tags.FlattenAndSet(d, resp.Tags); err != nil {
		return err
	}
	return nil
}

func resourceHealthcareApisMedTechServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	workspace, err := parse.WorkspaceID(d.Get("workspace_id").(string))
	if err != nil {
		return fmt.Errorf("parsing healthcare api workspace error: %+v", err)
	}
	medTechServiceId := parse.NewMedTechServiceID(workspace.SubscriptionId, workspace.ResourceGroup, workspace.Name, d.Get("name").(string))

	namespaceName := d.Get("eventhub_namespace_name").(string) + ".servicebus.windows.net"
	identity, err := expandMedTechServiceIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	parameters := healthcareapis.IotConnector{
		Location: utils.String(azure.NormalizeLocation(d.Get("location").(string))),
		Identity: identity,
		IotConnectorProperties: &healthcareapis.IotConnectorProperties{
			IngestionEndpointConfiguration: &healthcareapis.IotEventHubIngestionEndpointConfiguration{
				EventHubName:                    utils.String(d.Get("eventhub_name").(string)),
				ConsumerGroup:                   utils.String(d.Get("eventhub_consumer_group_name").(string)),
				FullyQualifiedEventHubNamespace: &namespaceName,
			},
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	deviceContentMap := healthcareapis.IotMappingProperties{}
	deviceMappingJson := fmt.Sprintf(`{ "content": %s }`, d.Get("device_mapping_json").(string))
	if err := json.Unmarshal([]byte(deviceMappingJson), &deviceContentMap); err != nil {
		return err
	}
	parameters.IotConnectorProperties.DeviceMapping = &deviceContentMap

	future, err := client.CreateOrUpdate(ctx, medTechServiceId.ResourceGroup, medTechServiceId.WorkspaceName, medTechServiceId.IotConnectorName, parameters)
	if err != nil {
		return fmt.Errorf("updating %s: %+v", medTechServiceId, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for update of %s: %+v", medTechServiceId, err)
	}

	stateConf := &pluginsdk.StateChangeConf{
		ContinuousTargetOccurence: 12,
		Delay:                     60 * time.Second,
		MinTimeout:                10 * time.Second,
		Pending:                   []string{"Creating", "Updating"},
		Target:                    []string{"Succeeded"},
		Refresh:                   medTechServiceCreateStateRefreshFunc(ctx, client, medTechServiceId),
		Timeout:                   d.Timeout(pluginsdk.TimeoutUpdate),
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for MedTech Service %s to settle down: %+v", medTechServiceId, err)
	}

	d.SetId(medTechServiceId.ID())

	return resourceHealthcareApisMedTechServiceRead(d, meta)
}

func resourceHealthcareApisMedTechServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).HealthCare.HealthcareWorkspaceMedTechServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MedTechServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.IotConnectorName, id.WorkspaceName)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	// NOTE: this can be removed when using `hashicorp/go-azure-sdk`'s base layer
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"Pending"},
		Target:                    []string{"Deleted"},
		Refresh:                   medTechServiceStateStatusCodeRefreshFunc(ctx, client, *id),
		Timeout:                   d.Timeout(pluginsdk.TimeoutDelete),
		ContinuousTargetOccurence: 3,
		PollInterval:              10 * time.Second,
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}
	return nil
}

func medTechServiceStateStatusCodeRefreshFunc(ctx context.Context, client *healthcareapis.IotConnectorsClient, id parse.MedTechServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id.ResourceGroup, id.WorkspaceName, id.IotConnectorName)

		if err != nil {
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Deleted", nil
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", id, err)
		}

		return res, "Pending", nil
	}
}

func expandMedTechServiceIdentity(input []interface{}) (*healthcareapis.ServiceManagedIdentityIdentity, error) {
	expanded, err := identity.ExpandSystemAssigned(input)
	if err != nil {
		return nil, err
	}

	return &healthcareapis.ServiceManagedIdentityIdentity{
		Type: healthcareapis.ServiceManagedIdentityType(string(expanded.Type)),
	}, nil
}

func flattenMedTechServiceIdentity(input *healthcareapis.ServiceManagedIdentityIdentity) []interface{} {
	var transition *identity.SystemAssigned

	if input != nil {
		transition = &identity.SystemAssigned{
			Type: identity.Type(string(input.Type)),
		}
		if input.PrincipalID != nil {
			principalID := *input.PrincipalID
			transition.PrincipalId = principalID.String()
		}
		if input.TenantID != nil {
			tenantID := *input.TenantID
			transition.TenantId = tenantID.String()
		}
	}

	return identity.FlattenSystemAssigned(transition)
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}

func medTechServiceCreateStateRefreshFunc(ctx context.Context, client *healthcareapis.IotConnectorsClient, medTechServiceId parse.MedTechServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, medTechServiceId.ResourceGroup, medTechServiceId.WorkspaceName, medTechServiceId.IotConnectorName)
		if err != nil {
			if utils.ResponseWasNotFound(resp.Response) {
				return nil, "", fmt.Errorf("unable to retrieve MedTech Service %q: %+v", medTechServiceId, err)
			}
			return nil, "Error", fmt.Errorf("polling for the status of %s: %+v", medTechServiceId, err)
		}

		return resp, string(resp.ProvisioningState), nil
	}
}
