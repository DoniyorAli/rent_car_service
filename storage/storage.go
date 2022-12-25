package storage

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/brand"
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/car"
)

type StorageInter interface {
	//* Brand
	AddNewBrand(id string, req *brand.CreateBrandRequest) (res *brand.Brand, err error)
	GetBrandById(req *brand.GetBrandByIDRequest) (*brand.Brand, error)
	GetBrandList(req *brand.GetBrandListRequest) (*brand.GetBrandListResponse, error)
	UpdateBrand(id string, req *brand.UpdateBrandRequest) error
	DeleteBrand(req *brand.DeleteBrandRequest) error

	//* Car
	AddCar(id string, req *car.CreateCarRequest) error
	GetCarById(id string) (*car.GetCarByIDResponse, error)
	GetCarList(limit, offset int, search string) (*car.GetCarListResponse, error)
	UpdateCar(id string, box *car.UpdateCarRequest) error
	DeleteCar(id string) error
}
