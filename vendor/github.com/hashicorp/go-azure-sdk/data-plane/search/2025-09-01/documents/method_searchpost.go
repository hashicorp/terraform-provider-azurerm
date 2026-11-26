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

type SearchPostOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SearchDocumentsResult
}

type SearchPostOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultSearchPostOperationOptions() SearchPostOperationOptions {
	return SearchPostOperationOptions{}
}

func (o SearchPostOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o SearchPostOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SearchPostOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// SearchPost ...
func (c DocumentsClient) SearchPost(ctx context.Context, input SearchRequest, options SearchPostOperationOptions) (result SearchPostOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
			http.StatusPartialContent,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          "/docs/Search.Post.Search",
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	if err = req.Marshal(input); err != nil {
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

	var model SearchDocumentsResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
