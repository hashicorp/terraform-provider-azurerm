package associationsinterface

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByTrafficControllerOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Association
}

type ListByTrafficControllerCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Association
}

type ListByTrafficControllerCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByTrafficControllerCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByTrafficController ...
func (c AssociationsInterfaceClient) ListByTrafficController(ctx context.Context, id TrafficControllerId) (result ListByTrafficControllerOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByTrafficControllerCustomPager{},
		Path:       fmt.Sprintf("%s/associations", id.ID()),
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
		Values *[]Association `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByTrafficControllerComplete retrieves all the results into a single object
func (c AssociationsInterfaceClient) ListByTrafficControllerComplete(ctx context.Context, id TrafficControllerId) (ListByTrafficControllerCompleteResult, error) {
	return c.ListByTrafficControllerCompleteMatchingPredicate(ctx, id, AssociationOperationPredicate{})
}

// ListByTrafficControllerCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AssociationsInterfaceClient) ListByTrafficControllerCompleteMatchingPredicate(ctx context.Context, id TrafficControllerId, predicate AssociationOperationPredicate) (result ListByTrafficControllerCompleteResult, err error) {
	items := make([]Association, 0)

	resp, err := c.ListByTrafficController(ctx, id)
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

	result = ListByTrafficControllerCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
