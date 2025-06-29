// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

import "encoding/json"

const Version = "2.0"

type ErrorCode int

const (
	ErrorCodeParseError                   ErrorCode = -32700
	ErrorCodeInvalidRequest               ErrorCode = -32600
	ErrorCodeMethodNotFound               ErrorCode = -32601
	ErrorCodeInvalidParams                ErrorCode = -32602
	ErrorCodeInternalError                ErrorCode = -32603
	ErrorCodeTaskNotFound                 ErrorCode = -32000
	ErrorCodeTaskNotCancelable            ErrorCode = -32001
	ErrorCodePushNotificationNotSupported ErrorCode = -32002
	ErrorCodeUnsupportedOperation         ErrorCode = -32003
)

type JSONRPCRequest struct {
	Id     string `json:"id,omitempty"`
	Method string `json:"method,omitempty"`
	Params any    `json:"params,omitempty"`
}

type JSONRPCError struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type JSONRPCResponse struct {
	Id      string        `json:"id,omitempty"`
	JSONRPC string        `json:"jsonrpc"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

func JSONParseError(err error) *JSONRPCError {
	return &JSONRPCError{
		Code:    int(ErrorCodeParseError),
		Message: err.Error(),
	}
}

func MethodNotFoundError() *JSONRPCError {
	return &JSONRPCError{
		Code:    -32601,
		Message: "Method not found",
	}
}

func InternalError() *JSONRPCError {
	return &JSONRPCError{
		Code:    -32603,
		Message: "Internal error",
	}
}

func JSONRPCSuccessResponse(id string, result any) JSONRPCResponse {
	return JSONRPCResponse{
		Id:      id,
		JSONRPC: Version,
		Result:  result,
	}
}

func JSONRPCErrorResponse(id string, jsonrpcError *JSONRPCError) JSONRPCResponse {
	return JSONRPCResponse{
		Id:      id,
		JSONRPC: Version,
		Error:   jsonrpcError,
	}
}

func MapTo[T any](result any) (T, error) {
	var value T
	bytes, err := json.Marshal(result)
	if err != nil {
		return value, err
	}
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}
