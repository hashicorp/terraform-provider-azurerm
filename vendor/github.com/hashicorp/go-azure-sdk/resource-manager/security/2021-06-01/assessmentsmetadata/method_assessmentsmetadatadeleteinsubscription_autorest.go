package assessmentsmetadata

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type AssessmentsMetadataDeleteInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
}

// AssessmentsMetadataDeleteInSubscription ...
func (c AssessmentsMetadataClient) AssessmentsMetadataDeleteInSubscription(ctx context.Context, id ProviderAssessmentMetadataId) (result AssessmentsMetadataDeleteInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForAssessmentsMetadataDeleteInSubscription(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataDeleteInSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataDeleteInSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForAssessmentsMetadataDeleteInSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "AssessmentsMetadataDeleteInSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForAssessmentsMetadataDeleteInSubscription prepares the AssessmentsMetadataDeleteInSubscription request.
func (c AssessmentsMetadataClient) preparerForAssessmentsMetadataDeleteInSubscription(ctx context.Context, id ProviderAssessmentMetadataId) (*http.Request, error) {
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

// responderForAssessmentsMetadataDeleteInSubscription handles the response to the AssessmentsMetadataDeleteInSubscription request. The method always
// closes the http.Response Body.
func (c AssessmentsMetadataClient) responderForAssessmentsMetadataDeleteInSubscription(resp *http.Response) (result AssessmentsMetadataDeleteInSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
