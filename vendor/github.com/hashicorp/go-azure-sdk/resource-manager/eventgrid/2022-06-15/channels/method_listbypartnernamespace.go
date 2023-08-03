package channels

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByPartnerNamespaceOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Channel
}

type ListByPartnerNamespaceCompleteResult struct {
	Items []Channel
}

type ListByPartnerNamespaceOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByPartnerNamespaceOperationOptions() ListByPartnerNamespaceOperationOptions {
	return ListByPartnerNamespaceOperationOptions{}
}

func (o ListByPartnerNamespaceOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByPartnerNamespaceOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByPartnerNamespaceOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByPartnerNamespace ...
func (c ChannelsClient) ListByPartnerNamespace(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOperationOptions) (result ListByPartnerNamespaceOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/channels", id.ID()),
		OptionsObject: options,
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
		Values *[]Channel `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByPartnerNamespaceComplete retrieves all the results into a single object
func (c ChannelsClient) ListByPartnerNamespaceComplete(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOperationOptions) (ListByPartnerNamespaceCompleteResult, error) {
	return c.ListByPartnerNamespaceCompleteMatchingPredicate(ctx, id, options, ChannelOperationPredicate{})
}

// ListByPartnerNamespaceCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ChannelsClient) ListByPartnerNamespaceCompleteMatchingPredicate(ctx context.Context, id PartnerNamespaceId, options ListByPartnerNamespaceOperationOptions, predicate ChannelOperationPredicate) (result ListByPartnerNamespaceCompleteResult, err error) {
	items := make([]Channel, 0)

	resp, err := c.ListByPartnerNamespace(ctx, id, options)
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

	result = ListByPartnerNamespaceCompleteResult{
		Items: items,
	}
	return
}
