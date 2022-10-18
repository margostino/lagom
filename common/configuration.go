package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
)

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}

func GetConfig(resource string) *Configuration {
	var configuration Configuration
	ymlFile, err := ioutil.ReadFile(resource)

	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	ymlFile = []byte(os.ExpandEnv(string(ymlFile)))
	err = yaml.Unmarshal(ymlFile, &configuration)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return &configuration
}
