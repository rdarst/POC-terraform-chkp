#!/usr/bin/python

import getopt
import requests
import sys
import os



# Set login info
url = os.environ['CHKP_SERVER'] + "/publish" 
payload = "{}"
headers = {
    'Content-Type': "application/json",
    'Cache-Control': "no-cache",
    'X-chkp-sid': "" + os.environ['CHKP_SID'] + "",
    }

# SSL Certificate Checking is disabled!!!
requests.packages.urllib3.disable_warnings()
response = requests.request("POST", url, data=payload, headers=headers, verify=False)

#Retrun the response
print(response.text)


