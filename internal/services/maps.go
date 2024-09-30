package services

import (
	"context"
	"ecampus/config"
	"github.com/gin-gonic/gin"
	"googlemaps.github.io/maps"
	"net/http"
)

type RouteInfo struct {
	Distance int      `json:"distance"`
	Duration string   `json:"duration"`
	Steps    []string `json:"steps"`
}

func GetRouteInfo(c *gin.Context) {
	origin := c.Query("origin")
	destination := c.Query("destination")

	if origin == "" || destination == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Origin and destination are required"})
		return
	}

	// Retrieve configuration and handle potential errors
	cfg, err := config.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Configuration error: " + err.Error()})
		return
	}

	// Validate that Google Maps API Key is present
	if cfg.GoogleMapsAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Google Maps API key not set"})
		return
	}

	// Create Google Maps client
	client, err := maps.NewClient(maps.WithAPIKey(cfg.GoogleMapsAPIKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Google Maps client"})
		return
	}

	// Prepare directions request
	r := &maps.DirectionsRequest{
		Origin:      origin,
		Destination: destination,
	}

	// Get directions
	route, _, err := client.Directions(context.Background(), r)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get directions: " + err.Error()})
		return
	}

	// Process the route information
	if len(route) > 0 {
		leg := route[0].Legs[0]
		steps := make([]string, len(leg.Steps))
		for i, step := range leg.Steps {
			steps[i] = step.HTMLInstructions
		}

		routeInfo := RouteInfo{
			Distance: int(leg.Distance.Meters),
			Duration: leg.Duration.String(),
			Steps:    steps,
		}

		c.JSON(http.StatusOK, routeInfo)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "No route found"})
	}
}
