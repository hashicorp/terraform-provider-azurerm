package sshpublickeys

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
)

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type GenerateKeyPairOperationResponse struct {
	HttpResponse *http.Response
	Model        *SshPublicKeyGenerateKeyPairResult
}

// GenerateKeyPair ...
func (c SshPublicKeysClient) GenerateKeyPair(ctx context.Context, id SshPublicKeyId) (result GenerateKeyPairOperationResponse, err error) {
	req, err := c.preparerForGenerateKeyPair(ctx, id)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sshpublickeys.SshPublicKeysClient", "GenerateKeyPair", nil, "Failure preparing request")
		return
	}

	result.HttpResponse, err = c.Client.Send(req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, "sshpublickeys.SshPublicKeysClient", "GenerateKeyPair", result.HttpResponse, "Failure sending request")
		return
	}

	result, err = c.responderForGenerateKeyPair(result.HttpResponse)
	if err != nil {
		err = autorest.NewErrorWithError(err, "sshpublickeys.SshPublicKeysClient", "GenerateKeyPair", result.HttpResponse, "Failure responding to request")
		return
	}

	return
}

// preparerForGenerateKeyPair prepares the GenerateKeyPair request.
func (c SshPublicKeysClient) preparerForGenerateKeyPair(ctx context.Context, id SshPublicKeyId) (*http.Request, error) {
	queryParameters := map[string]interface{}{
		"api-version": defaultApiVersion,
	}

	preparer := autorest.CreatePreparer(
		autorest.AsContentType("application/json; charset=utf-8"),
		autorest.AsPost(),
		autorest.WithBaseURL(c.baseUri),
		autorest.WithPath(fmt.Sprintf("%s/generateKeyPair", id.ID())),
		autorest.WithQueryParameters(queryParameters))
	return preparer.Prepare((&http.Request{}).WithContext(ctx))
}

// responderForGenerateKeyPair handles the response to the GenerateKeyPair request. The method always
// closes the http.Response Body.
func (c SshPublicKeysClient) responderForGenerateKeyPair(resp *http.Response) (result GenerateKeyPairOperationResponse, err error) {
	err = autorest.Respond(
		resp,
		azure.WithErrorUnlessStatusCode(http.StatusOK),
		autorest.ByUnmarshallingJSON(&result.Model),
		autorest.ByClosing())
	result.HttpResponse = resp
	return
}
