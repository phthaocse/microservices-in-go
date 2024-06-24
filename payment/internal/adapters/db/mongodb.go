package db

import (
	"context"
	"github.com/huseyinbabal/microservices/payment/internal/application/core/domain"
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
	log.Info(config.Address)
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

type PaymentMongo struct {
	ID         int64     `bson:"_id"`
	CustomerID int64     `bson:"customer_id"`
	Status     string    `bson:"status"`
	OrderID    int64     `bson:"order_id"`
	TotalPrice float32   `bson:"total_price"`
	CreatedAt  time.Time `bson:"created_at"`
}

func (a *MongoAdapter) Get(ctx context.Context, id string) (domain.Payment, error) {
	var paymentEntity PaymentMongo
	err := a.db.Collection("payment").FindOne(ctx, bson.M{"_id": id}).Decode(&paymentEntity)
	if err != nil {
		return domain.Payment{}, err
	}
	payment := domain.Payment{
		ID:         int64(paymentEntity.ID),
		CustomerID: paymentEntity.CustomerID,
		Status:     paymentEntity.Status,
		OrderId:    paymentEntity.OrderID,
		TotalPrice: paymentEntity.TotalPrice,
		CreatedAt:  paymentEntity.CreatedAt.UnixNano(),
	}
	return payment, nil
}

func (a *MongoAdapter) Save(ctx context.Context, payment *domain.Payment) error {
	paymentModel := PaymentMongo{
		ID:         time.Now().UnixNano(),
		CustomerID: payment.CustomerID,
		Status:     payment.Status,
		OrderID:    payment.OrderId,
		TotalPrice: payment.TotalPrice,
	}

	_, err := a.db.Collection("payment").InsertOne(ctx, &paymentModel)

	return err
}
