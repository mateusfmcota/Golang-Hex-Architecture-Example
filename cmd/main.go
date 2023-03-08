package main

import (
	"log"
	// "os"
	"hex/internal/adapters/app/api"
	"hex/internal/adapters/core/arithmetic"
	"hex/internal/adapters/framework/right/db"
	"hex/internal/ports"

	gRPC "hex/internal/adapters/framework/left/grpc"
)

func main() {

	var err error

	//ports
	var dbaseAdapter ports.DbPort
	var core ports.ArithmeticPort
	var appAdapter ports.APIPort
	var gRPCAdapter ports.GRPCPort

	// dbaseDriver := os.Getenv("DB_DRIVER")
	// dbaseSourceNAme := os.Getenv("DB_NAME")

	dbaseAdapter, err = db.NewAdapter("sqlite3", "./db.db")

	if err != nil {
		log.Fatal(err)
	}

	defer dbaseAdapter.CloseDbConnection()

	core = arithmetic.NewAdapter()

	appAdapter = api.NewAdapter(dbaseAdapter, core)

	gRPCAdapter = gRPC.NewAdapter(appAdapter)

	gRPCAdapter.Run()

}
