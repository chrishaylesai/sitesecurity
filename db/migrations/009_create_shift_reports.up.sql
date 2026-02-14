CREATE TABLE shift_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shift_id UUID NOT NULL REFERENCES shifts(id) ON DELETE CASCADE,
    worker_id UUID NOT NULL REFERENCES workers(id),
    template_id UUID REFERENCES shift_report_templates(id),
    data JSONB NOT NULL DEFAULT '{}',
    submitted_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shift_reports_shift_id ON shift_reports (shift_id);
CREATE INDEX idx_shift_reports_worker_id ON shift_reports (worker_id);
