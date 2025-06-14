package pgsql

import (
	"Stant/LestaGamesInternship/internal/domain/models"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type MetricStore struct {
	dbConn DBConn
}

func NewMetricStore(dbConn DBConn) *MetricStore {
	return &MetricStore{dbConn: dbConn}
}

const upsertMetric = `
	INSERT INTO lesta_start.metrics
	(timestamp, name, value)
	VALUES
	($1, $2, $3)
	ON CONFLICT (timestamp, name) DO UPDATE
	SET value = $3
	;
`

func (s *MetricStore) Track(ctx context.Context, metric models.Metric) error {
	if _, err := s.dbConn.Exec(ctx, upsertMetric, metric.Timestamp().UTC(), metric.Name(), metric.Value()); err != nil {
		return fmt.Errorf("pgsql/metricStore.Track: [%w]", err)
	}

	return nil
}

const checkMetric = `
	SELECT EXISTS
	(SELECT 1 FROM lesta_start.metrics
	WHERE timestamp = $1 AND name = $2
	LIMIT 1)
	;
`

func (s *MetricStore) IsTracked(ctx context.Context, timestamp time.Time, name models.MetricName) (bool, error) {
	isExist := false

	row := s.dbConn.QueryRow(ctx, checkMetric, timestamp.UTC(), name)
	if err := row.Scan(&isExist); err != nil {
		return false, fmt.Errorf("pgsql/metricStore.isTracked: [%w]", err)
	}

	return isExist, nil
}

const selectMetric = `
	SELECT timestamp, name, value
	FROM lesta_start.metrics
	WHERE timestamp = $1 AND name = $2
	LIMIT 1
	;
`

func (s *MetricStore) Find(ctx context.Context, timestamp time.Time, name models.MetricName) (models.Metric, error) {
	row := s.dbConn.QueryRow(ctx, selectMetric, timestamp.UTC(), name)

	metric, err := s.scanMetric(row)
	if err != nil {
		return models.Metric{}, fmt.Errorf("pgsql/metricStore.Find: [%w]", err)
	}

	return metric, nil
}

const selectMetricsByTimestamp = `
	SELECT timestamp, name, value
	FROM lesta_start.metrics
	WHERE timestamp = $1
	;
`

func (s *MetricStore) FindAllByTimestamp(ctx context.Context, timestamp time.Time) ([]models.Metric, error) {
	rows, err := s.dbConn.Query(ctx, selectMetricsByTimestamp, timestamp.UTC())
	if err != nil {
		return nil, fmt.Errorf("pgsql/metricStore.FindAllByTimestamp: [%w]", err)
	}

	metrics, err := s.scanMetrics(rows)
	if err != nil {
		return nil, fmt.Errorf("pgsql/metricStore.FindAllByTimestamp: [%w]", err)
	}

	return metrics, nil
}

const selectMetricsByName = `
	SELECT timestamp, name, value
	FROM lesta_start.metrics
	WHERE name = $1
	;
`

func (s *MetricStore) FindAllByName(ctx context.Context, name models.MetricName) ([]models.Metric, error) {
	rows, err := s.dbConn.Query(ctx, selectMetricsByName, name)
	if err != nil {
		return nil, fmt.Errorf("pgsql/metricStore.FindAllByName: [%w]", err)
	}

	metrics, err := s.scanMetrics(rows)
	if err != nil {
		return nil, fmt.Errorf("pgsql/metricStore.FindAllByName: [%w]", err)
	}

	return metrics, nil
}

func (s *MetricStore) scanMetric(row pgx.Row) (models.Metric, error) {
	var (
		timestamp time.Time
		name      models.MetricName
		value     float64
	)
	if err := row.Scan(&timestamp, &name, &value); err != nil {
		return models.Metric{}, fmt.Errorf("pgsql/metricStore.scanMetric: [%w]", err)
	}

	return models.NewMetric(timestamp, name, value), nil
}

func (s *MetricStore) scanMetrics(rows pgx.Rows) ([]models.Metric, error) {
	metrics := []models.Metric{}
	for rows.Next() {
		metric, err := s.scanMetric(rows)
		if err != nil {
			return nil, fmt.Errorf("pgsql/metricStore.scanMetrics: [%w]", err)
		}

		metrics = append(metrics, metric)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("pgsql/metricStore.scanMetrics: [%w]", err)
	}

	if rows.CommandTag().RowsAffected() == 0 {
		return metrics, fmt.Errorf("pgsql/collectionStore.scanMetrics: [%w]", pgx.ErrNoRows)
	}
	return metrics, nil
}
