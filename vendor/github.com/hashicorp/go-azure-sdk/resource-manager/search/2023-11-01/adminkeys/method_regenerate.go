package adminkeys

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type RegenerateOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AdminKeyResult
}

type RegenerateOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultRegenerateOperationOptions() RegenerateOperationOptions {
	return RegenerateOperationOptions{}
}

func (o RegenerateOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o RegenerateOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o RegenerateOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// Regenerate ...
func (c AdminKeysClient) Regenerate(ctx context.Context, id KeyKindId, options RegenerateOperationOptions) (result RegenerateOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
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

	if err = resp.Unmarshal(&result.Model); err != nil {
		return
	}

	return
}
