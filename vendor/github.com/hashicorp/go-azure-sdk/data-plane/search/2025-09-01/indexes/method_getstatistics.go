package indexes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetStatisticsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *GetIndexStatisticsResult
}

type GetStatisticsOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultGetStatisticsOperationOptions() GetStatisticsOperationOptions {
	return GetStatisticsOperationOptions{}
}

func (o GetStatisticsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o GetStatisticsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetStatisticsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// GetStatistics ...
func (c IndexesClient) GetStatistics(ctx context.Context, id IndexId, options GetStatisticsOperationOptions) (result GetStatisticsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("/indexes('%s')/search.stats", id.PathElements()...),
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

	var model GetIndexStatisticsResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
