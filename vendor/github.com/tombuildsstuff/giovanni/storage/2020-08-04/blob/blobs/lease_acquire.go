package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type AcquireLeaseInput struct {
	// The ID of the existing Lease, if leased
	LeaseID *string

	// Specifies the duration of the lease, in seconds, or negative one (-1) for a lease that never expires.
	// A non-infinite lease can be between 15 and 60 seconds
	LeaseDuration int

	// The Proposed new ID for the Lease
	ProposedLeaseID *string
}

type AcquireLeaseResponse struct {
	HttpResponse *client.Response

	LeaseID string
}

// AcquireLease establishes and manages a lock on a blob for write and delete operations.
func (c Client) AcquireLease(ctx context.Context, containerName, blobName string, input AcquireLeaseInput) (resp AcquireLeaseResponse, err error) {

	if containerName == "" {
		return resp, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return resp, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return resp, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if input.LeaseID != nil && *input.LeaseID == "" {
		return resp, fmt.Errorf("`input.LeaseID` cannot be an empty string, if specified")
	}

	if input.ProposedLeaseID != nil && *input.ProposedLeaseID == "" {
		return resp, fmt.Errorf("`input.ProposedLeaseID` cannot be an empty string, if specified")
	}
	// An infinite lease duration is -1 seconds. A non-infinite lease can be between 15 and 60 seconds
	if input.LeaseDuration != -1 && (input.LeaseDuration <= 15 || input.LeaseDuration >= 60) {
		return resp, fmt.Errorf("`input.LeaseDuration` must be -1 (infinite), or between 15 and 60 seconds")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: acquireLeaseOptions{
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

	if resp.HttpResponse != nil {
		if resp.HttpResponse.Header != nil {
			resp.LeaseID = resp.HttpResponse.Header.Get("x-ms-lease-id")
		}
	}

	return
}

type acquireLeaseOptions struct {
	input AcquireLeaseInput
}

func (a acquireLeaseOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	headers.Append("x-ms-lease-action", "acquire")
	headers.Append("x-ms-lease-duration", strconv.Itoa(a.input.LeaseDuration))

	if a.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *a.input.LeaseID)
	}

	if a.input.ProposedLeaseID != nil {
		headers.Append("x-ms-proposed-lease-id", *a.input.ProposedLeaseID)
	}

	return headers
}

func (a acquireLeaseOptions) ToOData() *odata.Query {
	return nil
}

func (a acquireLeaseOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "lease")
	return out
}
