// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package chaosstudio

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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	chaosstudioClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/chaosstudio/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ChaosStudioCapabilityResource{}

type ChaosStudioCapabilityResource struct{}

func (r ChaosStudioCapabilityResource) ModelObject() any {
	return &ChaosStudioTargetResourceSchema{}
}

type ChaosStudioCapabilityResourceSchema struct {
	CapabilityType      string `tfschema:"capability_type"`
	ChaosStudioTargetId string `tfschema:"chaos_studio_target_id"`
	Urn                 string `tfschema:"urn"`
}

func (r ChaosStudioCapabilityResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateChaosStudioCapabilityID
}

func (r ChaosStudioCapabilityResource) ResourceType() string {
	return "azurerm_chaos_studio_capability"
}

func (r ChaosStudioCapabilityResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"capability_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},
		"chaos_studio_target_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: commonids.ValidateChaosStudioTargetID,
		},
	}
}

func (r ChaosStudioCapabilityResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"urn": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r ChaosStudioCapabilityResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := chaosstudioClient.From(metadata.Client).V20231101.Capabilities
			capabilityTypesClient := chaosstudioClient.From(metadata.Client).V20231101.CapabilityTypes
			targetsClient := chaosstudioClient.From(metadata.Client).V20231101.Targets
			subscriptionId := account.From(metadata.Client).SubscriptionId

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

			if !features.From(metadata.Client).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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

func (r ChaosStudioCapabilityResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := chaosstudioClient.From(metadata.Client).V20231101.Capabilities
			schema := ChaosStudioCapabilityResourceSchema{}

			id, err := commonids.ParseChaosStudioCapabilityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			targetId := commonids.NewChaosStudioTargetID(id.Scope, id.TargetName)

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			schema.ChaosStudioTargetId = targetId.ID()
			schema.CapabilityType = id.CapabilityName

			if model := resp.Model; model != nil {
				if props := model.Properties; props != nil {
					schema.Urn = pointer.From(props.Urn)
				}
			}

			return metadata.Encode(&schema)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r ChaosStudioCapabilityResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := chaosstudioClient.From(metadata.Client).V20231101.Capabilities

			id, err := commonids.ParseChaosStudioCapabilityID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
