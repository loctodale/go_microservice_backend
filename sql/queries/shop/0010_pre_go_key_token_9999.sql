-- name: AddKeyToken :execresult
INSERT INTO pre_go_key_token_9999 (
    shop_id, refresh_token, key_created_at, key_updated_at, key_deleted_at
) VALUES (?, ?,NOW(), NOW(), NOW());