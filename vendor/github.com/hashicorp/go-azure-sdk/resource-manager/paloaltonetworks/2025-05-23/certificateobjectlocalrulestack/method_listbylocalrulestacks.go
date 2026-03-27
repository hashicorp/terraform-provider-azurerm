package certificateobjectlocalrulestack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByLocalRulestacksOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateObjectLocalRulestackResource
}

type ListByLocalRulestacksCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateObjectLocalRulestackResource
}

type ListByLocalRulestacksCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByLocalRulestacksCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByLocalRulestacks ...
func (c CertificateObjectLocalRulestackClient) ListByLocalRulestacks(ctx context.Context, id LocalRulestackId) (result ListByLocalRulestacksOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByLocalRulestacksCustomPager{},
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

// ListByLocalRulestacksComplete retrieves all the results into a single object
func (c CertificateObjectLocalRulestackClient) ListByLocalRulestacksComplete(ctx context.Context, id LocalRulestackId) (ListByLocalRulestacksCompleteResult, error) {
	return c.ListByLocalRulestacksCompleteMatchingPredicate(ctx, id, CertificateObjectLocalRulestackResourceOperationPredicate{})
}

// ListByLocalRulestacksCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificateObjectLocalRulestackClient) ListByLocalRulestacksCompleteMatchingPredicate(ctx context.Context, id LocalRulestackId, predicate CertificateObjectLocalRulestackResourceOperationPredicate) (result ListByLocalRulestacksCompleteResult, err error) {
	items := make([]CertificateObjectLocalRulestackResource, 0)

	resp, err := c.ListByLocalRulestacks(ctx, id)
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

	result = ListByLocalRulestacksCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
