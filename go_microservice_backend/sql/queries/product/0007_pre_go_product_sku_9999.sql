-- name: AddNewProductSKU :execresult
INSERT into pre_go_product_sku_9999(
	spu_id, sku_price, sku_stock, sku_attribute_value,
    sku_created_at, sku_updated_at
) VALUES (?, ?, ?, ?, NOW(), NOW())