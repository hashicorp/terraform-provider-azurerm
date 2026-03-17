# Provider Documentation Standards

In an effort to keep the [provider documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs) consistent, this page documents some standards that have been agreed on.

This document defines standards for resource and data source reference documentation in `website\docs\r` and `website\docs\d`. It does not define standards for other documentation types such as `guides`, `functions`, `actions`, `upgrade guides`, or `list` pages.

This page will grow over time, and suggestions are welcome!

## Documentation Locations

Resource and data source reference documentation is located under the `website\docs` directory in the repository. This documentation is split between `resources` and `data sources` which are kept in different sub-directories of the `website\docs` directory.

- Resource documentation is in the `website\docs\r` directory.
- Data source documentation is in the `website\docs\d` directory.

Reference documentation should follow the name of the Terraform resource or data source it is documenting.

- If you are documenting the resource `azurerm_example` the documentation should be named `example.html.markdown` and placed in the `website\docs\r` directory.
- If you are documenting the data source `azurerm_example` the documentation should be named `example.html.markdown` and placed in the `website\docs\d` directory.

## Front Matter

Each resource/data source must include the below Front Matter at the begining of the documentation file e.g., `example.html.markdown`.

The `subcategory` value should come from the website category defined for the service. To find the allowed values, check `website/allowed-subcategories`. If you scaffold the documentation using `make scaffold-website`, the generated front matter will also tell you which category to use. If the service supports multiple website categories, match the existing documentation for that service.

For resources the front matter should be defined as:

```markdown
---
subcategory: "ExampleService"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_example"
description: |-
  Manages an Example.
---
```

For data sources the front matter should be defined as:

```markdown
---
subcategory: "ExampleService"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_example"
description: |-
  Gets information about an existing Example.
---
```

## Examples

Each resource/data source must include an example, general guidelines for examples are as follows:

- Examples MUST be functional, meaning that if a user copies the example and runs `terraform plan`, no errors should be returned.
- Generally the resource instance name should simply be `example`. e.g. `resource "azurerm_resource_group" "example"`.
- All name arguments within the example configuration should use simple example values that match the resource being defined. Where naming restrictions and field validation allow, prefer values prefixed with `example-`. If a field's validation or naming restrictions do not allow that pattern, use the simplest valid value for that field. Avoid overly complex naming, and ensure any naming restrictions and validation are followed. e.g. `name = example-resource-group`.
- Avoid multiple examples unless a specific configuration is particularly difficult to configure. If there are many complex examples to document, consider using the `examples` folder in the repository instead.
- Examples don't need to include every argument, generally the same configuration as the basic acceptance test will suffice (including any resource dependencies, such as the configuration from the template).
- Resource/Data Source examples should not define a `terraform` or `provider` block.

## Code Fences

The following conventions apply to code fences:

- Use the most specific code fence language that matches the snippet.
- Terraform configuration should use `hcl` code fences. Do not use `terraform` code fences for HCL configuration blocks.
- Keep Terraform examples copy/pasteable and self-contained.

## Arguments

The ``## Arguments Reference`` section is used to document fields that can be set by the user in the Terraform configuration.

### Descriptions

The following conventions apply to argument descriptions:

- Descriptions should be concise, avoid adding too much detail, links to external documentation, etc. If more detail must be added, use a [note](#notes).
- If an argument has `ForceNew: true`, its description must end with `Changing this forces a new resource to be created.`
- If the argument has a default value, this must be documented using `` Defaults to `default1`. ``
- If the argument has validation allowing only specific inputs, e.g. `validation.StringInSlice()`, these must be documented using `` Possible values are `value1`, `value2`, and `value3. ``.
  * Other common entries include:
    - Arguments with a single allowed value: `` The only possible value is `value1`. ``
    - Arguments allowing a range of values, e.g. `validation.IntBetween()`: `` Possible values range between `1` and `100`. ``


Examples:

- `name` - (Required) The name which should be used for this resource.
- `argument_enabled` - (Optional) Should this argument be enabled? Possible values are `true` and `false`. Defaults to `false`.
- `argument_in_gb` - (Optional) The argument in gigabytes. Possible values range between `4` and `256`.

### Ordering

Arguments in the documentation are expected to be ordered as follows:

- Any arguments that make up the resource's ID, with the last user specified segment (usually `name`) first. e.g. `name` then `resource_group_name`, or `name` then `parent_resource_id`.
- The `location` field if present.
- Required arguments, sorted alphabetically.
- Optional arguments, sorted alphabetically, with the exception of `tags`, which must always be documented last.

> **Note:** This ordering applies to both `typed` and `untyped` implementations. Even when typed resources or data sources surface computed or optional fields via `Attributes()`/`model` structs, the published documentation must still follow the sequence described above.

### Block Arguments

Block arguments require two entries in the documentation:

- The initial entry, e.g. `` * `block_argument` - (Optional) A `block_argument` as defined below. ``, using the correct indefinite article for the block name (`A` or `An`, as appropriate).
- A subsection, added after all top-level arguments. If multiple blocks are present in the resource, these subsections should be ordered alphabetically.

Arguments within a block subsection are expected to be ordered as follows:

- Required arguments, sorted alphabetically.
- Optional arguments, sorted alphabetically.

Example:

```markdown
## Arguments Reference

* `name` - (Required) The name which should be used for this resource.

* `resource_group_name` - (Required) The name which should be used for this resource.

* `location` - (Required) The name which should be used for this resource.

* `argument` - (Required) This argument does something nifty.

* `block_argument` - (Optional) A `block_argument` as defined below.

* `some_other_argument` - (Optional) This argument does something magical.

---

A `block_argument` supports the following:

* `block_argument` - (Required) Specifies the block_argument of this nested item.

* `some_other_block_argument` - (Required) Specifies the some_other_block_argument of this nested item exists.

* `optional_block_argument` - (Optional) Specifies the optional_block_argument of this nested item exists.

* `some_other_optional_block_argument` - (Optional) Specifies the some_other_optional_block_argument of this nested item exists.

## Attributes Reference

...

```

## Attributes

The ``## Attributes Reference`` section is used to document `Computed` fields that are returned by Terraform after the resource or data source is read.

Directly after the initial heading of ``## Attributes Reference`` you must include the exact text ``In addition to the Arguments listed above - the following Attributes are exported:``.

### Descriptions

Attribute descriptions should be concise, and must not include `possible`, `default`, or `ForceNew` values.

### Ordering

Attributes in the documentation are expected to be ordered as follows:

- The `id` attribute.
- The remaining attributes are sorted alphabetically, no exceptions.

### Block Attributes

Block attributes must have two entries in the documentation:

- The initial entry in the ``## Attributes Reference``, e.g. ``* `block_attribute` - A `block_attribute` as defined below.``, using the correct indefinite article for the block name (`A` or `An`, as appropriate).
- A subsection, added after all top-level attributes. If multiple blocks are present in the resource, these subsections should be ordered alphabetically.

Attributes within a block subsection are expected to be ordered as follows:

- The required attributes, sorted alphabetically.
- The optional attributes, sorted alphabetically.

Each block subsection is ordered independently of the top-level `## Attributes Reference` list.

Example:

```markdown
## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of this resource.

* `block_attribute` - A `block_attribute` as defined below.

* `example_attribute` - This attribute returns an example value.

* `some_other_attribute` - This attribute returns something magical.

---

A `block_attribute` block exports the following:

* `required_attribute` - This attribute returns a required example value.

* `some_other_required_attribute` - This attribute returns another required example value.

* `optional_attribute` - This attribute returns an optional example value.

* `some_other_optional_attribute` - This attribute returns another optional example value.
...
```

## Notes

Note blocks are used to provide additional information to users beyond the basic description of a resource, argument or attribute.

In the past, there have been different approaches to how notes were formatted, some examples are:

- Different words to indicate level of importance, e.g. `Info`, `Important`, `Caution`, and `Be Aware`.
- Capitalization differences, e.g. `Note:` vs `NOTE:`.
- Whether or not a colon is included, e.g. `Note:` vs `Note`.

Going forward, all notes should follow the exact same format (`(->|~>|!>) **Note:**`) where level of importance is indicated through the different types of notes as documented below.

Breaking changes have previously been added as notes to the resource documentation.
These should no longer be included. Instead, follow these guidelines:

- Breaking changes in a minor version should be added to the top of the changelog.
- Breaking changes in a major version should be added to the upgrade guide.

> We may revisit the guidelines above and/or add a specific place in the documentation for all breaking changes in minor versions.

<!--
    - TODO: Considerations for when to add notes? We probably don't want to overdo it (More relevant to informational notes)
-->

### Informational Note

Informational note blocks should generally be used when a note provides additional useful information, recommendations and/or tips to the user.

To add an informational note, use `-> **Note:**`, within the Terraform registry documentation this will template as a block with an info icon.

For example, extra information on the supported values for an argument, possibly linking to external documentation for the resource/service:

```markdown
* `type` - (Required) The type. Possible values are `This`, `That`, and `Other`.

-> **Note:** More information on each of the supported types can be found in [type documentation](link-to-additional-info)
```

### Warning Note

Warning note blocks should generally be used when a note provides information that the user will need to avoid certain errors, however if these errors are encountered they should not break anything or cause irreversible changes.

To add a warning note, use `~> **Note:**`, within the Terraform registry documentation this will template as a block with a warning icon.

For example, an argument that is optional but required when another argument is set to `true`:

```markdown
* `optional_argument_enabled` - (Optional) Is the optional argument enabled? Defaults to `false`.

* `optional_argument` - (Optional) An optional argument.

~> **Note:** The argument `optional_argument` is required when `optional_argument_enabled` is set to `true`.
```

### Caution Note

Caution note blocks should generally be used when a note provides critical information on potential irreversible changes, data loss or other things that can negatively affect a user's environment.

To add a caution note, use `!> **Note:**`, within the Terraform registry documentation this will template as a block with a caution icon.

For example, an argument that when set to `true` cannot be reversed without recreating the resource:

```markdown
* `irreversible_argument_enabled` - (Optional) Is irreversible argument enabled? Defaults to `false`.

!> **Note:** The argument `irreversible_argument_enabled` cannot be disabled after being enabled.
```
