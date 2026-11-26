package openapis

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ServiceListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ServiceResource
}

type ServiceListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ServiceResource
}

type ServiceListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ServiceListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ServiceList ...
func (c OpenapisClient) ServiceList(ctx context.Context, id DatabaseAccountId) (result ServiceListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ServiceListCustomPager{},
		Path:       fmt.Sprintf("%s/services", id.ID()),
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
		Values *[]ServiceResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ServiceListComplete retrieves all the results into a single object
func (c OpenapisClient) ServiceListComplete(ctx context.Context, id DatabaseAccountId) (ServiceListCompleteResult, error) {
	return c.ServiceListCompleteMatchingPredicate(ctx, id, ServiceResourceOperationPredicate{})
}

// ServiceListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c OpenapisClient) ServiceListCompleteMatchingPredicate(ctx context.Context, id DatabaseAccountId, predicate ServiceResourceOperationPredicate) (result ServiceListCompleteResult, err error) {
	items := make([]ServiceResource, 0)

	resp, err := c.ServiceList(ctx, id)
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

	result = ServiceListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
