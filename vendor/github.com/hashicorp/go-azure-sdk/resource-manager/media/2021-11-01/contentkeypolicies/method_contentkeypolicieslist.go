package contentkeypolicies

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ContentKeyPoliciesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ContentKeyPolicy
}

type ContentKeyPoliciesListCompleteResult struct {
	Items []ContentKeyPolicy
}

type ContentKeyPoliciesListOperationOptions struct {
	Filter  *string
	Orderby *string
	Top     *int64
}

func DefaultContentKeyPoliciesListOperationOptions() ContentKeyPoliciesListOperationOptions {
	return ContentKeyPoliciesListOperationOptions{}
}

func (o ContentKeyPoliciesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ContentKeyPoliciesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ContentKeyPoliciesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Orderby != nil {
		out.Append("$orderby", fmt.Sprintf("%v", *o.Orderby))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ContentKeyPoliciesList ...
func (c ContentKeyPoliciesClient) ContentKeyPoliciesList(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions) (result ContentKeyPoliciesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/contentKeyPolicies", id.ID()),
		OptionsObject: options,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		return
	}

	var resp *client.Response
	resp, err = req.ExecutePaged(ctx)
	if resp != nil {
		result.OData = resp.OData
		result.HttpResponse = resp.Response
	}
	if err != nil {
		return
	}

	var values struct {
		Values *[]ContentKeyPolicy `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ContentKeyPoliciesListComplete retrieves all the results into a single object
func (c ContentKeyPoliciesClient) ContentKeyPoliciesListComplete(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions) (ContentKeyPoliciesListCompleteResult, error) {
	return c.ContentKeyPoliciesListCompleteMatchingPredicate(ctx, id, options, ContentKeyPolicyOperationPredicate{})
}

// ContentKeyPoliciesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ContentKeyPoliciesClient) ContentKeyPoliciesListCompleteMatchingPredicate(ctx context.Context, id MediaServiceId, options ContentKeyPoliciesListOperationOptions, predicate ContentKeyPolicyOperationPredicate) (result ContentKeyPoliciesListCompleteResult, err error) {
	items := make([]ContentKeyPolicy, 0)

	resp, err := c.ContentKeyPoliciesList(ctx, id, options)
	if err != nil {
		err = fmt.Errorf("loading results: %+v", err)
		return
	}
	if resp.Model != nil {
		for _, v := range *resp.Model {
			if predicate.Matches(v) {
				items = append(items, v)
			}
		}
	}

	result = ContentKeyPoliciesListCompleteResult{
		Items: items,
	}
	return
}
