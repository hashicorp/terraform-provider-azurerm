package springcloud

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"github.com/tombuildsstuff/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

type SpringCloudApplicationPerformanceMonitoringModel struct {
	Name                 string                 `tfschema:"name"`
	SpringCloudServiceId string                 `tfschema:"spring_cloud_service_id"`
	Type                 string                 `tfschema:"type"`
	GloballyEnabled      bool                   `tfschema:"globally_enabled"`
	Properties           map[string]interface{} `tfschema:"properties"`
	Secrets              map[string]interface{} `tfschema:"secrets"`
}

type SpringCloudApplicationPerformanceMonitoringResource struct{}

var _ sdk.ResourceWithUpdate = SpringCloudApplicationPerformanceMonitoringResource{}

func (s SpringCloudApplicationPerformanceMonitoringResource) ResourceType() string {
	return "azurerm_spring_cloud_application_performance_monitoring"
}

func (s SpringCloudApplicationPerformanceMonitoringResource) ModelObject() interface{} {
	return &SpringCloudApplicationPerformanceMonitoringModel{}
}

func (s SpringCloudApplicationPerformanceMonitoringResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return validate.SpringCloudApplicationPerformanceMonitoringID
}

func (s SpringCloudApplicationPerformanceMonitoringResource) Arguments() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"spring_cloud_service_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.SpringCloudServiceID,
		},

		"type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"globally_enabled": {
			Type:     pluginsdk.TypeBool,
			Optional: true,
			Default:  false,
		},

		"properties": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},

		"secrets": {
			Type:     pluginsdk.TypeMap,
			Optional: true,
			Elem: &pluginsdk.Schema{
				Type:         pluginsdk.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func (s SpringCloudApplicationPerformanceMonitoringResource) Attributes() map[string]*schema.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (s SpringCloudApplicationPerformanceMonitoringResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model SpringCloudApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			client := metadata.Client.AppPlatform.ApmClient
			springId, err := parse.SpringCloudServiceID(model.SpringCloudServiceId)
			if err != nil {
				return fmt.Errorf("parsing spring service ID: %+v", err)
			}
			id := parse.NewSpringCloudApplicationPerformanceMonitoringID(springId.SubscriptionId, springId.ResourceGroup, springId.SpringName, model.Name)

			existing, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApmName)
			if err != nil && !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
			if !utils.ResponseWasNotFound(existing.Response) {
				return metadata.ResourceRequiresImport(s.ResourceType(), id)
			}

			resource := appplatform.ApmResource{
				Properties: &appplatform.ApmProperties{
					Type:       utils.String(model.Type),
					Properties: utils.ExpandMapStringPtrString(model.Properties),
					Secrets:    utils.ExpandMapStringPtrString(model.Secrets),
				},
			}
			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApmName, resource)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for creation of %s: %+v", id, err)
			}

			if model.GloballyEnabled {
				serviceClient := metadata.Client.AppPlatform.ServicesClient
				apmReference := appplatform.ApmReference{
					ResourceID: utils.String(id.ID()),
				}
				future, err := serviceClient.EnableApmGlobally(ctx, springId.ResourceGroup, springId.SpringName, apmReference)
				if err != nil {
					return fmt.Errorf("enabling %s globally: %+v", id, err)
				}
				if err = future.WaitForCompletionRef(ctx, serviceClient.Client); err != nil {
					return fmt.Errorf("waiting for enabling %s globally: %+v", id, err)
				}
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (s SpringCloudApplicationPerformanceMonitoringResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.ApmClient

			id, err := parse.SpringCloudApplicationPerformanceMonitoringID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var model SpringCloudApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApmName)
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			properties := resp.Properties
			if properties == nil {
				return fmt.Errorf("retrieving %s: properties was nil", id)
			}

			if metadata.ResourceData.HasChange("type") {
				properties.Type = &model.Type
			}

			if metadata.ResourceData.HasChange("properties") {
				properties.Properties = utils.ExpandMapStringPtrString(model.Properties)
			}

			if metadata.ResourceData.HasChange("secrets") {
				properties.Secrets = utils.ExpandMapStringPtrString(model.Secrets)
			}

			resource := appplatform.ApmResource{
				Properties: properties,
			}

			future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.SpringName, id.ApmName, resource)
			if err != nil {
				return fmt.Errorf("updating %s: %+v", id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for update of %s: %+v", id, err)
			}

			if metadata.ResourceData.HasChange("globally_enabled") {
				serviceClient := metadata.Client.AppPlatform.ServicesClient
				apmReference := appplatform.ApmReference{
					ResourceID: utils.String(id.ID()),
				}
				if model.GloballyEnabled {
					future, err := serviceClient.EnableApmGlobally(ctx, id.ResourceGroup, id.SpringName, apmReference)
					if err != nil {
						return fmt.Errorf("enabling %s globally: %+v", id, err)
					}
					if err = future.WaitForCompletionRef(ctx, serviceClient.Client); err != nil {
						return fmt.Errorf("waiting for enabling %s globally: %+v", id, err)
					}
				} else {
					future, err := serviceClient.DisableApmGlobally(ctx, id.ResourceGroup, id.SpringName, apmReference)
					if err != nil {
						return fmt.Errorf("disabling %s globally: %+v", id, err)
					}
					if err = future.WaitForCompletionRef(ctx, serviceClient.Client); err != nil {
						return fmt.Errorf("waiting for disabling %s globally: %+v", id, err)
					}
				}
			}

			return nil
		},
	}
}

func (s SpringCloudApplicationPerformanceMonitoringResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.ApmClient
			serviceClient := metadata.Client.AppPlatform.ServicesClient

			id, err := parse.SpringCloudApplicationPerformanceMonitoringID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, err := client.Get(ctx, id.ResourceGroup, id.SpringName, id.ApmName)
			if err != nil {
				if utils.ResponseWasNotFound(resp.Response) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			result, err := serviceClient.ListGloballyEnabledApms(ctx, id.ResourceGroup, id.SpringName)
			if err != nil {
				return fmt.Errorf("listing globally enabled apms: %+v", err)
			}
			globallyEnabled := false
			if result.Value != nil {
				for _, apmID := range *result.Value {
					apmId, err := parse.SpringCloudApplicationPerformanceMonitoringIDInsensitively(apmID)
					if err == nil && apmId.ID() == id.ID() {
						globallyEnabled = true
						break
					}
				}
			}

			var model SpringCloudApplicationPerformanceMonitoringModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			state := SpringCloudApplicationPerformanceMonitoringModel{
				Name:                 id.ApmName,
				SpringCloudServiceId: parse.NewSpringCloudServiceID(id.SubscriptionId, id.ResourceGroup, id.SpringName).ID(),
				Secrets:              model.Secrets,
				GloballyEnabled:      globallyEnabled,
			}

			if props := resp.Properties; props != nil {
				if props.Type != nil {
					state.Type = *props.Type
				}
				if props.Properties != nil {
					state.Properties = utils.FlattenMapStringPtrString(props.Properties)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (s SpringCloudApplicationPerformanceMonitoringResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.AppPlatform.ApmClient

			id, err := parse.SpringCloudApplicationPerformanceMonitoringID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			future, err := client.Delete(ctx, id.ResourceGroup, id.SpringName, id.ApmName)
			if err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
				if !response.WasNotFound(future.Response()) {
					return fmt.Errorf("waiting for deletion of %s: %+v", *id, err)
				}
			}

			return nil
		},
	}
}
