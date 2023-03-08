package internal

import (
	"context"
)

type Service interface {
	CreateLogs(context.Context, CreateLogsIn) (CreateLogsOut, error)
}

type CreateLogsIn struct {
	Data []string
}

type CreateLogsOut struct {
}
