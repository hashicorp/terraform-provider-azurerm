// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package location

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	return Normalize(old) == Normalize(new)
}

func StateFunc(location interface{}) string {
	input := location.(string)
	return Normalize(input)
}
