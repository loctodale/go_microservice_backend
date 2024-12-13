-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_token_used_9999
ADD CONSTRAINT pre_go_key_token_used_9999_key_token_foreign
FOREIGN KEY(token_id)
REFERENCES pre_go_key_token_9999(token_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pre_go_token_used_9999
DROP FOREIGN KEY pre_go_key_token_used_9999_key_token_foreign;
-- +goose StatementEnd