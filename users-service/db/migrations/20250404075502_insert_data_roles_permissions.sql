-- +goose Up
-- +goose StatementBegin

INSERT INTO roles (id, name) VALUES
    (1, 'user'),
    (2, 'manager'),
    (3, 'admin')
ON CONFLICT (id) DO NOTHING;

INSERT INTO permissions (id, code, description) VALUES
    (1, 'category.create', 'Create a category'),
    (2, 'category.update', 'Update a category'),
    (3, 'category.delete', 'Delete a category'),
    (4, 'category.getbyid', 'Get category by ID'),
    (5, 'category.getall', 'Get all categories'),

    (6, 'product.create', 'Create a product'),
    (7, 'product.update', 'Update a product'),
    (8, 'product.delete', 'Delete a product'),
    (9, 'product.getbyid', 'Get product by ID'),
    (10, 'product.getall', 'Get all products'),

    (11, 'auth.refresh', 'Refresh authentication token')
ON CONFLICT (id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id) VALUES
    (1, 4),
    (1, 5),
    (1, 9),
    (1, 10);

INSERT INTO role_permissions (role_id, permission_id) VALUES
    (2, 1),
    (2, 2),
    (2, 3),
    (2, 4),
    (2, 5),

    (2, 6),
    (2, 7),
    (2, 8),
    (2, 9),
    (2, 10);

INSERT INTO role_permissions (role_id, permission_id) VALUES
    (3, 1),
    (3, 2),
    (3, 3),
    (3, 4),
    (3, 5),

    (3, 6),
    (3, 7),
    (3, 8),
    (3, 9),
    (3, 10),

    (3, 11);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM roles;
-- +goose StatementEnd
