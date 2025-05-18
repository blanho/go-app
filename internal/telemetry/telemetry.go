package telemetry

import (
	"context"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/yourusername/azure-go-app/internal/config"
)

type Telemetry struct {
	client appinsights.TelemetryClient
	config *config.Config
}

func NewTelemetry(cfg *config.Config) *Telemetry {
	if cfg.ApplicationInsightsKey == "" {
		return &Telemetry{
			client: appinsights.NewTelemetryClient(""),
			config: cfg,
		}
	}
	
	client := appinsights.NewTelemetryClient(cfg.ApplicationInsightsKey)
	
	client.Context().Tags.Cloud().SetRole(cfg.ServiceName)
	client.Context().Tags.Cloud().SetRoleInstance(cfg.PodName)
	
	return &Telemetry{
		client: client,
		config: cfg,
	}
}

func (t *Telemetry) TrackRequest(ctx context.Context, name string, startTime time.Time, duration time.Duration, responseCode string, success bool) {
	t.client.TrackRequest(name, startTime, duration, responseCode, success)
}

func (t *Telemetry) TrackEvent(name string, properties map[string]string, measurements map[string]float64) {
	t.client.TrackEvent(name, properties, measurements)
}

func (t *Telemetry) TrackMetric(name string, value float64, properties map[string]string) {
	t.client.TrackMetric(name, value, properties)
}

func (t *Telemetry) TrackException(err error, properties map[string]string) {
	t.client.TrackException(err, properties, nil)
}

func (t *Telemetry) Flush() {
	<-t.client.Channel().Flush()
}