package capi

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
	"time"
)

func (c BaseClient) Preparer(ctx context.Context, function, method, path string, decorators ...autorest.PrepareDecorator) (*http.Request, error) {

	//headers required by the API
	headers := map[string]interface{}{
		"x-ms-date":    time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
		"x-ms-version": c.Version,
	}

	//standard preparers
	preparers := []autorest.PrepareDecorator{
		autorest.WithMethod(method),
		autorest.WithBaseURL(c.BaseURI),
		autorest.WithPath(path),
		autorest.WithHeaders(headers),
	}

	//add custom ones (WithJSON ect)
	preparers = append(preparers, decorators...)

	//add authorization preparer
	preparers = append(preparers, WithAuthorizationHeader(c.AccountKey.Key, c.AccountKey.Type, c.AccountKey.Version))

	req, err := autorest.CreatePreparer(preparers...).Prepare((&http.Request{}).WithContext(ctx))
	if err != nil {
		err = autorest.NewErrorWithError(err, c.ID, function, nil, "Failure preparing request")
	}

	return req, err
}

func (c BaseClient) Sender(function string, req *http.Request) (autorest.Response, error) {
	resp, err := autorest.SendWithSender(c, req, azure.DoRetryWithRegistration(c.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, c.ID, function, resp, "Failure sending request")
	}

	return autorest.Response{Response: resp}, err
}

func (c BaseClient) Responder(function string, resp *http.Response, expectedStatus int, decorators ...autorest.RespondDecorator) (autorest.Response, error) {
	//standard responders
	responders := []autorest.RespondDecorator{
		c.ByInspecting(),
		azure.WithErrorUnlessStatusCode(expectedStatus),
	}

	//add custom ones (WithJSON ect)
	responders = append(responders, decorators...)

	//add closer
	responders = append(responders, autorest.ByClosing())

	err := autorest.Respond(resp, responders...)
	if err != nil {
		err = autorest.NewErrorWithError(err, c.ID, function, resp, "Failure responding to request")
	}

	return autorest.Response{Response: resp}, err
}
