# https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/resources/application

resource "azuread_application" "koki_app_ui" {
  display_name     = "koki_app_ui"
  sign_in_audience = "AzureADMyOrg"

  single_page_application {
    redirect_uris = ["https://koki.sheno.ca", "http://localhost:4200"]
  }
}