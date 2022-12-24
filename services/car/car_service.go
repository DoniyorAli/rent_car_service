package car

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/car"
	"MyProjects/RentCar_gRPC/rent_car_service/storage"
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CarService struct {
	Stg storage.StorageInter
	car.UnimplementedCarServiceServer
}

// ?===============================================================================================================
func (cs *CarService) CreateCar(ctx context.Context, req *car.CreateCarRequest) (*car.Car, error) {
	id := uuid.New()
	err := cs.Stg.AddCar(id.String(), req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.CreateCar: %s", err.Error())
	}

	res, err := cs.Stg.GetCarById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.GetBrandById: %s", err.Error())
	}

	return &car.Car{
		CarId: res.CarId,
		Model: res.Model,
		Color: res.Color,
		CarType: res.CarType,
		Mileage: res.Mileage,
		Year: res.Year,
		Price: res.Price,
		BrandId: res.BrandId, //! ??????????????????????? res.Brand.BramdId
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

// ?===============================================================================================================
func (cs *CarService) GetCarByID(ctx context.Context, req *car.GetCarByIDRequest) (*car.GetCarByIDResponse, error) {
	res, err := cs.Stg.GetCarById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.GetCarById: %s", err.Error())
	}
	return res, nil
}

// ?===============================================================================================================
func (cs *CarService) GetCarList(ctx context.Context, req *car.GetCarListRequest) (*car.GetCarListResponse, error) {
	res, err := cs.Stg.GetCarList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.GetCarList: %s", err.Error())
	}
	return res, nil
}

// ?===============================================================================================================
func (cs *CarService) UpdateCar(ctx context.Context, req *car.UpdateCarRequest) (*car.Car, error) {
	err := cs.Stg.UpdateCar(req.Id, req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.UpdateCar: %s", err.Error())
	}

	res, err := cs.Stg.GetCarById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "cs.Stg.UpdateCar: %s", err.Error())
	}

	return &car.Car{
		CarId: res.CarId,
		Model: res.Model,
		Color: res.Color,
		CarType: res.CarType,
		Mileage: res.Mileage,
		Year: res.Year,
		Price: res.Price,
		BrandId: res.BrandId,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}

// ?===============================================================================================================
func (cs *CarService) DeleteCar(ctx context.Context, req *car.DeleteCarRequest) (*car.Car, error) {
	res, err := cs.Stg.GetCarById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.GetCarById: %s", err.Error())
	}

	err = cs.Stg.DeleteCar(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cs.Stg.DeleteCar: %s", err.Error())
	}

	return &car.Car{
		CarId: res.CarId,
		Model: res.Model,
		Color: res.Color,
		CarType: res.CarType,
		Mileage: res.Mileage,
		Year: res.Year,
		Price: res.Price,
		BrandId: res.BrandId,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
	}, nil
}
