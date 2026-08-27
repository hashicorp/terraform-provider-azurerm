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

type GetCertificateVersionsOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateItem
}

type GetCertificateVersionsCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateItem
}

type GetCertificateVersionsOperationOptions struct {
	Maxresults *int64
}

func DefaultGetCertificateVersionsOperationOptions() GetCertificateVersionsOperationOptions {
	return GetCertificateVersionsOperationOptions{}
}

func (o GetCertificateVersionsOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetCertificateVersionsOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetCertificateVersionsOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetCertificateVersionsCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetCertificateVersionsCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetCertificateVersions ...
func (c CertificatesClient) GetCertificateVersions(ctx context.Context, id CertificateId, options GetCertificateVersionsOperationOptions) (result GetCertificateVersionsOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetCertificateVersionsCustomPager{},
		Path:          fmt.Sprintf("%s/versions", id.Path()),
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

// GetCertificateVersionsComplete retrieves all the results into a single object
func (c CertificatesClient) GetCertificateVersionsComplete(ctx context.Context, id CertificateId, options GetCertificateVersionsOperationOptions) (GetCertificateVersionsCompleteResult, error) {
	return c.GetCertificateVersionsCompleteMatchingPredicate(ctx, id, options, CertificateItemOperationPredicate{})
}

// GetCertificateVersionsCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificatesClient) GetCertificateVersionsCompleteMatchingPredicate(ctx context.Context, id CertificateId, options GetCertificateVersionsOperationOptions, predicate CertificateItemOperationPredicate) (result GetCertificateVersionsCompleteResult, err error) {
	items := make([]CertificateItem, 0)

	resp, err := c.GetCertificateVersions(ctx, id, options)
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

	result = GetCertificateVersionsCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
