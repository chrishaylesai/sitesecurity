CREATE TYPE alarm_status AS ENUM ('raised', 'acknowledged', 'resolved');

CREATE TABLE alarms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    worker_id UUID NOT NULL REFERENCES workers(id),
    shift_id UUID REFERENCES shifts(id) ON DELETE SET NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    message TEXT,
    status alarm_status NOT NULL DEFAULT 'raised',
    raised_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    acknowledged_at TIMESTAMPTZ,
    resolved_at TIMESTAMPTZ
);

CREATE INDEX idx_alarms_worker_id ON alarms (worker_id);
CREATE INDEX idx_alarms_status ON alarms (status);
CREATE INDEX idx_alarms_raised_at ON alarms (raised_at);
