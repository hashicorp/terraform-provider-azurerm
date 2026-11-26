package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CollectionPartitionListUsagesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PartitionUsage
}

type CollectionPartitionListUsagesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PartitionUsage
}

type CollectionPartitionListUsagesOperationOptions struct {
	Filter *string
}

func DefaultCollectionPartitionListUsagesOperationOptions() CollectionPartitionListUsagesOperationOptions {
	return CollectionPartitionListUsagesOperationOptions{}
}

func (o CollectionPartitionListUsagesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CollectionPartitionListUsagesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CollectionPartitionListUsagesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	return &out
}

type CollectionPartitionListUsagesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CollectionPartitionListUsagesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CollectionPartitionListUsages ...
func (c OpenapisClient) CollectionPartitionListUsages(ctx context.Context, id CollectionId, options CollectionPartitionListUsagesOperationOptions) (result CollectionPartitionListUsagesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CollectionPartitionListUsagesCustomPager{},
		Path:          fmt.Sprintf("%s/partitions/usages", id.ID()),
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
		Values *[]PartitionUsage `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CollectionPartitionListUsagesComplete retrieves all the results into a single object
func (c OpenapisClient) CollectionPartitionListUsagesComplete(ctx context.Context, id CollectionId, options CollectionPartitionListUsagesOperationOptions) (CollectionPartitionListUsagesCompleteResult, error) {
	return c.CollectionPartitionListUsagesCompleteMatchingPredicate(ctx, id, options, PartitionUsageOperationPredicate{})
}

// CollectionPartitionListUsagesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) CollectionPartitionListUsagesCompleteMatchingPredicate(ctx context.Context, id CollectionId, options CollectionPartitionListUsagesOperationOptions, predicate PartitionUsageOperationPredicate) (result CollectionPartitionListUsagesCompleteResult, err error) {
	items := make([]PartitionUsage, 0)

	resp, err := c.CollectionPartitionListUsages(ctx, id, options)
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

	result = CollectionPartitionListUsagesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
