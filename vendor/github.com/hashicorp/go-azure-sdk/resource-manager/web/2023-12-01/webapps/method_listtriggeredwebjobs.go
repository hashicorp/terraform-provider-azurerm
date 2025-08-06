package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListTriggeredWebJobsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]TriggeredWebJob
}

type ListTriggeredWebJobsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []TriggeredWebJob
}

type ListTriggeredWebJobsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListTriggeredWebJobsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListTriggeredWebJobs ...
func (c WebAppsClient) ListTriggeredWebJobs(ctx context.Context, id commonids.AppServiceId) (result ListTriggeredWebJobsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListTriggeredWebJobsCustomPager{},
		Path:       fmt.Sprintf("%s/triggeredWebJobs", id.ID()),
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
		Values *[]TriggeredWebJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListTriggeredWebJobsComplete retrieves all the results into a single object
func (c WebAppsClient) ListTriggeredWebJobsComplete(ctx context.Context, id commonids.AppServiceId) (ListTriggeredWebJobsCompleteResult, error) {
	return c.ListTriggeredWebJobsCompleteMatchingPredicate(ctx, id, TriggeredWebJobOperationPredicate{})
}

// ListTriggeredWebJobsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListTriggeredWebJobsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate TriggeredWebJobOperationPredicate) (result ListTriggeredWebJobsCompleteResult, err error) {
	items := make([]TriggeredWebJob, 0)

	resp, err := c.ListTriggeredWebJobs(ctx, id)
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

	result = ListTriggeredWebJobsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
