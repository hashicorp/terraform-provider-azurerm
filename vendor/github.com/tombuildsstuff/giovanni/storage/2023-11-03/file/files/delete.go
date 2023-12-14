package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type DeleteResponse struct {
	HttpResponse *client.Response
}

// Delete immediately deletes the file from the File Share.
func (c Client) Delete(ctx context.Context, shareName, path, fileName string) (resp DeleteResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if path != "" {
		path = fmt.Sprintf("/%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: nil,
		Path:          fmt.Sprintf("%s/%s%s", shareName, path, fileName),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	resp.HttpResponse, err = req.Execute(ctx)
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}
