-- name: VerifyUserCredentials :one
SELECT count(1) as n_count FROM auth_users WHERE username = @username AND password = @password limit 1;

-- name: GetValidUserInfo :one
SELECT * FROM auth_users WHERE username = @username AND password = @password limit 1;

-- name: GetRoleNamesByUsername :many
SELECT auth_roles.name FROM auth_user_roles
LEFT JOIN auth_users ON auth_users.id = auth_user_roles.user_id
LEFT JOIN auth_roles ON auth_roles.id = auth_user_roles.role_id
WHERE auth_users.username = @username;
