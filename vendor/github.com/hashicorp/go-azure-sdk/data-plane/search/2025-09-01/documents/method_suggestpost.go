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

type SuggestPostOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SuggestDocumentsResult
}

type SuggestPostOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultSuggestPostOperationOptions() SuggestPostOperationOptions {
	return SuggestPostOperationOptions{}
}

func (o SuggestPostOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o SuggestPostOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SuggestPostOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// SuggestPost ...
func (c DocumentsClient) SuggestPost(ctx context.Context, input SuggestRequest, options SuggestPostOperationOptions) (result SuggestPostOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          "/docs/Search.Post.Suggest",
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

	var model SuggestDocumentsResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
