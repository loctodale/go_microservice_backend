-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_product_sku_9999
ADD CONSTRAINT pre_go_product_sku_9999_sku_id_foreign
FOREIGN KEY(spu_id)
REFERENCES `pre_go_product_base_9999`(spu_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pre_go_product_sku_9999
DROP FOREIGN KEY pre_go_product_sku_9999_sku_id_foreign;
-- +goose StatementEnd
