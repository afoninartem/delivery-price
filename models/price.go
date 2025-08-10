package models

import "time"

type PriceResponse struct {
	RedirectURL        string `json:"redirect_url"`
	OfferFailureBlocks struct {
		CargoIntercity struct {
			Blocks []struct {
				FailureReason string `json:"failure_reason"`
				Button        struct {
					Title    string `json:"title"`
					Subtitle string `json:"subtitle"`
					Style    string `json:"style"`
					Action   struct {
						Type              string `json:"type"`
						AdditionalContext struct {
							Vertical string `json:"vertical"`
						} `json:"additional_context"`
					} `json:"action"`
				} `json:"button"`
				Type string `json:"type"`
			} `json:"blocks"`
		} `json:"cargo_intercity"`
	} `json:"offer_failure_blocks"`
	ClaimsOffers []struct {
		Payload    string `json:"payload,omitempty"`
		TariffInfo struct {
			Title           string `json:"title"`
			Tariff          string `json:"tariff"`
			Vertical        string `json:"vertical"`
			SurgeLimit      string `json:"surge_limit"`
			TariffExtraInfo struct {
				TotalRouteTimeSeconds                  int    `json:"total_route_time_seconds"`
				SourcePointFreeWaitingTimeSeconds      int    `json:"source_point_free_waiting_time_seconds"`
				DestinationPointFreeWaitingTimeSeconds int    `json:"destination_point_free_waiting_time_seconds"`
				SourcePointWaitingPricePerMinute       string `json:"source_point_waiting_price_per_minute"`
				DestinationPointWaitingPricePerMinute  string `json:"destination_point_waiting_price_per_minute"`
			} `json:"tariff_extra_info"`
		} `json:"tariff_info"`
		Price struct {
			TotalPrice            string `json:"total_price"`
			TotalPriceWithVat     string `json:"total_price_with_vat"`
			StrikeoutPrice        string `json:"strikeout_price"`
			StrikeoutPriceWithVat string `json:"strikeout_price_with_vat"`
			Currency              string `json:"currency"`
			SurgeLevel            string `json:"surge_level"`
			SurgeRatio            string `json:"surge_ratio"`
		} `json:"price,omitempty"`
		VisitingIntervals []struct {
			Type string    `json:"type"`
			From time.Time `json:"from"`
			To   time.Time `json:"to"`
		} `json:"visiting_intervals,omitempty"`
		FailureReason struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"failure_reason,omitempty"`
	} `json:"claims_offers"`
}
