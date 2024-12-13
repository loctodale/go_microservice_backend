-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_key_token_9999
ADD CONSTRAINT pre_go_key_token_9999_shop_id_unique
UNIQUE(shop_id),
ADD CONSTRAINT pre_go_key_token_9999_shop_id_foreign
FOREIGN KEY(shop_id)
REFERENCES pre_go_shop_base_9999(shop_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pre_go_key_token_9999
DROP FOREIGN KEY pre_go_key_token_9999_shop_id_foreign;
-- +goose StatementEnd
