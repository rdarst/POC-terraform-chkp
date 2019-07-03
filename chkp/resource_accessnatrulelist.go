package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
      //  "strings"
        "reflect"
)

func resourceAccessNatRuleList() *schema.Resource {
        return &schema.Resource{
                Create: resourceAccessNatRuleListCreate,
                Read:   resourceAccessNatRuleListRead,
                Update: resourceAccessNatRuleListUpdate,
                Delete: resourceAccessNatRuleListDelete,

                Schema: map[string]*schema.Schema{

                  "package": {
                          Type:     schema.TypeString,
                          Optional: true,
                            },

						            "natlist": {
                          Type:     schema.TypeList,
                          Optional: true,
                          Elem: &schema.Resource{
                            Schema: map[string]*schema.Schema{
                              "enabled": {
                                      Type:     schema.TypeBool,
                                      Optional: true,
                                      Default: true,
                                        },
                             "installon": {
                                       Type:     schema.TypeSet,
                                       Optional: true,
                                       Elem: &schema.Schema{ Type: schema.TypeString },
                                         },
                             "method": {
                                     Type:     schema.TypeString,
                                     Required: true,
                                       },
                             "originaldestination": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "originalservice": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "originalsource": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "translateddestination": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "translatedservice": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "translatedsource": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "comments": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "uid": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                     Computed: true,
                             },
                           },
                         },
                     },
               },
       }
}

func resourceAccessNatRuleListCreate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)

  policypackage := d.Get("package").(string)
  d.Set("package", policypackage)
  d.SetId(policypackage)
  // Pull in NAT Rulebase rules
  natlist := d.Get("natlist").([]interface{})

  // Create variable to hold uids from returned rules
  uidlist := make([]interface{}, 0, len(natlist))

  for q := range natlist {
    natelement := natlist[q].(map[string]interface{})
          var natrulebase = chkp.AccessRulebaseNATList{}
          // Place Nat Rules in the between the section titles


            positionreturn := chkp.Position{
              Above: "Terraform NAT Rules Above This Position",
          }

          natrulebase.Position = positionreturn
          natrulebase.Package = policypackage
          natrulebase.Enabled = natelement["enabled"].(bool)
          // Get the list of Install-on Targets
          installon := natelement["installon"].(*schema.Set).List()
          installlist := make([]string, 0, len(installon))
          for _, v := range installon {
            val, ok := v.(string)
            if ok && val != "" {
              installlist = append(installlist, v.(string))
            }
          }

          if len(installlist) > 0 {
                natrulebase.Installon = installlist
              }

          natrulebase.OriginalDestination = natelement["originaldestination"].(string)
          natrulebase.OriginalService = natelement["originalservice"].(string)
          natrulebase.OriginalSource = natelement["originalsource"].(string)
          natrulebase.TranslatedDestination = natelement["translateddestination"].(string)
          natrulebase.TranslatedService = natelement["translatedservice"].(string)
          natrulebase.TranslatedSource = natelement["translatedsource"].(string)
          natrulebase.Comments = natelement["comments"].(string)
          natrulebase.Method = natelement["method"].(string)
          id, err := client.CreateAccessNATRulebaseList(natrulebase)
          //Read in from the API Output
          readNATRulebase := chkp.NATRulebaseResult{}
          json.Unmarshal(id, &readNATRulebase)

          natreturn := make(map[string]interface{})
          natreturn["uid"] = readNATRulebase.Uid
          natreturn["enabled"] = readNATRulebase.Enabled
          natreturn["install-on"] = readNATRulebase.Installon
          natreturn["original-destination"] = readNATRulebase.OriginalDestination
          natreturn["original-service"] = readNATRulebase.OriginalService
          natreturn["original-source"] = readNATRulebase.OriginalSource
          natreturn["translated-destination"] = readNATRulebase.TranslatedDestination
          natreturn["translated-service"] = readNATRulebase.TranslatedService
          natreturn["translated-source"] = readNATRulebase.TranslatedSource
          natreturn["comments"] = readNATRulebase.Comments
          uidlist = append(uidlist, natreturn)
          if err != nil {
        		return err
        	}
}

d.Set("natlist", uidlist)
return nil
}

func resourceAccessNatRuleListRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)

  // Set API limit on the number of records returned default is 50 max is 500
  limit := 100
  packageuid := d.Get("package").(string)
  id, err := client.ShowNATRulebaseList(packageuid, limit, 0)
  if err != nil {
    status := err.Error()
    if (status == "404") {
          // If the object is not found remove it from state
          d.SetId("")
          return nil
    } else {
      return err
    }
  }
  readNATRulebase := chkp.NATRulebaseResultRead{}
  json.Unmarshal(id, &readNATRulebase)

  var result map[string]interface{}
  json.Unmarshal(id, &result)

  rules := result["rulebase"].([]interface{})
  total := readNATRulebase.Total
  to := readNATRulebase.To
  offset := to

  rulelistread := make([]interface{}, 0, total)
  for t := range rules {
    rulelistread = append(rulelistread, rules[t].(map[string]interface{}))
  }
  // Check to see if we need to gather more rules
  if total > to {
  done := 0
  for done < 1 {
    id3, err := client.ShowNATRulebaseList(packageuid, limit, offset)
    var resultappend map[string]interface{}
    readNATRulebaseappend := chkp.NATRulebaseResultRead{}
    json.Unmarshal(id3, &resultappend)
    json.Unmarshal(id3, &readNATRulebaseappend)
    rules := resultappend["rulebase"].([]interface{})
    for t := range rules {
      rulelistread = append(rulelistread, rules[t].(map[string]interface{}))
    }
    if err != nil {
  		return err
  	}
    total = readNATRulebaseappend.Total
    to = readNATRulebaseappend.To
    offset = to

    if total == to { done = 2}
    //done = 2
      }
    }

    // build the interface for the nat list
    natlist := make([]interface{}, 0, total)

    for q := range rulelistread {
      ruleelement := rulelistread[q].(map[string]interface{})
            if ruleelement["name"].(string) == "Terraform NAT Rules Below This Position" {

              terraformnatrules := ruleelement["rulebase"].([]interface{})
              for t := range terraformnatrules {
                natelement := terraformnatrules[t].(map[string]interface{})
                natreturn := make(map[string]interface{})
                natreturn["uid"] = natelement["uid"]
                natreturn["comments"] = natelement["comments"].(string)
                natreturn["enabled"] = natelement["enabled"].(bool)
                natreturn["method"] = natelement["method"].(string)
                natreturn["originalsource"] = natelement["original-source"].(map[string]interface{})["name"].(string)
                natreturn["originaldestination"] = natelement["original-destination"].(map[string]interface{})["name"].(string)
                natreturn["originalservice"] = natelement["original-service"].(map[string]interface{})["name"].(string)
                natreturn["translatedsource"] = natelement["translated-source"].(map[string]interface{})["name"].(string)
                natreturn["translateddestination"] = natelement["translated-destination"].(map[string]interface{})["name"].(string)
                natreturn["translatedservice"] = natelement["translated-service"].(map[string]interface{})["name"].(string)
                //Gather Install on Targets
                installon := natelement["install-on"]
                installonlist := installon.([]interface{})
                installontargets := make([]string, 0, len(installonlist))
                for i := range installonlist {
                      installonelements := installonlist[i].(map[string]interface{})
                      if installonelements["name"] != "Policy Targets" {
                      installontargets = append(installontargets, installonelements["name"].(string))
                        }
                      }
                natreturn["installon"] = client.ConvertListtoSet(installontargets)
                natlist = append(natlist, natreturn)
              }

            }

          }
  d.Set("natlist", natlist)
  return nil

}

func resourceAccessNatRuleListUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
packageuid := d.Get("package").(string)
//Get changes old and new and place them into a new interface
o, n := d.GetChange("natlist")
	if o == nil {
			o = new(schema.Set)
	}
		if n == nil {
			n = new(schema.Set)
		}
  _ = n
  // Pull in Rulebase rules
  natlist := d.Get("natlist").([]interface{})
  olist := o.([]interface{})
  olist_num := len(olist)
  nlist := n.([]interface{})
  nlist_num := len(nlist)
  // Create variable to hold values for rules that need to be updated
  uidlist := make([]interface{}, 0, len(natlist))
  if olist_num > nlist_num {
    //olist_num = olist_num - 1
    nlist_num = nlist_num - 1

    for w := range olist {
      deletenatrule := olist[w].(map[string]interface{})
      if w > nlist_num {
        deleteuid := deletenatrule["uid"].(string)
        client.DeleteAccessNATRule(deleteuid, packageuid)
      }
    }
    num_rules_delete := olist_num - nlist_num
    _ = num_rules_delete
  }
  for q := range natlist {
    var accessrulebase = chkp.AccessRulebaseNATListSet{}
    // Get new and old list of nat rules
    nelement := nlist[q].(map[string]interface{})
    if q < olist_num {
    oelement := olist[q].(map[string]interface{})
    _ = oelement
    _ = nelement
    //layerelement := layerlist[q].(map[string]interface{})

          // Starting checks to see if the rule needs updated.  If a update is needed we will set all aspects of the rule
          updateneeded := false
          if oelement["comments"] != nelement["comments"].(string) {
              updateneeded = true
            }

          if oelement["enabled"].(bool) != nelement["enabled"].(bool) {
              updateneeded = true
            }

          if oelement["method"].(string) != nelement["method"].(string) {
              updateneeded = true
            }

          if oelement["originalsource"].(string) != nelement["originalsource"].(string) {
              updateneeded = true
            }

          if oelement["originalservice"].(string) != nelement["originalservice"].(string) {
              updateneeded = true
            }
          if oelement["originaldestination"].(string) != nelement["originaldestination"].(string) {
              updateneeded = true
            }

          if oelement["translatedsource"].(string) != nelement["translatedsource"].(string) {
              updateneeded = true
            }

          if oelement["translatedservice"].(string) != nelement["translatedservice"].(string) {
              updateneeded = true
            }
          if oelement["translateddestination"].(string) != nelement["translateddestination"].(string) {
              updateneeded = true
            }
          installonnew := nelement["installon"].(*schema.Set).List()
          installonold := oelement["installon"].(*schema.Set).List()
          if !(reflect.DeepEqual(installonnew, installonold)) {
            updateneeded = true
          }

          if updateneeded == true  {
          accessrulebase.Uid = oelement["uid"].(string)
          accessrulebase.Comments = nelement["comments"].(string)
          accessrulebase.Enabled = nelement["enabled"].(bool)
          accessrulebase.Method = nelement["method"].(string)
          accessrulebase.OriginalSource = nelement["originalsource"].(string)
          accessrulebase.OriginalService = nelement["originalservice"].(string)
          accessrulebase.OriginalDestination = nelement["originaldestination"].(string)
          accessrulebase.TranslatedSource = nelement["translatedsource"].(string)
          accessrulebase.TranslatedService = nelement["translatedservice"].(string)
          accessrulebase.TranslatedDestination = nelement["translateddestination"].(string)
          accessrulebase.Package = packageuid
          installon := nelement["installon"].(*schema.Set).List()
          installonlist := make([]string, 0, len(installon))
          if len(installon) > 0 {
          for _, v := range installon {
            val, ok := v.(string)
            if ok && val != "" {
              installonlist = append(installonlist, v.(string))
            }
          }
        } else {
          installonlist = append(installonlist, "Policy Targets")
        }
          accessrulebase.Installon = installonlist
}

              if updateneeded == true {

                id, err := client.SetAccessNATRulebaseList(accessrulebase)
                //Read in from the API Output
                readNATRulebase := chkp.NATRulebaseResult{}
                json.Unmarshal(id, &readNATRulebase)

                natreturn := make(map[string]interface{})
                natreturn["uid"] = readNATRulebase.Uid
                natreturn["enabled"] = readNATRulebase.Enabled
                natreturn["install-on"] = readNATRulebase.Installon
                natreturn["original-destination"] = readNATRulebase.OriginalDestination
                natreturn["original-service"] = readNATRulebase.OriginalService
                natreturn["original-source"] = readNATRulebase.OriginalSource
                natreturn["translated-destination"] = readNATRulebase.TranslatedDestination
                natreturn["translated-service"] = readNATRulebase.TranslatedService
                natreturn["translated-source"] = readNATRulebase.TranslatedSource
                natreturn["comments"] = readNATRulebase.Comments
                uidlist = append(uidlist, natreturn)
                if err != nil {
              		return err
              	}

                    }
        }else {
          // If we have more rules in the list we will add them here and not use the api set command
        var natrulebase = chkp.AccessRulebaseNATList{}

          positionreturn := chkp.Position{
            Above: "Terraform NAT Rules Above This Position",
        }

        natrulebase.Position = positionreturn
        natrulebase.Package = packageuid
        natrulebase.Enabled = nelement["enabled"].(bool)
        // Get the list of Install-on Targets
        installon := nelement["installon"].(*schema.Set).List()
        installlist := make([]string, 0, len(installon))
        for _, v := range installon {
          val, ok := v.(string)
          if ok && val != "" {
            installlist = append(installlist, v.(string))
          }
        }

        if len(installlist) > 0 {
              natrulebase.Installon = installlist
            }

        natrulebase.OriginalDestination = nelement["originaldestination"].(string)
        natrulebase.OriginalService = nelement["originalservice"].(string)
        natrulebase.OriginalSource = nelement["originalsource"].(string)
        natrulebase.TranslatedDestination = nelement["translateddestination"].(string)
        natrulebase.TranslatedService = nelement["translatedservice"].(string)
        natrulebase.TranslatedSource = nelement["translatedsource"].(string)
        natrulebase.Comments = nelement["comments"].(string)
        natrulebase.Method = nelement["method"].(string)
        id, err := client.CreateAccessNATRulebaseList(natrulebase)
        //Read in from the API Output
        readNATRulebase := chkp.NATRulebaseResult{}
        json.Unmarshal(id, &readNATRulebase)

        natreturn := make(map[string]interface{})
        natreturn["uid"] = readNATRulebase.Uid
        natreturn["enabled"] = readNATRulebase.Enabled
        natreturn["install-on"] = readNATRulebase.Installon
        natreturn["original-destination"] = readNATRulebase.OriginalDestination
        natreturn["original-service"] = readNATRulebase.OriginalService
        natreturn["original-source"] = readNATRulebase.OriginalSource
        natreturn["translated-destination"] = readNATRulebase.TranslatedDestination
        natreturn["translated-service"] = readNATRulebase.TranslatedService
        natreturn["translated-source"] = readNATRulebase.TranslatedSource
        natreturn["comments"] = readNATRulebase.Comments
        uidlist = append(uidlist, natreturn)
        if err != nil {
          return err
        }

      }

}
d.Set("rulebase", uidlist)

  return nil
}

func resourceAccessNatRuleListDelete(d *schema.ResourceData, meta interface{}) error {
     client := meta.(*chkp.Client)
     client.DeletePolicyPackage(d.Id())

	return nil

}
