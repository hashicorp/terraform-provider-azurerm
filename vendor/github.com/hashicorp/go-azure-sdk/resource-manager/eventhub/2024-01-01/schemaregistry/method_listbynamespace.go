package schemaregistry

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByNamespaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SchemaGroup
}

type ListByNamespaceCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SchemaGroup
}

type ListByNamespaceOperationOptions struct {
	Skip *int64
	Top  *int64
}

func DefaultListByNamespaceOperationOptions() ListByNamespaceOperationOptions {
	return ListByNamespaceOperationOptions{}
}

func (o ListByNamespaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByNamespaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListByNamespaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ListByNamespaceCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByNamespaceCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByNamespace ...
func (c SchemaRegistryClient) ListByNamespace(ctx context.Context, id NamespaceId, options ListByNamespaceOperationOptions) (result ListByNamespaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListByNamespaceCustomPager{},
		Path:          fmt.Sprintf("%s/schemaGroups", id.ID()),
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
		Values *[]SchemaGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByNamespaceComplete retrieves all the results into a single object
func (c SchemaRegistryClient) ListByNamespaceComplete(ctx context.Context, id NamespaceId, options ListByNamespaceOperationOptions) (ListByNamespaceCompleteResult, error) {
	return c.ListByNamespaceCompleteMatchingPredicate(ctx, id, options, SchemaGroupOperationPredicate{})
}

// ListByNamespaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SchemaRegistryClient) ListByNamespaceCompleteMatchingPredicate(ctx context.Context, id NamespaceId, options ListByNamespaceOperationOptions, predicate SchemaGroupOperationPredicate) (result ListByNamespaceCompleteResult, err error) {
	items := make([]SchemaGroup, 0)

	resp, err := c.ListByNamespace(ctx, id, options)
	if err != nil {
		result.LatestHttpResponse = resp.HttpResponse
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

	result = ListByNamespaceCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
