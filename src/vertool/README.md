# About
This tool will compile your apprication and increment your version number.
# Usage
* Create a new file called **app.go** and past the code below into the file
```
package main

import "fmt"

var (
	//Version : go build/run -ldflags "-X main.Version=123"
	Version string = "Alpha"
	//Copyright : Build with -ldflags
	Copyright string = "Copyright Albie van der Merwe © 2020"
	//AppName : Name of application
	AppName string = "Application"
	//Build Date
	BuildDate string = "2021-09-09"
	//Name of the host the build was done
	BuildHost string = "localhost"
	//Service name
	ServiceName = "MyService"
	//Serving port
	PortNo = "8080"
)

func main() {
	fmt.Println("Version:\t", Version)
	fmt.Println("Copyright:\t", Copyright)
	fmt.Println("AppName:\t", AppName)
	fmt.Println("BuildDate:\t", BuildDate)
	fmt.Println("BuildHost:\t", BuildHost)
}
```
* Now create the json file, in this case the file name is app.json
```
vertool -file ./app.json -func create
```
* Now to build the application 
```
go-build -a "Version Tool" -s vertool -o vertool -p 8080
```
* Now you can run the app and it will output the new variable values
```
./app
```
* Output
```
Version:         0.3.20210923.233436
Copyright:       Copyright Albie van der Merwe © 2020
AppName:         app
BuildDate:       Thu 23 Sep 2021 23:34:36 BST
BuildHost:       MacBook-Pro.local
Service Port:    8080
```
