package context

import (
	"context"
	"github.com/mobamoh/snapsight/models"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func User(ctx context.Context) *models.User {
	value := ctx.Value(userKey)
	user, ok := value.(*models.User)
	if !ok {
		return nil
	}
	return user
}
