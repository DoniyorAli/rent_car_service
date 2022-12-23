package brand

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/blogpost"
	"MyProjects/RentCar_gRPC/rent_car_service/storage"
	"context"
	"log"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BrandService struct {
	Stg storage.StorageInter
	blogpost.UnimplementedArticleServiceServer
}

func NewArticleService(stg storage.StorageInter) *BrandService {
	return &BrandService{
		Stg: stg,
	}
}

func (s *BrandService) Ping(ctx context.Context, req *blogpost.Empty) (*blogpost.Pong, error) {
	log.Println("Ping")

	return &blogpost.Pong{
		Message: "OK",
	}, nil
}

//?==============================================================================================================

func (s *BrandService) CreateArticle(ctx context.Context, req *blogpost.CreateArticleRequest) (*blogpost.Article, error) {
	id := uuid.New()

	err := s.stg.AddNewArticle(id.String(), req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.AddNewArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleById(id.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

//?==============================================================================================================

func (s *BrandService) UpdateArticle(ctx context.Context, req *blogpost.UpdateArticleRequest) (*blogpost.Article, error) {
	err := s.stg.UpdateArticle(req)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.UpdateArticle: %s", err.Error())
	}

	article, err := s.stg.GetArticleById(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

//?==============================================================================================================

func (s *BrandService) DeleteArticle(ctx context.Context, req *blogpost.DeleteArticleRequest) (*blogpost.Article, error) {

	article, err := s.stg.GetArticleById(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	err = s.stg.DeleteArticle(article.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.DeleteArticle: %s", err.Error())
	}

	return &blogpost.Article{
		Id:        article.Id,
		Content:   article.Content,
		AuthorId:  article.Author.Id,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}, nil
}

//?==============================================================================================================

func (s *BrandService) GetArticleList(ctx context.Context, req *blogpost.GetArticleListRequest) (*blogpost.GetArticleListResponse, error) {
	res, err := s.stg.GetArticleList(int(req.Offset), int(req.Limit), req.Search)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleList: %s", err.Error())
	}

	return res, nil
}

//?==============================================================================================================

func (s *BrandService) GetArticleByID(ctx context.Context, req *blogpost.GetArticleByIDRequest) (*blogpost.GetArticleByIDResponse, error) {
	article, err := s.stg.GetArticleById(req.Id)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "s.stg.GetArticleById: %s", err.Error())
	}

	return article, nil
}
