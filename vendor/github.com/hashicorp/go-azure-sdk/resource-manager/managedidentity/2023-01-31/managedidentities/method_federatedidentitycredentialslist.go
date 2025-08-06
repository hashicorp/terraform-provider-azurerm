package managedidentities

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

type FederatedIdentityCredentialsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]FederatedIdentityCredential
}

type FederatedIdentityCredentialsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []FederatedIdentityCredential
}

type FederatedIdentityCredentialsListOperationOptions struct {
	Top *int64
}

func DefaultFederatedIdentityCredentialsListOperationOptions() FederatedIdentityCredentialsListOperationOptions {
	return FederatedIdentityCredentialsListOperationOptions{}
}

func (o FederatedIdentityCredentialsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o FederatedIdentityCredentialsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o FederatedIdentityCredentialsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

type FederatedIdentityCredentialsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *FederatedIdentityCredentialsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// FederatedIdentityCredentialsList ...
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsList(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions) (result FederatedIdentityCredentialsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &FederatedIdentityCredentialsListCustomPager{},
		Path:          fmt.Sprintf("%s/federatedIdentityCredentials", id.ID()),
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
		Values *[]FederatedIdentityCredential `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// FederatedIdentityCredentialsListComplete retrieves all the results into a single object
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsListComplete(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions) (FederatedIdentityCredentialsListCompleteResult, error) {
	return c.FederatedIdentityCredentialsListCompleteMatchingPredicate(ctx, id, options, FederatedIdentityCredentialOperationPredicate{})
}

// FederatedIdentityCredentialsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedIdentitiesClient) FederatedIdentityCredentialsListCompleteMatchingPredicate(ctx context.Context, id commonids.UserAssignedIdentityId, options FederatedIdentityCredentialsListOperationOptions, predicate FederatedIdentityCredentialOperationPredicate) (result FederatedIdentityCredentialsListCompleteResult, err error) {
	items := make([]FederatedIdentityCredential, 0)

	resp, err := c.FederatedIdentityCredentialsList(ctx, id, options)
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

	result = FederatedIdentityCredentialsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
