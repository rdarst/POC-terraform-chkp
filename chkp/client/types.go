package chkp

type Login struct {
	User string `json:"user"`
	Password     string `json:"password"`
	Domain string `json:"domain"`
}

type LoginResponse struct {
	Sid        string `json:"sid,omitempty"`
	Timeout int `json:"session-timeout,omitempty"`
}

type APIError struct {
	Code        string
	Message     string
}

type Host struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Ipv4address          string              `json:"ipv4-address,omitempty"`
	Color                string              `json:"color,omitempty"`
	Newname				       string			      	 `json:"new-name,omitempty"`
  NatSettings											         `json:"nat-settings,omitempty"`
}

type AddressRange struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Ipv4addressfirst     string              `json:"ipv4-address-first,omitempty"`
	Ipv4addresslast      string              `json:"ipv4-address-last,omitempty"`
	Color                string              `json:"color,omitempty"`
	Newname				       string			      	 `json:"new-name,omitempty"`
  NatSettings											         `json:"nat-settings,omitempty"`
}

type Network struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Subnet4              string              `json:"subnet4"`
	Masklength4          int                 `json:"mask-length4"`
	Color                string              `json:"color,omitempty"`
	Newname				       string			      	 `json:"new-name,omitempty"`
	NatSettings											         `json:"nat-settings,omitempty"`
}

type ServiceTcp struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Port                 string              `json:"port,omitempty"`
	Protocol             string              `json:"protocol,omitempty"`
	SessionTimeout       int                 `json:"session-timeout,omitempty"`
	MatchBySig           bool                `json:"match-by-protocol-signature,omitempty"`
	MatchForAny          bool                `json:"match-for-any"`
	Color                string              `json:"color,omitempty"`
	Newname				       string			      	 `json:"new-name,omitempty"`
}

type ServiceUdp struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Port                 string              `json:"port,omitempty"`
	Protocol             string              `json:"protocol,omitempty"`
	SessionTimeout       int                 `json:"session-timeout,omitempty"`
	MatchBySig           bool                `json:"match-by-protocol-signature,omitempty"`
	MatchForAny          bool                `json:"match-for-any"`
	Color                string              `json:"color,omitempty"`
	Newname				       string			      	 `json:"new-name,omitempty"`
}

type PolicyPackage struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Access               bool                `json:"access,omitempty"`
	Color                string              `json:"color,omitempty"`
	DesktopSecurity	     bool	  		      	 `json:"desktop-security"`
	Qos	                 bool	  		      	 `json:"qos"`
	QosPolicyType        string              `json:"qos-policy-type,omitempty"`
	ThreatPrevention		 bool                `json:"threat-prevention"`
	Newname				       string			      	 `json:"new-name,omitempty"`
}

type AccessLayer struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Newname				       string			      	 `json:"new-name,omitempty"`
	AppAndUrl            bool                `json:"application-and-url-filtering,omitempty"`
	ContentAwareness     bool			           `json:"content-awareness,omitempty"`
	Firewall             bool			           `json:"firewall,omitempty"`
	MobileAccess         bool			           `json:"mobile-access"`
	Shared							 bool                `json:"shared,omitempty"`
	Color                string              `json:"color,omitempty"`
	Comments             string              `json:"comments,omitempty"`
	AddDefaultRule       bool                `json:"add-default-rule,omitempty"`
}

type GroupMembers struct {
        Uid            string							`json:"uid,omitempty"`
        Name           string             `json:"name,omitempty"`
        Color            string             `json:"color,omitempty"`
        Members        []struct {
                Uid          string       `json:"uid,omitempty"`
                Name         string       `json:"name,omitempty"`
                Type         string       `json:"type,omitempty"`
        }
}

type Group struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Color                string              `json:"color,omitempty"`
	Members			       []string			 	       `json:"members,omitempty"`
	Newname				       string			 	       `json:"new-name,omitempty"`
}

type DynamicObject struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Color                string              `json:"color,omitempty"`
	Newname				       string			 	       `json:"new-name,omitempty"`
}

type AccessRulebase struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Layer                string              `json:"layer,omitempty"`
	Position			            			 	       `json:"position,omitempty"`
  Action			         string			 	       `json:"action,omitempty"`
	InlineLayer	         string			 	       `json:"inline-layer,omitempty"`
	Source			         []string			 	     `json:"source"`
  SourceNegate         bool                `json:"source-negate"`
  Destination			     []string		 	       `json:"destination"`
	DestinationNegate    bool                `json:"destination-negate"`
	Enabled	    		     bool		    	       `json:"enabled"`
  Track																     `json:"track,omitempty"`
	Newname				       string			 	       `json:"new-name,omitempty"`
	Service              []string            `json:"service,omitempty"`
}

type AccessRulebaseList struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Layer                string              `json:"layer,omitempty"`
	Position			       int     			    	 `json:"position,omitempty"`
  Action			         string			 	       `json:"action,omitempty"`
	InlineLayer	         string			 	       `json:"inline-layer,omitempty"`
	Source			         []string			 	     `json:"source"`
  SourceNegate         bool                `json:"source-negate"`
  Destination			     []string		 	       `json:"destination"`
	DestinationNegate    bool                `json:"destination-negate"`
	Enabled	    		     bool		    	       `json:"enabled"`
  Track																     `json:"track,omitempty"`
	Newname				       string			 	       `json:"new-name,omitempty"`
	Rulenumber           int							  `json:"rule-number,omitempty"`
	Service              []string            `json:"service,omitempty"`
}

type AccessRulebaseResultRead struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	From 	               int                 `json:"from,omitempty"`
	To 	                 int                 `json:"to,omitempty"`
	Total 	             int                 `json:"total,omitempty"`
	AccessRulebaseResult 										 `json:"rulebase,omitempty"`
}

type AccessSection struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Layer                string              `json:"layer,omitempty"`
	Position			            			 	       `json:"position,omitempty"`
	Newname				       string			 	       `json:"new-name,omitempty"`
}

type AccessSectionUpdate struct {
	Uid                  string              `json:"uid,omitempty"`
	Name                 string              `json:"name,omitempty"`
	Layer                string              `json:"layer,omitempty"`
	Newname				       string			 	       `json:"new-name,omitempty"`
}

type Track                struct				 	{
					Type						string				 `json:"type,omitempty"`
					Alert						string				 `json:"alert,omitempty"`
					Accounting      bool           `json:"accounting,omitempty"`
					PerConnection   bool           `json:"per-connection,omitempty"`
					PerSession      bool           `json:"per-session,omitempty"`
				}

type TrackReturn          struct				 	{
					Type						string				 `json:"type,omitempty"`
					Name						string				 `json:"name,omitempty"`
					Uid  						string				 `json:"uid,omitempty"`
				}

type NatSettings          struct				 	{
					Hidebehind			string				 `json:"hide-behind,omitempty"`
					Ipaddress			  string				 `json:"ip-address,omitempty"`
					Autorule        bool           `json:"auto-rule"`
					Installon       string         `json:"install-on,omitempty"`
					Method          string         `json:"method,omitempty"`
				}

type AccessRulebaseResult struct {
					Uid                  string              `json:"uid,omitempty"`
					Name                 string              `json:"name,omitempty"`
					Layer                string              `json:"layer,omitempty"`
					Position			            			 	       `json:"position,omitempty"`
				  Action			         string			 	       `json:"action,omitempty"`
					InlineLayer	              			 	       `json:"inline-layer,omitempty"`
					Source			                			 	     `json:"source,omitempty"`
          SourceNegate         bool                `json:"source-negate"`
					Destination			            		 	       `json:"destination,omitempty"`
					DestinationNegate    bool                `json:"destination-negate"`
					Enabled	    		     bool		    	       `json:"enabled"`
				  Track																     `json:"track,omitempty"`
        	Newname				       string			 	       `json:"new-name,omitempty"`
					Service                                  `json:"service,omitempty"`
				}

type AccessSectionResult struct {
					Uid                  string              `json:"uid,omitempty"`
					Name                 string              `json:"name,omitempty"`
					Layer                string              `json:"layer,omitempty"`
					Position			            			 	       `json:"position,omitempty"`
					Newname				       string			 	       `json:"new-name,omitempty"`
				}

type Source  []struct				 	{
					Name						string				 `json:"name,omitempty"`
					Uid						  string				 `json:"uid,omitempty"`
					Type            string         `json:"type,omitempty"`
									}
type Destination  []struct				 	{
					Name						string				 `json:"name,omitempty"`
					Uid						  string				 `json:"uid,omitempty"`
					Type            string         `json:"type,omitempty"`
																		}

type InlineLayer  []struct				 	{
					Name						string				 `json:"name,omitempty"`
				  Uid						  string				 `json:"uid,omitempty"`
					Type            string         `json:"type,omitempty"`
																		}

type Position    struct				 	{
					Above						string				 `json:"above,omitempty"`
					Below					  string				 `json:"below,omitempty"`
					Top             string         `json:"top,omitempty"`
					Bottom          string         `json:"bottom,omitempty"`																														}


type Service  []struct				 	{
					Name						string				 `json:"name,omitempty"`
				  Uid						  string				 `json:"uid,omitempty"`
				  Type            string         `json:"type,omitempty"`
																		}
type Taskid struct {
    Taskid               string              `json:"task-id"`
}

type ErrorMessage struct {
    Code               string              `json:"code"`
		Message              string              `json:"message"`
}
