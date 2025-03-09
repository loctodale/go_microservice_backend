-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_product_base_9999
ADD CONSTRAINT pre_go_product_base_9999_category_id_foreign
FOREIGN KEY(category_id)
REFERENCES pre_go_product_category_9999(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pre_go_product_base_9999
DROP FOREIGN KEY pre_go_product_base_9999_category_id_foreign;
-- +goose StatementEnd
