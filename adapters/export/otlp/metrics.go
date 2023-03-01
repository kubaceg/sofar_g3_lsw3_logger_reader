package otlp

import (
	"context"
	"fmt"
	"github.com/kubaceg/sofar_g3_lsw3_logger_reader/adapters/devices/sofar"
	grpc "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	http "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	sdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"log"
)

const (
	appName = "sofar.logger"
)

type Config struct {
	Http struct {
		Url string `yaml:"url"`
	} `yaml:"http"`
	Grpc struct {
		Url string `yaml:"url"`
	} `yaml:"grpc"`
}

type Service struct {
	m         metric.Meter
	reader    sdk.Reader
	exporters []sdk.Exporter
}

func New(c *Config) (*Service, error) {
	reader := sdk.NewManualReader()
	mp := sdk.NewMeterProvider(
		sdk.WithReader(reader),
		sdk.WithResource(newResource()),
	)

	global.SetMeterProvider(mp)
	m := global.Meter(appName)

	exporters := make([]sdk.Exporter, 0)
	if url := c.Grpc.Url; url != "" {
		e, err := grpc.New(context.Background(), grpc.WithEndpoint(url), grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, e)
	}

	if url := c.Http.Url; url != "" {
		e, err := http.New(context.Background(), http.WithEndpoint(url), http.WithInsecure())
		if err != nil {
			return nil, err
		}
		exporters = append(exporters, e)
	}

	s := Service{
		m:         m,
		reader:    reader,
		exporters: exporters,
	}

	err := s.initGauges()
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// initGauges creates Int64 gauges for all reply fields that will be read
func (s *Service) initGauges() error {
	for _, rr := range sofar.AllRegisterRanges {
		for _, f := range rr.ReplyFields {

			if f.Name == "" || f.ValueType == "" {
				// Measurements without a name or value type are ignored in replies
				continue
			}

			name := f.Name
			g := s.createGauge(name)
			_, err := s.m.RegisterCallback(
				// this function is called when a collection is triggered
				func(ctx context.Context, o metric.Observer) error {
					measurements := sofar.GetLastReading()
					if v, ok := measurements[name]; ok {
						o.ObserveInt64(*g, convertToInt64(v))
					} else {
						log.Printf("could not find measurement for %s", name)
					}
					return nil
				}, *g)
			if err != nil {
				log.Println("error registering gauge callback")
				return err
			}
		}
	}
	return nil
}

func (s *Service) createGauge(n string) *instrument.Int64ObservableGauge {
	newGauge, _ := s.m.Int64ObservableGauge(
		appName+"."+n,
		instrument.WithUnit("1"),
	)
	return &newGauge
}

// CollectAndPushMetrics triggers the collection and export of metrics over OTLP
func (s *Service) CollectAndPushMetrics(ctx context.Context, measurements map[string]interface{}) error {
	err := s.collectAndPushMetrics(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) collectAndPushMetrics(ctx context.Context) error {
	rm := metricdata.ResourceMetrics{}
	err := s.reader.Collect(ctx, &rm)
	if err != nil {
		return err
	}

	for _, e := range s.exporters {
		err = e.Export(ctx, rm)
		if err != nil {
			return err
		}
	}

	return nil
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(appName),
			semconv.ServiceVersion("v0.1.0"),
		),
	)
	return r
}

func convertToInt64(v interface{}) int64 {
	switch v.(type) {
	case uint32:
		u := v.(uint32)
		return int64(u)
	case uint16:
		u := v.(uint16)
		return int64(u)
	case int16:
		u := v.(int16)
		return int64(u)
	default:
		fmt.Println("unexpected type encountered")
		return 0
	}
}
