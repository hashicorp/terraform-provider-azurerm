// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package reflect

import (
	"context"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func IsGenericAttrValue(ctx context.Context, target interface{}) bool {
	return reflect.TypeOf((*attr.Value)(nil)) == reflect.TypeOf(target)
}
