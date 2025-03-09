package model

type CreateCategoryInput struct {
	ParentId            int    `json:"parent_id"`
	CategoryName        string `json:"category_name"`
	HasActiveChildren   bool   `json:"has_active_children"`
	CategorySPUCount    int    `json:"category_spu_count"`
	CategoryStatus      int    `json:"category_status"`
	CategoryDescription string `json:"category_description"`
	CategoryIcon        string `json:"category_icon"`
	CategorySort        int    `json:"category_sort"`
}

type CreateProductInput struct {
	CategoryID int    `json:"category_id"`
	BrandID    int    `json:"brand_id"`
	SPUName    string `json:"spu_name"`
	SPUDesc    string `json:"spu_desc"`
	SPUImg     string `json:"spu_img"`
	SPUVideo   string `json:"spu_video"`
	SPUSort    int    `json:"spu_sort"`
	SPUPrice   int    `json:"spu_price"`
}

type CreateProductSKUInput struct {
	SPUID             int    `json:"spu_id"`
	SKUPrice          int    `json:"sku_price"`
	SKUStock          int    `json:"sku_stock"`
	SKUAttributeValue string `json:"sku_attribute_value"`
}
