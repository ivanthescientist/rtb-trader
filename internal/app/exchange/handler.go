package exchange

import (
	"encoding/json"
	"github.com/bsm/openrtb"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"rtb-trader/internal/app/market"
)

type Handler struct {
	trader   *market.Trader
	logger   *log.Logger
	lockPool chan struct{}
}

func NewHandler(trader *market.Trader, workerCount int, logger *log.Logger) *Handler {
	lockPool := make(chan struct{}, workerCount)
	for i := 0; i < workerCount; i++ {
		lockPool <- struct{}{}
	}

	return &Handler{
		trader:   trader,
		logger:   logger,
		lockPool: lockPool,
	}
}

func (h *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		resp.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	select {
	case <-h.lockPool:
		defer func() { h.lockPool <- struct{}{} }()
	default:
		log.Info("Skipping responding to request due to worker pool exhaustion")
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	var requestBytes []byte
	var err error
	var bidRequest openrtb.BidRequest

	requestBytes, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Errorf("Failed to read request: %s", err)
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	err = json.Unmarshal(requestBytes, &bidRequest)
	if err != nil {
		log.Errorf("Failed to deserialize response: %s", err)
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	bidResponse := h.trader.Bid(&bidRequest)
	if bidResponse == nil {
		log.Info("No bids placed by trader")
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	var responseBytes []byte
	responseBytes, err = json.Marshal(bidResponse)
	if err != nil {
		log.Errorf("Failed to marshall bid response: %s", err)
		resp.WriteHeader(http.StatusNoContent)
		return
	}

	var written int
	written, err = resp.Write(responseBytes)
	if written != len(responseBytes) || err != nil {
		log.Errorf("Failed to write response: ")
	}
}
