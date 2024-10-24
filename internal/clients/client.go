// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package clients

import (
	"context"
	"fmt"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/validation"
	aadb2c_v2021_04_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/aadb2c/2021-04-01-preview"
	analysisservices_v2017_08_01 "github.com/hashicorp/go-azure-sdk/resource-manager/analysisservices/2017-08-01"
	azurestackhci_v2024_01_01 "github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01"
	datadog_v2021_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/datadog/2021-03-01"
	dns_v2018_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/dns/2018-05-01"
	eventgrid_v2022_06_15 "github.com/hashicorp/go-azure-sdk/resource-manager/eventgrid/2022-06-15"
	fluidrelay_2022_05_26 "github.com/hashicorp/go-azure-sdk/resource-manager/fluidrelay/2022-05-26"
	hdinsight_v2021_06_01 "github.com/hashicorp/go-azure-sdk/resource-manager/hdinsight/2021-06-01"
	nginx_2024_06_01_preview "github.com/hashicorp/go-azure-sdk/resource-manager/nginx/2024-06-01-preview"
	redis_2024_03_01 "github.com/hashicorp/go-azure-sdk/resource-manager/redis/2024-03-01"
	servicenetworking_2023_11_01 "github.com/hashicorp/go-azure-sdk/resource-manager/servicenetworking/2023-11-01"
	storagecache_2023_05_01 "github.com/hashicorp/go-azure-sdk/resource-manager/storagecache/2023-05-01"
	systemcentervirtualmachinemanager_2023_10_07 "github.com/hashicorp/go-azure-sdk/resource-manager/systemcentervirtualmachinemanager/2023-10-07"
	workloads_v2023_04_01 "github.com/hashicorp/go-azure-sdk/resource-manager/workloads/2023-04-01"
	"github.com/hashicorp/terraform-provider-azurerm/internal/common"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	aadb2c "github.com/hashicorp/terraform-provider-azurerm/internal/services/aadb2c/client"
	advisor "github.com/hashicorp/terraform-provider-azurerm/internal/services/advisor/client"
	analysisServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/analysisservices/client"
	apiManagement "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	appConfiguration "github.com/hashicorp/terraform-provider-azurerm/internal/services/appconfiguration/client"
	applicationInsights "github.com/hashicorp/terraform-provider-azurerm/internal/services/applicationinsights/client"
	appService "github.com/hashicorp/terraform-provider-azurerm/internal/services/appservice/client"
	arckubernetes "github.com/hashicorp/terraform-provider-azurerm/internal/services/arckubernetes/client"
	arcResourceBridge "github.com/hashicorp/terraform-provider-azurerm/internal/services/arcresourcebridge/client"
	attestation "github.com/hashicorp/terraform-provider-azurerm/internal/services/attestation/client"
	authorization "github.com/hashicorp/terraform-provider-azurerm/internal/services/authorization/client"
	automanage "github.com/hashicorp/terraform-provider-azurerm/internal/services/automanage/client"
	automation "github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/client"
	azureStackHCI "github.com/hashicorp/terraform-provider-azurerm/internal/services/azurestackhci/client"
	batch "github.com/hashicorp/terraform-provider-azurerm/internal/services/batch/client"
	blueprints "github.com/hashicorp/terraform-provider-azurerm/internal/services/blueprints/client"
	bot "github.com/hashicorp/terraform-provider-azurerm/internal/services/bot/client"
	cdn "github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/client"
	cognitiveServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/cognitive/client"
	communication "github.com/hashicorp/terraform-provider-azurerm/internal/services/communication/client"
	compute "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/client"
	confidentialledger "github.com/hashicorp/terraform-provider-azurerm/internal/services/confidentialledger/client"
	connections "github.com/hashicorp/terraform-provider-azurerm/internal/services/connections/client"
	consumption "github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/client"
	containerapps "github.com/hashicorp/terraform-provider-azurerm/internal/services/containerapps/client"
	containerServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/containers/client"
	cosmosdb "github.com/hashicorp/terraform-provider-azurerm/internal/services/cosmos/client"
	costmanagement "github.com/hashicorp/terraform-provider-azurerm/internal/services/costmanagement/client"
	customproviders "github.com/hashicorp/terraform-provider-azurerm/internal/services/customproviders/client"
	dashboard "github.com/hashicorp/terraform-provider-azurerm/internal/services/dashboard/client"
	datamigration "github.com/hashicorp/terraform-provider-azurerm/internal/services/databasemigration/client"
	databoxedge "github.com/hashicorp/terraform-provider-azurerm/internal/services/databoxedge/client"
	databricks "github.com/hashicorp/terraform-provider-azurerm/internal/services/databricks/client"
	datadog "github.com/hashicorp/terraform-provider-azurerm/internal/services/datadog/client"
	datafactory "github.com/hashicorp/terraform-provider-azurerm/internal/services/datafactory/client"
	dataprotection "github.com/hashicorp/terraform-provider-azurerm/internal/services/dataprotection/client"
	datashare "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	desktopvirtualization "github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/client"
	devtestlabs "github.com/hashicorp/terraform-provider-azurerm/internal/services/devtestlabs/client"
	digitaltwins "github.com/hashicorp/terraform-provider-azurerm/internal/services/digitaltwins/client"
	dns "github.com/hashicorp/terraform-provider-azurerm/internal/services/dns/client"
	domainservices "github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/client"
	elastic "github.com/hashicorp/terraform-provider-azurerm/internal/services/elastic/client"
	elasticsan "github.com/hashicorp/terraform-provider-azurerm/internal/services/elasticsan/client"
	eventgrid "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventgrid/client"
	eventhub "github.com/hashicorp/terraform-provider-azurerm/internal/services/eventhub/client"
	extendedlocation "github.com/hashicorp/terraform-provider-azurerm/internal/services/extendedlocation/client"
	fluidrelay "github.com/hashicorp/terraform-provider-azurerm/internal/services/fluidrelay/client"
	frontdoor "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/client"
	graph "github.com/hashicorp/terraform-provider-azurerm/internal/services/graphservices/client"
	hdinsight "github.com/hashicorp/terraform-provider-azurerm/internal/services/hdinsight/client"
	healthcare "github.com/hashicorp/terraform-provider-azurerm/internal/services/healthcare/client"
	hsm "github.com/hashicorp/terraform-provider-azurerm/internal/services/hsm/client"
	hybridcompute "github.com/hashicorp/terraform-provider-azurerm/internal/services/hybridcompute/client"
	iotcentral "github.com/hashicorp/terraform-provider-azurerm/internal/services/iotcentral/client"
	iothub "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/client"
	keyvault "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/client"
	kusto "github.com/hashicorp/terraform-provider-azurerm/internal/services/kusto/client"
	lighthouse "github.com/hashicorp/terraform-provider-azurerm/internal/services/lighthouse/client"
	loadbalancers "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadbalancer/client"
	loadtestservice "github.com/hashicorp/terraform-provider-azurerm/internal/services/loadtestservice/client"
	loganalytics "github.com/hashicorp/terraform-provider-azurerm/internal/services/loganalytics/client"
	logic "github.com/hashicorp/terraform-provider-azurerm/internal/services/logic/client"
	machinelearning "github.com/hashicorp/terraform-provider-azurerm/internal/services/machinelearning/client"
	maintenance "github.com/hashicorp/terraform-provider-azurerm/internal/services/maintenance/client"
	managedapplication "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedapplications/client"
	managedhsm "github.com/hashicorp/terraform-provider-azurerm/internal/services/managedhsm/client"
	managementgroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/client"
	maps "github.com/hashicorp/terraform-provider-azurerm/internal/services/maps/client"
	mixedreality "github.com/hashicorp/terraform-provider-azurerm/internal/services/mixedreality/client"
	mobilenetwork "github.com/hashicorp/terraform-provider-azurerm/internal/services/mobilenetwork/client"
	monitor "github.com/hashicorp/terraform-provider-azurerm/internal/services/monitor/client"
	mssql "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssql/client"
	mssqlmanagedinstance "github.com/hashicorp/terraform-provider-azurerm/internal/services/mssqlmanagedinstance/client"
	mysql "github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/client"
	netapp "github.com/hashicorp/terraform-provider-azurerm/internal/services/netapp/client"
	network "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	networkfunction "github.com/hashicorp/terraform-provider-azurerm/internal/services/networkfunction/client"
	newrelic "github.com/hashicorp/terraform-provider-azurerm/internal/services/newrelic/client"
	nginx "github.com/hashicorp/terraform-provider-azurerm/internal/services/nginx/client"
	notificationhub "github.com/hashicorp/terraform-provider-azurerm/internal/services/notificationhub/client"
	oracle "github.com/hashicorp/terraform-provider-azurerm/internal/services/oracle/client"
	orbital "github.com/hashicorp/terraform-provider-azurerm/internal/services/orbital/client"
	paloalto "github.com/hashicorp/terraform-provider-azurerm/internal/services/paloalto/client"
	policy "github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/client"
	portal "github.com/hashicorp/terraform-provider-azurerm/internal/services/portal/client"
	postgres "github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/client"
	powerBI "github.com/hashicorp/terraform-provider-azurerm/internal/services/powerbi/client"
	privatedns "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/client"
	dnsresolver "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatednsresolver/client"
	purview "github.com/hashicorp/terraform-provider-azurerm/internal/services/purview/client"
	recoveryServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/recoveryservices/client"
	redhatopenshift "github.com/hashicorp/terraform-provider-azurerm/internal/services/redhatopenshift/client"
	redis "github.com/hashicorp/terraform-provider-azurerm/internal/services/redis/client"
	redisenterprise "github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/client"
	relay "github.com/hashicorp/terraform-provider-azurerm/internal/services/relay/client"
	resource "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/client"
	search "github.com/hashicorp/terraform-provider-azurerm/internal/services/search/client"
	securityCenter "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	sentinel "github.com/hashicorp/terraform-provider-azurerm/internal/services/sentinel/client"
	serviceBus "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicebus/client"
	serviceConnector "github.com/hashicorp/terraform-provider-azurerm/internal/services/serviceconnector/client"
	serviceFabric "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabric/client"
	serviceFabricManaged "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicefabricmanaged/client"
	serviceNetworking "github.com/hashicorp/terraform-provider-azurerm/internal/services/servicenetworking/client"
	signalr "github.com/hashicorp/terraform-provider-azurerm/internal/services/signalr/client"
	appPlatform "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	storage "github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/client"
	storageCache "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagecache/client"
	storageMover "github.com/hashicorp/terraform-provider-azurerm/internal/services/storagemover/client"
	streamAnalytics "github.com/hashicorp/terraform-provider-azurerm/internal/services/streamanalytics/client"
	subscription "github.com/hashicorp/terraform-provider-azurerm/internal/services/subscription/client"
	synapse "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	systemCenterVirtualMachineManager "github.com/hashicorp/terraform-provider-azurerm/internal/services/systemcentervirtualmachinemanager/client"
	trafficManager "github.com/hashicorp/terraform-provider-azurerm/internal/services/trafficmanager/client"
	videoindexer "github.com/hashicorp/terraform-provider-azurerm/internal/services/videoindexer/client"
	vmware "github.com/hashicorp/terraform-provider-azurerm/internal/services/vmware/client"
	voiceServices "github.com/hashicorp/terraform-provider-azurerm/internal/services/voiceservices/client"
	web "github.com/hashicorp/terraform-provider-azurerm/internal/services/web/client"
	workloads "github.com/hashicorp/terraform-provider-azurerm/internal/services/workloads/client"
)

type Client struct {
	autoClient

	// StopContext is used for propagating control from Terraform Core (e.g. Ctrl/Cmd+C)
	StopContext context.Context

	Account  *ResourceManagerAccount
	Features features.UserFeatures

	AadB2c                            *aadb2c_v2021_04_01_preview.Client
	Advisor                           *advisor.Client
	AnalysisServices                  *analysisservices_v2017_08_01.Client
	ApiManagement                     *apiManagement.Client
	AppConfiguration                  *appConfiguration.Client
	AppInsights                       *applicationInsights.Client
	AppPlatform                       *appPlatform.Client
	AppService                        *appService.Client
	ArcKubernetes                     *arckubernetes.Client
	ArcResourceBridge                 *arcResourceBridge.Client
	Attestation                       *attestation.Client
	Authorization                     *authorization.Client
	Automanage                        *automanage.Client
	Automation                        *automation.Client
	AzureStackHCI                     *azurestackhci_v2024_01_01.Client
	Batch                             *batch.Client
	Blueprints                        *blueprints.Client
	Bot                               *bot.Client
	Cdn                               *cdn.Client
	Cognitive                         *cognitiveServices.Client
	Communication                     *communication.Client
	Compute                           *compute.Client
	ConfidentialLedger                *confidentialledger.Client
	Connections                       *connections.Client
	Consumption                       *consumption.Client
	ContainerApps                     *containerapps.Client
	Containers                        *containerServices.Client
	Cosmos                            *cosmosdb.Client
	CostManagement                    *costmanagement.Client
	CustomProviders                   *customproviders.Client
	Dashboard                         *dashboard.Client
	DatabaseMigration                 *datamigration.Client
	DataBricks                        *databricks.Client
	DataboxEdge                       *databoxedge.Client
	Datadog                           *datadog_v2021_03_01.Client
	DataFactory                       *datafactory.Client
	DataProtection                    *dataprotection.Client
	DataShare                         *datashare.Client
	DesktopVirtualization             *desktopvirtualization.Client
	DevTestLabs                       *devtestlabs.Client
	DigitalTwins                      *digitaltwins.Client
	Dns                               *dns_v2018_05_01.Client
	DomainServices                    *domainservices.Client
	Elastic                           *elastic.Client
	ElasticSan                        *elasticsan.Client
	EventGrid                         *eventgrid_v2022_06_15.Client
	Eventhub                          *eventhub.Client
	ExtendedLocation                  *extendedlocation.Client
	FluidRelay                        *fluidrelay_2022_05_26.Client
	Frontdoor                         *frontdoor.Client
	Graph                             *graph.Client
	HSM                               *hsm.Client
	HDInsight                         *hdinsight_v2021_06_01.Client
	HybridCompute                     *hybridcompute.Client
	HealthCare                        *healthcare.Client
	IoTCentral                        *iotcentral.Client
	IoTHub                            *iothub.Client
	KeyVault                          *keyvault.Client
	Kusto                             *kusto.Client
	Lighthouse                        *lighthouse.Client
	LoadBalancers                     *loadbalancers.Client
	LoadTestService                   *loadtestservice.AutoClient
	LogAnalytics                      *loganalytics.Client
	Logic                             *logic.Client
	MachineLearning                   *machinelearning.Client
	Maintenance                       *maintenance.Client
	ManagedApplication                *managedapplication.Client
	ManagementGroups                  *managementgroup.Client
	ManagedHSMs                       *managedhsm.Client
	Maps                              *maps.Client
	MixedReality                      *mixedreality.Client
	Monitor                           *monitor.Client
	MobileNetwork                     *mobilenetwork.Client
	MSSQL                             *mssql.Client
	MSSQLManagedInstance              *mssqlmanagedinstance.Client
	MySQL                             *mysql.Client
	NetApp                            *netapp.Client
	Network                           *network.Client
	NetworkFunction                   *networkfunction.Client
	NewRelic                          *newrelic.Client
	Nginx                             *nginx_2024_06_01_preview.Client
	NotificationHubs                  *notificationhub.Client
	Oracle                            *oracle.Client
	Orbital                           *orbital.Client
	PaloAlto                          *paloalto.Client
	Policy                            *policy.Client
	Portal                            *portal.Client
	Postgres                          *postgres.Client
	PowerBI                           *powerBI.Client
	PrivateDns                        *privatedns.Client
	PrivateDnsResolver                *dnsresolver.Client
	Purview                           *purview.Client
	RecoveryServices                  *recoveryServices.Client
	RedHatOpenShift                   *redhatopenshift.Client
	Redis                             *redis_2024_03_01.Client
	RedisEnterprise                   *redisenterprise.Client
	Relay                             *relay.Client
	Resource                          *resource.Client
	Search                            *search.Client
	SecurityCenter                    *securityCenter.Client
	Sentinel                          *sentinel.Client
	ServiceBus                        *serviceBus.Client
	ServiceConnector                  *serviceConnector.Client
	ServiceFabric                     *serviceFabric.Client
	ServiceFabricManaged              *serviceFabricManaged.Client
	ServiceNetworking                 *servicenetworking_2023_11_01.Client
	SignalR                           *signalr.Client
	Storage                           *storage.Client
	StorageCache                      *storagecache_2023_05_01.Client
	StorageMover                      *storageMover.Client
	StreamAnalytics                   *streamAnalytics.Client
	Subscription                      *subscription.Client
	Synapse                           *synapse.Client
	SystemCenterVirtualMachineManager *systemcentervirtualmachinemanager_2023_10_07.Client
	TrafficManager                    *trafficManager.Client
	VideoIndexer                      *videoindexer.Client
	Vmware                            *vmware.Client
	VoiceServices                     *voiceServices.Client
	Web                               *web.Client
	Workloads                         *workloads_v2023_04_01.Client
}

// NOTE: it should be possible for this method to become Private once the top level Client's removed

func (client *Client) Build(ctx context.Context, o *common.ClientOptions) error {
	autorest.Count429AsRetry = false
	// Disable the Azure SDK for Go's validation since it's unhelpful for our use-case
	validation.Disabled = true

	if err := buildAutoClients(&client.autoClient, o); err != nil {
		return fmt.Errorf("building auto-clients: %+v", err)
	}

	client.Features = o.Features
	client.StopContext = ctx

	var err error

	if client.AadB2c, err = aadb2c.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AadB2c: %+v", err)
	}
	if client.Advisor, err = advisor.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Advisor: %+v", err)
	}
	if client.AnalysisServices, err = analysisServices.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AnalysisServices: %+v", err)
	}
	if client.ApiManagement, err = apiManagement.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ApiManagement: %+v", err)
	}
	if client.AppConfiguration, err = appConfiguration.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AppConfiguration: %+v", err)
	}
	if client.AppInsights, err = applicationInsights.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ApplicationInsights: %+v", err)
	}
	if client.AppPlatform, err = appPlatform.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AppPlatform: %+v", err)
	}
	if client.AppService, err = appService.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AppService: %+v", err)
	}
	if client.ArcKubernetes, err = arckubernetes.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ArcKubernetes: %+v", err)
	}
	if client.ArcResourceBridge, err = arcResourceBridge.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Arc Resource Bridge: %+v", err)
	}
	if client.Attestation, err = attestation.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Attestation: %+v", err)
	}
	if client.Authorization, err = authorization.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Authorization: %+v", err)
	}
	if client.Automanage, err = automanage.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AutoManage: %+v", err)
	}
	if client.Automation, err = automation.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Automation: %+v", err)
	}
	if client.AzureStackHCI, err = azureStackHCI.NewClient(o); err != nil {
		return fmt.Errorf("building clients for AzureStackHCI: %+v", err)
	}
	if client.Batch, err = batch.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Batch: %+v", err)
	}
	if client.Blueprints, err = blueprints.NewClient(o); err != nil {
		return fmt.Errorf("building clients for BluePrints: %+v", err)
	}
	if client.Bot, err = bot.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Bot: %+v", err)
	}
	client.Cdn = cdn.NewClient(o)
	if client.Cognitive, err = cognitiveServices.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Cognitive: %+v", err)
	}
	if client.Communication, err = communication.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Communication: %+v", err)
	}
	if client.Compute, err = compute.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Compute: %+v", err)
	}
	if client.ConfidentialLedger, err = confidentialledger.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ConfidentialLedger: %+v", err)
	}
	if client.Connections, err = connections.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Connections: %+v", err)
	}
	if client.Consumption, err = consumption.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Consumption: %+v", err)
	}
	if client.Containers, err = containerServices.NewContainersClient(o); err != nil {
		return fmt.Errorf("building clients for Containers: %+v", err)
	}
	if client.ContainerApps, err = containerapps.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Container Apps: %+v", err)
	}
	if client.Cosmos, err = cosmosdb.NewClient(o); err != nil {
		return fmt.Errorf("building clients for CosmosDB: %+v", err)
	}
	if client.CostManagement, err = costmanagement.NewClient(o); err != nil {
		return fmt.Errorf("building clients for CostManagement: %+v", err)
	}
	if client.CustomProviders, err = customproviders.NewClient(o); err != nil {
		return fmt.Errorf("building clients for CustomProviders: %+v", err)
	}
	if client.Dashboard, err = dashboard.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Dashboard: %+v", err)
	}
	if client.DatabaseMigration, err = datamigration.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DatabaseMigration: %+v", err)
	}
	if client.DataBricks, err = databricks.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataBricks: %+v", err)
	}
	if client.DataboxEdge, err = databoxedge.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataboxEdge: %+v", err)
	}
	if client.Datadog, err = datadog.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Datadog: %+v", err)
	}
	if client.DataFactory, err = datafactory.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataFactory: %+v", err)
	}
	if client.DataProtection, err = dataprotection.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataProtection: %+v", err)
	}
	if client.DataShare, err = datashare.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DataShare: %+v", err)
	}
	if client.DesktopVirtualization, err = desktopvirtualization.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DesktopVirtualization: %+v", err)
	}
	if client.DevTestLabs, err = devtestlabs.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DevTestLabs: %+v", err)
	}
	if client.DigitalTwins, err = digitaltwins.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DigitalTwins: %+v", err)
	}
	if client.Dns, err = dns.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Dns: %+v", err)
	}
	if client.DomainServices, err = domainservices.NewClient(o); err != nil {
		return fmt.Errorf("building clients for DomainServices: %+v", err)
	}
	if client.Elastic, err = elastic.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Elastic: %+v", err)
	}
	if client.ElasticSan, err = elasticsan.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ElasticSan: %+v", err)
	}
	if client.EventGrid, err = eventgrid.NewClient(o); err != nil {
		return fmt.Errorf("building clients for EventGrid: %+v", err)
	}
	if client.Eventhub, err = eventhub.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Eventhub: %+v", err)
	}
	if client.ExtendedLocation, err = extendedlocation.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ExtendedLocation: %+v", err)
	}
	if client.FluidRelay, err = fluidrelay.NewClient(o); err != nil {
		return fmt.Errorf("building clients for FluidRelay: %+v", err)
	}
	client.Frontdoor = frontdoor.NewClient(o)
	if client.Graph, err = graph.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Graph: %+v", err)
	}
	if client.HSM, err = hsm.NewClient(o); err != nil {
		return fmt.Errorf("building clients for HSM: %+v", err)
	}
	if client.HDInsight, err = hdinsight.NewClient(o); err != nil {
		return fmt.Errorf("building clients for HDInsight: %+v", err)
	}
	if client.HealthCare, err = healthcare.NewClient(o); err != nil {
		return fmt.Errorf("building clients for HealthCare: %+v", err)
	}
	if client.HybridCompute, err = hybridcompute.NewClient(o); err != nil {
		return fmt.Errorf("building clients for HybridCompute: %+v", err)
	}
	if client.IoTCentral, err = iotcentral.NewClient(o); err != nil {
		return fmt.Errorf("building clients for IoTCentral: %+v", err)
	}
	if client.IoTHub, err = iothub.NewClient(o); err != nil {
		return fmt.Errorf("building clients for IoTHub: %+v", err)
	}
	if client.KeyVault, err = keyvault.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Key Vault: %+v", err)
	}
	if client.Kusto, err = kusto.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Kusto: %+v", err)
	}
	if client.Lighthouse, err = lighthouse.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Lighthouse: %+v", err)
	}
	if client.LogAnalytics, err = loganalytics.NewClient(o); err != nil {
		return fmt.Errorf("building clients for LogAnalytics: %+v", err)
	}
	if client.LoadBalancers, err = loadbalancers.NewClient(o); err != nil {
		return fmt.Errorf("building clients for LoadBalancers: %+v", err)
	}
	if client.LoadTestService, err = loadtestservice.NewClient(o); err != nil {
		return fmt.Errorf("building clients for LoadTestService: %+v", err)
	}
	if client.Logic, err = logic.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Logic: %+v", err)
	}
	if client.MachineLearning, err = machinelearning.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Machine Learning: %+v", err)
	}
	if client.Maintenance, err = maintenance.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Maintenance: %+v", err)
	}
	if client.ManagedApplication, err = managedapplication.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Managed Applications: %+v", err)
	}
	if client.ManagementGroups, err = managementgroup.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Management Groups: %+v", err)
	}
	if client.ManagedHSMs, err = managedhsm.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ManagedHSM: %+v", err)
	}
	if client.Maps, err = maps.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Maps: %+v", err)
	}
	if client.MixedReality, err = mixedreality.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Mixed Reality: %+v", err)
	}
	if client.Monitor, err = monitor.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Monitor: %+v", err)
	}
	if client.MobileNetwork, err = mobilenetwork.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Mobile Network: %+v", err)
	}
	if client.MSSQL, err = mssql.NewClient(o); err != nil {
		return fmt.Errorf("building clients for MSSQL: %+v", err)
	}
	if client.MSSQLManagedInstance, err = mssqlmanagedinstance.NewClient(o); err != nil {
		return fmt.Errorf("building clients for MSSQLManagedInstance: %+v", err)
	}
	if client.MySQL, err = mysql.NewClient(o); err != nil {
		return fmt.Errorf("building clients for MySQL: %+v", err)
	}
	if client.NetApp, err = netapp.NewClient(o); err != nil {
		return fmt.Errorf("building clients for NetApp: %+v", err)
	}
	if client.Network, err = network.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Network: %+v", err)
	}
	if client.NetworkFunction, err = networkfunction.NewClient(o); err != nil {
		return fmt.Errorf("building clients for NetworkFunction: %+v", err)
	}
	if client.NewRelic, err = newrelic.NewClient(o); err != nil {
		return fmt.Errorf("building clients for NewRelic: %+v", err)
	}
	if client.Nginx, err = nginx.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Nginx: %+v", err)
	}
	if client.NotificationHubs, err = notificationhub.NewClient(o); err != nil {
		return fmt.Errorf("building clients for NotificationHubs: %+v", err)
	}
	if client.Oracle, err = oracle.NewClient(o); err != nil {
		return fmt.Errorf("building clients for OracleDatabase: %+v", err)
	}
	if client.Orbital, err = orbital.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Orbital: %+v", err)
	}
	if client.Policy, err = policy.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Policy: %+v", err)
	}
	if client.PaloAlto, err = paloalto.NewClient(o); err != nil {
		return fmt.Errorf("building clients for PaloAlto: %+v", err)
	}
	if client.Portal, err = portal.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Portal: %+v", err)
	}
	if client.Postgres, err = postgres.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Postgres: %+v", err)
	}
	if client.PowerBI, err = powerBI.NewClient(o); err != nil {
		return fmt.Errorf("building clients for PowerBI: %+v", err)
	}
	if client.PrivateDns, err = privatedns.NewClient(o); err != nil {
		return fmt.Errorf("building clients for PrivateDns: %+v", err)
	}
	if client.PrivateDnsResolver, err = dnsresolver.NewClient(o); err != nil {
		return fmt.Errorf("building clients for PrivateDnsResolver: %+v", err)
	}
	if client.Purview, err = purview.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Purview: %+v", err)
	}
	if client.RecoveryServices, err = recoveryServices.NewClient(o); err != nil {
		return fmt.Errorf("building clients for RecoveryServices: %+v", err)
	}
	if client.RedHatOpenShift, err = redhatopenshift.NewClient(o); err != nil {
		return fmt.Errorf("building clients for RedHatOpenShift: %+v", err)
	}
	if client.Redis, err = redis.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Redis: %+v", err)
	}
	if client.RedisEnterprise, err = redisenterprise.NewClient(o); err != nil {
		return fmt.Errorf("building clients for RedisEnterprise: %+v", err)
	}
	if client.Relay, err = relay.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Relay: %+v", err)
	}
	if client.Resource, err = resource.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Resource: %+v", err)
	}
	if client.Search, err = search.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Search: %+v", err)
	}
	if client.SecurityCenter, err = securityCenter.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Security Center: %+v", err)
	}
	if client.Sentinel, err = sentinel.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Sentinel: %+v", err)
	}
	if client.ServiceBus, err = serviceBus.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ServiceBus: %+v", err)
	}
	if client.ServiceConnector, err = serviceConnector.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ServiceConnector: %+v", err)
	}
	if client.ServiceFabric, err = serviceFabric.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ServiceConnector: %+v", err)
	}
	if client.ServiceFabricManaged, err = serviceFabricManaged.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ServiceFabric: %+v", err)
	}
	if client.ServiceNetworking, err = serviceNetworking.NewClient(o); err != nil {
		return fmt.Errorf("building clients for ServiceNetworking: %+v", err)
	}
	if client.SignalR, err = signalr.NewClient(o); err != nil {
		return fmt.Errorf("building clients for SignalR: %+v", err)
	}
	if client.Storage, err = storage.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Storage: %+v", err)
	}
	if client.StorageCache, err = storageCache.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Storage Cache: %+v", err)
	}
	if client.StorageMover, err = storageMover.NewClient(o); err != nil {
		return fmt.Errorf("building clients for StorageMover: %+v", err)
	}
	if client.StreamAnalytics, err = streamAnalytics.NewClient(o); err != nil {
		return fmt.Errorf("building clients for StreamAnalytics: %+v", err)
	}
	if client.Subscription, err = subscription.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Subscription: %+v", err)
	}

	client.Synapse = synapse.NewClient(o)
	if client.SystemCenterVirtualMachineManager, err = systemCenterVirtualMachineManager.NewClient(o); err != nil {
		return fmt.Errorf("building clients for System Center Virtual Machine Manager: %+v", err)
	}
	if client.TrafficManager, err = trafficManager.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Traffic Manager: %+v", err)
	}

	if client.VideoIndexer, err = videoindexer.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Video Indexer: %+v", err)
	}

	if client.Vmware, err = vmware.NewClient(o); err != nil {
		return fmt.Errorf("building clients for VMWare: %+v", err)
	}
	if client.VoiceServices, err = voiceServices.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Voice Services: %+v", err)
	}
	client.Web = web.NewClient(o)

	if client.Workloads, err = workloads.NewClient(o); err != nil {
		return fmt.Errorf("building clients for Workloads: %+v", err)
	}

	return nil
}
