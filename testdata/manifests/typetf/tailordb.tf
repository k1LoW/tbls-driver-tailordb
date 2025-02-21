resource "tailor_tailordb" "my" {
  workspace_id = tailor_workspace.my.id
  namespace    = "my-tailordb"
}

locals {
  permission_everyone = {
    create = [
      {
        attribute = "everyone"
        permit    = "allow"
      }
    ]
    read = [
      {
        attribute = "everyone"
        permit    = "allow"
      }
    ]
    update = [
      {
        attribute = "everyone"
        permit    = "allow"
      }
    ]
    delete = [
      {
        attribute = "everyone"
        permit    = "allow"
      }
    ]
    admin = [
      {
        attributes = [
          "everyone",
        ]
        permit = "allow"
      }
    ]
  }
}
