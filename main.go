package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mongodb/grip"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	StatsDB *mongo.Client
	TestDB  *mongo.Client
	ctx     context.Context
)

func main() {
	ctx = context.Background()
	opts := options.Client().ApplyURI(
		"mongodb+srv://admin:<password>@skunkworks20210412.zgc06.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
	).SetAuth(options.Credential{Username: "admin", Password: "ddfda"}).SetServerAPIOptions(&options.ServerAPIOptions{ServerAPIVersion: options.ServerAPIVersion1})
	client, err := mongo.Connect(ctx, opts)
	grip.Error(err)
	StatsDB = client

	opts = options.Client().ApplyURI(
		"mongodb+srv://admin:<password>@testdb1.zgc06.mongodb.net/myFirstDatabase?retryWrites=true&w=majority",
	).SetAuth(options.Credential{Username: "admin", Password: "ddfda"}).SetServerAPIOptions(&options.ServerAPIOptions{ServerAPIVersion: options.ServerAPIVersion1})
	client, err = mongo.Connect(ctx, opts)
	grip.Error(err)
	TestDB = client

	// StoreSnapshot()
	// grip.Info(GetSnapshots("foo", "collection1"))

	serve()
}

func serve() {
	http.HandleFunc("/snapshots", GetSnapshotsHandler)
	http.ListenAndServe(":3001", nil)
}

// https://docs.mongodb.com/manual/reference/command/collStats/#mongodb-dbcommand-dbcmd.collStats
type CollStats struct {
	Collection  string `bson:"ns"`
	SizeData    int    `bson:"size"` // uncompressed
	SizeStorage int    `bson:"storageSize"`
	SizeIndices int    `bson:"totalIndexSize"`
	SizeTotal   int    `bson:"totalSize"` // SizeStorage + SizeIndices
	NumDocs     int    `bson:"count"`
	NumIndices  int    `bson:"nindexes"`
	// indexSizes
}

type Snapshot struct {
	DB           string    `bson:"db"`
	Collection   string    `bson:"collection"`
	GatheredTime time.Time `bson:"gathered_time"`
	Stats        CollStats `bson:"stats"`
}

func GetStats() CollStats {
	res := TestDB.Database("foo").RunCommand(ctx, bson.M{
		"collStats": "collection1",
	})
	grip.Error(res.Err())
	stats := CollStats{}
	grip.Error(res.Decode(&stats))
	return stats
}

func GetSnapshot() Snapshot {
	stats := GetStats()
	pieces := strings.Split(stats.Collection, ".")
	snapshot := Snapshot{
		DB:           pieces[0],
		Collection:   pieces[1],
		GatheredTime: time.Now(),
		Stats:        stats,
	}

	return snapshot
}

func StoreSnapshot() {
	snapshot := GetSnapshot()
	_, err := StatsDB.Database("snapshots").Collection("snapshots").InsertOne(ctx, snapshot)
	grip.Error(err)
}

func GetSnapshots(db, collection string) []Snapshot {
	cursor, err := StatsDB.Database("snapshots").Collection("snapshots").Find(ctx, bson.M{
		"db":         db,
		"collection": collection,
	}, options.Find().SetSort(bson.M{
		"gathered_time": -1,
	}))
	grip.Error(err)
	snapshots := []Snapshot{}
	grip.Error(cursor.All(ctx, &snapshots))
	return snapshots
}

func GetSnapshotsHandler(rw http.ResponseWriter, r *http.Request) {
	data := GetSnapshots("foo", "collection1")
	b, err := json.Marshal(data)
	grip.Error(err)
	rw.Write(b)
}
