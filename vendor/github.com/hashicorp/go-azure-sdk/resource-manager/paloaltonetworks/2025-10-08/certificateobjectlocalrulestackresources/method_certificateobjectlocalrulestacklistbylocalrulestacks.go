package certificateobjectlocalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateObjectLocalRulestackListByLocalRulestacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateObjectLocalRulestackResource
}

type CertificateObjectLocalRulestackListByLocalRulestacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateObjectLocalRulestackResource
}

type CertificateObjectLocalRulestackListByLocalRulestacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CertificateObjectLocalRulestackListByLocalRulestacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CertificateObjectLocalRulestackListByLocalRulestacks ...
func (c CertificateObjectLocalRulestackResourcesClient) CertificateObjectLocalRulestackListByLocalRulestacks(ctx context.Context, id LocalRulestackId) (result CertificateObjectLocalRulestackListByLocalRulestacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CertificateObjectLocalRulestackListByLocalRulestacksCustomPager{},
		Path:       fmt.Sprintf("%s/certificates", id.ID()),
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
		Values *[]CertificateObjectLocalRulestackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CertificateObjectLocalRulestackListByLocalRulestacksComplete retrieves all the results into a single object
func (c CertificateObjectLocalRulestackResourcesClient) CertificateObjectLocalRulestackListByLocalRulestacksComplete(ctx context.Context, id LocalRulestackId) (CertificateObjectLocalRulestackListByLocalRulestacksCompleteResult, error) {
	return c.CertificateObjectLocalRulestackListByLocalRulestacksCompleteMatchingPredicate(ctx, id, CertificateObjectLocalRulestackResourceOperationPredicate{})
}

// CertificateObjectLocalRulestackListByLocalRulestacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificateObjectLocalRulestackResourcesClient) CertificateObjectLocalRulestackListByLocalRulestacksCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, predicate CertificateObjectLocalRulestackResourceOperationPredicate) (result CertificateObjectLocalRulestackListByLocalRulestacksCompleteResult, err error) {
	items := make([]CertificateObjectLocalRulestackResource, 0)

	resp, err := c.CertificateObjectLocalRulestackListByLocalRulestacks(ctx, id)
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

	result = CertificateObjectLocalRulestackListByLocalRulestacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
