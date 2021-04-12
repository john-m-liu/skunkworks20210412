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
	Collection  string `bson:"ns" json:"collection"`
	SizeData    int    `bson:"size" json:"size_data"` // uncompressed
	SizeStorage int    `bson:"storageSize" json:"size_storage"`
	SizeIndices int    `bson:"totalIndexSize" json:"size_indices"`
	SizeTotal   int    `bson:"totalSize" json:"size_total"` // SizeStorage + SizeIndices
	NumDocs     int    `bson:"count" json:"num_docs"`
	NumIndices  int    `bson:"nindexes" json:"num_indices"`
	// indexSizes
}

type Snapshot struct {
	DB               string    `bson:"db" json:"db"`
	Collection       string    `bson:"collection" json:"collection"`
	GatheredTime     time.Time `bson:"gathered_time" json:"gathered_time"`
	GatheredTimeUnix int64     `json:"gathered_time_unix"`
	Stats            CollStats `bson:"stats" json:"stats"`
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
	now := time.Now()
	stats := GetStats()
	pieces := strings.Split(stats.Collection, ".")
	snapshot := Snapshot{
		DB:           pieces[0],
		Collection:   pieces[1],
		GatheredTime: now,
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
	for i, s := range snapshots {
		snapshots[i].GatheredTimeUnix = s.GatheredTime.Unix()
	}
	return snapshots
}

func GetSnapshotsHandler(rw http.ResponseWriter, r *http.Request) {
	data := GetSnapshots("foo", "collection1")
	b, err := json.Marshal(data)
	grip.Error(err)
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Write(b)
}
