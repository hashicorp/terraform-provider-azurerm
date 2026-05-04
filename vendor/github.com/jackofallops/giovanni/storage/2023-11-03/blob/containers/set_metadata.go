package containers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type SetMetaDataInput struct {
	MetaData map[string]string
	LeaseId  string
}

type SetMetaDataResponse struct {
	HttpResponse *http.Response
}

// SetMetaData sets the specified MetaData on the Container without a Lease ID
func (c Client) SetMetaData(ctx context.Context, containerName string, input SetMetaDataInput) (result SetMetaDataResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}
	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf("`input.MetaData` is not valid: %s", err)
		return
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: setMetaDataOptions{
			metaData: input.MetaData,
			leaseId:  input.LeaseId,
		},
		Path: fmt.Sprintf("/%s", containerName),
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

var _ client.Options = setMetaDataOptions{}

type setMetaDataOptions struct {
	metaData map[string]string
	leaseId  string
}

func (o setMetaDataOptions) ToHeaders() *client.Headers {
	headers := containerOptions{
		metaData: o.metaData,
	}.ToHeaders()

	// If specified, Get Container Properties only succeeds if the container’s lease is active and matches this ID.
	// If there is no active lease or the ID does not match, 412 (Precondition Failed) is returned.
	if o.leaseId != "" {
		headers.Append("x-ms-lease-id", o.leaseId)
	}

	return headers
}

func (o setMetaDataOptions) ToOData() *odata.Query {
	return nil
}

func (o setMetaDataOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "metadata")
	return query
}
