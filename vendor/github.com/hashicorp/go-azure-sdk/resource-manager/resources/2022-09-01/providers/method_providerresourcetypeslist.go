package providers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ProviderResourceTypesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ProviderResourceType
}

type ProviderResourceTypesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ProviderResourceType
}

type ProviderResourceTypesListOperationOptions struct {
	Expand *string
}

func DefaultProviderResourceTypesListOperationOptions() ProviderResourceTypesListOperationOptions {
	return ProviderResourceTypesListOperationOptions{}
}

func (o ProviderResourceTypesListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ProviderResourceTypesListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ProviderResourceTypesListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Expand != nil {
		out.Append("$expand", fmt.Sprintf("%v", *o.Expand))
	}
	return &out
}

type ProviderResourceTypesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ProviderResourceTypesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ProviderResourceTypesList ...
func (c ProvidersClient) ProviderResourceTypesList(ctx context.Context, id SubscriptionProviderId, options ProviderResourceTypesListOperationOptions) (result ProviderResourceTypesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ProviderResourceTypesListCustomPager{},
		Path:          fmt.Sprintf("%s/resourceTypes", id.ID()),
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
		Values *[]ProviderResourceType `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ProviderResourceTypesListComplete retrieves all the results into a single object
func (c ProvidersClient) ProviderResourceTypesListComplete(ctx context.Context, id SubscriptionProviderId, options ProviderResourceTypesListOperationOptions) (ProviderResourceTypesListCompleteResult, error) {
	return c.ProviderResourceTypesListCompleteMatchingPredicate(ctx, id, options, ProviderResourceTypeOperationPredicate{})
}

// ProviderResourceTypesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ProvidersClient) ProviderResourceTypesListCompleteMatchingPredicate(ctx context.Context, id SubscriptionProviderId, options ProviderResourceTypesListOperationOptions, predicate ProviderResourceTypeOperationPredicate) (result ProviderResourceTypesListCompleteResult, err error) {
	items := make([]ProviderResourceType, 0)

	resp, err := c.ProviderResourceTypesList(ctx, id, options)
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

	result = ProviderResourceTypesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
