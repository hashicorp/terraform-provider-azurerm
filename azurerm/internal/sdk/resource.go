package sdk

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
)

type resourceBase interface {
	// Arguments is a list of user-configurable (that is: Required, Optional, or Optional and Computed)
	// arguments for this Resource
	Arguments() map[string]*schema.Schema

	// Attributes is a list of read-only (e.g. Computed-only) attributes
	Attributes() map[string]*schema.Schema

	// ModelObject is an instance of the object the Schema is decoded/encoded into
	ModelObject() interface{}

	// ResourceType is the exposed name of this resource (e.g. `azurerm_example`)
	ResourceType() string
}

// A Data Source is an object which looks up information about an existing resource and returns
// this information for use elsewhere
//
// Notably not all Terraform Resources/Azure API's make sense as a Data Source - this information
// has to be available consistently since these are queried on-demand
type DataSource interface {
	resourceBase

	// Read is a ResourceFunc which looks up and sets field values into the Terraform State
	Read() ResourceFunc
}

// A Resource is an object which can be provisioned and managed by Terraform
// that is, Created, Retrieved, Deleted, Imported (and optionally, Updated, by implementing
// the 'ResourceWithUpdate' interface)
//
// It's worth calling out that not all Azure API's make sense as Terraform Resources - as a
// general rule if it supports CR(U)D it could, however.
type Resource interface {
	resourceBase

	// Create will provision this resource using the information from the Terraform Configuration
	// NOTE: the shim layer will automatically call the Read function once this has been created
	// so it's no longer necessary to call this explicitly
	Create() ResourceFunc

	// Read retrieves the latest values for this object and saves them into Terraform's State
	Read() ResourceFunc

	// Delete will remove an existing resource using the information available in Terraform's State
	Delete() ResourceFunc

	// IDValidationFunc returns the SchemaValidateFunc used to validate the ID is valid during
	// `terraform import` - ensuring users don't inadvertently specify the incorrect Resource ID
	IDValidationFunc() schema.SchemaValidateFunc
}

// TODO: ResourceWithCustomizeDiff
// TODO: ResourceWithStateMigration
// TODO: a generic state migration for updating ID's

// ResourceWithUpdate is an optional interface
//
// Notably the Arguments for Resources implementing this interface
// cannot be entirely ForceNew - else this interface implementation
// is superfluous.
type ResourceWithUpdate interface {
	Resource

	// Update will make changes to this resource using the information from the Terraform Configuration/Plan
	// NOTE: the shim layer will automatically call the Read function once this has been created
	// so it's no longer necessary to call this explicitly
	Update() ResourceFunc
}

// ResourceWithDeprecation is an optional interface
//
// Resources implementing this interface will be marked as Deprecated
// and output the DeprecationMessage during Terraform operations.
type ResourceWithDeprecation interface {
	Resource

	// DeprecationMessage returns the Deprecation message for this resource
	// NOTE: this must return a non-empty string
	DeprecationMessage() string
}

// ResourceRunFunc is the function which can be run
// ctx provides a Context instance with the user-provided timeout
// metadata is a reference to an object containing the Client, ResourceData and a Logger
type ResourceRunFunc func(ctx context.Context, metadata ResourceMetaData) error

type ResourceFunc struct {
	// Func is the function which should be called for this Resource Func
	// for example, during Read this is the Read function, during Update this is the Update function
	Func ResourceRunFunc

	// Timeout is the default timeout, which can be overridden by users
	// for this method - in-turn used for the Azure API
	Timeout time.Duration
}

type ResourceMetaData struct {
	// Client is a reference to the Azure Providers Client - providing a typed reference to this object
	Client *clients.Client

	// Logger provides a logger for debug purposes
	Logger Logger

	// ResourceData is a reference to the ResourceData object from Terraform's Plugin SDK
	// This is used to be able to call operations directly should Encode/Decode be insufficient
	// for example, to determine if a field has changes
	ResourceData *schema.ResourceData

	// serializationDebugLogger is used for testing purposes
	serializationDebugLogger Logger
}

// MarkAsGone marks this resource as removed in the Remote API, so this is no longer available
func (rmd ResourceMetaData) MarkAsGone() error {
	rmd.ResourceData.SetId("")
	return nil
}

// ResourceRequiresImport returns an error saying that this resource must be imported with instructions
// on how to do this (namely, using `terraform import`
func (rmd ResourceMetaData) ResourceRequiresImport(resourceName string, idFormatter resourceid.Formatter) error {
	resourceId := idFormatter.ID("") // TODO: remove the dependency on ID in the interface
	return tf.ImportAsExistsError(resourceName, resourceId)
}
