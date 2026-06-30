package http_server

import (
	"time"
)

type HTTPServerConfig struct {
	Addr                         string
	DisableGeneralOptionsHandler bool
	ReadTimeout                  time.Duration
	ReadHeaderTimeout            time.Duration
	WriteTimeout                 time.Duration
	IdleTimeout                  time.Duration
}

func NewHTTPServerConfig() HTTPServerConfig {
	return HTTPServerConfig{}
}

func (c HTTPServerConfig) WithAddr(addr string) HTTPServerConfig {
	c.Addr = addr
	return c
}

func (c HTTPServerConfig) WithDisableGeneralOptionsHandler(disableGeneralOptionsHandler bool) HTTPServerConfig {
	c.DisableGeneralOptionsHandler = disableGeneralOptionsHandler
	return c
}

func (c HTTPServerConfig) WithReadTimeout(readTimeout time.Duration) HTTPServerConfig {
	c.ReadTimeout = readTimeout
	return c
}

func (c HTTPServerConfig) WithReadHeaderTimeout(readHeaderTimeout time.Duration) HTTPServerConfig {
	c.ReadHeaderTimeout = readHeaderTimeout
	return c
}

func (c HTTPServerConfig) WithWriteTimeout(writeTimeout time.Duration) HTTPServerConfig {
	c.WriteTimeout = writeTimeout
	return c
}

func (c HTTPServerConfig) WithIdleTimeout(idleTimeout time.Duration) HTTPServerConfig {
	c.IdleTimeout = idleTimeout
	return c
}
