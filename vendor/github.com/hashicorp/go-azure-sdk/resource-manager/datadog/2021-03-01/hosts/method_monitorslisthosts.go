package hosts

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsListHostsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatadogHost
}

type MonitorsListHostsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DatadogHost
}

type MonitorsListHostsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListHostsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListHosts ...
func (c HostsClient) MonitorsListHosts(ctx context.Context, id MonitorId) (result MonitorsListHostsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &MonitorsListHostsCustomPager{},
		Path:       fmt.Sprintf("%s/listHosts", id.ID()),
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
		Values *[]DatadogHost `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MonitorsListHostsComplete retrieves all the results into a single object
func (c HostsClient) MonitorsListHostsComplete(ctx context.Context, id MonitorId) (MonitorsListHostsCompleteResult, error) {
	return c.MonitorsListHostsCompleteMatchingPredicate(ctx, id, DatadogHostOperationPredicate{})
}

// MonitorsListHostsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HostsClient) MonitorsListHostsCompleteMatchingPredicate(ctx context.Context, id MonitorId, predicate DatadogHostOperationPredicate) (result MonitorsListHostsCompleteResult, err error) {
	items := make([]DatadogHost, 0)

	resp, err := c.MonitorsListHosts(ctx, id)
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

	result = MonitorsListHostsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
