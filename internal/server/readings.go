package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type buoySnapshot struct {
	ID       string `json:"id"`
	BuoyID   string `json:"buoyId"`
	Status   string `json:"status"`
	GPS      string `json:"gps"`
	Compass  string `json:"compass"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type buoyState struct {
	Snapshot  buoySnapshot
	Latitude  float64
	Longitude float64
}

func (s *Server) BuoysHandler(w http.ResponseWriter, r *http.Request) {
	buoys := []buoySnapshot{}
	readings, err := s.db.GetLastKnownPosition(r.Context())
	if err == nil {
		for _, reading := range readings {
			imageURL := reading.ImageURL
			if imageURL == "" {
				latestImageURL, storageErr := s.storage.GetLatestImageURL(r.Context(), reading.BuoyID)
				if storageErr == nil {
					imageURL = latestImageURL
				}
			}

			buoys = append(buoys, buoySnapshot{
				ID:       reading.BuoyID,
				BuoyID:   reading.BuoyID,
				Status:   "online",
				GPS:      fmt.Sprintf("%f,%f", reading.Latitude, reading.Longitude),
				Compass:  "waiting",
				ImageURL: imageURL,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(struct {
		Buoys []buoySnapshot `json:"buoys"`
	}{
		Buoys: buoys,
	})
}
