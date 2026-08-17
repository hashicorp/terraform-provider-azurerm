package trafficmanagers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type HeatMapGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *HeatMapModel
}

type HeatMapGetOperationOptions struct {
	BotRight *string
	TopLeft  *string
}

func DefaultHeatMapGetOperationOptions() HeatMapGetOperationOptions {
	return HeatMapGetOperationOptions{}
}

func (o HeatMapGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o HeatMapGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o HeatMapGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.BotRight != nil {
		out.Append("botRight", fmt.Sprintf("%v", *o.BotRight))
	}
	if o.TopLeft != nil {
		out.Append("topLeft", fmt.Sprintf("%v", *o.TopLeft))
	}
	return &out
}

// HeatMapGet ...
func (c TrafficmanagersClient) HeatMapGet(ctx context.Context, id TrafficManagerProfileId, options HeatMapGetOperationOptions) (result HeatMapGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/heatMaps/default", id.ID()),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var model HeatMapModel
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
