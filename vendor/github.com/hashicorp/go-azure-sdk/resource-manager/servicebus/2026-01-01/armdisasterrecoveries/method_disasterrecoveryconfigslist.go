package armdisasterrecoveries

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DisasterRecoveryConfigsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ArmDisasterRecovery
}

type DisasterRecoveryConfigsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ArmDisasterRecovery
}

type DisasterRecoveryConfigsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *DisasterRecoveryConfigsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// DisasterRecoveryConfigsList ...
func (c ArmDisasterRecoveriesClient) DisasterRecoveryConfigsList(ctx context.Context, id NamespaceId) (result DisasterRecoveryConfigsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &DisasterRecoveryConfigsListCustomPager{},
		Path:       fmt.Sprintf("%s/disasterRecoveryConfigs", id.ID()),
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
		Values *[]ArmDisasterRecovery `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// DisasterRecoveryConfigsListComplete retrieves all the results into a single object
func (c ArmDisasterRecoveriesClient) DisasterRecoveryConfigsListComplete(ctx context.Context, id NamespaceId) (DisasterRecoveryConfigsListCompleteResult, error) {
	return c.DisasterRecoveryConfigsListCompleteMatchingPredicate(ctx, id, ArmDisasterRecoveryOperationPredicate{})
}

// DisasterRecoveryConfigsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ArmDisasterRecoveriesClient) DisasterRecoveryConfigsListCompleteMatchingPredicate(ctx context.Context, id NamespaceId, predicate ArmDisasterRecoveryOperationPredicate) (result DisasterRecoveryConfigsListCompleteResult, err error) {
	items := make([]ArmDisasterRecovery, 0)

	resp, err := c.DisasterRecoveryConfigsList(ctx, id)
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

	result = DisasterRecoveryConfigsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
