package globalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GlobalRulestacklistAdvancedSecurityObjectsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]AdvSecurityObjectModel
}

type GlobalRulestacklistAdvancedSecurityObjectsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []AdvSecurityObjectModel
}

type GlobalRulestacklistAdvancedSecurityObjectsOperationOptions struct {
	Skip *string
	Top  *int64
	Type *AdvSecurityObjectTypeEnum
}

func DefaultGlobalRulestacklistAdvancedSecurityObjectsOperationOptions() GlobalRulestacklistAdvancedSecurityObjectsOperationOptions {
	return GlobalRulestacklistAdvancedSecurityObjectsOperationOptions{}
}

func (o GlobalRulestacklistAdvancedSecurityObjectsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GlobalRulestacklistAdvancedSecurityObjectsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GlobalRulestacklistAdvancedSecurityObjectsOperationOptions) ToQuery() *client.QueryParams {
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

type GlobalRulestacklistAdvancedSecurityObjectsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GlobalRulestacklistAdvancedSecurityObjectsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GlobalRulestacklistAdvancedSecurityObjects ...
func (c GlobalRulestackResourcesClient) GlobalRulestacklistAdvancedSecurityObjects(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistAdvancedSecurityObjectsOperationOptions) (result GlobalRulestacklistAdvancedSecurityObjectsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &GlobalRulestacklistAdvancedSecurityObjectsCustomPager{},
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

// GlobalRulestacklistAdvancedSecurityObjectsComplete retrieves all the results into a single object
func (c GlobalRulestackResourcesClient) GlobalRulestacklistAdvancedSecurityObjectsComplete(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistAdvancedSecurityObjectsOperationOptions) (GlobalRulestacklistAdvancedSecurityObjectsCompleteResult, error) {
	return c.GlobalRulestacklistAdvancedSecurityObjectsCompleteMatchingPredicate(ctx, id, options, AdvSecurityObjectModelOperationPredicate{})
}

// GlobalRulestacklistAdvancedSecurityObjectsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackResourcesClient) GlobalRulestacklistAdvancedSecurityObjectsCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistAdvancedSecurityObjectsOperationOptions, predicate AdvSecurityObjectModelOperationPredicate) (result GlobalRulestacklistAdvancedSecurityObjectsCompleteResult, err error) {
	items := make([]AdvSecurityObjectModel, 0)

	resp, err := c.GlobalRulestacklistAdvancedSecurityObjects(ctx, id, options)
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

	result = GlobalRulestacklistAdvancedSecurityObjectsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
