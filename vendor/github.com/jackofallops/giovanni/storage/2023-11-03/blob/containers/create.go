package containers

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type CreateInput struct {
	// Specifies whether data in the container may be accessed publicly and the level of access
	AccessLevel AccessLevel

	// The encryption scope to set as the default on the container.
	DefaultEncryptionScope string

	// Setting this to ture indicates that every blob that's uploaded to this container uses the default encryption scope.
	EncryptionScopeOverrideDisabled bool

	// A name-value pair to associate with the container as metadata.
	MetaData map[string]string
}

type CreateResponse struct {
	HttpResponse *http.Response
}

// Create creates a new container under the specified account.
// If the container with the same name already exists, the operation fails.
func (c Client) Create(ctx context.Context, containerName string, input CreateInput) (result CreateResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}
	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf("`input.MetaData` is not valid: %+v", err)
		return
	}

	// Retry the container creation if a conflicting container is still in the process of being deleted
	retryFunc := func(resp *http.Response, _ *odata.OData) (bool, error) {
		if resp != nil {
			if response.WasStatusCode(resp, http.StatusConflict) {
				// TODO: move this error response parsing to a common helper function
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					return false, fmt.Errorf("could not parse response body")
				}
				resp.Body.Close()
				respBody = bytes.TrimPrefix(respBody, []byte("\xef\xbb\xbf"))
				res := ErrorResponse{}
				if err = xml.Unmarshal(respBody, &res); err != nil {
					return false, err
				}
				resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
				if res.Code != nil {
					return strings.Contains(*res.Code, "ContainerBeingDeleted"), nil
				}
			}
		}
		return false, nil
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: createOptions{
			input: input,
		},
		Path:      fmt.Sprintf("/%s", containerName),
		RetryFunc: retryFunc,
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

var _ client.Options = createOptions{}

type createOptions struct {
	input CreateInput
}

func (o createOptions) ToHeaders() *client.Headers {
	headers := containerOptions{
		metaData: o.input.MetaData,
	}.ToHeaders()

	// If this header is not included in the request, container data is private to the account owner.
	if o.input.AccessLevel != Private {
		headers.Append("x-ms-blob-public-access", string(o.input.AccessLevel))
	}

	if o.input.DefaultEncryptionScope != "" {
		// These two headers must be used together.
		headers.Append("x-ms-default-encryption-scope", o.input.DefaultEncryptionScope)
		headers.Append("x-ms-deny-encryption-scope-override", fmt.Sprintf("%t", o.input.EncryptionScopeOverrideDisabled))
	}

	return headers
}

func (createOptions) ToOData() *odata.Query {
	return nil
}

func (createOptions) ToQuery() *client.QueryParams {
	return containerOptions{}.ToQuery()
}
