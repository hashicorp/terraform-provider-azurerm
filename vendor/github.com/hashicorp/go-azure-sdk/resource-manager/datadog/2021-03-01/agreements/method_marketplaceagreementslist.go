package agreements

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

type MarketplaceAgreementsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DatadogAgreementResource
}

type MarketplaceAgreementsListCompleteResult struct {
	Items []DatadogAgreementResource
}

// MarketplaceAgreementsList ...
func (c AgreementsClient) MarketplaceAgreementsList(ctx context.Context, id commonids.SubscriptionId) (result MarketplaceAgreementsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/providers/Microsoft.Datadog/agreements", id.ID()),
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
		Values *[]DatadogAgreementResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// MarketplaceAgreementsListComplete retrieves all the results into a single object
func (c AgreementsClient) MarketplaceAgreementsListComplete(ctx context.Context, id commonids.SubscriptionId) (MarketplaceAgreementsListCompleteResult, error) {
	return c.MarketplaceAgreementsListCompleteMatchingPredicate(ctx, id, DatadogAgreementResourceOperationPredicate{})
}

// MarketplaceAgreementsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AgreementsClient) MarketplaceAgreementsListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, predicate DatadogAgreementResourceOperationPredicate) (result MarketplaceAgreementsListCompleteResult, err error) {
	items := make([]DatadogAgreementResource, 0)

	resp, err := c.MarketplaceAgreementsList(ctx, id)
	if err != nil {
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

	result = MarketplaceAgreementsListCompleteResult{
		Items: items,
	}
	return
}
