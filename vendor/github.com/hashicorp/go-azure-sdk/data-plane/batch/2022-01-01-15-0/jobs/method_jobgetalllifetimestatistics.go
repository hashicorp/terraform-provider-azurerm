package jobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type JobGetAllLifetimeStatisticsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *JobStatistics
}

type JobGetAllLifetimeStatisticsOperationOptions struct {
	ClientRequestId       *string
	OcpDate               *string
	ReturnClientRequestId *bool
	Timeout               *int64
}

func DefaultJobGetAllLifetimeStatisticsOperationOptions() JobGetAllLifetimeStatisticsOperationOptions {
	return JobGetAllLifetimeStatisticsOperationOptions{}
}

func (o JobGetAllLifetimeStatisticsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.ClientRequestId != nil {
		out.Append("client-request-id", fmt.Sprintf("%v", *o.ClientRequestId))
	}
	if o.OcpDate != nil {
		out.Append("ocp-date", fmt.Sprintf("%v", *o.OcpDate))
	}
	if o.ReturnClientRequestId != nil {
		out.Append("return-client-request-id", fmt.Sprintf("%v", *o.ReturnClientRequestId))
	}
	return &out
}

func (o JobGetAllLifetimeStatisticsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o JobGetAllLifetimeStatisticsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Timeout != nil {
		out.Append("timeout", fmt.Sprintf("%v", *o.Timeout))
	}
	return &out
}

// JobGetAllLifetimeStatistics ...
func (c JobsClient) JobGetAllLifetimeStatistics(ctx context.Context, options JobGetAllLifetimeStatisticsOperationOptions) (result JobGetAllLifetimeStatisticsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          "/lifetimejobstats",
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

	var model JobStatistics
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
