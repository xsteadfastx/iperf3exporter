//nolint:gochecknoglobals
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/VictoriaMetrics/metrics"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.xsfx.dev/logginghandler"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd is the command line interface.
var rootCmd = &cobra.Command{
	Use: "iperf3exporter",
	Run: func(cmd *cobra.Command, args []string) {
		if versionFlag {
			fmt.Printf("iperf3exporter %s, commit %s, %s", version, commit, date) //nolint:forbidigo

			return
		}

		http.Handle("/probe", logginghandler.Handler(http.HandlerFunc(probeHandler)))
		log.Info().Str("listen", c.Exporter.Listen).Msg("starting...")
		log.Fatal().Err(http.ListenAndServe(c.Exporter.Listen, nil)).Msg("goodbye")
	},
}

// config is a struct that a config file can get unmarshaled to.
type config struct {
	Exporter struct {
		Listen         string
		Timeout        time.Duration
		ProcessMetrics bool `mapstructure:"process_metrics"`
	}
	Log struct {
		JSON   bool
		Colors bool
	}
	Iperf3 struct {
		Time int
	}
}

// c is a global config struct instance.
var c config

var (
	cfgFile     string
	versionFlag bool
)

var (
	downloadSentBitsPerSecond     = metrics.NewFloatCounter("iperf3_download_sent_bits_per_second")
	downloadSentSeconds           = metrics.NewFloatCounter("iperf3_download_sent_seconds")
	downloadSentBytes             = metrics.NewFloatCounter("iperf3_download_sent_bytes")
	downloadReceivedBitsPerSecond = metrics.NewFloatCounter("iperf3_download_received_bits_per_second")
	downloadReceivedSeconds       = metrics.NewFloatCounter("iperf3_download_received_seconds")
	downloadReceivedBytes         = metrics.NewFloatCounter("iperf3_download_received_bytes")
)

var (
	uploadSentBitsPerSecond     = metrics.NewFloatCounter("iperf3_upload_sent_bits_per_second")
	uploadSentSeconds           = metrics.NewFloatCounter("iperf3_upload_sent_seconds")
	uploadSentBytes             = metrics.NewFloatCounter("iperf3_upload_sent_bytes")
	uploadReceivedBitsPerSecond = metrics.NewFloatCounter("iperf3_upload_received_bits_per_second")
	uploadReceivedSeconds       = metrics.NewFloatCounter("iperf3_upload_received_seconds")
	uploadReceivedBytes         = metrics.NewFloatCounter("iperf3_upload_received_bytes")
)

//nolint:tagliatelle
type iperfResult struct {
	End struct {
		SumSent struct {
			Seconds       float64 `json:"seconds"`
			Bytes         float64 `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
		} `json:"sum_sent"`
		SumReceived struct {
			Seconds       float64 `json:"seconds"`
			Bytes         float64 `json:"bytes"`
			BitsPerSecond float64 `json:"bits_per_second"`
		} `json:"sum_received"`
	} `json:"end"`
}

func runIperf(ctx context.Context, target string, cmdArgs []string, logger zerolog.Logger) (iperfResult, error) {
	args := []string{
		"-J",
		"-c",
		target,
		"-t", strconv.Itoa(c.Iperf3.Time),
	}

	args = append(args, cmdArgs...)

	cmd := exec.CommandContext(ctx, "iperf3", args...)

	logger.Debug().Str("cmd", cmd.String()).Msg("created command")

	// Buffers to store stdout and stderr.
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		logger.Debug().
			Str("stdout", outb.String()).
			Str("stderr", errb.String()).
			Msg("output from failed run")

		return iperfResult{}, fmt.Errorf("could not run command: %w", err)
	}

	// Unmarshal the output to the iperf struct.
	var p iperfResult
	if err := json.Unmarshal(outb.Bytes(), &p); err != nil {
		return iperfResult{}, fmt.Errorf("could not unmarshal result: %w", err)
	}

	return p, nil
}

func download(ctx context.Context, target string, logger zerolog.Logger) error {
	r, err := runIperf(
		ctx,
		target,
		[]string{"-R"},
		logger,
	)
	if err != nil {
		return fmt.Errorf("could not get download metrics: %w", err)
	}

	downloadSentBitsPerSecond.Set(r.End.SumSent.BitsPerSecond)
	downloadSentBytes.Set(r.End.SumSent.Bytes)
	downloadSentSeconds.Set(r.End.SumSent.Seconds)

	downloadReceivedBitsPerSecond.Set(r.End.SumReceived.BitsPerSecond)
	downloadReceivedBytes.Set(r.End.SumReceived.Bytes)
	downloadReceivedSeconds.Set(r.End.SumReceived.Seconds)

	return nil
}

func upload(ctx context.Context, target string, logger zerolog.Logger) error {
	r, err := runIperf(
		ctx,
		target,
		[]string{},
		logger,
	)
	if err != nil {
		return fmt.Errorf("could not get upload metrics: %w", err)
	}

	uploadSentBitsPerSecond.Set(r.End.SumSent.BitsPerSecond)
	uploadSentBytes.Set(r.End.SumSent.Bytes)
	uploadSentSeconds.Set(r.End.SumSent.Seconds)

	uploadReceivedBitsPerSecond.Set(r.End.SumReceived.BitsPerSecond)
	uploadReceivedBytes.Set(r.End.SumReceived.Bytes)
	uploadReceivedSeconds.Set(r.End.SumReceived.Seconds)

	return nil
}

func probeHandler(w http.ResponseWriter, r *http.Request) {
	logger := logginghandler.Logger(r)

	// Extract target.
	target := r.URL.Query().Get("target")

	if target == "" {
		logger.Error().Msg("could not find target in url params")
		http.Error(w, "could not find target in url params", http.StatusUnprocessableEntity)

		return
	}

	logger.Debug().Str("target", target).Msg("extracted target from url params")

	ctx, cancel := context.WithTimeout(context.Background(), c.Exporter.Timeout)
	defer cancel()

	logger.Info().Msg("getting download metrics")

	if err := download(ctx, target, logger); err != nil {
		logger.Error().Err(err).Msg("could not create download metrics")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	logger.Info().Msg("getting upload metrics")

	if err := upload(ctx, target, logger); err != nil {
		logger.Error().Err(err).Msg("could not create upload metrics")
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	logger.Info().Msg("done scraping")

	metrics.WritePrometheus(w, c.Exporter.ProcessMetrics)
}

func init() { //nolint:gochecknoinits,funlen
	cobra.OnInitialize(initConfig)

	// Version.
	rootCmd.PersistentFlags().BoolVarP(&versionFlag, "version", "v", false, "print version")

	// Config.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	// Exporter.Listen.
	rootCmd.PersistentFlags().String("listen", "127.0.0.1:9119", "listen string")

	if err := viper.BindPFlag("exporter.listen", rootCmd.PersistentFlags().Lookup("listen")); err != nil {
		log.Fatal().Err(err).Msg("could not bind flag")
	}

	viper.SetDefault("exporter.listen", "127.0.0.1:9119")

	// Exporter.Timeout.
	rootCmd.PersistentFlags().Duration("timeout", time.Minute, "scraping timeout")

	if err := viper.BindPFlag("exporter.timeout", rootCmd.PersistentFlags().Lookup("timeout")); err != nil {
		log.Fatal().Err(err).Msg("could not bind flag")
	}

	viper.SetDefault("exporter.timeout", time.Minute)

	// Exporter.ProcessMetrics.
	rootCmd.PersistentFlags().Bool("process-metrics", true, "exporter process metrics")

	if err := viper.BindPFlag("exporter.processmetrics", rootCmd.PersistentFlags().Lookup("process-metrics")); err != nil {
		log.Fatal().Err(err).Msg("could not bind flag")
	}

	viper.SetDefault("exporter.processmetrics", true)

	// Log.JSON.
	rootCmd.PersistentFlags().Bool("log-json", false, "JSON log output")

	if err := viper.BindPFlag("log.json", rootCmd.PersistentFlags().Lookup("log-json")); err != nil {
		log.Fatal().Err(err).Msg("could not bind flag")
	}

	viper.SetDefault("log.json", false)

	// Log.Colors.
	rootCmd.PersistentFlags().Bool("log-colors", true, "colorful log output")

	if err := viper.BindPFlag("log.colors", rootCmd.PersistentFlags().Lookup("log-colors")); err != nil {
		log.Fatal().Err(err).Msg("could not bind flag")
	}

	viper.SetDefault("log.colors", true)

	// Iperf3.Time.
	rootCmd.PersistentFlags().Int("time", 5, "time in seconds to transmit for") //nolint:gomnd

	if err := viper.BindPFlag("iperf3.time", rootCmd.PersistentFlags().Lookup("time")); err != nil {
		log.Fatal().Err(err).Msg("could not bind flag")
	}

	viper.SetDefault("iperf3.time", 5) //nolint:gomnd
}

func initConfig() {
	// Load config file.
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)

		if err := viper.ReadInConfig(); err != nil {
			log.Fatal().Err(err).Msg("error on reading config")
		}
	}

	// Environment variable handling.
	viper.SetEnvPrefix("IPERF3EXPORTER")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Unmarshal config.
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatal().Err(err).Msg("could not unmarshal config")
	}

	// Setting up logging.
	switch c.Log.JSON {
	case false:
		if c.Log.Colors {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		} else {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: true})
		}
	case true:
		log.Logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	}

	log.Logger = log.With().Caller().Logger()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("goodbye")
	}
}
