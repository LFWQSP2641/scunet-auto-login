package adapter

import "context"

type Authenticator interface {
	Login(ctx context.Context, username, password string, extra map[string]string) error
	Logout(ctx context.Context, username, password string, extra map[string]string) error
}
