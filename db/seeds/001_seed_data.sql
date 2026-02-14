-- Seed data for local development
-- Run after migrations have been applied

-- Companies
INSERT INTO companies (id, name, address, phone, email) VALUES
  ('a1b2c3d4-0001-4000-8000-000000000001', 'Sentinel Security Services', '45 Victoria Street, London SW1H 0EU', '+44 20 7946 0123', 'info@sentinel-security.co.uk'),
  ('a1b2c3d4-0001-4000-8000-000000000002', 'Guardian Protection Group', '12 Queen Street, Manchester M2 5HT', '+44 161 496 0456', 'contact@guardian-protection.co.uk'),
  ('a1b2c3d4-0001-4000-8000-000000000003', 'Apex Security Solutions', '8 George Square, Glasgow G2 1DY', '+44 141 352 0789', 'hello@apex-security.co.uk');

-- Worksites
INSERT INTO worksites (id, company_id, name, address, latitude, longitude) VALUES
  ('b2c3d4e5-0001-4000-8000-000000000001', 'a1b2c3d4-0001-4000-8000-000000000001', 'Canary Wharf Tower', '1 Canada Square, London E14 5AB', 51.5054, -0.0197),
  ('b2c3d4e5-0001-4000-8000-000000000002', 'a1b2c3d4-0001-4000-8000-000000000001', 'Westminster Shopping Centre', '15 Victoria Street, London SW1H 0HW', 51.4975, -0.1357),
  ('b2c3d4e5-0001-4000-8000-000000000003', 'a1b2c3d4-0001-4000-8000-000000000001', 'Kings Cross Development', '1 Pancras Square, London N1C 4AG', 51.5320, -0.1240),
  ('b2c3d4e5-0001-4000-8000-000000000004', 'a1b2c3d4-0001-4000-8000-000000000002', 'Arndale Centre', 'Market Street, Manchester M4 3AQ', 53.4831, -2.2380),
  ('b2c3d4e5-0001-4000-8000-000000000005', 'a1b2c3d4-0001-4000-8000-000000000002', 'MediaCityUK', 'MediaCityUK, Salford M50 2HF', 53.4727, -2.2984),
  ('b2c3d4e5-0001-4000-8000-000000000006', 'a1b2c3d4-0001-4000-8000-000000000003', 'Buchanan Galleries', '220 Buchanan Street, Glasgow G1 2GF', 55.8636, -4.2518),
  ('b2c3d4e5-0001-4000-8000-000000000007', 'a1b2c3d4-0001-4000-8000-000000000003', 'Glasgow Royal Infirmary', '84 Castle Street, Glasgow G4 0SF', 55.8625, -4.2381);

-- Workers (auth_subject matches Keycloak user IDs — will be updated on first login)
INSERT INTO workers (id, auth_subject, first_name, last_name, email, phone) VALUES
  ('c3d4e5f6-0001-4000-8000-000000000001', 'admin.user', 'Admin', 'User', 'admin@sitesecurity.local', '+44 7700 900001'),
  ('c3d4e5f6-0001-4000-8000-000000000002', 'site.manager', 'Site', 'Manager', 'manager@sitesecurity.local', '+44 7700 900002'),
  ('c3d4e5f6-0001-4000-8000-000000000003', 'john.smith', 'John', 'Smith', 'john.smith@sitesecurity.local', '+44 7700 900003'),
  ('c3d4e5f6-0001-4000-8000-000000000004', 'jane.doe', 'Jane', 'Doe', 'jane.doe@sitesecurity.local', '+44 7700 900004'),
  ('c3d4e5f6-0001-4000-8000-000000000005', 'bob.wilson', 'Bob', 'Wilson', 'bob.wilson@sitesecurity.local', '+44 7700 900005'),
  ('c3d4e5f6-0001-4000-8000-000000000006', 'sarah.jones', 'Sarah', 'Jones', 'sarah.jones@sitesecurity.local', '+44 7700 900006'),
  ('c3d4e5f6-0001-4000-8000-000000000007', 'david.brown', 'David', 'Brown', 'david.brown@sitesecurity.local', '+44 7700 900007'),
  ('c3d4e5f6-0001-4000-8000-000000000008', 'emma.taylor', 'Emma', 'Taylor', 'emma.taylor@sitesecurity.local', '+44 7700 900008'),
  ('c3d4e5f6-0001-4000-8000-000000000009', 'michael.clark', 'Michael', 'Clark', 'michael.clark@sitesecurity.local', '+44 7700 900009'),
  ('c3d4e5f6-0001-4000-8000-000000000010', 'lisa.white', 'Lisa', 'White', 'lisa.white@sitesecurity.local', '+44 7700 900010');

-- Worker-Company memberships
INSERT INTO worker_companies (worker_id, company_id, role, status) VALUES
  -- Admin User is company_admin at Sentinel
  ('c3d4e5f6-0001-4000-8000-000000000001', 'a1b2c3d4-0001-4000-8000-000000000001', 'company_admin', 'active'),
  -- Site Manager is site_admin at Sentinel
  ('c3d4e5f6-0001-4000-8000-000000000002', 'a1b2c3d4-0001-4000-8000-000000000001', 'site_admin', 'active'),
  -- Workers at Sentinel
  ('c3d4e5f6-0001-4000-8000-000000000003', 'a1b2c3d4-0001-4000-8000-000000000001', 'worker', 'active'),
  ('c3d4e5f6-0001-4000-8000-000000000004', 'a1b2c3d4-0001-4000-8000-000000000001', 'worker', 'active'),
  ('c3d4e5f6-0001-4000-8000-000000000005', 'a1b2c3d4-0001-4000-8000-000000000001', 'worker', 'active'),
  ('c3d4e5f6-0001-4000-8000-000000000006', 'a1b2c3d4-0001-4000-8000-000000000001', 'worker', 'active'),
  -- Workers at Guardian
  ('c3d4e5f6-0001-4000-8000-000000000003', 'a1b2c3d4-0001-4000-8000-000000000002', 'worker', 'active'),
  ('c3d4e5f6-0001-4000-8000-000000000007', 'a1b2c3d4-0001-4000-8000-000000000002', 'worker', 'active'),
  ('c3d4e5f6-0001-4000-8000-000000000008', 'a1b2c3d4-0001-4000-8000-000000000002', 'worker', 'active'),
  -- Workers at Apex
  ('c3d4e5f6-0001-4000-8000-000000000009', 'a1b2c3d4-0001-4000-8000-000000000003', 'worker', 'active'),
  ('c3d4e5f6-0001-4000-8000-000000000010', 'a1b2c3d4-0001-4000-8000-000000000003', 'worker', 'active');

-- Certificates
INSERT INTO certificates (id, worker_id, name, issuing_body, certificate_number, issued_date, expiry_date) VALUES
  ('d4e5f6a7-0001-4000-8000-000000000001', 'c3d4e5f6-0001-4000-8000-000000000003', 'SIA Door Supervisor', 'Security Industry Authority', 'DS-2024-00123', '2024-01-15', '2027-01-15'),
  ('d4e5f6a7-0001-4000-8000-000000000002', 'c3d4e5f6-0001-4000-8000-000000000003', 'First Aid at Work', 'St John Ambulance', 'FA-2024-04567', '2024-03-01', '2027-03-01'),
  ('d4e5f6a7-0001-4000-8000-000000000003', 'c3d4e5f6-0001-4000-8000-000000000004', 'SIA Door Supervisor', 'Security Industry Authority', 'DS-2023-00456', '2023-06-20', '2026-06-20'),
  ('d4e5f6a7-0001-4000-8000-000000000004', 'c3d4e5f6-0001-4000-8000-000000000004', 'CCTV Operator', 'Security Industry Authority', 'CC-2024-00789', '2024-02-10', '2027-02-10'),
  ('d4e5f6a7-0001-4000-8000-000000000005', 'c3d4e5f6-0001-4000-8000-000000000005', 'SIA Security Guard', 'Security Industry Authority', 'SG-2024-00234', '2024-04-01', '2027-04-01'),
  ('d4e5f6a7-0001-4000-8000-000000000006', 'c3d4e5f6-0001-4000-8000-000000000006', 'SIA Close Protection', 'Security Industry Authority', 'CP-2023-00567', '2023-09-15', '2026-09-15'),
  ('d4e5f6a7-0001-4000-8000-000000000007', 'c3d4e5f6-0001-4000-8000-000000000007', 'SIA Door Supervisor', 'Security Industry Authority', 'DS-2024-00890', '2024-05-01', '2027-05-01'),
  ('d4e5f6a7-0001-4000-8000-000000000008', 'c3d4e5f6-0001-4000-8000-000000000008', 'SIA Security Guard', 'Security Industry Authority', 'SG-2023-00345', '2023-11-01', '2026-11-01'),
  -- Expired certificate
  ('d4e5f6a7-0001-4000-8000-000000000009', 'c3d4e5f6-0001-4000-8000-000000000009', 'SIA Door Supervisor', 'Security Industry Authority', 'DS-2021-00111', '2021-08-01', '2024-08-01'),
  ('d4e5f6a7-0001-4000-8000-000000000010', 'c3d4e5f6-0001-4000-8000-000000000010', 'SIA Security Guard', 'Security Industry Authority', 'SG-2024-00678', '2024-01-10', '2027-01-10');

-- Shift report templates
INSERT INTO shift_report_templates (id, company_id, name, fields) VALUES
  ('e5f6a7b8-0001-4000-8000-000000000001', 'a1b2c3d4-0001-4000-8000-000000000001', 'Standard Patrol Report', '[
    {"label": "Patrol route completed", "type": "boolean", "required": true},
    {"label": "Number of patrols", "type": "number", "required": true},
    {"label": "Incidents observed", "type": "text", "required": false},
    {"label": "Condition of premises", "type": "select", "required": true, "options": ["Good", "Fair", "Poor", "Requires attention"]},
    {"label": "Additional notes", "type": "textarea", "required": false}
  ]'),
  ('e5f6a7b8-0001-4000-8000-000000000002', 'a1b2c3d4-0001-4000-8000-000000000001', 'Incident Report', '[
    {"label": "Incident type", "type": "select", "required": true, "options": ["Trespass", "Theft", "Vandalism", "Anti-social behaviour", "Medical", "Fire", "Other"]},
    {"label": "Description", "type": "textarea", "required": true},
    {"label": "Police notified", "type": "boolean", "required": true},
    {"label": "Police reference number", "type": "text", "required": false},
    {"label": "Witnesses present", "type": "boolean", "required": true},
    {"label": "CCTV footage available", "type": "boolean", "required": true}
  ]'),
  ('e5f6a7b8-0001-4000-8000-000000000003', 'a1b2c3d4-0001-4000-8000-000000000002', 'Night Shift Handover', '[
    {"label": "All areas secured", "type": "boolean", "required": true},
    {"label": "Keys accounted for", "type": "boolean", "required": true},
    {"label": "Outstanding issues", "type": "textarea", "required": false},
    {"label": "Handover notes", "type": "textarea", "required": true}
  ]');

-- Shifts (using future dates relative to seed time)
INSERT INTO shifts (id, worksite_id, created_by, title, description, start_time, end_time, status) VALUES
  ('f6a7b8c9-0001-4000-8000-000000000001', 'b2c3d4e5-0001-4000-8000-000000000001', 'c3d4e5f6-0001-4000-8000-000000000002', 'Night Watch - Canary Wharf', 'Standard overnight security patrol', NOW() + INTERVAL '1 day', NOW() + INTERVAL '1 day 8 hours', 'open'),
  ('f6a7b8c9-0001-4000-8000-000000000002', 'b2c3d4e5-0001-4000-8000-000000000002', 'c3d4e5f6-0001-4000-8000-000000000002', 'Day Shift - Westminster', 'Daytime retail security', NOW() + INTERVAL '2 days', NOW() + INTERVAL '2 days 10 hours', 'open'),
  ('f6a7b8c9-0001-4000-8000-000000000003', 'b2c3d4e5-0001-4000-8000-000000000003', 'c3d4e5f6-0001-4000-8000-000000000002', 'Evening Security - Kings Cross', 'Evening patrol of development site', NOW() + INTERVAL '1 day 14 hours', NOW() + INTERVAL '1 day 22 hours', 'assigned'),
  ('f6a7b8c9-0001-4000-8000-000000000004', 'b2c3d4e5-0001-4000-8000-000000000004', 'c3d4e5f6-0001-4000-8000-000000000001', 'Weekend Cover - Arndale', 'Weekend shopping centre security', NOW() + INTERVAL '3 days', NOW() + INTERVAL '3 days 12 hours', 'open'),
  ('f6a7b8c9-0001-4000-8000-000000000005', 'b2c3d4e5-0001-4000-8000-000000000001', 'c3d4e5f6-0001-4000-8000-000000000002', 'Morning Patrol', 'Early morning building check', NOW() - INTERVAL '1 day', NOW() - INTERVAL '16 hours', 'completed');

-- Shift assignments
INSERT INTO shift_assignments (id, shift_id, worker_id, status, assigned_at, responded_at) VALUES
  ('a7b8c9d0-0001-4000-8000-000000000001', 'f6a7b8c9-0001-4000-8000-000000000003', 'c3d4e5f6-0001-4000-8000-000000000003', 'accepted', NOW() - INTERVAL '2 hours', NOW() - INTERVAL '1 hour'),
  ('a7b8c9d0-0001-4000-8000-000000000002', 'f6a7b8c9-0001-4000-8000-000000000001', 'c3d4e5f6-0001-4000-8000-000000000004', 'offered', NOW() - INTERVAL '30 minutes', NULL),
  ('a7b8c9d0-0001-4000-8000-000000000003', 'f6a7b8c9-0001-4000-8000-000000000005', 'c3d4e5f6-0001-4000-8000-000000000005', 'completed', NOW() - INTERVAL '2 days', NOW() - INTERVAL '2 days');

-- Shift reports
INSERT INTO shift_reports (id, shift_id, worker_id, template_id, data, submitted_at) VALUES
  ('b8c9d0e1-0001-4000-8000-000000000001', 'f6a7b8c9-0001-4000-8000-000000000005', 'c3d4e5f6-0001-4000-8000-000000000005', 'e5f6a7b8-0001-4000-8000-000000000001', '{
    "Patrol route completed": true,
    "Number of patrols": 4,
    "Incidents observed": "None",
    "Condition of premises": "Good",
    "Additional notes": "All clear. No issues to report."
  }', NOW() - INTERVAL '16 hours');

-- Location check-ins
INSERT INTO location_check_ins (id, worker_id, shift_id, latitude, longitude, recorded_at) VALUES
  ('c9d0e1f2-0001-4000-8000-000000000001', 'c3d4e5f6-0001-4000-8000-000000000005', 'f6a7b8c9-0001-4000-8000-000000000005', 51.5054, -0.0197, NOW() - INTERVAL '1 day'),
  ('c9d0e1f2-0001-4000-8000-000000000002', 'c3d4e5f6-0001-4000-8000-000000000005', 'f6a7b8c9-0001-4000-8000-000000000005', 51.5055, -0.0195, NOW() - INTERVAL '22 hours'),
  ('c9d0e1f2-0001-4000-8000-000000000003', 'c3d4e5f6-0001-4000-8000-000000000005', 'f6a7b8c9-0001-4000-8000-000000000005', 51.5053, -0.0198, NOW() - INTERVAL '20 hours'),
  ('c9d0e1f2-0001-4000-8000-000000000004', 'c3d4e5f6-0001-4000-8000-000000000005', 'f6a7b8c9-0001-4000-8000-000000000005', 51.5054, -0.0196, NOW() - INTERVAL '18 hours');

-- Alarms (one resolved, one recent)
INSERT INTO alarms (id, worker_id, shift_id, latitude, longitude, message, status, raised_at, acknowledged_at, resolved_at) VALUES
  ('d0e1f2a3-0001-4000-8000-000000000001', 'c3d4e5f6-0001-4000-8000-000000000005', 'f6a7b8c9-0001-4000-8000-000000000005', 51.5054, -0.0197, 'Suspicious individual near loading bay', 'resolved', NOW() - INTERVAL '20 hours', NOW() - INTERVAL '19 hours 50 minutes', NOW() - INTERVAL '19 hours');
