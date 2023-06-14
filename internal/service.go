package internal

import (
	"context"
	"time"
)

type Service interface {
	CreateLog(context.Context, CreateLogIn) (CreateLogOut, error)
	CreateLogs(context.Context, CreateLogsIn) (CreateLogsOut, error)
	Logs(context.Context, LogsIn) ([]LogOut, error)
	Device(context.Context, DeviceIn) (DeviceOut, error)
}

// try to keep these as close to direct json encode/decode

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

type LogOut struct {
	DeviceID     string    `json:"-"`
	Date         time.Time `json:"logDate"`
	TemperatureF int       `json:"temperature"`
	// inferred
	Alert bool `json:"alert"`
}

type LogsIn struct {
	DeviceID string
}

type DeviceIn struct {
	ID string
}

type DeviceOut struct {
	ID   string   `json:"deviceId"`
	Logs []LogOut `json:"logs"`

	// inferred
	AverageTemperature int       `json:"averageTemperature"`
	MostRecentLogDate  time.Time `json:"mostRecentLogDate"`
	TotalAlerts        int       `json:"totalAlerts"`
}
