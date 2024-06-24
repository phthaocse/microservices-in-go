package db

import (
	"context"
	"github.com/huseyinbabal/microservices/order/internal/application/core/domain"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoAdapter struct {
	db *mongo.Database
}

type MongoConfig struct {
	Name       string        `yaml:"name"`
	Address    []string      `yaml:"address"`
	RepSetName *string       `yaml:"repset_name"`
	DBAuthen   string        `yaml:"dbauthen"`
	User       string        `yaml:"user"`
	Pass       string        `yaml:"pass"`
	Timeout    time.Duration `yaml:"timeout"`
	Database   string        `yaml:"dbname"`
}

func NewMongoAdapter(config *MongoConfig) (*MongoAdapter, error) {

	clientOption := &options.ClientOptions{
		Hosts:   config.Address,
		Timeout: &config.Timeout,
	}

	// set monitoring

	// set Authen
	if config.User != "" && config.Pass != "" {
		clientAuth := options.Credential{
			AuthSource: config.DBAuthen,
			Username:   config.User,
			Password:   config.Pass,
		}

		clientOption.SetAuth(clientAuth)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Info("connect to mongodb\n")
	log.Infof("mongo config: %v\n", config)
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	log.Info("Connected to MongoDB")

	return &MongoAdapter{db: client.Database(config.Database)}, nil
}

type OrderMongo struct {
	ID         int64            `bson:"_id"`
	CustomerID int64            `bson:"customer_id"`
	Status     string           `bson:"status"`
	OrderItems []OrderItemMongo `bson:"order_items"`

	CreatedAt time.Time `bson:"created_at"`
}

type OrderItemMongo struct {
	ProductCode string  `bson:"product_code"`
	UnitPrice   float32 `bson:"unit_price"`
	Quantity    int32   `bson:"quantity"`
	OrderID     uint    `bson:"order_id"`
}

func (a *MongoAdapter) Get(ctx context.Context, id int64) (domain.Order, error) {
	var orderEntity OrderMongo
	var orderItems []domain.OrderItem
	err := a.db.Collection("order").FindOne(ctx, bson.M{"_id": id}).Decode(&orderEntity)
	if err != nil {
		return domain.Order{}, err
	}

	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}
	return order, nil
}

func (a *MongoAdapter) Save(ctx context.Context, order *domain.Order) error {
	var orderItems []OrderItemMongo
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItemMongo{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}
	orderModel := OrderMongo{
		ID:         time.Now().UnixNano(),
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,
	}

	_, err := a.db.Collection("order").InsertOne(ctx, &orderModel)

	return err
}
