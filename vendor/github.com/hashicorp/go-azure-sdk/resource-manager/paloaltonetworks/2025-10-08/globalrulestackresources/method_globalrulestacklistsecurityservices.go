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

type GlobalRulestacklistSecurityServicesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SecurityServicesTypeList
}

type GlobalRulestacklistSecurityServicesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SecurityServicesTypeList
}

type GlobalRulestacklistSecurityServicesOperationOptions struct {
	Skip *string
	Top  *int64
	Type *SecurityServicesTypeEnum
}

func DefaultGlobalRulestacklistSecurityServicesOperationOptions() GlobalRulestacklistSecurityServicesOperationOptions {
	return GlobalRulestacklistSecurityServicesOperationOptions{}
}

func (o GlobalRulestacklistSecurityServicesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GlobalRulestacklistSecurityServicesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GlobalRulestacklistSecurityServicesOperationOptions) ToQuery() *client.QueryParams {
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

type GlobalRulestacklistSecurityServicesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GlobalRulestacklistSecurityServicesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GlobalRulestacklistSecurityServices ...
func (c GlobalRulestackResourcesClient) GlobalRulestacklistSecurityServices(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistSecurityServicesOperationOptions) (result GlobalRulestacklistSecurityServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &GlobalRulestacklistSecurityServicesCustomPager{},
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

// GlobalRulestacklistSecurityServicesComplete retrieves all the results into a single object
func (c GlobalRulestackResourcesClient) GlobalRulestacklistSecurityServicesComplete(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistSecurityServicesOperationOptions) (GlobalRulestacklistSecurityServicesCompleteResult, error) {
	return c.GlobalRulestacklistSecurityServicesCompleteMatchingPredicate(ctx, id, options, SecurityServicesTypeListOperationPredicate{})
}

// GlobalRulestacklistSecurityServicesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackResourcesClient) GlobalRulestacklistSecurityServicesCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, options GlobalRulestacklistSecurityServicesOperationOptions, predicate SecurityServicesTypeListOperationPredicate) (result GlobalRulestacklistSecurityServicesCompleteResult, err error) {
	items := make([]SecurityServicesTypeList, 0)

	resp, err := c.GlobalRulestacklistSecurityServices(ctx, id, options)
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

	result = GlobalRulestacklistSecurityServicesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
