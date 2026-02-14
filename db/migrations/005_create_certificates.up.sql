CREATE TABLE certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    issuing_body VARCHAR(255),
    certificate_number VARCHAR(255),
    issued_date DATE,
    expiry_date DATE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_certificates_worker_id ON certificates (worker_id);
CREATE INDEX idx_certificates_expiry_date ON certificates (expiry_date);
