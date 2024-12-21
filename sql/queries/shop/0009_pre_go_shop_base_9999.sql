-- name: CheckShopBaseIsExists :one
SELECT COUNT(shop_account)
FROM pre_go_shop_base_9999
WHERE shop_account = ?;

-- name: AddIntoShopBase :execresult
INSERT into pre_go_shop_base_9999 (
    shop_account, shop_password, shop_salt , shop_status, shop_created_at, shop_updated_at, shop_deleted_at
) values (?,?,?,0, NOW(), NOW(), NOW());

-- name: ChangePassword :execresult
UPDATE pre_go_shop_base_9999
SET shop_password = ?
WHERE shop_account = ?;

-- name: GetShopByAccount :one
SELECT shop_salt, shop_id, shop_password
FROM pre_go_shop_base_9999
WHERE shop_account = ?;


