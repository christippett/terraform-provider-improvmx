---
page_title: "ImprovMX Domain Data Source - terraform-provider-improvmx"
subcategory: ""
description: |-
  Sample data source in the Terraform provider scaffolding.
---

# Data Source `improvmx_domain`

Sample data source in the Terraform provider scaffolding.

## Example Usage

```terraform
data "improvmx_domain" "example" {
  sample_attribute = "foo"
}
```

## Schema

### Required

- **sample_attribute** (String, Required) Sample attribute.

### Optional

- **id** (String, Optional) The ID of this resource.
