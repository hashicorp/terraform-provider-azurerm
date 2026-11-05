// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_kubernetes_provisioned_cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ArcKubernetesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ArcKubernetesProvisionedClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := connectedclusters.NewConnectedClusterID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.ConnectedClusterGet(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			arcAgentAutoUpgrade := connectedclusters.AutoUpgradeOptionsDisabled
			if model.ArcAgentAutoUpgradeEnabled {
				arcAgentAutoUpgrade = connectedclusters.AutoUpgradeOptionsEnabled
			}

			expandedIdentity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			payload := &connectedclusters.ConnectedCluster{
				Identity: pointer.From(expandedIdentity),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Kind:     pointer.To(connectedclusters.ConnectedClusterKindProvisionedCluster),
				Properties: connectedclusters.ConnectedClusterProperties{
					ArcAgentProfile: &connectedclusters.ArcAgentProfile{
						AgentAutoUpgrade: pointer.To(arcAgentAutoUpgrade),
					},
				},
			}

			if aadProfileVal := model.AzureActiveDirectory; len(aadProfileVal) != 0 {
				payload.Properties.AadProfile = expandArcKubernetesClusterAadProfile(aadProfileVal)
			}

			if desiredVersion := model.ArcAgentDesiredVersion; desiredVersion != "" {
				payload.Properties.ArcAgentProfile.DesiredAgentVersion = pointer.To(desiredVersion)
			}

			if err := client.ConnectedClusterCreateCallbackThenPoll(ctx, id, *payload, metadata.SetIDAndIdentityCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			if err := pluginsdk.SetResourceIdentityData(metadata.ResourceData, &id); err != nil {
				return err
			}

			return nil
		},
	}
}

func expandArcKubernetesClusterAadProfile(input []AzureActiveDirectoryModel) *connectedclusters.AadProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := connectedclusters.AadProfile{
		EnableAzureRBAC: pointer.To(v.AzureRbacEnabled),
	}

	if tenantIdVal := v.TenantId; tenantIdVal != "" {
		output.TenantID = pointer.To(tenantIdVal)
	}

	if groupVal := v.AdminGroupObjectIds; len(groupVal) != 0 {
		output.AdminGroupObjectIDs = pointer.To(groupVal)
	}

	return &output
}
