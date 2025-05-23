package main

import (
	"encoding/json"
	"time"

	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm"
	"github.com/proxy-wasm/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {}

func init() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	delay *time.Duration
}

// Override types.DefaultPluginContext.
func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{contextID: contextID, pluginContext: p}
}

type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID     uint32
	pluginContext *pluginContext
}

var additionalHeaders = map[string]string{
	"x-envoy-wasm-plugin": "delay",
}

func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	var message map[string]string
	if err := json.Unmarshal(data, &message); err == nil {
		d, err := time.ParseDuration(message["delay"])
		if err == nil {
			ctx.delay = &d
		}
	} else {
		proxywasm.LogCriticalf("error unmarshal configuration: %v", err)
	}

	return types.OnPluginStartStatusOK
}

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogErrorf("OnHttpResponseHeaders")
	for key, value := range additionalHeaders {
		proxywasm.AddHttpResponseHeader(key, value)
	}

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogErrorf("OnHttpStreamDone")

	reqHeaders, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogErrorf("failed to get request headers: %v", err)
		return
	}
	for _, h := range reqHeaders {
		proxywasm.LogErrorf("request header: <%s: %s>", h[0], h[1])
	}
}

func (ctx *httpContext) OnHttpRequestHeaders(int, bool) types.Action {
	proxywasm.LogErrorf("OnHttpRequestHeaders")

	// if ctx.pluginContext.delay != nil {
	// 	time.Sleep(*ctx.pluginContext.delay)
	// }

	return types.ActionContinue
}
