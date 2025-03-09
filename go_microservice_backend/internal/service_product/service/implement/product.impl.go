package implement

import (
	"context"
	"fmt"
	"go_microservice_backend_api/internal/model"
	"go_microservice_backend_api/internal/service_product/database"
	"go_microservice_backend_api/pkg/response"
	"strconv"
)

type sProductService struct {
	r *database.Queries
}

func NewProductService(r *database.Queries) *sProductService {
	return &sProductService{r: r}
}

func (s *sProductService) AddNewProduct(ctx context.Context, in model.CreateProductInput, shopId string) (int, error) {
	newShopID, err := strconv.Atoi(shopId)
	if err != nil {
		return response.ErrCodeParamInvalid, err
	}
	result, err := s.r.AddNewProduct(ctx, database.AddNewProductParams{
		BrandID:        int64(in.BrandID),
		CategoryID:     uint64(in.CategoryID),
		ShopID:         int64(newShopID),
		SpuDescription: in.SPUDesc,
		SpuImgUrl:      in.SPUImg,
		SpuName:        in.SPUName,
		SpuPrice:       strconv.Itoa(in.SPUPrice),
		SpuSort:        int64(in.SPUSort),
		SpuVideoUrl:    in.SPUVideo,
	})
	if err != nil {
		return response.ErrCreateFailed, err
	}
	_, err = result.LastInsertId()
	return 1, nil
}

func (s *sProductService) AddNewSKUProduct(ctx context.Context, in model.CreateProductSKUInput) (statusCode int, err error) {
	result, err := s.r.AddNewProductSKU(ctx, database.AddNewProductSKUParams{
		SkuAttributeValue: in.SKUAttributeValue,
		SkuPrice:          strconv.Itoa(in.SKUPrice),
		SkuStock:          int64(in.SKUStock),
		SpuID:             uint64(in.SPUID),
	})
	if err != nil {
		return response.ErrCreateFailed, fmt.Errorf(err.Error())
	}
	_, err = result.LastInsertId()
	return 1, nil
}
