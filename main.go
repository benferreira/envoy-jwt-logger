package main

import (
	"encoding/base64"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/valyala/fastjson"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	claimsToLog []string
}

// OnPluginStart sets the plugin context with the claims it is expected to log
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	claims, err := parseConfigClaims(string(data))

	if err != nil {
		proxywasm.LogCriticalf("failed to parse plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	if len(claims) == 0 || claims[0] == "" {
		proxywasm.LogCritical("no claims to log, check config.configuration.value")
		return types.OnPluginStartStatusFailed
	}

	ctx.claimsToLog = claims
	proxywasm.LogInfof("jwt claim logging for claims: %v", ctx.claimsToLog)

	return types.OnPluginStartStatusOK
}

func parseConfigClaims(data string) ([]string, error) {
	parsedJson, err := fastjson.Parse(string(data))

	if err != nil {
		return nil, err
	}

	v := parsedJson.GetArray("claims")
	claims := make([]string, 0)

	for _, claim := range v {
		claims = append(claims, string(claim.GetStringBytes()))
	}

	return claims, nil
}

type httpHeaders struct {
	types.DefaultHttpContext
	contextID   uint32
	claimsToLog *[]string
}

func (ctx *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpHeaders{contextID: contextID, claimsToLog: &ctx.claimsToLog}
}

// OnHttpRequestHeaders iterates over http request headers looking for the authorization header.
// If an authorization header is found and it contains a JWT, the JWT is decoded and the claims are logged
func (ctx *httpHeaders) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	headers, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
		return types.ActionContinue
	}

	//Locate the authorization header and log jwt claims, if present
	for _, header := range headers {
		if header[0] == "authorization" {
			logJWTClaims(header[1], *ctx.claimsToLog)
		}
	}

	return types.ActionContinue
}

// logJWTClaims will log the contents of each claim in the provided jwt
func logJWTClaims(jwt string, claims []string) {
	jwtParts := strings.Split(jwt, ".")

	if len(jwtParts) < 2 || len(jwtParts) > 3 {
		proxywasm.LogErrorf("invalid jwt structure has %d parts; only 2-3 are allowed", len(jwtParts))
		return
	}

	payload, err := base64.RawURLEncoding.DecodeString(jwtParts[1])

	if err != nil {
		proxywasm.LogErrorf("failed to decode jwt payload: %v", err)
		return
	}

	//Parse payload
	parsedJson, err := fastjson.Parse(string(payload))

	if err != nil {
		proxywasm.LogErrorf("failed to parse jwt payload to json: %v", err)
		return
	}

	//Log claims
	for _, claim := range claims {
		proxywasm.LogInfof("jwt[%s]:'%s'", claim, parsedJson.GetStringBytes(claim))
	}
}
