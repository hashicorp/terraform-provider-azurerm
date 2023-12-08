package nginxcertificate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CertificatesListOperationResponse struct {
	HttpResponse *http.Response
	OData        *odata.OData
	Model        *[]NginxCertificate
}

type CertificatesListCompleteResult struct {
	Items []NginxCertificate
}

// CertificatesList ...
func (c NginxCertificateClient) CertificatesList(ctx context.Context, id NginxDeploymentId) (result CertificatesListOperationResponse, err error) {
	opts := client.RequestOptions{
		ContentType: "application/json; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
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
		Values *[]NginxCertificate `json:"value"`
	}
	if err = resp.Unmarshal(&values); err != nil {
		return
	}

	result.Model = values.Values

	return
}

// CertificatesListComplete retrieves all the results into a single object
func (c NginxCertificateClient) CertificatesListComplete(ctx context.Context, id NginxDeploymentId) (CertificatesListCompleteResult, error) {
	return c.CertificatesListCompleteMatchingPredicate(ctx, id, NginxCertificateOperationPredicate{})
}

// CertificatesListCompleteMatchingPredicate retrieves all the results and then applies the predicate
func (c NginxCertificateClient) CertificatesListCompleteMatchingPredicate(ctx context.Context, id NginxDeploymentId, predicate NginxCertificateOperationPredicate) (result CertificatesListCompleteResult, err error) {
	items := make([]NginxCertificate, 0)

	resp, err := c.CertificatesList(ctx, id)
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

	result = CertificatesListCompleteResult{
		Items: items,
	}
	return
}
