package utils

import (
	"encoding/json"
	"os"
)

type AWSConfig struct {
	AWSRegion   string `json:"aws_region"`
	DynamoTable string `json:"dynamo_table"`
}

type RdsConfig struct {
}

func ReadAWSConfig(filePath string) (AWSConfig, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return AWSConfig{}, err
	}
	defer file.Close()

	var config AWSConfig
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	return config, err
}
