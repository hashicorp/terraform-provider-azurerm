package resourceguards

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetDeleteResourceGuardProxyRequestsObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DppBaseResource
}

type GetDeleteResourceGuardProxyRequestsObjectsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DppBaseResource
}

type GetDeleteResourceGuardProxyRequestsObjectsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDeleteResourceGuardProxyRequestsObjectsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDeleteResourceGuardProxyRequestsObjects ...
func (c ResourceGuardsClient) GetDeleteResourceGuardProxyRequestsObjects(ctx context.Context, id ResourceGuardId) (result GetDeleteResourceGuardProxyRequestsObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &GetDeleteResourceGuardProxyRequestsObjectsCustomPager{},
		Path:       fmt.Sprintf("%s/deleteResourceGuardProxyRequests", id.ID()),
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
		Values *[]DppBaseResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDeleteResourceGuardProxyRequestsObjectsComplete retrieves all the results into a single object
func (c ResourceGuardsClient) GetDeleteResourceGuardProxyRequestsObjectsComplete(ctx context.Context, id ResourceGuardId) (GetDeleteResourceGuardProxyRequestsObjectsCompleteResult, error) {
	return c.GetDeleteResourceGuardProxyRequestsObjectsCompleteMatchingPredicate(ctx, id, DppBaseResourceOperationPredicate{})
}

// GetDeleteResourceGuardProxyRequestsObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceGuardsClient) GetDeleteResourceGuardProxyRequestsObjectsCompleteMatchingPredicate(ctx context.Context, id ResourceGuardId, predicate DppBaseResourceOperationPredicate) (result GetDeleteResourceGuardProxyRequestsObjectsCompleteResult, err error) {
	items := make([]DppBaseResource, 0)

	resp, err := c.GetDeleteResourceGuardProxyRequestsObjects(ctx, id)
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

	result = GetDeleteResourceGuardProxyRequestsObjectsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
