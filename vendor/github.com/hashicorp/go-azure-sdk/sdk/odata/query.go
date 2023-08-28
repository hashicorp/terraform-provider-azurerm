// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MPL-2.0 License. See NOTICE.txt in the project root for license information.

package odata

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// ConsistencyLevel is included in the API request headers when making advanced data queries
type ConsistencyLevel string

const (
	ConsistencyLevelEventual ConsistencyLevel = "eventual"
)

// Metadata specifies the level of control information desired in the response for an API request and is appended to the Accept header
type Metadata string

const (
	MetadataFull    Metadata = "full"
	MetadataMinimal Metadata = "minimal"
	MetadataNone    Metadata = "none"
)

// Query describes OData query parameters that can be included in an API request.
type Query struct {
	// ConsistencyLevel sets the corresponding http header
	ConsistencyLevel ConsistencyLevel

	// Metadata indicates how much control information is requested (services assume "minimal" when not specified)
	Metadata Metadata

	// Count includes a count of the total number of items in a collection alongside the page of data values
	Count bool

	// Expand includes the expanded resource or collection referenced by a single relationship
	Expand Expand

	// Filter retrieves just a subset of a collection, or relationships like members, memberOf, transitiveMembers, and transitiveMemberOf
	Filter string

	// Format specifies the media format of the items returned
	Format Format

	// OrderBy specify the sort order of the items returned
	OrderBy OrderBy

	// Search restricts the results of a request to match a search criterion
	Search string // complicated

	// Select returns a set of properties that are different than the default set for an individual resource or a collection of resources
	Select []string

	// Skip sets the number of items to skip at the start of a collection
	Skip int

	// Top specifies the page size of the result set
	Top int

	// DeltaToken is used to query a delta endpoint
	DeltaToken string
}

// Headers returns a http.Header map containing OData specific headers
func (q Query) Headers() http.Header {
	// Take extra care over canonicalization of header names
	headers := http.Header{}
	headers.Set("Odata-Maxversion", ODataVersion)
	headers.Set("Odata-Version", ODataVersion)

	accept := "application/json; charset=utf-8; IEEE754Compatible=false"
	if q.Metadata != "" {
		accept = fmt.Sprintf("%s; odata.metadata=%s", accept, q.Metadata)
	}
	headers.Set("Accept", accept)

	if q.ConsistencyLevel != "" {
		headers.Set("Consistencylevel", string(q.ConsistencyLevel))
	}

	return headers
}

// AppendHeaders returns the provided http.Header map with OData specific headers appended, for use in requests
func (q Query) AppendHeaders(header http.Header) http.Header {
	if header == nil {
		header = http.Header{}
	}
	for k, v := range q.Headers() {
		if len(v) > 0 {
			header.Set(k, v[0])
		}
	}
	return header
}

// Values returns a url.Values map containing OData specific query parameters
func (q Query) Values() url.Values {
	p := url.Values{}

	if q.Count {
		p.Add("$count", fmt.Sprintf("%t", q.Count))
	}
	if expand := q.Expand.String(); expand != "" {
		p.Add("$expand", expand)
	}
	if q.Filter != "" {
		p.Add("$filter", q.Filter)
	}
	if format := string(q.Format); format != "" {
		p.Add("$format", format)
	}
	if orderBy := q.OrderBy.String(); orderBy != "" {
		p.Add("$orderby", orderBy)
	}
	if q.Search != "" {
		p.Add("$search", fmt.Sprintf(`"%s"`, q.Search))
	}
	if len(q.Select) > 0 {
		p.Add("$select", strings.Join(q.Select, ","))
	}
	if q.Skip > 0 {
		p.Add("$skip", strconv.Itoa(q.Skip))
	}
	if q.Top > 0 {
		p.Add("$top", strconv.Itoa(q.Top))
	}
	if q.DeltaToken != "" {
		p.Add("$deltatoken", q.DeltaToken)
	}

	return p
}

// AppendValues returns the provided url.Values map with OData specific query parameters appended, for use in requests
func (q Query) AppendValues(values url.Values) url.Values {
	if values == nil {
		values = url.Values{}
	}
	for k, v := range q.Values() {
		if len(v) > 0 {
			values.Set(k, v[0])
		}
	}
	return values
}

type Expand struct {
	Relationship string
	Select       []string
}

func (e Expand) String() (val string) {
	val = e.Relationship
	if len(e.Select) > 0 {
		val = fmt.Sprintf("%s($select=%s)", val, strings.Join(e.Select, ","))
	}
	return
}

type Format string

const (
	FormatJson Format = "json"
	FormatAtom Format = "atom"
	FormatXml  Format = "xml"
)

type Direction string

const (
	Ascending  Direction = "asc"
	Descending Direction = "desc"
)

type OrderBy struct {
	Field     string
	Direction Direction
}

func (o OrderBy) String() (val string) {
	val = o.Field
	if val != "" && o.Direction != "" {
		val = fmt.Sprintf("%s %s", val, o.Direction)
	}
	return
}
