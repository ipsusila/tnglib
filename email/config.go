package email

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/ipsusila/tnglib"
	mail "github.com/xhit/go-simple-mail/v2"
)

// client auth type
var vmClientAuthType = map[string]tls.ClientAuthType{
	"NoClientCert":               tls.NoClientCert,
	"RequestClientCert":          tls.RequestClientCert,
	"RequireAnyClientCert":       tls.RequireAnyClientCert,
	"VerifyClientCertIfGiven":    tls.VerifyClientCertIfGiven,
	"RequireAndVerifyClientCert": tls.RequireAndVerifyClientCert,
}

var vmEncryptType = map[string]mail.Encryption{
	"NONE":     mail.EncryptionNone,
	"SSL":      mail.EncryptionSSL,
	"TLS":      mail.EncryptionTLS,
	"SSL/TLS":  mail.EncryptionSSLTLS,
	"STARTTLS": mail.EncryptionSTARTTLS,
}

var vmAuthType = map[string]mail.AuthType{
	"PLAIN":    mail.AuthPlain,
	"NONE":     mail.AuthNone,
	"LOGIN":    mail.AuthLogin,
	"CRAM-MD5": mail.AuthCRAMMD5,
}

// ServerConfig stores connection configuration
type ServerConfig struct {
	Authentication string `json:"authentication"`
	Encryption     string `json:"encryption"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Helo           string `json:"helo"`
	ConnectTimeout string `json:"connectTimeout"`
	SendTimeout    string `json:"sendTimeout"`
	Host           string `json:"host"`
	Port           int    `json:"port"`
	KeepAlive      bool   `json:"keepAlive"`
}

// TlsConfig definition
type TlsConfig struct {
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	ServerName         string `json:"serverName"`
	CertFile           string `json:"certFile"`
	KeyFile            string `json:"keyFile"`
	ClientAuth         string `json:"clientAuth"`
}

// Config stores configuration for the application
type Config struct {
	Server *ServerConfig `json:"server"`
	Tls    *TlsConfig    `json:"tls"`
}

// MakeTlsConfig return tls.Config from given configuration
func (t *TlsConfig) MakeTlsConfig() (*tls.Config, error) {
	tc := &tls.Config{
		InsecureSkipVerify: t.InsecureSkipVerify,
	}
	if t.CertFile != "" && t.KeyFile != "" {
		cer, err := tls.LoadX509KeyPair(t.CertFile, t.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("load certificate %s/%s error: %w", t.CertFile, t.KeyFile, err)
		}
		tc.Certificates = []tls.Certificate{cer}
	}
	if t.ServerName != "" {
		tc.ServerName = t.ServerName
	}

	// setup client auth type
	if cliAuth, found := vmClientAuthType[t.ClientAuth]; found {
		tc.ClientAuth = cliAuth
	}

	return tc, nil
}

// Configure server
func (s *ServerConfig) Configure(srv *mail.SMTPServer) error {
	srv.Host = s.Host
	srv.Port = s.Port
	srv.Username = s.Username
	srv.Password = s.Password
	srv.Helo = s.Helo
	srv.KeepAlive = s.KeepAlive

	if s.ConnectTimeout != "" {
		to, err := time.ParseDuration(s.ConnectTimeout)
		if err != nil {
			return fmt.Errorf("parsing connect timeout `%s` error: %w", s.ConnectTimeout, err)
		}
		srv.ConnectTimeout = to
	}

	if s.SendTimeout != "" {
		to, err := time.ParseDuration(s.SendTimeout)
		if err != nil {
			return fmt.Errorf("parsing send timeout `%s` error: %w", s.ConnectTimeout, err)
		}
		srv.SendTimeout = to
	}

	var ok bool
	srv.Authentication, ok = vmAuthType[s.Authentication]
	if !ok {
		srv.Authentication = mail.AuthNone
	}

	srv.Encryption, ok = vmEncryptType[s.Encryption]
	if !ok {
		srv.Encryption = mail.EncryptionNone
	}

	return nil
}

// DefaultConfig return default configuration
func DefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Helo: "localhost",
		},
		Tls: &TlsConfig{
			InsecureSkipVerify: true,
		},
	}
}

// LoadConfig loads configuration from file, either in JSON or HJSON
func LoadConfig(filename string) (*Config, error) {
	conf := DefaultConfig()
	if err := tnglib.LoadConfigTo(filename, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
