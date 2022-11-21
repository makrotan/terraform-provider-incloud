---
page_title: "incloud_app Resource - terraform-provider-incloud"
subcategory: ""
description: |-
  
---

# Resource `incloud_app`




## Example Usage


```terraform
resource "incloud_app" "test" {
  identifier = "from-tf"
  name = "from-tf"
  git_url = "https://gitlab.com/klamar2/incloud-demo-repo.git"
}
```




## Argument Reference

The following arguments are supported:

- `identifier` - (Required) [string] 
- `name` - (Required) [string] 
- `git_url` - (Required) [string] 
- `branch` - [string] 

## Attributes Reference

In addition to all the arguments above, the following attributes are exported:

- `status` - (Required) [string] 