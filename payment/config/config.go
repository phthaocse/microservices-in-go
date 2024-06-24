package config

import (
	"fmt"
	"github.com/huseyinbabal/microservices/payment/internal/adapters/db"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetEnv() string {
	return getEnvironmentValue("ENV")
}

func GetDataSourceURL() string {
	return getEnvironmentValue("DATA_SOURCE_URL")
}

func GetApplicationPort() int {
	portStr := getEnvironmentValue("APPLICATION_PORT")
	port, err := strconv.Atoi(portStr)

	if err != nil {
		log.Fatalf("port: %s is invalid", portStr)
	}

	return port
}
func getEnvironmentValue(key string) string {
	if os.Getenv(key) == "" {
		log.Fatalf("%s environment variable is missing.", key)
	}

	return os.Getenv(key)
}

func ReadMongoConfig() *db.MongoConfig {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(dir, "config")
	fmt.Println("config path dir: ", path)
	//Get base config
	config := new(db.MongoConfig)
	readConfig(path, "mongodb.yaml", config)
	return config
}

func envMapper(input string) string {

	if len(input) < 3 || input[0] != ':' || input[len(input)-1] != ':' {
		return input
	}

	input = input[1 : len(input)-1]

	split := strings.SplitN(input, ":", 2)
	result, ok := os.LookupEnv(split[0])
	if !ok && len(split) == 2 {
		return split[1]
	}
	return result
}

func readConfig(configPath, configName string, result interface{}) {
	file := filepath.Join(configPath, configName)
	log.Printf("Reading config file : %v", file)
	data, err := os.ReadFile(file)
	if err == nil {
		data = []byte(os.Expand(string(data), envMapper))
		err = yaml.Unmarshal(data, result)
		if err != nil {
			panic(fmt.Errorf("fatal error config file: %s - %s", configName, err))
		}
	}
	log.Printf("read file error: %v", err)
}
