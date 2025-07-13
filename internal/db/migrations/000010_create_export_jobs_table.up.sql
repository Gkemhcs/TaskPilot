
CREATE TYPE export_job_status AS ENUM (
  'pending',
  'processing',
  'completed',
  'failed'
);

CREATE TYPE export_type AS ENUM (
  'project_excel',
  'task_excel'
);

CREATE TABLE export_jobs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  status export_job_status NOT NULL DEFAULT 'pending',
  export_type export_type NOT NULL,
  url TEXT,
  error_message TEXT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

