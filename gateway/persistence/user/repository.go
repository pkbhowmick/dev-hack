package user

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"

	sqlcdb "github.com/pkbhowmick/dev-hack/gateway/infra/sqlc/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
	DB *sqlcdb.Queries
}

type UserMongo struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

// MongoDB setup
var mongoClient *mongo.Client
var mongoCtx context.Context

func initMongo() {
	var err error
	mongoCtx = context.Background()
	// Update the URI with your MongoDB credentials
	mongoURI := "mongodb://mongouser:mongopassword@localhost:27017"
	mongoClient, err = mongo.Connect(mongoCtx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = mongoClient.Ping(mongoCtx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

func (r *Repository) Save(ctx context.Context, name, email, password string) error {
	arg := sqlcdb.CreateUserParams{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: password,
	}

	err := r.DB.CreateUser(ctx, arg)
	if err != nil {
		return err
	}

	// store the data with name, email to mongodb
	initMongo()

	// Now, store the name and email in MongoDB
	userMongo := UserMongo{Name: name, Email: email}
	collection := mongoClient.Database("mydb").Collection("users")
	_, mongoErr := collection.InsertOne(ctx, userMongo)
	if mongoErr != nil {
		// Handle error (you might also want to delete the SQL entry to keep data consistent)
		return mongoErr
	}

	return nil
}

func (r *Repository) GetByName(ctx context.Context, name string) (sqlcdb.User, error) {
	return r.DB.GetUserByName(ctx, name)
}

func (r *Repository) GetByID(ctx context.Context, id string) (*sqlcdb.User, error) {
	user, err := r.DB.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
