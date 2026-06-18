// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package synapse

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/data-plane/synapse/2021-06-01-preview/linkedservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
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
				// NOTE: the go-azure-sdk `linkedservices` package does not expose an enum for the
				// linked service discriminator `type`, so this list is maintained locally. The values
				// are the exact discriminator strings previously sourced from the Synapse SDK.
				ValidateFunc: validation.StringInSlice([]string{
					"AmazonMWS",
					"AmazonRdsForOracle",
					"AmazonRdsForSqlServer",
					"AmazonRedshift",
					"AmazonS3",
					"AzureBatch",
					"AzureBlobFS",
					"AzureBlobStorage",
					"AzureDataExplorer",
					"AzureDataLakeAnalytics",
					"AzureDataLakeStore",
					"AzureDatabricks",
					"AzureDatabricksDeltaLake",
					"AzureFileStorage",
					"AzureFunction",
					"AzureKeyVault",
					"AzureML",
					"AzureMLService",
					"AzureMariaDB",
					"AzureMySql",
					"AzurePostgreSql",
					"AzureSqlDW",
					"AzureSqlDatabase",
					"AzureSqlMI",
					"AzureSearch",
					"AzureStorage",
					"AzureTableStorage",
					"Cassandra",
					"CommonDataServiceForApps",
					"Concur",
					"CosmosDb",
					"CosmosDbMongoDbApi",
					"Couchbase",
					"CustomDataSource",
					"Db2",
					"Drill",
					"Dynamics",
					"DynamicsAX",
					"DynamicsCrm",
					"Eloqua",
					"FileServer",
					"FtpServer",
					"GoogleAdWords",
					"GoogleBigQuery",
					"GoogleCloudStorage",
					"Greenplum",
					"HBase",
					"HDInsight",
					"HDInsightOnDemand",
					"HttpServer",
					"Hdfs",
					"Hive",
					"Hubspot",
					"Impala",
					"Informix",
					"Jira",
					"LinkedService",
					"Magento",
					"MariaDB",
					"Marketo",
					"MicrosoftAccess",
					"MongoDb",
					"MongoDbAtlas",
					"MongoDbV2",
					"MySql",
					"Netezza",
					"OData",
					"Odbc",
					"Office365",
					"Oracle",
					"OracleServiceCloud",
					"Paypal",
					"Phoenix",
					"PostgreSql",
					"Presto",
					"QuickBooks",
					"Responsys",
					"RestService",
					"SqlServer",
					"Salesforce",
					"SalesforceMarketingCloud",
					"SalesforceServiceCloud",
					"SapBW",
					"SapCloudForCustomer",
					"SapEcc",
					"SapHana",
					"SapOpenHub",
					"SapTable",
					"ServiceNow",
					"Sftp",
					"SharePointOnlineList",
					"Shopify",
					"Snowflake",
					"Spark",
					"Square",
					"Sybase",
					"Teradata",
					"Vertica",
					"Web",
					"Xero",
					"Zoho",
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

	endpoint := fmt.Sprintf("https://%s.%s", workspaceId.Name, *synapseDomainSuffix)

	id := parse.NewLinkedServiceID(workspaceId.SubscriptionId, workspaceId.ResourceGroup, workspaceId.Name, d.Get("name").(string))
	linkedServiceId := linkedservices.NewLinkedServiceID(endpoint, id.Name)
	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.LinkedServiceGetLinkedService(ctx, linkedServiceId, linkedservices.DefaultLinkedServiceGetLinkedServiceOperationOptions())
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_synapse_linked_service", id.ID())
			}
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

	// The linked service resource is an opaque-JSON contract: route the assembled `props` map through
	// `RawLinkedServiceImpl`, whose `MarshalJSON` re-emits `Values` verbatim, so unmodeled/arbitrary
	// `typeProperties` and sibling fields are sent exactly as supplied (see audit-synapse §4).
	input := linkedservices.LinkedServiceResource{
		Properties: linkedservices.RawLinkedServiceImpl{
			Type:   d.Get("type").(string),
			Values: props,
		},
	}

	if err := client.LinkedServiceCreateOrUpdateLinkedServiceThenPoll(ctx, linkedServiceId, input, linkedservices.DefaultLinkedServiceCreateOrUpdateLinkedServiceOperationOptions()); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

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

	endpoint := fmt.Sprintf("https://%s.%s", id.WorkspaceName, *synapseDomainSuffix)

	resp, err := client.LinkedServiceGetLinkedService(ctx, linkedservices.NewLinkedServiceID(endpoint, id.Name), linkedservices.DefaultLinkedServiceGetLinkedServiceOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", parse.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())

	// Parse the raw response body rather than the typed `resp.Model.Properties`: for known discriminator
	// types the typed model has no catch-all and would drop unmodeled `additional_properties` siblings.
	// The base client buffers and resets the body, so it is re-readable here (see audit-synapse §4.4).
	m := make(map[string]*json.RawMessage)
	if resp.HttpResponse != nil && resp.HttpResponse.Body != nil {
		rawBody, err := io.ReadAll(resp.HttpResponse.Body)
		if err != nil {
			return fmt.Errorf("reading response body for %s: %+v", id, err)
		}

		var envelope struct {
			Properties *json.RawMessage `json:"properties"`
		}
		if err := json.Unmarshal(rawBody, &envelope); err != nil {
			return fmt.Errorf("unmarshaling response for %s: %+v", id, err)
		}
		if envelope.Properties != nil {
			if err := json.Unmarshal(*envelope.Properties, &m); err != nil {
				return fmt.Errorf("unmarshaling properties for %s: %+v", id, err)
			}
		}
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

	parameters := make(map[string]*linkedservices.ParameterSpecification)
	if v, ok := m["parameters"]; ok && v != nil {
		if err := json.Unmarshal(*v, &parameters); err != nil {
			return err
		}
		delete(m, "parameters")
	}
	if err := d.Set("parameters", flattenSynapseParameters(parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	var integrationRuntime *linkedservices.IntegrationRuntimeReference
	if v, ok := m["connectVia"]; ok && v != nil {
		integrationRuntime = &linkedservices.IntegrationRuntimeReference{}
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

	endpoint := fmt.Sprintf("https://%s.%s", id.WorkspaceName, *synapseDomainSuffix)

	if err := client.LinkedServiceDeleteLinkedServiceThenPoll(ctx, linkedservices.NewLinkedServiceID(endpoint, id.Name)); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandSynapseParameters(input map[string]interface{}) map[string]*linkedservices.ParameterSpecification {
	output := make(map[string]*linkedservices.ParameterSpecification)

	for k, v := range input {
		output[k] = &linkedservices.ParameterSpecification{
			Type:         linkedservices.ParameterTypeString,
			DefaultValue: pointer.To[interface{}](v.(string)),
		}
	}

	return output
}

func expandSynapseLinkedServiceIntegrationRuntimeV2(input []interface{}) *linkedservices.IntegrationRuntimeReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	return &linkedservices.IntegrationRuntimeReference{
		ReferenceName: v["name"].(string),
		Type:          linkedservices.IntegrationRuntimeReferenceTypeIntegrationRuntimeReference,
		Parameters:    pointer.To(v["parameters"].(map[string]interface{})),
	}
}

func flattenSynapseParameters(input map[string]*linkedservices.ParameterSpecification) map[string]interface{} {
	output := make(map[string]interface{})

	for k, v := range input {
		if v != nil && v.DefaultValue != nil {
			// we only support string parameters at this time
			val, ok := (*v.DefaultValue).(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func flattenSynapseLinkedServiceIntegrationRuntimeV2(input *linkedservices.IntegrationRuntimeReference) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"name":       input.ReferenceName,
			"parameters": pointer.From(input.Parameters),
		},
	}
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return utils.NormalizeJson(old) == utils.NormalizeJson(new)
}
