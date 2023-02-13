package environments

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type Authorization struct {
	Audiences     []string
	LoginEndpoint string
}

type Environment struct {
	Name string

	Authorization   *Authorization
	MicrosoftGraph  Api
	ResourceManager Api

	AnalysisServices                  Api
	ApiManagement                     Api
	AppConfiguration                  Api
	ApplicationInsights               Api
	AppService                        Api
	Attestation                       Api
	AzureADIdentityGovernanceInsights Api
	AzureADIntegratedApp              Api
	AzureADNotification               Api
	AzureCLI                          Api
	AzureDevOps                       Api
	AzureVPN                          Api
	Batch                             Api
	Bing                              Api
	BotFrameworkDevPortal             Api
	BranchConnectWebService           Api
	CDNFrontDoor                      Api
	Cognitive                         Api
	ComputeRecommendations            Api
	ContainerRegistry                 Api
	Connections                       Api
	CortanaAtWork                     Api
	CortanaAtWorkBing                 Api
	CortanaRuntime                    Api
	CosmosDB                          Api
	CustomerInsights                  Api
	DataBricks                        Api
	DataCatalog                       Api
	DataLake                          Api
	DataMigrations                    Api
	DigitalTwins                      Api
	DomainController                  Api
	Dynamic365BusinessCentral         Api
	Dynamics365DataExportService      Api
	DynamicsCRM                       Api
	DynamicsERP                       Api
	EventHubs                         Api
	Flow                              Api
	GraphConnector                    Api
	HDInsight                         Api
	HealthCare                        Api
	IamSupportability                 Api
	ImportExport                      Api
	InformationProtectionSyncService  Api
	InTune                            Api
	IoTCentral                        Api
	IoTHubDeviceProvisioning          Api
	KeyVault                          Api
	KubernetesServiceAADServer        Api
	Kusto                             Api
	KustoMFA                          Api
	LogAnalytics                      Api
	ManagedHSM                        Api
	Maps                              Api
	MariaDB                           Api
	MediaServices                     Api
	Microsoft365DataAtRestEncryption  Api
	MicrosoftInvoicing                Api
	MileIqAdminCenter                 Api
	MileIqDashboard                   Api
	MileIqRestService                 Api
	MixedReality                      Api
	MySql                             Api
	Office365Connectors               Api
	Office365Demeter                  Api
	Office365DwEngineV2               Api
	Office365ExchangeOnline           Api
	Office365ExchangeOnlineProtection Api
	Office365InformationProtection    Api
	Office365Management               Api
	Office365Portal                   Api
	Office365SharePointOnline         Api
	Office365Zoom                     Api
	OneNote                           Api
	OneProfile                        Api
	OperationalInsights               Api
	OSSRDMBS                          Api
	PeopleCards                       Api
	PolicyAdministration              Api
	Portal                            Api
	Postgresql                        Api
	PowerAppsRuntime                  Api
	PowerAppsRuntimeService           Api
	PowerBiService                    Api
	Purview                           Api
	RightsManagement                  Api
	SecurityInsights                  Api
	ServiceBus                        Api
	ServiceDeploy                     Api
	ServiceTrust                      Api
	SkypeForBusinessOnline            Api
	Signup                            Api
	SpeechRecognition                 Api
	Sql                               Api
	StackHCI                          Api
	StreamAnalytics                   Api
	Storage                           Api
	StorageSync                       Api
	Synapse                           Api
	SynapseGateway                    Api
	SynapseStudio                     Api
	TargetedMessaging                 Api
	TimeSeriesInsights                Api
	Teams                             Api
	ThreatProtection                  Api
	TrafficManager                    Api
	UniversalPrint                    Api
	WindowsDefenderATP                Api
	WindowsVirtualDesktop             Api
	Yammer                            Api
}

var _ Api = &ApiEndpoint{}

type ApiEndpoint struct {
	domainSuffix        *string
	endpoint            *string
	microsoftGraphAppId *string
	name                string
	resourceIdentifier  *string
}

func NewApiEndpoint(name, endpoint string, microsoftGraphAppId *string) *ApiEndpoint {
	return &ApiEndpoint{
		endpoint:            pointer.To(endpoint),
		microsoftGraphAppId: microsoftGraphAppId,
		name:                name,
	}
}

func (e *ApiEndpoint) withResourceIdentifier(identifier string) *ApiEndpoint {
	e.resourceIdentifier = pointer.To(identifier)
	return e
}

func (e *ApiEndpoint) DomainSuffix() (*string, bool) {
	if e == nil {
		return nil, false
	}
	return e.domainSuffix, e.domainSuffix != nil
}

func (e *ApiEndpoint) Endpoint() (*string, bool) {
	if e == nil {
		return nil, false
	}
	return e.endpoint, e.endpoint != nil
}

func (e *ApiEndpoint) MicrosoftGraphAppId() (*string, bool) {
	if e == nil || e.microsoftGraphAppId == nil {
		return nil, false
	}
	return e.microsoftGraphAppId, true
}

func (e *ApiEndpoint) Name() string {
	return e.name
}

func (e *ApiEndpoint) ResourceIdentifier() (*string, bool) {
	if e == nil {
		return nil, false
	}
	return e.resourceIdentifier, e.resourceIdentifier != nil
}
