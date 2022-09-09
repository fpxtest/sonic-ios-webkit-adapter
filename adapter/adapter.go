/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package adapters

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"log"
	"sonic-ios-webkit-adapter/entity"
	"sonic-ios-webkit-adapter/protocols"
	"strings"
)

type Adapter struct {
	targetID          string
	messageFilters    map[string]protocols.MessageAdapters
	isTargetBased     bool
	applicationID     *string
	pageID            *int
	waitingForID      int
	adapterRequestMap map[int]func(message []byte)
	// 给iOS
	sendWebkit func([]byte)
	// 给devtool
	sendDevTool func([]byte)
	// recv for IOS
	receiveWebKit func([]byte)
	// recv for devtool
	receiveDevTool func([]byte)
}

func (a *Adapter) AddMessageFilter(method string, filter protocols.MessageAdapters) {
	if a.messageFilters == nil {
		a.messageFilters = make(map[string]protocols.MessageAdapters)
	}
	a.messageFilters[method] = filter
}

func (a *Adapter) CallTarget(method string, params interface{}, callFunc func(message []byte)) {
	a.waitingForID -= 1
	var message = &entity.TargetProtocol{}
	arr, err := json.Marshal(params)
	if err != nil {
		log.Fatal(err)
	}
	println(string(arr))
	message.ID = a.waitingForID
	message.Method = method
	message.Params = params
	a.adapterRequestMap[a.waitingForID] = callFunc
	a.sendToTarget(message)
}

func (a *Adapter) sendToTarget(message *entity.TargetProtocol) {
	log.Println("origin send message:")
	arr, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(arr))
	if a.isTargetBased {
		if !strings.Contains(message.Method, "Target") {
			var newMessage = &entity.TargetProtocol{}

			newMessage.ID = message.ID
			newMessage.Method = "Target.sendMessageToTarget"
			newMessage.Params = &entity.TargetParams{
				TargetId: a.targetID,
				Message:  string(arr),
				ID:       message.ID,
			}
			message = newMessage
			log.Println("new send message:")
		}
	}
	arr, err = json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	a.sendWebkit(arr)
}

func (a *Adapter) FireEventToTools(method string, params interface{}) {
	response := map[string]interface{}{
		"method": method,
		"params": params,
	}
	arr, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}
	a.sendDevTool(arr)
}

func (a *Adapter) FireResultToTools(id int, params interface{}) {
	response := map[string]interface{}{
		"id":     id,
		"result": params,
	}
	arr, err := json.Marshal(response)
	if err != nil {
		log.Panic(err)
	}
	a.sendDevTool(arr)
}

func (a *Adapter) ReplyWithEmpty(msg string) []byte {
	a.FireResultToTools(int(gjson.Get(msg, "id").Int()), map[string]interface{}{})
	return nil
}

func (a *Adapter) SetTargetID(targetID string) {
	a.targetID = targetID
}

func (a *Adapter) SetSendWebkit(sendWebkit func([]byte)) {
	a.sendWebkit = sendWebkit
}

func (a *Adapter) SetSendDevTool(sendDevTool func([]byte)) {
	a.sendDevTool = sendDevTool
}

func (a *Adapter) SetReceiveWebkit(receiveWebkit func([]byte)) {
	a.receiveWebKit = receiveWebkit
}

func (a *Adapter) SetReceiveDevTool(receiveDevTool func([]byte)) {
	a.receiveDevTool = receiveDevTool
}
