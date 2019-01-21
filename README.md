# terraform-chkp
Proof of concept Terraform Provider for the Check Point R80 API
This provider is for testing purposes only.

Only access based rules can be modifed with the current POC provider

Build the provider by installing go and compiling the provider

Example
```
sudo apt-get install golang-go
/usr/lib/go-1.10/bin/go get github.com/hashicorp/terraform
/usr/lib/go-1.10/bin/go build -o terraform-provider-chkp
```

Then in the directory with your terraform example .tf files do the following
Login to a R80 server using the login utility to set a valid session id
This python script can be found in the utils directory

login.py -u admin -p vpn123 -s 10.10.10.10

Otherwise setup the two environment variables to allow terraform to pick up your R80 session details. 
For example
export CHKP_SID="oH9f7BaC-63kcF2fg3qokliwHPrXtCEIf4V8zvIpTmE"
export CHKP_SERVER="https://10.10.10.10/web_api"


terraform init

terraform apply

publish.py

Use destroy to remove what was created

terraform destroy

Notes -

When removing or adding rules via the rule list, this provider will not preserve UUID's of rules that were modifed.  This also applies to NAT rules as well.

