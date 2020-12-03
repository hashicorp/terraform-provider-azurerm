package example

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/resourceid"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
)

type Logger interface {
	Info(message string)
	InfoF(format string, args ...interface{})
	Warn(message string)
	WarnF(format string, args ...interface{})
}

type ResourceRunFunc func(ctx context.Context, metadata ResourceMetaData) error

type ResourceFunc struct {
	Func    ResourceRunFunc
	Timeout time.Duration
}

type ResourceMetaData struct {
	Client       *clients.Client
	Logger       Logger
	ResourceData *schema.ResourceData
}

func (rmd ResourceMetaData) Decode(input interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		if val, exists := field.Tag.Lookup("hcl"); exists {
			hclValue := rmd.ResourceData.Get(val)

			if v, ok := hclValue.(string); ok {
				reflect.ValueOf(input).Elem().Field(i).SetString(v)
			}
			if v, ok := hclValue.(int); ok {
				reflect.ValueOf(input).Elem().Field(i).SetInt(int64(v))
			}

			// TODO: other types
		}
	}
	return nil
}

func (rmd *ResourceMetaData) Encode(input interface{}) error {
	objType := reflect.TypeOf(input).Elem()
	objVal := reflect.ValueOf(input).Elem()
	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		fieldVal := objVal.Field(i)
		if hclTag, exists := field.Tag.Lookup("hcl"); exists {
			// TODO: make this better
			switch field.Type.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				iv := fieldVal.Int()
				log.Printf("[TOMTOM] Setting %q to %d", hclTag, iv)

				if err := rmd.ResourceData.Set(hclTag, iv); err != nil {
					return err
				}

			case reflect.String:
				sv := fieldVal.String()
				log.Printf("[TOMTOM] Setting %q to %q", hclTag, sv)
				if err := rmd.ResourceData.Set(hclTag, sv); err != nil {
					return err
				}

			default:
				return fmt.Errorf("unknown type %+v for key %q", field.Type.Kind(), hclTag)
			}
		}
	}
	return nil
}

func (rmd ResourceMetaData) SetID(formatter resourceid.Formatter) {
	subscriptionId := rmd.Client.Account.SubscriptionId
	rmd.ResourceData.SetId(formatter.ID(subscriptionId))
}

type Resource interface {
	Arguments() map[string]*schema.Schema
	Attributes() map[string]*schema.Schema

	ResourceType() string

	Create() ResourceFunc
	Read() ResourceFunc
	Delete() ResourceFunc
	IDValidationFunc() schema.SchemaValidateFunc
}

type ResourceWithUpdate interface {
	Update() ResourceFunc
}

type ResourceWrapper struct {
	resource Resource
}

func NewResourceWrapper(resource Resource) ResourceWrapper {
	return ResourceWrapper{
		resource: resource,
	}
}

func (rw ResourceWrapper) Resource() (*schema.Resource, error) {
	resourceSchema, err := rw.schema()
	if err != nil {
		return nil, fmt.Errorf("building Schema: %+v", err)
	}

	var d = func(duration time.Duration) *time.Duration {
		return &duration
	}
	logger := ExampleLogger{}

	resource := schema.Resource{
		Schema: *resourceSchema,

		Create: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta, logger)
			err := rw.resource.Create().Func(ctx, metaData)
			if err != nil {
				return err
			}
			return rw.resource.Read().Func(ctx, metaData)
		},

		// looks like these could be reused, easiest if they're not
		Read: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta, logger)
			return rw.resource.Read().Func(ctx, metaData)
		},
		Delete: func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta, logger)
			return rw.resource.Delete().Func(ctx, metaData)
		},

		Timeouts: &schema.ResourceTimeout{
			Create: d(rw.resource.Create().Timeout),
			Read:   d(rw.resource.Read().Timeout),
			Delete: d(rw.resource.Delete().Timeout),
		},
		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			fn := rw.resource.IDValidationFunc()
			warnings, errors := fn(id, "id")
			if len(warnings) > 0 {
				for _, warning := range warnings {
					logger.Warn(warning)
				}
			}
			if len(errors) > 0 {
				out := ""
				for _, error := range errors {
					out += error.Error()
				}
				return fmt.Errorf(out)
			}

			return err
		}),
	}

	if v, ok := rw.resource.(ResourceWithUpdate); ok {
		resource.Update = func(d *schema.ResourceData, meta interface{}) error {
			ctx, metaData := rw.runArgs(d, meta, logger)
			err := v.Update().Func(ctx, metaData)
			if err != nil {
				return err
			}
			return rw.resource.Read().Func(ctx, metaData)
		}
		resource.Timeouts.Update = d(v.Update().Timeout)
	}

	return &resource, nil
}

func (rw ResourceWrapper) run(in ResourceFunc, logger Logger) func(d *schema.ResourceData, meta interface{}) error {
	return func(d *schema.ResourceData, meta interface{}) error {
		ctx, metaData := rw.runArgs(d, meta, logger)
		err := in.Func(ctx, metaData)
		// TODO: ensure the logger is drained/processed
		return err
	}
}

func (rw ResourceWrapper) schema() (*map[string]*schema.Schema, error) {
	out := make(map[string]*schema.Schema, 0)

	for k, v := range rw.resource.Arguments() {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		// TODO: if readonly

		out[k] = v
	}

	for k, v := range rw.resource.Attributes() {
		if _, alreadyExists := out[k]; alreadyExists {
			return nil, fmt.Errorf("%q already exists in the schema", k)
		}

		// TODO: if editable

		// every attribute has to be computed
		v.Computed = true
		out[k] = v
	}

	return &out, nil
}

func (rw ResourceWrapper) runArgs(d *schema.ResourceData, meta interface{}, logger Logger) (context.Context, ResourceMetaData) {
	ctx := meta.(*clients.Client).StopContext
	client := meta.(*clients.Client)
	metaData := ResourceMetaData{
		Client:       client,
		Logger:       logger,
		ResourceData: d,
	}

	return ctx, metaData
}
