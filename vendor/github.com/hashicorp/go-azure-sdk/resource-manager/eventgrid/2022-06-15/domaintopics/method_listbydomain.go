package domaintopics

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByDomainOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DomainTopic
}

type ListByDomainCompleteResult struct {
	Items []DomainTopic
}

type ListByDomainOperationOptions struct {
	Filter *string
	Top    *int64
}

func DefaultListByDomainOperationOptions() ListByDomainOperationOptions {
	return ListByDomainOperationOptions{}
}

func (o ListByDomainOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListByDomainOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}
	return &out
}

func (o ListByDomainOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Top != nil {
		out.Append("$top", fmt.Sprintf("%v", *o.Top))
	}
	return &out
}

// ListByDomain ...
func (c DomainTopicsClient) ListByDomain(ctx context.Context, id DomainId, options ListByDomainOperationOptions) (result ListByDomainOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		Path:          fmt.Sprintf("%s/topics", id.ID()),
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
		Values *[]DomainTopic `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByDomainComplete retrieves all the results into a single object
func (c DomainTopicsClient) ListByDomainComplete(ctx context.Context, id DomainId, options ListByDomainOperationOptions) (ListByDomainCompleteResult, error) {
	return c.ListByDomainCompleteMatchingPredicate(ctx, id, options, DomainTopicOperationPredicate{})
}

// ListByDomainCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c DomainTopicsClient) ListByDomainCompleteMatchingPredicate(ctx context.Context, id DomainId, options ListByDomainOperationOptions, predicate DomainTopicOperationPredicate) (result ListByDomainCompleteResult, err error) {
	items := make([]DomainTopic, 0)

	resp, err := c.ListByDomain(ctx, id, options)
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

	result = ListByDomainCompleteResult{
		Items: items,
	}
	return
}
