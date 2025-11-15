# Provider Documentation Standards

In an effort to keep the [provider documentation](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs) consistent, this page documents some standards that have been agreed on. 

This page will grow over time, and suggestions are welcome!

## Examples

Each resource/data source must include an example, general guidelines for examples are as follows:

- Examples MUST be functional, i.e. if a user copies the example and runs `terraform plan` no errors should be returned.
- Generally the resource instance name should simply be `example`. E.g. `resource "azurerm_resource_group" "example"`.
- All name arguments within the example configuration should be prefixed with `example-` (unless this is disallowed by the naming restrictions), avoid overly complex naming, and ensure any naming restrictions are followed. E.g. `name = example-server`.
- Avoid multiple examples unless a specific configuration is particularly difficult to configure. If there are many complex examples to document, consider using the `examples` folder in the repository instead.
- Examples don't need to include every argument, generally the same configuration as the basic acceptance test will suffice (including any resource dependencies, i.e. the configuration from the template).
- Resource/Data Source examples should not define a `terraform` or `provider` block.

## Arguments

### Ordering

Arguments in the documentation are expected to be ordered as follows:

1. Any arguments that make up the resource's ID, with the last user specified segment (usually `name`) first. E.g. `name` then `resource_group_name`, or `name` then `parent_resource_id`.
2. The `location` field if present.
3. Required arguments, sorted alphabetically.
4. Optional arguments, sorted alphabetically.

### Descriptions

The following conventions apply to argument descriptions:

- Descriptions should be concise, avoid adding too much detail, links to external documentation, etc. If more detail must be added, use a [note](#notes).
- If an argument has `ForceNew: true`, its description must end with `Changing this forces a new <resource name> to be created.`
- If the argument has validation allowing only specific inputs, e.g. `validation.StringInSlice()`, these must be documented using `` Possible values are `value1`, `value2`, and `value3. ``. Other common entries include:
  - Arguments with a single allowed value: `` The only possible values is `value1`. ``
  - Arguments allowing a range of values, e.g. `validation.IntBetween()`: `` Possible values range between `1` and `100`. ``
- If the argument has a default value, this must be documented using `` Defaults to `default1`. ``

Examples:

- `` * `name` - (Required) The name which should be used for this resource. Changing this forces a new resource to be created.``
- `` * `public_network_access` - (Optional) The public network access setting for this resource. Possible values are `Enabled` and `Disabled`. Defaults to `Enabled`. ``
- `` * `disk_size_in_gb` - (Optional) The disk size in gigabytes. Possible values range between `4` and `256`. ``

### Block Arguments

Block arguments must have two entries in the documentation:

1. The initial entry, e.g. `` * `block_argument` - (Optional) A `block_argument` as defined below. ``
2. A subsection, added after all top-level arguments. If multiple blocks are present in the resource, these subsections should be ordered alphabetically. 

Example:

```
## Arguments Reference

`name` - (Required) The name which should be used for this resource.

`block_argument` - (Optional) A `block_argument` as defined below.

`some_other_argument` - (Optional) This argument does something magical.

---

A `block_argument` supports the following:

* `nested_argument_1` - (Required) A nested argument that must be specified.

* `nested_argument_2` - (Optional) A nested argument that may be specified.

## Attributes References

...

```

## Attributes

### Ordering

Attributes in the documentation are expected to be ordered as follows:

1. the `id` attribute.
2. The remaining attributes, sorted alphabetically

### Descriptions

Attribute descriptions should be concise, and must not include possible or default values.

### Block Attributes

Block attributes must have two entries in the documentation:

1. The initial entry, e.g. `` * `block_attribute` - A `block_attribute` as defined below. ``
2. A subsection, added after all top-level attributes. If multiple blocks are present in the resource, these subsections should be ordered alphabetically. 

Example:

```
## Attributes Reference

`id` - The ID of this resource.

`block_attribute` - A `block_attribute` as defined below.

`some_other_attribute` - This attribute returns something magical.

---

A `block_attribute` exports the following:

* `nested_attribute_1` - A very whimsical attribute.

* `nested_attribute_2` - A much more monotonous attribute.

## Timeouts

When documenting timeouts, use the updated link format for all new resources:

- **New resources**: Use `https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts`
- **Existing resources**: Continue using `https://www.terraform.io/language/resources/syntax#operation-timeouts` to maintain consistency

```

## Notes

Note blocks are used to provide additional information to users beyond the basic description of a resource, argument or attribute.

In the past, there have been different approaches to how notes were formatted, some examples are:

- Different words to indicate level of importance, e.g. `Info`, `Important`, `Caution`, and `Be Aware`.
- Capitalisation differences, e.g. `Note:` vs `NOTE:`.
- Whether or not a colon is included, e.g. `Note:` vs `Note`.

Going forward, all notes should follow the exact same format (`(->|~>|!>) **Note:**`) where level of importance is indicated through the different types of notes as documented below.

Breaking changes have previously been added as notes to the resource documentation.
These should no longer be included, instead follow these guidelines:

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
* `type` - (Required) The type. Possible values include `This`, `That`, and `Other`.

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