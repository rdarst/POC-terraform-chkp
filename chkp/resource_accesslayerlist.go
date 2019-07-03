package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
      //  "strings"
        "reflect"
)

func resourceAccessLayerList() *schema.Resource {
        return &schema.Resource{
                Create: resourceAccessLayerListCreate,
                Read:   resourceAccessLayerListRead,
                Update: resourceAccessLayerListUpdate,
                Delete: resourceAccessLayerListDelete,

                Schema: map[string]*schema.Schema{

                  "appandurl": {
                          Type:     schema.TypeBool,
                          Optional: true,
                            },
                  "contentawareness": {
                          Type:     schema.TypeBool,
                          Optional: true,
                          Default: false,
                  },
                  "adddefaultrule": {
                          Type:     schema.TypeBool,
                          Required: true,
                  },
                  "firewall": {
                          Type:     schema.TypeBool,
                          Optional: true,
                          Default: true,
                  },
                  "mobileaccess": {
                          Type:     schema.TypeBool,
                          Optional: true,
                          Default: false,
                  },
                  "name": {
                          Type:     schema.TypeString,
                          Required: true,
                  },
                  "color": {
                          Type:     schema.TypeString,
                          Optional: true,
                          Default: "black",
                  },
                  "comments": {
                          Type:     schema.TypeString,
                          Optional: true,
                  },
                  "shared": {
                          Type:     schema.TypeBool,
                          Optional: true,
                          Default: false,
                  },
                  "uid": {
                          Type:     schema.TypeString,
                          Optional: true,
                          Computed: true,
                  },

						            "rulebase": {
                          Type:     schema.TypeList,
                          Optional: true,
                          Elem: &schema.Resource{
                            Schema: map[string]*schema.Schema{
                              "name": {
                                      Type:     schema.TypeString,
                                      Optional: true,
                                        },
                             "action": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "inlinelayer": {
                                     Type:     schema.TypeString,
                                     Optional: true,
                                       },
                             "adddefaultrule": {
                                     Type:     schema.TypeBool,
                                     Optional: true,
                                       },
                             "source": {
                                     Type:     schema.TypeSet,
                                     Optional: true,
                                     Elem: &schema.Schema{ Type: schema.TypeString },
                                       },
                             "destination": {
                                     Type:     schema.TypeSet,
                                     Optional: true,
                                     Elem: &schema.Schema{ Type: schema.TypeString },
                                       },
                             "service": {
                                     Type:     schema.TypeSet,
                                     Optional: true,
                                     Elem: &schema.Schema{ Type: schema.TypeString },
                                       },
                             "track": {
                                     Type:     schema.TypeList,
                                     Required: true,
                                     MaxItems: 1,
                                     Elem: &schema.Resource{
                                       Schema: map[string]*schema.Schema{
                                           "type": {
                                                 Type:     schema.TypeString,
                                                 Optional: true,
                                                 Default: "Log",
                                           },
                                           "alert": {
                                                 Type:     schema.TypeString,
                                                 Optional: true,
                                                 Default: "none",
                                           },
                                           "accounting": {
                                                 Type:     schema.TypeBool,
                                                 Optional: true,
                                                 //Default:  false,
                                           },
                                           "perconnection": {
                                                 Type:     schema.TypeBool,
                                                 Optional: true,
                                                 Default:  true,
                                           },
                                           "persession": {
                                                 Type:     schema.TypeBool,
                                                 Optional: true,
                                           },
                                         },
                                       },
                                     },
                             "enabled": {
                                     Type:     schema.TypeBool,
                                     Optional: true,
                                     Default: true,
                                       },

                             "sourcenegate": {
                                     Type:     schema.TypeBool,
                                     Optional: true,
                                             },

                             "destinationnegate": {
                                     Type:     schema.TypeBool,
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

func resourceAccessLayerListCreate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)

  var accesslayer = chkp.AccessLayer{}
	accesslayer.Name = d.Get("name").(string)
	accesslayer.AppAndUrl = d.Get("appandurl").(bool)
  accesslayer.ContentAwareness = d.Get("contentawareness").(bool)
  accesslayer.Firewall = d.Get("firewall").(bool)
  accesslayer.MobileAccess = d.Get("mobileaccess").(bool)
  accesslayer.Shared = d.Get("shared").(bool)
  accesslayer.AddDefaultRule = d.Get("adddefaultrule").(bool)
  accesslayer.Color = d.Get("color").(string)
  accesslayer.Comments = d.Get("comments").(string)

	id, err := client.CreateAccessLayer(accesslayer)
  //Read in from the API Output
  readAccessLayer := chkp.AccessLayer{}
  json.Unmarshal(id, &readAccessLayer)
	d.SetId(readAccessLayer.Uid)
  d.Set("uid", readAccessLayer.Uid)
  d.Set("name", readAccessLayer.Name)
	if err != nil {
		return err
	}


  // Pull in Rulebase rules
  layerlist :=d.Get("rulebase").([]interface{})

  // Create variable to hold uids from returned rules
  uidlist := make([]interface{}, 0, len(layerlist))

  for q := range layerlist {
    layerelement := layerlist[q].(map[string]interface{})
          var accessrulebase = chkp.AccessRulebaseList{}
          accessrulebase.Name = layerelement["name"].(string)
          accessrulebase.Layer = readAccessLayer.Name
          accessrulebase.Action = layerelement["action"].(string)
          accessrulebase.Enabled = layerelement["enabled"].(bool)
          accessrulebase.InlineLayer = layerelement["inlinelayer"].(string)
          accessrulebase.DestinationNegate = layerelement["destinationnegate"].(bool)
          accessrulebase.SourceNegate = layerelement["sourcenegate"].(bool)
          accessrulebase.Position = q + 1

          tracklist := layerelement["track"].([]interface{})
          for i := range tracklist {
                trackelements := tracklist[i].(map[string]interface{})
                trackreturn := chkp.Track{
                  Type:       trackelements["type"].(string),
                  Alert:       trackelements["alert"].(string),
                  Accounting:   trackelements["accounting"].(bool),
                }
                accessrulebase.Track = trackreturn
            }

            source := layerelement["source"].(*schema.Set).List()
            sourcelist := make([]string, 0, len(source))
          	for _, v := range source {
          		val, ok := v.(string)
          		if ok && val != "" {
          			sourcelist = append(sourcelist, v.(string))
          		}
            }
            //Check to see if source list is empty.  If it is set the Source to Any
            if len(sourcelist) > 0 {
                  accessrulebase.Source = sourcelist
                }

            destination := layerelement["destination"].(*schema.Set).List()
            destinationlist := make([]string, 0, len(destination))
            for _, v := range destination {
              val, ok := v.(string)
              if ok && val != "" {
                destinationlist = append(destinationlist, v.(string))
              }
            }
            if len(destinationlist) > 0 {
            accessrulebase.Destination = destinationlist
               }
            service := layerelement["service"].(*schema.Set).List()
            servicelist := make([]string, 0, len(service))
          	for _, v := range service {
          		val, ok := v.(string)
          		if ok && val != "" {
          			servicelist = append(servicelist, v.(string))
          		}
            }
            if len(servicelist) > 0 {
            accessrulebase.Service = servicelist
              }

          id, err := client.CreateAccessRulebaseList(accessrulebase)
          //Read in from the API Output
          readAccessRulebase := chkp.AccessRulebaseResult{}
          json.Unmarshal(id, &readAccessRulebase)
        	//d.SetId(readAccessRulebase.Uid)
          //d.Set("uid", readAccessRulebase.Uid)
          layerreturn := make(map[string]interface{})
          layerreturn["uid"] = readAccessRulebase.Uid
          layerreturn["service"] = readAccessRulebase.Service
          layerreturn["source"] = readAccessRulebase.Source
          layerreturn["destination"] = readAccessRulebase.Destination
          layerreturn["track"] = readAccessRulebase.Track
          layerreturn["name"] = readAccessRulebase.Name
          uidlist = append(uidlist, layerreturn)
          if err != nil {
        		return err
        	}
}

d.Set("rulebase", uidlist)
return nil
}

func resourceAccessLayerListRead(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  id, err := client.ShowAccessLayer(d.Id())
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
  readAccessLayer := chkp.AccessLayer{}
  json.Unmarshal(id, &readAccessLayer)
	d.SetId(readAccessLayer.Uid)
	d.Set("color", readAccessLayer.Color)
	d.Set("name", readAccessLayer.Name)
	d.Set("appandurl", readAccessLayer.AppAndUrl)
  d.Set("contentawareness", readAccessLayer.ContentAwareness)
  d.Set("firewall", readAccessLayer.Firewall)
  d.Set("mobileaccess", readAccessLayer.MobileAccess)
  d.Set("shared", readAccessLayer.Shared)
  d.Set("comments", readAccessLayer.Comments)
  if err != nil {
		return err
	}

  // Set API limit on the number of records returned default is 50 max is 500
  limit := 100
  layer := readAccessLayer.Name
  id2, err := client.ShowAccessRulebaseList(layer, limit, 0)

  readAccessRulebase := chkp.AccessRulebaseResultRead{}
  json.Unmarshal(id2, &readAccessRulebase)

  var result map[string]interface{}
  json.Unmarshal(id2, &result)

  rules := result["rulebase"].([]interface{})
  total := readAccessRulebase.Total
  to := readAccessRulebase.To
  //from := readAccessRulebase.From
  offset := to

  rulelistread := make([]interface{}, 0, total)
  for t := range rules {
    rulelistread = append(rulelistread, rules[t].(map[string]interface{}))
  }
  // Check to see if we need to gather more rules
  if total > to {
  done := 0
  for done < 1 {
    id3, err := client.ShowAccessRulebaseList(layer, limit, offset)
    var resultappend map[string]interface{}
    readAccessRulebaseappend := chkp.AccessRulebaseResultRead{}
    json.Unmarshal(id3, &resultappend)
    json.Unmarshal(id3, &readAccessRulebaseappend)
    rules := resultappend["rulebase"].([]interface{})
    for t := range rules {
      rulelistread = append(rulelistread, rules[t].(map[string]interface{}))
    }
    if err != nil {
  		return err
  	}
    total = readAccessRulebaseappend.Total
    to = readAccessRulebaseappend.To
    offset = to

    if total == to { done = 2}
    //done = 2
      }
    }

  // build the interface for the rule list
  rulelist := make([]interface{}, 0, len(rules))

  for q := range rulelistread {
    ruleelement := rulelistread[q].(map[string]interface{})

          layerreturn := make(map[string]interface{})
          layerreturn["uid"] = ruleelement["uid"]
          // If name is empty don't populate array for name
          if ruleelement["name"] != nil {
            layerreturn["name"] = ruleelement["name"].(string)
          }
        //  layerreturn["inlinelayer"] = ruleelement["inline-layer"].(string)
          layerreturn["enabled"] = ruleelement["enabled"].(bool)
          layerreturn["sourcenegate"] = ruleelement["source-negate"].(bool)
          layerreturn["destinationnegate"] = ruleelement["destination-negate"].(bool)

          //Gather Action settings
          action := ruleelement["action"]
          actionlist := action.(map[string]interface{})
          actionresults  := actionlist["name"]
          if actionresults == "Inner Layer" {
              actionresults = "Apply Layer"
              inlinelayerresult := ruleelement["inline-layer"]
              inlinelayerresultlist := inlinelayerresult.(map[string]interface{})
              layerreturn["inlinelayer"] = inlinelayerresultlist["name"]
          }
          //layerreturn["action"] = strings.ToLower(actionresults.(string))
          layerreturn["action"] = actionresults.(string)

          //Gather Sources
          source := ruleelement["source"]
          sourcelist := source.([]interface{})
          sources := make([]string, 0, len(sourcelist))
          for i := range sourcelist {
                sourceelements := sourcelist[i].(map[string]interface{})
                if sourceelements["name"] != "Any" {
                sources = append(sources, sourceelements["name"].(string))
                  }
                }
          layerreturn["source"] = client.ConvertListtoSet(sources)

          //Gather destination
          destination := ruleelement["destination"]
          destinationlist := destination.([]interface{})
          destinations := make([]string, 0, len(destinationlist))
          for i := range destinationlist {
                destinationelements := destinationlist[i].(map[string]interface{})
                if destinationelements["name"] != "Any" {
                    destinations = append(destinations, destinationelements["name"].(string))
                  }
                }
          layerreturn["destination"] = client.ConvertListtoSet(destinations)

          //Gather service
          service := ruleelement["service"]
          servicelist := service.([]interface{})
          services := make([]string, 0, len(servicelist))
          for i := range servicelist {
                serviceelements := servicelist[i].(map[string]interface{})
                if serviceelements["name"] != "Any" {
                    services = append(services, serviceelements["name"].(string))
                  }
                }
          layerreturn["service"] = client.ConvertListtoSet(services)

          // Tracklist
          track := ruleelement["track"]
          tracklist := track.(map[string]interface{})
          tracks := &chkp.Track{
                    Type:       tracklist["type"].(map[string]interface{})["name"].(string),
                    Accounting:  tracklist["accounting"].(bool),
                    Alert:  tracklist["alert"].(string),
                    PerConnection:  tracklist["per-connection"].(bool),
                    PerSession:  tracklist["per-session"].(bool),
                }
          layerreturn["track"] = flattenTrackSettings(tracks)
          rulelist = append(rulelist, layerreturn)
  }

  d.Set("rulebase", rulelist)

  return nil

}



func resourceAccessLayerListUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  updatepolicy := false
  var accesslayer = chkp.AccessLayerUpdate{}
  oldname, newname := d.GetChange("name")
  if oldname.(string) != newname.(string) {
    accesslayer.Newname = newname.(string)
    accesslayer.Name = oldname.(string)
    updatepolicy = true
   } else {
     accesslayer.Name = d.Get("name").(string)
   }
  if d.HasChange("appandurl") {
   accesslayer.AppAndUrl = d.Get("appandurl").(bool)
   updatepolicy = true
  }
  if d.HasChange("contentawareness") {
    accesslayer.ContentAwareness = d.Get("contentawareness").(bool)
    updatepolicy = true
   }
   if d.HasChange("firewall") {
     accesslayer.Firewall = d.Get("firewall").(bool)
     updatepolicy = true
    }
    if d.HasChange("mobileaccess") {
      accesslayer.MobileAccess = d.Get("mobileaccess").(bool)
      updatepolicy = true
     }
   if d.HasChange("shared") {
     accesslayer.Shared = d.Get("shared").(bool)
     updatepolicy = true
    }
    if d.HasChange("color") {
      accesslayer.Color = d.Get("color").(string)
      updatepolicy = true
     }
   if d.HasChange("comments") {
    accesslayer.Comments = d.Get("comments").(string)
     updatepolicy = true
    }

	if updatepolicy  {
    id, err := client.SetAccessLayer(accesslayer)
    if err != nil {return err}
    //Read in from the API Output
    readAccessLayer := chkp.AccessLayer{}
    if err := json.Unmarshal(id, &readAccessLayer); err != nil {return err}
    d.SetId(readAccessLayer.Uid)
    d.Set("uid", readAccessLayer.Uid)
    d.Set("name", readAccessLayer.Name)
  	if err != nil {
  		return err
  	}
  }
  layer := d.Get("name").(string)

//Get changes old and new and place them into a new interface
o, n := d.GetChange("rulebase")
	if o == nil {
			o = new(schema.Set)
	}
		if n == nil {
			n = new(schema.Set)
		}
  _ = n
  // Pull in Rulebase rules
  layerlist := d.Get("rulebase").([]interface{})
  olist := o.([]interface{})
  olist_num := len(olist)
  nlist := n.([]interface{})
  nlist_num := len(nlist)
  // Create variable to hold values for rules that need to be updated
  uidlist := make([]interface{}, 0, len(layerlist))
  if olist_num > nlist_num {
    num_rules_delete := olist_num - nlist_num
    //for w := olist_num; w = nlist_num; w-- {
    for w := num_rules_delete; w > 0; w-- {
          position := olist_num
          olist_num = olist_num - 1
      client.DeleteAccessRuleByRuleNum(position, layer)

    }
  }
  for q := range layerlist {
    var accessrulebase = chkp.AccessRulebaseList{}
    //if q <= olist_num { oelement := olist[q].(map[string]interface{})}
    nelement := nlist[q].(map[string]interface{})
    if q < olist_num {
    oelement := olist[q].(map[string]interface{})
    _ = oelement
    _ = nelement
    //layerelement := layerlist[q].(map[string]interface{})

          // Starting checks to see if the rule needs updated.  If a update is needed we will set all aspects of the rule
          updateneeded := false
          if oelement["name"] != nelement["name"].(string) {
              updateneeded = true
            }
          if oelement["action"] != nelement["action"].(string) {
              updateneeded = true
            }
          if oelement["enabled"].(bool) != nelement["enabled"].(bool) {
              updateneeded = true
            }
          if oelement["destinationnegate"].(bool) != nelement["destinationnegate"].(bool) {
              updateneeded = true
            }
          if oelement["sourcenegate"].(bool) != nelement["sourcenegate"].(bool) {
              updateneeded = true
            }
          if oelement["inlinelayer"].(string) != nelement["inlinelayer"].(string) {
              updateneeded = true
            }
            tracklistnew := nelement["track"].([]interface{})
            tracklistold := oelement["track"].([]interface{})
          if !(reflect.DeepEqual(tracklistnew, tracklistold)) {
            updateneeded = true
          }
          sourcenew := nelement["source"].(*schema.Set).List()
          sourceold := oelement["source"].(*schema.Set).List()
          if !(reflect.DeepEqual(sourcenew, sourceold)) {
            updateneeded = true
          }
          destinationnew := nelement["destination"].(*schema.Set).List()
          destinationold := oelement["destination"].(*schema.Set).List()
          if !(reflect.DeepEqual(destinationnew, destinationold)) {
            updateneeded = true
          }
          servicenew := nelement["service"].(*schema.Set).List()
          serviceold := oelement["service"].(*schema.Set).List()
          if !(reflect.DeepEqual(servicenew, serviceold)) {
            updateneeded = true
          }

          if updateneeded == true  {
          accessrulebase.Rulenumber = q + 1
          accessrulebase.Layer = layer
          accessrulebase.Newname = nelement["name"].(string)
          accessrulebase.Action = nelement["action"].(string)
          accessrulebase.Enabled = nelement["enabled"].(bool)
          accessrulebase.DestinationNegate = nelement["destinationnegate"].(bool)
          accessrulebase.SourceNegate = nelement["sourcenegate"].(bool)
          accessrulebase.InlineLayer = nelement["inlinelayer"].(string)
          tracklist := nelement["track"].([]interface{})
          for i := range tracklist {
                trackelements := tracklist[i].(map[string]interface{})
                trackreturn := chkp.Track{
                  Type:       trackelements["type"].(string),
                  Alert:       trackelements["alert"].(string),
                  Accounting:   trackelements["accounting"].(bool),
                }
                accessrulebase.Track = trackreturn
            }
            source := nelement["source"].(*schema.Set).List()
            sourcelist := make([]string, 0, len(source))
            if len(source) > 0 {
            for _, v := range source {
              val, ok := v.(string)
              if ok && val != "" {
                sourcelist = append(sourcelist, v.(string))
              }
            }
          } else {
            sourcelist = append(sourcelist, "Any")
          }

            accessrulebase.Source = sourcelist

            destination := nelement["destination"].(*schema.Set).List()
            destinationlist := make([]string, 0, len(destination))
            if len(destination) > 0 {
            for _, v := range destination {
              val, ok := v.(string)
              if ok && val != "" {
                destinationlist = append(destinationlist, v.(string))
              }
            }
          } else {
               destinationlist = append(destinationlist, "Any")
          }

            accessrulebase.Destination = destinationlist

            service := nelement["service"].(*schema.Set).List()
            servicelist := make([]string, 0, len(service))
            if len(service) > 0 {
            for _, v := range service {
              val, ok := v.(string)
              if ok && val != "" {
                servicelist = append(servicelist, v.(string))
              }
            }
          } else {
            servicelist = append(servicelist, "Any")
          }

            accessrulebase.Service = servicelist
}

              if updateneeded == true {

                idset, err := client.SetAccessRulebaseList(accessrulebase)
                //Read in from the API Output
                readAccessRulebase := chkp.AccessRulebaseResult{}
                json.Unmarshal(idset, &readAccessRulebase)

                layerreturn := make(map[string]interface{})
                layerreturn["uid"] = readAccessRulebase.Uid
                layerreturn["service"] = readAccessRulebase.Service
                layerreturn["source"] = readAccessRulebase.Source
                layerreturn["destination"] = readAccessRulebase.Destination
                layerreturn["track"] = readAccessRulebase.Track
                layerreturn["name"] = readAccessRulebase.Name
                uidlist = append(uidlist, layerreturn)
                if err != nil {
                  return err
                }

                    }
        }else {
          // If we have more rules in the list we will add them here and not use the api set command
          accessrulebase.Position = q + 1
          accessrulebase.Layer = layer
          accessrulebase.Name = nelement["name"].(string)
          accessrulebase.Action = nelement["action"].(string)
          accessrulebase.Enabled = nelement["enabled"].(bool)
          accessrulebase.DestinationNegate = nelement["destinationnegate"].(bool)
          accessrulebase.SourceNegate = nelement["sourcenegate"].(bool)
          tracklist := nelement["track"].([]interface{})
          for i := range tracklist {
                trackelements := tracklist[i].(map[string]interface{})
                trackreturn := chkp.Track{
                  Type:       trackelements["type"].(string),
                  Alert:       trackelements["alert"].(string),
                  Accounting:   trackelements["accounting"].(bool),
                }
                accessrulebase.Track = trackreturn
            }
            source := nelement["source"].(*schema.Set).List()
            sourcelist := make([]string, 0, len(source))
            if len(source) > 0 {
            for _, v := range source {
              val, ok := v.(string)
              if ok && val != "" {
                sourcelist = append(sourcelist, v.(string))
              }
            }
          } else {
            sourcelist = append(sourcelist, "Any")
          }

            accessrulebase.Source = sourcelist

            destination := nelement["destination"].(*schema.Set).List()
            destinationlist := make([]string, 0, len(destination))
            if len(destination) > 0 {
            for _, v := range destination {
              val, ok := v.(string)
              if ok && val != "" {
                destinationlist = append(destinationlist, v.(string))
              }
            }
          } else {
               destinationlist = append(destinationlist, "Any")
          }

            accessrulebase.Destination = destinationlist

            service := nelement["service"].(*schema.Set).List()
            servicelist := make([]string, 0, len(service))
            if len(service) > 0 {
            for _, v := range service {
              val, ok := v.(string)
              if ok && val != "" {
                servicelist = append(servicelist, v.(string))
              }
            }
          } else {
            servicelist = append(servicelist, "Any")
          }

            accessrulebase.Service = servicelist

              idadd, err := client.CreateAccessRulebaseList(accessrulebase)
              //Read in from the API Output
              readAccessRulebase := chkp.AccessRulebaseResult{}
              json.Unmarshal(idadd, &readAccessRulebase)
              layerreturn := make(map[string]interface{})
              layerreturn["uid"] = readAccessRulebase.Uid
              layerreturn["service"] = readAccessRulebase.Service
              layerreturn["source"] = readAccessRulebase.Source
              layerreturn["destination"] = readAccessRulebase.Destination
              layerreturn["track"] = readAccessRulebase.Track
              layerreturn["name"] = readAccessRulebase.Name
              uidlist = append(uidlist, layerreturn)
              if err != nil {
                return err
            }
        }
      }

d.Set("rulebase", uidlist)

  return nil

}

  func resourceAccessLayerListDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
    client.DeleteAccessLayer(d.Id())
    return nil

  }
