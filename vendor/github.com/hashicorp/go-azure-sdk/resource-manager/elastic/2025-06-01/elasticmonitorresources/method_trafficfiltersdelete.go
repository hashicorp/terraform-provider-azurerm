package elasticmonitorresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TrafficFiltersDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type TrafficFiltersDeleteOperationOptions struct {
	RulesetId *string
}

func DefaultTrafficFiltersDeleteOperationOptions() TrafficFiltersDeleteOperationOptions {
	return TrafficFiltersDeleteOperationOptions{}
}

func (o TrafficFiltersDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TrafficFiltersDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o TrafficFiltersDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RulesetId != nil {
		out.Append("rulesetId", fmt.Sprintf("%v", *o.RulesetId))
	}
	return &out
}

// TrafficFiltersDelete ...
func (c ElasticMonitorResourcesClient) TrafficFiltersDelete(ctx context.Context, id MonitorId, options TrafficFiltersDeleteOperationOptions) (result TrafficFiltersDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/deleteTrafficFilter", id.ID()),
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

	return
}
