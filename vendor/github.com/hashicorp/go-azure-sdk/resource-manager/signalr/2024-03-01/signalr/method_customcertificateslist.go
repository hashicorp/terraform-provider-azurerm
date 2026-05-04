package signalr

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CustomCertificatesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]CustomCertificate
}

type CustomCertificatesListCompleteResult struct {
	LatestHttpResponse *http.Response
	Items              []CustomCertificate
}

type CustomCertificatesListCustomPager struct {
	NextLink *odata.Link `json:"nextLink"`
}

func (p *CustomCertificatesListCustomPager) NextPageLink() *odata.Link {
	defer func() {
		p.NextLink = nil
	}()

	return p.NextLink
}

// CustomCertificatesList ...
func (c SignalRClient) CustomCertificatesList(ctx context.Context, id SignalRId) (result CustomCertificatesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		Pager:      &CustomCertificatesListCustomPager{},
		Path:       fmt.Sprintf("%s/customCertificates", id.ID()),
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
		Values *[]CustomCertificate `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CustomCertificatesListComplete retrieves all the results into a single object
func (c SignalRClient) CustomCertificatesListComplete(ctx context.Context, id SignalRId) (CustomCertificatesListCompleteResult, error) {
	return c.CustomCertificatesListCompleteMatchingPredicate(ctx, id, CustomCertificateOperationPredicate{})
}

// CustomCertificatesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c SignalRClient) CustomCertificatesListCompleteMatchingPredicate(ctx context.Context, id SignalRId, predicate CustomCertificateOperationPredicate) (result CustomCertificatesListCompleteResult, err error) {
	items := make([]CustomCertificate, 0)

	resp, err := c.CustomCertificatesList(ctx, id)
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

	result = CustomCertificatesListCompleteResult{
		LatestHttpResponse: resp.HttpResponse,
		Items:              items,
	}
	return
}
