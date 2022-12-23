package car

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/blogpost"
	"MyProjects/RentCar_gRPC/rent_car_service/storage"
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authorService struct {
	stg storage.StorageInter
	blogpost.UnimplementedAuthorServiceServer
}

func NewAuthorService(stg storage.StorageInter) *authorService {
	return &authorService{
		stg: stg,
	}
}

func (s *authorService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")

	return &blogpost.Pong{
		Message: "OK",
	}, nil
}

//?==============================================================================================================

func (s *authorService) CreateAuthor(ctx context.Context, req *blogpost.CreateAuthorRequest) (*blogpost.CreateAuthorResponse, error) {
	id := uuid.New()

	err := s.stg.AddAuthor(id.String(), req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.CreateAuthor: %s", err.Error())
	}

	res, err := s.stg.GetAuthorById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	return &blogpost.CreateAuthorResponse{
		Id: res.Id,
	}, nil
}

func (s *authorService) GetAuthorByID(ctx context.Context, req *blogpost.GetAuthorByIDRequest) (*blogpost.Author, error) {
	author, err := s.stg.GetAuthorById(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorById: %s", err.Error())
	}
	return author, nil
}

func (s *authorService) GetAuthorList(ctx context.Context, req *blogpost.GetAuthorListRequest) (*blogpost.GetAuthorListResponse, error) {
	res, err := s.stg.GetAuthorList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetAuthorList: %s", err.Error())
	}

	return res, nil
}

func (s *authorService) UpdateAuthor(ctx context.Context, req *blogpost.UpdateAuthorRequest) (*blogpost.UpdateAuthorResponse, error) {

	err := s.stg.UpdateAuthor(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateAuthor: %s", err.Error())
	}

	_, err = s.stg.GetAuthorById(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateAuthor: %s", err.Error())
	}

	return &blogpost.UpdateAuthorResponse{
		Status: "Updated",
	}, nil
}

func (s *authorService) DeleteAuthor(ctx context.Context, req *blogpost.DeleteAuthorRequest) (*blogpost.DeleteAuthorResponse, error) {

	err := s.stg.DeleteAuthor(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteAuthor: %s", err.Error())
	}

	return &blogpost.DeleteAuthorResponse{
		Status: "Deleted",
	}, nil
}
