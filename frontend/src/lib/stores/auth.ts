import { get, writable } from 'svelte/store';
import type { ApiResponse, User } from '../types';
import { api } from '../api/client';

type AuthState = {
	accessToken: string;
	user: User | null;
};

const initial: AuthState = {
	accessToken: '',
	user: null,
};

function createAuthStore() {
	const { subscribe, set, update } = writable<AuthState>(initial);

	return {
		subscribe,
		setToken: (token: string) => {
			if (typeof window !== 'undefined') {
				window.localStorage.setItem('df_access_token', token);
			}
			update((s) => ({ ...s, accessToken: token }));
		},
		setUser: (user: User | null) => update((s) => ({ ...s, user })),
		reset: () => {
			if (typeof window !== 'undefined') {
				window.localStorage.removeItem('df_access_token');
			}
			set(initial);
		},
		fetchMe: async () => {
			const me = await api.get<ApiResponse<User>>('/me');
			set({ accessToken: getAccessToken(), user: me.data });
		},
	};
}

function getAccessToken(): string {
	if (typeof window === 'undefined') return '';
	return window.localStorage.getItem('df_access_token') || '';
}

export const authStore = createAuthStore();
