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

type CertificateListByBatchAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]Certificate
}

type CertificateListByBatchAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []Certificate
}

type CertificateListByBatchAccountOperationOptions struct {
	Filter     *string
	Maxresults *int64
	Select     *string
}

func DefaultCertificateListByBatchAccountOperationOptions() CertificateListByBatchAccountOperationOptions {
	return CertificateListByBatchAccountOperationOptions{}
}

func (o CertificateListByBatchAccountOperationOptions) ToHeaders() *client.Headers {
	out := client.Headers{}

	return &out
}

func (o CertificateListByBatchAccountOperationOptions) ToOData() *odata.Query {
	out := odata.Query{}

	return &out
}

func (o CertificateListByBatchAccountOperationOptions) ToQuery() *client.QueryParams {
	out := client.QueryParams{}
	if o.Filter != nil {
		out.Append("$filter", fmt.Sprintf("%v", *o.Filter))
	}
	if o.Maxresults != nil {
		out.Append("maxresults", fmt.Sprintf("%v", *o.Maxresults))
	}
	if o.Select != nil {
		out.Append("$select", fmt.Sprintf("%v", *o.Select))
	}
	return &out
}

type CertificateListByBatchAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CertificateListByBatchAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CertificateListByBatchAccount ...
func (c CertificatesClient) CertificateListByBatchAccount(ctx context.Context, id BatchAccountId, options CertificateListByBatchAccountOperationOptions) (result CertificateListByBatchAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodGet,
		OptionsObject: options,
		Pager:         &CertificateListByBatchAccountCustomPager{},
		Path:          fmt.Sprintf("%s/certificates", id.ID()),
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
		Values *[]Certificate `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CertificateListByBatchAccountComplete retrieves all the results into a single object
func (c CertificatesClient) CertificateListByBatchAccountComplete(ctx context.Context, id BatchAccountId, options CertificateListByBatchAccountOperationOptions) (CertificateListByBatchAccountCompleteResult, error) {
	return c.CertificateListByBatchAccountCompleteMatchingPredicate(ctx, id, options, CertificateOperationPredicate{})
}

// CertificateListByBatchAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificatesClient) CertificateListByBatchAccountCompleteMatchingPredicate(ctx context.Context, id BatchAccountId, options CertificateListByBatchAccountOperationOptions, predicate CertificateOperationPredicate) (result CertificateListByBatchAccountCompleteResult, err error) {
	items := make([]Certificate, 0)

	resp, err := c.CertificateListByBatchAccount(ctx, id, options)
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

	result = CertificateListByBatchAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
