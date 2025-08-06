package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSlotDifferencesFromProductionOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SlotDifference
}

type ListSlotDifferencesFromProductionCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SlotDifference
}

type ListSlotDifferencesFromProductionCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSlotDifferencesFromProductionCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSlotDifferencesFromProduction ...
func (c WebAppsClient) ListSlotDifferencesFromProduction(ctx context.Context, id commonids.AppServiceId, input CsmSlotEntity) (result ListSlotDifferencesFromProductionOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListSlotDifferencesFromProductionCustomPager{},
		Path:       fmt.Sprintf("%s/slotsdiffs", id.ID()),
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
		Values *[]SlotDifference `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListSlotDifferencesFromProductionComplete retrieves all the results into a single object
func (c WebAppsClient) ListSlotDifferencesFromProductionComplete(ctx context.Context, id commonids.AppServiceId, input CsmSlotEntity) (ListSlotDifferencesFromProductionCompleteResult, error) {
	return c.ListSlotDifferencesFromProductionCompleteMatchingPredicate(ctx, id, input, SlotDifferenceOperationPredicate{})
}

// ListSlotDifferencesFromProductionCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListSlotDifferencesFromProductionCompleteMatchingPredicate(ctx context.Context, id commonids.AppServiceId, input CsmSlotEntity, predicate SlotDifferenceOperationPredicate) (result ListSlotDifferencesFromProductionCompleteResult, err error) {
	items := make([]SlotDifference, 0)

	resp, err := c.ListSlotDifferencesFromProduction(ctx, id, input)
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

	result = ListSlotDifferencesFromProductionCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
