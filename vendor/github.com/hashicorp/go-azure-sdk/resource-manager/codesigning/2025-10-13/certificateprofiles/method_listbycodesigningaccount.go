package certificateprofiles

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListByCodeSigningAccountOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateProfile
}

type ListByCodeSigningAccountCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateProfile
}

type ListByCodeSigningAccountCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListByCodeSigningAccountCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListByCodeSigningAccount ...
func (c CertificateProfilesClient) ListByCodeSigningAccount(ctx context.Context, id CodeSigningAccountId) (result ListByCodeSigningAccountOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListByCodeSigningAccountCustomPager{},
		Path:       fmt.Sprintf("%s/certificateProfiles", id.ID()),
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
		Values *[]CertificateProfile `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListByCodeSigningAccountComplete retrieves all the results into a single object
func (c CertificateProfilesClient) ListByCodeSigningAccountComplete(ctx context.Context, id CodeSigningAccountId) (ListByCodeSigningAccountCompleteResult, error) {
	return c.ListByCodeSigningAccountCompleteMatchingPredicate(ctx, id, CertificateProfileOperationPredicate{})
}

// ListByCodeSigningAccountCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c CertificateProfilesClient) ListByCodeSigningAccountCompleteMatchingPredicate(ctx context.Context, id CodeSigningAccountId, predicate CertificateProfileOperationPredicate) (result ListByCodeSigningAccountCompleteResult, err error) {
	items := make([]CertificateProfile, 0)

	resp, err := c.ListByCodeSigningAccount(ctx, id)
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

	result = ListByCodeSigningAccountCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
