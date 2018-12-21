# terraform-chkp
Demo Terraform Provider for the Check Point R80 API

Build the provider by installing go and compiling the provider

Example - /usr/lib/go-1.10/bin/go build -o terraform-provider-chkp

Then in the directory with your terraform example .tf files do the following
Login to a R80 server using the login utility to set a valid session id
This python script can be found in the utils directory

login.py -u admin -p vpn123 -s 10.10.10.10

terraform init

terraform apply

publish.py
