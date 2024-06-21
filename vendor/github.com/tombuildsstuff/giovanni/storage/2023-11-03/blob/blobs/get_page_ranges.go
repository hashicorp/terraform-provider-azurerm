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

type GetPageRangesInput struct {
	LeaseID *string

	StartByte *int64
	EndByte   *int64
}

type GetPageRangesResponse struct {
	HttpResponse *http.Response

	// The size of the blob in bytes
	ContentLength *int64

	// The Content Type of the blob
	ContentType string

	// The ETag associated with this blob
	ETag string

	PageRanges []PageRange `xml:"PageRange"`
}

type PageRange struct {
	// The start byte offset for this range, inclusive
	Start int64 `xml:"Start"`

	// The end byte offset for this range, inclusive
	End int64 `xml:"End"`
}

// GetPageRanges returns the list of valid page ranges for a page blob or snapshot of a page blob.
func (c Client) GetPageRanges(ctx context.Context, containerName, blobName string, input GetPageRangesInput) (result GetPageRangesResponse, err error) {
	if containerName == "" {
		return result, fmt.Errorf("`containerName` cannot be an empty string")
	}

	if strings.ToLower(containerName) != containerName {
		return result, fmt.Errorf("`containerName` must be a lower-cased string")
	}

	if blobName == "" {
		return result, fmt.Errorf("`blobName` cannot be an empty string")
	}

	if (input.StartByte != nil && input.EndByte == nil) || (input.StartByte == nil && input.EndByte != nil) {
		return result, fmt.Errorf("`input.StartByte` and `input.EndByte` must both be specified, or both be nil")
	}

	opts := client.RequestOptions{
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: getPageRangesOptions{
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
			if resp.Header != nil {
				result.ContentType = resp.Header.Get("Content-Type")
				result.ETag = resp.Header.Get("ETag")

				if v := resp.Header.Get("x-ms-blob-content-length"); v != "" {
					i, innerErr := strconv.Atoi(v)
					if innerErr != nil {
						err = fmt.Errorf("parsing `x-ms-blob-content-length` header value %q: %+v", v, innerErr)
						return
					}

					i64 := int64(i)
					result.ContentLength = &i64
				}
			}

			err = resp.Unmarshal(&result)
			if err != nil {
				err = fmt.Errorf("unmarshalling response: %+v", err)
				return
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}

type getPageRangesOptions struct {
	input GetPageRangesInput
}

func (g getPageRangesOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}

	if g.input.LeaseID != nil {
		headers.Append("x-ms-lease-id", *g.input.LeaseID)
	}

	if g.input.StartByte != nil && g.input.EndByte != nil {
		headers.Append("x-ms-range", fmt.Sprintf("bytes=%d-%d", *g.input.StartByte, *g.input.EndByte))
	}

	return headers
}

func (g getPageRangesOptions) ToOData() *odata.Query {
	return nil
}

func (g getPageRangesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("comp", "pagelist")
	return out
}
