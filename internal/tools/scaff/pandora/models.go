package pandora

// This file contains Go representations of the JSON returned by the Pandora
// Data API (resource-manager endpoints). The struct tags mirror the exact JSON
// field names/casing returned by the server so decoding is deterministic.

// ServicesResponse is returned by GET /v1/resource-manager/services
type ServicesResponse struct {
	Services map[string]ServiceSummary `json:"services"`
}

// ServiceSummary is a single entry in the services list.
type ServiceSummary struct {
	Generate bool   `json:"generate"`
	URI      string `json:"uri"`
}

// ServiceResponse is returned by GET /v1/resource-manager/services/{Service}
type ServiceResponse struct {
	ResourceProvider string                    `json:"resourceProvider"`
	TerraformURI     string                    `json:"terraformUri"`
	Versions         map[string]VersionSummary `json:"versions"`
}

// VersionSummary is a single entry in a service's versions map.
type VersionSummary struct {
	Generate bool   `json:"generate"`
	Preview  bool   `json:"preview"`
	URI      string `json:"uri"`
}

// VersionResponse is returned by GET /v1/resource-manager/services/{Service}/{version}
type VersionResponse struct {
	Resources map[string]ResourceSummary `json:"resources"`
	Source    string                     `json:"source"`
}

// ResourceSummary points at the schema and operations documents for a resource.
type ResourceSummary struct {
	OperationsURI string `json:"operationsUri"`
	SchemaURI     string `json:"schemaUri"`
}

// SchemaResponse is returned by GET .../{Resource}/schema and is the source of
// truth for the model graph, enums (constants) and the resource ID structure.
type SchemaResponse struct {
	Constants   map[string]Constant   `json:"constants"`
	Models      map[string]Model      `json:"models"`
	ResourceIds map[string]ResourceID `json:"resourceIds"`
}

// Constant is an enum type - a set of allowed string (or other) values.
type Constant struct {
	Type   string            `json:"type"`
	Values map[string]string `json:"values"`
}

// Model is a single object/struct in the API model graph.
type Model struct {
	TypeHintIn     *string          `json:"typeHintIn"`
	Fields         map[string]Field `json:"fields"`
	IsParent       bool             `json:"IsParent"`
	ParentTypeName *string          `json:"parentTypeName"`
}

// Field is a single property on a Model.
type Field struct {
	IsTypeHint       bool             `json:"isTypeHint"`
	Description      string           `json:"description"`
	JSONName         string           `json:"jsonName"`
	ObjectDefinition ObjectDefinition `json:"objectDefinition"`
	Optional         bool             `json:"optional"`
	ReadOnly         bool             `json:"readOnly"`
	Required         bool             `json:"required"`
	Sensitive        bool             `json:"sensitive"`
}

// ObjectDefinition describes the type of a Field (or the item type of a List).
// Type is one of the Pandora primitive/complex type names, e.g. String, Integer,
// Boolean, Float, DateTime, Location, Tags, Reference, List, Dictionary,
// UserAssignedIdentityMap, SystemData, RawFile, etc.
// ReferenceName is set when Type == "Reference" (points at a Model or Constant).
// NestedItem is set when Type == "List" or "Dictionary".
type ObjectDefinition struct {
	Type          string            `json:"type"`
	ReferenceName *string           `json:"referenceName"`
	NestedItem    *ObjectDefinition `json:"nestedItem"`
	Nullable      bool              `json:"nullable"`
}

// ResourceID describes a parsed Azure Resource ID and its segments.
type ResourceID struct {
	ConstantNames []string            `json:"constantNames"`
	Constants     map[string]Constant `json:"constants"`
	ID            string              `json:"id"`
	Segments      []Segment           `json:"segments"`
}

// Segment is a single component of a Resource ID.
// Type is one of: Static, ResourceProvider, SubscriptionId, ResourceGroup,
// UserSpecified, Constant, Scope.
type Segment struct {
	ExampleValue string `json:"exampleValue"`
	FixedValue   string `json:"fixedValue"`
	Type         string `json:"type"`
	Name         string `json:"name"`
}

// OperationsResponse is returned by GET .../{Resource}/operations
type OperationsResponse struct {
	Operations map[string]Operation `json:"operations"`
}

// Operation is a single API operation (Create, Read, Update, Delete, List, ...).
type Operation struct {
	ContentType                     string            `json:"contentType"`
	Description                     string            `json:"description"`
	ExpectedStatusCodes             []int             `json:"expectedStatusCodes"`
	LongRunning                     bool              `json:"longRunning"`
	Method                          string            `json:"method"`
	Options                         map[string]Option `json:"options"`
	RequestObject                   *ObjectDefinition `json:"requestObject"`
	ResourceIDName                  *string           `json:"resourceIdName"`
	ResponseObject                  *ObjectDefinition `json:"responseObject"`
	URISuffix                       *string           `json:"uriSuffix"`
	FieldContainingPaginationDetail *string           `json:"fieldContainingPaginationDetails"`
}

// Option is an operation option (e.g. query parameters / OperationOptions).
type Option struct {
	Field            string            `json:"field"`
	ObjectDefinition *ObjectDefinition `json:"objectDefinition"`
	Required         bool              `json:"required"`
}
