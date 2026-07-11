// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package ir

// ResourceIR is the fully-resolved, deterministic representation of a resource,
// built from the Pandora schema + operations documents and consumed by the code
// generation templates.
type ResourceIR struct {
	// Provider / naming
	ProviderName      string // e.g. "azurerm"
	ProviderGithubOrg string // e.g. "hashicorp"
	ServicePackage    string // service package directory, e.g. "redhatopenshift"
	ServiceName       string // Pandora service / client field name, e.g. "RedHatOpenShift"
	ResourceProvider  string // ARM RP namespace, e.g. "Microsoft.RedHatOpenShift"
	APIVersion        string // e.g. "2025-07-25"
	PandoraResource   string // Pandora resource key, e.g. "OpenShiftClusters"

	// User-facing resource identity
	Name          string // Go Camel name, e.g. "RedHatOpenShiftCluster"
	TerraformType string // e.g. "azurerm_redhat_openshift_cluster"

	// SDK
	SDKPackage    string // go-azure-sdk package, e.g. "openshiftclusters"
	SDKImportPath string // full import path
	ClientField   string // client accessor field, e.g. "OpenShiftClustersClient"

	// Resource ID
	IDTypeName       string   // e.g. "OpenShiftClusterId"
	IDParseFunc      string   // e.g. "ParseOpenShiftClusterID"
	IDNewFunc        string   // e.g. "NewOpenShiftClusterID"
	IDValidateFunc   string   // e.g. "ValidateOpenShiftClusterID"
	IDNameSegment    string   // Go name of the final user-specified segment, e.g. "OpenShiftClusterName"
	IDNewArgs        []string // ordered Go expressions for New*ID() sourced from the create model var
	HasSubscription  bool     // ID contains a SubscriptionId segment
	HasResourceGroup bool     // ID contains a ResourceGroup segment

	// CRUD operations
	Updatable       bool
	CreateOp        string // operation name, e.g. "CreateOrUpdate"
	ReadOp          string // e.g. "Get"
	UpdateOp        string // e.g. "Update"
	DeleteOp        string // e.g. "Delete"
	CreateLRO       bool
	UpdateLRO       bool
	DeleteLRO       bool
	CreateModel     string // request model name, e.g. "OpenShiftCluster"
	UpdateModel     string // update request model name, e.g. "OpenShiftClusterUpdate"
	ReadModel       string // response model name, e.g. "OpenShiftCluster"
	PropertiesModel string // the Properties sub-model name, e.g. "OpenShiftClusterProperties"
	// Update model details (drive Update() body generation)
	UpdateHasTags          bool            // update model carries a Tags field
	UpdatePropertiesModel  string          // Properties sub-model name in the update model
	UpdatablePropSDKFields map[string]bool // SDK field names present in the update Properties model

	// List resource support
	IsListable            bool   // has subscription/resource-group or parent list operations
	ListBySubscriptionOp  string // e.g. "List"
	ListByResourceGroupOp string // e.g. "ListByResourceGroup"
	// Parent-scoped list (child resources)
	ListByParentOp     string // list operation scoped to a parent, e.g. "List"
	ParentIDType       string // parent scope ID type, e.g. "StorageMoverId"
	ParentListAttr     string // config attribute name, e.g. "storage_mover_id"
	ParentParseFunc    string // e.g. "ParseStorageMoverID"
	ParentValidateFunc string // e.g. "ValidateStorageMoverID"
	ParentPackage      string // package of the parent ID funcs (untyped); e.g. "virtualwans"
	ParentImportPath   string // full import path of the parent ID package (untyped)

	// Untyped list support (populated when upgrading a native Plugin SDK resource).
	Untyped         bool   // render the native Plugin SDK list resource shape
	ConstructorFunc string // e.g. "resourceVirtualHubIP"
	FlattenFunc     string // e.g. "resourceVirtualHubIPFlatten"
	FlattenIDValue  bool   // flatten takes its id by value, so the list passes *id
	IDPackage       string // ID type package; may differ from SDKPackage (e.g. "commonids")
	IDImportPath    string // full import path of the ID package
	TestStructName  string // acceptance-test struct name (may differ in casing from Name+"Resource")

	// Schema (populated in the model-graph walk)
	ModelStructName string        // top-level model struct, e.g. "RedHatOpenShiftClusterModel"
	TopLevel        []*Property   // flattened top-level schema properties
	Blocks          []*BlockModel // nested block model structs to emit (deduped)
}

// Property is a single schema property (attribute or nested block) in the IR.
type Property struct {
	// Terraform
	TFName      string // snake_case schema key, e.g. "vm_size"
	TFType      string // pluginsdk type, e.g. "TypeString", "TypeList", "TypeBool", "TypeInt"
	Description string // schema Description text, sourced from the Pandora field description
	Required    bool
	Optional    bool
	Computed    bool
	ForceNew    bool
	Sensitive   bool
	MaxItems    int // 1 for single-nested Reference blocks, 0 otherwise

	// Model (Go struct field)
	GoField string // Go struct field name, e.g. "VmSize"
	GoType  string // Go type for the model struct field, e.g. "string", "int64", "bool", "[]MainProfile"

	// Nesting
	IsBlock   bool   // true if this is a nested block (Reference/List of Reference)
	BlockName string // referenced BlockModel name if IsBlock

	// Enum
	IsEnum     bool     // backed by a Pandora constant
	EnumType   string   // constant reference name, e.g. "EncryptionAtHost"
	EnumValues []string // allowed values, for validation

	// Source
	SDKField        string // originating SDK model field name, e.g. "VMSize"
	JSONName        string // original JSON name
	UnderProperties bool   // true if this property lives under the ARM Properties envelope
}

// BlockModel is a nested object model to emit as a Go struct + expand/flatten pair.
type BlockModel struct {
	Name       string      // Go struct name, e.g. "MainProfile"
	SDKModel   string      // originating SDK model name, e.g. "MasterProfile"
	Properties []*Property // fields of the block
}
