package utils

import (
	"context"
	"time"
)

func GetDBContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}
