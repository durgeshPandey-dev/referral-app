package logger

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
)

type ctxKey string

const RequestIDKey ctxKey = "request_id"

func getReqID(ctx context.Context) string {
	if ctx == nil {
		return "unknown"
	}
	if v := ctx.Value(RequestIDKey); v != nil {
		return v.(string)
	}
	return "unknown"
}

func log(color string, level string, ctx context.Context, msg string, kv map[string]interface{}) {
	reqID := getReqID(ctx)

	parts := []string{
		fmt.Sprintf("time=%s", time.Now().Format(time.RFC3339)),
		fmt.Sprintf("level=%s", level),
		fmt.Sprintf("req_id=%s", reqID),
		fmt.Sprintf("msg=%q", msg),
	}

	keys := make([]string, 0, len(kv))
	for k := range kv {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", k, kv[k]))
	}

	line := strings.Join(parts, " ")
	fmt.Printf("%s%s\033[0m\n", color, line)
}

func Info(ctx context.Context, msg string, kv map[string]interface{}) {
	log("\033[32m", "INFO", ctx, msg, kv)
}

func Error(ctx context.Context, msg string, kv map[string]interface{}) {
	log("\033[31m", "ERROR", ctx, msg, kv)
}

func Warn(ctx context.Context, msg string, kv map[string]interface{}) {
	log("\033[33m", "WARN", ctx, msg, kv)
}
