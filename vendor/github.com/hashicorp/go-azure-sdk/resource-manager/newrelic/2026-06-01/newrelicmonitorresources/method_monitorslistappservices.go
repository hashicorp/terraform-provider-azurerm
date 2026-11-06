package newrelicmonitorresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type MonitorsListAppServicesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AppServiceInfo
}

type MonitorsListAppServicesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AppServiceInfo
}

type MonitorsListAppServicesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *MonitorsListAppServicesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// MonitorsListAppServices ...
func (c NewRelicMonitorResourcesClient) MonitorsListAppServices(ctx context.Context, id MonitorId, input AppServicesGetRequest) (result MonitorsListAppServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &MonitorsListAppServicesCustomPager{},
		Path:       fmt.Sprintf("%s/listAppServices", id.ID()),
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
		Values *[]AppServiceInfo `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MonitorsListAppServicesComplete retrieves all the results into a single object
func (c NewRelicMonitorResourcesClient) MonitorsListAppServicesComplete(ctx context.Context, id MonitorId, input AppServicesGetRequest) (MonitorsListAppServicesCompleteResult, error) {
	return c.MonitorsListAppServicesCompleteMatchingPredicate(ctx, id, input, AppServiceInfoOperationPredicate{})
}

// MonitorsListAppServicesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NewRelicMonitorResourcesClient) MonitorsListAppServicesCompleteMatchingPredicate(ctx context.Context, id MonitorId, input AppServicesGetRequest, predicate AppServiceInfoOperationPredicate) (result MonitorsListAppServicesCompleteResult, err error) {
	items := make([]AppServiceInfo, 0)

	resp, err := c.MonitorsListAppServices(ctx, id, input)
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

	result = MonitorsListAppServicesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
