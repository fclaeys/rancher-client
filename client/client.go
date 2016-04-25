package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Sirupsen/logrus"
)

type RancherClient struct {
	url       string
	accessKey string
	secretKey string
}

type RancherPagination struct {
	First    int  `json:"first"`
	Previous int  `json:"previous"`
	Next     int  `json:"next"`
	Limit    int  `json:"limit"`
	Total    int  `json:"total"`
	Partial  bool `json:"partial"`
}

type RancherSortLinks struct {
}

type RancherFilters struct {
}

type RancherCreateDefaults struct {
}

type RancherResponse struct {
	SortLinks      RancherSortLinks      `json:"sortLinks"`
	Pagination     RancherPagination     `json:"pagination"`
	Sort           string                `json:"sort"`
	Filters        RancherFilters        `json:"filters"`
	CreateDefaults RancherCreateDefaults `json:"createDefaults"`
}

type Container struct {
	ID          string            `json:"id"`
	UUID        string            `json:"uuid"`
	Name        string            `json:"name"`
	State       string            `json:"state"`
	Labels      map[string]string `json:"labels"`
	Env         map[string]string `json:"environment"`
	Ports       []string          `json:"ports"`
	DataVolumes []string          `json:"dataVolumes"`
}

type ContainersResponse struct {
	RancherResponse
	Data []Container `json:"data"`
}

func NewClient(url string, accessKey string, secretKey string) *RancherClient {
	return &RancherClient{url, accessKey, secretKey}
}

func (r *RancherClient) sendRequest(path string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", r.url+path, nil)
	req.Header.Add("Accept", "application/json")

	if r.accessKey != "" && r.secretKey != "" {
		base64Token := base64.StdEncoding.EncodeToString([]byte(r.accessKey + ":" + r.secretKey))
		req.Header.Add("Authorization", "Basic "+base64Token)
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error %v %v", resp.Status, path)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func (r *RancherClient) GetRunningContainers() ([]Container, error) {
	var containers ContainersResponse
	resp, err := r.sendRequest("/containers?state_eq=Running")
	if err != nil {
		return containers.Data, err
	}
	logrus.Infof("Get containers response %v", containers)

	if err = json.Unmarshal(resp, &containers); err != nil {
		return containers.Data, err
	}

	return containers.Data, nil
}
