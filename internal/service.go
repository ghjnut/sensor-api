package internal

import (
	"context"

	"ghjnut/sensor"
)

type Service interface {
	CreateLog(context.Context, CreateLogIn) (CreateLogOut, error)
	CreateLogs(context.Context, CreateLogsIn) (CreateLogsOut, error)
	Logs(context.Context, LogsIn) ([]sensor.Log, error)
	Device(context.Context, DeviceIn) (sensor.Device, error)
}

type CreateLogIn struct {
	EventDate string
	DeviceID  string
	TempF     string
}

type CreateLogOut struct {
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
