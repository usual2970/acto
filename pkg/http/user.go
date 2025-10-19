package http

import "context"

// UserInfo 认证后的用户信息
type UserInfo struct {
	Username string
	Role     string
	Claims   map[string]any
}

type ctxKeyUserInfo struct{}

// GetUserFromContext 获取认证后的用户信息
func GetUserFromContext(ctx context.Context) (*UserInfo, bool) {
	v := ctx.Value(ctxKeyUserInfo{})
	if v == nil {
		return nil, false
	}
	info, ok := v.(*UserInfo)
	return info, ok
}

// WithUserInfo 写入用户信息到 context
func WithUserInfo(ctx context.Context, user *UserInfo) context.Context {
	return context.WithValue(ctx, ctxKeyUserInfo{}, user)
}
