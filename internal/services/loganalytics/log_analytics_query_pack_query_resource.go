package loganalytics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	queryPacks "github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/operationalinsights"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsQueryPackQueryModel struct {
	Name           string                 `tfschema:"name"`
	QueryPackId    string                 `tfschema:"query_pack_id"`
	Body           string                 `tfschema:"body"`
	Description    *string                `tfschema:"description"`
	DisplayName    string                 `tfschema:"display_name"`
	PropertiesJson string                 `tfschema:"properties_json"`
	Related        []Related              `tfschema:"related"`
	Tags           map[string]interface{} `tfschema:"tags"`
}

type Related struct {
	Categories    []interface{} `tfschema:"categories"`
	ResourceTypes []interface{} `tfschema:"resource_types"`
	Solutions     []interface{} `tfschema:"solutions"`
}

type LogAnalyticsQueryPackQueryResource struct{}

var _ sdk.ResourceWithUpdate = LogAnalyticsQueryPackQueryResource{}

func (r LogAnalyticsQueryPackQueryResource) ResourceType() string {
	return "azurerm_log_analytics_query_pack_query"
}

func (r LogAnalyticsQueryPackQueryResource) ModelObject() interface{} {
	return &LogAnalyticsQueryPackQueryModel{}
}

func (r LogAnalyticsQueryPackQueryResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return queryPacks.ValidateQueriesID
}

func (r LogAnalyticsQueryPackQueryResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			Computed:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},

		"query_pack_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: queryPacks.ValidateQueryPackID,
		},

		"body": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"properties_json": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			StateFunc: utils.NormalizeJson,
		},

		"related": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"categories": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"applications",
								"audit",
								"container",
								"databases",
								"desktopanalytics",
								"management",
								"monitor",
								"network",
								"resources",
								"security",
								"virtualmachines",
								"windowsvirtualdesktop",
								"workloads",
							}, false),
						},
					},

					"resource_types": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"microsoft.agfoodplatform/farmbeats",
								"microsoft.informationprotection/datasecuritymanagement",
								"microsoft.appconfiguration/configurationstores",
								"microsoft.web/sites",
								"microsoft.authorization/tenants",
								"microsoft.autonomousdevelopmentplatform/accounts",
								"microsoft.resources/azureactivity",
								"microsoft.attestation/attestationproviders",
								"microsoft.cache/redis",
								"microsoft.communication/communicationservices",
								"microsoft.documentdb/databaseaccounts",
								"microsoft.datacollaboration/workspaces",
								"microsoft.security/antimalwaresettings",
								"microsoft.digitaltwins/digitaltwinsinstances",
								"microsoft.eventgrid/topics",
								"microsoft.network/azurefirewalls",
								"microsoft.dashboard/grafana",
								"microsoft.loadtestservice/loadtests",
								"microsoft.documentdb/cassandraclusters",
								"microsoft.containerservice/managedclusters",
								"microsoft.insights/workloadmonitoring",
								"microsoft.netapp/netappaccounts/capacitypools",
								"microsoft.securityinsights/purview",
								"microsoft.purview/accounts",
								"microsoft.networkfunction/azuretrafficcollectors",
								"microsoft.botservice/botservices",
								"microsoft.connectedcache/cachenodes",
								"microsoft.connectedvehicle/platformaccounts",
								"microsoft.network/networkwatchers/connectionmonitors",
								"microsoft.d365customerinsights/instances",
								"microsoft.dynamics/fraudprotection/purchase",
								"microsoft.experimentation/experimentworkspaces",
								"microsoft.hdinsight/clusters",
								"microsoft.intune/operations",
								"microsoft.aadiam/tenants",
								"microsoft.machinelearningservices/workspaces",
								"microsoft.network/networksecurityperimeters",
								"microsoft.openenergyplatform/energyservices",
								"microsoft.openlogisticsplatform/workspaces",
								"microsoft.compute/virtualmachines",
								"microsoft.operationalinsights/workspaces",
								"microsoft.powerbi/tenants",
								"microsoft.powerbi/tenants/workspaces",
								"microsoft.securityinsights/cef",
								"microsoft.securityinsights/datacollection",
								"microsoft.securityinsights/anomalies",
								"microsoft.securityinsights/dnsnormalized",
								"microsoft.securityinsights/networksessionnormalized",
								"microsoft.securityinsights/amazon",
								"microsoft.securityinsights/securityinsights/mcas",
								"microsoft.securityinsights/mda",
								"microsoft.securityinsights/mde",
								"microsoft.securityinsights/mdi",
								"microsoft.securityinsights/mdo",
								"microsoft.securityinsights/office365",
								"microsoft.securityinsights/securityinsights",
								"microsoft.securityinsights/tvm",
								"microsoft.securityinsights/watchlists",
								"microsoft.storagecache/caches",
								"microsoft.synapse/workspaces",
								"microsoft.network/networkwatchers/trafficanalytics",
								"microsoft.videoindexer/accounts",
								"microsoft.desktopvirtualization/hostpools",
								"default",
								"subscription",
								"resourcegroup",
								"microsoft.signalrservice/webpubsub",
								"microsoft.insights/components",
								"microsoft.desktopvirtualization/applicationgroups",
								"microsoft.desktopvirtualization/workspaces",
								"microsoft.timeseriesinsights/environments",
								"microsoft.workloadmonitor/monitors",
								"microsoft.analysisservices/servers",
								"microsoft.batch/batchaccounts",
								"microsoft.cdn/profiles",
								"microsoft.appplatform/spring",
								"microsoft.media/mediaservices",
								"microsoft.cognitiveservices/accounts",
								"microsoft.keyvault/vaults",
								"microsoft.storage/storageaccounts",
								"microsoft.signalrservice/signalr",
								"microsoft.containerregistry/registries",
								"microsoft.kusto/clusters",
								"microsoft.aad/domainservices",
								"microsoft.blockchain/blockchainmembers",
								"microsoft.eventgrid/domains",
								"microsoft.eventgrid/partnernamespaces",
								"microsoft.eventgrid/partnertopics",
								"microsoft.eventgrid/systemtopics",
								"microsoft.conenctedvmwarevsphere/virtualmachines",
								"microsoft.azurestackhci/virtualmachines",
								"microsoft.scvmm/virtualmachines",
								"microsoft.compute/virtualmachinescalesets",
								"microsoft.kubernetes/connectedclusters",
								"microsoft.databricks/workspaces",
								"microsoft.insights/autoscalesettings",
								"microsoft.devices/iothubs",
								"microsoft.servicefabric/clusters",
								"microsoft.logic/workflows",
								"microsoft.apimanagement/service",
								"microsoft.automation/automationaccounts",
								"microsoft.datafactory/factories",
								"microsoft.recoveryservices/vaults",
								"microsoft.datalakestore/accounts",
								"microsoft.datalakeanalytics/accounts",
								"microsoft.powerbidedicated/capacities",
								"microsoft.datashare/accounts",
								"microsoft.sql/managedinstances",
								"microsoft.sql/servers",
								"microsoft.sql/servers/databases",
								"microsoft.dbformysql/servers",
								"microsoft.dbforpostgresql/servers",
								"microsoft.dbforpostgresql/serversv2",
								"microsoft.dbforpostgresql/flexibleservers",
								"microsoft.dbformariadb/servers",
								"microsoft.devices/provisioningservices",
								"microsoft.eventhub/namespaces",
								"microsoft.network/applicationgateways",
								"microsoft.network/expressroutecircuits",
								"microsoft.network/frontdoors",
								"microsoft.network/loadbalancers",
								"microsoft.network/networkinterfaces",
								"microsoft.network/networksecuritygroups",
								"microsoft.network/publicipaddresses",
								"microsoft.network/trafficmanagerprofiles",
								"microsoft.network/virtualnetworkgateways",
								"microsoft.network/vpngateways",
								"microsoft.network/virtualnetworks",
								"microsoft.search/searchservices",
								"microsoft.streamanalytics/streamingjobs",
								"microsoft.network/bastionhosts",
								"microsoft.healthcareapis/services",
								"microsoft.servicebus/namespaces",
							}, false),
						},
					},

					"solutions": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice([]string{
								"AADDomainServices",
								"ADAssessment",
								"ADAssessmentPlus",
								"ADReplication",
								"ADSecurityAssessment",
								"AlertManagement",
								"AntiMalware",
								"ApplicationInsights",
								"AzureAssessment",
								"AzureSecurityOfThings",
								"AzureSentinelDSRE",
								"AzureSentinelPrivatePreview",
								"BehaviorAnalyticsInsights",
								"ChangeTracking",
								"CompatibilityAssessment",
								"ContainerInsights",
								"Containers",
								"CustomizedWindowsEventsFiltering",
								"DeviceHealthProd",
								"DnsAnalytics",
								"ExchangeAssessment",
								"ExchangeOnlineAssessment",
								"IISAssessmentPlus",
								"InfrastructureInsights",
								"InternalWindowsEvent",
								"LogManagement",
								"Microsoft365Analytics",
								"NetworkMonitoring",
								"SCCMAssessmentPlus",
								"SCOMAssessment",
								"SCOMAssessmentPlus",
								"SPAssessment",
								"SQLAdvancedThreatProtection",
								"SQLAssessment",
								"SQLAssessmentPlus",
								"SQLDataClassification",
								"SQLThreatDetection",
								"SQLVulnerabilityAssessment",
								"Security",
								"SecurityCenter",
								"SecurityCenterFree",
								"SecurityInsights",
								"ServiceMap",
								"SfBAssessment",
								"SfBOnlineAssessment",
								"SharePointOnlineAssessment",
								"SurfaceHub",
								"Updates",
								"VMInsights",
								"WEFInternalUat",
								"WEF_10x",
								"WEF_10xDSRE",
								"WaaSUpdateInsights",
								"WinLog",
								"WindowsClientAssessmentPlus",
								"WindowsEventForwarding",
								"WindowsFirewall",
								"WindowsServerAssessment",
								"WireData",
								"WireData2",
							}, false),
						},
					},
				},
			},
		},

		"tags": commonschema.Tags(),
	}
}

func (r LogAnalyticsQueryPackQueryResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LogAnalyticsQueryPackQueryResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LogAnalyticsQueryPackQueryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.LogAnalytics.QueryPacksClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			queryPackId, _ := queryPacks.ParseQueryPackID(model.QueryPackId)

			if model.Name == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Log Analytics Query Pack Query: %+v", err)
				}

				model.Name = uuid
			}

			id := queryPacks.NewQueriesID(subscriptionId, queryPackId.ResourceGroupName, queryPackId.QueryPackName, model.Name)

			existing, err := client.QueriesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			properties := &queryPacks.LogAnalyticsQueryPackQuery{
				Properties: &queryPacks.LogAnalyticsQueryPackQueryProperties{
					Body:        model.Body,
					DisplayName: model.DisplayName,
				},
			}

			if model.Description != nil {
				properties.Properties.Description = model.Description
			}

			if model.PropertiesJson != "" {
				var propertiesJson interface{}
				if err := json.Unmarshal([]byte(model.PropertiesJson), &propertiesJson); err != nil {
					return fmt.Errorf("parsing JSON: %+v", err)
				}
				properties.Properties.Properties = &propertiesJson
			}

			if model.Related != nil {
				properties.Properties.Related = expandLogAnalyticsQueryPackQueryRelated(model.Related)
			}

			if model.Tags != nil {
				properties.Properties.Tags = expandLogAnalyticsQueryPackQueryTags(model.Tags)
			}

			if _, err := client.QueriesPut(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r LogAnalyticsQueryPackQueryResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := queryPacks.ParseQueriesID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LogAnalyticsQueryPackQueryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.QueriesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("body") {
				properties.Properties.Body = model.Body
			}

			if metadata.ResourceData.HasChange("description") {
				properties.Properties.Description = model.Description
			}

			if metadata.ResourceData.HasChange("display_name") {
				properties.Properties.DisplayName = model.DisplayName
			}

			if metadata.ResourceData.HasChange("properties_json") {
				var propertiesJson interface{}
				if err := json.Unmarshal([]byte(model.PropertiesJson), &propertiesJson); err != nil {
					return fmt.Errorf("parsing JSON: %+v", err)
				}
				properties.Properties.Properties = &propertiesJson
			}

			if metadata.ResourceData.HasChange("related") {
				properties.Properties.Related = expandLogAnalyticsQueryPackQueryRelated(model.Related)
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Properties.Tags = expandLogAnalyticsQueryPackQueryTags(model.Tags)
			}

			if _, err := client.QueriesUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r LogAnalyticsQueryPackQueryResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := queryPacks.ParseQueriesID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.QueriesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			propsJson, jsonErr := json.Marshal(model.Properties.Properties)
			if jsonErr != nil {
				return fmt.Errorf("parsing JSON for Log Analytics Query Pack Query Properties: %+v", jsonErr)
			}

			state := LogAnalyticsQueryPackQueryModel{
				Name:           id.QueryPackName,
				QueryPackId:    queryPacks.NewQueryPackID(id.SubscriptionId, id.ResourceGroupName, id.QueryPackName).ID(),
				Body:           model.Properties.Body,
				DisplayName:    model.Properties.DisplayName,
				PropertiesJson: string(propsJson),
				Related:        flattenLogAnalyticsQueryPackQueryRelated(model.Properties.Related),
				Tags:           flattenLogAnalyticsQueryPackQueryTags(*model.Properties.Tags),
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogAnalyticsQueryPackQueryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPacksClient

			id, err := queryPacks.ParseQueriesID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.QueriesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandLogAnalyticsQueryPackQueryTags(input map[string]interface{}) *map[string][]string {
	output := make(map[string][]string)
	for k, v := range input {
		output[k] = strings.Split(v.(string), ",")
	}
	return &output
}

func flattenLogAnalyticsQueryPackQueryTags(input map[string][]string) map[string]interface{} {
	results := make(map[string]interface{})
	if input == nil {
		return results
	}
	for k, v := range input {
		results[k] = strings.Join(v, ",")
	}
	return results
}

func expandLogAnalyticsQueryPackQueryRelated(input []Related) *queryPacks.LogAnalyticsQueryPackQueryPropertiesRelated {
	if len(input) == 0 {
		return nil
	}

	result := &queryPacks.LogAnalyticsQueryPackQueryPropertiesRelated{}

	if input[0].Categories != nil {
		result.Categories = utils.ExpandStringSlice(input[0].Categories)
	}

	if input[0].ResourceTypes != nil {
		result.ResourceTypes = utils.ExpandStringSlice(input[0].ResourceTypes)
	}

	if input[0].Solutions != nil {
		result.Solutions = utils.ExpandStringSlice(input[0].Solutions)
	}

	return result
}

func flattenLogAnalyticsQueryPackQueryRelated(input *queryPacks.LogAnalyticsQueryPackQueryPropertiesRelated) []Related {
	if input == nil {
		return make([]Related, 0)
	}

	result := make([]Related, 0)
	result = append(result, Related{
		Categories:    utils.FlattenStringSlice(input.Categories),
		ResourceTypes: utils.FlattenStringSlice(input.ResourceTypes),
		Solutions:     utils.FlattenStringSlice(input.Solutions),
	})
	return result
}
