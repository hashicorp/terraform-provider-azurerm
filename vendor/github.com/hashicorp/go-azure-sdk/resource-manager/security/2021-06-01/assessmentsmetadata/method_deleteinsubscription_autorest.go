package assessmentsmetadata

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type DeleteInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
}

// DeleteInSubscription ...
func (c AssessmentsMetadataClient) DeleteInSubscription(ctx context.Context, id ProviderAssessmentMetadataId) (result DeleteInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForDeleteInSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "DeleteInSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "DeleteInSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForDeleteInSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "DeleteInSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForDeleteInSubscription prepares the DeleteInSubscription request.
func (c AssessmentsMetadataClient) preparerForDeleteInSubscription(ctx context.Context, id ProviderAssessmentMetadataId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsDelete(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForDeleteInSubscription handles the response to the DeleteInSubscription request. The method always
// closes the http.Response Body.
func (c AssessmentsMetadataClient) responderForDeleteInSubscription(resp *http.Response) (result DeleteInSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
