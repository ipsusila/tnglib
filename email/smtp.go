package email

import (
	"github.com/d5/tengo/v2"
	"github.com/ipsusila/registry"
	"github.com/ipsusila/tnglib"
	mail "github.com/xhit/go-simple-mail/v2"
)

var (
	smtpSrvRegistry = registry.NewImmutableSyncMapRegistry[string, *mail.SMTPServer]()
)

// RegisterSmtpServer with given key
func RegisterSmtpServer(key string, conf *Config) error {
	// 1. Configure server
	server := mail.NewSMTPClient()
	if err := conf.Server.Configure(server); err != nil {
		return err
	}
	tlsCfg, err := conf.Tls.MakeTlsConfig()
	if err != nil {
		return err
	}
	server.TLSConfig = tlsCfg

	return smtpSrvRegistry.Register(key, server)
}

func smtpConnectFunc() tengo.CallableFunc {
	return func(args ...tengo.Object) (tengo.Object, error) {
		s, err := tnglib.ArgIToString(0, args...)
		if err != nil {
			return nil, err
		}
		srv, err := smtpSrvRegistry.Entry(s)
		if err != nil {
			return tnglib.WrapError(err), nil
		}

		// setup keep alive
		csrv := *srv
		if len(args) > 1 {
			keepAlive, err := tnglib.ArgIToBool(1, args...)
			if err != nil {
				return nil, err
			}
			csrv.KeepAlive = keepAlive
		}

		// get timeout
		timeout, err, _ := tnglib.ArgIToDuration(2, args...)
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		csrv.ConnectTimeout = timeout

		// Wrap connection to object
		cli, err := csrv.Connect()
		if err != nil {
			return tnglib.WrapError(err), nil
		}
		return makeConnection(cli), nil
	}
}
