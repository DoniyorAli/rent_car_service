package main

import (
	"MyProjects/RentCar_gRPC/rent_car_service/config"

	"MyProjects/RentCar_gRPC/rent_car_service/services/brand"
	"MyProjects/RentCar_gRPC/rent_car_service/services/car"

	"MyProjects/RentCar_gRPC/rent_car_service/storage"
	"MyProjects/RentCar_gRPC/rent_car_service/storage/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// * @license.name  Apache 2.0
// * @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

func main() {

	cfg := config.Load()
	psqlAUTH := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDatabase,
	)

	var err error
	var Stg storage.StorageInter
	Stg, err = postgres.InitDB(psqlAUTH)
	if err != nil {
		panic(err)
	}

	log.Printf("\ngRPC server running port%s with tcp protocol!\n", cfg.GRPCPort)

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()

	brandService := brand.NewBrandService(Stg)
	brand.RegisterAuthorServiceServer(srv, brandService)

	carService := car.NewBrandService(Stg)
	brand.RegisterBrandServiceServer(srv, carService)

	reflection.Register(srv)

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
