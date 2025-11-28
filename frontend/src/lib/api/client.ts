import { browser } from '$app/environment';
import { get } from 'svelte/store';
import type { ApiResponse, ErrorResponse } from '../types';
import { authStore } from '../stores/auth';

const API_BASE = '/api/v1';

async function request<T>(path: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
	const token = browser ? get(authStore).accessToken : '';
	const headers: Record<string, string> = {
		'Content-Type': 'application/json',
		...(options.headers as Record<string, string>),
	};
	if (token) {
		headers.Authorization = `Bearer ${token}`;
	}

	const resp = await fetch(`${API_BASE}${path}`, {
		...options,
		headers,
	});

	const body = await resp.json().catch(() => ({}));
	if (!resp.ok) {
		const err = body as ErrorResponse;
		throw new Error(err.error?.message ?? 'Request failed');
	}

	return body as ApiResponse<T>;
}

export const api = {
	get: request,
	post: <T>(path: string, body: unknown) =>
		request<T>(path, { method: 'POST', body: JSON.stringify(body) }),
};
