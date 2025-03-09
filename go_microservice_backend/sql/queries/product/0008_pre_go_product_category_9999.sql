-- name: AddNewCategory :execresult
INSERT into pre_go_product_category_9999 (
	parent_id, category_name, has_active_children, category_spu_count, category_status, category_description, category_icon, category_sort,
	category_created_at, category_updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW());