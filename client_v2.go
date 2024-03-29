/*
HAPROXY DATAPLANE API CLIENT (V2)
AUTHOR: Omar Aouini & Allan Nava
mail: aouini.omar93@gmail.com
Date: 19/04/2022
Update: 21/10/2022

MIT License

Copyright (c) 2022 OmarAouini | Allan Nava

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
//
package haproxy
//
import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

// the Haproxy client type Interface
type IHaproxyClient interface {
	//GetBasicInfo() (*HaproxyInfo, error)
	HealthCheck() error
	GetVersion() (*string, error) // get haproxy version
	GetSites() (*HaproxySites, error)
	GetStats() (*HaproxyStats, error)
	GetReloads() (*HaproxyReloads, error)
	GetTransactions() (*HaproxyTransactions, error)
	GetConfigurationGlobal() (*HaproxyConfigurationGlobal, error)
	GetConfigurationDefaults() (*HaproxyConfigurationDefaults, error)
	GetBackends() (*HaproxyBackends, error)
	GetFrontends() (*HaproxyFrontends, error)
	GetBackendSwitchingRules(frontend string) (*HaproxyBackendSwitchingRules, error)
	GetServers(backend string) (*HaproxyFrontends, error)
	GetAcls(parentType string, parentName string) (*HaproxyAcls, error) //parentType eg: "backend" or "frontend"
	GetServerSwitchingRules(backend string) (*HaproxyServerSwitchingRules, error)
	GetHttpRequestRules(parentType string, parentName string) (*HaproxyHttpRequestRules, error)
	AddBackend(transactionId string, backend *HaproxyAddBackend) error
	AddFrontend(transactionId string, addFrontend *HaproxyAddFrontend) error
	AddAcl(parenttype string, parentName string, transactionId string, addAcl *HaproxyAddAcl) error
	AddServer(backend string, transactionId string, addServer *HaproxyAddServer) error
	AddHttpRequestRule(parentType string, parentName string, transactionId string, addRule *HaproxyAddHttpRequestRule) error
	AddBackendSwitchingRule(frontend string, transactionId string, addRule *HaproxyAddBackendSwitchingRule) error
	StartTransaction(haproxyVersion string) (*string, error)
	CommitTransaction(transactionId string) error
	CheckDuplicateDefinitions() (*HaproxyDuplicateDefinitionsResult, error) // check for duplicate definitions in the haproxy cfg
}

type haproxyClient struct {
	Url   string
	Rest  *resty.Client
	Debug bool
}

// create a new Haproxy client
//
// the debug argument is used in case you want to print the request for debug purposes.
//
// example usage:
//
// client, error := haproxy.NewHaproxyClient("http://127.0.0.1", "user", "password", false)
func NewHaproxyClient(haproxyUrl string, basicAuthUsername string, basicAuthPassword string, debug bool) (IHaproxyClient, error) {
	client := haproxyClient{
		Url:  haproxyUrl,
		Rest: resty.New().SetBasicAuth(basicAuthUsername, basicAuthPassword),
	}
	if debug {
		client.Rest.SetDebug(true)
		client.Debug = true
		log.Println("Debug mode is enabled for the Haproxy client ")
	}
	err := client.HealthCheck()
	if err != nil {
		return nil, err
	}
	return &client, nil
}
/*func (h *haproxyClient) GetBasicInfo() (*HaproxyInfo, error) {
	if h.Debug {
		log.Println("GetBasicInfo called() ", h.Url)
	}
	url := h.Url + "/v2/services/haproxy/info"
	response := HaproxyInfo{}
	resp, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		if h.Debug {
			log.Println("resp ", resp)
		}
		return nil, err
	}
	return &response, nil
}*/
//
func (h *haproxyClient) HealthCheck() error {
	if h.Debug {
		log.Println("GetBasicInfo called() ", h.Url)
	}
	url := h.Url + "/v2/"
	resp, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  Get(url)
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status(), "200") {
		if h.Debug {
			log.Println("resp ", resp)
		}
		return errors.New("Could not connect haproxy")
	}
	return nil
}

func (h *haproxyClient) GetVersion() (*string, error) {
	if h.Debug {
		log.Println("GetVersion called()")
	}
	url := h.Url + "/v2/services/haproxy/info"
	response :=  HaproxyInfo{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response.Haproxy.Version, nil
}

func (h *haproxyClient) GetSites() (*HaproxySites, error) {
	if h.Debug {
		log.Println("GetSites called()")
	}
	url := h.Url + "/v2/services/haproxy/sites"
	response := HaproxySites{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) GetStats() (*HaproxyStats, error) {
	if h.Debug {
		log.Println("GetStats called()")
	}
	url := h.Url + "/v2/services/haproxy/sites"
	response := HaproxyStats{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) GetReloads() (*HaproxyReloads, error) {
	if h.Debug {
		log.Println("GetReloads called()")
	}
	url := h.Url + "/v2/services/haproxy/reloads"
	response := HaproxyReloads{}
	_, err := h.Rest.R().
		SetHeader("Accept", "application/json").
		SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) GetTransactions() (*HaproxyTransactions, error) {
	if h.Debug {
		log.Println("GetTransactions called()")
	}
	url := h.Url + "/v2/services/haproxy/transactions"
	response := HaproxyTransactions{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) StartTransaction(haproxyVersion string) (*string, error) {
	if h.Debug {
		log.Println("StartTransaction called()")
	}
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/transactions?version=%s", haproxyVersion)
	response := HaproxyTransaction{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Post(url)
	if err != nil {
		return nil, err
	}
	return &response.ID, nil
}

func (h *haproxyClient) GetConfigurationGlobal() (*HaproxyConfigurationGlobal, error) {
	if h.Debug {
		log.Println("GetConfigurationGlobal called()")
	}
	url := h.Url + "/v2/services/haproxy/configuration/global"
	response := HaproxyConfigurationGlobal{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) GetConfigurationDefaults() (*HaproxyConfigurationDefaults, error) {
	url := h.Url + "/v2/services/haproxy/configuration/defaults"
	response := HaproxyConfigurationDefaults{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) GetBackends() (*HaproxyBackends, error) {
	url := h.Url + "/v2/services/haproxy/configuration/backends"
	response := HaproxyBackends{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
		//.(*HaproxyErrorResponse)
	}
	return &response, nil
}

func (h *haproxyClient) GetFrontends() (*HaproxyFrontends, error) {
	url := h.Url + "/v2/services/haproxy/configuration/frontends"
	response := HaproxyFrontends{}
	_, err := h.Rest.R().
		SetHeader("Accept", "application/json").
		SetResult(&response).Get(url)
	if err != nil {
		return nil, err
		//.(*HaproxyErrorResponse)
	}
	return &response, nil
}

func (h *haproxyClient) GetBackendSwitchingRules(frontend string) (*HaproxyBackendSwitchingRules, error) {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/backend_switching_rules?frontend=%s", frontend)
	response := HaproxyBackendSwitchingRules{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
		//.(*HaproxyErrorResponse)
	}
	return &response, nil
}

func (h *haproxyClient) GetServers(backend string) (*HaproxyFrontends, error) {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/servers?backend=%s", backend)
	response := HaproxyFrontends{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
		//.(*HaproxyErrorResponse)
	}
	return &response, nil
}

// parent type: backend or frontend
func (h *haproxyClient) GetAcls(parentType string, parentName string) (*HaproxyAcls, error) {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/acls?parent_type=%s&parent_name=%s", parentType, parentName)
	response := HaproxyAcls{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (h *haproxyClient) GetServerSwitchingRules(backend string) (*HaproxyServerSwitchingRules, error) {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/server_switching_rules?backend=%s", backend)
	response := HaproxyServerSwitchingRules{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (h *haproxyClient) GetHttpRequestRules(parentType string, parentName string) (*HaproxyHttpRequestRules, error) {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/http_request_rules?parent_type=%s&parent_name=%s", parentType, parentName)
	response := HaproxyHttpRequestRules{}
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&response).Get(url)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (h *haproxyClient) AddFrontend(transactionId string, addFrontend *HaproxyAddFrontend) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/frontends?transaction_id=%s", transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyAddFrontend{}).SetBody(addFrontend).Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (h *haproxyClient) AddBackend(transactionId string, addBackend *HaproxyAddBackend) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/backends?transaction_id=%s", transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyAddBackend{}).SetBody(addBackend).Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (h *haproxyClient) AddAcl(parenttype string, parentName string, transactionId string, addAcl *HaproxyAddAcl) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/acls?parent_type=%s&parent_name=%s&transaction_id=%s", parenttype, parentName, transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyAddAcl{}).SetBody(addAcl).Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (h *haproxyClient) AddServer(backend string, transactionId string, addServer *HaproxyAddServer) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/servers?backend=%s&transaction_id=%s", backend, transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyAddServer{}).SetBody(addServer).Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (h *haproxyClient) AddHttpRequestRule(parentType string, parentName string, transactionId string, addRule *HaproxyAddHttpRequestRule) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/http_request_rules?parent_type=%s&parent_name=%s&transaction_id=%s", parentType, parentName, transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyAddHttpRequestRule{}).SetBody(addRule).Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (h *haproxyClient) AddBackendSwitchingRule(frontend string, transactionId string, addRule *HaproxyAddBackendSwitchingRule) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/configuration/backend_switching_rules?frontend=%s&transaction_id=%s", frontend, transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyAddBackendSwitchingRule{}).SetBody(addRule).Post(url)
	if err != nil {
		return err
	}

	return nil
}

func (h *haproxyClient) CommitTransaction(transactionId string) error {
	url := h.Url + fmt.Sprintf("/v2/services/haproxy/transactions/%s", transactionId)
	_, err := h.Rest.R().
      SetHeader("Accept", "application/json").
	  SetResult(&HaproxyCommitTransaction{}).Put(url)
	if err != nil {
		return err.(*HaproxyErrorResponse)
	}

	return nil
}

func (h *haproxyClient) CheckDuplicateDefinitions() (*HaproxyDuplicateDefinitionsResult, error) {
	result := HaproxyDuplicateDefinitionsResult{}

	// //BACKENDS
	// backends, err := h.GetBackends()
	// if err != nil {
	// 	return nil, err
	// }
	// backendsNames := []string{}
	// for _, v := range backends.Data {
	// 	backendsNames = append(backendsNames, v.Name)
	// }
	// result.Backends = dupesCheck(backendsNames)

	// //FRONTENDS
	// frontends, err := h.GetFrontends()
	// if err != nil {
	// 	return nil, err
	// }
	// frontendNames := []string{}
	// for _, v := range frontends.Data {
	// 	frontendNames = append(frontendNames, v.Name)
	// }
	// result.Frontends = dupesCheck(frontendNames)

	// //ACLS
	//FOR EACH BACKEND and FRONTEND
	// acls, err := h.GetAcls("backend", "gpu-f-c0-1")
	// if err != nil {
	// 	return nil, err
	// }
	// aclsnames := []string{}
	// for _, v := range acls.Data {
	// 	aclsnames = append(aclsnames, v.AclName)
	// }
	// result.Acls = dupesCheck(aclsnames)

	//SERVERS
	//FOR EACH BACKEND
	// servers, err := h.GetServers("gpu-f-c0-1")
	// if err != nil {
	// 	return nil, err
	// }
	// serversNames := []string{}
	// for _, v := range servers.Data {
	// 	serversNames = append(serversNames, v.Name)
	// }
	// result.Servers = dupesCheck(serversNames)

	return &result, nil
}

func dupesCheck(list []string) []DuplicateCount {
	duplicate_frequency := make(map[string]int)
	result := []DuplicateCount{}

	for _, item := range list {
		// check if the item/element exist in the duplicate_frequency map
		_, exist := duplicate_frequency[item]

		if exist {
			duplicate_frequency[item] += 1 // increase counter by 1 if already in the map
		} else {
			duplicate_frequency[item] = 1 // else start counting from 1
		}
	}

	for k, v := range duplicate_frequency {
		if v > 1 { //only add duplicate to result (more than 1 count)
			res := DuplicateCount{}
			res.Name = k
			res.Count = v
			result = append(result, res)
		}
	}

	return result
}
