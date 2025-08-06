package containers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/hashicorp/go-azure-sdk/sdk/odata"
)

type ListBlobsInput struct {
	Delimiter  *string
	Include    *[]Dataset
	Marker     *string
	MaxResults *int
	Prefix     *string
}

type ListBlobsResponse struct {
	ListBlobsResult

	HttpResponse *http.Response
}

type ListBlobsResult struct {
	Delimiter  string  `xml:"Delimiter"`
	Marker     string  `xml:"Marker"`
	MaxResults int     `xml:"MaxResults"`
	NextMarker *string `xml:"NextMarker,omitempty"`
	Prefix     string  `xml:"Prefix"`
	Blobs      Blobs   `xml:"Blobs"`
}

type Blobs struct {
	Blobs      []BlobDetails `xml:"Blob"`
	BlobPrefix *BlobPrefix   `xml:"BlobPrefix"`
}

type BlobDetails struct {
	Name       string                 `xml:"Name"`
	Deleted    bool                   `xml:"Deleted,omitempty"`
	MetaData   map[string]interface{} `map:"Metadata,omitempty"`
	Properties *BlobProperties        `xml:"Properties,omitempty"`
	Snapshot   *string                `xml:"Snapshot,omitempty"`
}

type BlobProperties struct {
	AccessTier             *string `xml:"AccessTier,omitempty"`
	AccessTierInferred     *bool   `xml:"AccessTierInferred,omitempty"`
	AccessTierChangeTime   *string `xml:"AccessTierChangeTime,omitempty"`
	BlobType               *string `xml:"BlobType,omitempty"`
	BlobSequenceNumber     *string `xml:"x-ms-blob-sequence-number,omitempty"`
	CacheControl           *string `xml:"Cache-Control,omitempty"`
	ContentEncoding        *string `xml:"ContentEncoding,omitempty"`
	ContentLanguage        *string `xml:"Content-Language,omitempty"`
	ContentLength          *int64  `xml:"Content-Length,omitempty"`
	ContentMD5             *string `xml:"Content-MD5,omitempty"`
	ContentType            *string `xml:"Content-Type,omitempty"`
	CopyCompletionTime     *string `xml:"CopyCompletionTime,omitempty"`
	CopyId                 *string `xml:"CopyId,omitempty"`
	CopyStatus             *string `xml:"CopyStatus,omitempty"`
	CopySource             *string `xml:"CopySource,omitempty"`
	CopyProgress           *string `xml:"CopyProgress,omitempty"`
	CopyStatusDescription  *string `xml:"CopyStatusDescription,omitempty"`
	CreationTime           *string `xml:"CreationTime,omitempty"`
	ETag                   *string `xml:"Etag,omitempty"`
	DeletedTime            *string `xml:"DeletedTime,omitempty"`
	IncrementalCopy        *bool   `xml:"IncrementalCopy,omitempty"`
	LastModified           *string `xml:"Last-Modified,omitempty"`
	LeaseDuration          *string `xml:"LeaseDuration,omitempty"`
	LeaseState             *string `xml:"LeaseState,omitempty"`
	LeaseStatus            *string `xml:"LeaseStatus,omitempty"`
	RemainingRetentionDays *string `xml:"RemainingRetentionDays,omitempty"`
	ServerEncrypted        *bool   `xml:"ServerEncrypted,omitempty"`
}

type BlobPrefix struct {
	Name string `xml:"Name"`
}

// ListBlobs lists the blobs matching the specified query within the specified Container
func (c Client) ListBlobs(ctx context.Context, containerName string, input ListBlobsInput) (result ListBlobsResponse, err error) {
	if containerName == "" {
		err = fmt.Errorf("`containerName` cannot be an empty string")
		return
	}
	if input.MaxResults != nil && (*input.MaxResults <= 0 || *input.MaxResults > 5000) {
		err = fmt.Errorf("`input.MaxResults` can either be nil or between 0 and 5000")
		return
	}
	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod: http.MethodGet,
		OptionsObject: listBlobsOptions{
			delimiter:  input.Delimiter,
			include:    input.Include,
			marker:     input.Marker,
			maxResults: input.MaxResults,
			prefix:     input.Prefix,
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

		if err == nil {
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

var _ client.Options = listBlobsOptions{}

type listBlobsOptions struct {
	delimiter  *string
	include    *[]Dataset
	marker     *string
	maxResults *int
	prefix     *string
}

func (o listBlobsOptions) ToHeaders() *client.Headers {
	return nil
}

func (o listBlobsOptions) ToOData() *odata.Query {
	return nil
}

func (o listBlobsOptions) ToQuery() *client.QueryParams {
	query := containerOptions{}.ToQuery()
	query.Append("comp", "list")

	if o.delimiter != nil {
		query.Append("delimiter", *o.delimiter)
	}
	if o.include != nil {
		vals := make([]string, 0)
		for _, v := range *o.include {
			vals = append(vals, string(v))
		}
		include := strings.Join(vals, ",")
		query.Append("include", include)
	}
	if o.marker != nil {
		query.Append("marker", *o.marker)
	}
	if o.maxResults != nil {
		query.Append("maxresults", fmt.Sprintf("%d", *o.maxResults))
	}
	if o.prefix != nil {
		query.Append("prefix", *o.prefix)
	}
	return query
}
