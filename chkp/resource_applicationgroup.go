package chkp

import (
        "github.com/hashicorp/terraform/helper/schema"
        "./client"
    		"encoding/json"
)

func resourceApplicationGroup() *schema.Resource {
        return &schema.Resource{
                Create: resourceApplicationGroupCreate,
                Read:   resourceApplicationGroupRead,
                Update: resourceApplicationGroupUpdate,
                Delete: resourceApplicationGroupDelete,

                Schema: map[string]*schema.Schema{

						            "name": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
                        "color": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Default: "black",
                        },
                        "members": &schema.Schema{
                                Type:     schema.TypeSet,
                                Optional: true,
                                Elem: &schema.Schema{Type: schema.TypeString},
                        },
                        "uid": {
                                Type:     schema.TypeString,
                                Optional: true,
                                Computed: true,
                            },

                },
        }
}

func resourceApplicationGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chkp.Client)
	var applicationgroup = chkp.ApplicationGroup{}
	applicationgroup.Name = d.Get("name").(string)
  applicationgroup.Color = d.Get("color").(string)

  // Pull in the list of ApplicationGroup Members
  members := d.Get("members").(*schema.Set).List()
  memberlist := make([]string, 0, len(members))
	for _, v := range members {
		val, ok := v.(string)
		if ok && val != "" {
			memberlist = append(memberlist, v.(string))
		}
  }
  applicationgroup.Members = memberlist

  id, err := client.CreateApplicationGroup(applicationgroup)
  readApplicationGroup := chkp.ApplicationGroup{}
  json.Unmarshal(id, &readApplicationGroup)
	d.SetId(readApplicationGroup.Uid)
  d.Set("uid", readApplicationGroup.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceApplicationGroupRead(d *schema.ResourceData, meta interface{}) error {
        client := meta.(*chkp.Client)
  // Call the API to get ApplicationGroup info
  id, err := client.ShowApplicationGroup(d.Id())
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
	readApplicationGroup := chkp.ApplicationGroup{}
  json.Unmarshal(id, &readApplicationGroup)
	d.SetId(readApplicationGroup.Uid)
	d.Set("color", readApplicationGroup.Color)
	d.Set("name", readApplicationGroup.Name)
  var applicationgroup chkp.ApplicationGroupMembers
  json.Unmarshal(id, &applicationgroup)

  // Read the list of Current ApplicationGroup Members
  i := 0
  mem := make([]string, 0, len(applicationgroup.Members))
  for _, member := range readApplicationGroup.Members {
    mem = append(mem, applicationgroup.Members[i].Name)
    _ = member
    i += 1
  }

  d.Set("members", mem)

	return nil
}

func resourceApplicationGroupUpdate(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
	var applicationgroup = chkp.ApplicationGroup{}
	applicationgroup.Color = d.Get("color").(string)
  if d.HasChange("name") {applicationgroup.Newname = d.Get("name").(string)}
	applicationgroup.Uid = d.Get("uid").(string)
  // Pull in the list of ApplicationGroup Members
  members := d.Get("members").(*schema.Set).List()
  memberlist := make([]string, 0, len(members))
  for _, v := range members {
    val, ok := v.(string)
    if ok && val != "" {
      memberlist = append(memberlist, v.(string))
    }
  }
  applicationgroup.Members = memberlist

	id, err := client.SetApplicationGroup(applicationgroup)
  readApplicationGroup := chkp.ApplicationGroupMembers{}
  if err := json.Unmarshal(id, &readApplicationGroup); err != nil {return err}
  //json.Unmarshal(id, &readApplicationGroup)
	d.SetId(readApplicationGroup.Uid)
  d.Set("uid", readApplicationGroup.Uid)
	if err != nil {
		return err
	}
	return nil
}

func resourceApplicationGroupDelete(d *schema.ResourceData, meta interface{}) error {
  client := meta.(*chkp.Client)
  used, err := client.CheckWhereUsed(d.Id())
  if used > 0 {
    client.WaitUntilNotUsed(d.Id())
  }
	result, err := client.DeleteApplicationGroup(d.Id())
  _ = result
  if err != nil {
		return err
	}
  return nil
}
