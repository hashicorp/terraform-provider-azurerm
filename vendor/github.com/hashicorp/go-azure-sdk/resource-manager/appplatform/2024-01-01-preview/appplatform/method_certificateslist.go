package appplatform

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CertificateResource
}

type CertificatesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CertificateResource
}

type CertificatesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CertificatesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CertificatesList ...
func (c AppPlatformClient) CertificatesList(ctx context.Context, id commonids.SpringCloudServiceId) (result CertificatesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CertificatesListCustomPager{},
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
		Values *[]CertificateResource `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CertificatesListComplete retrieves all the results into a single object
func (c AppPlatformClient) CertificatesListComplete(ctx context.Context, id commonids.SpringCloudServiceId) (CertificatesListCompleteResult, error) {
	return c.CertificatesListCompleteMatchingPredicate(ctx, id, CertificateResourceOperationPredicate{})
}

// CertificatesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c AppPlatformClient) CertificatesListCompleteMatchingPredicate(ctx context.Context, id commonids.SpringCloudServiceId, predicate CertificateResourceOperationPredicate) (result CertificatesListCompleteResult, err error) {
	items := make([]CertificateResource, 0)

	resp, err := c.CertificatesList(ctx, id)
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

	result = CertificatesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
