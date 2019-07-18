#!/usr/bin/python

import getopt
import requests
import sys
import json
import os

user = ''
password = ''
mgmtserverip = ''

if len(sys.argv) <= 6:
   print 'Error - Format should be - login.py -u <username> -p <password> -s <mgmt_server_ip>'
   exit(1)

try:
   opts, args = getopt.getopt(sys.argv[1:],"u:p:s:", ['usr=','pass=','serverip=', 'help'])
except getopt.GetoptError:
   print 'Error - Format should be - login.py -u <username> -p <password> -s <mgmt_server_ip>'
   sys.exit(2)
for opt, arg in opts:
   if opt in ('-h', '--help'):
      print 'login.py -u <username> -p <password> -s <mgmt_server_ip>'
      sys.exit()
   elif opt in ("-u", "--usr"):
      user = arg
   elif opt in ("-p", "--pass"):
      password = arg
   elif opt in ('-s', '--serverip'):
      mgmtserverip = arg

# Set login info
url = "https://" + mgmtserverip + "/web_api/login"
payload = "{\r\n  \"user\" : \"" + user + "\",\r\n  \"password\" : \"" + password + "\" ,\r\n  \"session-timeout\" : \"3600\"\r\n}"
headers = {
    'Content-Type': "application/json",
    'Cache-Control': "no-cache",
    }

# SSL Certificate Checking is disabled!!!
requests.packages.urllib3.disable_warnings()
response = requests.request("POST", url, data=payload, headers=headers, verify=False)

#Retrun the response
#print(response.text)
sid_json = json.loads(response.text)
print "export CHKP_SID=\"" + sid_json['sid'] + "\""
print "export CHKP_SERVER=\"https://" + mgmtserverip + "/web_api\""
