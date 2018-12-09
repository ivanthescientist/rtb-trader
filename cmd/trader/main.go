package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"rtb-trader/internal/app/exchange"
	"rtb-trader/internal/app/market"
	"runtime"
	"syscall"
)

var logger *log.Logger

func init() {
	logger = log.New()
}

func main() {
	var err error
	var trader *market.Trader
	var exchangeHandler *exchange.Handler
	var server *http.Server

	trader = createTrader()
	exchangeHandler = exchange.NewHandler(trader, runtime.NumCPU(), logger)
	http.Handle("/", exchangeHandler)

	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	server = &http.Server{
		Addr:    "0.0.0.0:9000",
		Handler: exchangeHandler,
	}

	go func() {
		<-signals
		logger.Info("Shutting down...")

		err = server.Shutdown(context.Background())
		if err != nil {
			logger.Error(err)
		}
	}()

	err = server.ListenAndServe()
	if err != nil {
		logger.Error(err)
	}
}

func createTrader() *market.Trader {
	var trader *market.Trader
	var segments = []market.Segment{
		{
			ID:       "android_segment",
			Matchers: []market.SegmentMatcher{market.CreateStrictOSMatcher("android")},
		},
		{
			ID:       "ios_segment",
			Matchers: []market.SegmentMatcher{market.CreateStrictOSMatcher("ios")},
		},
		{
			ID:       "windows_segment",
			Matchers: []market.SegmentMatcher{market.CreateStrictOSMatcher("windows")},
		},
	}

	trader = market.NewTrader(segments, logger)

	return trader
}
