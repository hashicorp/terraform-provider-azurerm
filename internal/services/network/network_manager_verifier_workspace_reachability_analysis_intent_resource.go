// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cosmosdb/2024-08-15/cosmosdb"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/reachabilityanalysisintent"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-01-01/reachabilityanalysisintents"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var _ sdk.Resource = ManagerVerifierWorkspaceReachabilityAnalysisIntentResource{}

type ManagerVerifierWorkspaceReachabilityAnalysisIntentResource struct{}

func (ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return reachabilityanalysisintent.ValidateReachabilityAnalysisIntentID
}

func (ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) ResourceType() string {
	return "azurerm_network_manager_verifier_workspace_reachability_analysis_intent"
}

func (ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) ModelObject() interface{} {
	return &ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel{}
}

type ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel struct {
	Description           string                                                        `tfschema:"description"`
	DestinationResourceId string                                                        `tfschema:"destination_resource_id"`
	IpTraffic             []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic `tfschema:"ip_traffic"`
	Name                  string                                                        `tfschema:"name"`
	SourceResourceId      string                                                        `tfschema:"source_resource_id"`
	VerifierWorkspaceId   string                                                        `tfschema:"verifier_workspace_id"`
}

type ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic struct {
	DestinationIps   []string                                      `tfschema:"destination_ips"`
	DestinationPorts []string                                      `tfschema:"destination_ports"`
	Protocols        []reachabilityanalysisintents.NetworkProtocol `tfschema:"protocols"`
	SourceIps        []string                                      `tfschema:"source_ips"`
	SourcePorts      []string                                      `tfschema:"source_ports"`
}

func (ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile(`^[a-zA-Z0-9\_\.\-]{1,64}$`),
				"`name` must be between 1 and 64 characters long and can only contain letters, numbers, underscores(_), periods(.), and hyphens(-).",
			),
		},

		"verifier_workspace_id": commonschema.ResourceIDReferenceRequiredForceNew(&reachabilityanalysisintents.VerifierWorkspaceId{}),

		"source_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				commonids.ValidatePublicIPAddressID,
				commonids.ValidateSubnetID,
				commonids.ValidateVirtualMachineID,
			),
		},

		"destination_resource_id": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.Any(
				commonids.ValidatePublicIPAddressID,
				commonids.ValidateSqlServerID,
				commonids.ValidateStorageAccountID,
				commonids.ValidateSubnetID,
				commonids.ValidateVirtualMachineID,
				cosmosdb.ValidateDatabaseAccountID,
			),
		},

		"ip_traffic": {
			Type:     pluginsdk.TypeList,
			Required: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"destination_ips": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsCIDR,
								validation.IsIPAddress,
							),
						},
					},
					"destination_ports": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.IpTrafficPort,
						},
					},
					"source_ips": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
							ValidateFunc: validation.Any(
								validation.IsCIDR,
								validation.IsIPAddress,
							),
						},
					},
					"source_ports": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validate.IpTrafficPort,
						},
					},
					"protocols": {
						Type:     pluginsdk.TypeList,
						Required: true,
						ForceNew: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringInSlice(reachabilityanalysisintents.PossibleValuesForNetworkProtocol(), false),
						},
					},
				},
			},
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	}
}

func (ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ReachabilityAnalysisIntents
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			workspaceId, err := reachabilityanalysisintents.ParseVerifierWorkspaceID(config.VerifierWorkspaceId)
			if err != nil {
				return err
			}

			id := reachabilityanalysisintents.NewReachabilityAnalysisIntentID(subscriptionId, workspaceId.ResourceGroupName, workspaceId.NetworkManagerName, workspaceId.VerifierWorkspaceName, config.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			payload := reachabilityanalysisintents.ReachabilityAnalysisIntent{
				Name: pointer.To(config.Name),
				Properties: reachabilityanalysisintents.ReachabilityAnalysisIntentProperties{
					Description:           pointer.To(config.Description),
					SourceResourceId:      config.SourceResourceId,
					DestinationResourceId: config.DestinationResourceId,
					IPTraffic:             expandReachabilityAnalysisIntentIPTraffic(config.IpTraffic),
				},
			}

			if _, err := client.Create(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ReachabilityAnalysisIntents

			id, err := reachabilityanalysisintents.ParseReachabilityAnalysisIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id, err)
			}

			workspaceId := reachabilityanalysisintents.NewVerifierWorkspaceID(id.SubscriptionId, id.ResourceGroupName, id.NetworkManagerName, id.VerifierWorkspaceName).ID()
			schema := ManagerVerifierWorkspaceReachabilityAnalysisIntentResourceModel{
				Name:                id.VerifierWorkspaceName,
				VerifierWorkspaceId: workspaceId,
			}

			if model := resp.Model; model != nil {
				props := model.Properties
				schema.Description = pointer.From(props.Description)
				schema.DestinationResourceId = props.DestinationResourceId
				schema.SourceResourceId = props.SourceResourceId
				schema.IpTraffic = flattenReachabilityAnalysisIntentIPTraffic(props.IPTraffic)
			}

			return metadata.Encode(&schema)
		},
	}
}

func (r ManagerVerifierWorkspaceReachabilityAnalysisIntentResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Network.ReachabilityAnalysisIntent

			id, err := reachabilityanalysisintent.ParseReachabilityAnalysisIntentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandReachabilityAnalysisIntentIPTraffic(input []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic) reachabilityanalysisintents.IPTraffic {
	if len(input) == 0 {
		return reachabilityanalysisintents.IPTraffic{}
	}

	return reachabilityanalysisintents.IPTraffic{
		DestinationIPs:   input[0].DestinationIps,
		DestinationPorts: input[0].DestinationPorts,
		Protocols:        input[0].Protocols,
		SourceIPs:        input[0].SourceIps,
		SourcePorts:      input[0].SourcePorts,
	}
}

func flattenReachabilityAnalysisIntentIPTraffic(input reachabilityanalysisintents.IPTraffic) []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic {
	return []ManagerVerifierWorkspaceReachabilityAnalysisIntentIpTraffic{
		{
			DestinationIps:   input.DestinationIPs,
			DestinationPorts: input.DestinationPorts,
			Protocols:        input.Protocols,
			SourceIps:        input.SourceIPs,
			SourcePorts:      input.SourcePorts,
		},
	}
}
