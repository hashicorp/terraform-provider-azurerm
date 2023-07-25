// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package loganalytics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/operationalinsights/2019-09-01/querypackqueries"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/suppress"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LogAnalyticsQueryPackQueryModel struct {
	Name                   string            `tfschema:"name"`
	QueryPackId            string            `tfschema:"query_pack_id"`
	Body                   string            `tfschema:"body"`
	DisplayName            string            `tfschema:"display_name"`
	Categories             []string          `tfschema:"categories"`
	Description            string            `tfschema:"description"`
	AdditionalSettingsJson string            `tfschema:"additional_settings_json"`
	ResourceTypes          []string          `tfschema:"resource_types"`
	Solutions              []string          `tfschema:"solutions"`
	Tags                   map[string]string `tfschema:"tags"`
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
	return querypackqueries.ValidateQueryID
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
			ValidateFunc: querypackqueries.ValidateQueryPackID,
		},

		"body": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

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

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"additional_settings_json": {
			Type:      pluginsdk.TypeString,
			Optional:  true,
			StateFunc: utils.NormalizeJson,
		},

		"resource_types": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
				ValidateFunc: validation.StringInSlice([]string{
					"default",
					"microsoft.aad/domainservices",
					"microsoft.aadiam/tenants",
					"microsoft.agfoodplatform/farmbeats",
					"microsoft.analysisservices/servers",
					"microsoft.apimanagement/service",
					"microsoft.appconfiguration/configurationstores",
					"microsoft.appplatform/spring",
					"microsoft.attestation/attestationproviders",
					"microsoft.authorization/tenants",
					"microsoft.automation/automationaccounts",
					"microsoft.autonomousdevelopmentplatform/accounts",
					"microsoft.azurestackhci/virtualmachines",
					"microsoft.batch/batchaccounts",
					"microsoft.blockchain/blockchainmembers",
					"microsoft.botservice/botservices",
					"microsoft.cache/redis",
					"microsoft.cdn/profiles",
					"microsoft.cognitiveservices/accounts",
					"microsoft.communication/communicationservices",
					"microsoft.compute/virtualmachines",
					"microsoft.compute/virtualmachinescalesets",
					"microsoft.connectedcache/cachenodes",
					"microsoft.connectedvehicle/platformaccounts",
					"microsoft.conenctedvmwarevsphere/virtualmachines",
					"microsoft.containerregistry/registries",
					"microsoft.containerservice/managedclusters",
					"microsoft.d365customerinsights/instances",
					"microsoft.dashboard/grafana",
					"microsoft.databricks/workspaces",
					"microsoft.datacollaboration/workspaces",
					"microsoft.datafactory/factories",
					"microsoft.datalakeanalytics/accounts",
					"microsoft.datalakestore/accounts",
					"microsoft.datashare/accounts",
					"microsoft.dbformariadb/servers",
					"microsoft.dbformysql/servers",
					"microsoft.dbforpostgresql/flexibleservers",
					"microsoft.dbforpostgresql/servers",
					"microsoft.dbforpostgresql/serversv2",
					"microsoft.digitaltwins/digitaltwinsinstances",
					"microsoft.documentdb/cassandraclusters",
					"microsoft.documentdb/databaseaccounts",
					"microsoft.desktopvirtualization/applicationgroups",
					"microsoft.desktopvirtualization/hostpools",
					"microsoft.desktopvirtualization/workspaces",
					"microsoft.devices/iothubs",
					"microsoft.devices/provisioningservices",
					"microsoft.dynamics/fraudprotection/purchase",
					"microsoft.eventgrid/domains",
					"microsoft.eventgrid/topics",
					"microsoft.eventgrid/partnernamespaces",
					"microsoft.eventgrid/partnertopics",
					"microsoft.eventgrid/systemtopics",
					"microsoft.eventhub/namespaces",
					"microsoft.experimentation/experimentworkspaces",
					"microsoft.hdinsight/clusters",
					"microsoft.healthcareapis/services",
					"microsoft.informationprotection/datasecuritymanagement",
					"microsoft.intune/operations",
					"microsoft.insights/autoscalesettings",
					"microsoft.insights/components",
					"microsoft.insights/workloadmonitoring",
					"microsoft.keyvault/vaults",
					"microsoft.kubernetes/connectedclusters",
					"microsoft.kusto/clusters",
					"microsoft.loadtestservice/loadtests",
					"microsoft.logic/workflows",
					"microsoft.machinelearningservices/workspaces",
					"microsoft.media/mediaservices",
					"microsoft.netapp/netappaccounts/capacitypools",
					"microsoft.network/applicationgateways",
					"microsoft.network/azurefirewalls",
					"microsoft.network/bastionhosts",
					"microsoft.network/expressroutecircuits",
					"microsoft.network/frontdoors",
					"microsoft.network/loadbalancers",
					"microsoft.network/networkinterfaces",
					"microsoft.network/networksecuritygroups",
					"microsoft.network/networksecurityperimeters",
					"microsoft.network/networkwatchers/connectionmonitors",
					"microsoft.network/networkwatchers/trafficanalytics",
					"microsoft.network/publicipaddresses",
					"microsoft.network/trafficmanagerprofiles",
					"microsoft.network/virtualnetworks",
					"microsoft.network/virtualnetworkgateways",
					"microsoft.network/vpngateways",
					"microsoft.networkfunction/azuretrafficcollectors",
					"microsoft.openenergyplatform/energyservices",
					"microsoft.openlogisticsplatform/workspaces",
					"microsoft.operationalinsights/workspaces",
					"microsoft.powerbi/tenants",
					"microsoft.powerbi/tenants/workspaces",
					"microsoft.powerbidedicated/capacities",
					"microsoft.purview/accounts",
					"microsoft.recoveryservices/vaults",
					"microsoft.resources/azureactivity",
					"microsoft.scvmm/virtualmachines",
					"microsoft.search/searchservices",
					"microsoft.security/antimalwaresettings",
					"microsoft.securityinsights/amazon",
					"microsoft.securityinsights/anomalies",
					"microsoft.securityinsights/cef",
					"microsoft.securityinsights/datacollection",
					"microsoft.securityinsights/dnsnormalized",
					"microsoft.securityinsights/mda",
					"microsoft.securityinsights/mde",
					"microsoft.securityinsights/mdi",
					"microsoft.securityinsights/mdo",
					"microsoft.securityinsights/networksessionnormalized",
					"microsoft.securityinsights/office365",
					"microsoft.securityinsights/purview",
					"microsoft.securityinsights/securityinsights",
					"microsoft.securityinsights/securityinsights/mcas",
					"microsoft.securityinsights/tvm",
					"microsoft.securityinsights/watchlists",
					"microsoft.servicebus/namespaces",
					"microsoft.servicefabric/clusters",
					"microsoft.signalrservice/signalr",
					"microsoft.signalrservice/webpubsub",
					"microsoft.sql/managedinstances",
					"microsoft.sql/servers",
					"microsoft.sql/servers/databases",
					"microsoft.storage/storageaccounts",
					"microsoft.storagecache/caches",
					"microsoft.streamanalytics/streamingjobs",
					"microsoft.synapse/workspaces",
					"microsoft.timeseriesinsights/environments",
					"microsoft.videoindexer/accounts",
					"microsoft.web/sites",
					"microsoft.workloadmonitor/monitors",
					"resourcegroup",
					"subscription",
				}, false),
			},
		},

		"solutions": {
			Type:             pluginsdk.TypeList,
			Optional:         true,
			DiffSuppressFunc: suppress.CaseDifference,
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
					"Security",
					"SecurityCenter",
					"SecurityCenterFree",
					"SecurityInsights",
					"ServiceMap",
					"SfBAssessment",
					"SfBOnlineAssessment",
					"SharePointOnlineAssessment",
					"SPAssessment",
					"SQLAdvancedThreatProtection",
					"SQLAssessment",
					"SQLAssessmentPlus",
					"SQLDataClassification",
					"SQLThreatDetection",
					"SQLVulnerabilityAssessment",
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
				}, true),
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

			client := metadata.Client.LogAnalytics.QueryPackQueriesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			queryPackId, _ := querypackqueries.ParseQueryPackID(model.QueryPackId)

			if model.Name == "" {
				uuid, err := uuid.GenerateUUID()
				if err != nil {
					return fmt.Errorf("generating UUID for Log Analytics Query Pack Query: %+v", err)
				}

				model.Name = uuid
			}

			id := querypackqueries.NewQueryID(subscriptionId, queryPackId.ResourceGroupName, queryPackId.QueryPackName, model.Name)

			existing, err := client.QueriesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := &querypackqueries.LogAnalyticsQueryPackQuery{
				Properties: &querypackqueries.LogAnalyticsQueryPackQueryProperties{
					Body:        model.Body,
					DisplayName: model.DisplayName,
					Related:     &querypackqueries.LogAnalyticsQueryPackQueryPropertiesRelated{},
				},
			}

			if model.Description != "" {
				parameters.Properties.Description = &model.Description
			}

			if model.AdditionalSettingsJson != "" {
				var additionalSettingsJson interface{}
				if err := json.Unmarshal([]byte(model.AdditionalSettingsJson), &additionalSettingsJson); err != nil {
					return fmt.Errorf("parsing JSON: %+v", err)
				}
				parameters.Properties.Properties = &additionalSettingsJson
			}

			if model.Categories != nil {
				parameters.Properties.Related.Categories = &model.Categories
			}

			if model.ResourceTypes != nil {
				parameters.Properties.Related.ResourceTypes = &model.ResourceTypes
			}

			if model.Solutions != nil {
				parameters.Properties.Related.Solutions = &model.Solutions
			}

			if model.Tags != nil {
				parameters.Properties.Tags = expandLogAnalyticsQueryPackQueryTags(model.Tags)
			}

			if _, err := client.QueriesPut(ctx, id, *parameters); err != nil {
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
			client := metadata.Client.LogAnalytics.QueryPackQueriesClient

			id, err := querypackqueries.ParseQueryID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model LogAnalyticsQueryPackQueryModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			parameters := &querypackqueries.LogAnalyticsQueryPackQuery{
				Properties: &querypackqueries.LogAnalyticsQueryPackQueryProperties{
					Body:        model.Body,
					DisplayName: model.DisplayName,
					Related:     &querypackqueries.LogAnalyticsQueryPackQueryPropertiesRelated{},
				},
			}

			if metadata.ResourceData.HasChange("description") {
				parameters.Properties.Description = &model.Description
			}

			if metadata.ResourceData.HasChange("additional_settings_json") {
				var additionalSettingsJson interface{}
				if err := json.Unmarshal([]byte(model.AdditionalSettingsJson), &additionalSettingsJson); err != nil {
					return fmt.Errorf("parsing JSON: %+v", err)
				}
				parameters.Properties.Properties = &additionalSettingsJson
			}

			if metadata.ResourceData.HasChange("categories") {
				parameters.Properties.Related.Categories = &model.Categories
			}

			if metadata.ResourceData.HasChange("resource_types") {
				parameters.Properties.Related.ResourceTypes = &model.ResourceTypes
			}

			if metadata.ResourceData.HasChange("solutions") {
				parameters.Properties.Related.Solutions = &model.Solutions
			}

			if metadata.ResourceData.HasChange("tags") {
				parameters.Properties.Tags = expandLogAnalyticsQueryPackQueryTags(model.Tags)
			}

			if _, err := client.QueriesUpdate(ctx, *id, *parameters); err != nil {
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
			client := metadata.Client.LogAnalytics.QueryPackQueriesClient

			id, err := querypackqueries.ParseQueryID(metadata.ResourceData.Id())
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

			props := model.Properties
			if props == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			state := LogAnalyticsQueryPackQueryModel{
				Name:        id.QueryName,
				QueryPackId: querypackqueries.NewQueryPackID(id.SubscriptionId, id.ResourceGroupName, id.QueryPackName).ID(),
				Body:        props.Body,
				DisplayName: props.DisplayName,
			}

			if props.Properties != nil {
				propsJson, jsonErr := json.Marshal(props.Properties)
				if jsonErr != nil {
					return fmt.Errorf("parsing JSON for Log Analytics Query Pack Query Properties: %+v", jsonErr)
				}
				state.AdditionalSettingsJson = string(propsJson)
			}

			if props.Description != nil {
				state.Description = *props.Description
			}

			if additionalSettings := props.Related; additionalSettings != nil {
				if additionalSettings.Categories != nil {
					state.Categories = *additionalSettings.Categories
				}

				if additionalSettings.ResourceTypes != nil {
					state.ResourceTypes = *additionalSettings.ResourceTypes
				}

				if additionalSettings.Solutions != nil {
					state.Solutions = *additionalSettings.Solutions
				}
			}

			if tags := props.Tags; tags != nil {
				state.Tags = flattenLogAnalyticsQueryPackQueryTags(*tags)
			}

			return metadata.Encode(&state)
		},
	}
}

func (r LogAnalyticsQueryPackQueryResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LogAnalytics.QueryPackQueriesClient

			id, err := querypackqueries.ParseQueryID(metadata.ResourceData.Id())
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

func expandLogAnalyticsQueryPackQueryTags(input map[string]string) *map[string][]string {
	if input == nil {
		return nil
	}

	output := make(map[string][]string)
	for k, v := range input {
		output[k] = strings.Split(v, ",")
	}

	return &output
}

func flattenLogAnalyticsQueryPackQueryTags(input map[string][]string) map[string]string {
	if input == nil {
		return nil
	}

	results := make(map[string]string)
	for k, v := range input {
		results[k] = strings.Join(v, ",")
	}

	return results
}
