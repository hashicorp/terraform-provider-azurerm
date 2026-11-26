package accountconnectionresource

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AccountConnectionsListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ConnectionPropertiesV2BasicResource
}

type AccountConnectionsListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ConnectionPropertiesV2BasicResource
}

type AccountConnectionsListOperationOptions struct {
	Category   *string
	IncludeAll *bool
	Target     *string
}

func DefaultAccountConnectionsListOperationOptions() AccountConnectionsListOperationOptions {
	return AccountConnectionsListOperationOptions{}
}

func (o AccountConnectionsListOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o AccountConnectionsListOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o AccountConnectionsListOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Category != nil {
		out.Append("category", fmt.Sprintf("%v", *o.Category))
	}
	if o.IncludeAll != nil {
		out.Append("includeAll", fmt.Sprintf("%v", *o.IncludeAll))
	}
	if o.Target != nil {
		out.Append("target", fmt.Sprintf("%v", *o.Target))
	}
	return &out
}

type AccountConnectionsListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *AccountConnectionsListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// AccountConnectionsList ...
func (c AccountConnectionResourceClient) AccountConnectionsList(ctx context.Context, id AccountId, options AccountConnectionsListOperationOptions) (result AccountConnectionsListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &AccountConnectionsListCustomPager{},
		Path:          fmt.Sprintf("%s/connections", id.ID()),
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
		Values *[]ConnectionPropertiesV2BasicResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// AccountConnectionsListComplete retrieves all the results into a single object
func (c AccountConnectionResourceClient) AccountConnectionsListComplete(ctx context.Context, id AccountId, options AccountConnectionsListOperationOptions) (AccountConnectionsListCompleteResult, error) {
	return c.AccountConnectionsListCompleteMatchingPredicate(ctx, id, options, ConnectionPropertiesV2BasicResourceOperationPredicate{})
}

// AccountConnectionsListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AccountConnectionResourceClient) AccountConnectionsListCompleteMatchingPredicate(ctx context.Context, id AccountId, options AccountConnectionsListOperationOptions, predicate ConnectionPropertiesV2BasicResourceOperationPredicate) (result AccountConnectionsListCompleteResult, err error) {
	items := make([]ConnectionPropertiesV2BasicResource, 0)

	resp, err := c.AccountConnectionsList(ctx, id, options)
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

	result = AccountConnectionsListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
