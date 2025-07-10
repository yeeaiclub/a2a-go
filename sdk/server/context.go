// Copyright 2025 yeeaiclub
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

package server

import (
	"net/http"

	"github.com/yeeaiclub/a2a-go/sdk/auth"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type CallContext struct {
	User     auth.User
	Metadata map[string]any
}

func (s CallContext) SetRequest(req *http.Request) {
	panic("implement me")
}

func (s CallContext) Request() *http.Request {
	panic("implement me")
}

func (s CallContext) Set(key string, value any) {
	panic("implement me")
}

func (s CallContext) Get(key string) any {
	panic("implement me")
}

func (s CallContext) GetSecurityRequirement() types.SecurityRequirement {
	panic("implement me")
}

func (s CallContext) GetSecuritySchemes(key string) types.SecurityScheme {
	panic("implement me")
}

func (s CallContext) SetSecurityConfig(security types.SecurityRequirement, schemes map[string]types.SecurityScheme) {
	panic("implement me")
}
