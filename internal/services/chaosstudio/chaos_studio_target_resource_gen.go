package chaosstudio

// NOTE: this file is generated - manual changes will be overwritten.
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.
import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/chaosstudio/2023-11-01/targets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.Resource = ChaosStudioTargetResource{}

type ChaosStudioTargetResource struct{}

func (r ChaosStudioTargetResource) ModelObject() interface{} {
	return &ChaosStudioTargetResourceSchema{}
}

type ChaosStudioTargetResourceSchema struct {
	Location         string `tfschema:"location"`
	TargetResourceId string `tfschema:"target_resource_id"`
	TargetType       string `tfschema:"target_type"`
}

func (r ChaosStudioTargetResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return commonids.ValidateChaosStudioTargetID
}
func (r ChaosStudioTargetResource) ResourceType() string {
	return "azurerm_chaos_studio_target"
}
func (r ChaosStudioTargetResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"location": commonschema.Location(),
		"target_resource_id": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
		"target_type": {
			ForceNew: true,
			Required: true,
			Type:     pluginsdk.TypeString,
		},
	}
}
func (r ChaosStudioTargetResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}
func (r ChaosStudioTargetResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Targets
			schema := ChaosStudioTargetResourceSchema{}

			id, err := commonids.ParseChaosStudioTargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(*id)
				}
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model := resp.Model; model != nil {
				schema.TargetResourceId = id.Scope
				schema.TargetType = id.TargetName
				if err := r.mapTargetToChaosStudioTargetResourceSchema(*model, &schema); err != nil {
					return fmt.Errorf("flattening model: %+v", err)
				}
			}

			return metadata.Encode(&schema)
		},
	}
}
func (r ChaosStudioTargetResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.ChaosStudio.V20231101.Targets

			id, err := commonids.ParseChaosStudioTargetID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ChaosStudioTargetResource) mapChaosStudioTargetResourceSchemaToTarget(input ChaosStudioTargetResourceSchema, output *targets.Target) error {
	output.Location = pointer.To(location.Normalize(input.Location))
	return nil
}

func (r ChaosStudioTargetResource) mapTargetToChaosStudioTargetResourceSchema(input targets.Target, output *ChaosStudioTargetResourceSchema) error {
	output.Location = location.NormalizeNilable(input.Location)
	return nil
}
