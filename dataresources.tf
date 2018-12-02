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

data "chkp_addressrange" "allinternet" {
          name = "All_Internet"
}
