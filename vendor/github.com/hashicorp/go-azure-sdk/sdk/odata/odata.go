// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package odata

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"

	"github.com/hashicorp/go-uuid"
)

// ODataVersion describes the highest OData spec version supported by this package
const ODataVersion = "4.0" // TODO: support 4.01 - https://docs.oasis-open.org/odata/odata-json-format/v4.01/cs01/odata-json-format-v4.01-cs01.html#_Toc499720587

// Id describes the ID of an OData entity.
type Id string

func (id Id) MarshalJSON() ([]byte, error) {
	id2 := regexp.MustCompile(`/v2/`).ReplaceAllString(string(id), `/v1.0/`)

	// For MS Graph, fix invalid IDs by attempting to parse a UUID from them and constructing a valid URI
	// TODO: improve logic here, currently assumes all invalid IDs are graph entities
	u, err := url.Parse(id2)
	if err != nil || u.Scheme == "" || u.Host == "" {
		matches := regexp.MustCompile(`([^()'"]+)\(['"]([^'"]+)['"]\)`).FindStringSubmatch(id2)
		if len(matches) != 3 {
			return nil, fmt.Errorf("Marshaling odata.Id: could not match a GUID")
		}

		objectType := matches[1]
		guid := matches[2]
		if _, err = uuid.ParseUUID(guid); err != nil {
			return nil, fmt.Errorf("Marshaling odata.Id: %+v", err)
		}

		// Although we're hard-coding `graph.microsoft.com` here, this doesn't _appear_ to be a problem
		// The host can seemingly be anything, even complete nonsense, and the API will accept it provided
		// it can parse out a version number, an object type and a GUID.
		id2 = fmt.Sprintf("https://graph.microsoft.com/v1.0/%s/%s", objectType, guid)
	}

	return json.Marshal(id2)
}

func (id *Id) UnmarshalJSON(data []byte) error {
	if id == nil {
		return nil
	}

	var id2 string
	if err := json.Unmarshal(data, &id2); err != nil {
		return err
	}
	*id = Id(regexp.MustCompile(`/v2/`).ReplaceAllString(id2, `/v1.0/`))

	return nil
}

// Link describes a URI obtained from an OData annotation.
type Link string

func (l *Link) UnmarshalJSON(data []byte) error {
	if l == nil {
		return nil
	}

	var link string
	if err := json.Unmarshal(data, &link); err != nil {
		return err
	}

	// Fix unescaped URLs
	// https://github.com/Azure/azure-sdk-for-go/issues/18809
	u, err := url.ParseRequestURI(link)
	if err != nil {
		// When an invalid URI is returned, we'll return nil instead of raising the error
		return nil
	}
	u.RawQuery = u.Query().Encode()

	// For MS Graph, "v2" is a dev/internal version that sometimes leaks out
	*l = Link(regexp.MustCompile(`/v2/`).ReplaceAllString(u.String(), `/v1.0/`))

	return nil
}

// OData is used to unmarshal OData metadata from an API response.
type OData struct {
	Context      *string `json:"@odata.context"`
	MetadataEtag *string `json:"@odata.metadataEtag"`
	Type         *Type   `json:"@odata.type"`
	Count        *int    `json:"@odata.count"`
	NextLink     *Link   `json:"@odata.nextLink"`
	Delta        *string `json:"@odata.delta"`
	DeltaLink    *Link   `json:"@odata.deltaLink"`
	Id           *Id     `json:"@odata.id"`
	EditLink     *Link   `json:"@odata.editLink"`
	Etag         *string `json:"@odata.etag"`

	Error *Error `json:"-"`

	Value interface{} `json:"value"`
}

func (o *OData) UnmarshalJSON(data []byte) error {
	// Unmarshal using a local type
	type odata OData
	var o2 odata
	if err := json.Unmarshal(data, &o2); err != nil {
		return err
	}
	*o = OData(o2)

	// Look for errors in the "error" and "odata.error" fields
	var e map[string]json.RawMessage
	if err := json.Unmarshal(data, &e); err != nil {
		return err
	}

	for _, k := range []string{"error", "odata.error"} {
		if v, ok := e[k]; ok {
			var e2 Error
			if err := json.Unmarshal(v, &e2); err != nil {
				return err
			}
			o.Error = &e2
			break
		}
	}

	return nil
}
