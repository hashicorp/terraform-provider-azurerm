package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
)

type DeleteResponse struct {
	HttpResponse *client.Response
}

// Delete marks the specified container for deletion.
// The container and any blobs contained within it are later deleted during garbage collection.
func (c Client) Delete(ctx context.Context, containerName string) (resp DeleteResponse, err error) {
	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodDelete,
		OptionsObject: containerOptions{},
		Path:          fmt.Sprintf("/%s", containerName),
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
