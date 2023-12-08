package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomDomainsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CustomDomain
}

type CustomDomainsListCompleteResult struct {
	Items []CustomDomain
}

// CustomDomainsList ...
func (c SignalRClient) CustomDomainsList(ctx context.Context, id SignalRId) (result CustomDomainsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/customDomains", id.ID()),
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
		Values *[]CustomDomain `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CustomDomainsListComplete retrieves all the results into a single object
func (c SignalRClient) CustomDomainsListComplete(ctx context.Context, id SignalRId) (CustomDomainsListCompleteResult, error) {
	return c.CustomDomainsListCompleteMatchingPredicate(ctx, id, CustomDomainOperationPredicate{})
}

// CustomDomainsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SignalRClient) CustomDomainsListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate CustomDomainOperationPredicate) (result CustomDomainsListCompleteResult, err error) {
	items := make([]CustomDomain, 0)

	resp, err := c.CustomDomainsList(ctx, id)
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

	result = CustomDomainsListCompleteResult{
		Items: items,
	}
	return
}
