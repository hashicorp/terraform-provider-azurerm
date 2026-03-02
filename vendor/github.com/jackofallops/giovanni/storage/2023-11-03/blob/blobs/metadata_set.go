package blobs

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type SetMetaDataInput struct {
	// The ID of the Lease
	// This must be specified if a Lease is present on the Blob, else a 403 is returned
	LeaseID *string

	// Any metadata which should be added to this blob
	MetaData map[string]string

	// The encryption scope for the blob.
	EncryptionScope *string
}

type SetMetaDataResponse struct {
	HttpResponse *http.Response
}

// SetMetaData marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
func (c Client) SetMetaData(ctx context.Context, containerName, blobName string, input SetMetaDataInput) (result SetMetaDataResponse, err error) {
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

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf(fmt.Sprintf("`input.MetaData` is not valid: %s.", err))
		return
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: setMetadataOptions{
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

type setMetadataOptions struct {
	input SetMetaDataInput
}

func (s setMetadataOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if s.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *s.input.LeaseID)
	}
	if s.input.EncryptionScope != nil {
		headers.Append("x-ms-encryption-scope", *s.input.EncryptionScope)
	}
	headers.Merge(metadata.SetMetaDataHeaders(s.input.MetaData))
	return headers
}

func (s setMetadataOptions) ToOData() *odata.Query {
	return nil
}

func (s setMetadataOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "metadata")
	return out
}
