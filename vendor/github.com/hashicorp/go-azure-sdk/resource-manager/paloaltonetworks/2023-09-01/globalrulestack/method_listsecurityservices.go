package globalrulestack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListSecurityServicesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]SecurityServicesTypeList
}

type ListSecurityServicesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []SecurityServicesTypeList
}

type ListSecurityServicesOperationOptions struct {
	Skip *string
	Top  *int64
	Type *SecurityServicesTypeEnum
}

func DefaultListSecurityServicesOperationOptions() ListSecurityServicesOperationOptions {
	return ListSecurityServicesOperationOptions{}
}

func (o ListSecurityServicesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o ListSecurityServicesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o ListSecurityServicesOperationOptions) ToQuery() *client.QueryParams {
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

type ListSecurityServicesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListSecurityServicesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListSecurityServices ...
func (c GlobalRulestackClient) ListSecurityServices(ctx context.Context, id GlobalRulestackId, options ListSecurityServicesOperationOptions) (result ListSecurityServicesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPost,
		OptionsObject: options,
		Pager:         &ListSecurityServicesCustomPager{},
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

// ListSecurityServicesComplete retrieves all the results into a single object
func (c GlobalRulestackClient) ListSecurityServicesComplete(ctx context.Context, id GlobalRulestackId, options ListSecurityServicesOperationOptions) (ListSecurityServicesCompleteResult, error) {
	return c.ListSecurityServicesCompleteMatchingPredicate(ctx, id, options, SecurityServicesTypeListOperationPredicate{})
}

// ListSecurityServicesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c GlobalRulestackClient) ListSecurityServicesCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, options ListSecurityServicesOperationOptions, predicate SecurityServicesTypeListOperationPredicate) (result ListSecurityServicesCompleteResult, err error) {
	items := make([]SecurityServicesTypeList, 0)

	resp, err := c.ListSecurityServices(ctx, id, options)
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

	result = ListSecurityServicesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
