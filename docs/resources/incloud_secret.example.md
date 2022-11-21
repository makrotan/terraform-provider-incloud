
```terraform
resource "incloud_secret" "my" {
  identifier = "my"
  data = {
    env = "development"
  }
}
```
