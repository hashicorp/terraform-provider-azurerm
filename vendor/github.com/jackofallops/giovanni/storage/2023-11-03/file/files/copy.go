package files

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
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
	HttpResponse *http.Response

	// The CopyID, which can be passed to AbortCopy to abort the copy.
	CopyID string

	// Either `success` or `pending`
	CopySuccess string
}

// Copy copies a blob or file to a destination file within the storage account asynchronously.
func (c Client) Copy(ctx context.Context, shareName, path, fileName string, input CopyInput) (result CopyResponse, err error) {

	if shareName == "" {
		err = fmt.Errorf("`shareName` cannot be an empty string")
		return
	}

	if strings.ToLower(shareName) != shareName {
		err = fmt.Errorf("`shareName` must be a lower-cased string")
		return
	}

	if fileName == "" {
		err = fmt.Errorf("`fileName` cannot be an empty string")
		return
	}

	if input.CopySource == "" {
		err = fmt.Errorf("`input.CopySource` cannot be an empty string")
		return
	}

	if err = metadata.Validate(input.MetaData); err != nil {
		err = fmt.Errorf("`input.MetaData` is not valid: %s", err)
		return
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

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response

		if err == nil {
			if resp.Header != nil {
				result.CopyID = resp.Header.Get("x-ms-copy-id")
				result.CopySuccess = resp.Header.Get("x-ms-copy-status")
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
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
