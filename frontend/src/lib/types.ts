export interface RepoInfo {
  id?: number;
  name?: string;
  url?: string;
}

export interface Project {
  id: number;
  user_id: number;
  name: string;
  slug: string;
  description?: string;
  project_type: string;
  github_repo?: RepoInfo;
}

export interface AgentRun {
  id: number;
  project_id: number;
  agent_type: string;
  trigger: string;
  status: string;
  results?: unknown;
  error_message?: string;
  started_at?: string;
  completed_at?: string;
  created_at?: string;
}

export interface User {
  id: number;
  github_id: number;
  username: string;
  email?: string | null;
  avatar_url?: string;
}

export interface ApiResponse<T> {
  data: T;
  meta?: Record<string, unknown>;
}

export interface ErrorResponse {
  error?: {
    message: string;
    code?: string;
    details?: unknown;
  };
}
