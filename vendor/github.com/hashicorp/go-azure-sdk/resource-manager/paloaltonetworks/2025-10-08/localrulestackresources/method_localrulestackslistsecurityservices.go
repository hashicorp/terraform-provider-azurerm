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

type LocalRulestackslistSecurityServicesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SecurityServicesTypeList
}

type LocalRulestackslistSecurityServicesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SecurityServicesTypeList
}

type LocalRulestackslistSecurityServicesOperationOptions struct {
	Skip *string
	Top  *int64
	Type *SecurityServicesTypeEnum
}

func DefaultLocalRulestackslistSecurityServicesOperationOptions() LocalRulestackslistSecurityServicesOperationOptions {
	return LocalRulestackslistSecurityServicesOperationOptions{}
}

func (o LocalRulestackslistSecurityServicesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o LocalRulestackslistSecurityServicesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o LocalRulestackslistSecurityServicesOperationOptions) ToQuery() *client.QueryParams {
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

type LocalRulestackslistSecurityServicesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *LocalRulestackslistSecurityServicesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// LocalRulestackslistSecurityServices ...
func (c LocalRulestackResourcesClient) LocalRulestackslistSecurityServices(ctx context.Context, id LocalRulestackId, options LocalRulestackslistSecurityServicesOperationOptions) (result LocalRulestackslistSecurityServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &LocalRulestackslistSecurityServicesCustomPager{},
		Path:          fmt.Sprintf("%s/listSecurityServices", id.ID()),
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
		Values *[]SecurityServicesTypeList `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// LocalRulestackslistSecurityServicesComplete retrieves all the results into a single object
func (c LocalRulestackResourcesClient) LocalRulestackslistSecurityServicesComplete(ctx context.Context, id LocalRulestackId, options LocalRulestackslistSecurityServicesOperationOptions) (LocalRulestackslistSecurityServicesCompleteResult, error) {
	return c.LocalRulestackslistSecurityServicesCompleteMatchingPredicate(ctx, id, options, SecurityServicesTypeListOperationPredicate{})
}

// LocalRulestackslistSecurityServicesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c LocalRulestackResourcesClient) LocalRulestackslistSecurityServicesCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, options LocalRulestackslistSecurityServicesOperationOptions, predicate SecurityServicesTypeListOperationPredicate) (result LocalRulestackslistSecurityServicesCompleteResult, err error) {
	items := make([]SecurityServicesTypeList, 0)

	resp, err := c.LocalRulestackslistSecurityServices(ctx, id, options)
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

	result = LocalRulestackslistSecurityServicesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
