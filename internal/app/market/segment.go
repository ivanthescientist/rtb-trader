package market

import (
	"rtb-trader/internal/app/bidrequest"
	"strings"
)

// Segment represents an Ad market segment and its corresponding bidding strategy
type Segment struct {
	ID       string
	Matchers []SegmentMatcher
}

// Match checks if all matchers match to provided bidReuqest's FieldMap
func (s *Segment) Match(fieldMap bidrequest.FieldMap) bool {
	for _, matcher := range s.Matchers {
		if !matcher(fieldMap) {
			return false
		}
	}

	return true
}

func (s Segment) String() string {
	return "Segment[" + s.ID + "]"
}

// SegmentMatcher represents a single segment matching function, recommended to only target one single field if possible
type SegmentMatcher func(fieldMap bidrequest.FieldMap) bool

// CreateOSMatcher makes a strict value segment matcher
func CreateStrictOSMatcher(os string) SegmentMatcher {
	desiredValue := strings.ToLower(os)
	return func(fieldMap bidrequest.FieldMap) bool {
		var requestValue, isPresent = fieldMap[bidrequest.FieldOSName]
		if !isPresent {
			return false
		}

		return strings.ToLower(requestValue) == desiredValue
	}
}
