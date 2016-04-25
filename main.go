package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/fclaeys/rancher-client/client"
)

const (
	accessKey   = "356AF61DBBEB86E0EE8E"
	secretKey   = "AJj8ENEq2CgberXw2axwC3Ygw2qVSNKKee7mg7FS"
	metadataURL = "http://192.168.64.5:8080/v1"
)

func main() {
	logrus.Infof("Plop")
	client := client.NewClient(metadataURL, accessKey, secretKey)

	containers, err := client.GetRunningContainers()
	if err != nil {
		logrus.Errorf("Error %v", err)
	}

	for _, container := range containers {
		logrus.Infof("Container %v", container)
	}

}
