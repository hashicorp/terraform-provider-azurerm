package assessmentsmetadata

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GetInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *SecurityAssessmentMetadataResponse
}

// GetInSubscription ...
func (c AssessmentsMetadataClient) GetInSubscription(ctx context.Context, id ProviderAssessmentMetadataId) (result GetInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForGetInSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "GetInSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "GetInSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGetInSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "GetInSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGetInSubscription prepares the GetInSubscription request.
func (c AssessmentsMetadataClient) preparerForGetInSubscription(ctx context.Context, id ProviderAssessmentMetadataId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsGet(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGetInSubscription handles the response to the GetInSubscription request. The method always
// closes the http.Response Body.
func (c AssessmentsMetadataClient) responderForGetInSubscription(resp *http.Response) (result GetInSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
