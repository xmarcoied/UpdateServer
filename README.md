[![Build Status](https://travis-ci.org/xmarcoied/UpdateServer.svg?branch=master)](https://travis-ci.org/xmarcoied/UpdateServer)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmarcoied/updateserver)](https://goreportcard.com/report/github.com/xmarcoied/updateserver)
[![GoDoc](https://godoc.org/code.videolan.org/GSoC2017/Marco/UpdateServer?status.svg)](https://godoc.org/code.videolan.org/GSoC2017/Marco/UpdateServer)

# UpdateServer

An UpdateServer written in golang to ship releases and manage update requests


## Installing
Nothing fancy, Just use the golang ```go get```

```
go get code.videolan.org/GSoC2017/Marco/UpdateServer
```

## Configuration
  The configuration comes in JSON format, you either edit the default config.json attached or give the path of the configuration file through flags 
  Example : 
  ```
  {
    "psqlinfo": {
        "psqlhost"     : "host" , 
        "psqlname"     : "updater" ,
        "psqluser"     : "postgres" , 
        "psqlpassword" : "postgres" ,
        "psqlport"     : "5432"
    }
  }
  
  ```

## Flags
  The UpdateServer uses ```-port``` flag to customize port web application will be running on (8080 default)
  
  and ```-config``` to add a path for the configuration file 
  
## Usage
  - ```cd $GOPATH/src/code.videolan.org/GSoC2017/Marco/UpdateServer```
  
  - ```go build``` to build a binary UpdateServer

  - ```./UpdateServer``` to run the Server
  
      Optional flags: 
      ```
      ./UpdateServer -port 80
      ./UpdateServer -config $HOME/config.json
      ```
      Default admin authentication : username:admin , password:admin
      
  -  ```<host>/admin/dashboard/channels``` to add new channel .
  
      - [secure] Add only public key, and sign the metadata at every release's action.
      
      - [less secure] Add both public and private key, and the server would auto-sign the releases .

  - ```<host>/admin/dashboard/releases```  to add new release

       insert the needed fields and choose a channel and you can add rules against it later .
       
     
  - add ```<host>/u/update``` to you client with a querystring for the update_request paramaters.
    
    Example : ```<host>/u/update?product=vlc&channel=stable&os=Win&os_ver=7&os_arch=64&product_ver=2.1.4```
    
    Expected release content :
      ```
      {"id":17,"channel":"stable","os":"Linux","os_ver":"Linux","os_arch":"32","product_ver":"2.1.1","url":"localhost","title":"Title","desc":"Description","product":"vlc"}
  
      ```
      
    Note: the signed status is the JSON release without ID field.
    
    so, at verifying the releases at client remove the ```"id":<id>, ``` part .
    
  - also add ```<host>/u/signature``` to get the signature for the associated release
  
    Example : ```<host>/u/signature?id=<release_id> ```
  
  - Monitor the incoming update_requests through ```<host>/u/requests```
  
  
