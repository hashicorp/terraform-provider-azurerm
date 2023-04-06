package kubernetesconfiguration

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/kubernetesconfiguration/2022-11-01/extensions"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

type PlanModel struct {
	Name          string `tfschema:"name"`
	Product       string `tfschema:"product"`
	PromotionCode string `tfschema:"promotion_code"`
	Publisher     string `tfschema:"publisher"`
	Version       string `tfschema:"version"`
}

func commonArguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringMatch(
				regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-.]{0,252}$"),
				"name must be between 1 and 253 characters in length and may contain only letters, numbers, periods (.), hyphens (-), and must begin with a letter or number.",
			),
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"cluster_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"extension_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"configuration_protected_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"configuration_settings": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"plan": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			ForceNew: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"product": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"publisher": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"promotion_code": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"version": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ForceNew:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
		},

		"release_train": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"version"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"release_namespace": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"target_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"target_namespace": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			Computed:      true,
			ForceNew:      true,
			ConflictsWith: []string{"release_namespace"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},

		"version": {
			Type:          pluginsdk.TypeString,
			Optional:      true,
			ForceNew:      true,
			ConflictsWith: []string{"release_train"},
			ValidateFunc:  validation.StringIsNotEmpty,
		},
	}
}

func deleteExtension() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.KubernetesConfiguration.ExtensionsClient

			id, err := extensions.ParseExtensionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if err := client.DeleteThenPoll(ctx, *id, extensions.DeleteOperationOptions{}); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandPlanModel(inputList []PlanModel) *extensions.Plan {
	if len(inputList) == 0 {
		return nil
	}
	input := &inputList[0]
	output := extensions.Plan{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}
	if input.PromotionCode != "" {
		output.PromotionCode = &input.PromotionCode
	}

	if input.Version != "" {
		output.Version = &input.Version
	}

	return &output
}

func flattenPlanModel(input *extensions.Plan) []PlanModel {
	var outputList []PlanModel
	if input == nil {
		return outputList
	}
	output := PlanModel{
		Name:      input.Name,
		Product:   input.Product,
		Publisher: input.Publisher,
	}
	if input.PromotionCode != nil {
		output.PromotionCode = *input.PromotionCode
	}

	if input.Version != nil {
		output.Version = *input.Version
	}

	return append(outputList, output)
}
