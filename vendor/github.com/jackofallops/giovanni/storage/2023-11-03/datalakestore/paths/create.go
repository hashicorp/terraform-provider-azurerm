package paths

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type PathResource string

const PathResourceFile PathResource = "file"
const PathResourceDirectory PathResource = "directory"

type CreateInput struct {
	Resource PathResource
}

type CreateResponse struct {
	HttpResponse *http.Response
}

// Create creates a Data Lake Store Gen2 Path within a Storage Account
func (c Client) Create(ctx context.Context, fileSystemName string, path string, input CreateInput) (result CreateResponse, err error) {

	if fileSystemName == "" {
		return result, fmt.Errorf("`fileSystemName` cannot be an empty string")
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusCreated,
		},
		HttpMethod: http.MethodPut,
		OptionsObject: CreateInput{
			Resource: input.Resource,
		},

		Path: fmt.Sprintf("/%s/%s", fileSystemName, path),
	}

	req, err := c.Client.NewRequest(ctx, opts)
	if err != nil {
		err = fmt.Errorf("building request: %+v", err)
		return result, err
	}

	var resp *client.Response
	resp, err = req.Execute(ctx)
	if resp != nil && resp.Response != nil {
		result.HttpResponse = resp.Response
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return result, err
	}

	return
}

func (c CreateInput) ToHeaders() *client.Headers {
	return nil
}

func (c CreateInput) ToOData() *odata.Query {
	return nil
}

func (c CreateInput) ToQuery() *client.QueryParams {
	out := &client.QueryParams{}
	out.Append("resource", string(c.Resource))
	return out
}
