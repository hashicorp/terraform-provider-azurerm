package authorizers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

// buildCanonicalizedHeader builds the Canonicalized Header required to sign Storage Requests
func buildCanonicalizedHeader(headers http.Header) string {
	cm := make(map[string]string)

	for k, v := range headers {
		headerName := strings.TrimSpace(strings.ToLower(k))
		if strings.HasPrefix(headerName, "x-ms-") {
			cm[headerName] = v[0]
		}
	}

	if len(cm) == 0 {
		return ""
	}

	var keys []string
	for key := range cm {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	ch := bytes.NewBufferString("")

	for _, key := range keys {
		ch.WriteString(key)
		ch.WriteRune(':')
		ch.WriteString(cm[key])
		ch.WriteRune('\n')
	}

	return strings.TrimSuffix(ch.String(), "\n")
}

// buildCanonicalizedResource builds the Canonical Resource required for to sign Storage Account requests
func buildCanonicalizedResource(uri, accountName string) (*string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	cr := bytes.NewBufferString("")
	if accountName != StorageEmulatorAccountName {
		cr.WriteString("/")
		cr.WriteString(primaryStorageAccountName(accountName))
	}

	if len(u.Path) > 0 {
		// Any portion of the CanonicalizedResource string that is derived from
		// the resource's URI should be encoded exactly as it is in the URI.
		// -- https://msdn.microsoft.com/en-gb/library/azure/dd179428.aspx
		cr.WriteString(u.EscapedPath())
	}

	// TODO: replace this with less of a hack
	if comp := u.Query().Get("comp"); comp != "" {
		cr.WriteString(fmt.Sprintf("?comp=%s", comp))
	}

	out := cr.String()
	return &out, nil
}

func formatSharedKeyLiteAuthorizationHeader(accountName, key string) string {
	canonicalizedAccountName := primaryStorageAccountName(accountName)
	return fmt.Sprintf("SharedKeyLite %s:%s", canonicalizedAccountName, key)
}

// hmacValue base-64 decodes the storageAccountKey, then signs the string with it
// as outlined here: https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-shared-key
func hmacValue(storageAccountKey, canonicalizedString string) string {
	key, err := base64.StdEncoding.DecodeString(storageAccountKey)
	if err != nil {
		return ""
	}

	encr := hmac.New(sha256.New, key)
	_, _ = encr.Write([]byte(canonicalizedString))
	return base64.StdEncoding.EncodeToString(encr.Sum(nil))
}

// prepareHeadersForRequest prepares a request so that it can be signed
// by ensuring the `date` and `x-ms-date` headers are set
func prepareHeadersForRequest(r *http.Request) {
	if r.Header == nil {
		r.Header = http.Header{}
	}

	date := time.Now().UTC().Format(http.TimeFormat)

	// a date must be set, X-Ms-Date should be used when both are set; but let's set both for completeness
	r.Header.Set("date", date)
	r.Header.Set("x-ms-date", date)
}

// primaryStorageAccountName returns the name of the primary for a given Storage Account
func primaryStorageAccountName(input string) string {
	// from https://docs.microsoft.com/en-us/rest/api/storageservices/authorize-with-shared-key
	// If you are accessing the secondary location in a storage account for which
	// read-access geo-replication (RA-GRS) is enabled, do not include the
	// -secondary designation in the authorization header.
	// For authorization purposes, the account name is always the name of the primary location,
	// even for secondary access.
	return strings.TrimSuffix(input, "-secondary")
}
