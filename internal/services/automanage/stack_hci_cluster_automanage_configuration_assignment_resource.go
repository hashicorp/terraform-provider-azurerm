// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automanage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehciassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/go-azure-sdk/resource-manager/azurestackhci/2024-01-01/clusters"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type StackHCIClusterConfigurationAssignment struct {
	StackHCIClusterId string `tfschema:"stack_hci_cluster_id"`
	ConfigurationId   string `tfschema:"configuration_id"`
}

func (v StackHCIClusterConfigurationAssignment) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"stack_hci_cluster_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: clusters.ValidateClusterID,
		},
		"configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationprofiles.ValidateConfigurationProfileID,
		},
	}
}

func (v StackHCIClusterConfigurationAssignment) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (v StackHCIClusterConfigurationAssignment) ModelObject() interface{} {
	return &StackHCIClusterConfigurationAssignment{}
}

func (v StackHCIClusterConfigurationAssignment) ResourceType() string {
	return "azurerm_stack_hci_cluster_automanage_configuration_assignment"
}

func (v StackHCIClusterConfigurationAssignment) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileHCIAssignmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model StackHCIClusterConfigurationAssignment
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			stackHciClusterId, err := clusters.ParseClusterID(model.StackHCIClusterId)
			if err != nil {
				return err
			}

			configurationId, err := configurationprofiles.ParseConfigurationProfileID(model.ConfigurationId)
			if err != nil {
				return err
			}

			// Currently, the configuration profile assignment name has to be hardcoded to "default" by API requirement.
			id := configurationprofilehciassignments.NewConfigurationProfileAssignmentID(subscriptionId, stackHciClusterId.ResourceGroupName, stackHciClusterId.ClusterName, "default")
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(v.ResourceType(), id)
			}

			properties := configurationprofilehciassignments.ConfigurationProfileAssignment{
				Name: pointer.To(id.ConfigurationProfileAssignmentName),
				Properties: &configurationprofilehciassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: pointer.To(configurationId.ID()),
					TargetId:             pointer.To(stackHciClusterId.ID()),
				},
			}

			if _, respErr := client.CreateOrUpdate(ctx, id, properties); respErr != nil {
				return fmt.Errorf("creating %s: %+v", id.String(), respErr)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (v StackHCIClusterConfigurationAssignment) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileHCIAssignmentsClient
			id, err := configurationprofilehciassignments.ParseConfigurationProfileAssignmentID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := StackHCIClusterConfigurationAssignment{}

			if model := resp.Model; model != nil {
				configurationId, err := configurationprofiles.ParseConfigurationProfileID(*model.Properties.ConfigurationProfile)
				if err != nil {
					return err
				}
				state.ConfigurationId = configurationId.ID()

				stackHciClusterId, err := clusters.ParseClusterID(*model.Properties.TargetId)
				if err != nil {
					return err
				}
				state.StackHCIClusterId = stackHciClusterId.ID()
			}

			return metadata.Encode(&state)
		},
	}
}

func (v StackHCIClusterConfigurationAssignment) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileHCIAssignmentsClient

			id, err := configurationprofilehciassignments.ParseConfigurationProfileAssignmentID(metadata.ResourceData.Id())
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

func (v StackHCIClusterConfigurationAssignment) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurationprofilehciassignments.ValidateConfigurationProfileAssignmentID
}

var _ sdk.Resource = &StackHCIClusterConfigurationAssignment{}
