-- name: ListProducts :many
SELECT id, name, price, category, img_thumb, img_mobile, img_tablet, img_desktop
FROM products
ORDER BY id;

-- name: GetProduct :one
SELECT id, name, price, category, img_thumb, img_mobile, img_tablet, img_desktop
FROM products
WHERE id = $1;

-- name: GetProductsByIDs :many
SELECT id, name, price, category, img_thumb, img_mobile, img_tablet, img_desktop
FROM products
WHERE id = ANY($1::text[]);
