package models

import "go.mongodb.org/mongo-driver/v2/bson"

// FlowerShop represents a shop that sells flowers. Flowers holds the array of
// flowers available in the shop.
type FlowerShop struct {
	ShopID   bson.ObjectID `json:"shopid,omitempty" bson:"shopid,omitempty"`
	Name     string        `json:"name" bson:"name"`
	Owner    string        `json:"owner" bson:"owner"`
	Location string        `json:"location" bson:"location"`
	Phone    string        `json:"phone" bson:"phone"`
	Flowers  []Flower      `json:"flowers" bson:"flowers"`
}

// FlowerShopRequest is the payload accepted when creating or updating a shop.
// The flowers array is managed through the add/remove flower endpoints.
type FlowerShopRequest struct {
	Name     string `json:"name" binding:"required"`
	Owner    string `json:"owner" binding:"required"`
	Location string `json:"location" binding:"required"`
	Phone    string `json:"phone"`
}

// FlowerShopFlowerRequest is the payload used to add a flower to a shop.
type FlowerShopFlowerRequest struct {
	FlowerID string `json:"flowerid" binding:"required"`
}
