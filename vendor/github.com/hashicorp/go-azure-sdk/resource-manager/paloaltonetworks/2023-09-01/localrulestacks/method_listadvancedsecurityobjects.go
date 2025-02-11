package localrulestacks

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListAdvancedSecurityObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AdvSecurityObjectModel
}

type ListAdvancedSecurityObjectsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AdvSecurityObjectModel
}

type ListAdvancedSecurityObjectsOperationOptions struct {
	Skip *string
	Top  *int64
	Type *AdvSecurityObjectTypeEnum
}

func DefaultListAdvancedSecurityObjectsOperationOptions() ListAdvancedSecurityObjectsOperationOptions {
	return ListAdvancedSecurityObjectsOperationOptions{}
}

func (o ListAdvancedSecurityObjectsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListAdvancedSecurityObjectsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListAdvancedSecurityObjectsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Skip != nil {
		out.Append("skip", fmt.Sprintf("%v", *o.Skip))
	}
	if o.Top != nil {
		out.Append("top", fmt.Sprintf("%v", *o.Top))
	}
	if o.Type != nil {
		out.Append("type", fmt.Sprintf("%v", *o.Type))
	}
	return &out
}

type ListAdvancedSecurityObjectsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListAdvancedSecurityObjectsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListAdvancedSecurityObjects ...
func (c LocalRulestacksClient) ListAdvancedSecurityObjects(ctx context.Context, id LocalRulestackId, options ListAdvancedSecurityObjectsOperationOptions) (result ListAdvancedSecurityObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListAdvancedSecurityObjectsCustomPager{},
		Path:          fmt.Sprintf("%s/listAdvancedSecurityObjects", id.ID()),
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
		Values *[]AdvSecurityObjectModel `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListAdvancedSecurityObjectsComplete retrieves all the results into a single object
func (c LocalRulestacksClient) ListAdvancedSecurityObjectsComplete(ctx context.Context, id LocalRulestackId, options ListAdvancedSecurityObjectsOperationOptions) (ListAdvancedSecurityObjectsCompleteResult, error) {
	return c.ListAdvancedSecurityObjectsCompleteMatchingPredicate(ctx, id, options, AdvSecurityObjectModelOperationPredicate{})
}

// ListAdvancedSecurityObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulestacksClient) ListAdvancedSecurityObjectsCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, options ListAdvancedSecurityObjectsOperationOptions, predicate AdvSecurityObjectModelOperationPredicate) (result ListAdvancedSecurityObjectsCompleteResult, err error) {
	items := make([]AdvSecurityObjectModel, 0)

	resp, err := c.ListAdvancedSecurityObjects(ctx, id, options)
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

	result = ListAdvancedSecurityObjectsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
