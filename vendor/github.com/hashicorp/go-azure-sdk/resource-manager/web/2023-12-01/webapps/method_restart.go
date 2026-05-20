package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RestartOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type RestartOperationOptions struct {
	SoftRestart *bool
	Synchronous *bool
}

func DefaultRestartOperationOptions() RestartOperationOptions {
	return RestartOperationOptions{}
}

func (o RestartOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o RestartOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o RestartOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.SoftRestart != nil {
		out.Append("softRestart", fmt.Sprintf("%v", *o.SoftRestart))
	}
	if o.Synchronous != nil {
		out.Append("synchronous", fmt.Sprintf("%v", *o.Synchronous))
	}
	return &out
}

// Restart ...
func (c WebAppsClient) Restart(ctx context.Context, id commonids.AppServiceId, options RestartOperationOptions) (result RestartOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/restart", id.ID()),
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
