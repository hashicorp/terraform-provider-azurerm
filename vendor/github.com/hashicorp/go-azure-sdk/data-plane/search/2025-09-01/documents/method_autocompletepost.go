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

type AutocompletePostOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AutocompleteResult
}

type AutocompletePostOperationOptions struct {
	XMsClientRequestId *string
}

func DefaultAutocompletePostOperationOptions() AutocompletePostOperationOptions {
	return AutocompletePostOperationOptions{}
}

func (o AutocompletePostOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o AutocompletePostOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AutocompletePostOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}

	return &out
}

// AutocompletePost ...
func (c DocumentsClient) AutocompletePost(ctx context.Context, input AutocompleteRequest, options AutocompletePostOperationOptions) (result AutocompletePostOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Path:          "/docs/Search.Post.Autocomplete",
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

	var model AutocompleteResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
