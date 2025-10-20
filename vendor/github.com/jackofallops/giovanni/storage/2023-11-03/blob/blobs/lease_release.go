package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ReleaseLeaseResponse struct {
	HttpResponse *http.Response
}

type ReleaseLeaseInput struct {
	LeaseID string
}

// ReleaseLease releases a lock based on the Lease ID.
func (c Client) ReleaseLease(ctx context.Context, containerName, blobName string, input ReleaseLeaseInput) (result ReleaseLeaseResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}

	if strings.ToLower(containerName) != containerName {
		err = fmt.Errorf("`containerName` must be a lower-cased string")
		return
	}

	if blobName == "" {
		err = fmt.Errorf("`blobName` cannot be an empty string")
		return
	}

	if input.LeaseID == "" {
		err = fmt.Errorf("`input.LeaseID` cannot be an empty string")
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: releaseLeaseOptions{
			leaseID: input.LeaseID,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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

type releaseLeaseOptions struct {
	leaseID string
}

func (r releaseLeaseOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-lease-action", "release")
	headers.Append("x-ms-lease-id", r.leaseID)
	return headers
}

func (r releaseLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (r releaseLeaseOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "lease")
	return out
}
