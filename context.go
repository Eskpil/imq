package imq

import "context"

type Context struct {
}

func upgradeContext(ctx context.Context) Context {
	return Context{}
}
