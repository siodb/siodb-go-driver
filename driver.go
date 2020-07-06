// Copyright (C) 2019-2020 Siodb GmbH. All rights reserved.
// Use of this source code is governed by a license that can be found
// in the LICENSE file.

package siodb

import (
	"crypto/rsa"
	"crypto/tls"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"net"
	"net/url"
	"os/user"
	"strconv"
)

type Config struct {
	protocol       string // https://golang.org/pkg/net/#Dial
	host           string // Network address
	port           string // Siodb port number
	user           string // Username
	identityFile   string // Public key for user
	privateKey     *rsa.PrivateKey
	unixSocketPath string // Unix socket path
	trace          bool   // Trace Siodb protol?
}

type SiodbDriver struct{}

func init() {
	sql.Register("siodb", &SiodbDriver{})
}

func ParseURI(URI string) (cfg Config, err error) {

	// Set default
	cfg.protocol = "siodbs"
	cfg.host = "localhost"
	cfg.port = "50000"
	cfg.identityFile = "~/.ssh/id_rsa"
	cfg.trace = false
	cfg.unixSocketPath = "/run/siodb/siodb.socket"
	if usr, err := user.Current(); err == nil {
		cfg.user = usr.Username
	}

	// Overwrite default with provided URI
	uri, err := url.Parse(URI)
	if uri.Scheme != "siodbs" && uri.Scheme != "siodb" && uri.Scheme != "siodbu" {
		return cfg, &SiodbDriverError{"Paring URI: unknown scheme '" + uri.Scheme + "'"}
	}
	cfg.protocol = uri.Scheme

	if len(uri.User.Username()) > 0 {
		cfg.user = uri.User.Username()
	}

	if len(uri.Hostname()) > 0 {
		cfg.host = uri.Hostname()
		if len(uri.Port()) > 0 {
			cfg.port = uri.Port()
		}
	} else {
		cfg.unixSocketPath, err = url.PathUnescape(uri.EscapedPath())
	}

	// Parse Options
	var options url.Values
	if options, err = url.ParseQuery(uri.RawQuery); err != nil {
		return cfg, &SiodbDriverError{"Error while paring options from URI: '" + err.Error() + "'."}
	}

	if len(options.Get("identity_file")) > 0 {
		if cfg.privateKey, err = loadPrivateKey(options.Get("identity_file"), options.Get("identity_file_password")); err != nil {
			return cfg, err
		}
	}

	if len(options.Get("trace")) > 0 {
		if trc, err := strconv.ParseBool(options.Get("trace")); err == nil {
			cfg.trace = trc
		} else {
			return cfg, &SiodbDriverError{"Paring URI: option 'trace' can be 'true' or 'false'."}
		}
	}

	if cfg.trace {
		fmt.Printf("## SIODB DRIVER | config used: %v.\n", cfg)
	}

	return cfg, err
}

func (d SiodbDriver) Open(dsn string) (driver.Conn, error) {

	var cfg Config
	var err error

	if cfg, err = ParseURI(dsn); err != nil {
		return nil, err
	}
	c := &connector{
		cfg: cfg,
	}

	// New siodbConn
	sc := &siodbConn{
		cfg: c.cfg,
	}

	// Plain connection
	if sc.cfg.protocol == "siodbu" {
		if sc.netConn, err = net.Dial("unix", sc.cfg.unixSocketPath); err != nil {
			return nil, &SiodbDriverError{"Unable to connect to " + sc.cfg.unixSocketPath + "."}
		}
	}

	// Plain connection
	if sc.cfg.protocol == "siodb" {
		if sc.netConn, err = net.Dial("tcp", sc.cfg.host+":"+sc.cfg.port); err != nil {
			return nil, &SiodbDriverError{"Unable to connect to " + sc.cfg.host + "."}
		}
	}

	// TLS connection
	if sc.cfg.protocol == "siodbs" {
		config := &tls.Config{InsecureSkipVerify: true}

		if sc.netConn, err = tls.Dial("tcp", sc.cfg.host+":"+sc.cfg.port, config); err != nil {
			return nil, &SiodbDriverError{"Unable to connect to " + sc.cfg.host + "."}
		}
	}

	// Authentification
	if err := sc.authenticate(); err != nil {
		return sc, err
	}

	return sc, nil

}
