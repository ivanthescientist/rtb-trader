package market

import (
	"github.com/bsm/openrtb"
	log "github.com/sirupsen/logrus"
	"rtb-trader/internal/app/bidrequest"
	"rtb-trader/internal/app/bidresponse"
	"sync"
)

type Trader struct {
	segmentsLock *sync.RWMutex
	segments     []Segment
	logger       *log.Logger
}

func NewTrader(segments []Segment, logger *log.Logger) *Trader {
	return &Trader{
		segments:     segments,
		segmentsLock: &sync.RWMutex{},
		logger:       logger,
	}
}

func (b *Trader) Bid(request *openrtb.BidRequest) *openrtb.BidResponse {
	b.segmentsLock.RLock()
	defer b.segmentsLock.RUnlock()

	var fieldMap, err = bidrequest.ExtractFields(request)
	if err != nil {
		return nil
	}

	b.logger.Info("Extracted fields from request: \n", fieldMap)

	var segment *Segment

	for _, s := range b.segments {
		if s.Match(fieldMap) {
			segment = &s
			break
		}
	}

	if segment == nil {
		return nil
	}

	b.logger.Infof("Found segment: %s", segment)

	bids := createBids(segment, request)
	if len(bids) == 0 {
		return nil
	}

	return bidresponse.CreateResponse(request, bids)
}

func (b *Trader) AddSegment(segment Segment) {
	b.segmentsLock.Lock()
	defer b.segmentsLock.Unlock()

	b.segments = append(b.segments, segment)
}

func createBids(segment *Segment, request *openrtb.BidRequest) []bidresponse.Bid {
	var bids = make([]bidresponse.Bid, len(request.Imp))

	for i, imp := range request.Imp {
		bids[i].ID = imp.ID
		bids[i].ImpressionID = imp.ID
		bids[i].Price = 1e-15
		bids[i].Content = ""
		bids[i].SegmentID = segment.ID
	}

	return bids
}
