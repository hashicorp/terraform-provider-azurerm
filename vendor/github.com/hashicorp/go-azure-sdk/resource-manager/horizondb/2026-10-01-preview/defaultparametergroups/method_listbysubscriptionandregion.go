package defaultparametergroups

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListBySubscriptionAndRegionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DefaultParameterGroup
}

type ListBySubscriptionAndRegionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DefaultParameterGroup
}

type ListBySubscriptionAndRegionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListBySubscriptionAndRegionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListBySubscriptionAndRegion ...
func (c DefaultParameterGroupsClient) ListBySubscriptionAndRegion(ctx context.Context, id LocationId) (result ListBySubscriptionAndRegionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListBySubscriptionAndRegionCustomPager{},
		Path:       fmt.Sprintf("%s/parameterGroups", id.ID()),
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
		Values *[]DefaultParameterGroup `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListBySubscriptionAndRegionComplete retrieves all the results into a single object
func (c DefaultParameterGroupsClient) ListBySubscriptionAndRegionComplete(ctx context.Context, id LocationId) (ListBySubscriptionAndRegionCompleteResult, error) {
	return c.ListBySubscriptionAndRegionCompleteMatchingPredicate(ctx, id, DefaultParameterGroupOperationPredicate{})
}

// ListBySubscriptionAndRegionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DefaultParameterGroupsClient) ListBySubscriptionAndRegionCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate DefaultParameterGroupOperationPredicate) (result ListBySubscriptionAndRegionCompleteResult, err error) {
	items := make([]DefaultParameterGroup, 0)

	resp, err := c.ListBySubscriptionAndRegion(ctx, id)
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

	result = ListBySubscriptionAndRegionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
