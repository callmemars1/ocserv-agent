package logger

import (
	"log/slog"
	"log/syslog"
)

func BuildForSyslog() (*slog.Logger, error) {
	writer, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "ocserv-agent")
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewTextHandler(writer, nil))
	slog.SetDefault(logger)

	return logger, nil
}
