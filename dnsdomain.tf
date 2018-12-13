resource "chkp_dnsdomain" "amazondotcom" {
  name = ".amazon.com"
  issubdomain = false
  color = "orange"
  }

  resource "chkp_dnsdomain" "azuredotcom" {
    name = ".azure.com"
    issubdomain = true
    color = "blue"
    }
