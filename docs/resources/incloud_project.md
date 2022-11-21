---
page_title: "incloud_project Resource - terraform-provider-incloud"
subcategory: ""
description: |-
  Projects are containers which contain multiple processes.
---

# Resource `incloud_project`

Projects are containers which contain multiple processes.


## Example Usage


```terraform
resource "incloud_project" "asd" {
  identifier = "asdf"
  name = "from-tf"
}
```




## Argument Reference

The following arguments are supported:

- `identifier` - (Required) [string] 
- `name` - (Required) [string] 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:
