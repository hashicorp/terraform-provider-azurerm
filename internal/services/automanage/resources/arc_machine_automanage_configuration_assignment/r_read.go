// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package arc_machine_automanage_configuration_assignment

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofileassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofilehcrpassignments"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automanage/2022-05-04/configurationprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (v Resource) Read() sdk.ResourceFunc {
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

			state := Resource{}

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
