# About
This utility provides the capability to dynamically update GoDaddy DNS records and seamlessly add any missing records with your ip address your isp allocates you.

The Dockerfile establishes a builder that leverages the latest Golang version to build the dns-updater application. Upon completing the application build, Docker generates a container image. 

For further information on creating the API key and secret, please visit https://developer.godaddy.com/getstarted.

# Usage
* Create a .env file and add the following environment variables
```
DDNS=<replace with apikey>
DDNS_SEC=<replace with api secret>
TZ=<your timezone>
DDNS_INTERVAL=<replace with desired value e.g. 30m>
DDNS_FILE=<location and file name of the dns records>
```
* To create the api key and secret goto https://developer.godaddy.com/getstarted

* To build the docker image run the following command
```
docker build -t dns-updater:latest .
```

* After the build of the container image you're ready to create the dns-updater container. You can use a mapped volume if needed.
```
 docker run --restart unless-stopped --name ddns --hostname ddns -dit --env-file -v /mnt/docker_vol/ddns/conf:/conf ./.env dns-updater
```
 
# JSON
To proceed with the setup, create a file named domains.json using the following structure, and place it in the same folder as your main.go and Dockerfile. Pay close attention to these crucial fields:
* domain: the domain needs to exist and owned by the api key and secret
* name: this will be the name of the record
* type: this is the type of record to add/update e.g. A, CNAME etc.
* data: this is the value of the record e.g. you want to update a CNAME record 'home' for domain another_example.com to point to ip address 1.1.1.1 and the ttl must be 600 seconds
```
[
  {
    "domain": "example.com",
    "type": "A",
    "name": "@",
    "ttl": 600,
    "data": "",
    "port": null,
    "service": "",
    "protocol": "",
    "priority": null
  },
  {
    "domain": "example.com",
    "type": "CNAME",
    "data": "@",
    "name": "dns-updater",
    "ttl": 600
  },
  {
    "domain": "another_example.com",
    "type": "CNAME",
    "data": "1.1.1.1",
    "name": "home",
    "ttl": 600
  }
]

```


