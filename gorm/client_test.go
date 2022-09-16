package gorm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/nubunto/maeve"
	"github.com/nubunto/maeve/gorm"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
)

type ClientSuite struct {
	suite.Suite

	dockerResource *dockertest.Resource
	dockerPool *dockertest.Pool

	client *gorm.Client
}

func (cs *ClientSuite) SetupSuite() {
	containerPool, err := dockertest.NewPool("")
	if err != nil {
		panic(err)
	}

	pg, err := containerPool.Run("postgres", "13", []string{"POSTGRES_PASSWORD=mytestpw", "POSTGRES_DB=testdb"})
	if err != nil {
		panic(err)
	}

	dsn := fmt.Sprintf("user=postgres password=mytestpw dbname=testdb host=localhost port=%s sslmode=disable", pg.GetPort("5432/tcp"))
	err = containerPool.Retry(func() error {
		cs.client, err = gorm.New(dsn)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	cs.dockerResource = pg
	cs.dockerPool = containerPool
}

func (cs *ClientSuite) TestGormClientFetch() {
	cs.NoError(cs.client.Put(context.Background(), maeve.KV("users/1", "my user here")))
	keys, err := cs.client.Fetch(context.Background(), maeve.Path("users/*"))
	cs.NoError(err)
	cs.NotEmpty(keys)
}

func (cs *ClientSuite) TestGormClientDelete() {
	cs.NoError(cs.client.Put(context.Background(), maeve.KV("users/1", "my user here")))
	cs.NoError(cs.client.Put(context.Background(), maeve.KV("users/2", "my 2nd user here")))
	cs.NoError(cs.client.Put(context.Background(), maeve.KV("users/3", "my 3rd user here")))

	keys, err := cs.client.Fetch(context.Background(), maeve.Path("users/*"))
	cs.NotEmpty(keys)
	cs.NoError(err)

	err = cs.client.Delete(context.Background(), maeve.Path("users/*"))
	cs.NoError(err)

	keys, err = cs.client.Fetch(context.Background(), maeve.Path("users/*"))
	cs.NoError(err)

	keys, err = cs.client.Fetch(context.Background(), maeve.Path("users/*"))
	cs.NoError(err)
	cs.Empty(keys)
}

func (cs *ClientSuite) TearDownSuite() {
	cs.NoError(cs.dockerPool.Purge(cs.dockerResource))
}

func TestClient(t *testing.T) {
	t.Parallel()

	suite.Run(t, new(ClientSuite))
}
