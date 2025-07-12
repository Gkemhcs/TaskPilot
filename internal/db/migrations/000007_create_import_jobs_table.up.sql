-- Create the ENUM type first
CREATE TYPE import_job_status AS ENUM ('pending', 'in_progress', 'completed', 'failed');

CREATE TYPE import_job_type AS ENUM ('project_excel', 'task_excel' );

-- Then create the table
CREATE TABLE import_jobs (
    id UUID PRIMARY KEY,
    file_path TEXT NOT NULL,
    importer_type import_job_type NOT NULL DEFAULT 'project_excel', -- e.g. 'project_excel'
    status import_job_status NOT NULL DEFAULT 'pending',
    error_message TEXT DEFAULT 'no error',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
