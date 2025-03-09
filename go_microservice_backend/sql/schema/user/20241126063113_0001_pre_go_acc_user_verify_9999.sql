-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pre_go_acc_user_verify_9999 (
	verify_id INT AUTO_INCREMENT PRIMARY KEY, -- ID of the OTP record
	verify_otp VARCHAR(6) NOT NULL,           -- OTP code (verification code)
	verify_key VARCHAR(255) NOT NULL,        -- User's email or phone to identify the OTP recipient
	verify_key_hash VARCHAR(255) NOT NULL,   -- Hash of verify_key
	verify_type INT DEFAULT 1,               -- 1: Email, 2: Phone, etc. (Type of verification)
	is_verified INT DEFAULT 0,               -- 0: No, 1: Yes (OTP verification status)
	is_deleted INT DEFAULT 0,                -- 0: No, 1: Yes (Deletion status)
	verify_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Record creation time
	verify_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Record update time
	INDEX idx_verify_otp (verify_otp),       -- Index on the verify_otp field
	UNIQUE KEY unique_verify_key (verify_key) -- Ensure verify_key is unique
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='account_user_verify';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS `pre_go_acc_user_verify_9999`;
-- +goose StatementEnd
