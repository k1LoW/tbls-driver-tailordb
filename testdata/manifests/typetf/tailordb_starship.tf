resource "tailor_tailordb_type" "starship" {
  workspace_id = tailor_workspace.my.id
  namespace    = tailor_tailordb.my.namespace
  name         = "Starship"
  description  = <<EOF
A single transport craft that has hyperdrive capability.
Designed to include all the types supported by TailorDB.
EOF

  fields = {
    name = {
      type        = "string"
      description = "Name of the starship"
      required    = true
    }
    serialNumber = {
      type        = "uuid"
      description = "Serial number of the starship"
      required    = true
      index       = true
    }
    status = {
      type        = "enum"
      description = "Status of the starship"
      allowed_values = [
        { value = "ACTIVE", description = "Active starship" },
        { value = "INACTIVE", description = "Inactive starship" },
      ]
    }
    isArmed = {
      type        = "boolean"
      description = "Whether the starship is armed"
    }
    crew = {
      type        = "integer"
      description = "Number of crew members"
    }
    length = {
      type        = "float"
      description = "Length of the starship in meters"
    }
    manufactureDate = {
      type        = "date"
      description = "Date of manufacture"
    }
    maintenanceTime = {
      type        = "time"
      description = "Time required for maintenance"
    }
    lastLaunchDateTime = {
      type        = "datetime"
      description = "Last launch date and time"
    }
    manufacturers = {
      type        = "string"
      description = "Manufacturers of the starship"
      array       = true
    }
    hyperdriveRating = {
      type        = "nested"
      description = "Hyperdrive rating of the starship"
      fields = {
        main = {
          type        = "float"
          description = "Main hyperdrive rating"
        }
        backup = {
          type        = "float"
          description = "Backup hyperdrive rating"
        }
      }
    }
  }

  type_permission = local.permission_everyone
}
