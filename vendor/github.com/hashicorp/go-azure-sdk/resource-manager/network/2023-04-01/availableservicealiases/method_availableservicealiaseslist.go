package availableservicealiases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AvailableServiceAliasesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AvailableServiceAlias
}

type AvailableServiceAliasesListCompleteResult struct {
	Items []AvailableServiceAlias
}

// AvailableServiceAliasesList ...
func (c AvailableServiceAliasesClient) AvailableServiceAliasesList(ctx context.Context, id LocationId) (result AvailableServiceAliasesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Path:       fmt.Sprintf("%s/availableServiceAliases", id.ID()),
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
		Values *[]AvailableServiceAlias `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AvailableServiceAliasesListComplete retrieves all the results into a single object
func (c AvailableServiceAliasesClient) AvailableServiceAliasesListComplete(ctx context.Context, id LocationId) (AvailableServiceAliasesListCompleteResult, error) {
	return c.AvailableServiceAliasesListCompleteMatchingPredicate(ctx, id, AvailableServiceAliasOperationPredicate{})
}

// AvailableServiceAliasesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AvailableServiceAliasesClient) AvailableServiceAliasesListCompleteMatchingPredicate(ctx context.Context, id LocationId, predicate AvailableServiceAliasOperationPredicate) (result AvailableServiceAliasesListCompleteResult, err error) {
	items := make([]AvailableServiceAlias, 0)

	resp, err := c.AvailableServiceAliasesList(ctx, id)
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

	result = AvailableServiceAliasesListCompleteResult{
		Items: items,
	}
	return
}
