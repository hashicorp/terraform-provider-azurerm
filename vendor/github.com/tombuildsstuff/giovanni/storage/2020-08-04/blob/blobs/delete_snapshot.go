package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteSnapshotInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string

	// The DateTime of the Snapshot which should be marked for Deletion
	SnapshotDateTime string
}

type DeleteSnapshotResponse struct {
	HttpResponse *client.Response
}

// DeleteSnapshot marks a single Snapshot of a Blob for Deletion based on it's DateTime, which will be deleted during the next Garbage Collection cycle.
func (c Client) DeleteSnapshot(ctx context.Context, containerName, blobName string, input DeleteSnapshotInput) (resp DeleteSnapshotResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.SnapshotDateTime == "" {
		return resp, fmt.Errorf("`input.SnapshotDateTime` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodDelete,
		OptionsObject: deleteSnapshotOptions{
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

type deleteSnapshotOptions struct {
	input DeleteSnapshotInput
}

func (d deleteSnapshotOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if d.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *d.input.LeaseID)
	}

	return headers
}

func (d deleteSnapshotOptions) ToOData() *odata.Query {
	return nil
}

func (d deleteSnapshotOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("snapshot", d.input.SnapshotDateTime)
	return out
}
