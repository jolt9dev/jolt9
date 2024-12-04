package ctxs

import "context"

type ExecContext struct {
	Env     map[string]string
	Secrets map[string]string
	Context context.Context
}
