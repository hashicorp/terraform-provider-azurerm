package localuseroperationgroup

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

type LocalUsersListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]LocalUser
}

type LocalUsersListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []LocalUser
}

type LocalUsersListOperationOptions struct {
	Filter      *string
	Include     *ListLocalUserIncludeParam
	Maxpagesize *int64
}

func DefaultLocalUsersListOperationOptions() LocalUsersListOperationOptions {
	return LocalUsersListOperationOptions{}
}

func (o LocalUsersListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocalUsersListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocalUsersListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Include != nil {
		out.Append("$include", fmt.Sprintf("%v", *o.Include))
	}
	if o.Maxpagesize != nil {
		out.Append("$maxpagesize", fmt.Sprintf("%v", *o.Maxpagesize))
	}
	return &out
}

type LocalUsersListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalUsersListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalUsersList ...
func (c LocalUserOperationGroupClient) LocalUsersList(ctx context.Context, id commonids.StorageAccountId, options LocalUsersListOperationOptions) (result LocalUsersListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &LocalUsersListCustomPager{},
		Path:          fmt.Sprintf("%s/localUsers", id.ID()),
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
		Values *[]LocalUser `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocalUsersListComplete retrieves all the results into a single object
func (c LocalUserOperationGroupClient) LocalUsersListComplete(ctx context.Context, id commonids.StorageAccountId, options LocalUsersListOperationOptions) (LocalUsersListCompleteResult, error) {
	return c.LocalUsersListCompleteMatchingPredicate(ctx, id, options, LocalUserOperationPredicate{})
}

// LocalUsersListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalUserOperationGroupClient) LocalUsersListCompleteMatchingPredicate(ctx context.Context, id commonids.StorageAccountId, options LocalUsersListOperationOptions, predicate LocalUserOperationPredicate) (result LocalUsersListCompleteResult, err error) {
	items := make([]LocalUser, 0)

	resp, err := c.LocalUsersList(ctx, id, options)
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

	result = LocalUsersListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
