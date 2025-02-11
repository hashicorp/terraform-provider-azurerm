package monitors

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAppServicesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AppServiceInfo
}

type ListAppServicesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AppServiceInfo
}

type ListAppServicesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAppServicesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAppServices ...
func (c MonitorsClient) ListAppServices(ctx context.Context, id MonitorId, input AppServicesGetRequest) (result ListAppServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListAppServicesCustomPager{},
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

// ListAppServicesComplete retrieves all the results into a single object
func (c MonitorsClient) ListAppServicesComplete(ctx context.Context, id MonitorId, input AppServicesGetRequest) (ListAppServicesCompleteResult, error) {
	return c.ListAppServicesCompleteMatchingPredicate(ctx, id, input, AppServiceInfoOperationPredicate{})
}

// ListAppServicesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c MonitorsClient) ListAppServicesCompleteMatchingPredicate(ctx context.Context, id MonitorId, input AppServicesGetRequest, predicate AppServiceInfoOperationPredicate) (result ListAppServicesCompleteResult, err error) {
	items := make([]AppServiceInfo, 0)

	resp, err := c.ListAppServices(ctx, id, input)
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

	result = ListAppServicesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
