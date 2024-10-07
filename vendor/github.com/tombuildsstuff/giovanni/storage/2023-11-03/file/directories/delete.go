package directories

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
)

type DeleteResponse struct {
	HttpResponse *http.Response
}

// Delete removes the specified empty directory
// Note that the directory must be empty before it can be deleted.
func (c Client) Delete(ctx context.Context, shareName, path string) (result DeleteResponse, err error) {

	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if path == "" {
		err = fmt.Errorf("`path` cannot be an empty string")
		return
	}

	// Retry the directory deletion if the directory is not empty (deleted files take a little while to disappear)
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
					return strings.Contains(*res.Code, "DirectoryNotEmpty"), nil
				}
			}
		}
		return false, nil
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: directoriesOptions{},
		Path:          fmt.Sprintf("/%s/%s", shareName, path),
		RetryFunc:     retryFunc,
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
