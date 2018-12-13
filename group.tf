resource "chkp_group" "mygroup" {
      name = "My_Terraform_Group"
      color = "light green"
      members = ["${chkp_network.vnet_local.name}",  "${chkp_host.jumphost.name}"]
}
