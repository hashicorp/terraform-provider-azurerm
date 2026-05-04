package linkedresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsListLinkedResourcesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LinkedResource
}

type MonitorsListLinkedResourcesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LinkedResource
}

type MonitorsListLinkedResourcesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListLinkedResourcesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListLinkedResources ...
func (c LinkedResourcesClient) MonitorsListLinkedResources(ctx context.Context, id MonitorId) (result MonitorsListLinkedResourcesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &MonitorsListLinkedResourcesCustomPager{},
		Path:       fmt.Sprintf("%s/listLinkedResources", id.ID()),
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
		Values *[]LinkedResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MonitorsListLinkedResourcesComplete retrieves all the results into a single object
func (c LinkedResourcesClient) MonitorsListLinkedResourcesComplete(ctx context.Context, id MonitorId) (MonitorsListLinkedResourcesCompleteResult, error) {
	return c.MonitorsListLinkedResourcesCompleteMatchingPredicate(ctx, id, LinkedResourceOperationPredicate{})
}

// MonitorsListLinkedResourcesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LinkedResourcesClient) MonitorsListLinkedResourcesCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate LinkedResourceOperationPredicate) (result MonitorsListLinkedResourcesCompleteResult, err error) {
	items := make([]LinkedResource, 0)

	resp, err := c.MonitorsListLinkedResources(ctx, id)
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

	result = MonitorsListLinkedResourcesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
