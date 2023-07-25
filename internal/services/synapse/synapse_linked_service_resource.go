// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	artifacts "github.com/tombuildsstuff/kermit/sdk/synapse/2021-06-01-preview/synapse"
)

func resourceSynapseLinkedService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSynapseLinkedServiceCreateUpdate,
		Read:   resourceSynapseLinkedServiceRead,
		Update: resourceSynapseLinkedServiceCreateUpdate,
		Delete: resourceSynapseLinkedServiceDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.LinkedServiceID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.SynapseLinkedServiceV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"synapse_workspace_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.WorkspaceID,
			},

			"type": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(artifacts.TypeBasicLinkedServiceTypeAmazonMWS),
					string(artifacts.TypeBasicLinkedServiceTypeAmazonRdsForOracle),
					string(artifacts.TypeBasicLinkedServiceTypeAmazonRdsForSQLServer),
					string(artifacts.TypeBasicLinkedServiceTypeAmazonRedshift),
					string(artifacts.TypeBasicLinkedServiceTypeAmazonS3),
					string(artifacts.TypeBasicLinkedServiceTypeAzureBatch),
					string(artifacts.TypeBasicLinkedServiceTypeAzureBlobFS),
					string(artifacts.TypeBasicLinkedServiceTypeAzureBlobStorage),
					string(artifacts.TypeBasicLinkedServiceTypeAzureDataExplorer),
					string(artifacts.TypeBasicLinkedServiceTypeAzureDataLakeAnalytics),
					string(artifacts.TypeBasicLinkedServiceTypeAzureDataLakeStore),
					string(artifacts.TypeBasicLinkedServiceTypeAzureDatabricks),
					string(artifacts.TypeBasicLinkedServiceTypeAzureDatabricksDeltaLake),
					string(artifacts.TypeBasicLinkedServiceTypeAzureFileStorage),
					string(artifacts.TypeBasicLinkedServiceTypeAzureFunction),
					string(artifacts.TypeBasicLinkedServiceTypeAzureKeyVault),
					string(artifacts.TypeBasicLinkedServiceTypeAzureML),
					string(artifacts.TypeBasicLinkedServiceTypeAzureMLService),
					string(artifacts.TypeBasicLinkedServiceTypeAzureMariaDB),
					string(artifacts.TypeBasicLinkedServiceTypeAzureMySQL),
					string(artifacts.TypeBasicLinkedServiceTypeAzurePostgreSQL),
					string(artifacts.TypeBasicLinkedServiceTypeAzureSQLDW),
					string(artifacts.TypeBasicLinkedServiceTypeAzureSQLDatabase),
					string(artifacts.TypeBasicLinkedServiceTypeAzureSQLMI),
					string(artifacts.TypeBasicLinkedServiceTypeAzureSearch),
					string(artifacts.TypeBasicLinkedServiceTypeAzureStorage),
					string(artifacts.TypeBasicLinkedServiceTypeAzureTableStorage),
					string(artifacts.TypeBasicLinkedServiceTypeCassandra),
					string(artifacts.TypeBasicLinkedServiceTypeCommonDataServiceForApps),
					string(artifacts.TypeBasicLinkedServiceTypeConcur),
					string(artifacts.TypeBasicLinkedServiceTypeCosmosDb),
					string(artifacts.TypeBasicLinkedServiceTypeCosmosDbMongoDbAPI),
					string(artifacts.TypeBasicLinkedServiceTypeCouchbase),
					string(artifacts.TypeBasicLinkedServiceTypeCustomDataSource),
					string(artifacts.TypeBasicLinkedServiceTypeDb2),
					string(artifacts.TypeBasicLinkedServiceTypeDrill),
					string(artifacts.TypeBasicLinkedServiceTypeDynamics),
					string(artifacts.TypeBasicLinkedServiceTypeDynamicsAX),
					string(artifacts.TypeBasicLinkedServiceTypeDynamicsCrm),
					string(artifacts.TypeBasicLinkedServiceTypeEloqua),
					string(artifacts.TypeBasicLinkedServiceTypeFileServer),
					string(artifacts.TypeBasicLinkedServiceTypeFtpServer),
					string(artifacts.TypeBasicLinkedServiceTypeGoogleAdWords),
					string(artifacts.TypeBasicLinkedServiceTypeGoogleBigQuery),
					string(artifacts.TypeBasicLinkedServiceTypeGoogleCloudStorage),
					string(artifacts.TypeBasicLinkedServiceTypeGreenplum),
					string(artifacts.TypeBasicLinkedServiceTypeHBase),
					string(artifacts.TypeBasicLinkedServiceTypeHDInsight),
					string(artifacts.TypeBasicLinkedServiceTypeHDInsightOnDemand),
					string(artifacts.TypeBasicLinkedServiceTypeHTTPServer),
					string(artifacts.TypeBasicLinkedServiceTypeHdfs),
					string(artifacts.TypeBasicLinkedServiceTypeHive),
					string(artifacts.TypeBasicLinkedServiceTypeHubspot),
					string(artifacts.TypeBasicLinkedServiceTypeImpala),
					string(artifacts.TypeBasicLinkedServiceTypeInformix),
					string(artifacts.TypeBasicLinkedServiceTypeJira),
					string(artifacts.TypeBasicLinkedServiceTypeLinkedService),
					string(artifacts.TypeBasicLinkedServiceTypeMagento),
					string(artifacts.TypeBasicLinkedServiceTypeMariaDB),
					string(artifacts.TypeBasicLinkedServiceTypeMarketo),
					string(artifacts.TypeBasicLinkedServiceTypeMicrosoftAccess),
					string(artifacts.TypeBasicLinkedServiceTypeMongoDb),
					string(artifacts.TypeBasicLinkedServiceTypeMongoDbAtlas),
					string(artifacts.TypeBasicLinkedServiceTypeMongoDbV2),
					string(artifacts.TypeBasicLinkedServiceTypeMySQL),
					string(artifacts.TypeBasicLinkedServiceTypeNetezza),
					string(artifacts.TypeBasicLinkedServiceTypeOData),
					string(artifacts.TypeBasicLinkedServiceTypeOdbc),
					string(artifacts.TypeBasicLinkedServiceTypeOffice365),
					string(artifacts.TypeBasicLinkedServiceTypeOracle),
					string(artifacts.TypeBasicLinkedServiceTypeOracleServiceCloud),
					string(artifacts.TypeBasicLinkedServiceTypePaypal),
					string(artifacts.TypeBasicLinkedServiceTypePhoenix),
					string(artifacts.TypeBasicLinkedServiceTypePostgreSQL),
					string(artifacts.TypeBasicLinkedServiceTypePresto),
					string(artifacts.TypeBasicLinkedServiceTypeQuickBooks),
					string(artifacts.TypeBasicLinkedServiceTypeResponsys),
					string(artifacts.TypeBasicLinkedServiceTypeRestService),
					string(artifacts.TypeBasicLinkedServiceTypeSQLServer),
					string(artifacts.TypeBasicLinkedServiceTypeSalesforce),
					string(artifacts.TypeBasicLinkedServiceTypeSalesforceMarketingCloud),
					string(artifacts.TypeBasicLinkedServiceTypeSalesforceServiceCloud),
					string(artifacts.TypeBasicLinkedServiceTypeSapBW),
					string(artifacts.TypeBasicLinkedServiceTypeSapCloudForCustomer),
					string(artifacts.TypeBasicLinkedServiceTypeSapEcc),
					string(artifacts.TypeBasicLinkedServiceTypeSapHana),
					string(artifacts.TypeBasicLinkedServiceTypeSapOpenHub),
					string(artifacts.TypeBasicLinkedServiceTypeSapTable),
					string(artifacts.TypeBasicLinkedServiceTypeServiceNow),
					string(artifacts.TypeBasicLinkedServiceTypeSftp),
					string(artifacts.TypeBasicLinkedServiceTypeSharePointOnlineList),
					string(artifacts.TypeBasicLinkedServiceTypeShopify),
					string(artifacts.TypeBasicLinkedServiceTypeSnowflake),
					string(artifacts.TypeBasicLinkedServiceTypeSpark),
					string(artifacts.TypeBasicLinkedServiceTypeSquare),
					string(artifacts.TypeBasicLinkedServiceTypeSybase),
					string(artifacts.TypeBasicLinkedServiceTypeTeradata),
					string(artifacts.TypeBasicLinkedServiceTypeVertica),
					string(artifacts.TypeBasicLinkedServiceTypeWeb),
					string(artifacts.TypeBasicLinkedServiceTypeXero),
					string(artifacts.TypeBasicLinkedServiceTypeZoho),
				}, false),
			},

			"type_properties_json": {
				Type:             pluginsdk.TypeString,
				Required:         true,
				StateFunc:        utils.NormalizeJson,
				DiffSuppressFunc: suppressJsonOrderingDifference,
			},

			"description": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"integration_runtime": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},

						"parameters": {
							Type:     pluginsdk.TypeMap,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
					},
				},
			},

			"parameters": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"annotations": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"additional_properties": {
				Type:     pluginsdk.TypeMap,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},
		},
	}
}

func resourceSynapseLinkedServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	workspaceId, err := parse.WorkspaceID(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.LinkedServiceClient(workspaceId.Name, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.GetLinkedService(ctx, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_synapse_linked_service", id.ID())
		}
	}

	props := map[string]interface{}{
		"type":       d.Get("type").(string),
		"connectVia": expandSynapseLinkedServiceIntegrationRuntimeV2(d.Get("integration_runtime").([]interface{})),
	}

	jsonDataStr := fmt.Sprintf(`{ "typeProperties": %s }`, d.Get("type_properties_json").(string))
	if err = json.Unmarshal([]byte(jsonDataStr), &props); err != nil {
		return err
	}

	if v, ok := d.GetOk("description"); ok {
		props["description"] = v.(string)
	}

	if v, ok := d.GetOk("parameters"); ok {
		props["parameters"] = expandSynapseParameters(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("annotations"); ok {
		props["annotations"] = v.([]interface{})
	}

	additionalProperties := d.Get("additional_properties").(map[string]interface{})
	for k, v := range additionalProperties {
		props[k] = v
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"properties": props,
	})
	if err != nil {
		return err
	}

	linkedService := &artifacts.LinkedServiceResource{}
	if err := linkedService.UnmarshalJSON(jsonData); err != nil {
		return err
	}

	future, err := client.CreateOrUpdateLinkedService(ctx, id.Name, *linkedService, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	// Sometimes this resource fails to create but Azure is returning a 200. We'll check if the last response failed or not before moving on
	// todo remove this once https://github.com/hashicorp/go-azure-sdk/pull/122 is merged
	if err = checkLinkedServiceResponse(future.Response()); err != nil {
		return err
	}

	d.SetId(id.ID())

	return resourceSynapseLinkedServiceRead(d, meta)
}

func resourceSynapseLinkedServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.LinkedServiceClient(id.WorkspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetLinkedService(ctx, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())

	byteArr, err := json.Marshal(resp.Properties)
	if err != nil {
		return err
	}

	var m map[string]*json.RawMessage
	if err = json.Unmarshal(byteArr, &m); err != nil {
		return err
	}

	description := ""
	if v, ok := m["description"]; ok && v != nil {
		if err := json.Unmarshal(*v, &description); err != nil {
			return err
		}
		delete(m, "description")
	}
	d.Set("description", description)

	t := ""
	if v, ok := m["type"]; ok && v != nil {
		if err := json.Unmarshal(*v, &t); err != nil {
			return err
		}
		delete(m, "type")
	}
	d.Set("type", t)

	annotations := make([]interface{}, 0)
	if v, ok := m["annotations"]; ok && v != nil {
		if err := json.Unmarshal(*v, &annotations); err != nil {
			return err
		}
		delete(m, "annotations")
	}
	d.Set("annotations", annotations)

	parameters := make(map[string]*artifacts.ParameterSpecification)
	if v, ok := m["parameters"]; ok && v != nil {
		if err := json.Unmarshal(*v, &parameters); err != nil {
			return err
		}
		delete(m, "parameters")
	}
	if err := d.Set("parameters", flattenSynapseParameters(parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	var integrationRuntime *artifacts.IntegrationRuntimeReference
	if v, ok := m["connectVia"]; ok && v != nil {
		integrationRuntime = &artifacts.IntegrationRuntimeReference{}
		if err := json.Unmarshal(*v, &integrationRuntime); err != nil {
			return err
		}
		delete(m, "connectVia")
	}
	if err := d.Set("integration_runtime", flattenSynapseLinkedServiceIntegrationRuntimeV2(integrationRuntime)); err != nil {
		return fmt.Errorf("setting `integration_runtime`: %+v", err)
	}

	delete(m, "typeProperties")

	// set "additional_properties"
	additionalProperties := make(map[string]interface{})
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &additionalProperties); err != nil {
		return err
	}
	d.Set("additional_properties", additionalProperties)

	return nil
}

func resourceSynapseLinkedServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.LinkedServiceClient(id.WorkspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	future, err := client.DeleteLinkedService(ctx, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for %s to be deleted: %+v", id, err)
	}

	return nil
}

func expandSynapseParameters(input map[string]interface{}) map[string]*artifacts.ParameterSpecification {
	output := make(map[string]*artifacts.ParameterSpecification)

	for k, v := range input {
		output[k] = &artifacts.ParameterSpecification{
			Type:         artifacts.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func expandSynapseLinkedServiceIntegrationRuntimeV2(input []interface{}) *artifacts.IntegrationRuntimeReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &artifacts.IntegrationRuntimeReference{
		ReferenceName: utils.String(v["name"].(string)),
		Type:          utils.String("IntegrationRuntimeReference"),
		Parameters:    v["parameters"].(map[string]interface{}),
	}
}

func flattenSynapseParameters(input map[string]*artifacts.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func flattenSynapseLinkedServiceIntegrationRuntimeV2(input *artifacts.IntegrationRuntimeReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	name := ""
	if input.ReferenceName != nil {
		name = *input.ReferenceName
	}

	return []interface{}{
		map[string]interface{}{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}

func checkLinkedServiceResponse(response *http.Response) error {
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading status response body: %+v", err)
	}
	defer response.Body.Close()

	body := make(map[string]interface{})
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return fmt.Errorf("could not parse status response: %+v", err)
	}

	if statusRaw, ok := body["status"]; ok && statusRaw != nil {
		if status, ok := statusRaw.(string); ok {
			if status == "Failed" {
				if errorRaw, ok := body["error"]; ok && errorRaw != nil {
					if responseError, ok := errorRaw.(map[string]interface{}); ok {
						if messageRaw, ok := responseError["message"]; ok && messageRaw != nil {
							if message, ok := messageRaw.(string); ok {
								return fmt.Errorf("creating/updating Linked Service: %s", message)
							}
						}
					}
				}
				// we are specifically checking for `error` in the payload but if the status is Failed, we should return what we know
				return fmt.Errorf("creating/updating Linked Service: %+v", body)
			}
		}
	}

	return nil
}
