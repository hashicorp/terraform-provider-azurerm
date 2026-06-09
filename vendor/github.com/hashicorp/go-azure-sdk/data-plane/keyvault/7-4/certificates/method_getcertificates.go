package certificates

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetCertificatesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateItem
}

type GetCertificatesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateItem
}

type GetCertificatesOperationOptions struct {
	IncludePending *bool
	Maxresults     *int64
}

func DefaultGetCertificatesOperationOptions() GetCertificatesOperationOptions {
	return GetCertificatesOperationOptions{}
}

func (o GetCertificatesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCertificatesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetCertificatesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludePending != nil {
		out.Append("includePending", fmt.Sprintf("%v", *o.IncludePending))
	}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetCertificatesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetCertificatesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetCertificates ...
func (c CertificatesClient) GetCertificates(ctx context.Context, options GetCertificatesOperationOptions) (result GetCertificatesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetCertificatesCustomPager{},
		Path:          "/certificates",
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
		Values *[]CertificateItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetCertificatesComplete retrieves all the results into a single object
func (c CertificatesClient) GetCertificatesComplete(ctx context.Context, options GetCertificatesOperationOptions) (GetCertificatesCompleteResult, error) {
	return c.GetCertificatesCompleteMatchingPredicate(ctx, options, CertificateItemOperationPredicate{})
}

// GetCertificatesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificatesClient) GetCertificatesCompleteMatchingPredicate(ctx context.Context, options GetCertificatesOperationOptions, predicate CertificateItemOperationPredicate) (result GetCertificatesCompleteResult, err error) {
	items := make([]CertificateItem, 0)

	resp, err := c.GetCertificates(ctx, options)
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

	result = GetCertificatesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
