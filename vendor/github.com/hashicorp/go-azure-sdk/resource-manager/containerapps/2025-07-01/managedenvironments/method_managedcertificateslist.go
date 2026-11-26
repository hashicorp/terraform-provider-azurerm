package managedenvironments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedCertificatesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]ManagedCertificate
}

type ManagedCertificatesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []ManagedCertificate
}

type ManagedCertificatesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ManagedCertificatesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ManagedCertificatesList ...
func (c ManagedEnvironmentsClient) ManagedCertificatesList(ctx context.Context, id ManagedEnvironmentId) (result ManagedCertificatesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ManagedCertificatesListCustomPager{},
		Path:       fmt.Sprintf("%s/managedCertificates", id.ID()),
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
		Values *[]ManagedCertificate `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ManagedCertificatesListComplete retrieves all the results into a single object
func (c ManagedEnvironmentsClient) ManagedCertificatesListComplete(ctx context.Context, id ManagedEnvironmentId) (ManagedCertificatesListCompleteResult, error) {
	return c.ManagedCertificatesListCompleteMatchingPredicate(ctx, id, ManagedCertificateOperationPredicate{})
}

// ManagedCertificatesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c ManagedEnvironmentsClient) ManagedCertificatesListCompleteMatchingPredicate(ctx context.Context, id ManagedEnvironmentId, predicate ManagedCertificateOperationPredicate) (result ManagedCertificatesListCompleteResult, err error) {
	items := make([]ManagedCertificate, 0)

	resp, err := c.ManagedCertificatesList(ctx, id)
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

	result = ManagedCertificatesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
