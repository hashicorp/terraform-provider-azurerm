package restorables

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestorableGremlinGraphsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RestorableGremlinGraphsListResult
}

type RestorableGremlinGraphsListOperationOptions struct {
	EndTime                      *string
	RestorableGremlinDatabaseRid *string
	StartTime                    *string
}

func DefaultRestorableGremlinGraphsListOperationOptions() RestorableGremlinGraphsListOperationOptions {
	return RestorableGremlinGraphsListOperationOptions{}
}

func (o RestorableGremlinGraphsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableGremlinGraphsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableGremlinGraphsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.EndTime != nil {
		out.Append("endTime", fmt.Sprintf("%v", *o.EndTime))
	}
	if o.RestorableGremlinDatabaseRid != nil {
		out.Append("restorableGremlinDatabaseRid", fmt.Sprintf("%v", *o.RestorableGremlinDatabaseRid))
	}
	if o.StartTime != nil {
		out.Append("startTime", fmt.Sprintf("%v", *o.StartTime))
	}
	return &out
}

// RestorableGremlinGraphsList ...
func (c RestorablesClient) RestorableGremlinGraphsList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableGremlinGraphsListOperationOptions) (result RestorableGremlinGraphsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restorableGraphs", id.ID()),
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

	var model RestorableGremlinGraphsListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
