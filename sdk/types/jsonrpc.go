// Copyright 2025 yumosx
//
// Licensed under the Apache License, Version 2.0 (the \"License\");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an \"AS IS\" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

const Version = "2.0"

type JSONRPCError struct {
	Code    int64  `json:"code,omitempty"`
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
		Code:    -32700,
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
