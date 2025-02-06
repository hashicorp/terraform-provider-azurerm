package resourceproviders

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

type ListCustomHostNameSitesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CustomHostnameSites
}

type ListCustomHostNameSitesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CustomHostnameSites
}

type ListCustomHostNameSitesOperationOptions struct {
	Hostname *string
}

func DefaultListCustomHostNameSitesOperationOptions() ListCustomHostNameSitesOperationOptions {
	return ListCustomHostNameSitesOperationOptions{}
}

func (o ListCustomHostNameSitesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListCustomHostNameSitesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListCustomHostNameSitesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Hostname != nil {
		out.Append("hostname", fmt.Sprintf("%v", *o.Hostname))
	}
	return &out
}

type ListCustomHostNameSitesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListCustomHostNameSitesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListCustomHostNameSites ...
func (c ResourceProvidersClient) ListCustomHostNameSites(ctx context.Context, id commonids.SubscriptionId, options ListCustomHostNameSitesOperationOptions) (result ListCustomHostNameSitesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &ListCustomHostNameSitesCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Web/customhostnameSites", id.ID()),
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
		Values *[]CustomHostnameSites `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListCustomHostNameSitesComplete retrieves all the results into a single object
func (c ResourceProvidersClient) ListCustomHostNameSitesComplete(ctx context.Context, id commonids.SubscriptionId, options ListCustomHostNameSitesOperationOptions) (ListCustomHostNameSitesCompleteResult, error) {
	return c.ListCustomHostNameSitesCompleteMatchingPredicate(ctx, id, options, CustomHostnameSitesOperationPredicate{})
}

// ListCustomHostNameSitesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ResourceProvidersClient) ListCustomHostNameSitesCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options ListCustomHostNameSitesOperationOptions, predicate CustomHostnameSitesOperationPredicate) (result ListCustomHostNameSitesCompleteResult, err error) {
	items := make([]CustomHostnameSites, 0)

	resp, err := c.ListCustomHostNameSites(ctx, id, options)
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

	result = ListCustomHostNameSitesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
