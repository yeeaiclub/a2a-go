package auth

import (
	"net/http"
	"strings"

	"github.com/yeeaiclub/a2a-go/sdk/client"
	"github.com/yeeaiclub/a2a-go/sdk/types"
)

type Interceptor struct {
	CredentialService Credential
}

func NewInterceptor(credentialService Credential) *Interceptor {
	return &Interceptor{CredentialService: credentialService}
}

func (i *Interceptor) Intercept(request *http.Request, ctx *client.CallContext, agentCard types.AgentCard) {
	for _, requirement := range agentCard.Security {
		for key := range requirement {
			credential, err := i.CredentialService.GetCredentials(key, ctx)
			if err != nil {
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
				return
			case types.OAUTH2, types.OPENIDConnect:
				request.Header.Set("Authorization", "Bearer "+credential)
			case types.HTTP:
				s := scheme.(types.HTTPAuthSecurityScheme)
				if strings.ToLower(s.Scheme) == "bearer" {
					request.Header.Set("Authorization", "Bearer "+credential)
				}
			}
		}
	}
}
