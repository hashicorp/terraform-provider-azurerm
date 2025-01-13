package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type DeleteInput struct {
	// Should any Snapshots for this Blob also be deleted?
	// If the Blob has Snapshots and this is set to False a 409 Conflict will be returned
	DeleteSnapshots bool

	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string
}

type DeleteResponse struct {
	HttpResponse *http.Response
}

// Delete marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
func (c Client) Delete(ctx context.Context, containerName, blobName string, input DeleteInput) (result DeleteResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return result, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return result, fmt.Errorf("`blobName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodDelete,
		OptionsObject: deleteOptions{
			input: input,
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

type deleteOptions struct {
	input DeleteInput
}

func (d deleteOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if d.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *d.input.LeaseID)
	}

	if d.input.DeleteSnapshots {
		headers.Append("x-ms-delete-snapshots", "include")
	}

	return headers
}

func (d deleteOptions) ToOData() *odata.Query {
	return nil
}

func (d deleteOptions) ToQuery() *client.QueryParams {
	return nil
}
