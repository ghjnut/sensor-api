package service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"ghjnut/sensor/internal"
)

type service struct {
	db *sql.DB
}

func New(db *sql.DB) internal.Service {
	return &service{db: db}
}

// TODO should we instead be calling CreateLog multiple times?
//func (s *service) CreateLog(context.Context, CreateLogIn) (CreateLogOut, error) {
//}

// TODO wrap in transaction and rollback
// TODO CreateLogsIn should be a struct with the 3 variables
func (s *service) CreateLogs(ctx context.Context, in internal.CreateLogsIn) (out internal.CreateLogsOut, err error) {
	for i := 0; i < len(in.Data); i++ {
		fields := strings.Split(in.Data[i], "|")
		if len(fields) != 3 {
			return out, errors.New("bad payload")
		}

		//2022-03-28T18:55:35.433Z
		t, err := time.Parse(time.RFC3339, fields[1])
		if err != nil {
			return out, err
		}

		//`TempF = (2*TempC) + 30`
		temp_c, err := strconv.Atoi(fields[2])
		if err != nil {
			return out, err
		}
		temp_f := 2*temp_c + 30

		//l := &Log{
		//  DeviceID:     fields[0],
		//  Date:         t,
		//  TemperatureF: temp_f,
		//}
		//_, err := db.Exec("INSERT INTO logs (event_date, device_id, temp_farenheit) VALUES ($1, $2, $3)", l.Date, l.DeviceID, l.TemperatureF)
		_, err = s.db.Exec("INSERT INTO logs (event_date, device_id, temp_farenheit) VALUES ($1, $2, $3)", t, fields[0], temp_f)
		if err != nil {
			return out, err
		}
	}
	return out, err
}
