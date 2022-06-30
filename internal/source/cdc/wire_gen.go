// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cdc

import (
	"github.com/cockroachdb/cdc-sink/internal/target/apply"
	"github.com/cockroachdb/cdc-sink/internal/target/apply/fan"
	"github.com/cockroachdb/cdc-sink/internal/target/auth/trust"
	"github.com/cockroachdb/cdc-sink/internal/target/resolve"
	"github.com/cockroachdb/cdc-sink/internal/target/schemawatch"
	"github.com/cockroachdb/cdc-sink/internal/target/sinktest"
	"github.com/cockroachdb/cdc-sink/internal/target/stage"
	"github.com/cockroachdb/cdc-sink/internal/target/timekeeper"
)

// Injectors from test_fixture.go:

func newTestFixture() (*testFixture, func(), error) {
	context, cleanup, err := sinktest.ProvideContext()
	if err != nil {
		return nil, nil, err
	}
	dbInfo, err := sinktest.ProvideDBInfo(context)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	pool := sinktest.ProvidePool(dbInfo)
	stagingDB, cleanup2, err := sinktest.ProvideStagingDB(context, pool)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	testDB, cleanup3, err := sinktest.ProvideTestDB(context, pool)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	baseFixture := sinktest.BaseFixture{
		Context:   context,
		DBInfo:    dbInfo,
		Pool:      pool,
		StagingDB: stagingDB,
		TestDB:    testDB,
	}
	watchers, cleanup4 := schemawatch.ProvideFactory(pool)
	appliers, cleanup5 := apply.ProvideFactory(watchers)
	fans := &fan.Fans{
		Appliers: appliers,
		Pool:     pool,
	}
	metaTable := sinktest.ProvideMetaTable(stagingDB, testDB)
	stagers := stage.ProvideFactory(pool, stagingDB)
	targetTable := sinktest.ProvideTimestampTable(stagingDB, testDB)
	timeKeeper, cleanup6, err := timekeeper.ProvideTimeKeeper(context, pool, targetTable)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	resolvers, cleanup7, err := resolve.ProvideFactory(context, appliers, metaTable, pool, stagers, timeKeeper, watchers)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	watcher, err := sinktest.ProvideWatcher(context, testDB, watchers)
	if err != nil {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	fixture := &sinktest.Fixture{
		BaseFixture: baseFixture,
		Appliers:    appliers,
		Fans:        fans,
		Resolvers:   resolvers,
		Stagers:     stagers,
		TimeKeeper:  timeKeeper,
		Watchers:    watchers,
		MetaTable:   metaTable,
		Watcher:     watcher,
	}
	authenticator := trust.New()
	handler := &Handler{
		Appliers:      appliers,
		Authenticator: authenticator,
		Pool:          pool,
		Resolvers:     resolvers,
		Stores:        stagers,
	}
	cdcTestFixture := &testFixture{
		Fixture: fixture,
		Handler: handler,
	}
	return cdcTestFixture, func() {
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// test_fixture.go:

type testFixture struct {
	*sinktest.Fixture
	Handler *Handler
}
