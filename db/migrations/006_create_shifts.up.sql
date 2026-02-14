CREATE TYPE shift_status AS ENUM ('open', 'assigned', 'in_progress', 'completed', 'cancelled');

CREATE TABLE shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    worksite_id UUID NOT NULL REFERENCES worksites(id) ON DELETE CASCADE,
    created_by UUID NOT NULL REFERENCES workers(id),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    status shift_status NOT NULL DEFAULT 'open',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shifts_worksite_id ON shifts (worksite_id);
CREATE INDEX idx_shifts_status ON shifts (status);
CREATE INDEX idx_shifts_start_time ON shifts (start_time);
