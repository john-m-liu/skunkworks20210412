package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/mongodb/amboy"
	"github.com/mongodb/grip"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	StatsDB *mongo.Client
	TestDB  *mongo.Client
	ctx     context.Context
	q       amboy.Queue
)

const (
	SnapshotDB         = "snapshots"
	SnapshotCollection = "snapshots"
	TestDBName         = "foo"
	TestCollName       = "collection1"
)

func main() {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()
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

	go func() {
		t := time.NewTimer(0)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				StoreSnapshot()
				t.Reset(10 * time.Second)
			}
		}
	}()

	go DeleteLoad(ctx)
	go WriteLoad(ctx)

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
	res := TestDB.Database(TestDBName).RunCommand(ctx, bson.M{
		"collStats": TestCollName,
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
	_, err := StatsDB.Database(SnapshotDB).Collection(SnapshotCollection).InsertOne(ctx, snapshot)
	grip.Error(err)
}

func GetSnapshots(db, collection string) []Snapshot {
	cursor, err := StatsDB.Database(SnapshotDB).Collection(SnapshotCollection).Find(ctx, bson.M{
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

func DeleteLoad(ctx context.Context) {
	t := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			_, err := TestDB.Database(TestDBName).Collection(TestCollName).DeleteOne(ctx, bson.M{})
			grip.Error(err)
			n := rand.Intn(1000)
			t.Reset(time.Duration(n) * time.Millisecond)
		}
	}
}

func WriteLoad(ctx context.Context) {
	t := time.NewTimer(0)
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
			size := rand.Intn(5000)
			s := make([]rune, size)
			for i := range s {
				s[i] = letters[rand.Intn(len(letters))]
			}
			doc := bson.M{
				"a": string(s),
			}
			_, err := TestDB.Database(TestDBName).Collection(TestCollName).InsertOne(ctx, doc)
			grip.Error(err)
			n := rand.Intn(900)
			t.Reset(time.Duration(n) * time.Millisecond)
		}
	}
}
