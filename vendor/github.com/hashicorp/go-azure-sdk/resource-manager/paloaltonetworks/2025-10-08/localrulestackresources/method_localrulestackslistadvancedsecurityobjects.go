package localrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type LocalRulestackslistAdvancedSecurityObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AdvSecurityObjectModel
}

type LocalRulestackslistAdvancedSecurityObjectsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AdvSecurityObjectModel
}

type LocalRulestackslistAdvancedSecurityObjectsOperationOptions struct {
	Skip *string
	Top  *int64
	Type *AdvSecurityObjectTypeEnum
}

func DefaultLocalRulestackslistAdvancedSecurityObjectsOperationOptions() LocalRulestackslistAdvancedSecurityObjectsOperationOptions {
	return LocalRulestackslistAdvancedSecurityObjectsOperationOptions{}
}

func (o LocalRulestackslistAdvancedSecurityObjectsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocalRulestackslistAdvancedSecurityObjectsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocalRulestackslistAdvancedSecurityObjectsOperationOptions) ToQuery() *client.QueryParams {
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

type LocalRulestackslistAdvancedSecurityObjectsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulestackslistAdvancedSecurityObjectsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulestackslistAdvancedSecurityObjects ...
func (c LocalRulestackResourcesClient) LocalRulestackslistAdvancedSecurityObjects(ctx context.Context, id LocalRulestackId, options LocalRulestackslistAdvancedSecurityObjectsOperationOptions) (result LocalRulestackslistAdvancedSecurityObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &LocalRulestackslistAdvancedSecurityObjectsCustomPager{},
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

// LocalRulestackslistAdvancedSecurityObjectsComplete retrieves all the results into a single object
func (c LocalRulestackResourcesClient) LocalRulestackslistAdvancedSecurityObjectsComplete(ctx context.Context, id LocalRulestackId, options LocalRulestackslistAdvancedSecurityObjectsOperationOptions) (LocalRulestackslistAdvancedSecurityObjectsCompleteResult, error) {
	return c.LocalRulestackslistAdvancedSecurityObjectsCompleteMatchingPredicate(ctx, id, options, AdvSecurityObjectModelOperationPredicate{})
}

// LocalRulestackslistAdvancedSecurityObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulestackResourcesClient) LocalRulestackslistAdvancedSecurityObjectsCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, options LocalRulestackslistAdvancedSecurityObjectsOperationOptions, predicate AdvSecurityObjectModelOperationPredicate) (result LocalRulestackslistAdvancedSecurityObjectsCompleteResult, err error) {
	items := make([]AdvSecurityObjectModel, 0)

	resp, err := c.LocalRulestackslistAdvancedSecurityObjects(ctx, id, options)
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

	result = LocalRulestackslistAdvancedSecurityObjectsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
