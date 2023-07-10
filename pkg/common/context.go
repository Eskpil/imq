package common

import "context"

type Context struct {
}

func UpgradeContext(ctx context.Context) Context {
	return Context{}
}
