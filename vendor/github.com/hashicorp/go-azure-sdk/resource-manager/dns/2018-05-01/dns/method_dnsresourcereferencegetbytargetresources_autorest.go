package dns

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DnsResourceReferenceGetByTargetResourcesOperationResponse struct {
	HttpResponse *http.Response
	Model        *DnsResourceReferenceResult
}

// DnsResourceReferenceGetByTargetResources ...
func (c DNSClient) DnsResourceReferenceGetByTargetResources(ctx context.Context, id commonids.SubscriptionId, input DnsResourceReferenceRequest) (result DnsResourceReferenceGetByTargetResourcesOperationResponse, err error) {
	req, err := c.preparerForDnsResourceReferenceGetByTargetResources(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dns.DNSClient", "DnsResourceReferenceGetByTargetResources", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "dns.DNSClient", "DnsResourceReferenceGetByTargetResources", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDnsResourceReferenceGetByTargetResources(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "dns.DNSClient", "DnsResourceReferenceGetByTargetResources", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDnsResourceReferenceGetByTargetResources prepares the DnsResourceReferenceGetByTargetResources request.
func (c DNSClient) preparerForDnsResourceReferenceGetByTargetResources(ctx context.Context, id commonids.SubscriptionId, input DnsResourceReferenceRequest) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/providers/Microsoft.Network/getDnsResourceReference", id.ID())),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDnsResourceReferenceGetByTargetResources handles the response to the DnsResourceReferenceGetByTargetResources request. The method always
// closes the http.Response Body.
func (c DNSClient) responderForDnsResourceReferenceGetByTargetResources(resp *http.Response) (result DnsResourceReferenceGetByTargetResourcesOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
