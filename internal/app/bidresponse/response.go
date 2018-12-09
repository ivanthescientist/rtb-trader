package bidresponse

import "github.com/bsm/openrtb"

func CreateResponse(request *openrtb.BidRequest, bids []Bid) *openrtb.BidResponse {
	var response openrtb.BidResponse

	response.ID = request.ID

	currency := "USD"
	if len(request.Cur) > 0 {
		currency = request.Cur[0]
	}
	response.Currency = currency

	response.SeatBid = []openrtb.SeatBid{
		{
			Bid: convertBids(bids),
		},
	}

	if len(request.WSeat) > 0 {
		response.SeatBid[0].Seat = request.WSeat[0]
	}

	return &response
}

func convertBids(internalBids []Bid) []openrtb.Bid {
	var bids = make([]openrtb.Bid, len(internalBids))

	for i := range internalBids {
		bids[i].ID = internalBids[i].ID
		bids[i].ImpID = internalBids[i].ImpressionID
		bids[i].AdMarkup = internalBids[i].Content
		bids[i].Price = internalBids[i].Price
		bids[i].CampaignID = openrtb.StringOrNumber(internalBids[i].SegmentID)
	}

	return bids
}
