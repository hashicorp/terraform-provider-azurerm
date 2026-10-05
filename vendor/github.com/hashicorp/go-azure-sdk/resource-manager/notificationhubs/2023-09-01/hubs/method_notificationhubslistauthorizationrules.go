package hubs

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type NotificationHubsListAuthorizationRulesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SharedAccessAuthorizationRuleResource
}

type NotificationHubsListAuthorizationRulesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SharedAccessAuthorizationRuleResource
}

type NotificationHubsListAuthorizationRulesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *NotificationHubsListAuthorizationRulesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// NotificationHubsListAuthorizationRules ...
func (c HubsClient) NotificationHubsListAuthorizationRules(ctx context.Context, id NotificationHubId) (result NotificationHubsListAuthorizationRulesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &NotificationHubsListAuthorizationRulesCustomPager{},
		Path:       fmt.Sprintf("%s/authorizationRules", id.ID()),
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
		Values *[]SharedAccessAuthorizationRuleResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// NotificationHubsListAuthorizationRulesComplete retrieves all the results into a single object
func (c HubsClient) NotificationHubsListAuthorizationRulesComplete(ctx context.Context, id NotificationHubId) (NotificationHubsListAuthorizationRulesCompleteResult, error) {
	return c.NotificationHubsListAuthorizationRulesCompleteMatchingPredicate(ctx, id, SharedAccessAuthorizationRuleResourceOperationPredicate{})
}

// NotificationHubsListAuthorizationRulesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HubsClient) NotificationHubsListAuthorizationRulesCompleteMatchingPredicate(ctx context.Context, id NotificationHubId, predicate SharedAccessAuthorizationRuleResourceOperationPredicate) (result NotificationHubsListAuthorizationRulesCompleteResult, err error) {
	items := make([]SharedAccessAuthorizationRuleResource, 0)

	resp, err := c.NotificationHubsListAuthorizationRules(ctx, id)
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

	result = NotificationHubsListAuthorizationRulesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
