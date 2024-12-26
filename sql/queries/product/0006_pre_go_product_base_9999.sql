-- name: AddNewProduct :execresult
INSERT into pre_go_product_base_9999 (
	category_id, shop_id, brand_id, spu_name, spu_description, spu_img_url, spu_video_url, spu_sort, spu_price, spu_status, spu_created_at,spu_updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, 0, NOW(), NOW())