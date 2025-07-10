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

package middleware

import (
	"strings"

	"github.com/yeeaiclub/a2a-go/sdk/types"
	"github.com/yeeaiclub/a2a-go/sdk/web"
)

func Intercept(credential Credential) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx web.Context) error {
			for _, requirement := range ctx.GetSecurityRequirement() {
				for key := range requirement {
					credentials, err := credential.GetCredentials(key, ctx)
					if err != nil {
						continue
					}
					if credentials == "" {
						continue
					}
					scheme := ctx.GetSecuritySchemes(key)
					if scheme == nil {
						continue
					}
					switch scheme.GetType() {
					case types.APIKEY:
						o := scheme.(types.APIKeySecurityScheme)
						if o.In == "header" {
							ctx.Request().Header.Set(o.Name, credentials)
						}
					case types.HTTP:
						s := scheme.(types.HTTPAuthSecurityScheme)
						if strings.ToLower(s.Scheme) == "bearer" {
							ctx.Request().Header.Set("Authorization", "Bearer "+credentials)
						}
					case types.OAUTH2, types.OPENIDConnect:
						ctx.Request().Header.Set("Authorization", "Bearer "+credentials)
					}
				}
			}
			return nil
		}
	}
}
