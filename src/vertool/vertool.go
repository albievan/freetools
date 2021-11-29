package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	//Version : go build/run -ldflags "-X main.Version=123"
	Version string = "alpha"
	//Copyright : Build with -ldflags
	Copyright string = "Copyright Albie van der Merwe © 2020"
	//AppName : Name of application
	AppName string = "Vertool"
	//Build Date
	BuildDate string = "2021-09-09"
	//Name of the host the build was done
	BuildHost string = "localhost"
	//Service name
	ServiceName = "Export"
	//Serving port
	PortNo = 8080

//To build an application do the following
//go build -ldflags "-X 'main.Version=${APP_VERSION}' -X 'main.Copyright=${APP_COPYRIGHT}'"  -o appname app.go
)

//VersionInfo : information about the build
type VersionInfo struct {
	Major     int64  `json:"major" comment:"Major version"`
	Minor     int64  `json:"minor" comment:"Minor version"`
	Build     string `json:"build" comment:"Build version"`
	Copyright string `json:"copyright" comment:"Copyright notice"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getVersionInfo(fileName string) VersionInfo {
	dat, err := ioutil.ReadFile(fileName)
	check(err)
	v := VersionInfo{
		Major:     0,
		Minor:     0,
		Build:     time.Now().Format("20060102.150405"),
		Copyright: "Copyright Albie van der Merwe © 2020",
	}
	v.ToVersionInfo(string(dat))
	return v
}

func getSection(section string, v VersionInfo) string {

	switch strings.ToLower(section) {
	case "build":
		return fmt.Sprintf("%s", v.Build)
	case "major":
		return fmt.Sprintf("%d", v.Major)
	case "minor":
		return fmt.Sprintf("%d", v.Minor)
	case "copyright":
		return v.Copyright
	case "version":
		return fmt.Sprintf("%d.%d.%s", v.Major, v.Minor, v.Build)
	default:
		return ""
	}

}

func saveToFile(fileName string, v VersionInfo) {
	err := ioutil.WriteFile(fileName, []byte(v.String()), 0644)
	check(err)
}

func main() {

	verFilePtr := flag.String("file", "?", "Path to json version file")
	funcPtr := flag.String("func", "help", "[create , read, inc, version]")
	sectionSectionPtr := flag.String("section", "build", "[major, minor, build, copyright, load]")
	copyrightPtr := flag.String("crnotice", "Copyright Albie van der Merwe © 2020", "Copyright notice")
	flag.Parse()

	if *verFilePtr == "?" {
		if strings.ToLower(*funcPtr) != "version" {
			*funcPtr = "help"
		}
	}

	switch strings.ToLower(*funcPtr) {
	case "version":
		fmt.Println(fmt.Sprintf("\n%s version %s %s\n", AppName, Version, Copyright))
	case "create":
		v := VersionInfo{
			Major:     0,
			Minor:     0,
			Build:     time.Now().Format("20060102.150405"),
			Copyright: "Copyright Albie van der Merwe © 2020",
		}
		saveToFile(*verFilePtr, v)
	case "read":
		v := getVersionInfo(*verFilePtr)
		fmt.Println(fmt.Sprintf("Version: %d.%d.%s", v.Major, v.Minor, v.Build))
		fmt.Println(fmt.Sprintf("%s", v.Copyright))
	case "load":
		v := getVersionInfo(*verFilePtr)
		fmt.Print(getSection(*sectionSectionPtr, v))
	case "inc":
		v := getVersionInfo(*verFilePtr)
		switch strings.ToLower(*sectionSectionPtr) {
		case "build":
			v.Build = time.Now().Format("20060102.150405")
			saveToFile(*verFilePtr, v)
		case "major":
			v.Major = v.Major + 1
			saveToFile(*verFilePtr, v)
		case "minor":
			v.Minor = v.Minor + 1
			saveToFile(*verFilePtr, v)
		case "copyright":
			v.Copyright = *copyrightPtr
			saveToFile(*verFilePtr, v)
		default:
			fmt.Println("Nothing to do")
		}

	default:
		fmt.Println("-file [path to file version] -func [reate, read, inc] -section [major, minor, build]")
		fmt.Println(fmt.Sprintf("\nVersion tool version %s %s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
			Version, Copyright, "Usage",
			"Set Copyright notice\n\t-file ./version.json -func inc -section copyright -crnotice \"Copyright Albie van der Merwe © 2020\"",
			"Incrament major version\n\t-file ./version.json -func inc -section major",
			"Incrament minor version\n\t-file ./version.json -func inc -section minor",
			"Incrament build version\n\t-file ./version.json -func inc -section build",
			"Read version file\n\t-file ./version.json -func read",
			"Load a section\n\t-file ./version.json -func load -section [major, minor, build, version, copyright]",
		))
	}
}

//String Return structure as json string
func (v *VersionInfo) String() string {
	//ctime := time.Now()
	e, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		return err.Error()
	}
	//log.Println(fmt.Sprintf("VersionInfo.String(): %s", time.Now().Sub(ctime)))
	return string(e)
}

// ToVersionInfo Return a JSON string
// jsonstr : string to struct
// return error of the string does not containe valid json
func (v *VersionInfo) ToVersionInfo(jsonstr string) error {
	//	ctime := time.Now()
	err := json.Unmarshal([]byte(jsonstr), v)
	if err != nil {
		log.Println(err)
		return err
	}
	//log.Println(fmt.Sprintf("VersionInfo.ToVersionInfo(): %s", time.Now().Sub(ctime)))
	return nil
}
