package diagnostic

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type WorkspaceDiagnosticDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type WorkspaceDiagnosticDeleteOperationOptions struct {
	IfMatch *string
}

func DefaultWorkspaceDiagnosticDeleteOperationOptions() WorkspaceDiagnosticDeleteOperationOptions {
	return WorkspaceDiagnosticDeleteOperationOptions{}
}

func (o WorkspaceDiagnosticDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o WorkspaceDiagnosticDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o WorkspaceDiagnosticDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// WorkspaceDiagnosticDelete ...
func (c DiagnosticClient) WorkspaceDiagnosticDelete(ctx context.Context, id WorkspaceDiagnosticId, options WorkspaceDiagnosticDeleteOperationOptions) (result WorkspaceDiagnosticDeleteOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusNoContent,
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
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

	return
}
