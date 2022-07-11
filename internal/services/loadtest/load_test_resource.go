package loadtest

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/loadtestservice/2021-12-01-preview/loadtests"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type LoadTestResource struct{}

var _ sdk.ResourceWithUpdate = LoadTestResource{}

type LoadTestResourceModel struct {
	Name          string            `tfschema:"name"`
	ResourceGroup string            `tfschema:"resource_group_name"`
	Location      string            `tfschema:"location"`
	Tags          map[string]string `tfschema:"tags"`
	DataPlaneURI  string            `tfschema:"dataplane_uri"`
}

func (r LoadTestResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"resource_group_name": azure.SchemaResourceGroupName(),

		"location": commonschema.Location(),

		"tags": tags.Schema(),
	}
}

func (r LoadTestResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"dataplane_uri": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func (r LoadTestResource) ModelObject() interface{} {
	return &LoadTestResourceModel{}
}

func (r LoadTestResource) ResourceType() string {
	return "azurerm_load_test"
}

func (r LoadTestResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return loadtests.ValidateLoadTestID
}

func (r LoadTestResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			var model LoadTestResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding %+v", err)
			}

			client := metadata.Client.LoadTest.LoadTestsClient
			subscriptionId := metadata.Client.Account.SubscriptionId
			id := loadtests.NewLoadTestID(subscriptionId, model.ResourceGroup, model.Name)

			existing, err := client.Get(ctx, id)
			if err != nil && !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing Linux %s: %+v", id, err)
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			loadTest := loadtests.LoadTestResource{
				Name:     &model.Name,
				Location: model.Location,
				Tags:     &model.Tags,
			}

			_, err = client.CreateOrUpdate(ctx, id, loadTest)
			if err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r LoadTestResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.LoadTest.LoadTestsClient
			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state LoadTestResourceModel
			if err := metadata.Decode(&state); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			existing, err := client.Get(ctx, *id)
			if err != nil {
				return fmt.Errorf("reading Load Test %s: %v", id, err)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = &state.Tags
			}

			_, err = client.CreateOrUpdate(ctx, *id, *existing.Model)
			if err != nil {
				return fmt.Errorf("updating Load Test %s: %+v", id, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}

func (r LoadTestResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.LoadTest.LoadTestsClient

			resp, err := client.Get(ctx, *id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("while checking for Load Test's %q existence: %+v", id.LoadTestName, err)
			}

			state := LoadTestResourceModel{
				Name:          id.LoadTestName,
				Location:      location.NormalizeNilable(utils.String(resp.Model.Location)),
				ResourceGroup: id.ResourceGroupName,
			}

			if model := resp.Model; model != nil {
				if model.Tags != nil {
					state.Tags = *model.Tags
				}
				if props := model.Properties; props != nil {
					state.DataPlaneURI = *props.DataPlaneURI
				}
			}
			return metadata.Encode(&state)
		},
		Timeout: 5 * time.Minute,
	}
}

func (r LoadTestResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			id, err := loadtests.ParseLoadTestID(metadata.ResourceData.Id())
			if err != nil {
				return fmt.Errorf("while parsing resource ID: %+v", err)
			}

			client := metadata.Client.LoadTest.LoadTestsClient

			_, err = client.Delete(ctx, *id)
			if err != nil {
				return fmt.Errorf("while removing Load Test %q: %+v", id.LoadTestName, err)
			}

			return nil
		},
		Timeout: 30 * time.Minute,
	}
}
