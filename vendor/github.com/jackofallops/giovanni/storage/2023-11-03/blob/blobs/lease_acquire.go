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
	HttpResponse *http.Response

	LeaseID string
}

// AcquireLease establishes and manages a lock on a blob for write and delete operations.
func (c Client) AcquireLease(ctx context.Context, containerName, blobName string, input AcquireLeaseInput) (result AcquireLeaseResponse, err error) {
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

	if input.LeaseID != nil && *input.LeaseID == "" {
		err = fmt.Errorf("`input.LeaseID` cannot be an empty string, if specified")
		return
	}

	if input.ProposedLeaseID != nil && *input.ProposedLeaseID == "" {
		err = fmt.Errorf("`input.ProposedLeaseID` cannot be an empty string, if specified")
		return
	}
	// An infinite lease duration is -1 seconds. A non-infinite lease can be between 15 and 60 seconds
	if input.LeaseDuration != -1 && (input.LeaseDuration <= 15 || input.LeaseDuration >= 60) {
		err = fmt.Errorf("`input.LeaseDuration` must be -1 (infinite), or between 15 and 60 seconds")
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Response != nil && resp.Header != nil {
				result.LeaseID = resp.Header.Get("x-ms-lease-id")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
