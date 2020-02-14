package webserver

const (
	logFileNameTimeFormat = "2006-01-02-15"
	ServerName            = "EPG-xcache-service"
	Version               = "0.1.1-20180418"
	MaxHTTPConnect        = 30000
	ReadBufferSize        = 1024 * 2
	MaxConnsPerIP         = 5
	MaxRequestsPerConn    = 100
	MaxRequestBodySize    = 1024 * 2
	Concurrency           = 3000
	ServerAddr            = ":3001"
)
