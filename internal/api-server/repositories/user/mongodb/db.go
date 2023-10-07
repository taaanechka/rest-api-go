package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"
	"github.com/taaanechka/rest-api-go/internal/apperror"
	"github.com/taaanechka/rest-api-go/pkg/client/mongodb"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

type DB struct {
	collection *mongo.Collection
	lg         *logging.Logger
}

func NewStorage(lg *logging.Logger, cfg userstorage.Config) (*DB, error) {
	ctx := context.Background()
	db, err := mongodb.NewClient(ctx, cfg.Host, cfg.Port,
		cfg.Username, cfg.Password, cfg.Database, cfg.AuthDB)
	if err != nil {
		return nil, fmt.Errorf("failed to init client: %w", err)
	}

	coll := db.Collection(cfg.Collection)

	for _, idx := range cfg.Indexes {
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{idx, 1}},
			Options: options.Index().SetUnique(true),
		}
		_, err = coll.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return nil, fmt.Errorf("failed to create index: %w", err)
		}
	}

	return &DB{
		collection: coll,
		lg:         lg,
	}, nil
}

func (d *DB) Create(ctx context.Context, user userstorage.User) (string, error) {
	d.lg.Debug("create user")
	res, err := d.collection.InsertOne(ctx, convertBLToDB(&user))
	if err != nil {
		return "", apperror.ErrCreate
	}

	d.lg.Debug("convert InsertedID to ObjectID")
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	d.lg.Trace(user)
	return "", fmt.Errorf("failed to convert objectid to hex. probably oid: %s", oid)
}

func (d *DB) FindAll(ctx context.Context) ([]userstorage.User, error) {
	cursor, err := d.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to read documents")
	}

	var us []User
	if err = cursor.All(ctx, &us); err != nil {
		return nil, apperror.ErrNotFound
	}

	u := make([]userstorage.User, 0, len(us))
	for _, cur := range us {
		u = append(u, convertDBToBL(&cur))
	}
	return u, nil
}

func (d *DB) FindOne(ctx context.Context, id string) (userstorage.User, error) {
	oid, ok := primitive.ObjectIDFromHex(id)
	if ok != nil {
		return userstorage.User{}, apperror.ErrBadID
	}

	filter := bson.M{"_id": oid}

	res := d.collection.FindOne(ctx, filter)
	if err := res.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = apperror.ErrNotFound
		}

		return userstorage.User{}, err
	}

	var uDB User
	if err := res.Decode(&uDB); err != nil {
		return userstorage.User{}, fmt.Errorf("failed to decode user: %w", err)
	}

	return convertDBToBL(&uDB), nil
}

func (d *DB) Update(ctx context.Context, id string, user userstorage.User) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrBadID
	}

	filter := bson.M{"_id": oid}

	uBytes, err := bson.Marshal(convertBLToDB(&user))
	if err != nil {
		return fmt.Errorf("failed to marshal user: %w", err)
	}

	var updUserObj bson.M
	err = bson.Unmarshal(uBytes, &updUserObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %w", err)
	}

	delete(updUserObj, "_id")

	upd := bson.M{
		"$set": updUserObj,
	}

	res, err := d.collection.UpdateOne(ctx, filter, upd)
	if err != nil {
		return apperror.ErrUpdate
	}

	if res.MatchedCount == 0 {
		return apperror.ErrNotFound
	}

	d.lg.Tracef("Matched %d documents and Modified %d documents", res.MatchedCount, res.ModifiedCount)

	return nil
}

func (d *DB) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return apperror.ErrBadID
	}

	filter := bson.M{"_id": oid}

	res, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return apperror.ErrDelete
	}
	if res.DeletedCount == 0 {
		return apperror.ErrNotFound
	}

	d.lg.Tracef("Deleted %d documents", res.DeletedCount)

	return nil
}
