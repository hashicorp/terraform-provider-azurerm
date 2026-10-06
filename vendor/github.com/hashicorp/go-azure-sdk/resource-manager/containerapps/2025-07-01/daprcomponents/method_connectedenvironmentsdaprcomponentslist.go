package daprcomponents

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ConnectedEnvironmentsDaprComponentsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DaprComponent
}

type ConnectedEnvironmentsDaprComponentsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DaprComponent
}

type ConnectedEnvironmentsDaprComponentsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ConnectedEnvironmentsDaprComponentsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ConnectedEnvironmentsDaprComponentsList ...
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsList(ctx context.Context, id ConnectedEnvironmentId) (result ConnectedEnvironmentsDaprComponentsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ConnectedEnvironmentsDaprComponentsListCustomPager{},
		Path:       fmt.Sprintf("%s/daprComponents", id.ID()),
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
		Values *[]DaprComponent `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ConnectedEnvironmentsDaprComponentsListComplete retrieves all the results into a single object
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsListComplete(ctx context.Context, id ConnectedEnvironmentId) (ConnectedEnvironmentsDaprComponentsListCompleteResult, error) {
	return c.ConnectedEnvironmentsDaprComponentsListCompleteMatchingPredicate(ctx, id, DaprComponentOperationPredicate{})
}

// ConnectedEnvironmentsDaprComponentsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DaprComponentsClient) ConnectedEnvironmentsDaprComponentsListCompleteMatchingPredicate(ctx context.Context, id ConnectedEnvironmentId, predicate DaprComponentOperationPredicate) (result ConnectedEnvironmentsDaprComponentsListCompleteResult, err error) {
	items := make([]DaprComponent, 0)

	resp, err := c.ConnectedEnvironmentsDaprComponentsList(ctx, id)
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

	result = ConnectedEnvironmentsDaprComponentsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
