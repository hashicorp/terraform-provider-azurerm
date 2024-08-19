package deviceupdates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type InstancesListByAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Instance
}

type InstancesListByAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Instance
}

type InstancesListByAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *InstancesListByAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// InstancesListByAccount ...
func (c DeviceupdatesClient) InstancesListByAccount(ctx context.Context, id AccountId) (result InstancesListByAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &InstancesListByAccountCustomPager{},
		Path:       fmt.Sprintf("%s/instances", id.ID()),
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
		Values *[]Instance `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// InstancesListByAccountComplete retrieves all the results into a single object
func (c DeviceupdatesClient) InstancesListByAccountComplete(ctx context.Context, id AccountId) (InstancesListByAccountCompleteResult, error) {
	return c.InstancesListByAccountCompleteMatchingPredicate(ctx, id, InstanceOperationPredicate{})
}

// InstancesListByAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DeviceupdatesClient) InstancesListByAccountCompleteMatchingPredicate(ctx context.Context, id AccountId, predicate InstanceOperationPredicate) (result InstancesListByAccountCompleteResult, err error) {
	items := make([]Instance, 0)

	resp, err := c.InstancesListByAccount(ctx, id)
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

	result = InstancesListByAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
