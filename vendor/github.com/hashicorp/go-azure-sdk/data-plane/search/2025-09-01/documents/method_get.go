package documents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *interface{}
}

type GetOperationOptions struct {
	Select             *[]string
	XMsClientRequestId *string
}

func DefaultGetOperationOptions() GetOperationOptions {
	return GetOperationOptions{}
}

func (o GetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o GetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	return &out
}

// Get ...
func (c DocumentsClient) Get(ctx context.Context, id DocId, options GetOperationOptions) (result GetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          fmt.Sprintf("/docs('%s')", id.PathElements()...),
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

	var model interface{}
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
