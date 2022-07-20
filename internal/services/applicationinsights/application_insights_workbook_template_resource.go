package applicationinsights

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/applicationinsights/2020-11-20/applicationinsights"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type ApplicationInsightsWorkbookTemplateModel struct {
	Name              string                         `tfschema:"name"`
	ResourceGroupName string                         `tfschema:"resource_group_name"`
	Author            string                         `tfschema:"author"`
	Galleries         []WorkbookTemplateGalleryModel `tfschema:"galleries"`
	Localized         string                         `tfschema:"localized"`
	Location          string                         `tfschema:"location"`
	Priority          int64                          `tfschema:"priority"`
	Tags              map[string]string              `tfschema:"tags"`
	TemplateData      string                         `tfschema:"template_data"`
}

type WorkbookTemplateGalleryModel struct {
	Name         string `tfschema:"name"`
	Category     string `tfschema:"category"`
	Order        int64  `tfschema:"order"`
	ResourceType string `tfschema:"resource_type"`
	Type         string `tfschema:"type"`
}

type ApplicationInsightsWorkbookTemplateResource struct{}

var _ sdk.ResourceWithUpdate = ApplicationInsightsWorkbookTemplateResource{}

func (r ApplicationInsightsWorkbookTemplateResource) ResourceType() string {
	return "azurerm_application_insights_workbook_template"
}

func (r ApplicationInsightsWorkbookTemplateResource) ModelObject() interface{} {
	return &ApplicationInsightsWorkbookTemplateModel{}
}

func (r ApplicationInsightsWorkbookTemplateResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return applicationinsights.ValidateWorkbookTemplateID
}

func (r ApplicationInsightsWorkbookTemplateResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"resource_group_name": commonschema.ResourceGroupName(),

		"location": commonschema.Location(),

		"template_data": {
			Type:             pluginsdk.TypeString,
			Required:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"galleries": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"category": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"order": {
						Type:     pluginsdk.TypeInt,
						Optional: true,
						Default:  0,
					},

					"resource_type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Default:      "Azure Monitor",
					},

					"type": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
						Default:      "workbook",
					},
				},
			},
		},

		"author": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"localized": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"priority": {
			Type:     pluginsdk.TypeInt,
			Optional: true,
			Default:  0,
		},

		"tags": commonschema.Tags(),
	}
}

func (r ApplicationInsightsWorkbookTemplateResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r ApplicationInsightsWorkbookTemplateResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model ApplicationInsightsWorkbookTemplateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppInsights.WorkbookTemplateClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := applicationinsights.NewWorkbookTemplateID(subscriptionId, model.ResourceGroupName, model.Name)
			existing, err := client.WorkbookTemplatesGet(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			var templateDataValue interface{}
			err = json.Unmarshal([]byte(model.TemplateData), &templateDataValue)
			if err != nil {
				return err
			}

			properties := &applicationinsights.WorkbookTemplate{
				Location: location.Normalize(model.Location),
				Properties: &applicationinsights.WorkbookTemplateProperties{
					Priority:     &model.Priority,
					TemplateData: templateDataValue,
				},

				Tags: &model.Tags,
			}

			if model.Author != "" {
				properties.Properties.Author = &model.Author
			}

			if model.Localized != "" {
				var localizedValue map[string][]applicationinsights.WorkbookTemplateLocalizedGallery
				if err := json.Unmarshal([]byte(model.Localized), &localizedValue); err != nil {
					return err
				}

				properties.Properties.Localized = &localizedValue
			}

			galleriesValue, err := expandWorkbookTemplateGalleryModel(model.Galleries)
			if err != nil {
				return err
			}

			if galleriesValue != nil {
				properties.Properties.Galleries = *galleriesValue
			}

			if _, err := client.WorkbookTemplatesCreateOrUpdate(ctx, id, *properties); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r ApplicationInsightsWorkbookTemplateResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookTemplateClient

			id, err := applicationinsights.ParseWorkbookTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model ApplicationInsightsWorkbookTemplateModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.WorkbookTemplatesGet(ctx, *id)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Model
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("author") {
				properties.Properties.Author = &model.Author
			}

			if metadata.ResourceData.HasChange("galleries") {
				galleriesValue, err := expandWorkbookTemplateGalleryModel(model.Galleries)
				if err != nil {
					return err
				}

				if galleriesValue != nil {
					properties.Properties.Galleries = *galleriesValue
				}
			}

			if metadata.ResourceData.HasChange("priority") {
				properties.Properties.Priority = &model.Priority
			}

			if metadata.ResourceData.HasChange("template_data") {
				var templateDataValue interface{}
				err := json.Unmarshal([]byte(model.TemplateData), &templateDataValue)
				if err != nil {
					return err
				}

				properties.Properties.TemplateData = templateDataValue
			}

			if metadata.ResourceData.HasChange("localized") {
				var localizedValue map[string][]applicationinsights.WorkbookTemplateLocalizedGallery
				if err := json.Unmarshal([]byte(model.Localized), &localizedValue); err != nil {
					return err
				}

				properties.Properties.Localized = &localizedValue
			}

			if metadata.ResourceData.HasChange("tags") {
				properties.Tags = &model.Tags
			}

			if _, err := client.WorkbookTemplatesCreateOrUpdate(ctx, *id, *properties); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r ApplicationInsightsWorkbookTemplateResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookTemplateClient

			id, err := applicationinsights.ParseWorkbookTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.WorkbookTemplatesGet(ctx, *id)
			if err != nil {
				if response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			model := resp.Model
			if model == nil {
				return fmt.Errorf("retrieving %s: model was nil", id)
			}

			state := ApplicationInsightsWorkbookTemplateModel{
				Name:              id.ResourceName,
				ResourceGroupName: id.ResourceGroupName,
				Location:          location.Normalize(model.Location),
			}

			if properties := model.Properties; properties != nil {
				if properties.Author != nil {
					state.Author = *properties.Author
				}

				galleriesValue, err := flattenWorkbookTemplateGalleryModel(&properties.Galleries)
				if err != nil {
					return err
				}

				state.Galleries = galleriesValue

				if properties.Priority != nil {
					state.Priority = *properties.Priority
				}

				if properties.TemplateData != nil {
					templateDataValue, err := json.Marshal(properties.TemplateData)
					if err != nil {
						return err
					}

					state.TemplateData = string(templateDataValue)
				}

				if properties.Localized != nil {
					localizedValue, err := json.Marshal(properties.Localized)
					if err != nil {
						return err
					}

					state.Localized = string(localizedValue)
				}
			}

			if model.Tags != nil {
				state.Tags = *model.Tags
			}

			return metadata.Encode(&state)
		},
	}
}

func (r ApplicationInsightsWorkbookTemplateResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppInsights.WorkbookTemplateClient

			id, err := applicationinsights.ParseWorkbookTemplateID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.WorkbookTemplatesDelete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", id, err)
			}

			return nil
		},
	}
}

func expandWorkbookTemplateGalleryModel(inputList []WorkbookTemplateGalleryModel) (*[]applicationinsights.WorkbookTemplateGallery, error) {
	var outputList []applicationinsights.WorkbookTemplateGallery
	for _, input := range inputList {
		output := applicationinsights.WorkbookTemplateGallery{
			Category:     utils.String(input.Category),
			Name:         utils.String(input.Name),
			Order:        utils.Int64(input.Order),
			ResourceType: utils.String(input.ResourceType),
			Type:         utils.String(input.Type),
		}

		outputList = append(outputList, output)
	}

	return &outputList, nil
}

func flattenWorkbookTemplateGalleryModel(inputList *[]applicationinsights.WorkbookTemplateGallery) ([]WorkbookTemplateGalleryModel, error) {
	var outputList []WorkbookTemplateGalleryModel
	if inputList == nil {
		return outputList, nil
	}

	for _, input := range *inputList {
		output := WorkbookTemplateGalleryModel{}

		if input.Category != nil {
			output.Category = *input.Category
		}

		if input.Name != nil {
			output.Name = *input.Name
		}

		if input.Order != nil {
			output.Order = *input.Order
		}

		if input.ResourceType != nil {
			output.ResourceType = *input.ResourceType
		}

		if input.Type != nil {
			output.Type = *input.Type
		}

		outputList = append(outputList, output)
	}

	return outputList, nil
}
