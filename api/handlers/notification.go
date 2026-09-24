package handlers

import (
	"fmt"
	"net/http"

	"practice/api/services"

	"github.com/gin-gonic/gin"
)

// StreamRideNotifications keeps an SSE connection open for ride-booked events.
func StreamRideNotifications(c *gin.Context) {
	customerID := c.Param("id")
	client := services.SubscribeRideNotifications(customerID)
	defer services.UnsubscribeRideNotifications(customerID, client)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	if _, err := fmt.Fprint(c.Writer, ": connected\n\n"); err != nil {
		return
	}
	if flusher, ok := c.Writer.(http.Flusher); ok {
		flusher.Flush()
	}

	for {
		select {
		case payload, ok := <-client:
			if !ok {
				return
			}
			if _, err := fmt.Fprintf(c.Writer, "event: ride_booked\ndata: %s\n\n", payload); err != nil {
				return
			}
			if flusher, ok := c.Writer.(http.Flusher); ok {
				flusher.Flush()
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}
