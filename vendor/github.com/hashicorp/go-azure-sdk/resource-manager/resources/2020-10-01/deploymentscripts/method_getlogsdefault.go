package deploymentscripts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetLogsDefaultOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *ScriptLog
}

type GetLogsDefaultOperationOptions struct {
	Tail *int64
}

func DefaultGetLogsDefaultOperationOptions() GetLogsDefaultOperationOptions {
	return GetLogsDefaultOperationOptions{}
}

func (o GetLogsDefaultOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetLogsDefaultOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o GetLogsDefaultOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Tail != nil {
		out.Append("tail", fmt.Sprintf("%v", *o.Tail))
	}
	return &out
}

// GetLogsDefault ...
func (c DeploymentScriptsClient) GetLogsDefault(ctx context.Context, id DeploymentScriptId, options GetLogsDefaultOperationOptions) (result GetLogsDefaultOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/logs/default", id.ID()),
		OptionsObject: options,
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
