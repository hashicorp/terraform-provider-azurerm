package webapps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ListPublicCertificatesSlotOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]PublicCertificate
}

type ListPublicCertificatesSlotCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []PublicCertificate
}

type ListPublicCertificatesSlotCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *ListPublicCertificatesSlotCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// ListPublicCertificatesSlot ...
func (c WebAppsClient) ListPublicCertificatesSlot(ctx context.Context, id SlotId) (result ListPublicCertificatesSlotOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &ListPublicCertificatesSlotCustomPager{},
		Path:       fmt.Sprintf("%s/publicCertificates", id.ID()),
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
		Values *[]PublicCertificate `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// ListPublicCertificatesSlotComplete retrieves all the results into a single object
func (c WebAppsClient) ListPublicCertificatesSlotComplete(ctx context.Context, id SlotId) (ListPublicCertificatesSlotCompleteResult, error) {
	return c.ListPublicCertificatesSlotCompleteMatchingPredicate(ctx, id, PublicCertificateOperationPredicate{})
}

// ListPublicCertificatesSlotCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebAppsClient) ListPublicCertificatesSlotCompleteMatchingPredicate(ctx context.Context, id SlotId, predicate PublicCertificateOperationPredicate) (result ListPublicCertificatesSlotCompleteResult, err error) {
	items := make([]PublicCertificate, 0)

	resp, err := c.ListPublicCertificatesSlot(ctx, id)
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

	result = ListPublicCertificatesSlotCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
