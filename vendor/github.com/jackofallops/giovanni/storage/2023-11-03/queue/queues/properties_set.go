package queues

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type SetStorageServicePropertiesResponse struct {
	HttpResponse *http.Response
}

type SetStorageServicePropertiesInput struct {
	Properties StorageServiceProperties
}

// SetServiceProperties sets the properties for this queue
func (c Client) SetServiceProperties(ctx context.Context, input SetStorageServicePropertiesInput) (result SetStorageServicePropertiesResponse, err error) {

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusAccepted,
		},
		HttpMethod:    http.MethodPut,
		OptionsObject: setStorageServicePropertiesOptions{},
		Path:          "/",
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return
	}

	marshalledProps, err := xml.Marshal(&input.Properties)
	if err != nil {
		return result, fmt.Errorf("marshalling request: %+v", err)
	}
	body := xml.Header + string(marshalledProps)
	req.Body = io.NopCloser(bytes.NewReader([]byte(body)))
	req.ContentLength = int64(len(body))
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))

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

type setStorageServicePropertiesOptions struct{}

func (s setStorageServicePropertiesOptions) ToHeaders() *client.Headers {
	return nil
}

func (s setStorageServicePropertiesOptions) ToOData() *odata.Query {
	return nil
}

func (s setStorageServicePropertiesOptions) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("restype", "service")
	out.Append("comp", "properties")
	return out
}
