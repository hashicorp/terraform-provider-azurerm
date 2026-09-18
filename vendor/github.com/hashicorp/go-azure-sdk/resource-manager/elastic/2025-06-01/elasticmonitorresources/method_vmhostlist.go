package elasticmonitorresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type VMHostListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]VMResources
}

type VMHostListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []VMResources
}

type VMHostListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *VMHostListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// VMHostList ...
func (c ElasticMonitorResourcesClient) VMHostList(ctx context.Context, id MonitorId) (result VMHostListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &VMHostListCustomPager{},
		Path:       fmt.Sprintf("%s/listVMHost", id.ID()),
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
		Values *[]VMResources `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// VMHostListComplete retrieves all the results into a single object
func (c ElasticMonitorResourcesClient) VMHostListComplete(ctx context.Context, id MonitorId) (VMHostListCompleteResult, error) {
	return c.VMHostListCompleteMatchingPredicate(ctx, id, VMResourcesOperationPredicate{})
}

// VMHostListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ElasticMonitorResourcesClient) VMHostListCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate VMResourcesOperationPredicate) (result VMHostListCompleteResult, err error) {
	items := make([]VMResources, 0)

	resp, err := c.VMHostList(ctx, id)
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

	result = VMHostListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
