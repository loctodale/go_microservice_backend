-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_product_base_9999 (
	spu_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	category_id BIGINT UNSIGNED NOT NULL ,
	shop_id BIGINT NOT NULL,
	brand_id BIGINT NOT NULL,
	spu_name VARCHAR(255) NOT NULL UNIQUE,
	spu_description VARCHAR(255) NOT NULL,
	spu_img_url VARCHAR(255) NOT NULL,
	spu_video_url VARCHAR(255) NOT NULL,
	spu_sort BIGINT NOT NULL,
	spu_price DECIMAL(8, 2) NOT NULL,
	spu_status BIGINT NOT NULL,
	spu_created_at TIMESTAMP NOT NULL,
	spu_updated_at TIMESTAMP NOT NULL,
	spu_deleted_at TIMESTAMP
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_product_base_9999`;
-- +goose StatementEnd
