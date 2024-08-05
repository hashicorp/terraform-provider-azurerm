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

type ListWebJobsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]WebJob
}

type ListWebJobsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []WebJob
}

type ListWebJobsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListWebJobsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListWebJobs ...
func (c WebAppsClient) ListWebJobs(ctx context.Context, id commonids.AppServiceId) (result ListWebJobsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListWebJobsCustomPager{},
		Path:       fmt.Sprintf("%s/webJobs", id.ID()),
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
		Values *[]WebJob `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListWebJobsComplete retrieves all the results into a single object
func (c WebAppsClient) ListWebJobsComplete(ctx context.Context, id commonids.AppServiceId) (ListWebJobsCompleteResult, error) {
	return c.ListWebJobsCompleteMatchingPredicate(ctx, id, WebJobOperationPredicate{})
}

// ListWebJobsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListWebJobsCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, predicate WebJobOperationPredicate) (result ListWebJobsCompleteResult, err error) {
	items := make([]WebJob, 0)

	resp, err := c.ListWebJobs(ctx, id)
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

	result = ListWebJobsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
