package value

import (
	"context"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	DefaultTimestamp = Timestamp(time.RFC3339)

	DefaultCaller = Caller()
)

type Valuer func(ctx context.Context) any

// Value return the function value.
func Value(ctx context.Context, v any) any {
	if v, ok := v.(Valuer); ok {
		return v(ctx)
	}
	return v
}

func BindValues(ctx context.Context, keyvals []any) {
	for i := 1; i < len(keyvals); i += 2 {
		if v, ok := keyvals[i].(Valuer); ok {
			keyvals[i] = v(ctx)
		}
	}
}

func ContainsValuer(keyvals []any) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}

// Timestamp returns a timestamp Valuer with a custom time format.
func Timestamp(layout string) Valuer {
	return func(context.Context) any {
		return time.Now().Format(layout)
	}
}

// Caller returns a Valuer that returns a pkg/file:line description of the caller.
func Caller() Valuer {
	return func(context.Context) any {
		depth := 3
		maxDepth := 7
		for {
			if depth > maxDepth {
				return ""
			}
			_, file, line, _ := runtime.Caller(depth)
			idx := strings.LastIndexByte(file, '/')
			if idx == -1 {
				return file[idx+1:] + ":" + strconv.Itoa(line)
			}
			idx = strings.LastIndexByte(file[:idx], '/')
			if strings.HasPrefix(file[idx+1:], "holog") && (strings.Contains(file[idx+1:], "log.go") || strings.Contains(file[idx+1:], "global.go")) {
				depth++
			} else {
				return file[idx+1:] + ":" + strconv.Itoa(line)
			}
		}
	}
}
