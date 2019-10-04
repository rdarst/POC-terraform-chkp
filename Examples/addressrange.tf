resource "chkp_addressrange" "testaddressrange" {
      ipv4addressfirst = "9.9.9.0"
      ipv4addresslast = "9.9.9.11"
      name = "testaddressrange"
      color = "pink"
      natsettings {
      autorule = true
      hidebehind = "gateway"
      method = "hide"
      installon = "All" 
            }
        }

resource "chkp_addressrange" "mytestaddressrange" {
      ipv4addressfirst = "9.9.9.12"
      ipv4addresslast = "9.9.9.22"
      name = "mytestaddressrange"
      color = "pink"
      }
