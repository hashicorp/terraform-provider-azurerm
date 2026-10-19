// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_automanage_configuration_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (v Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Automanage.ConfigurationProfileArcMachineAssignmentsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model Resource
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

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(v.ResourceType(), id)
				}
			}

			properties := configurationprofilehcrpassignments.ConfigurationProfileAssignment{
				Properties: &configurationprofilehcrpassignments.ConfigurationProfileAssignmentProperties{
					ConfigurationProfile: pointer.To(configurationId.ID()),
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
