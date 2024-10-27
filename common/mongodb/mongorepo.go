package mongodb

import (
	"build-service-gin/common/logger"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrContextNotFoundKeyRegion = errors.New("mongo multi conn: context not found key region")
	ErrNotFoundRegion           = errors.New("mongo multi conn: mapping collections not found region")
)

// ModelInterface defines methods for models metadata, like collection name and indexes.
type ModelInterface interface {
	CollectionName() string
	IndexModels() []mongo.IndexModel
}

// Repository manages MongoDB collection and connection mapping for models.
type Repository[T ModelInterface] struct {
	*mongo.Collection
	*FilterPlayer
	mappingCollections map[string]*mongo.Collection
	err                error
}

// NewRepository initializes a new repository with MongoDB connections and collections.
func NewRepository[T ModelInterface](dbStorage *DatabaseStorage, opts ...*options.CollectionOptions) *Repository[T] {
	log := logger.GetLogger()

	var t T
	collectionName := t.CollectionName()
	indexModels := t.IndexModels()

	// Single DB connection
	if dbStorage.db != nil {
		collection, err := newRepository(dbStorage.db, collectionName, indexModels, opts...)
		if err != nil {
			log.Fatal().Msgf("new repository error: %v", err)
		}

		return &Repository[T]{
			Collection: collection,
		}
	}

	// Multi-connection setup
	connNames := t.CollectionName()
	if len(connNames) == 0 {
		log.Fatal().Msgf("mongo multi conn: connNames not found for collectionName=%s", collectionName)
	}

	mappingDatabases := make(map[string]*mongo.Database)
	mappingNames := make(map[string]string)

	for _, connName := range connNames {
		splitConnName := strings.Split(string(connName), "::")
		if len(splitConnName) != 2 {
			log.Fatal().Msgf("mongo multi conn: connName=%s invalid", string(connName))
		}

		db, ok := dbStorage.mappingDB[string(connName)]
		if !ok {
			log.Fatal().Msgf("mongo multi conn: connName=%s not found in mappingDB", string(connName))
		}

		region := splitConnName[0]
		name := splitConnName[1]

		mappingDatabases[region] = db
		mappingNames[region] = name
	}

	mappingCollections := make(map[string]*mongo.Collection)
	for region, db := range mappingDatabases {
		collection, err := newRepository(db, collectionName, indexModels, opts...)
		if err != nil {
			connName := fmt.Sprintf("%s::%s", region, mappingNames[region])
			log.Fatal().Msgf("mongo multi conn - connName=%s: new repository error: %v", connName, err)
		}

		mappingCollections[region] = collection
	}

	return &Repository[T]{
		mappingCollections: mappingCollections,
	}
}

// newRepository creates a collection and sets up indexes if necessary.
func newRepository(db *mongo.Database, collectionName string, indexModels []mongo.IndexModel, opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	collection := db.Collection(collectionName, opts...)

	if len(indexModels) > 0 {
		ctx, cc := context.WithTimeout(context.Background(), 30*time.Second)
		defer cc()
		_, err := collection.Indexes().CreateMany(ctx, indexModels)
		if err != nil {
			return nil, fmt.Errorf("create index collectionName=%v error: %v", collectionName, err)
		}
	}

	return collection, nil
}

// NewFilterPlayer initializes a new FilterPlayer in the repository.
func (r *Repository[T]) NewFilterPlayer() *Repository[T] {
	return &Repository[T]{
		Collection:   r.Collection,
		FilterPlayer: NewFilterPlayer(),
		err:          r.err,
	}
}

// FindOneDoc finds and returns a single document.
func (r *Repository[T]) FindOneDoc(ctx context.Context) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	opts := r.optsFindOne
	if len(r.sortOne) > 0 {
		opts.Sort = r.sortOne
	}

	var m T
	err := r.Collection.FindOne(ctx, r.filter, &opts).Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

// FindDocs finds multiple documents based on the filter.
func (r *Repository[T]) FindDocs(ctx context.Context) ([]*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	opts := r.optsFind
	if len(r.sort) > 0 {
		opts.Sort = r.sort
	}

	cs, err := r.Collection.Find(ctx, r.filter, &opts)
	if err != nil {
		return nil, err
	}

	ms := make([]*T, 0)
	err = cs.All(ctx, &ms)
	if err != nil {
		return nil, err
	}

	return ms, nil
}

// CreateOneDocument inserts a new document into the collection.
func (r *Repository[T]) CreateOneDocument(ctx context.Context, document *T) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	t := time.Now()
	doc, err := r.convertToBson(document)
	if err != nil {
		return nil, err
	}

	doc["created_at"] = &t
	doc["updated_at"] = &t
	result, err := r.Collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	doc["_id"] = result.InsertedID
	entity, _ := r.convertToObject(doc)
	return entity, nil
}

// CreateManyDocs inserts multiple documents into the collection.
func (r *Repository[T]) CreateManyDocs(ctx context.Context, documents []*T) ([]*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	t := time.Now()
	var docsProcessed []interface{}
	for _, document := range documents {
		docP, err := r.convertToBson(document)
		if err != nil {
			return nil, err
		}

		docP["created_at"] = &t
		docP["updated_at"] = &t
		docsProcessed = append(docsProcessed, docP)
	}

	result, err := r.Collection.InsertMany(ctx, docsProcessed)
	if err != nil {
		return nil, err
	}

	var entities []*T
	for i, doc := range docsProcessed {
		doc.(bson.M)["_id"] = result.InsertedIDs[i]
		entity, _ := r.convertToObject(doc.(bson.M))
		entities = append(entities, entity)
	}

	return entities, nil
}

// Helper methods for converting between Go objects and BSON
func (r *Repository[T]) convertToObject(b bson.M) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	if b == nil {
		return nil, nil
	}

	bytes, err := bson.Marshal(b)
	if err != nil {
		return nil, err
	}

	var doc T
	if err := bson.Unmarshal(bytes, &doc); err != nil {
		return nil, err
	}

	return &doc, nil
}

func (r *Repository[T]) convertToBson(document *T) (bson.M, error) {
	if r.err != nil {
		return nil, r.err
	}

	if document == nil {
		return nil, nil
	}

	bytes, err := bson.Marshal(document)
	if err != nil {
		return nil, err
	}

	var doc bson.M
	if err := bson.Unmarshal(bytes, &doc); err != nil {
		return nil, err
	}

	return doc, nil
}
func (r *Repository[T]) UpdateOneDoc(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.Collection.UpdateOne(ctx, r.filter, update, opts...)
}

func (r *Repository[T]) UpsertDoc(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	optUpsert := options.Update().SetUpsert(true)
	opts = append(opts, optUpsert)

	return r.Collection.UpdateOne(ctx, r.filter, update, opts...)
}

func (r *Repository[T]) UpdateManyDocs(ctx context.Context, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.Collection.UpdateMany(ctx, r.filter, update, opts...)
}

func (r *Repository[T]) FindOneAndUpdateDoc(ctx context.Context, update interface{}, opts ...*options.FindOneAndUpdateOptions) (*T, error) {
	if r.err != nil {
		return nil, r.err
	}

	res := r.Collection.FindOneAndUpdate(ctx, r.filter, update, opts...)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var m T
	err := res.Decode(&m)
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *Repository[T]) CountDocs(ctx context.Context, opts ...*options.CountOptions) (int64, error) {
	if r.err != nil {
		return 0, r.err
	}

	return r.Collection.CountDocuments(ctx, r.filter, opts...)
}

func (r *Repository[T]) DeleteOneDoc(ctx context.Context, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.Collection.DeleteOne(ctx, r.filter, opts...)
}

func (r *Repository[T]) DeleteManyDocs(ctx context.Context, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.Collection.DeleteMany(ctx, r.filter, opts...)
}

func (r *Repository[T]) DistinctDocs(ctx context.Context, fieldName string, opts ...*options.DistinctOptions) ([]interface{}, error) {
	if r.err != nil {
		return nil, r.err
	}

	return r.Collection.Distinct(ctx, fieldName, r.filter, opts...)
}

func (r *Repository[T]) SetLimit(limit int64) *Repository[T] {
	r.optsFind.Limit = &limit
	return r
}

func (r *Repository[T]) SetSkip(skip int64) *Repository[T] {
	r.optsFind.Skip = &skip
	return r
}

func (r *Repository[T]) SetSkipOne(skip int64) *Repository[T] {
	r.optsFindOne.Skip = &skip
	return r
}

func (r *Repository[T]) SetProjection(projection bson.M) *Repository[T] {
	r.optsFind.Projection = projection
	return r
}

func (r *Repository[T]) SetProjectionOne(projection bson.M) *Repository[T] {
	r.optsFindOne.Projection = projection
	return r
}

func (r *Repository[T]) SetHint(hint bson.M) *Repository[T] {
	r.optsFind.Hint = hint
	return r
}

func (r *Repository[T]) SetHintOne(hint bson.M) *Repository[T] {
	r.optsFindOne.Hint = hint
	return r
}
