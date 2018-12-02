package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)

func resourceAccessRulebase() *schema.Resource {
        return &schema.Resource{
                Create: resourceAccessRulebaseCreate,
                Read:   resourceAccessRulebaseRead,
                Update: resourceAccessRulebaseUpdate,
                Delete: resourceAccessRulebaseDelete,

                Schema: map[string]*schema.Schema{

						            "name": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "position": {
                                Type:     schema.TypeList,
                                Optional: true,
                                Elem: &schema.Resource{
                                  Schema: map[string]*schema.Schema{
                                      "top": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "bottom": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "above": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "below": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                    },
                                  },
                                },
						            "layer": {
                                Type:     schema.TypeString,
                                Required: true,
                                  },
                        "action": {
                                Type:     schema.TypeString,
                                Optional: true,
                                  },
                        "inlinelayer": {
                                Type:     schema.TypeString,
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
                                Elem: &schema.Resource{
                                  Schema: map[string]*schema.Schema{
                                      "type": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                            Default: "log",
                                      },
                                      "alert": {
                                            Type:     schema.TypeString,
                                            Optional: true,
                                      },
                                      "accounting": {
                                            Type:     schema.TypeBool,
                                            Optional: true,
                                      },
                                      "perconnection": {
                                            Type:     schema.TypeBool,
                                            Optional: true,
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
        }
}

func resourceAccessRulebaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var accessrulebase = chkp.AccessRulebase{}
	accessrulebase.Name = d.Get("name").(string)
  accessrulebase.Layer = d.Get("layer").(string)
  accessrulebase.Action = d.Get("action").(string)
  accessrulebase.Enabled = d.Get("enabled").(bool)
  accessrulebase.InlineLayer = d.Get("inlinelayer").(string)
  accessrulebase.DestinationNegate = d.Get("destinationnegate").(bool)
  accessrulebase.SourceNegate = d.Get("sourcenegate").(bool)

  positionlist := d.Get("position").([]interface{})
  for i := range positionlist {
        positionelements := positionlist[i].(map[string]interface{})
        positionreturn := chkp.Position{
          Above:       positionelements["above"].(string),
          Below:       positionelements["below"].(string),
          Top:         positionelements["top"].(string),
          Bottom:      positionelements["bottom"].(string),
        }
        accessrulebase.Position = positionreturn
    }

  tracklist := d.Get("track").([]interface{})
  for i := range tracklist {
        trackelements := tracklist[i].(map[string]interface{})
        trackreturn := chkp.Track{
          Type:       trackelements["type"].(string),
          Alert:       trackelements["alert"].(string),
          Accounting:   trackelements["accounting"].(bool),
        }
        accessrulebase.Track = trackreturn
    }

    source := d.Get("source").(*schema.Set).List()
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

    destination := d.Get("destination").(*schema.Set).List()
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
    service := d.Get("service").(*schema.Set).List()
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

  id, err := client.CreateAccessRulebase(accessrulebase)
  //Read in from the API Output
  readAccessRulebase := chkp.AccessRulebaseResult{}
  json.Unmarshal(id, &readAccessRulebase)
	d.SetId(readAccessRulebase.Uid)
  d.Set("uid", readAccessRulebase.Uid)

  i :=0
  serviceread := make([]string, 0, len(readAccessRulebase.Service))
  for _, test := range readAccessRulebase.Service {
       serviceread = append(serviceread, readAccessRulebase.Service[i].Uid)
       _ = test
       i +=1
    }

  d.Set("service", serviceread)

	if err != nil {
		return err
	}
	return nil
}

func resourceAccessRulebaseRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
  layer := d.Get("layer").(string)
  id, err := client.ShowAccessRulebase(d.Id(),layer)

	readAccessRulebase := chkp.AccessRulebaseResult{}
  json.Unmarshal(id, &readAccessRulebase)
	d.SetId(readAccessRulebase.Uid)
	d.Set("uid", readAccessRulebase.Uid)
  if len(readAccessRulebase.Layer) != 0 {
    layername, err := client.ReadLayerUIDtoName(readAccessRulebase.Layer)
    if err == nil {
      d.Set("layer", layername)
    }
  }

  d.Set("enabled", readAccessRulebase.Enabled)
  d.Set("name", readAccessRulebase.Name)
  d.Set("sourcenegate", readAccessRulebase.SourceNegate)
  d.Set("destinationnegate", readAccessRulebase.DestinationNegate)

  i := 0
  source := make([]string, 0, len(readAccessRulebase.Source))
  for _, test := range readAccessRulebase.Source {
    // Handle cases where we have any in the rule - don't provide it in the string
    if readAccessRulebase.Source[i].Name != "Any" {
        source = append(source, readAccessRulebase.Source[i].Uid)
        }
       _ = test
       i += 1
    }

  d.Set("source", source)

  i = 0
  destination := make([]string, 0, len(readAccessRulebase.Destination))
  for _, test := range readAccessRulebase.Destination {
    // Handle cases where we have any in the rule - don't provide it in the string
    if readAccessRulebase.Destination[i].Name != "Any" {
       destination = append(destination, readAccessRulebase.Destination[i].Uid)
     }
       _ = test
       i += 1
    }

  d.Set("destination", destination)

  i = 0
  service := make([]string, 0, len(readAccessRulebase.Service))
  for _, test := range readAccessRulebase.Service {
    // Handle cases where we have any in the rule - don't provide it in the string
    if readAccessRulebase.Service[i].Name != "Any" {
       service = append(service, readAccessRulebase.Service[i].Uid)
     }
       _ = test
       i += 1
    }

  d.Set("service", service)

  i = 0
  inlinelayer := make([]string, 0, len(readAccessRulebase.InlineLayer))
  for _, test := range readAccessRulebase.InlineLayer {
       inlinelayer = append(inlinelayer, readAccessRulebase.InlineLayer[i].Name)
       _ = test
       i += 1
    }

  d.Set("inlinelayer", inlinelayer)


	if err != nil {
		return err
	}
	return nil
}

func resourceAccessRulebaseUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var accessrulebase = chkp.AccessRulebase{}
  accessrulebase.Layer = d.Get("layer").(string)
  accessrulebase.Enabled = d.Get("enabled").(bool)
  accessrulebase.Action = d.Get("action").(string)
  accessrulebase.InlineLayer = d.Get("inlinelayer").(string)
  accessrulebase.DestinationNegate = d.Get("destinationnegate").(bool)
  accessrulebase.SourceNegate = d.Get("sourcenegate").(bool)
  //Update name if it has changed
	if d.HasChange("name") {accessrulebase.Newname = d.Get("name").(string)}

  positionlist := d.Get("position").([]interface{})

  for i := range positionlist {
        positionelements := positionlist[i].(map[string]interface{})
        positionreturn := chkp.Position{
          Above:       positionelements["above"].(string),
          Below:       positionelements["below"].(string),
          Top:         positionelements["top"].(string),
          Bottom:      positionelements["bottom"].(string),
        }
        accessrulebase.Position = positionreturn
    }

  tracklist := d.Get("track").([]interface{})

  for i := range tracklist {
        trackelements := tracklist[i].(map[string]interface{})
        trackreturn := chkp.Track{
          Type:       trackelements["type"].(string),
          Alert:       trackelements["alert"].(string),
          Accounting:   trackelements["accounting"].(bool),
        }
        accessrulebase.Track = trackreturn
    }

    source := d.Get("source").(*schema.Set).List()
    sourcelist := make([]string, 0, len(source))
    for _, v := range source {
      val, ok := v.(string)
      if ok && val != "" {
        sourcelist = append(sourcelist, v.(string))
      }
    }
    //Check to see if source list is empty.  Leave empty to set to Any
    if len(sourcelist) > 0 {
          accessrulebase.Source = sourcelist
        } else {
          any := []string{"Any"}
          accessrulebase.Source = any
        }

    destination := d.Get("destination").(*schema.Set).List()
    destinationlist := make([]string, 0, len(destination))
    for _, v := range destination {
      val, ok := v.(string)
      if ok && val != "" {
        destinationlist = append(destinationlist, v.(string))
      }
    }
    //Check to see if destination list is empty.  Leave empty to set to Any
    if len(destinationlist) > 0 {
    accessrulebase.Destination = destinationlist
      } else {
        any := []string{"Any"}
        accessrulebase.Destination = any
      }

    service := d.Get("service").(*schema.Set).List()
    servicelist := make([]string, 0, len(service))
    for _, v := range service {
      val, ok := v.(string)
      if ok && val != "" {
        servicelist = append(servicelist, v.(string))
      }
    }
    //Check to see if service list is empty.  Leave empty to set to Any
    if len(servicelist) > 0 {
    accessrulebase.Service = servicelist
      } else {
        any := []string{"Any"}
        accessrulebase.Service = any
      }

	accessrulebase.Uid = d.Get("uid").(string)

  id,  err := client.SetAccessRulebase(accessrulebase)
  if err != nil {return err}
  _ = id
  //Read in from the API Output
  //readAccessRulebase := chkp.AccessRulebase{}
  //if err := json.Unmarshal(id, &readAccessRulebase); err != nil {return err}
  //d.SetId(readAccessRulebase.Uid)
  //d.Set("uid", readAccessRulebase.Uid)

	//if err != nil {
	//	return err
	//}
	return nil
}

func resourceAccessRulebaseDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
    layer := d.Get("layer").(string)
    uid := d.Get("uid").(string)
	client.DeleteAccessRulebase(uid, layer)
	return nil

}
