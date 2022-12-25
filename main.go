package main

import (
	"MyProjects/RentCar_gRPC/rent_car_service/config"
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/brand"
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/car"

	brandService "MyProjects/RentCar_gRPC/rent_car_service/services/brand"
	carService "MyProjects/RentCar_gRPC/rent_car_service/services/car"

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
	var interStg storage.StorageInter
	interStg, err = postgres.InitDB(psqlAUTH)
	if err != nil {
		panic(err)
	}

	log.Printf("\ngRPC server running port%s with tcp protocol!\n", cfg.GRPCPort)

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		panic(err)
	}

	newS := grpc.NewServer()

	b := &brandService.BrandService{
		Stg: interStg,
	}
	brand.RegisterBrandServiceServer(newS, b)

	c := &carService.CarService{
		Stg: interStg,
	}
	car.RegisterCarServiceServer(newS, c)

	// srv := grpc.NewServer()

	// brandService := brand.NewBrandService(Stg)
	// brand.RegisterBrandServiceServer(srv, brandService)

	// carService := car.NewCarService(Stg)
	// brand.RegisterCarServiceServer(srv, carService)

	reflection.Register(newS)

	if err := newS.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
