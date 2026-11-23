// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaos_studio_capability

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/capabilities"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/capabilitytypes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Capabilities
			capabilityTypesClient := metadata.Client.ChaosStudio.V20231101.CapabilityTypes
			targetsClient := metadata.Client.ChaosStudio.V20231101.Targets
			subscriptionId := metadata.Client.Account.SubscriptionId

			var config ChaosStudioCapabilityResourceSchema

			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			// need to get the location of the target to be able to validate the capability type
			targetId, err := commonids.ParseChaosStudioTargetID(config.ChaosStudioTargetId)
			if err != nil {
				return err
			}
			existingTarget, err := targetsClient.Get(ctx, *targetId)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", targetId, err)
			}

			if existingTarget.Model == nil && existingTarget.Model.Location == nil {
				return fmt.Errorf("location for %s was nil", targetId)
			}

			targetTypeId := capabilitytypes.NewTargetTypeID(subscriptionId, *existingTarget.Model.Location, targetId.TargetName)

			// validate capability since valid values are dependent on the target type
			capabilityTypes := make([]string, 0)
			resp, err := capabilityTypesClient.ListComplete(ctx, targetTypeId, capabilitytypes.DefaultListOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving list of chaos capability types: %+v", err)
			}
			typeIsValid := false
			for _, item := range resp.Items {
				if name := item.Name; name != nil {
					if strings.EqualFold(config.CapabilityType, *item.Name) {
						typeIsValid = true
					}
					capabilityTypes = append(capabilityTypes, pointer.From(item.Name))
				}
			}

			if !typeIsValid {
				return fmt.Errorf("%q is not a valid `capability_type` for the target type %q, must be one of %+v", config.CapabilityType, targetId.TargetName, capabilityTypes)
			}

			id := commonids.NewChaosStudioCapabilityID(targetId.Scope, targetId.TargetName, config.CapabilityType)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id)
				if err != nil {
					if !response.WasNotFound(existing.HttpResponse) {
						return fmt.Errorf("checking for the presence of an existing %s: %+v", id, err)
					}
				}
				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			var payload capabilities.Capability

			// The API only accepts requests with an empty body for Properties
			props := capabilities.CapabilityProperties{}

			payload.Properties = pointer.To(props)

			if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
