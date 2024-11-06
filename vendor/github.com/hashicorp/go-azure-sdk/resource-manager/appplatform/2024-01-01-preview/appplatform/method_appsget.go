package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AppsGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AppResource
}

type AppsGetOperationOptions struct {
	SyncStatus *string
}

func DefaultAppsGetOperationOptions() AppsGetOperationOptions {
	return AppsGetOperationOptions{}
}

func (o AppsGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AppsGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AppsGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SyncStatus != nil {
		out.Append("syncStatus", fmt.Sprintf("%v", *o.SyncStatus))
	}
	return &out
}

// AppsGet ...
func (c AppPlatformClient) AppsGet(ctx context.Context, id AppId, options AppsGetOperationOptions) (result AppsGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          id.ID(),
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

	var model AppResource
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
