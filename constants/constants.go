package constants

type HeaderKey string

const (
	HEADER_CORRELATION_ID HeaderKey = "X-Correlation-Id"
)

type ContextKey string

const (
	CORRELATION_ID_CONTEXT_KEY ContextKey = "correlationId"
)
