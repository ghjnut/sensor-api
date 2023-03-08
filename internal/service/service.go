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

func (s *service) CreateLog(ctx context.Context, in internal.CreateLogIn) (out internal.CreateLogOut, err error) {
	//2022-03-28T18:55:35.433Z
	t, err := time.Parse(time.RFC3339, in.EventDate)
	if err != nil {
		return out, err
	}

	//`TempF = (2*TempC) + 30`
	temp_c, err := strconv.Atoi(in.TempF)
	if err != nil {
		return out, err
	}
	temp_f := 2*temp_c + 30

	_, err = s.db.Exec("INSERT INTO logs (event_date, device_id, temp_farenheit) VALUES ($1, $2, $3)", t, in.DeviceID, temp_f)
	return out, err
}

// TODO wrap in transaction and rollback
func (s *service) CreateLogs(ctx context.Context, in internal.CreateLogsIn) (out internal.CreateLogsOut, err error) {
	for i := 0; i < len(in.Data); i++ {
		fields := strings.Split(in.Data[i], "|")
		if len(fields) != 3 {
			return out, errors.New("bad payload")
		}
		clin := internal.CreateLogIn{
			DeviceID:  fields[0],
			EventDate: fields[1],
			TempF:     fields[2],
		}
		_, err = s.CreateLog(ctx, clin)
		if err != nil {
			return out, err
		}
	}
	return out, err
}

func (s *service) Logs(ctx context.Context, in internal.LogsIn) (logs []internal.LogOut, err error) {
	rows, err := s.db.Query("SELECT event_date, device_id, temp_farenheit FROM logs WHERE device_id = $1", in.DeviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l internal.LogOut
		if err := rows.Scan(&l.Date, &l.DeviceID, &l.TemperatureF); err != nil {
			return logs, err
		}
		if l.TemperatureF > 32 {
			l.Alert = true
		}
		logs = append(logs, l)
	}
	// redundant, but explicit
	if err = rows.Err(); err != nil {
		return logs, err
	}

	return logs, err
}

func (s *service) Device(ctx context.Context, in internal.DeviceIn) (dev internal.DeviceOut, err error) {
	dev.ID = in.ID

	var lin internal.LogsIn
	lin.DeviceID = in.ID
	logs, err := s.Logs(ctx, lin)
	if err != nil {
		return dev, err
	}
	dev.Logs = logs

	total_temp := 0
	latest_date := time.Time{}
	alert_cnt := 0
	for _, l := range dev.Logs {
		total_temp += l.TemperatureF

		// work for first element since t is zeroed
		if l.Date.After(latest_date) {
			latest_date = l.Date
		}

		if l.Alert {
			alert_cnt++
		}
	}
	dev.AverageTemperature = total_temp / len(dev.Logs)
	dev.MostRecentLogDate = latest_date
	dev.TotalAlerts = alert_cnt

	return dev, err
}
