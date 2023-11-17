package utils

import (
	"context"
	"runtime"

	"go.elastic.co/apm/v2"
)

func StartAPMSpan(ctx context.Context, useAPM bool, name string, spanType string) (*apm.Span, context.Context) {
	if useAPM {
		span, ctx := apm.StartSpan(ctx, name, spanType)
		return span, ctx
	}
	return nil, ctx
}

func EndAPMSpan(span *apm.Span) {
	if span != nil {
		span.End()
	}
}

func SetAPMLabel(span *apm.Span, key string, value interface{}) {
	if span != nil {
		span.Context.SetLabel(key, value)
	}
}

func GetFunctionName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}

	return fn.Name()
}

func SetAPMUsername(span *apm.Span, c context.Context, username string) {
	if span != nil {
		if tx := apm.TransactionFromContext(c); tx != nil {
			tx.Context.SetUsername(username)
		}
	}
}
