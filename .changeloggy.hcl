archive_entries = true
changelog_file            = "CHANGELOG.md"
default_version_increment = "minor"

format {
  date = "January 2, 2006"

  template = <<-TEMPLATE
## {{ .Version }} ({{ .Date }})
{{ if .Notes }}
{{ range .Notes }}* {{ .Body }} ([#{{ .PR }}](https://github.com/hashicorp/terraform-provider-azurerm/issues/{{ .PR }}))
{{ end }}{{ end }}{{ range .Kinds }}
{{ .Heading }}:

{{ range .Entries }}* {{ .Body }} ([#{{ .PR }}](https://github.com/hashicorp/terraform-provider-azurerm/issues/{{ .PR }}))
{{ end }}{{ end }}
TEMPLATE
}

kind "breaking-change" {
  heading = "BREAKING CHANGES"
  priority = 1

  type "breaking" {
    priority = 1
  }
}

kind "feature" {
  heading  = "FEATURES"
  priority = 2

  type "new-action" {
    regex    = "^\\*\\*New Action\\*\\*: `azurerm_\\w+`$"
    example  = "**New Action**: `azurerm_example`"
    priority = 1
  }

  type "new-data-source" {
    regex    = "^\\*\\*New Data Source\\*\\*: `azurerm_\\w+`$"
    example  = "**New Data Source**: `azurerm_example`"
    priority = 2
  }

  type "new-list-resource" {
    regex    = "^\\*\\*New List Resource\\*\\*: `azurerm_\\w+`$"
    example  = "**New List Resource**: `azurerm_example`"
    priority = 3
  }

  type "new-resource" {
    regex    = "^\\*\\*New Resource\\*\\*: `azurerm_\\w+`$"
    example  = "**New Resource**: `azurerm_example`"
    priority = 4
  }

  type "new-function" {
    regex    = "^\\*\\*New Function\\*\\*: `\\w+`$"
    example  = "**New Function**: `example_function`"
    priority = 5
  }
}

kind "enhancement" {
  heading  = "ENHANCEMENTS"
  priority = 3

  type "dependency" {
    regex    = "^dependencies: `[^`]+` - update( API version)? to `[^`]+`$"
    example  = "dependencies: `sql` - update API version to `2025-01-01`"
    priority = 1
  }

  type "new-property" {
    regex    = "^`(Data Source: |List Resource: |Action: )?azurerm_[^`]+` - add support for the (`[^`]+`(, `[^`]+`)*(,? and `[^`]+`)?) (propert(y|ies)|blocks?)$"
    example  = "`azurerm_example` - add support for the `new_field` property"
    priority = 2
  }

  type "resource-enhancement" {
    regex    = "^`(Data Source: |List Resource: |Action: )?azurerm_[^`]+` - .+$"
    example  = "`azurerm_example` - improve validation for the `name` property"
    priority = 3
  }

  type "other-enhancement" {
    regex    = "^.+$"
    example  = "`go-azure-sdk` - `Delete` operations now poll on asynchronous operation URLs"
    priority = 4
  }
}

kind "bug" {
  heading  = "BUG FIXES"
  priority = 4

  type "provider-fix" {
    regex   = "^provider: `[^`]+` - .+$"
    example = "provider: `provider_feature` - now correctly reads the default value from an environment variable"
    priority = 1
  }

  type "resource-fix" {
    regex   = "^`(Data Source: |List Resource: |Action: )?azurerm_[^`]+` - .+$"
    example = "`azurerm_example` - fix panic when `name` is empty"
    priority = 2
  }

  type "other-fix" {
    regex   = "^.+$"
    example = "`go-azure-sdk` - `Delete` operations now poll on asynchronous operation URLs"
    priority = 3
  }
}
