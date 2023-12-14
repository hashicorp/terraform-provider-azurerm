package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type UndeleteResponse struct {
	HttpResponse *client.Response
}

// Undelete restores the contents and metadata of soft deleted blob and any associated soft deleted snapshots.
func (c Client) Undelete(ctx context.Context, containerName, blobName string) (resp UndeleteResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: undeleteOptions{},
		Path:          fmt.Sprintf("/%s/%s", containerName, blobName),
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

type undeleteOptions struct{}

func (u undeleteOptions) ToHeaders() *client.Headers {
	return nil
}

func (u undeleteOptions) ToOData() *odata.Query {
	return nil
}

func (u undeleteOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "undelete")
	return out
}
