// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package environments

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
)

type Authorization struct {
	Audiences        []string
	IdentityProvider string
	LoginEndpoint    string
	Tenant           string
}

type Environment struct {
	Name string

	Authorization   *Authorization
	MicrosoftGraph  Api
	ResourceManager Api

	AnalysisServices                                  Api
	ApiManagement                                     Api
	AppConfiguration                                  Api
	AppService                                        Api
	ApplicationInsights                               Api
	Attestation                                       Api
	AzureADIdentityGovernanceInsights                 Api
	AzureADIntegratedApp                              Api
	AzureADNotification                               Api
	AzureDevOps                                       Api
	AzureServiceManagement                            Api
	AzureVPN                                          Api
	Batch                                             Api
	Bing                                              Api
	BotFrameworkDevPortal                             Api
	BranchConnectWebService                           Api
	CDNFrontDoor                                      Api
	Cognitive                                         Api
	ComputeRecommendations                            Api
	Connections                                       Api
	ContainerRegistry                                 Api
	CortanaAtWork                                     Api
	CortanaAtWorkBing                                 Api
	CortanaRuntime                                    Api
	CosmosDB                                          Api
	CustomerInsights                                  Api
	DataBricks                                        Api
	DataCatalog                                       Api
	DataLake                                          Api
	DataMigrations                                    Api
	DigitalTwins                                      Api
	DomainController                                  Api
	Dynamic365BusinessCentral                         Api
	Dynamics365DataExportService                      Api
	DynamicsCRM                                       Api
	DynamicsERP                                       Api
	EventHubs                                         Api
	Flow                                              Api
	GraphConnector                                    Api
	HDInsight                                         Api
	HealthCare                                        Api
	IamSupportability                                 Api
	ImportExport                                      Api
	InTune                                            Api
	InformationProtectionSyncService                  Api
	IoTCentral                                        Api
	IoTHubDeviceProvisioning                          Api
	KeyVault                                          Api
	KubernetesServiceAADServer                        Api
	Kusto                                             Api
	KustoMFA                                          Api
	LogAnalytics                                      Api
	ManagedHSM                                        Api
	Maps                                              Api
	MariaDB                                           Api
	MediaServices                                     Api
	Microsoft365DataAtRestEncryption                  Api
	MicrosoftAzureCli                                 Api
	MicrosoftInvoicing                                Api
	MicrosoftOffice                                   Api
	MicrosoftStorageSync                              Api
	MicrosoftTeams                                    Api
	MicrosoftTeamsWebClient                           Api
	MileIqAdminCenter                                 Api
	MileIqDashboard                                   Api
	MileIqRestService                                 Api
	MixedReality                                      Api
	MySql                                             Api
	OSSRDBMSPostgreSQLFlexibleServerAadAuthentication Api
	OSSRDMBS                                          Api
	Office365Connectors                               Api
	Office365Demeter                                  Api
	Office365DwEngineV2                               Api
	Office365ExchangeOnline                           Api
	Office365ExchangeOnlineProtection                 Api
	Office365InformationProtection                    Api
	Office365Management                               Api
	Office365Portal                                   Api
	Office365SharePointOnline                         Api
	Office365SuiteUx                                  Api
	Office365Zoom                                     Api
	OfficeHome                                        Api
	OfficeUwpPwa                                      Api
	OneNote                                           Api
	OneProfile                                        Api
	OperationalInsights                               Api
	PeopleCards                                       Api
	PolicyAdministration                              Api
	Portal                                            Api
	Postgresql                                        Api
	PowerAppsRuntime                                  Api
	PowerAppsRuntimeService                           Api
	PowerBiService                                    Api
	Purview                                           Api
	RightsManagement                                  Api
	SecurityInsights                                  Api
	ServiceBus                                        Api
	ServiceDeploy                                     Api
	ServiceTrust                                      Api
	Signup                                            Api
	SkypeForBusinessOnline                            Api
	SpeechRecognition                                 Api
	Sql                                               Api
	StackHCI                                          Api
	Storage                                           Api
	StorageSync                                       Api
	StreamAnalytics                                   Api
	Synapse                                           Api
	SynapseGateway                                    Api
	SynapseStudio                                     Api
	TargetedMessaging                                 Api
	Teams                                             Api
	ThreatProtection                                  Api
	TimeSeriesInsights                                Api
	TrafficManager                                    Api
	UniversalPrint                                    Api
	WindowsDefenderATP                                Api
	WindowsVirtualDesktop                             Api
	Yammer                                            Api
}

var _ Api = &ApiEndpoint{}

type ApiEndpoint struct {
	domainSuffix       *string
	endpoint           *string
	appId              *string
	name               string
	resourceIdentifier *string
}

func NewApiEndpoint(name, endpoint string, appId *string) *ApiEndpoint {
	return &ApiEndpoint{
		appId:    appId,
		endpoint: pointer.To(endpoint),
		name:     name,
	}
}

func (e *ApiEndpoint) WithResourceIdentifier(identifier string) Api {
	newApi := *e
	newApi.resourceIdentifier = pointer.To(identifier)
	return &newApi
}

func (e *ApiEndpoint) Available() bool {
	return e != nil && (e.resourceIdentifier != nil || e.endpoint != nil)
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

func (e *ApiEndpoint) AppId() (*string, bool) {
	if e == nil || e.appId == nil {
		return nil, false
	}
	return e.appId, true
}

func (e *ApiEndpoint) Name() string {
	if e == nil {
		return "(nil)"
	}
	return e.name
}

func (e *ApiEndpoint) ResourceIdentifier() (*string, bool) {
	if e == nil {
		return nil, false
	}
	return e.resourceIdentifier, e.resourceIdentifier != nil
}
