package apisbytag

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ApiListByTagsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagResourceContract
}

type ApiListByTagsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagResourceContract
}

type ApiListByTagsOperationOptions struct {
	Filter               *string
	IncludeNotTaggedApis *bool
	Skip                 *int64
	Top                  *int64
}

func DefaultApiListByTagsOperationOptions() ApiListByTagsOperationOptions {
	return ApiListByTagsOperationOptions{}
}

func (o ApiListByTagsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ApiListByTagsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ApiListByTagsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.IncludeNotTaggedApis != nil {
		out.Append("includeNotTaggedApis", fmt.Sprintf("%v", *o.IncludeNotTaggedApis))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type ApiListByTagsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ApiListByTagsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ApiListByTags ...
func (c ApisByTagClient) ApiListByTags(ctx context.Context, id ServiceId, options ApiListByTagsOperationOptions) (result ApiListByTagsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ApiListByTagsCustomPager{},
		Path:          fmt.Sprintf("%s/apisByTags", id.ID()),
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

// ApiListByTagsComplete retrieves all the results into a single object
func (c ApisByTagClient) ApiListByTagsComplete(ctx context.Context, id ServiceId, options ApiListByTagsOperationOptions) (ApiListByTagsCompleteResult, error) {
	return c.ApiListByTagsCompleteMatchingPredicate(ctx, id, options, TagResourceContractOperationPredicate{})
}

// ApiListByTagsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApisByTagClient) ApiListByTagsCompleteMatchingPredicate(ctx context.Context, id ServiceId, options ApiListByTagsOperationOptions, predicate TagResourceContractOperationPredicate) (result ApiListByTagsCompleteResult, err error) {
	items := make([]TagResourceContract, 0)

	resp, err := c.ApiListByTags(ctx, id, options)
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

	result = ApiListByTagsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
