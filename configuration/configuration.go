package configuration

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tkanos/gonfig"
)

type Configuration struct {
	Remote struct {
		Verb      string `json:"Verb"`
 		Endpoint  string `json:"Endpoint"`
		Service   string `json:"Service"`
		Operation string `json:"Operation"`
		Format    string `json:"Format"`
		ResponseGroup string `json:"ResponseGroup"`
	} `json:"Remote"`
	Credentials struct {
		AssociateTag string `json:"AssociateTag"`
		AWSSecretKey string `json:"AWSSecretKey"`
		AWSAccessKeyId string `json:"AWSAccessKeyId"`
	} `json:"Credentials"`
	Api struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Endpoints struct {
			Amazon string `json:"amazon"`
			Item string `json:"item"`
		}
	} `json:"Api"`
}

func InitConfiguration() Configuration {
	configuration := Configuration{}
	err := gonfig.GetConf(getFileName(), &configuration)
	if err != nil {
		fmt.Println("error " + err.Error())
		os.Exit(1)
	}
	return configuration
}

func getFileName() string {
	filename := []string{"conf", ".json"}
	_, dirname, _, _ := runtime.Caller(0)
	filePath := path.Join(filepath.Dir(dirname), strings.Join(filename, ""))

	return filePath
}
