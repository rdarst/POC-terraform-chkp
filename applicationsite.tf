resource "chkp_applicationsite" "myapplicationsite" {
      name = "My_Terraform_ApplicationSite"
      color = "light green"
      urllist = ["www.google.com",  "www.ebay.com", "www.amazon.com", "www.core.com"]
      primarycategory = "Social Networking"
      additionalcategories = ["Instant Chat", "Supports Streaming"]
      #tags = ["MyTagTest","MyOtherTag"]
}
