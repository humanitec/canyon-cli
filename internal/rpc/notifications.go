package rpc

import "context"

type ctxKeyNtype struct {
}

var NotificationChannelKey = &ctxKeyNtype{}

func GetNotificationChannel(ctx context.Context) chan<- JsonRpcNotification {
	v, _ := ctx.Value(NotificationChannelKey).(chan<- JsonRpcNotification)
	return v
}
