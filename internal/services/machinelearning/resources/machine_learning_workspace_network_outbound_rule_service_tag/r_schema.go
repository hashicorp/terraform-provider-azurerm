// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package machine_learning_workspace_network_outbound_rule_service_tag

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2025-06-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func (r Resource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"workspace_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: managednetwork.ValidateWorkspaceID,
		},

		"service_tag": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ValidateFunc: validation.StringInSlice([]string{
				"AppConfiguration",
				"AppService",
				"AzureActiveDirectory",
				"AzureAdvancedThreatProtection",
				"AzureArcInfrastructure",
				"AzureAttestation",
				"AzureBackup",
				"AzureBotService",
				"AzureContainerRegistry",
				"AzureCosmosDB",
				"AzureDataLake",
				"AzureDevSpaces",
				"AzureInformationProtection",
				"AzureIoTHub",
				"AzureKeyVault",
				"AzureManagedGrafana",
				"AzureMonitor",
				"AzureOpenDatasets",
				"AzurePlatformDNS",
				"AzurePlatformIMDS",
				"AzurePlatformLKM",
				"AzureResourceManager",
				"AzureSignalR",
				"AzureSiteRecovery",
				"AzureSpringCloud",
				"AzureStack",
				"AzureUpdateDelivery",
				"DataFactoryManagement",
				"EventHub",
				"GuestAndHybridManagement",
				"M365ManagementActivityApi",
				"M365ManagementActivityApi",
				"MicrosoftAzureFluidRelay",
				"MicrosoftCloudAppSecurity",
				"MicrosoftContainerRegistry",
				"PowerPlatformInfra",
				"ServiceBus",
				"Sql",
				"Storage",
				"WindowsAdminCenter",
				"AppServiceManagement",
				"AutonomousDevelopmentPlatform",
				"AzureActiveDirectoryDomainServices",
				"AzureCloud",
				"AzureConnectors",
				"AzureContainerAppsService",
				"AzureDatabricks",
				"AzureDeviceUpdate",
				"AzureEventGrid",
				"AzureFrontDoor.Frontend",
				"AzureFrontDoor.Backend",
				"AzureFrontDoor.FirstParty",
				"AzureHealthcareAPIs",
				"AzureLoadBalancer",
				"AzureMachineLearning",
				"AzureSphere",
				"AzureWebPubSub",
				"BatchNodeManagement",
				"ChaosStudio",
				"CognitiveServicesFrontend",
				"CognitiveServicesManagement",
				"DataFactory",
				"Dynamics365ForMarketingEmail",
				"Dynamics365BusinessCentral",
				"EOPExternalPublishedIPs",
				"Internet",
				"LogicApps",
				"Marketplace",
				"MicrosoftDefenderForEndpoint",
				"PowerBI",
				"PowerQueryOnline",
				"ServiceFabric",
				"SqlManagement",
				"StorageSyncService",
				"WindowsVirtualDesktop",
				"VirtualNetwork",
			}, false),
		},

		"protocol": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"*", "TCP", "UDP", "ICMP"}, false),
		},

		"port_ranges": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (r Resource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
