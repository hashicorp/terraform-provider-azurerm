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

type RestorableTableResourcesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *RestorableTableResourcesListResult
}

type RestorableTableResourcesListOperationOptions struct {
	RestoreLocation       *string
	RestoreTimestampInUtc *string
}

func DefaultRestorableTableResourcesListOperationOptions() RestorableTableResourcesListOperationOptions {
	return RestorableTableResourcesListOperationOptions{}
}

func (o RestorableTableResourcesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestorableTableResourcesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestorableTableResourcesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.RestoreLocation != nil {
		out.Append("restoreLocation", fmt.Sprintf("%v", *o.RestoreLocation))
	}
	if o.RestoreTimestampInUtc != nil {
		out.Append("restoreTimestampInUtc", fmt.Sprintf("%v", *o.RestoreTimestampInUtc))
	}
	return &out
}

// RestorableTableResourcesList ...
func (c RestorablesClient) RestorableTableResourcesList(ctx context.Context, id RestorableDatabaseAccountId, options RestorableTableResourcesListOperationOptions) (result RestorableTableResourcesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restorableTableResources", id.ID()),
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

	var model RestorableTableResourcesListResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
