CREATE TABLE location_check_ins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    shift_id UUID REFERENCES shifts(id) ON DELETE SET NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    recorded_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_location_check_ins_worker_id ON location_check_ins (worker_id);
CREATE INDEX idx_location_check_ins_shift_id ON location_check_ins (shift_id);
CREATE INDEX idx_location_check_ins_recorded_at ON location_check_ins (recorded_at);
