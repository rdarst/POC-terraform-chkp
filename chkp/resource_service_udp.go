package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
        "strings"
)


func resourceServiceUdp() *schema.Resource {
        return &schema.Resource{
                Create: resourceServiceUdpCreate,
                Read:   resourceServiceUdpRead,
                Update: resourceServiceUdpUpdate,
                Delete: resourceServiceUdpDelete,

                Schema: map[string]*schema.Schema{
							        	"port": {
                                Type:     schema.TypeString,
                                Required: true,
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
                        "protocol": {
                                Type:     schema.TypeString,
                                Optional: true,
                        },
                        "sessiontimeout": {
                                Type:     schema.TypeInt,
                                Optional: true,
                                Default: 3600,
                        },
                        "matchbysig": {
                                Type:     schema.TypeBool,
                                Optional: true,
                        },
                        "matchforany": {
                                Type:     schema.TypeBool,
                                Optional: true,
                                Default: true,
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                        },
                },
        }
}

func resourceServiceUdpCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var serviceudp = chkp.ServiceUdp{}
	serviceudp.Name = d.Get("name").(string)
	serviceudp.Port = d.Get("port").(string)
  serviceudp.Color = d.Get("color").(string)
  serviceudp.Protocol = strings.ToUpper(d.Get("protocol").(string))
  serviceudp.SessionTimeout = d.Get("sessiontimeout").(int)
  serviceudp.MatchBySig = d.Get("matchbysig").(bool)
  serviceudp.MatchForAny = d.Get("matchforany").(bool)
	id, err := client.CreateServiceUdp(serviceudp)
  //Read in from the API Output
  readServiceUdp := chkp.ServiceUdp{}
  json.Unmarshal(id, &readServiceUdp)
	d.SetId(readServiceUdp.Uid)
  d.Set("uid", readServiceUdp.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceServiceUdpRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
	id, err := client.ShowServiceUdp(d.Id())

	readServiceUdp := chkp.ServiceUdp{}
  json.Unmarshal(id, &readServiceUdp)
	d.SetId(readServiceUdp.Uid)
	d.Set("color", readServiceUdp.Color)
	d.Set("name", readServiceUdp.Name)
	d.Set("port", readServiceUdp.Port)
  d.Set("protocol", strings.ToUpper(readServiceUdp.Protocol))
  d.Set("sessiontimeout", readServiceUdp.SessionTimeout)
  d.Set("matchbysig", readServiceUdp.MatchBySig)
  d.Set("matchforany", readServiceUdp.MatchForAny)
	if err != nil {
		return err
	}
	return nil
}

func resourceServiceUdpUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var serviceudp = chkp.ServiceUdp{}
	serviceudp.Port = d.Get("port").(string)
	serviceudp.Color = d.Get("color").(string)
  serviceudp.Protocol = d.Get("protocol").(string)
  serviceudp.SessionTimeout = d.Get("sessiontimeout").(int)
  serviceudp.MatchBySig = d.Get("matchbysig").(bool)
  serviceudp.MatchForAny = d.Get("matchforany").(bool)
	if d.HasChange("name") {serviceudp.Newname = d.Get("name").(string)}
	serviceudp.Uid = d.Get("uid").(string)
	id, err := client.SetServiceUdp(serviceudp)
  if err != nil {return err}
  //Read in from the API Output
  readServiceUdp := chkp.ServiceUdp{}
  if err := json.Unmarshal(id, &readServiceUdp); err != nil {return err}
  d.SetId(readServiceUdp.Uid)
  d.Set("uid", readServiceUdp.Uid)

	if err != nil {
		return err
	}
	return nil
}

func resourceServiceUdpDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*chkp.Client)
	client.DeleteServiceUdp(d.Id())
	return nil

}
