package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Счетчик авторизации
	LoginAttempts = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_login_attempts_total",
		Help: "Total number of login attempts",
	}, []string{"status"})

	// Время процесса авторизации
	LoginDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "auth_login_duration_seconds",
		Help:    "Duration of login attempts",
		Buckets: prometheus.DefBuckets,
	})

	// Кол-во созданных пользователей
	UserCreated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_users_created_total",
		Help: "Total number of users created",
	}, []string{"status"})

	// Длительность создания пользователя
	UserCreationDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "auth_user_creation_duration_second",
		Help:    "Duration of user creation",
		Buckets: prometheus.DefBuckets,
	})

	// Кол-во сгенерированных токенов
	TokenGenerated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_token_generated_total",
		Help: "Total number of tokens generated",
	}, []string{"type"})

	// Кол-во валидированных токенов (сколько прошло проверок)
	TokensValidated = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "auth_tokens_validated_total",
		Help: "Total number of token validations",
	}, []string{"status"})

	// Кол-во запросов к gRPC методам
	GRPCRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "grpc_request_total",
		Help: "Total gRPC requests",
	}, []string{"method", "status"})

	// длительность gRPC запросов
	GRPCRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "gRPC request duration",
		Buckets: prometheus.DefBuckets,
	}, []string{"method"})
)
