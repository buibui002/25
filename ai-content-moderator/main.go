package main

import (
	"github.com/buibui002/ai-content-moderator/moderator"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

type httpContext struct {
	types.DefaultHttpContext
	mod moderator.Moderator
}
func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	body, err := proxywasm.GetHttpRequestBody(0, bodySize)
	if err != nil {
		proxywasm.LogError("failed to read request body: " + err.Error())
		return types.ActionContinue
	}

	result, err := ctx.mod.Check(string(body))
	if err != nil {
		proxywasm.LogError("moderator error: " + err.Error())
		return types.ActionContinue
	}

	switch result.Action {
	case moderator.ActionBlock:
		err := proxywasm.SendHttpResponse(
			403,
			[][2]string{{"Content-Type", "text/plain"}},
			[]byte("Blocked by content moderator"),
			-1,
		)
		if err != nil {
			return 0
		}
		return types.ActionPause

	case moderator.ActionModify:
		err := proxywasm.ReplaceHttpRequestBody(
			[]byte(result.ModifiedContent),
			len(result.ModifiedContent),
		)
		if err != nil {
			proxywasm.LogError("failed to replace body: " + err.Error())
		}
	}

	return types.ActionContinue
}

// 拦截 AI 响应体
func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	if endOfStream {
		body, err := proxywasm.GetHttpResponseBody(0, bodySize)
		if err != nil {
			proxywasm.LogError("failed to read response body: " + err.Error())
			return types.ActionContinue
		}

		result, err := ctx.mod.Check(string(body))
		if err != nil {
			proxywasm.LogError("moderator error: " + err.Error())
			return types.ActionContinue
		}

		switch result.Action {
		case moderator.ActionBlock:
			proxywasm.ReplaceHttpResponseBody([]byte("Blocked by content moderator"))
		case moderator.ActionModify:
			proxywasm.ReplaceHttpResponseBody([]byte(result.ModifiedContent))
		}
	}

	return types.ActionContinue
}
