CREATE TYPE worker_role AS ENUM ('worker', 'company_admin', 'site_admin');
CREATE TYPE membership_status AS ENUM ('active', 'inactive');

CREATE TABLE worker_companies (
    worker_id UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    role worker_role NOT NULL DEFAULT 'worker',
    status membership_status NOT NULL DEFAULT 'active',
    joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (worker_id, company_id)
);

CREATE INDEX idx_worker_companies_company_id ON worker_companies (company_id);
