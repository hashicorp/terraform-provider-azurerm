// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package markdown

import (
	"regexp"
	"strings"
)

type ImportSection struct {
	heading Heading
	content []string
}

var _ SectionWithTemplate = &ImportSection{}

func (s *ImportSection) Match(line string) bool {
	return regexp.MustCompile(`#+(\s)*import.*`).MatchString(strings.ToLower(line))
}

func (s *ImportSection) SetHeading(line string) {
	s.heading = NewHeading(line)
}

func (s *ImportSection) GetHeading() Heading {
	return s.heading
}

func (s *ImportSection) SetContent(content []string) {
	s.content = content
}

func (s *ImportSection) GetContent() []string {
	return s.content
}

func (s *ImportSection) Template() string {
	return `## Import

[bt]{{ .Name }}[bt] resources can be imported using one of the following methods:

* The [bt]terraform import[bt] CLI command with an [bt]id[bt] string:

  [bt][bt][bt]shell
  terraform import {{ .Name }}.example {{ .ResourceID | id }}
  [bt][bt][bt]

* An [bt]import[bt] block with an [bt]id[bt] argument:
  
  [bt][bt][bt]hcl
  import {
    to = {{ .Name }}.example
    id = "{{ .ResourceID | id }}"
  }
  [bt][bt][bt]

* An [bt]import[bt] block with an [bt]identity[bt] argument:

  [bt][bt][bt]hcl
  import {
    to       = {{ .Name }}.example
    identity = {
      {{ .ResourceID | identity }}
    }
  }
  [bt][bt][bt]
`
}

/*
`azurerm_subnet` resources can be imported using one of the following methods:

* The `terraform import` CLI command with an `id` string

```shell
terraform import azurerm_subnet.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1
```

* An `import` block with an `id` argument

```hcl
import {
  to = azurerm_subnet.example
  id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Network/virtualNetworks/myvnet1/subnets/mysubnet1"
}
```

* An `import` block with an `identity` argument

```hcl
import {
  to = azurerm_subnet.example
  identity = {
    subscription_id      = "00000000-0000-0000-0000-000000000000"
    resource_group_name  = "mygroup1"
    virtual_network_name = "myvnet1"
    subnet_name          = "mysubnet1"
  }
}
*/
