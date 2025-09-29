package apioperationtag

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type TagListByOperationOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TagContract
}

type TagListByOperationCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TagContract
}

type TagListByOperationOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultTagListByOperationOperationOptions() TagListByOperationOperationOptions {
	return TagListByOperationOperationOptions{}
}

func (o TagListByOperationOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o TagListByOperationOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o TagListByOperationOperationOptions) ToQuery() *client.QueryParams {
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

type TagListByOperationCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *TagListByOperationCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// TagListByOperation ...
func (c ApiOperationTagClient) TagListByOperation(ctx context.Context, id OperationId, options TagListByOperationOperationOptions) (result TagListByOperationOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &TagListByOperationCustomPager{},
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

// TagListByOperationComplete retrieves all the results into a single object
func (c ApiOperationTagClient) TagListByOperationComplete(ctx context.Context, id OperationId, options TagListByOperationOperationOptions) (TagListByOperationCompleteResult, error) {
	return c.TagListByOperationCompleteMatchingPredicate(ctx, id, options, TagContractOperationPredicate{})
}

// TagListByOperationCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ApiOperationTagClient) TagListByOperationCompleteMatchingPredicate(ctx context.Context, id OperationId, options TagListByOperationOperationOptions, predicate TagContractOperationPredicate) (result TagListByOperationCompleteResult, err error) {
	items := make([]TagContract, 0)

	resp, err := c.TagListByOperation(ctx, id, options)
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

	result = TagListByOperationCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
