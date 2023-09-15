package identityprovider

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByServiceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]IdentityProviderContract
}

type ListByServiceCompleteResult struct {
	Items []IdentityProviderContract
}

// ListByService ...
func (c IdentityProviderClient) ListByService(ctx context.Context, id ServiceId) (result ListByServiceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/identityProviders", id.ID()),
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
		Values *[]IdentityProviderContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByServiceComplete retrieves all the results into a single object
func (c IdentityProviderClient) ListByServiceComplete(ctx context.Context, id ServiceId) (ListByServiceCompleteResult, error) {
	return c.ListByServiceCompleteMatchingPredicate(ctx, id, IdentityProviderContractOperationPredicate{})
}

// ListByServiceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c IdentityProviderClient) ListByServiceCompleteMatchingPredicate(ctx context.Context, id ServiceId, predicate IdentityProviderContractOperationPredicate) (result ListByServiceCompleteResult, err error) {
	items := make([]IdentityProviderContract, 0)

	resp, err := c.ListByService(ctx, id)
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

	result = ListByServiceCompleteResult{
		Items: items,
	}
	return
}
