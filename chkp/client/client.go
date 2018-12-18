package chkp

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
  "crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

type Client struct {
	client *http.Client
	base   string
	sid  string
}

func NewClientWith(server string, sid string) (*Client, error) {
      http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	  var netClient = &http.Client{
	    Timeout: time.Second * 10,
	}
	//login := Login{
	//	User: user,
	//	Password:     password,
	//	Domain: "CMA1",
	//	}
	//loginBytes, err := json.Marshal(login)
	//if err != nil {
	//	return nil, err
	//}
	//loginReader := bytes.NewReader(loginBytes)

	//response, err := netClient.Post(server+"/login", "application/json", loginReader)
	//if response.StatusCode == 400 {
		//return nil, errors.New("Sorry but the log in info is incorrect")
	//}
	//if err != nil {
	//	return nil, err
	//}
	//defer response.Body.Close()
	//contents, err := ioutil.ReadAll(response.Body)
	//if err != nil {
		//return nil, err
	//}
	//loginResponse := LoginResponse{}
	//err = json.Unmarshal(contents, &loginResponse)
	//if err != nil {
	//	return nil, err
	//}
	return &Client{netClient, server, sid}, nil
}

//func(c *Client) GetNewSession(user string, password string, server string) (string, error) {
//      http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
//	    var netClient = &http.Client{
//	    Timeout: time.Second * 10,
//	}
//	login := Login{
//		User: user,
//		Password:     password,
//	//	Domain: "CMA1",
//		}
//	loginBytes, err := json.Marshal(login)
//	if err != nil {
//		return "", err
//	}
//	loginReader := bytes.NewReader(loginBytes)
//  response, err := netClient.Post(server+"/login", "application/json", loginReader)
//	if response.StatusCode == 400 {
//	  return "", errors.New("Sorry but the log in info is incorrect")
//	}
//	if err != nil {
//		return "", err
//	}
//	defer response.Body.Close()
//	contents, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//	  return "", err
//	}
//	loginResponse := LoginResponse{}
//	err = json.Unmarshal(contents, &loginResponse)
//  if err != nil {
//		return "", err
//	}
//  sid := loginResponse.Sid
//	return sid, err
//}

func(c *Client) CreateHost(host Host) ([]byte, error) {

	spotBytes, _ := json.Marshal(host)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-host", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetHost(host Host) ([]byte, error) {

	spotBytes, _ := json.Marshal(host)
	spotReader := bytes.NewReader(spotBytes)
	req, err := http.NewRequest("POST", c.base+"/set-host", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowHost(hostuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, hostuid))
	req, err := http.NewRequest("POST", c.base+"/show-host", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)
	return body, err
}

func(c *Client) DeleteHost(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-host", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) CreateAddressRange(addressrange AddressRange) ([]byte, error) {

	spotBytes, _ := json.Marshal(addressrange)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-address-range", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetAddressRange(addressrange AddressRange) ([]byte, error) {

	spotBytes, _ := json.Marshal(addressrange)
	spotReader := bytes.NewReader(spotBytes)
	req, err := http.NewRequest("POST", c.base+"/set-address-range", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowAddressRange(addressrangeuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, addressrangeuid))
	req, err := http.NewRequest("POST", c.base+"/show-address-range", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteAddressRange(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-address-range", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) CreateGroup(group Group) ([]byte, error) {

	spotBytes, _ := json.Marshal(group)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-group", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetGroup(group Group) ([]byte, error) {

	spotBytes, _ := json.Marshal(group)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-group", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowGroup(groupuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, groupuid))
	req, err := http.NewRequest("POST", c.base+"/show-group", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteGroup(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-group", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}
////////////////////
func(c *Client) CreateApplicationGroup(applicationgroup ApplicationGroup) ([]byte, error) {

	spotBytes, _ := json.Marshal(applicationgroup)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-application-site-group", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetApplicationGroup(applicationgroup ApplicationGroup) ([]byte, error) {

	spotBytes, _ := json.Marshal(applicationgroup)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-application-site-group", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowApplicationGroup(applicationgroupuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, applicationgroupuid))
	req, err := http.NewRequest("POST", c.base+"/show-application-site-group", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteApplicationGroup(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-application-site-group", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}
////////////////////
func(c *Client) CreateApplicationSite(group ApplicationSite) ([]byte, error) {

	spotBytes, _ := json.Marshal(group)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-application-site", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetApplicationSite(group ApplicationSite) ([]byte, error) {

	spotBytes, _ := json.Marshal(group)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-application-site", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetApplicationSiteUpdateTag(group ApplicationSiteTagAddRemove) ([]byte, error) {

	spotBytes, _ := json.Marshal(group)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-application-site", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowApplicationSite(groupuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, groupuid))
	req, err := http.NewRequest("POST", c.base+"/show-application-site", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteApplicationSite(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-application-site", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

///////////////////
func(c *Client) CreateServiceGroup(servicegroup ServiceGroup) ([]byte, error) {

	spotBytes, _ := json.Marshal(servicegroup)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-service-group", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetServiceGroup(servicegroup ServiceGroup) ([]byte, error) {

	spotBytes, _ := json.Marshal(servicegroup)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-service-group", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowServiceGroup(servicegroupuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, servicegroupuid))
	req, err := http.NewRequest("POST", c.base+"/show-service-group", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteServiceGroup(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-service-group", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

//////////////////////////////////////////
func(c *Client) CreateDynamicObject(dynamicobject DynamicObject) ([]byte, error) {

	spotBytes, _ := json.Marshal(dynamicobject)
	spotReader := bytes.NewReader(spotBytes)
	req, err := http.NewRequest("POST", c.base+"/add-dynamic-object", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetDynamicObject(dynamicobject DynamicObject) ([]byte, error) {

	spotBytes, _ := json.Marshal(dynamicobject)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-dynamic-object", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowDynamicObject(dynamicobjectuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, dynamicobjectuid))
	req, err := http.NewRequest("POST", c.base+"/show-dynamic-object", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteDynamicObject(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-dynamic-object", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}
/////////////////////////////////////////
func(c *Client) CreateTag(tag Tag) ([]byte, error) {

	spotBytes, _ := json.Marshal(tag)
	spotReader := bytes.NewReader(spotBytes)
	req, err := http.NewRequest("POST", c.base+"/add-tag", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetTag(tag Tag) ([]byte, error) {

	spotBytes, _ := json.Marshal(tag)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-tag", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowTag(taguid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, taguid))
	req, err := http.NewRequest("POST", c.base+"/show-tag", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteTag(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-tag", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}
/////////////////////////////////////////
func(c *Client) CreateSecurityZone(securityzone SecurityZone) ([]byte, error) {

	spotBytes, _ := json.Marshal(securityzone)
	spotReader := bytes.NewReader(spotBytes)
	req, err := http.NewRequest("POST", c.base+"/add-security-zone", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetSecurityZone(securityzone SecurityZone) ([]byte, error) {

	spotBytes, _ := json.Marshal(securityzone)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-security-zone", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowSecurityZone(securityzoneuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, securityzoneuid))
	req, err := http.NewRequest("POST", c.base+"/show-security-zone", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteSecurityZone(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-security-zone", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}
/////////////////////////////////////////
////////////////////////////////////////
func(c *Client) CreateDNSDomain(dnsdomain DNSDomain) ([]byte, error) {

	spotBytes, _ := json.Marshal(dnsdomain)
	spotReader := bytes.NewReader(spotBytes)
	req, err := http.NewRequest("POST", c.base+"/add-dns-domain", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetDNSDomain(dnsdomain DNSDomain) ([]byte, error) {

	spotBytes, _ := json.Marshal(dnsdomain)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-dns-domain", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowDNSDomain(dnsdomainuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, dnsdomainuid))
	req, err := http.NewRequest("POST", c.base+"/show-dns-domain", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteDNSDomain(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-dns-domain", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}
/////////////////////////////////////////

////////////////////////////////////////

func(c *Client) CreateNetwork(network Network) ([]byte, error) {

	spotBytes, _ := json.Marshal(network)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-network", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetNetwork(network Network) ([]byte, error) {

	spotBytes, _ := json.Marshal(network)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-network", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowNetwork(networkuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, networkuid))
	req, err := http.NewRequest("POST", c.base+"/show-network", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) DeleteNetwork(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-network", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) CreateServiceTcp(servicetcp ServiceTcp) ([]byte, error) {

	spotBytes, _ := json.Marshal(servicetcp)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-service-tcp", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) SetServiceTcp(servicetcp ServiceTcp) ([]byte, error) {

	spotBytes, _ := json.Marshal(servicetcp)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-service-tcp", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400  {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowServiceTcp(servicetcpuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, servicetcpuid))
	req, err := http.NewRequest("POST", c.base+"/show-service-tcp", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteServiceTcp(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-service-tcp", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) CreateServiceUdp(serviceudp ServiceUdp) ([]byte, error) {

	spotBytes, _ := json.Marshal(serviceudp)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-service-udp", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) SetServiceUdp(serviceudp ServiceUdp) ([]byte, error) {

	spotBytes, _ := json.Marshal(serviceudp)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-service-udp", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400  {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowServiceUdp(serviceudpuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, serviceudpuid))
	req, err := http.NewRequest("POST", c.base+"/show-service-udp", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)

	return body, err
}

func(c *Client) DeleteServiceUdp(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-service-udp", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}
	defer resp.Body.Close()
	return "", err
}

func(c *Client) CreatePolicyPackage(policypackage PolicyPackage) ([]byte, error) {

	spotBytes, _ := json.Marshal(policypackage)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-package", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetPolicyPackage(policypackage PolicyPackage) ([]byte, error) {

	spotBytes, _ := json.Marshal(policypackage)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-package", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400  {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowPolicyPackage(policypackageuid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, policypackageuid))
	req, err := http.NewRequest("POST", c.base+"/show-package", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)
	return body, err
}

func(c *Client) DeletePolicyPackage(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-package", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	//c.Publish(c.sid)
	return "", err
}

func(c *Client) CreateAccessLayer(accesslayer AccessLayer) ([]byte, error) {

	spotBytes, _ := json.Marshal(accesslayer)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-access-layer", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetAccessLayer(accesslayer AccessLayer) ([]byte, error) {

	spotBytes, _ := json.Marshal(accesslayer)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-access-layer", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400  {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowAccessLayer(accesslayeruid string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, accesslayeruid))
	req, err := http.NewRequest("POST", c.base+"/show-access-layer", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	type Uid struct {
	                  Uid string
					  			}
	var uid Uid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &uid)
	return body, err
}

func(c *Client) DeleteAccessLayer(uid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, uid))
	req, err := http.NewRequest("POST", c.base+"/delete-access-layer", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	//c.Publish(c.sid)
	return "", err
}

func(c *Client) CreateAccessRulebase(accessrulebase AccessRulebase) ([]byte, error) {

	spotBytes, _ := json.Marshal(accessrulebase)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-access-rule", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) SetAccessRulebase(accessrulebase AccessRulebase) ([]byte, error) {

	spotBytes, _ := json.Marshal(accessrulebase)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-access-rule", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400  {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowAccessRulebase(accessrulebaseuid string, accessrulebaselayer string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v","layer": "%v"}`, accessrulebaseuid, accessrulebaselayer))

	req, err := http.NewRequest("POST", c.base+"/show-access-rule", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)

	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowAccessRulebaseByName(accessrulebasename string, accessrulebaselayer string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v","layer": "%v"}`, accessrulebasename, accessrulebaselayer))

	req, err := http.NewRequest("POST", c.base+"/show-access-rule", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
  //If object is not found 404 is sent
	if resp.StatusCode == 404 {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) DeleteAccessRulebase(accessrulebaseuid string, accessrulebaselayer string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v","layer": "%v"}`, accessrulebaseuid, accessrulebaselayer))
	req, err := http.NewRequest("POST", c.base+"/delete-access-rule", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	//c.Publish(c.sid)
	return "", err
}

func(c *Client) DeleteAccessRuleByName(accessrulebasename string, accessrulebaselayer string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v","layer": "%v"}`, accessrulebasename, accessrulebaselayer))
	req, err := http.NewRequest("POST", c.base+"/delete-access-rule", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) DeleteAccessRuleByRuleNum(accessrulebaserulenum int, accessrulebaselayer string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"rule-number": "%v","layer": "%v"}`, accessrulebaserulenum, accessrulebaselayer))
	req, err := http.NewRequest("POST", c.base+"/delete-access-rule", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) CreateAccessRulebaseList(accessrulebase AccessRulebaseList) ([]byte, error) {

	spotBytes, _ := json.Marshal(accessrulebase)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-access-rule", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) SetAccessRulebaseList(accessrulebase AccessRulebaseList) ([]byte, error) {

	spotBytes, _ := json.Marshal(accessrulebase)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-access-rule", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}


func(c *Client) ShowAccessRulebaseList(accessrulebaselayer string,limit int, offset int) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v", "details-level" : "standard", "use-object-dictionary" : false, "limit" : %v, "offset" : %v}`, accessrulebaselayer, limit, offset))

	req, err := http.NewRequest("POST", c.base+"/show-access-rulebase", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
  //If object is not found 404 is sent
	if resp.StatusCode == 404 {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowNATRulebaseList(packageuid string,limit int, offset int) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"package": "%v", "details-level" : "standard", "use-object-dictionary" : false, "limit" : %v, "offset" : %v}`, packageuid, limit, offset))

	req, err := http.NewRequest("POST", c.base+"/show-nat-rulebase", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
  //If object is not found 404 is sent
	if resp.StatusCode == 404 {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}


func(c *Client) CreateAccessSection(accesssection AccessSection) ([]byte, error) {

	spotBytes, _ := json.Marshal(accesssection)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-access-section", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) SetAccessSection(accesssection AccessSectionUpdate) ([]byte, error) {

	spotBytes, _ := json.Marshal(accesssection)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-access-section", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400  {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ShowAccessSection(accesssectionuid string, accesssectionlayer string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v","layer": "%v"}`, accesssectionuid, accesssectionlayer))

	req, err := http.NewRequest("POST", c.base+"/show-access-section", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)

	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) DeleteAccessSection(accesssectionuid string, accesssectionlayer string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v","layer": "%v"}`, accesssectionuid, accesssectionlayer))
	req, err := http.NewRequest("POST", c.base+"/delete-access-section", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	return "", err
}

func(c *Client) ReadHostData(hostname string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v"}`, hostname))
	req, err := http.NewRequest("POST", c.base+"/show-host", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ReadAddressRangeData(hostname string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v"}`, hostname))
	req, err := http.NewRequest("POST", c.base+"/show-address-range", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ReadServiceTcpData(servicetcpname string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v"}`, servicetcpname))
	req, err := http.NewRequest("POST", c.base+"/show-service-tcp", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ReadServiceUdpData(serviceudpname string) ([]byte, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v"}`, serviceudpname))
	req, err := http.NewRequest("POST", c.base+"/show-service-udp", bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return body, err
}

func(c *Client) ReadLayerUIDtoName(layeruid string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v"}`, layeruid))
	req, err := http.NewRequest("POST", c.base+"/show-access-layer", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	accesslayer := AccessLayer{}
  json.Unmarshal(body, &accesslayer)
	accesslayername := accesslayer.Name

	return accesslayername, err
}

func(c *Client) ReadLayerNametoUID(layername string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"name": "%v"}`, layername))
	req, err := http.NewRequest("POST", c.base+"/show-access-layer", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	accesslayer := AccessLayer{}
  json.Unmarshal(body, &accesslayer)
	accesslayeruid := accesslayer.Uid
	return accesslayeruid, err
}

func(c *Client) CreateNATSection(packagename string, position string, name string) (string, error) {

  var jsonStr = []byte(fmt.Sprintf(`{"package": "%v", "position": "%v", "name": "%v"}`, packagename, position, name))
	req, err := http.NewRequest("POST", c.base+"/add-nat-section", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return "", fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	natsection := NATSection{}
  json.Unmarshal(body, &natsection)
	natsectionuid := natsection.Uid
	return natsectionuid, err
}

func(c *Client) Publish(sid string) (string, error) {

	var jsonStr = []byte(`{}`)
	req, err := http.NewRequest("POST", c.base+"/publish", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)

	defer resp.Body.Close()

  var taskid Taskid
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal([]byte(body), &taskid)

	//currenttask, _ := task
	status, _ := c.CheckTaskid(taskid.Taskid)
  _ = status
	return "status", err
}

func(c *Client) CheckTaskid(taskid string) (string, error) {

	var jsonStrTask = []byte(fmt.Sprintf(`{"task-id": "%v"}`, taskid))

	req, err := http.NewRequest("POST", c.base+"/show-task", bytes.NewBuffer(jsonStrTask))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	//time.Sleep(500 * time.Millisecond)
	resp, err := c.client.Do(req)
    defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	_ = body

	return "failed", err
}

func(c *Client) Logout() (string, error) {

	var jsonStr = []byte(`{}`)
	req, err := http.NewRequest("POST", c.base+"/logout", bytes.NewBuffer(jsonStr))

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)

	defer resp.Body.Close()

	return "1", err
}

func(c *Client) CreateAccessNATRulebaseList(natrulebase AccessRulebaseNATList) ([]byte, error) {

	spotBytes, _ := json.Marshal(natrulebase)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/add-nat-rule", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) SetAccessNATRulebaseList(natrulebase AccessRulebaseNATListSet) ([]byte, error) {

	spotBytes, _ := json.Marshal(natrulebase)
	spotReader := bytes.NewReader(spotBytes)

	req, err := http.NewRequest("POST", c.base+"/set-nat-rule", spotReader)

	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) DeleteAccessNATRule(uid string, packageuid string) ([]byte, error) {

	var jsonStr = []byte(fmt.Sprintf(`{"uid": "%v", "package": "%v"}`, uid, packageuid))
	req, err := http.NewRequest("POST", c.base+"/delete-nat-rule", bytes.NewBuffer(jsonStr))


	if err != nil {
		return nil, errors.New("Sorry something went wrong.  API busy??")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-chkp-sid", c.sid)
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if resp.StatusCode >= 400 {
		var errorreturn APIError
		body, _ := ioutil.ReadAll(resp.Body)
		json.Unmarshal([]byte(body), &errorreturn)
		return nil, fmt.Errorf("%d Error returned from R80API.  %v", resp.StatusCode, errorreturn)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body, err
}

func(c *Client) ConvertListtoSet(listelements []string) *schema.Set {
	values := make([]interface{}, len(listelements))
	for q := range listelements {
		values[q] = listelements[q]
	}
	return schema.NewSet(schema.HashString, values)
}
