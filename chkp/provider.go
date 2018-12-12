package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
		    "github.com/hashicorp/terraform/terraform"
		    "./client"
)

func Provider() terraform.ResourceProvider {
        return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHKP_SERVER", nil),
				Description: "server for chkp",
			},

			"sid": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHKP_SID", nil),
				Description: "session for chkp",
			},
		},
                ResourcesMap: map[string]*schema.Resource{
                        "chkp_host": resourceHost(),
                        "chkp_group": resourceGroup(),
                        "chkp_network": resourceNetwork(),
                        "chkp_servicetcp": resourceServiceTcp(),
                        "chkp_serviceudp": resourceServiceUdp(),
                        "chkp_policypackage": resourcePolicyPackage(),
                        "chkp_accesslayer": resourceAccessLayer(),
                        "chkp_accessrulebase": resourceAccessRulebase(),
                        "chkp_accesssection": resourceAccessSection(),
                        "chkp_dyanmicobject": resourceDynamicObject(),
                        "chkp_accessrulebaselist": resourceAccessRulebaseList(),
                        "chkp_accessnatlist": resourceAccessNatRuleList(),
                        "chkp_addressrange": resourceAddressRange(),

                },
                DataSourcesMap:map[string]*schema.Resource{
                       "chkp_host": dataSourceHost(),
                       "chkp_servicetcp": dataSourceServiceTcp(),
                       "chkp_addressrange": dataSourceAddressRange(),
                },

				ConfigureFunc: configureProvider,

        }
}
func configureProvider(d *schema.ResourceData) (interface{}, error) {
	server := d.Get("server").(string)
	sid := d.Get("sid").(string)
	return chkp.NewClientWith(server, sid)
}
