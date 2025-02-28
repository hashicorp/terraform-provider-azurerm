package apioperationsbytag

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type OperationListByTagsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagResourceContract
}

type OperationListByTagsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagResourceContract
}

type OperationListByTagsOperationOptions struct {
	Filter                     *string
	IncludeNotTaggedOperations *bool
	Skip                       *int64
	Top                        *int64
}

func DefaultOperationListByTagsOperationOptions() OperationListByTagsOperationOptions {
	return OperationListByTagsOperationOptions{}
}

func (o OperationListByTagsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o OperationListByTagsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o OperationListByTagsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IncludeNotTaggedOperations != nil {
		out.Append("includeNotTaggedOperations", fmt.Sprintf("%v", *o.IncludeNotTaggedOperations))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type OperationListByTagsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *OperationListByTagsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// OperationListByTags ...
func (c ApiOperationsByTagClient) OperationListByTags(ctx context.Context, id ApiId, options OperationListByTagsOperationOptions) (result OperationListByTagsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &OperationListByTagsCustomPager{},
		Path:          fmt.Sprintf("%s/operationsByTags", id.ID()),
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
		Values *[]TagResourceContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// OperationListByTagsComplete retrieves all the results into a single object
func (c ApiOperationsByTagClient) OperationListByTagsComplete(ctx context.Context, id ApiId, options OperationListByTagsOperationOptions) (OperationListByTagsCompleteResult, error) {
	return c.OperationListByTagsCompleteMatchingPredicate(ctx, id, options, TagResourceContractOperationPredicate{})
}

// OperationListByTagsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiOperationsByTagClient) OperationListByTagsCompleteMatchingPredicate(ctx context.Context, id ApiId, options OperationListByTagsOperationOptions, predicate TagResourceContractOperationPredicate) (result OperationListByTagsCompleteResult, err error) {
	items := make([]TagResourceContract, 0)

	resp, err := c.OperationListByTags(ctx, id, options)
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

	result = OperationListByTagsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
