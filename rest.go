/*
Copyright (c) 2015, Alberto Corsín Lafuente
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice,
   this list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its contributors
   may be used to endorse or promote products derived from this software
   without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
package main

import (
	"net/http"
)

func APIReturn(code int, info string, w http.ResponseWriter) {
	output := make(map[string]map[string]interface{})
	output["result"] = make(map[string]interface{})
	output["result"]["code"] = code
	output["result"]["info"] = info
	WriteJSON(output, code, w)
}

func APISingleResult(resultCode int, resultInfo string, data map[string]interface{}, w http.ResponseWriter) {
	output := make(map[string]map[string]interface{})
	output["result"] = make(map[string]interface{})
	output["result"]["code"] = resultCode
	output["result"]["info"] = resultInfo
	output["data"] = data
	WriteJSON(output, resultCode, w)
}

type APIMultipleOutput struct {
	Result map[string]interface{}   `json:"result"`
	Data   []map[string]interface{} `json:"data"`
	Paging map[string]interface{}   `json:"paging"`
}

func APIMultipleResults(resultCode int, resultInfo string, data APIMultipleOutput, w http.ResponseWriter) {
	data.Result = make(map[string]interface{})
	data.Result["code"] = resultCode
	data.Result["info"] = resultInfo
	WriteJSON(data, resultCode, w)
}
