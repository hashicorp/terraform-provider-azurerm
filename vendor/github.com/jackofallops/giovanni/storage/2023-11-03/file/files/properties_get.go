package files

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/go-azure-sdk/sdk/client"
	"github.com/jackofallops/giovanni/storage/internal/metadata"
)

type GetResponse struct {
	HttpResponse *http.Response

	CacheControl          string
	ContentDisposition    string
	ContentEncoding       string
	ContentLanguage       string
	ContentLength         *int64
	ContentMD5            string
	ContentType           string
	CopyID                string
	CopyStatus            string
	CopySource            string
	CopyProgress          string
	CopyStatusDescription string
	CopyCompletionTime    string
	Encrypted             bool

	MetaData map[string]string
}

// GetProperties returns the Properties for the specified file
func (c Client) GetProperties(ctx context.Context, shareName, path, fileName string) (result GetResponse, err error) {
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

	if path != "" {
		path = fmt.Sprintf("%s/", path)
	}

	opts := client.RequestOptions{
		ContentType: "application/xml; charset=utf-8",
		ExpectedStatusCodes: []int{
			http.StatusOK,
		},
		HttpMethod:    http.MethodHead,
		OptionsObject: nil,
		Path:          fmt.Sprintf("%s/%s%s", shareName, path, fileName),
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
				result.CacheControl = resp.Header.Get("Cache-Control")
				result.ContentDisposition = resp.Header.Get("Content-Disposition")
				result.ContentEncoding = resp.Header.Get("Content-Encoding")
				result.ContentLanguage = resp.Header.Get("Content-Language")
				result.ContentMD5 = resp.Header.Get("Content-MD5")
				result.ContentType = resp.Header.Get("Content-Type")
				result.CopyCompletionTime = resp.Header.Get("x-ms-copy-completion-time")
				result.CopyID = resp.Header.Get("x-ms-copy-id")
				result.CopyProgress = resp.Header.Get("x-ms-copy-progress")
				result.CopySource = resp.Header.Get("x-ms-copy-source")
				result.CopyStatus = resp.Header.Get("x-ms-copy-status")
				result.CopyStatusDescription = resp.Header.Get("x-ms-copy-status-description")
				result.Encrypted = strings.EqualFold(resp.Header.Get("x-ms-server-encrypted"), "true")
				result.MetaData = metadata.ParseFromHeaders(resp.Header)

				contentLengthRaw := resp.Header.Get("Content-Length")
				if contentLengthRaw != "" {
					var contentLength int
					contentLength, err = strconv.Atoi(contentLengthRaw)
					if err != nil {
						err = fmt.Errorf("parsing `Content-Length` header value %q: %s", contentLengthRaw, err)
						return
					}
					contentLengthI64 := int64(contentLength)
					result.ContentLength = &contentLengthI64
				}
			}
		}
	}
	if err != nil {
		err = fmt.Errorf("executing request: %+v", err)
		return
	}

	return
}
