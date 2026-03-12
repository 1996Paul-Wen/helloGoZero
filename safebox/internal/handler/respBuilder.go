package handler

import (
	"context"

	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/types"
	"github.com/1996Paul-Wen/helloGoZero/safebox/internal/util"
)

func BuildSuccessResp(ctx context.Context, data any) types.GeneralResponse {

	return types.GeneralResponse{
		Code:    0,
		Msg:     "",
		Data:    data,
		TraceID: util.LoadTraceFrom(ctx), // todo
	}
}

func BuildFailResp(ctx context.Context, code int, err error) types.GeneralResponse {
	return types.GeneralResponse{
		Code:    code,
		Msg:     err.Error(),
		Data:    nil,
		TraceID: util.LoadTraceFrom(ctx),
	}
}
