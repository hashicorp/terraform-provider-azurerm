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

type GetDeletedCertificatesOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]DeletedCertificateItem
}

type GetDeletedCertificatesCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []DeletedCertificateItem
}

type GetDeletedCertificatesOperationOptions struct {
	IncludePending *bool
	Maxresults     *int64
}

func DefaultGetDeletedCertificatesOperationOptions() GetDeletedCertificatesOperationOptions {
	return GetDeletedCertificatesOperationOptions{}
}

func (o GetDeletedCertificatesOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o GetDeletedCertificatesOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o GetDeletedCertificatesOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.IncludePending != nil {
		out.Append("includePending", fmt.Sprintf("%v", *o.IncludePending))
	}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	return &out
}

type GetDeletedCertificatesCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *GetDeletedCertificatesCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// GetDeletedCertificates ...
func (c CertificatesClient) GetDeletedCertificates(ctx context.Context, options GetDeletedCertificatesOperationOptions) (result GetDeletedCertificatesOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &GetDeletedCertificatesCustomPager{},
		Path:          "/deletedcertificates",
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
		Values *[]DeletedCertificateItem `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// GetDeletedCertificatesComplete retrieves all the results into a single object
func (c CertificatesClient) GetDeletedCertificatesComplete(ctx context.Context, options GetDeletedCertificatesOperationOptions) (GetDeletedCertificatesCompleteResult, error) {
	return c.GetDeletedCertificatesCompleteMatchingPredicate(ctx, options, DeletedCertificateItemOperationPredicate{})
}

// GetDeletedCertificatesCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificatesClient) GetDeletedCertificatesCompleteMatchingPredicate(ctx context.Context, options GetDeletedCertificatesOperationOptions, predicate DeletedCertificateItemOperationPredicate) (result GetDeletedCertificatesCompleteResult, err error) {
	items := make([]DeletedCertificateItem, 0)

	resp, err := c.GetDeletedCertificates(ctx, options)
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

	result = GetDeletedCertificatesCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
