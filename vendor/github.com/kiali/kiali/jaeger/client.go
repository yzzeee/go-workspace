package jaeger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/kiali/kiali/config"
	jaegerModel "github.com/kiali/kiali/jaeger/model"
	jsonConv "github.com/kiali/kiali/jaeger/model/converter/json"
	jsonModel "github.com/kiali/kiali/jaeger/model/json"
	"github.com/kiali/kiali/log"
	"github.com/kiali/kiali/models"
	"github.com/kiali/kiali/util/grpcutil"
	"github.com/kiali/kiali/util/httputil"
)

// ClientInterface for mocks (only mocked function are necessary here)
type ClientInterface interface {
	GetAppTraces(ns, app string, query models.TracingQuery) (traces *JaegerResponse, err error)
	GetTraceDetail(traceId string) (*JaegerSingleTrace, error)
	GetErrorTraces(ns, app string, duration time.Duration) (errorTraces int, err error)
	GetServiceStatus() (available bool, err error)
}

// Client for Jaeger API.
type Client struct {
	ClientInterface
	grpcClient jaegerModel.QueryServiceClient
	httpClient http.Client
	baseURL    *url.URL
	ctx        context.Context
}

func NewClient(token string) (*Client, error) {
	cfg := config.Get()
	cfgTracing := cfg.ExternalServices.Tracing

	if !cfgTracing.Enabled {
		return nil, errors.New("jaeger is not enabled")
	} else {
		auth := cfgTracing.Auth
		if auth.UseKialiToken {
			auth.Token = token
		}
		ctx := context.Background()

		u, errParse := url.Parse(cfgTracing.InClusterURL)
		if !cfg.InCluster {
			u, errParse = url.Parse(cfgTracing.URL)
		}
		if errParse != nil {
			log.Errorf("Error parsing Jaeger URL: %s", errParse)
			return nil, errParse
		}

		if cfgTracing.UseGRPC {
			// GRPC client

			// Note: jaeger-query does not have built-in secured communication, at the moment it is only achieved through reverse proxies (cf https://github.com/jaegertracing/jaeger/issues/1718).
			// When using the GRPC client, if a proxy is used it has to support GRPC.
			// Basic and Token auth are in theory implemented for the GRPC client (see package grpcutil) but were not tested because openshift's oauth-proxy doesn't support GRPC at the time.
			// Leaving some commented-out code below -- perhaps useful, perhaps not -- to consider when testing secured GRPC.
			// if auth.Token != "" {
			// 	requestMetadata := metadata.New(map[string]string{
			// 		spanstore.BearerTokenKey: auth.Token,
			// 	})
			// 	ctx = metadata.NewOutgoingContext(ctx, requestMetadata)
			// }

			port := u.Port()
			if port == "" {
				p, _ := net.LookupPort("tcp", u.Scheme)
				port = strconv.Itoa(p)
			}
			opts, err := grpcutil.GetAuthDialOptions(u.Scheme == "https", &auth)
			if err != nil {
				log.Errorf("Error while building GRPC dial options: %v", err)
				return nil, err
			}
			address := fmt.Sprintf("%s:%s", u.Hostname(), port)
			log.Tracef("Jaeger GRPC client info: address=%s, auth.type=%s", address, auth.Type)
			conn, err := grpc.Dial(address, opts...)
			if err != nil {
				log.Errorf("Error while establishing GRPC connection: %v", err)
				return nil, err
			}
			client := jaegerModel.NewQueryServiceClient(conn)
			log.Infof("Create Jaeger GRPC client %s", address)
			return &Client{grpcClient: client, ctx: ctx}, nil
		} else {
			// Legacy HTTP client
			log.Tracef("Using legacy HTTP client for Jaeger: url=%v, auth.type=%s", u, auth.Type)
			timeout := time.Duration(config.Get().ExternalServices.Tracing.QueryTimeout) * time.Second
			transport, err := httputil.CreateTransport(&auth, &http.Transport{}, timeout, nil)
			if err != nil {
				return nil, err
			}
			client := http.Client{Transport: transport, Timeout: timeout}
			log.Infof("Create Jaeger HTTP client %s", u)
			return &Client{httpClient: client, baseURL: u, ctx: ctx}, nil
		}
	}
}

// GetAppTraces fetches traces of an app
func (in *Client) GetAppTraces(namespace, app string, q models.TracingQuery) (*JaegerResponse, error) {
	if in.grpcClient == nil {
		return getAppTracesHTTP(in.httpClient, in.baseURL, namespace, app, q)
	}
	jaegerServiceName := buildJaegerServiceName(namespace, app)
	findTracesRQ := &jaegerModel.FindTracesRequest{
		Query: &jaegerModel.TraceQueryParameters{
			ServiceName:  jaegerServiceName,
			StartTimeMin: timestamppb.New(q.Start),
			StartTimeMax: timestamppb.New(q.End),
			Tags:         q.Tags,
			DurationMin:  durationpb.New(q.MinDuration),
			SearchDepth:  int32(q.Limit),
		},
	}
	ctx, cancel := context.WithTimeout(in.ctx, time.Duration(config.Get().ExternalServices.Tracing.QueryTimeout)*time.Second)
	defer cancel()

	stream, err := in.grpcClient.FindTraces(ctx, findTracesRQ)
	if err != nil {
		err = fmt.Errorf("GetAppTraces, Jaeger GRPC client error: %v", err)
		log.Error(err.Error())
		return nil, err
	}

	tracesMap, err := readSpansStream(stream)
	if err != nil {
		return nil, err
	}
	r := JaegerResponse{
		Data:              []jsonModel.Trace{},
		JaegerServiceName: jaegerServiceName,
	}
	for _, t := range tracesMap {
		converted := jsonConv.FromDomain(t)
		r.Data = append(r.Data, *converted)
	}

	return &r, nil
}

// GetTraceDetail fetches a specific trace from its ID
func (in *Client) GetTraceDetail(strTraceID string) (*JaegerSingleTrace, error) {
	if in.grpcClient == nil {
		return getTraceDetailHTTP(in.httpClient, in.baseURL, strTraceID)
	}
	traceID, err := jaegerModel.TraceIDFromString(strTraceID)
	if err != nil {
		return nil, fmt.Errorf("GetTraceDetail, invalid trace ID: %v", err)
	}
	bTraceId := make([]byte, 16)
	_, err = traceID.MarshalTo(bTraceId)
	if err != nil {
		return nil, fmt.Errorf("GetTraceDetail, invalid marshall: %v", err)
	}
	getTraceRQ := &jaegerModel.GetTraceRequest{TraceId: bTraceId}

	ctx, cancel := context.WithTimeout(in.ctx, 4*time.Second)
	defer cancel()

	stream, err := in.grpcClient.GetTrace(ctx, getTraceRQ)
	if err != nil {
		return nil, fmt.Errorf("GetTraceDetail, Jaeger GRPC client error: %v", err)
	}
	tracesMap, err := readSpansStream(stream)
	if err != nil {
		return nil, err
	}
	if trace, ok := tracesMap[traceID]; ok {
		converted := jsonConv.FromDomain(trace)
		return &JaegerSingleTrace{Data: *converted}, nil
	}
	// Not found
	return nil, nil
}

// GetErrorTraces fetches number of traces in error for the given app
func (in *Client) GetErrorTraces(ns, app string, duration time.Duration) (int, error) {
	// Note: grpc vs http switch is performed in subsequent call 'GetAppTraces'
	now := time.Now()
	query := models.TracingQuery{
		Start: now.Add(-duration),
		End:   now,
		Tags:  map[string]string{"error": "true"},
	}
	for key, value := range config.Get().ExternalServices.Tracing.QueryScope {
		query.Tags[key] = value
	}

	traces, err := in.GetAppTraces(ns, app, query)
	if err != nil {
		return 0, err
	}
	return len(traces.Data), nil
}

func (in *Client) GetServiceStatus() (bool, error) {
	// Check Service Status using HTTP when gRPC is not enabled
	if in.grpcClient == nil {
		return getServiceStatusHTTP(in.httpClient, in.baseURL)
	}

	ctx, cancel := context.WithTimeout(in.ctx, 4*time.Second)
	defer cancel()

	_, err := in.grpcClient.GetServices(ctx, &jaegerModel.GetServicesRequest{})
	return err == nil, err
}

type SpansStreamer interface {
	Recv() (*jaegerModel.SpansResponseChunk, error)
	grpc.ClientStream
}

func readSpansStream(stream SpansStreamer) (map[jaegerModel.TraceID]*jaegerModel.Trace, error) {
	tracesMap := make(map[jaegerModel.TraceID]*jaegerModel.Trace)
	for received, err := stream.Recv(); err != io.EOF; received, err = stream.Recv() {
		if err != nil {
			if status.Code(err) == codes.DeadlineExceeded {
				log.Trace("Jaeger GRPC client timeout")
				break
			}
			log.Errorf("jaeger GRPC client, stream error: %v", err)
			return nil, fmt.Errorf("Jaeger GRPC client, stream error: %v", err)
		}
		for i, span := range received.Spans {
			traceId := jaegerModel.TraceID{}
			err := traceId.Unmarshal(span.TraceId)
			if err != nil {
				log.Errorf("Jaeger TraceId unmarshall error: %v", err)
				continue
			}
			if trace, ok := tracesMap[traceId]; ok {
				trace.Spans = append(trace.Spans, received.Spans[i])
			} else {
				tracesMap[traceId] = &jaegerModel.Trace{
					Spans: []*jaegerModel.Span{received.Spans[i]},
				}
			}
		}
	}
	return tracesMap, nil
}

func buildJaegerServiceName(namespace, app string) string {
	conf := config.Get()
	if conf.ExternalServices.Tracing.NamespaceSelector {
		return app + "." + namespace
	}
	return app
}
