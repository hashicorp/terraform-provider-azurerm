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

type AutocompleteGetOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *AutocompleteResult
}

type AutocompleteGetOperationOptions struct {
	AutocompleteMode   *AutocompleteMode
	Filter             *string
	Fuzzy              *bool
	HighlightPostTag   *string
	HighlightPreTag    *string
	MinimumCoverage    *float64
	Search             *string
	SearchFields       *[]string
	SuggesterName      *string
	Top                *int64
	XMsClientRequestId *string
}

func DefaultAutocompleteGetOperationOptions() AutocompleteGetOperationOptions {
	return AutocompleteGetOperationOptions{}
}

func (o AutocompleteGetOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}
	if o.XMsClientRequestId != nil {
		out.Append("x-ms-client-request-id", fmt.Sprintf("%v", *o.XMsClientRequestId))
	}
	return &out
}

func (o AutocompleteGetOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AutocompleteGetOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.AutocompleteMode != nil {
		out.Append("autocompleteMode", fmt.Sprintf("%v", *o.AutocompleteMode))
	}
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
	if o.Search != nil {
		out.Append("search", fmt.Sprintf("%v", *o.Search))
	}
	if o.SearchFields != nil {
		out.Append("searchFields", fmt.Sprintf("%v", *o.SearchFields))
	}
	if o.SuggesterName != nil {
		out.Append("suggesterName", fmt.Sprintf("%v", *o.SuggesterName))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// AutocompleteGet ...
func (c DocumentsClient) AutocompleteGet(ctx context.Context, options AutocompleteGetOperationOptions) (result AutocompleteGetOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Path:          "/docs/Search.Autocomplete",
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

	var model AutocompleteResult
	result.Model = &model
	if err = resp.Unmarshal(result.Model); err != nil {
		return
	}

	return
}
