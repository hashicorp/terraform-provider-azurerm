package webpubsub

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
	Items []CustomCertificate
}

// CustomCertificatesList ...
func (c WebPubSubClient) CustomCertificatesList(ctx context.Context, id WebPubSubId) (result CustomCertificatesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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
func (c WebPubSubClient) CustomCertificatesListComplete(ctx context.Context, id WebPubSubId) (CustomCertificatesListCompleteResult, error) {
	return c.CustomCertificatesListCompleteMatchingPredicate(ctx, id, CustomCertificateOperationPredicate{})
}

// CustomCertificatesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c WebPubSubClient) CustomCertificatesListCompleteMatchingPredicate(ctx context.Context, id WebPubSubId, predicate CustomCertificateOperationPredicate) (result CustomCertificatesListCompleteResult, err error) {
	items := make([]CustomCertificate, 0)

	resp, err := c.CustomCertificatesList(ctx, id)
	if err != nil {
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
		Items: items,
	}
	return
}
