package networkconnections

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListHealthDetailsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]HealthCheckStatusDetails
}

type ListHealthDetailsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []HealthCheckStatusDetails
}

type ListHealthDetailsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListHealthDetailsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListHealthDetails ...
func (c NetworkConnectionsClient) ListHealthDetails(ctx context.Context, id NetworkConnectionId) (result ListHealthDetailsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListHealthDetailsCustomPager{},
		Path:       fmt.Sprintf("%s/healthChecks", id.ID()),
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
		Values *[]HealthCheckStatusDetails `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListHealthDetailsComplete retrieves all the results into a single object
func (c NetworkConnectionsClient) ListHealthDetailsComplete(ctx context.Context, id NetworkConnectionId) (ListHealthDetailsCompleteResult, error) {
	return c.ListHealthDetailsCompleteMatchingPredicate(ctx, id, HealthCheckStatusDetailsOperationPredicate{})
}

// ListHealthDetailsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NetworkConnectionsClient) ListHealthDetailsCompleteMatchingPredicate(ctx context.Context, id NetworkConnectionId, predicate HealthCheckStatusDetailsOperationPredicate) (result ListHealthDetailsCompleteResult, err error) {
	items := make([]HealthCheckStatusDetails, 0)

	resp, err := c.ListHealthDetails(ctx, id)
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

	result = ListHealthDetailsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
