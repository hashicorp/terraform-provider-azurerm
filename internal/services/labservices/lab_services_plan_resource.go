package labservices

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservices/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservices/sdk/2021-10-01-preview/labplan"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/labservices/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LabServicesPlanResource struct {
}

var _ sdk.Resource = LabServicesPlanResource{}

type LabServicesPlanResourceModel struct {
	Name          string                    `tfschema:"name"`
	ResourceGroup string                    `tfschema:"resource_group_name"`
	Location      string                    `tfschema:"location"`
	Type          string                    `tfschema:"type"`
	Tags          map[string]interface{}    `tfschema:"tags"`
	Properties    labplan.LabPlanProperties `tfschema:"properties"`
}

func (r LabServicesPlanResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: azure.ValidateResourceID,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": location.Schema(),

		"type": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"tags": tags.Schema(),
	}
}

func (r LabServicesPlanResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r LabServicesPlanResource) ModelObject() interface{} {
	return &LabServicesPlanResourceModel{}
}

func (r LabServicesPlanResource) ResourceType() string {
	return "azurerm_lab_services_plan"
}

func (r LabServicesPlanResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LabServicesPlanResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.LabServices.LabPlanClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := parse.NewLabPlanID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			labplan := labplan.LabPlan{
				Name:     &model.Name,
				Location: model.Location,
			}

			_, err = client.CreateOrUpdate(ctx, id, labplan)
			if err != nil {
				return fmt.Errorf("creating Lab Plan %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r LabServicesPlanResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.LabPlanID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.LabServices.LabPlanClient

			labPlan, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(labPlan.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("while checking for Lab Plans's %q existence: %+v", id.Name, err)
			}

			model := LabServicesPlanResourceModel{
				Name:     id.Name,
				Location: location.NormalizeNilable(utils.String(labPlan.Model.Location)),
			}

			return metadata.Encode(&model)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r LabServicesPlanResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := parse.LabPlanID(metadata.ResourceData.Id())

			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.LabServices.LabPlanClient

			_, err = client.Delete(ctx, id)
			if err != nil {
				return fmt.Errorf("while removing Lab Plan %q: %+v", id.Name, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r LabServicesPlanResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.LabPlanID
}
