package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/btcsuite/btcd/rpcclient"
	"github.com/platanus/bitcoin-prometheus-exporter/client"
	"github.com/platanus/bitcoin-prometheus-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

func getEnvBool(key string, defaultValue bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Fatalf("Environment Variable value for %s must be a boolean", key)
	}
	return b
}

var (
	// Set during go build
	version   string
	gitCommit string

	// Defaults values
	defaultNamespace       = getEnv("NAMESPACE", "bitcoind")
	defaultListenAddress   = getEnv("LISTEN_ADDRESS", ":9113")
	defaultMetricsPath     = getEnv("TELEMETRY_PATH", "/metrics")
	defaultRPCHost         = getEnv("RPC_HOST", "localhost:8332")
	defaultRPCUser         = getEnv("RPC_USER", "")
	defaultRPCPass         = getEnv("RPC_PASS", "")
	defaultRPCHTTPPostMode = getEnvBool("RPC_POST_MODE", true)
	defaultRPCDisableTLS   = getEnvBool("RPC_DISABLE_TLS", true)

	// Command-line flags
	namespace = flag.String("namespace", defaultNamespace,
		"The namespace or prefix to use in the exported metrics. The default value can be overwritten by NAMESPACE environment variable.")
	listenAddr = flag.String("web.listen-address", defaultListenAddress,
		"An address to listen on for web interface and telemetry. The default value can be overwritten by LISTEN_ADDRESS environment variable.")
	metricsPath = flag.String("web.telemetry-path", defaultMetricsPath,
		"A path under which to expose metrics. The default value can be overwritten by TELEMETRY_PATH environment variable.")
	rpcHost = flag.String("rpc.host", defaultRPCHost,
		"Bitcoin node RPC host. The default value can be overwritten by RPC_HOST environment variable.")
	rpcUser = flag.String("rpc.user", defaultRPCUser,
		"Bitcoin node RPC username. The default value can be overwritten by RPC_USER environment variable.")
	rpcPass = flag.String("rpc.pass", defaultRPCPass,
		"Bitcoin node RPC password. The default value can be overwritten by RPC_PASS environment variable.")
	rpcHTTPPostMode = flag.Bool("rpc.post-mode", defaultRPCHTTPPostMode,
		"If the RPC requests should be only in post mode. The default value can be overwritten by RPC_POST_MODE environment variable.")
	rpcDisableTLS = flag.Bool("rpc.disable-tls", defaultRPCDisableTLS,
		"Disable TLS when connecting to the RPC api. The default value can be overwritten by RPC_DISABLE_TLS environment variable.")
)

func main() {
	flag.Parse()

	log.Printf("Starting Bitcoin Prometheus Exporter Version=%v GitCommit=%v", version, gitCommit)

	registry := prometheus.NewRegistry()

	connCfg := &rpcclient.ConnConfig{
		Host:         *rpcHost,
		User:         *rpcUser,
		Pass:         *rpcPass,
		HTTPPostMode: *rpcHTTPPostMode,
		DisableTLS:   *rpcDisableTLS,
	}
	rpcclient, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer rpcclient.Shutdown()

	client, err := client.NewBitcoinClient(rpcclient)
	if err != nil {
		log.Fatalf("Could not create Bitcoin Rpc Client: %v", err)
	}

	registry.MustRegister(collector.NewBitcoinCollector(client, *namespace))

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Bitcoin Exporter</title></head>
			<body>
			<h1>Bitcoin Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}
