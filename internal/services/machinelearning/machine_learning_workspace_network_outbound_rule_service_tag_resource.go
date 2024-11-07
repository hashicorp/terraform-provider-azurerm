// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package machinelearning

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/machinelearningservices/2024-04-01/managednetwork"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.ResourceWithUpdate = WorkspaceNetworkOutboundRuleServiceTag{}

type machineLearningWorkspaceServiceTagOutboundRuleModel struct {
	Name        string `tfschema:"name"`
	WorkspaceId string `tfschema:"workspace_id"`
	ServiceTag  string `tfschema:"service_tag"`
	Protocol    string `tfschema:"protocol"`
	PortRanges  string `tfschema:"port_ranges"`
}

type WorkspaceNetworkOutboundRuleServiceTag struct{}

var _ sdk.Resource = WorkspaceNetworkOutboundRuleServiceTag{}

func (r WorkspaceNetworkOutboundRuleServiceTag) ResourceType() string {
	return "azurerm_machine_learning_workspace_network_outbound_rule_service_tag"
}

func (r WorkspaceNetworkOutboundRuleServiceTag) ModelObject() interface{} {
	return &machineLearningWorkspaceServiceTagOutboundRuleModel{}
}

func (r WorkspaceNetworkOutboundRuleServiceTag) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return managednetwork.ValidateOutboundRuleID
}

func (r WorkspaceNetworkOutboundRuleServiceTag) Arguments() map[string]*pluginsdk.Schema {
	arguments := map[string]*pluginsdk.Schema{
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
	return arguments
}

func (r WorkspaceNetworkOutboundRuleServiceTag) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r WorkspaceNetworkOutboundRuleServiceTag) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceServiceTagOutboundRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MachineLearning.ManagedNetwork
			subscriptionId := metadata.Client.Account.SubscriptionId

			workspaceId, err := managednetwork.ParseWorkspaceID(model.WorkspaceId)
			if err != nil {
				return err
			}
			id := managednetwork.NewOutboundRuleID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, model.Name)
			existing, err := client.SettingsRuleGet(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_machine_learning_workspace_network_outbound_rule_service_tag", id.ID())
			}

			outboundRule := managednetwork.OutboundRuleBasicResource{
				Name: pointer.To(model.Name),
				Type: pointer.To(string(managednetwork.RuleTypeServiceTag)),
				Properties: managednetwork.ServiceTagOutboundRule{
					Category: pointer.To(managednetwork.RuleCategoryUserDefined),
					Destination: &managednetwork.ServiceTagDestination{
						PortRanges: &model.PortRanges,
						Protocol:   &model.Protocol,
						ServiceTag: &model.ServiceTag,
					},
				},
			}

			if err = client.SettingsRuleCreateOrUpdateThenPoll(ctx, id, outboundRule); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r WorkspaceNetworkOutboundRuleServiceTag) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model machineLearningWorkspaceServiceTagOutboundRuleModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.MachineLearning.ManagedNetwork
			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			existing, err := client.SettingsRuleGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}

			payload := existing.Model

			serviceTagOutboundRule := managednetwork.ServiceTagOutboundRule{
				Type:        managednetwork.RuleTypeServiceTag,
				Category:    pointer.To(managednetwork.RuleCategoryUserDefined),
				Destination: &managednetwork.ServiceTagDestination{},
			}

			if metadata.ResourceData.HasChange("service_tag") {
				serviceTagOutboundRule.Destination.ServiceTag = pointer.To(model.ServiceTag)
			}

			if metadata.ResourceData.HasChange("protocol") {
				serviceTagOutboundRule.Destination.Protocol = pointer.To(model.Protocol)
			}

			if metadata.ResourceData.HasChange("port_ranges") {
				serviceTagOutboundRule.Destination.PortRanges = pointer.To(model.PortRanges)
			}

			payload.Properties = serviceTagOutboundRule
			if err := client.SettingsRuleCreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}
			return nil
		},
	}
}

func (r WorkspaceNetworkOutboundRuleServiceTag) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.SettingsRuleGet(ctx, *id)
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

			state := machineLearningWorkspaceServiceTagOutboundRuleModel{
				Name: id.OutboundRuleName,
			}

			if props := model.Properties; props != nil {
				if prop, ok := props.(managednetwork.ServiceTagOutboundRule); ok && prop.Destination != nil {
					if prop.Destination.ServiceTag != nil {
						state.ServiceTag = *prop.Destination.ServiceTag
					}

					if prop.Destination.Protocol != nil {
						state.Protocol = *prop.Destination.Protocol
					}

					if prop.Destination.PortRanges != nil {
						state.PortRanges = *prop.Destination.PortRanges
					}
				}
			}
			state.WorkspaceId = managednetwork.NewWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName).ID()
			return metadata.Encode(&state)
		},
	}
}

func (r WorkspaceNetworkOutboundRuleServiceTag) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MachineLearning.ManagedNetwork

			id, err := managednetwork.ParseOutboundRuleID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.SettingsRuleDelete(ctx, *id)
			if err != nil {
				return fmt.Errorf("deleting Machine Learning Workspace Service Tag Network Outbound Rule %q (Resource Group %q, Workspace %q): %+v", id.OutboundRuleName, id.ResourceGroupName, id.WorkspaceName, err)
			}

			if err = future.Poller.PollUntilDone(ctx); err != nil {
				return fmt.Errorf("waiting for deletion of Machine Learning Workspace Service Tag Network Outbound Rule %q (Resource Group %q, Workspace %q): %+v", id.OutboundRuleName, id.ResourceGroupName, id.WorkspaceName, err)
			}

			return nil
		},
	}
}
