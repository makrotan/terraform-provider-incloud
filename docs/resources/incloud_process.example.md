
```terraform
resource "incloud_process" "foo" {
  project_id = incloud_project.asd.identifier
  identifier = "xcyvxcvyxc"
  name = "from-tf"

  spec = <<-EOL
    send_slack_notification:
      app: incloud-demo-repo
      function: slack_notify
      parameters:
        text: Hello {{ parameters.name }}, This is a test
        webhook_url: '{{ fetch(''secret/makrotan-slack'').slack_webhook_url }}'
    send_slack_notification2:
      after: send_slack_notification
      app: incloud-demo-repo
      function: slack_notify
      parameters:
        webhook_url: "{{ fetch('secret/makrotan-slack').slack_webhook_url }}"
        text: "the previous webhook returned http status: {{ step.send_slack_notification.result }}"

  EOL
}
```
