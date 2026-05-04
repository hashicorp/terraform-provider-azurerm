package hostpool

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListRegistrationTokensOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]RegistrationTokenMinimal
}

type ListRegistrationTokensCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []RegistrationTokenMinimal
}

type ListRegistrationTokensCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListRegistrationTokensCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListRegistrationTokens ...
func (c HostPoolClient) ListRegistrationTokens(ctx context.Context, id HostPoolId) (result ListRegistrationTokensOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPost,
		Pager:      &ListRegistrationTokensCustomPager{},
		Path:       fmt.Sprintf("%s/listRegistrationTokens", id.ID()),
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
		Values *[]RegistrationTokenMinimal `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListRegistrationTokensComplete retrieves all the results into a single object
func (c HostPoolClient) ListRegistrationTokensComplete(ctx context.Context, id HostPoolId) (ListRegistrationTokensCompleteResult, error) {
	return c.ListRegistrationTokensCompleteMatchingPredicate(ctx, id, RegistrationTokenMinimalOperationPredicate{})
}

// ListRegistrationTokensCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c HostPoolClient) ListRegistrationTokensCompleteMatchingPredicate(ctx context.Context, id HostPoolId, predicate RegistrationTokenMinimalOperationPredicate) (result ListRegistrationTokensCompleteResult, err error) {
	items := make([]RegistrationTokenMinimal, 0)

	resp, err := c.ListRegistrationTokens(ctx, id)
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

	result = ListRegistrationTokensCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
