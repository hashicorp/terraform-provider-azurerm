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

type GetCertificateIssuersOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateIssuerItem
}

type GetCertificateIssuersCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateIssuerItem
}

type GetCertificateIssuersOperationOptions struct {
	Maxresults *int64
}

func DefaultGetCertificateIssuersOperationOptions() GetCertificateIssuersOperationOptions {
	return GetCertificateIssuersOperationOptions{}
}

func (o GetCertificateIssuersOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCertificateIssuersOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetCertificateIssuersOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetCertificateIssuersCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetCertificateIssuersCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetCertificateIssuers ...
func (c CertificatesClient) GetCertificateIssuers(ctx context.Context, options GetCertificateIssuersOperationOptions) (result GetCertificateIssuersOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetCertificateIssuersCustomPager{},
		Path:          "/certificates/issuers",
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
		Values *[]CertificateIssuerItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetCertificateIssuersComplete retrieves all the results into a single object
func (c CertificatesClient) GetCertificateIssuersComplete(ctx context.Context, options GetCertificateIssuersOperationOptions) (GetCertificateIssuersCompleteResult, error) {
	return c.GetCertificateIssuersCompleteMatchingPredicate(ctx, options, CertificateIssuerItemOperationPredicate{})
}

// GetCertificateIssuersCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificatesClient) GetCertificateIssuersCompleteMatchingPredicate(ctx context.Context, options GetCertificateIssuersOperationOptions, predicate CertificateIssuerItemOperationPredicate) (result GetCertificateIssuersCompleteResult, err error) {
	items := make([]CertificateIssuerItem, 0)

	resp, err := c.GetCertificateIssuers(ctx, options)
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

	result = GetCertificateIssuersCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
