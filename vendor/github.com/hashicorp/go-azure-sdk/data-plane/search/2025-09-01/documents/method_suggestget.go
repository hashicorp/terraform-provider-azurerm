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

type SuggestGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *SuggestDocumentsResult
}

type SuggestGetOperationOptions struct {
	Filter             *string
	Fuzzy              *bool
	HighlightPostTag   *string
	HighlightPreTag    *string
	MinimumCoverage    *float64
	Orderby            *[]string
	Search             *string
	SearchFields       *[]string
	Select             *[]string
	SuggesterName      *string
	Top                *int64
	XMsClientRequestId *string
}

func DefaultSuggestGetOperationOptions() SuggestGetOperationOptions {
	return SuggestGetOperationOptions{}
}

func (o SuggestGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o SuggestGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o SuggestGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Fuzzy != nil {
		out.Append("fuzzy", fmt.Sprintf("%v", *o.Fuzzy))
	}
	if o.HighlightPostTag != nil {
		out.Append("highlightPostTag", fmt.Sprintf("%v", *o.HighlightPostTag))
	}
	if o.HighlightPreTag != nil {
		out.Append("highlightPreTag", fmt.Sprintf("%v", *o.HighlightPreTag))
	}
	if o.MinimumCoverage != nil {
		out.Append("minimumCoverage", fmt.Sprintf("%v", *o.MinimumCoverage))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Search != nil {
		out.Append("search", fmt.Sprintf("%v", *o.Search))
	}
	if o.SearchFields != nil {
		out.Append("searchFields", fmt.Sprintf("%v", *o.SearchFields))
	}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	if o.SuggesterName != nil {
		out.Append("suggesterName", fmt.Sprintf("%v", *o.SuggesterName))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// SuggestGet ...
func (c DocumentsClient) SuggestGet(ctx context.Context, options SuggestGetOperationOptions) (result SuggestGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          "/docs/Search.Suggest",
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

	var model SuggestDocumentsResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
