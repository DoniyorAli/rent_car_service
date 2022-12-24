package brand

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/brand"
	"MyProjects/RentCar_gRPC/rent_car_service/storage"
	"context"

	"github.com/google/uuid"
)

type BrandService struct {
	Stg storage.StorageInter
	brand.UnimplementedBrandServiceServer
}

// ?==============================================================================================================
func (b *BrandService) CreateBrand(ctx context.Context, req *brand.CreateBrandRequest) (*brand.Brand, error) {
	id := uuid.New()

	res, err := b.Stg.AddNewBrand(id.String(), req)

	if err != nil {
		return res, err
	}
	return res, nil
}

// ?==============================================================================================================
func (b *BrandService) GetBrandByID(ctx context.Context, req *brand.GetBrandByIDRequest) (*brand.Brand, error) {
	res, err := b.Stg.GetBrandById(req)
	if err != nil {
		return res, nil
	}
	return res, nil
}

// ?==============================================================================================================
func (b *BrandService) GetBrandList(ctx context.Context, req *brand.GetBrandListRequest) (*brand.GetBrandListResponse, error) {
	res, err := b.Stg.GetBrandList(req)
	if err != nil {
		return res, nil
	}
	return res, nil
}

// ?==============================================================================================================
func (b *BrandService) UpdateBrand(ctx context.Context, req *brand.UpdateBrandRequest) (*brand.Brand, error) {
	err := b.Stg.UpdateBrand(req.Id, req)
	if err != nil {
		return &brand.Brand{}, nil
	}

	box := &brand.GetBrandByIDRequest{BrandId: req.Id}

	res, err := b.Stg.GetBrandById(box)
	if err != nil {
		return &brand.Brand{}, nil
	}
	return res, nil
}

// ?==============================================================================================================
func (b *BrandService) DeleteBrand(ctx context.Context, req *brand.DeleteBrandRequest) (*brand.Brand, error) {
	res, err := b.Stg.GetBrandById((*brand.GetBrandByIDRequest)(req))
	if err != nil {
		return &brand.Brand{}, err
	}

	err = b.Stg.DeleteBrand(req)
	if err != nil {
		return &brand.Brand{}, err
	}
	return res, nil
}
