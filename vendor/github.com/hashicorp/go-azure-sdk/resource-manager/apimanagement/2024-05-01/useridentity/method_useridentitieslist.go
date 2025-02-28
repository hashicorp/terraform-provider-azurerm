package useridentity

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type UserIdentitiesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]UserIdentityContract
}

type UserIdentitiesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []UserIdentityContract
}

type UserIdentitiesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *UserIdentitiesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// UserIdentitiesList ...
func (c UserIdentityClient) UserIdentitiesList(ctx context.Context, id UserId) (result UserIdentitiesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &UserIdentitiesListCustomPager{},
		Path:       fmt.Sprintf("%s/identities", id.ID()),
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
		Values *[]UserIdentityContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// UserIdentitiesListComplete retrieves all the results into a single object
func (c UserIdentityClient) UserIdentitiesListComplete(ctx context.Context, id UserId) (UserIdentitiesListCompleteResult, error) {
	return c.UserIdentitiesListCompleteMatchingPredicate(ctx, id, UserIdentityContractOperationPredicate{})
}

// UserIdentitiesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c UserIdentityClient) UserIdentitiesListCompleteMatchingPredicate(ctx context.Context, id UserId, predicate UserIdentityContractOperationPredicate) (result UserIdentitiesListCompleteResult, err error) {
	items := make([]UserIdentityContract, 0)

	resp, err := c.UserIdentitiesList(ctx, id)
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

	result = UserIdentitiesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
