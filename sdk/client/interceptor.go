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

package client

import (
	"net/http"
	"strings"

	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type Interceptor struct {
	CredentialService Credential
}

func NewInterceptor(credentialService Credential) *Interceptor {
	return &Interceptor{CredentialService: credentialService}
}

func (i *Interceptor) Intercept(request *http.Request, ctx *CallContext, agentCard types.AgentCard) {
	for _, requirement := range agentCard.Security {
		for key := range requirement {
			credential, err := i.CredentialService.GetCredentials(key, ctx)
			if err != nil {
				continue
			}
			if credential == "" {
				continue
			}
			scheme, ok := agentCard.SecuritySchemes[key]
			if !ok || scheme == nil {
				continue
			}
			switch scheme.GetType() {
			case types.APIKEY:
				o := scheme.(types.APIKeySecurityScheme)
				if o.In == "header" {
					request.Header.Set(o.Name, credential)
				}
			case types.HTTP:
				s := scheme.(types.HTTPAuthSecurityScheme)
				if strings.ToLower(s.Scheme) == "bearer" {
					request.Header.Set("Authorization", "Bearer "+credential)
				}
			case types.OAUTH2, types.OPENIDConnect:
				request.Header.Set("Authorization", "Bearer "+credential)
			}
		}
	}
}
