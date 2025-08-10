package models

type UserState struct {
	Step     string
	Location *Location
}

type LastPrices map[uint]string // map[location_id]last_seen_price
