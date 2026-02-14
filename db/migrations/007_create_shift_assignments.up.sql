CREATE TYPE assignment_status AS ENUM ('offered', 'accepted', 'declined', 'completed');

CREATE TABLE shift_assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shift_id UUID NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    worker_id UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    status assignment_status NOT NULL DEFAULT 'offered',
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    responded_at TIMESTAMPTZ,
    UNIQUE (shift_id, worker_id)
);

CREATE INDEX idx_shift_assignments_worker_id ON shift_assignments (worker_id);
CREATE INDEX idx_shift_assignments_shift_id ON shift_assignments (shift_id);
