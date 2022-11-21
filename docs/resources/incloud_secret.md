---
page_title: "incloud_secret Resource - terraform-provider-incloud"
subcategory: ""
description: |-
  
---

# Resource `incloud_secret`




## Example Usage


```terraform
resource "incloud_secret" "my" {
  identifier = "my"
  data = {
    env = "development"
  }
}
```




## Argument Reference

The following arguments are supported:

- `identifier` - (Required) [string] 
- `data` - (Required) [map] 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
