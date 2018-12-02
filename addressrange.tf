resource "chkp_addressrange" "testaddressrange" {
      ipv4addressfirst = "9.9.9.0"
      ipv4addresslast = "9.9.9.11"
      name = "testaddressrange"
      color = "yellow"
      natsettings {
        autorule = false }
        }
