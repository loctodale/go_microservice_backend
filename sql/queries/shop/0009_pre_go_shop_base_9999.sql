-- name: CheckShopBaseIsExists :one
SELECT COUNT(shop_account)
FROM pre_go_shop_base_9999
WHERE shop_account = ?;

