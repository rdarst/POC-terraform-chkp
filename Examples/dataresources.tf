data "chkp_servicetcp" "http" {
          name = "http"
}

data "chkp_servicetcp" "https" {
          name = "https"
}

data "chkp_servicetcp" "ssh" {
          name = "ssh"
}

data "chkp_servicetcp" "ftp" {
          name = "ftp"
}

data "chkp_serviceudp" "domain-udp" {
          name = "domain-udp"
}

data "chkp_addressrange" "allinternet" {
          name = "All_Internet"
}

data "chkp_securityzone" "DMZZone" {
          name = "DMZZone"
}

data "chkp_dynamicobject" "DMZNet" {
          name = "DMZNet"
}

#data "chkp_dnsdomain" "googledotcom" {
#          name = ".google.com"
#}
