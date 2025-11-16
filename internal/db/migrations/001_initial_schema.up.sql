-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    github_id INTEGER UNIQUE NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255),
    avatar_url TEXT,
    access_token TEXT,
    refresh_token TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_github_id ON users(github_id);
CREATE INDEX idx_users_username ON users(username);

-- Create projects table
CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    description TEXT,
    project_type VARCHAR(50) NOT NULL, -- 'novel', 'screenplay', 'technical-book', 'non-fiction'
    github_repo_id INTEGER,
    github_repo_name VARCHAR(255),
    github_repo_url TEXT,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, slug)
);

CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_github_repo_id ON projects(github_repo_id);

-- Create ai_credits table
CREATE TABLE IF NOT EXISTS ai_credits (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    credits_remaining INTEGER NOT NULL DEFAULT 0,
    credits_total INTEGER NOT NULL DEFAULT 0,
    subscription_tier VARCHAR(50) DEFAULT 'free', -- 'free', 'creator', 'professional', 'enterprise'
    subscription_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Create ai_usage_log table
CREATE TABLE IF NOT EXISTS ai_usage_log (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
    agent_type VARCHAR(50) NOT NULL, -- 'continuity', 'style', 'timeline', 'fact'
    model_name VARCHAR(100) NOT NULL,
    tokens_used INTEGER NOT NULL,
    cost_cents INTEGER NOT NULL, -- Store cost in cents
    status VARCHAR(20) NOT NULL, -- 'success', 'failed', 'pending'
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_ai_usage_user_id ON ai_usage_log(user_id);
CREATE INDEX idx_ai_usage_project_id ON ai_usage_log(project_id);
CREATE INDEX idx_ai_usage_created_at ON ai_usage_log(created_at);

-- Create agent_runs table
CREATE TABLE IF NOT EXISTS agent_runs (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    agent_type VARCHAR(50) NOT NULL,
    trigger VARCHAR(50) NOT NULL, -- 'manual', 'commit', 'pr', 'scheduled'
    status VARCHAR(20) NOT NULL, -- 'queued', 'running', 'completed', 'failed'
    results JSONB,
    error_message TEXT,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_agent_runs_project_id ON agent_runs(project_id);
CREATE INDEX idx_agent_runs_status ON agent_runs(status);

-- Create project_stats table
CREATE TABLE IF NOT EXISTS project_stats (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    word_count INTEGER DEFAULT 0,
    chapter_count INTEGER DEFAULT 0,
    last_commit_sha VARCHAR(255),
    last_commit_at TIMESTAMP WITH TIME ZONE,
    stats_data JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id)
);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ai_credits_updated_at BEFORE UPDATE ON ai_credits
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_project_stats_updated_at BEFORE UPDATE ON project_stats
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
