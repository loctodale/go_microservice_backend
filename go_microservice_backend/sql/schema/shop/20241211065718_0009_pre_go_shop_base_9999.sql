-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_shop_base_9999 (
	shop_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	shop_account VARCHAR(255) NOT NULL UNIQUE,
	shop_password VARCHAR(255) NOT NULL,
	shop_salt VARCHAR(255) NOT NULL,
	shop_status INT default 0, -- 1 is 'active' , 0 is 'inactive'
	shop_created_at TIMESTAMP NOT NULL,
	shop_updated_at TIMESTAMP NOT NULL,
	shop_deleted_at TIMESTAMP NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_go_shop_base_9999;
-- +goose StatementEnd
