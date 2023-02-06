package digitaltwins

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/digitaltwins/2020-12-01/endpoints"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func validateEndpointType(validate func(input endpoints.DigitalTwinsEndpointResourceProperties) error) pluginsdk.ImporterFunc {
	return func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
		id, err := endpoints.ParseEndpointID(d.Id())
		if err != nil {
			return nil, err
		}

		client := meta.(*clients.Client).DigitalTwins.EndpointClient
		resp, err := client.DigitalTwinsEndpointGet(ctx, *id)
		if err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		if resp.Model == nil {
			return nil, fmt.Errorf("retrieving %s: model was nil", *id)
		}
		if resp.Model.Properties == nil {
			return nil, fmt.Errorf("retrieving %s: model.Properties was nil", *id)
		}
		if resp.Model.Properties == nil {
			return nil, fmt.Errorf("retrieving %s: model.Properties was nil", *id)
		}
		if err := validate(resp.Model.Properties); err != nil {
			return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
		}
		return []*pluginsdk.ResourceData{d}, nil
	}
}
