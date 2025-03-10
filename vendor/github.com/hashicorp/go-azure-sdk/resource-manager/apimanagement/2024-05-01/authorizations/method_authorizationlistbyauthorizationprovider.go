package authorizations

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AuthorizationListByAuthorizationProviderOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AuthorizationContract
}

type AuthorizationListByAuthorizationProviderCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AuthorizationContract
}

type AuthorizationListByAuthorizationProviderOperationOptions struct {
	Filter *string
	Skip   *int64
	Top    *int64
}

func DefaultAuthorizationListByAuthorizationProviderOperationOptions() AuthorizationListByAuthorizationProviderOperationOptions {
	return AuthorizationListByAuthorizationProviderOperationOptions{}
}

func (o AuthorizationListByAuthorizationProviderOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AuthorizationListByAuthorizationProviderOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AuthorizationListByAuthorizationProviderOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Skip != nil {
		out.Append("$skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type AuthorizationListByAuthorizationProviderCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AuthorizationListByAuthorizationProviderCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AuthorizationListByAuthorizationProvider ...
func (c AuthorizationsClient) AuthorizationListByAuthorizationProvider(ctx context.Context, id AuthorizationProviderId, options AuthorizationListByAuthorizationProviderOperationOptions) (result AuthorizationListByAuthorizationProviderOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &AuthorizationListByAuthorizationProviderCustomPager{},
		Path:          fmt.Sprintf("%s/authorizations", id.ID()),
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
		Values *[]AuthorizationContract `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AuthorizationListByAuthorizationProviderComplete retrieves all the results into a single object
func (c AuthorizationsClient) AuthorizationListByAuthorizationProviderComplete(ctx context.Context, id AuthorizationProviderId, options AuthorizationListByAuthorizationProviderOperationOptions) (AuthorizationListByAuthorizationProviderCompleteResult, error) {
	return c.AuthorizationListByAuthorizationProviderCompleteMatchingPredicate(ctx, id, options, AuthorizationContractOperationPredicate{})
}

// AuthorizationListByAuthorizationProviderCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AuthorizationsClient) AuthorizationListByAuthorizationProviderCompleteMatchingPredicate(ctx context.Context, id AuthorizationProviderId, options AuthorizationListByAuthorizationProviderOperationOptions, predicate AuthorizationContractOperationPredicate) (result AuthorizationListByAuthorizationProviderCompleteResult, err error) {
	items := make([]AuthorizationContract, 0)

	resp, err := c.AuthorizationListByAuthorizationProvider(ctx, id, options)
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

	result = AuthorizationListByAuthorizationProviderCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
