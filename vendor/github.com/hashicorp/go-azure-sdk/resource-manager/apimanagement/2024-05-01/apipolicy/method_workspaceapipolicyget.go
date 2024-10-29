package apipolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiPolicyGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *PolicyContract
}

type WorkspaceApiPolicyGetOperationOptions struct {
	Format *PolicyExportFormat
}

func DefaultWorkspaceApiPolicyGetOperationOptions() WorkspaceApiPolicyGetOperationOptions {
	return WorkspaceApiPolicyGetOperationOptions{}
}

func (o WorkspaceApiPolicyGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o WorkspaceApiPolicyGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceApiPolicyGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Format != nil {
		out.Append("format", fmt.Sprintf("%v", *o.Format))
	}
	return &out
}

// WorkspaceApiPolicyGet ...
func (c ApiPolicyClient) WorkspaceApiPolicyGet(ctx context.Context, id WorkspaceApiId, options WorkspaceApiPolicyGetOperationOptions) (result WorkspaceApiPolicyGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("%s/policies/policy", id.ID()),
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

	var model PolicyContract
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
