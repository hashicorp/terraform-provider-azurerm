package apioperationpolicy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceApiOperationPolicyDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type WorkspaceApiOperationPolicyDeleteOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceApiOperationPolicyDeleteOperationOptions() WorkspaceApiOperationPolicyDeleteOperationOptions {
	return WorkspaceApiOperationPolicyDeleteOperationOptions{}
}

func (o WorkspaceApiOperationPolicyDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceApiOperationPolicyDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceApiOperationPolicyDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceApiOperationPolicyDelete ...
func (c ApiOperationPolicyClient) WorkspaceApiOperationPolicyDelete(ctx context.Context, id ApiOperationId, options WorkspaceApiOperationPolicyDeleteOperationOptions) (result WorkspaceApiOperationPolicyDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		Path:          fmt.Sprintf("%s/policies/policy", id.ID()),
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

	return
}
