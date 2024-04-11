package filesystems

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type DeleteResponse struct {
	HttpResponse *http.Response
}

// Delete deletes a Data Lake Store Gen2 FileSystem within a Storage Account
func (c Client) Delete(ctx context.Context, fileSystemName string) (result DeleteResponse, err error) {
	if fileSystemName == "" {
		err = fmt.Errorf("`fileSystemName` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: fileSystemOptions{},
		Path:          fmt.Sprintf("/%s", fileSystemName),
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
