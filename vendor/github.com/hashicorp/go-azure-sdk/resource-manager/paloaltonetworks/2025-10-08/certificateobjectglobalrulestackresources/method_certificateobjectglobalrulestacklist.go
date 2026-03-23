package certificateobjectglobalrulestackresources

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificateObjectGlobalRulestackListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateObjectGlobalRulestackResource
}

type CertificateObjectGlobalRulestackListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateObjectGlobalRulestackResource
}

type CertificateObjectGlobalRulestackListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CertificateObjectGlobalRulestackListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CertificateObjectGlobalRulestackList ...
func (c CertificateObjectGlobalRulestackResourcesClient) CertificateObjectGlobalRulestackList(ctx context.Context, id GlobalRulestackId) (result CertificateObjectGlobalRulestackListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CertificateObjectGlobalRulestackListCustomPager{},
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
		Values *[]CertificateObjectGlobalRulestackResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CertificateObjectGlobalRulestackListComplete retrieves all the results into a single object
func (c CertificateObjectGlobalRulestackResourcesClient) CertificateObjectGlobalRulestackListComplete(ctx context.Context, id GlobalRulestackId) (CertificateObjectGlobalRulestackListCompleteResult, error) {
	return c.CertificateObjectGlobalRulestackListCompleteMatchingPredicate(ctx, id, CertificateObjectGlobalRulestackResourceOperationPredicate{})
}

// CertificateObjectGlobalRulestackListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificateObjectGlobalRulestackResourcesClient) CertificateObjectGlobalRulestackListCompleteMatchingPredicate(ctx context.Context, id GlobalRulestackId, predicate CertificateObjectGlobalRulestackResourceOperationPredicate) (result CertificateObjectGlobalRulestackListCompleteResult, err error) {
	items := make([]CertificateObjectGlobalRulestackResource, 0)

	resp, err := c.CertificateObjectGlobalRulestackList(ctx, id)
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

	result = CertificateObjectGlobalRulestackListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
