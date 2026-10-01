INSERT INTO clients (id, name, email, created_at) VALUES
  ('5e1b7a90-1111-4c2d-8e3f-000000000001', 'John Doe', 'john@j.com', NOW()),
  ('5e1b7a90-2222-4c2d-8e3f-000000000002', 'Jane Doe', 'jane@j.com', NOW());

INSERT INTO accounts (id, client_id, balance, created_at) VALUES
  ('8f2a3c4d-1111-4a5b-9c6d-000000000001', '5e1b7a90-1111-4c2d-8e3f-000000000001', 1000, NOW()),
  ('8f2a3c4d-2222-4a5b-9c6d-000000000002', '5e1b7a90-2222-4c2d-8e3f-000000000002', 1000, NOW());
