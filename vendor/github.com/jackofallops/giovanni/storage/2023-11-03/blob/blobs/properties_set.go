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

type SetPropertiesInput struct {
	CacheControl         *string
	ContentType          *string
	ContentMD5           *string
	ContentEncoding      *string
	ContentLanguage      *string
	LeaseID              *string
	ContentDisposition   *string
	ContentLength        *int64
	SequenceNumberAction *SequenceNumberAction
	BlobSequenceNumber   *string
}

type SetPropertiesResponse struct {
	HttpResponse *http.Response

	BlobSequenceNumber string
	Etag               string
}

// SetProperties sets system properties on the blob.
func (c Client) SetProperties(ctx context.Context, containerName, blobName string, input SetPropertiesInput) (result SetPropertiesResponse, err error) {
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
			http.StatusOK,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: setPropertiesOptions{
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

type SequenceNumberAction string

var (
	Increment SequenceNumberAction = "increment"
	Max       SequenceNumberAction = "max"
	Update    SequenceNumberAction = "update"
)

type setPropertiesOptions struct {
	input SetPropertiesInput
}

func (s setPropertiesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if s.input.CacheControl != nil {
		headers.Append("x-ms-blob-cache-control", *s.input.CacheControl)
	}
	if s.input.ContentDisposition != nil {
		headers.Append("x-ms-blob-content-disposition", *s.input.ContentDisposition)
	}
	if s.input.ContentEncoding != nil {
		headers.Append("x-ms-blob-content-encoding", *s.input.ContentEncoding)
	}
	if s.input.ContentLanguage != nil {
		headers.Append("x-ms-blob-content-language", *s.input.ContentLanguage)
	}
	if s.input.ContentMD5 != nil {
		headers.Append("x-ms-blob-content-md5", *s.input.ContentMD5)
	}
	if s.input.ContentType != nil {
		headers.Append("x-ms-blob-content-type", *s.input.ContentType)
	}
	if s.input.ContentLength != nil {
		headers.Append("x-ms-blob-content-length", strconv.Itoa(int(*s.input.ContentLength)))
	}
	if s.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *s.input.LeaseID)
	}
	if s.input.SequenceNumberAction != nil {
		headers.Append("x-ms-sequence-number-action", string(*s.input.SequenceNumberAction))
	}
	if s.input.BlobSequenceNumber != nil {
		headers.Append("x-ms-blob-sequence-number", *s.input.BlobSequenceNumber)
	}

	return headers
}

func (s setPropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (s setPropertiesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "properties")
	return out
}
