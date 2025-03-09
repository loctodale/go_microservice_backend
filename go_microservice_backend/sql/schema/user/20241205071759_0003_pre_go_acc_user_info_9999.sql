-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_info_9999 (
	user_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY COMMENT 'User ID',
	user_account VARCHAR(255) NOT NULL COMMENT 'User account', -- Account username
	user_nickname VARCHAR(255) COMMENT 'User nickname', -- Nickname of the user
	user_avatar VARCHAR(255) COMMENT 'User avatar', -- Avatar image URL for the user
	user_state TINYINT UNSIGNED NOT NULL COMMENT 'User state: 0=Locked, 1=Active',
	user_mobile VARCHAR(20) COMMENT 'Mobile phone number', -- User's mobile phone number
	user_gender TINYINT UNSIGNED COMMENT 'User gender: 0=Secret, 1=Male, 2=Female',
	user_birthday DATE COMMENT 'User birthday', -- Date of birth
	user_email VARCHAR(255) COMMENT 'User email address', -- Email address
	user_is_authentication TINYINT UNSIGNED NOT NULL COMMENT 'Authentication status',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp',
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp'
	) COMMENT='Table for user account information';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_info_9999`
-- +goose StatementEnd
