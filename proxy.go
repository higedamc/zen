package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elazarl/goproxy"
)

type Proxy struct{}

func NewProxy() *Proxy {
	return &Proxy{}
}

// ConfigureTLS configures the proxy to use the given certificate and key for
// TLS connections.
func (p *Proxy) ConfigureTLS(certFile, keyFile string) error {
	caCert, err := ioutil.ReadFile(certFile)
	if err != nil {
		return fmt.Errorf("failed to read CA certificate: %v", err)
	}
	caKey, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("failed to read CA key: %v", err)
	}
	goproxyCa, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return fmt.Errorf("failed to parse CA certificate and key: %v", err)
	}
	if goproxyCa.Leaf, err = x509.ParseCertificate(goproxyCa.Certificate[0]); err != nil {
		return fmt.Errorf("failed to parse CA certificate: %v", err)
	}

	goproxy.GoproxyCa = goproxyCa
	goproxy.OkConnect = &goproxy.ConnectAction{Action: goproxy.ConnectAccept, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.MitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.HTTPMitmConnect = &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	goproxy.RejectConnect = &goproxy.ConnectAction{Action: goproxy.ConnectReject, TLSConfig: goproxy.TLSConfigFromCA(&goproxyCa)}
	return nil
}

// Start starts the proxy on the given address.
func (p *Proxy) Start(addr string) error {
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	return http.ListenAndServe(addr, proxy)
}
