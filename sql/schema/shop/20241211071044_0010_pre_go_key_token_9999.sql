-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_key_token_9999 (
	token_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	shop_id BIGINT UNSIGNED NOT NULL,
	public_key VARCHAR(255) NOT NULL,
	private_key VARCHAR(255) NOT NULL,
	refresh_token VARCHAR(255) NOT NULL,
	key_created_at TIMESTAMP NOT NULL,
	key_updated_at TIMESTAMP NOT NULL,
	key_deleted_at TIMESTAMP NOT NULL
	);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_go_key_token_9999;
-- +goose StatementEnd