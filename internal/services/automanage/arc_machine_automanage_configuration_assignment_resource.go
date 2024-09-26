// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package automanage

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type ArcMachineConfigurationAssignment struct {
	ArcMachineId    string `tfschema:"arc_machine_id"`
	ConfigurationId string `tfschema:"configuration_id"`
}

func (v ArcMachineConfigurationAssignment) Arguments() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"arc_machine_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationprofileassignments.ValidateMachineID,
		},
		"configuration_id": {
			Type:         schema.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: configurationprofiles.ValidateConfigurationProfileID,
		},
	}
}

func (v ArcMachineConfigurationAssignment) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (v ArcMachineConfigurationAssignment) ModelObject() interface{} {
	return &ArcMachineConfigurationAssignment{}
}

func (v ArcMachineConfigurationAssignment) ResourceType() string {
	return "azurerm_arc_machine_automanage_configuration_assignment"
}

func (v ArcMachineConfigurationAssignment) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileArcMachineAssignmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model ArcMachineConfigurationAssignment
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			arcMachineId, err := configurationprofileassignments.ParseMachineID(model.ArcMachineId)
			if err != nil {
				return err
			}

			configurationId, err := configurationprofiles.ParseConfigurationProfileID(model.ConfigurationId)
			if err != nil {
				return err
			}

			// Currently, the configuration profile assignment name has to be hardcoded to "default" by API requirement.
			id := configurationprofilehcrpassignments.NewProviders2ConfigurationProfileAssignmentID(subscriptionId, arcMachineId.ResourceGroupName, arcMachineId.MachineName, "default")
			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(v.ResourceType(), id)
			}

			properties := configurationprofilehcrpassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehcrpassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: pointer.FromString(configurationId.ID()),
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

func (v ArcMachineConfigurationAssignment) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileArcMachineAssignmentsClient
			id, err := configurationprofilehcrpassignments.ParseProviders2ConfigurationProfileAssignmentID(metadata.ResourceData.Id())
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

			state := ArcMachineConfigurationAssignment{}

			if model := resp.Model; model != nil {
				configurationId, err := configurationprofiles.ParseConfigurationProfileID(*model.Properties.ConfigurationProfile)
				if err != nil {
					return err
				}
				state.ConfigurationId = configurationId.ID()

				arcMachineId, err := configurationprofileassignments.ParseMachineID(*model.Properties.TargetId)
				if err != nil {
					return err
				}
				state.ArcMachineId = arcMachineId.ID()
			}

			return metadata.Encode(&state)
		},
	}
}

func (v ArcMachineConfigurationAssignment) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileArcMachineAssignmentsClient

			id, err := configurationprofilehcrpassignments.ParseProviders2ConfigurationProfileAssignmentID(metadata.ResourceData.Id())
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

func (v ArcMachineConfigurationAssignment) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return configurationprofilehcrpassignments.ValidateProviders2ConfigurationProfileAssignmentID
}

var _ sdk.Resource = &ArcMachineConfigurationAssignment{}
