// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package desktopvirtualization

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/desktopvirtualization/2024-04-03/scalingplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/desktopvirtualization/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var (
	_ sdk.Resource           = DesktopVirtualizationScalingPlanHostPoolAssociationResource{}
	_ sdk.ResourceWithUpdate = DesktopVirtualizationScalingPlanHostPoolAssociationResource{}
)

type DesktopVirtualizationScalingPlanHostPoolAssociationResource struct{}

func (DesktopVirtualizationScalingPlanHostPoolAssociationResource) ModelObject() interface{} {
	return &DesktopVirtualizationScalingPlanHostPoolAssociationResourceModel{}
}

func (DesktopVirtualizationScalingPlanHostPoolAssociationResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return func(input interface{}, key string) (warnings []string, errors []error) {
		v, ok := input.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", key))
			return
		}

		if _, err := parse.ScalingPlanHostPoolAssociationID(v); err != nil {
			errors = append(errors, err)
		}

		return
	}
}

func (DesktopVirtualizationScalingPlanHostPoolAssociationResource) ResourceType() string {
	return "azurerm_virtual_desktop_scaling_plan_host_pool_association"
}

type DesktopVirtualizationScalingPlanHostPoolAssociationResourceModel struct {
	ScalingPlanId string `tfschema:"scaling_plan_id"`
	HostPoolId    string `tfschema:"host_pool_id"`
	Enabled       bool   `tfschema:"enabled"`
}

func (r DesktopVirtualizationScalingPlanHostPoolAssociationResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"scaling_plan_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: scalingplan.ValidateScalingPlanID,
		},

		"host_pool_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: scalingplan.ValidateHostPoolID,
		},

		"enabled": {
			Type:     pluginsdk.TypeBool,
			Required: true,
		},
	}
}

func (r DesktopVirtualizationScalingPlanHostPoolAssociationResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r DesktopVirtualizationScalingPlanHostPoolAssociationResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			var model DesktopVirtualizationScalingPlanHostPoolAssociationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			log.Printf("[INFO] preparing arguments for Virtual Desktop Scaling Plan <-> Host Pool Association creation.")
			scalingPlanId, err := scalingplan.ParseScalingPlanID(model.ScalingPlanId)
			if err != nil {
				return err
			}
			hostPoolId, err := scalingplan.ParseHostPoolID(model.HostPoolId)
			if err != nil {
				return err
			}
			associationId := parse.NewScalingPlanHostPoolAssociationId(*scalingPlanId, *hostPoolId)

			locks.ByName(scalingPlanId.ScalingPlanName, scalingPlanResourceType)
			defer locks.UnlockByName(scalingPlanId.ScalingPlanName, scalingPlanResourceType)

			locks.ByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())
			defer locks.UnlockByName(hostPoolId.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())

			existing, err := client.Get(ctx, *scalingPlanId)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", *scalingPlanId)
				}

				return fmt.Errorf("retrieving %s: %+v", *scalingPlanId, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", *scalingPlanId)
			}
			scalingPlanModel := *existing.Model

			hostPoolAssociations := []scalingplan.ScalingHostPoolReference{}
			if v := scalingPlanModel.Properties.HostPoolReferences; v != nil {
				hostPoolAssociations = *v
			}

			hostPoolStr := hostPoolId.ID()
			if scalingPlanHostPoolAssociationExists(scalingPlanModel.Properties, hostPoolStr) {
				return metadata.ResourceRequiresImport(r.ResourceType(), associationId)
			}
			hostPoolAssociations = append(hostPoolAssociations, scalingplan.ScalingHostPoolReference{
				HostPoolArmPath:    &hostPoolStr,
				ScalingPlanEnabled: pointer.To(model.Enabled),
			})

			payload := scalingplan.ScalingPlanPatch{
				Properties: &scalingplan.ScalingPlanPatchProperties{
					HostPoolReferences: &hostPoolAssociations,
					Schedules:          scalingPlanModel.Properties.Schedules,
				},
				Tags: scalingPlanModel.Tags,
			}
			if _, err = client.Update(ctx, *scalingPlanId, payload); err != nil {
				return fmt.Errorf("creating association between %s and %s: %+v", *scalingPlanId, *hostPoolId, err)
			}

			metadata.SetID(associationId)
			return nil
		},
	}
}

func (r DesktopVirtualizationScalingPlanHostPoolAssociationResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			state := DesktopVirtualizationScalingPlanHostPoolAssociationResourceModel{}

			id, err := parse.ScalingPlanHostPoolAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			scalingPlan, err := client.Get(ctx, id.ScalingPlan)
			if err != nil {
				if response.WasNotFound(scalingPlan.HttpResponse) {
					log.Printf("[DEBUG] %s was not found - removing from state!", id.ScalingPlan)
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
			}
			if model := scalingPlan.Model; model != nil {
				hostPoolId := id.HostPool.ID()
				exists := scalingPlanHostPoolAssociationExists(model.Properties, hostPoolId)
				if !exists {
					log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.ScalingPlan, id.HostPool)
					return metadata.MarkAsGone(id)
				}
				if v := model.Properties.HostPoolReferences; v != nil {
					for _, referenceId := range *v {
						if referenceId.HostPoolArmPath != nil {
							if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
								state.Enabled = pointer.From(referenceId.ScalingPlanEnabled)
							}
						}
					}
				}

				state.ScalingPlanId = id.ScalingPlan.ID()
				state.HostPoolId = hostPoolId
			}

			return metadata.Encode(&state)
		},
	}
}

func (r DesktopVirtualizationScalingPlanHostPoolAssociationResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			var model DesktopVirtualizationScalingPlanHostPoolAssociationResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id, err := parse.ScalingPlanHostPoolAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)
			defer locks.UnlockByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)

			locks.ByName(id.HostPool.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())
			defer locks.UnlockByName(id.HostPool.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())

			existing, err := client.Get(ctx, id.ScalingPlan)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id.ScalingPlan)
				}

				return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id.ScalingPlan)
			}
			scalingPlanModel := *existing.Model
			if !scalingPlanHostPoolAssociationExists(scalingPlanModel.Properties, id.HostPool.ID()) {
				log.Printf("[DEBUG] Association between %s and %s was not found - removing from state!", id.ScalingPlan, id.HostPool)
				return metadata.MarkAsGone(id)
			}

			hostPoolReferences := []scalingplan.ScalingHostPoolReference{}
			hostPoolId := id.HostPool.ID()
			if v := scalingPlanModel.Properties.HostPoolReferences; v != nil {
				for _, referenceId := range *v {
					if referenceId.HostPoolArmPath != nil {
						if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
							referenceId.ScalingPlanEnabled = pointer.To(model.Enabled)
						}
					}
					hostPoolReferences = append(hostPoolReferences, referenceId)
				}
			}

			payload := scalingplan.ScalingPlanPatch{
				Properties: &scalingplan.ScalingPlanPatchProperties{
					HostPoolReferences: &hostPoolReferences,
					Schedules:          scalingPlanModel.Properties.Schedules,
				},
				Tags: scalingPlanModel.Tags,
			}
			if _, err = client.Update(ctx, id.ScalingPlan, payload); err != nil {
				return fmt.Errorf("updating association between %s and %s: %+v", id.ScalingPlan, id.HostPool, err)
			}

			return nil
		},
	}
}

func (r DesktopVirtualizationScalingPlanHostPoolAssociationResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 60 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.DesktopVirtualization.ScalingPlansClient

			id, err := parse.ScalingPlanHostPoolAssociationID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			locks.ByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)
			defer locks.UnlockByName(id.ScalingPlan.ScalingPlanName, scalingPlanResourceType)

			locks.ByName(id.HostPool.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())
			defer locks.UnlockByName(id.HostPool.HostPoolName, DesktopVirtualizationHostPoolResource{}.ResourceType())

			existing, err := client.Get(ctx, id.ScalingPlan)
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("%s was not found", id.ScalingPlan)
				}

				return fmt.Errorf("retrieving %s: %+v", id.ScalingPlan, err)
			}
			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id.ScalingPlan)
			}
			model := *existing.Model

			hostPoolReferences := []scalingplan.ScalingHostPoolReference{}
			hostPoolId := id.HostPool.ID()
			if v := model.Properties.HostPoolReferences; v != nil {
				for _, referenceId := range *v {
					if referenceId.HostPoolArmPath != nil {
						if strings.EqualFold(*referenceId.HostPoolArmPath, hostPoolId) {
							continue
						}
					}

					hostPoolReferences = append(hostPoolReferences, referenceId)
				}
			}

			payload := scalingplan.ScalingPlanPatch{
				Properties: &scalingplan.ScalingPlanPatchProperties{
					HostPoolReferences: &hostPoolReferences,
					Schedules:          model.Properties.Schedules,
				},
				Tags: model.Tags,
			}
			if _, err = client.Update(ctx, id.ScalingPlan, payload); err != nil {
				return fmt.Errorf("removing association between %s and %s: %+v", id.ScalingPlan, id.HostPool, err)
			}

			return nil
		},
	}
}

func scalingPlanHostPoolAssociationExists(props scalingplan.ScalingPlanProperties, applicationGroupId string) bool {
	if props.HostPoolReferences == nil {
		return false
	}

	for _, id := range *props.HostPoolReferences {
		if id.HostPoolArmPath != nil {
			if strings.EqualFold(*id.HostPoolArmPath, applicationGroupId) {
				return true
			}
		}
	}

	return false
}
