


-- Step 1: Create ENUM type for color (run this first)
CREATE TYPE project_color AS ENUM ('GREEN', 'YELLOW', 'RED');

-- Step 2: Now create the table using that ENUM
CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    color project_color, -- Optional for frontend UI
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
