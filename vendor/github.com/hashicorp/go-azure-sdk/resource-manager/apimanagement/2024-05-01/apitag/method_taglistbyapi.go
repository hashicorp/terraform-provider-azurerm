package apitag

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagListByApiOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagContract
}

type TagListByApiCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagContract
}

type TagListByApiOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultTagListByApiOperationOptions() TagListByApiOperationOptions {
	return TagListByApiOperationOptions{}
}

func (o TagListByApiOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TagListByApiOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o TagListByApiOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type TagListByApiCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TagListByApiCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TagListByApi ...
func (c ApiTagClient) TagListByApi(ctx context.Context, id ApiId, options TagListByApiOperationOptions) (result TagListByApiOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &TagListByApiCustomPager{},
		Path:          fmt.Sprintf("%s/tags", id.ID()),
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
		Values *[]TagContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// TagListByApiComplete retrieves all the results into a single object
func (c ApiTagClient) TagListByApiComplete(ctx context.Context, id ApiId, options TagListByApiOperationOptions) (TagListByApiCompleteResult, error) {
	return c.TagListByApiCompleteMatchingPredicate(ctx, id, options, TagContractOperationPredicate{})
}

// TagListByApiCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiTagClient) TagListByApiCompleteMatchingPredicate(ctx context.Context, id ApiId, options TagListByApiOperationOptions, predicate TagContractOperationPredicate) (result TagListByApiCompleteResult, err error) {
	items := make([]TagContract, 0)

	resp, err := c.TagListByApi(ctx, id, options)
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

	result = TagListByApiCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
