package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type DeleteResponse struct {
	HttpResponse *http.Response
}

// Delete immediately deletes the file from the File Share.
func (c Client) Delete(ctx context.Context, shareName, path, fileName string) (result DeleteResponse, err error) {

	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if fileName == "" {
		err = fmt.Errorf("`fileName` cannot be an empty string")
		return
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
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
