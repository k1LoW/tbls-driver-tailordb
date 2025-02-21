resource "tailor_workspace" "my" {
  name   = "my"
  region = "us-west"
}

output "workspace" {
  value = tailor_workspace.my.id
}
