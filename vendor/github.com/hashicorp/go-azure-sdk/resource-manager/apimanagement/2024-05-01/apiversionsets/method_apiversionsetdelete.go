package apiversionsets

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiVersionSetDeleteOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
}

type ApiVersionSetDeleteOperationOptions struct {
	IfMatch *string
}

func DefaultApiVersionSetDeleteOperationOptions() ApiVersionSetDeleteOperationOptions {
	return ApiVersionSetDeleteOperationOptions{}
}

func (o ApiVersionSetDeleteOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.IfMatch != nil {
		out.Append("If-Match", fmt.Sprintf("%v", *o.IfMatch))
	}
	return &out
}

func (o ApiVersionSetDeleteOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ApiVersionSetDeleteOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// ApiVersionSetDelete ...
func (c ApiVersionSetsClient) ApiVersionSetDelete(ctx context.Context, id ApiVersionSetId, options ApiVersionSetDeleteOperationOptions) (result ApiVersionSetDeleteOperationResponse, err error) {
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
