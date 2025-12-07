package natsx

type NatsConfig interface {
	GetNatsURL() string
	GetNatsJWTUserFilePath() string
	GetNatsEstimatorStream() string
	GetNatsEstimatorStreamTopics() []string
	GetNatsStreamReplicas() int
	GetClientName() string
}
