package polardb

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
)

// UpgradeDBClusterMinorVersion invokes the polardb.UpgradeDBClusterMinorVersion API synchronously
func (client *Client) UpgradeDBClusterMinorVersion(request *UpgradeDBClusterMinorVersionRequest) (response *UpgradeDBClusterMinorVersionResponse, err error) {
	response = CreateUpgradeDBClusterMinorVersionResponse()
	err = client.DoAction(request, response)
	return
}

// UpgradeDBClusterMinorVersionWithChan invokes the polardb.UpgradeDBClusterMinorVersion API asynchronously
func (client *Client) UpgradeDBClusterMinorVersionWithChan(request *UpgradeDBClusterMinorVersionRequest) (<-chan *UpgradeDBClusterMinorVersionResponse, <-chan error) {
	responseChan := make(chan *UpgradeDBClusterMinorVersionResponse, 1)
	errChan := make(chan error, 1)
	err := client.AddAsyncTask(func() {
		defer close(responseChan)
		defer close(errChan)
		response, err := client.UpgradeDBClusterMinorVersion(request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- response
		}
	})
	if err != nil {
		errChan <- err
		close(responseChan)
		close(errChan)
	}
	return responseChan, errChan
}

// UpgradeDBClusterMinorVersionWithCallback invokes the polardb.UpgradeDBClusterMinorVersion API asynchronously
func (client *Client) UpgradeDBClusterMinorVersionWithCallback(request *UpgradeDBClusterMinorVersionRequest, callback func(response *UpgradeDBClusterMinorVersionResponse, err error)) <-chan int {
	result := make(chan int, 1)
	err := client.AddAsyncTask(func() {
		var response *UpgradeDBClusterMinorVersionResponse
		var err error
		defer close(result)
		response, err = client.UpgradeDBClusterMinorVersion(request)
		callback(response, err)
		result <- 1
	})
	if err != nil {
		defer close(result)
		callback(nil, err)
		result <- 0
	}
	return result
}

// UpgradeDBClusterMinorVersionRequest is the request struct for api UpgradeDBClusterMinorVersion
type UpgradeDBClusterMinorVersionRequest struct {
	*requests.RpcRequest
	ResourceOwnerId      requests.Integer `position:"Query" name:"ResourceOwnerId"`
	PlannedEndTime       string           `position:"Query" name:"PlannedEndTime"`
	ResourceOwnerAccount string           `position:"Query" name:"ResourceOwnerAccount"`
	DBClusterId          string           `position:"Query" name:"DBClusterId"`
	OwnerAccount         string           `position:"Query" name:"OwnerAccount"`
	OwnerId              requests.Integer `position:"Query" name:"OwnerId"`
	PlannedStartTime     string           `position:"Query" name:"PlannedStartTime"`
	FromTimeService      requests.Boolean `position:"Query" name:"FromTimeService"`
}

// UpgradeDBClusterMinorVersionResponse is the response struct for api UpgradeDBClusterMinorVersion
type UpgradeDBClusterMinorVersionResponse struct {
	*responses.BaseResponse
	RequestId string `json:"RequestId" xml:"RequestId"`
}

// CreateUpgradeDBClusterMinorVersionRequest creates a request to invoke UpgradeDBClusterMinorVersion API
func CreateUpgradeDBClusterMinorVersionRequest() (request *UpgradeDBClusterMinorVersionRequest) {
	request = &UpgradeDBClusterMinorVersionRequest{
		RpcRequest: &requests.RpcRequest{},
	}
	request.InitWithApiInfo("polardb", "2017-08-01", "UpgradeDBClusterMinorVersion", "polardb", "openAPI")
	request.Method = requests.POST
	return
}

// CreateUpgradeDBClusterMinorVersionResponse creates a response to parse from UpgradeDBClusterMinorVersion response
func CreateUpgradeDBClusterMinorVersionResponse() (response *UpgradeDBClusterMinorVersionResponse) {
	response = &UpgradeDBClusterMinorVersionResponse{
		BaseResponse: &responses.BaseResponse{},
	}
	return
}