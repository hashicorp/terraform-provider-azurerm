package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteSnapshotsInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string
}

type DeleteSnapshotsResponse struct {
	HttpResponse *client.Response
}

// DeleteSnapshots marks all Snapshots of a Blob for Deletion, which will be deleted during the next Garbage Collection Cycle.
func (c Client) DeleteSnapshots(ctx context.Context, containerName, blobName string, input DeleteSnapshotsInput) (resp DeleteSnapshotsResponse, err error) {

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
			http.StatusAccepted,
		},
		HttpMethod: http.MethodDelete,
		OptionsObject: deleteSnapshotsOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s", containerName, blobName),
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

type deleteSnapshotsOptions struct {
	input DeleteSnapshotsInput
}

func (d deleteSnapshotsOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	headers.Append("x-ms-delete-snapshots", "only")

	if d.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *d.input.LeaseID)
	}
	return headers
}

func (d deleteSnapshotsOptions) ToOData() *odata.Query {
	return nil
}

func (d deleteSnapshotsOptions) ToQuery() *client.QueryParams {
	return nil
}
