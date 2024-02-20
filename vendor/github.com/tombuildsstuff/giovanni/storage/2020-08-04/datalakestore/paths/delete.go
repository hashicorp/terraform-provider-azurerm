package paths

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type DeleteResponse struct {
	HttpResponse *client.Response
}

// Delete deletes a Data Lake Store Gen2 FileSystem within a Storage Account
func (c Client) Delete(ctx context.Context, fileSystemName string, path string) (resp DeleteResponse, err error) {

	if fileSystemName == "" {
		return resp, fmt.Errorf("`fileSystemName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: nil,
		Path:          fmt.Sprintf("/%s/%s", fileSystemName, path),
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
