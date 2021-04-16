package pluginsdk

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ResourceTimeout = schema.ResourceTimeout

func DefaultTimeout(tx interface{}) *time.Duration {
	return schema.DefaultTimeout(tx)
}

const (
	TimeoutCreate  = schema.TimeoutCreate
	TimeoutRead    = schema.TimeoutRead
	TimeoutUpdate  = schema.TimeoutUpdate
	TimeoutDelete  = schema.TimeoutDelete
	TimeoutDefault = schema.TimeoutDefault
)
