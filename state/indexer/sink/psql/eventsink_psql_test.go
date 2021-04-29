package psqlsink

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	schema "github.com/adlio/schema"
	_ "github.com/lib/pq"
	dockertest "github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
)

var db *sql.DB

var (
	user     = "postgres"
	password = "secret"
	port     = "5432"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
	dbName   = "postgres"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: DriverName,
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
			"listen_addresses = '*'",
		},
		ExposedPorts: []string{port},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// Set the container to expire in a minute to avoid orphaned containers
	// hanging around
	_ = resource.Expire(60)

	dsn = fmt.Sprintf(dsn, user, password, port, dbName)
	if err = pool.Retry(func() error {
		var err error

		es, err := NewPSQLEventSink(dsn)
		if err != nil {
			return err
		}

		db = es.store

		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resetDB()

	sm, err := readSchema()
	if err != nil {
		db.Close()
		log.Fatalf("Could not read schema: %s", err)
	}

	err = schema.NewMigrator().Apply(db, sm)
	if err != nil {
		db.Close()
		log.Fatalf("Could not apply schema to db: %s", err)
	}

	code := m.Run()

	// When you're done, kill and remove the container
	if err = pool.Purge(resource); err != nil {
		db.Close()
		log.Fatalf("Could not purge resource: %s", err)
	}

	db.Close()
	os.Exit(code)
}

func TestBlockFuncs(t *testing.T) {
	indexer := &PSQLEventSink{store: db}
	assert.NoError(t, indexer.IndexBlockEvents(getTestBlockHeader()))

	r, err := indexer.HasBlock(1)
	assert.Nil(t, err)
	assert.True(t, r)

	r, err = indexer.HasBlock(2)
	assert.Nil(t, err)
	assert.False(t, r)

	r2, err := indexer.SearchBlockEvents(context.TODO(), nil)
	assert.Nil(t, r2)
	assert.Equal(t, errors.New("block search is not supported via the postgres event sink"), err)
}

func TestTxFuncs(t *testing.T) {
	indexer := &PSQLEventSink{store: db}

	txResult := txResultWithEvents([]abci.Event{
		{Type: "account", Attributes: []abci.EventAttribute{{Key: []byte("number"), Value: []byte("1"), Index: true}}},
		{Type: "account", Attributes: []abci.EventAttribute{{Key: []byte("owner"), Value: []byte("Ivan"), Index: true}}},
		{Type: "", Attributes: []abci.EventAttribute{{Key: []byte("not_allowed"), Value: []byte("Vlad"), Index: true}}},
	})
	err := indexer.IndexTxEvents(txResult)
	assert.NoError(t, err)

	tx, err := indexer.GetTxByHash(types.Tx(txResult.Tx).Hash())
	assert.NoError(t, err)
	assert.Equal(t, txResult, tx)

	r2, err := indexer.SearchTxEvents(context.TODO(), nil)
	assert.Nil(t, r2)
	assert.Equal(t, errors.New("tx search is not supported via the postgres event sink"), err)
}

func getTestBlockHeader() types.EventDataNewBlockHeader {
	return types.EventDataNewBlockHeader{
		Header: types.Header{Height: 1},
		ResultBeginBlock: abci.ResponseBeginBlock{
			Events: []abci.Event{
				{
					Type: "begin_event",
					Attributes: []abci.EventAttribute{
						{
							Key:   []byte("proposer"),
							Value: []byte("FCAA001"),
							Index: true,
						},
					},
				},
			},
		},
		ResultEndBlock: abci.ResponseEndBlock{
			Events: []abci.Event{
				{
					Type: "end_event",
					Attributes: []abci.EventAttribute{
						{
							Key:   []byte("foo"),
							Value: []byte("100"),
							Index: true,
						},
					},
				},
			},
		},
	}
}

func readSchema() ([]*schema.Migration, error) {

	filename := "schema.sql"
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read sql file from '%s': %w", filename, err)
	}

	mg := &schema.Migration{}
	mg.ID = time.Now().Local().String() + " db schema"
	mg.Script = string(contents)
	return append([]*schema.Migration{}, mg), nil
}

func resetDB() {
	q := "DROP TABLE IF EXISTS block_events,tx_events,tx_results"
	_, err := db.Exec(q)
	if err != nil {
		db.Close()
		log.Fatalf("Could not reset TABLE: %s", err)
	}

	q = "DROP TYPE IF EXISTS block_event_type"
	_, err = db.Exec(q)
	if err != nil {
		db.Close()
		log.Fatalf("Could not reset TYPE: %s", err)
	}
}

func txResultWithEvents(events []abci.Event) *abci.TxResult {
	tx := types.Tx("HELLO WORLD")
	return &abci.TxResult{
		Height: 1,
		Index:  0,
		Tx:     tx,
		Result: abci.ResponseDeliverTx{
			Data:   []byte{0},
			Code:   abci.CodeTypeOK,
			Log:    "",
			Events: events,
		},
	}
}
