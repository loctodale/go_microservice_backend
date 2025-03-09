-- name: AddKeyToken :execresult
INSERT INTO pre_go_key_token_9999 (
    shop_id, refresh_token, shop_credential_id, key_created_at, key_updated_at, key_deleted_at
) VALUES (?, ?,?,NOW(), NOW(), NOW());

-- name: GetKeyTokenByShopId :one
SELECT shop_credential_id,refresh_token, token_id
FROM pre_go_key_token_9999
WHERE shop_id = ?;

-- name: UpdateKeyToken :execresult
UPDATE pre_go_key_token_9999
SET refresh_token = ?, shop_credential_id = ?, key_updated_at = NOW()
WHERE token_id = ?;