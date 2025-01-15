package createresource

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

type CreationSupportedListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CreateResourceSupportedResponse
}

type CreationSupportedListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CreateResourceSupportedResponse
}

type CreationSupportedListOperationOptions struct {
	DatadogOrganizationId *string
}

func DefaultCreationSupportedListOperationOptions() CreationSupportedListOperationOptions {
	return CreationSupportedListOperationOptions{}
}

func (o CreationSupportedListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CreationSupportedListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CreationSupportedListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.DatadogOrganizationId != nil {
		out.Append("datadogOrganizationId", fmt.Sprintf("%v", *o.DatadogOrganizationId))
	}
	return &out
}

type CreationSupportedListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CreationSupportedListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CreationSupportedList ...
func (c CreateResourceClient) CreationSupportedList(ctx context.Context, id commonids.SubscriptionId, options CreationSupportedListOperationOptions) (result CreationSupportedListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CreationSupportedListCustomPager{},
		Path:          fmt.Sprintf("%s/providers/Microsoft.Datadog/subscriptionStatuses", id.ID()),
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
		Values *[]CreateResourceSupportedResponse `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CreationSupportedListComplete retrieves all the results into a single object
func (c CreateResourceClient) CreationSupportedListComplete(ctx context.Context, id commonids.SubscriptionId, options CreationSupportedListOperationOptions) (CreationSupportedListCompleteResult, error) {
	return c.CreationSupportedListCompleteMatchingPredicate(ctx, id, options, CreateResourceSupportedResponseOperationPredicate{})
}

// CreationSupportedListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CreateResourceClient) CreationSupportedListCompleteMatchingPredicate(ctx context.Context, id commonids.SubscriptionId, options CreationSupportedListOperationOptions, predicate CreateResourceSupportedResponseOperationPredicate) (result CreationSupportedListCompleteResult, err error) {
	items := make([]CreateResourceSupportedResponse, 0)

	resp, err := c.CreationSupportedList(ctx, id, options)
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

	result = CreationSupportedListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
