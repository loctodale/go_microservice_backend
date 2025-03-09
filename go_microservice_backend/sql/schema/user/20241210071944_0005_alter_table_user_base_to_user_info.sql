-- +goose Up
-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_info_9999
ADD CONSTRAINT pre_go_acc_user_info_9999_verify_key_foreign
FOREIGN KEY (user_account) REFERENCES pre_go_acc_user_base_9999(user_account);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pre_go_acc_user_info_9999
DROP FOREIGN KEY pre_go_acc_user_info_9999_verify_key_foreign;
-- +goose StatementEnd
