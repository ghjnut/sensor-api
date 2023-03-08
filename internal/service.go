package internal

import (
	"context"

	"ghjnut/sensor"
)

type Service interface {
	CreateLogs(context.Context, CreateLogsIn) (CreateLogsOut, error)
	Logs(context.Context, LogsIn) ([]sensor.Log, error)
	Device(context.Context, DeviceIn) (sensor.Device, error)
}

type CreateLogsIn struct {
	Data []string
}

type CreateLogsOut struct {
}

type LogsIn struct {
	DeviceID string
}

type DeviceIn struct {
	ID string
}
