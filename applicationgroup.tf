resource "chkp_applicationgroup" "myappsitegroup" {
      name = "My_Terraform_AppsiteGroup"
      color = "light green"
      members = ["Facebook",  "Ebay Desktop"]
}
