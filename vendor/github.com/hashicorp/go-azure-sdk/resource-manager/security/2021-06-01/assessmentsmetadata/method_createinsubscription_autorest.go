package assessmentsmetadata

import (
	"context"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type CreateInSubscriptionOperationResponse struct {
	HttpResponse *http.Response
	Model        *SecurityAssessmentMetadataResponse
}

// CreateInSubscription ...
func (c AssessmentsMetadataClient) CreateInSubscription(ctx context.Context, id ProviderAssessmentMetadataId, input SecurityAssessmentMetadataResponse) (result CreateInSubscriptionOperationResponse, err error) {
	req, err := c.preparerForCreateInSubscription(ctx, id, input)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "CreateInSubscription", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "CreateInSubscription", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForCreateInSubscription(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "assessmentsmetadata.AssessmentsMetadataClient", "CreateInSubscription", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForCreateInSubscription prepares the CreateInSubscription request.
func (c AssessmentsMetadataClient) preparerForCreateInSubscription(ctx context.Context, id ProviderAssessmentMetadataId, input SecurityAssessmentMetadataResponse) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPut(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(id.ID()),
		autorest.WithJSON(input),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForCreateInSubscription handles the response to the CreateInSubscription request. The method always
// closes the http.Response Body.
func (c AssessmentsMetadataClient) responderForCreateInSubscription(resp *http.Response) (result CreateInSubscriptionOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp

	return
}
