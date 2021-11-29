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
	AppName string = "Version Tool"
	//Build date
	BuildDate string = "Today"
	//Build Host
	BuildHost string = "localhost"
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

vertool -file ./app.json -func inc -section build
vertool -file ./app.json -func inc -section minor
export APP_VERSION=`vertool -file ./app.json -func load -section version`
export APP_COPYRIGHT=`vertool -file ./app.json -func load -section copyright`
export APP_BUILD_DATE=`date`
export APP_BUILD_HOST=`hostname`
export APP_NAME="app"

go build -ldflags "-X 'main.Version=${APP_VERSION}' -X 'main.Copyright=${APP_COPYRIGHT}' -X 'main.AppName=${APP_NAME}' -X 'main.BuildDate=${APP_BUILD_DATE}' -X 'main.BuildHost=${APP_BUILD_HOST}'" -o $APP_NAME app.go
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
```
