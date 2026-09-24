package services

import (
	"encoding/json"
	"sync"
)

// RideBookedNotification is delivered to a customer's active SSE connection.
type RideBookedNotification struct {
	Type string      `json:"type"`
	Ride interface{} `json:"ride"`
}

type notificationBroker struct {
	mu      sync.RWMutex
	clients map[string]map[chan []byte]struct{}
}

var rideNotifications = notificationBroker{
	clients: make(map[string]map[chan []byte]struct{}),
}

func (b *notificationBroker) subscribe(customerID string) chan []byte {
	b.mu.Lock()
	defer b.mu.Unlock()

	client := make(chan []byte, 1)
	if b.clients[customerID] == nil {
		b.clients[customerID] = make(map[chan []byte]struct{})
	}
	b.clients[customerID][client] = struct{}{}
	return client
}

func (b *notificationBroker) unsubscribe(customerID string, client chan []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if clients := b.clients[customerID]; clients != nil {
		if _, exists := clients[client]; exists {
			delete(clients, client)
			close(client)
		}
		if len(clients) == 0 {
			delete(b.clients, customerID)
		}
	}
}

func (b *notificationBroker) publish(customerID string, notification RideBookedNotification) {
	payload, err := json.Marshal(notification)
	if err != nil {
		return
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	for client := range b.clients[customerID] {
		select {
		case client <- payload:
		default:
			// Do not let a slow frontend block ride booking.
		}
	}
}

func SubscribeRideNotifications(customerID string) chan []byte {
	return rideNotifications.subscribe(customerID)
}

func UnsubscribeRideNotifications(customerID string, client chan []byte) {
	rideNotifications.unsubscribe(customerID, client)
}

func PublishRideBooked(customerID string, ride interface{}) {
	rideNotifications.publish(customerID, RideBookedNotification{
		Type: "ride_booked",
		Ride: ride,
	})
}
