// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package base

import (
	"github.com/cockroachdb/cdc-sink/internal/util/diag"
)

// Injectors from injector.go:

// NewFixture constructs a self-contained test fixture.
func NewFixture() (*Fixture, func(), error) {
	context, cleanup, err := ProvideContext()
	if err != nil {
		return nil, nil, err
	}
	diagnostics, cleanup2 := diag.New(context)
	sourcePool, cleanup3, err := ProvideSourcePool(context, diagnostics)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	sourceSchema, cleanup4, err := ProvideSourceSchema(context, sourcePool)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	stagingPool, cleanup5, err := ProvideStagingPool(context)
	if err != nil {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	stagingSchema, cleanup6, err := ProvideStagingSchema(context, stagingPool)
	if err != nil {
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	targetPool, cleanup7, err := ProvideTargetPool(context, sourcePool, diagnostics)
	if err != nil {
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	targetStatements, cleanup8 := ProvideTargetStatements(targetPool)
	targetSchema, cleanup9, err := ProvideTargetSchema(context, diagnostics, targetPool, targetStatements)
	if err != nil {
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	fixture := &Fixture{
		Context:      context,
		SourcePool:   sourcePool,
		SourceSchema: sourceSchema,
		StagingPool:  stagingPool,
		StagingDB:    stagingSchema,
		TargetCache:  targetStatements,
		TargetPool:   targetPool,
		TargetSchema: targetSchema,
	}
	return fixture, func() {
		cleanup9()
		cleanup8()
		cleanup7()
		cleanup6()
		cleanup5()
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
