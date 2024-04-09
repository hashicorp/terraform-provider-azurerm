package product

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceProductDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type WorkspaceProductDeleteOperationOptions struct {
	DeleteSubscriptions *bool
	IfMatch             *string
}

func DefaultWorkspaceProductDeleteOperationOptions() WorkspaceProductDeleteOperationOptions {
	return WorkspaceProductDeleteOperationOptions{}
}

func (o WorkspaceProductDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceProductDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o WorkspaceProductDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DeleteSubscriptions != nil {
		out.Append("deleteSubscriptions", fmt.Sprintf("%v", *o.DeleteSubscriptions))
	}
	return &out
}

// WorkspaceProductDelete ...
func (c ProductClient) WorkspaceProductDelete(ctx context.Context, id WorkspaceProductId, options WorkspaceProductDeleteOperationOptions) (result WorkspaceProductDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		Path:          id.ID(),
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
