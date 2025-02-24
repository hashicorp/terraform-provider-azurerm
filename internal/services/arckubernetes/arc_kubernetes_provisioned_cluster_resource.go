// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package arckubernetes

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	arckubernetes "github.com/hashicorp/go-azure-sdk/resource-manager/hybridkubernetes/2024-01-01/connectedclusters"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

var (
	_ sdk.Resource           = ArcKubernetesProvisionedClusterResource{}
	_ sdk.ResourceWithUpdate = ArcKubernetesProvisionedClusterResource{}
)

// This resource is same type as the ArcKubernetesClusterResource but with kind="ProvisionedCluster".
type ArcKubernetesProvisionedClusterResource struct{}

type ArcKubernetesProvisionedClusterModel struct {
	AgentVersion               string                         `tfschema:"agent_version"`
	ArcAgentAutoUpgradeEnabled bool                           `tfschema:"arc_agent_auto_upgrade_enabled"`
	ArcAgentDesiredVersion     string                         `tfschema:"arc_agent_desired_version"`
	AzureActiveDirectory       []AzureActiveDirectoryModel    `tfschema:"azure_active_directory"`
	Distribution               string                         `tfschema:"distribution"`
	Identity                   []identity.ModelSystemAssigned `tfschema:"identity"`
	Infrastructure             string                         `tfschema:"infrastructure"`
	KubernetesVersion          string                         `tfschema:"kubernetes_version"`
	Location                   string                         `tfschema:"location"`
	Name                       string                         `tfschema:"name"`
	Offering                   string                         `tfschema:"offering"`
	ResourceGroupName          string                         `tfschema:"resource_group_name"`
	Tags                       map[string]string              `tfschema:"tags"`
	TotalCoreCount             int64                          `tfschema:"total_core_count"`
	TotalNodeCount             int64                          `tfschema:"total_node_count"`
}

type AzureActiveDirectoryModel struct {
	AdminGroupObjectIds []string `tfschema:"admin_group_object_ids"`
	AzureRbacEnabled    bool     `tfschema:"azure_rbac_enabled"`
	TenantId            string   `tfschema:"tenant_id"`
}

func (r ArcKubernetesProvisionedClusterResource) ResourceType() string {
	return "azurerm_arc_kubernetes_provisioned_cluster"
}

func (r ArcKubernetesProvisionedClusterResource) ModelObject() interface{} {
	return &ArcKubernetesProvisionedClusterModel{}
}

func (r ArcKubernetesProvisionedClusterResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return arckubernetes.ValidateConnectedClusterID
}

func (r ArcKubernetesProvisionedClusterResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[-_a-zA-Z0-9]{1,260}$"),
				"The name of Arc Kubernetes Provisioned Cluster can only include alphanumeric characters, underscores, hyphens, has a maximum length of 260 characters, and must be unique.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"identity": commonschema.SystemAssignedIdentityRequiredForceNew(),

		"azure_active_directory": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"admin_group_object_ids": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsUUID,
						},
					},

					"azure_rbac_enabled": {
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  false,
					},

					"tenant_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.IsUUID,
					},
				},
			},
		},

		"arc_agent_desired_version": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"arc_agent_auto_upgrade_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  true,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ArcKubernetesProvisionedClusterResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"agent_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"distribution": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"infrastructure": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"kubernetes_version": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"offering": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"total_core_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},

		"total_node_count": {
			Type:     pluginsdk.TypeInt,
			Computed: true,
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ArcKubernetesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ArcKubernetesProvisionedClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := arckubernetes.NewConnectedClusterID(subscriptionId, model.ResourceGroupName, model.Name)

			existing, err := client.ConnectedClusterGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			arcAgentAutoUpgrade := arckubernetes.AutoUpgradeOptionsDisabled
			if model.ArcAgentAutoUpgradeEnabled {
				arcAgentAutoUpgrade = arckubernetes.AutoUpgradeOptionsEnabled
			}

			expandedIdentity, err := identity.ExpandSystemAssignedFromModel(model.Identity)
			if err != nil {
				return fmt.Errorf("expanding `identity`: %+v", err)
			}

			payload := &arckubernetes.ConnectedCluster{
				Identity: pointer.From(expandedIdentity),
				Location: location.Normalize(model.Location),
				Tags:     pointer.To(model.Tags),
				Kind:     pointer.To(arckubernetes.ConnectedClusterKindProvisionedCluster),
				Properties: arckubernetes.ConnectedClusterProperties{
					ArcAgentProfile: &arckubernetes.ArcAgentProfile{
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

			if err := client.ConnectedClusterCreateThenPoll(ctx, id, *payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)

			return nil
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ArcKubernetesClient

			id, err := arckubernetes.ParseConnectedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ArcKubernetesProvisionedClusterModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.ConnectedClusterGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if resp.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}
			payload := resp.Model

			if metadata.ResourceData.HasChange("azure_active_directory") {
				payload.Properties.AadProfile = expandArcKubernetesClusterAadProfile(model.AzureActiveDirectory)
			}

			if metadata.ResourceData.HasChange("tags") {
				payload.Tags = pointer.To(model.Tags)
			}

			if metadata.ResourceData.HasChange("arc_agent_desired_version") {
				if desiredVersion := model.ArcAgentDesiredVersion; desiredVersion != "" {
					payload.Properties.ArcAgentProfile.DesiredAgentVersion = pointer.To(desiredVersion)
				} else {
					payload.Properties.ArcAgentProfile.DesiredAgentVersion = nil
				}
			}

			if metadata.ResourceData.HasChange("arc_agent_auto_upgrade_enabled") {
				autoUpgradeOption := arckubernetes.AutoUpgradeOptionsEnabled
				if !model.ArcAgentAutoUpgradeEnabled {
					autoUpgradeOption = arckubernetes.AutoUpgradeOptionsDisabled
				}

				payload.Properties.ArcAgentProfile.AgentAutoUpgrade = pointer.To(autoUpgradeOption)
			}

			if err := client.ConnectedClusterCreateThenPoll(ctx, *id, *payload); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ArcKubernetesClient

			id, err := arckubernetes.ParseConnectedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.ConnectedClusterGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := ArcKubernetesProvisionedClusterModel{
				Name:              id.ConnectedClusterName,
				ResourceGroupName: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				state.Identity = identity.FlattenSystemAssignedToModel(&model.Identity)
				state.Location = location.Normalize(model.Location)
				state.Tags = pointer.From(model.Tags)

				props := model.Properties
				state.AzureActiveDirectory = flattenArcKubernetesClusterAadProfile(props.AadProfile)
				state.AgentVersion = pointer.From(props.AgentVersion)
				state.Distribution = pointer.From(props.Distribution)
				state.Infrastructure = pointer.From(props.Infrastructure)
				state.KubernetesVersion = pointer.From(props.KubernetesVersion)
				state.Offering = pointer.From(props.Offering)
				state.TotalCoreCount = pointer.From(props.TotalCoreCount)
				state.TotalNodeCount = pointer.From(props.TotalNodeCount)

				arcAgentAutoUpgradeEnabled := true
				arcAgentdesiredVersion := ""
				if arcAgentProfile := props.ArcAgentProfile; arcAgentProfile != nil {
					arcAgentdesiredVersion = pointer.From(arcAgentProfile.DesiredAgentVersion)
					if arcAgentProfile.AgentAutoUpgrade != nil && *arcAgentProfile.AgentAutoUpgrade == arckubernetes.AutoUpgradeOptionsDisabled {
						arcAgentAutoUpgradeEnabled = false
					}
				}
				state.ArcAgentAutoUpgradeEnabled = arcAgentAutoUpgradeEnabled
				state.ArcAgentDesiredVersion = arcAgentdesiredVersion
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ArcKubernetesProvisionedClusterResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ArcKubernetes.ArcKubernetesClient

			id, err := arckubernetes.ParseConnectedClusterID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.ConnectedClusterDeleteThenPoll(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandArcKubernetesClusterAadProfile(input []AzureActiveDirectoryModel) *arckubernetes.AadProfile {
	if len(input) == 0 {
		return nil
	}

	v := input[0]
	output := arckubernetes.AadProfile{
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

func flattenArcKubernetesClusterAadProfile(input *arckubernetes.AadProfile) []AzureActiveDirectoryModel {
	if input == nil || (input.EnableAzureRBAC == nil && input.AdminGroupObjectIDs == nil && input.TenantID == nil) {
		return make([]AzureActiveDirectoryModel, 0)
	}

	return []AzureActiveDirectoryModel{
		{
			AzureRbacEnabled:    pointer.From(input.EnableAzureRBAC),
			AdminGroupObjectIds: pointer.From(input.AdminGroupObjectIDs),
			TenantId:            pointer.From(input.TenantID),
		},
	}
}
