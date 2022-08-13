# https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application
# https://docs.microsoft.com/en-us/graph/migrate-azure-ad-graph-configure-permissions
# pre-requisite: API permissions of the SP :Microsoft Graph (2) >Application.ReadWrite.All Directory.ReadWrite.All

data "azuread_client_config" "current" {}

resource "random_uuid" "app_role_id" {}

resource "azuread_application" "koki_app_ui" {
  display_name     = "koki-app-ui"
  sign_in_audience = "AzureADMyOrg"
  owners           = [data.azuread_client_config.current.object_id, local.tenant_owner_object_id]


  app_role {
    allowed_member_types = ["User", "Application"]
    description          = "Approver has the ability to approve privileges assignment/removal"
    display_name         = "Approver"
    enabled              = true
    id                   = random_uuid.app_role_id.result
    value                = "Privilege.Approve"
  }

  single_page_application {
    redirect_uris = ["https://koki.sheno.ca/", "http://localhost:4200/"]
  }
}