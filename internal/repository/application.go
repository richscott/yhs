package repository

import (
	"context"
	"fmt"

	"github.com/apache/yunikorn-core/pkg/webservice/dao"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (s *RepoPostgres) UpsertApplications(ctx context.Context, apps []*dao.ApplicationDAOInfo) error {
	upsertSQL := `INSERT INTO applications (id, app_id, used_resource, max_used_resource, pending_resource,
			partition, queue_name, submission_time, finished_time, requests, allocations, state,
			"user", groups, rejected_message, state_log, place_holder_data, has_reserved, reservations,
			max_request_priority)
			VALUES (@id, @app_id,@used_resource, @max_used_resource, @pending_resource, @partition, @queue_name,
			@submission_time, @finished_time, @requests, @allocations, @state, @user, @groups,
			@rejected_message, @state_log, @place_holder_data, @has_reserved, @reservations, @max_request_priority)
		ON CONFLICT (partition, queue_name, app_id) DO UPDATE SET
			used_resource = EXCLUDED.used_resource,
			max_used_resource = EXCLUDED.max_used_resource,
			pending_resource = EXCLUDED.pending_resource,
			finished_time = EXCLUDED.finished_time,
			requests = EXCLUDED.requests,
			allocations = EXCLUDED.allocations,
			state = EXCLUDED.state,
			rejected_message = EXCLUDED.rejected_message,
			state_log = EXCLUDED.state_log,
			place_holder_data = EXCLUDED.place_holder_data,
			has_reserved = EXCLUDED.has_reserved,
			reservations = EXCLUDED.reservations,
			max_request_priority = EXCLUDED.max_request_priority`

	for _, a := range apps {
		_, err := s.dbpool.Exec(ctx, upsertSQL,
			pgx.NamedArgs{
				"id":                   uuid.NewString(),
				"app_id":               a.ApplicationID,
				"used_resource":        a.UsedResource,
				"max_used_resource":    a.MaxUsedResource,
				"pending_resource":     a.PendingResource,
				"partition":            a.Partition,
				"queue_name":           a.QueueName,
				"submission_time":      a.SubmissionTime,
				"finished_time":        a.FinishedTime,
				"requests":             a.Requests,
				"allocations":          a.Allocations,
				"state":                a.State,
				"user":                 a.User,
				"groups":               a.Groups,
				"rejected_message":     a.RejectedMessage,
				"state_log":            a.StateLog,
				"place_holder_data":    a.PlaceholderData,
				"has_reserved":         a.HasReserved,
				"reservations":         a.Reservations,
				"max_request_priority": a.MaxRequestPriority,
			})
		if err != nil {
			return fmt.Errorf("could not insert application into DB: %v", err)
		}
	}
	return nil
}

func (s *RepoPostgres) GetAllApplications(ctx context.Context) ([]*dao.ApplicationDAOInfo, error) {
	selectSQL := `SELECT * FROM applications`

	var apps []*dao.ApplicationDAOInfo

	rows, err := s.dbpool.Query(ctx, selectSQL)
	if err != nil {
		return nil, fmt.Errorf("could not get applications from DB: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var app dao.ApplicationDAOInfo
		var id string
		err := rows.Scan(&id, &app.ApplicationID, &app.UsedResource, &app.MaxUsedResource, &app.PendingResource,
			&app.Partition, &app.QueueName, &app.SubmissionTime, &app.FinishedTime, &app.Requests, &app.Allocations,
			&app.State, &app.User, &app.Groups, &app.RejectedMessage, &app.StateLog, &app.PlaceholderData,
			&app.HasReserved, &app.Reservations, &app.MaxRequestPriority)
		if err != nil {
			return nil, fmt.Errorf("could not scan application from DB: %v", err)
		}
		apps = append(apps, &app)
	}
	return apps, nil
}

func (s *RepoPostgres) GetAppsPerPartitionPerQueue(ctx context.Context, partition, queue string) ([]*dao.ApplicationDAOInfo, error) {
	selectSQL := `SELECT * FROM applications WHERE queue_name = $1 AND partition = $2`

	var apps []*dao.ApplicationDAOInfo

	rows, err := s.dbpool.Query(ctx, selectSQL, queue, partition)
	if err != nil {
		return nil, fmt.Errorf("could not get applications from DB: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var app dao.ApplicationDAOInfo
		var id string
		err := rows.Scan(&id, &app.ApplicationID, &app.UsedResource, &app.MaxUsedResource, &app.PendingResource,
			&app.Partition, &app.QueueName, &app.SubmissionTime, &app.FinishedTime, &app.Requests, &app.Allocations,
			&app.State, &app.User, &app.Groups, &app.RejectedMessage, &app.StateLog, &app.PlaceholderData,
			&app.HasReserved, &app.Reservations, &app.MaxRequestPriority)
		if err != nil {
			return nil, fmt.Errorf("could not scan application from DB: %v", err)
		}
		apps = append(apps, &app)
	}
	return apps, nil
}
