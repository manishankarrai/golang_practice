package services

import (
	"context"
	"errors"
	"time"

	"practice/api/db"
	"practice/api/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const shopOrderCollection = "shoporders"

var (
	ErrShopOrderNotFound = errors.New("shop order not found")
	ErrInvalidShopOrder = errors.New("invalid shop order")
)

func shopOrderColl() *mongo.Collection { return db.Database("golang_practice").Collection(shopOrderCollection) }

func shopOrderID(id string) (bson.ObjectID, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil { return bson.NilObjectID, ErrShopOrderNotFound }
	return objectID, nil
}

func buildShopOrder(ctx context.Context, req models.ShopOrderRequest) (models.ShopOrder, error) {
	buyer, err := GetFlowerBuyerByID(ctx, req.BuyerID)
	if err != nil { return models.ShopOrder{}, err }
	shop, err := GetFlowerShopByID(ctx, req.ShopID)
	if err != nil { return models.ShopOrder{}, err }
	_ = buyer

	items := make([]models.ShopOrderItem, 0, len(req.Items))
	total := 0.0
	for _, requested := range req.Items {
		flowerID, err := bson.ObjectIDFromHex(requested.FlowerID)
		if err != nil { return models.ShopOrder{}, ErrInvalidShopOrder }
		var flower models.Flower
		found := false
		for _, available := range shop.Flowers {
			if available.FlowerID == flowerID { flower, found = available, true; break }
		}
		if !found { return models.ShopOrder{}, ErrInvalidShopOrder }
		subtotal := flower.Pricing * float64(requested.Quantity)
		items = append(items, models.ShopOrderItem{FlowerID: flower.FlowerID, Name: flower.Name, Quantity: requested.Quantity, Price: flower.Pricing, Subtotal: subtotal})
		total += subtotal
	}
	status := req.Status
	if status == "" { status = "pending" }
	now := time.Now().UTC()
	return models.ShopOrder{BuyerID: buyer.BuyerID, ShopID: shop.ShopID, Items: items, Total: total, Status: status, CreatedAt: now, UpdatedAt: now}, nil
}

func CreateShopOrder(ctx context.Context, req models.ShopOrderRequest) (models.ShopOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	order, err := buildShopOrder(ctx, req)
	if err != nil { return models.ShopOrder{}, err }
	order.OrderID = bson.NewObjectID()
	if _, err := shopOrderColl().InsertOne(ctx, order); err != nil { return models.ShopOrder{}, err }
	return order, nil
}

func GetShopOrders(ctx context.Context) ([]models.ShopOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cursor, err := shopOrderColl().Find(ctx, bson.M{})
	if err != nil { return nil, err }
	defer cursor.Close(ctx)
	orders := make([]models.ShopOrder, 0)
	if err := cursor.All(ctx, &orders); err != nil { return nil, err }
	return orders, nil
}

func GetShopOrderByID(ctx context.Context, id string) (models.ShopOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	orderID, err := shopOrderID(id)
	if err != nil { return models.ShopOrder{}, err }
	var order models.ShopOrder
	if err := shopOrderColl().FindOne(ctx, bson.M{"orderid": orderID}).Decode(&order); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return models.ShopOrder{}, ErrShopOrderNotFound }
		return models.ShopOrder{}, err
	}
	return order, nil
}

func UpdateShopOrder(ctx context.Context, id string, req models.ShopOrderRequest) (models.ShopOrder, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	orderID, err := shopOrderID(id)
	if err != nil { return models.ShopOrder{}, err }
	order, err := buildShopOrder(ctx, req)
	if err != nil { return models.ShopOrder{}, err }
	order.OrderID = orderID
	var updated models.ShopOrder
	err = shopOrderColl().FindOneAndUpdate(ctx, bson.M{"orderid": orderID}, bson.M{"$set": bson.M{"buyerid": order.BuyerID, "shopid": order.ShopID, "items": order.Items, "total": order.Total, "status": order.Status, "updated_at": time.Now().UTC()}}, options.FindOneAndUpdate().SetReturnDocument(options.After)).Decode(&updated)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) { return models.ShopOrder{}, ErrShopOrderNotFound }
		return models.ShopOrder{}, err
	}
	return updated, nil
}

func DeleteShopOrder(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	orderID, err := shopOrderID(id)
	if err != nil { return err }
	result, err := shopOrderColl().DeleteOne(ctx, bson.M{"orderid": orderID})
	if err != nil { return err }
	if result.DeletedCount == 0 { return ErrShopOrderNotFound }
	return nil
}
