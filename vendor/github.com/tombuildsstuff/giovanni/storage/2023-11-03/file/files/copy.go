package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/tombuildsstuff/giovanni/storage/internal/metadata"
)

type CopyInput struct {
	// Specifies the URL of the source file or blob, up to 2 KB in length.
	//
	// To copy a file to another file within the same storage account, you may use Shared Key to authenticate
	// the source file. If you are copying a file from another storage account, or if you are copying a blob from
	// the same storage account or another storage account, then you must authenticate the source file or blob using a
	// shared access signature. If the source is a public blob, no authentication is required to perform the copy
	// operation. A file in a share snapshot can also be specified as a copy source.
	CopySource string

	MetaData map[string]string
}

type CopyResponse struct {
	HttpResponse *client.Response

	// The CopyID, which can be passed to AbortCopy to abort the copy.
	CopyID string

	// Either `success` or `pending`
	CopySuccess string
}

// Copy copies a blob or file to a destination file within the storage account asynchronously.
func (c Client) Copy(ctx context.Context, shareName, path, fileName string, input CopyInput) (resp CopyResponse, err error) {

	if shareName == "" {
		return resp, fmt.Errorf("`shareName` cannot be an empty string")
	}

	if strings.ToLower(shareName) != shareName {
		return resp, fmt.Errorf("`shareName` must be a lower-cased string")
	}

	if fileName == "" {
		return resp, fmt.Errorf("`fileName` cannot be an empty string")
	}

	if input.CopySource == "" {
		return resp, fmt.Errorf("`input.CopySource` cannot be an empty string")
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		return resp, fmt.Errorf("`input.MetaData` is not valid: %s", err)
	}

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: CopyOptions{
			input: input,
		},
		Path: fmt.Sprintf("/%s/%s%s", shareName, path, fileName),
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
			resp.CopyID = resp.HttpResponse.Header.Get("x-ms-copy-id")
			resp.CopySuccess = resp.HttpResponse.Header.Get("x-ms-copy-status")
		}
	}

	return
}

type CopyOptions struct {
	input CopyInput
}

func (c CopyOptions) ToHeaders() *client.Headers {
	headers := &client.Headers{}
	if len(c.input.MetaData) > 0 {
		headers.Merge(metadata.SetMetaDataHeaders(c.input.MetaData))
	}
	headers.Append("x-ms-copy-source", c.input.CopySource)
	return headers
}

func (c CopyOptions) ToOData() *odata.Query {
	return nil
}

func (c CopyOptions) ToQuery() *client.QueryParams {
	return nil
}
