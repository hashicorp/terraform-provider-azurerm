package streamingjobs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type SkuListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]GetStreamingJobSkuResult
}

type SkuListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []GetStreamingJobSkuResult
}

// SkuList ...
func (c StreamingJobsClient) SkuList(ctx context.Context, id StreamingJobId) (result SkuListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/skus", id.ID()),
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
		Values *[]GetStreamingJobSkuResult `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// SkuListComplete retrieves all the results into a single object
func (c StreamingJobsClient) SkuListComplete(ctx context.Context, id StreamingJobId) (SkuListCompleteResult, error) {
	return c.SkuListCompleteMatchingPredicate(ctx, id, GetStreamingJobSkuResultOperationPredicate{})
}

// SkuListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c StreamingJobsClient) SkuListCompleteMatchingPredicate(ctx context.Context, id StreamingJobId, predicate GetStreamingJobSkuResultOperationPredicate) (result SkuListCompleteResult, err error) {
	items := make([]GetStreamingJobSkuResult, 0)

	resp, err := c.SkuList(ctx, id)
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

	result = SkuListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
