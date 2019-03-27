package cosmos

import (
	"context"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"net/http"
	"time"
)

func (client BaseClient) Preparer(ctx context.Context, function, method, path string, pathParameters map[string]interface{}, decorators ...autorest.PrepareDecorator) (*http.Request, error) {

	//headers required by the API
	headers := map[string]interface{}{
		"x-ms-date":    time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"),
		"x-ms-version": client.Version,
	}

	//standard preparers
	preparers := []autorest.PrepareDecorator{
		autorest.WithMethod(method),
		autorest.WithBaseURL(client.BaseURI),
		autorest.WithPathParameters(path, pathParameters),
		autorest.WithHeaders(headers),
	}

	//add custom ones (WithJSON ect)
	preparers = append(preparers, decorators...)

	//add authorization preparer
	preparers = append(preparers, WithAuthorizationHeader(client.AccountKey.Key, client.AccountKey.Type, client.AccountKey.Version))

	req, err := autorest.CreatePreparer(preparers...).Prepare((&http.Request{}).WithContext(ctx))
	if err != nil {
		err = autorest.NewErrorWithError(err, client.ID, function, nil, "Failure preparing request")
	}

	return req, err
}

func (client BaseClient) Sender(function string, req *http.Request) (autorest.Response, error) {
	resp, err := autorest.SendWithSender(client, req, azure.DoRetryWithRegistration(client.Client))
	if err != nil {
		err = autorest.NewErrorWithError(err, client.ID, function, resp, "Failure sending request")
	}

	return autorest.Response{Response: resp}, err
}

func (client BaseClient) Responder(function string, resp *http.Response, expectedStatus int, decorators ...autorest.RespondDecorator) (autorest.Response, error) {
	//standard responders
	responders := []autorest.RespondDecorator{
		client.ByInspecting(),
		azure.WithErrorUnlessStatusCode(expectedStatus),
	}

	//add custom ones (WithJSON ect)
	responders = append(responders, decorators...)

	//add closer
	responders = append(responders, autorest.ByClosing())

	err := autorest.Respond(resp, responders...)
	if err != nil {
		err = autorest.NewErrorWithError(err, client.ID, function, resp, "Failure responding to request")
	}

	return autorest.Response{Response: resp}, err
}

type APIResponse struct {
	autorest.Response `json:"-"`
	Path              string `json:"-"` //unique identifier for this object, cosmos/{accountName}/dbs/{databaseName} ect

	ID   *string `json:"id,omitempty"`
	Self *string `json:"_self,omitempty"`
	ETag *string `json:"_etag,omitempty"`
	RID  *string `json:"_rid,omitempty"`
	TS   *int    `json:"_ts,omitempty"`
}

func (e APIResponse) PopulateCommon(r autorest.Response) {
	e.Response = r

}
