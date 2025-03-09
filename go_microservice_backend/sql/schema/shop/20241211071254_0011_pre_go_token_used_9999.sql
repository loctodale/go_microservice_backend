-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_token_used_9999 (
	token_used_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	token_id BIGINT UNSIGNED NOT NULL,
	refresh_token_used VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pre_go_token_used_9999;
-- +goose StatementEnd
