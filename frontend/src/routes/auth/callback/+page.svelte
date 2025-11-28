<script lang="ts">
import { onMount } from 'svelte';
import { page } from '$app/stores';
import { goto } from '$app/navigation';
import { api } from '$lib/api/client';
import { authStore } from '$lib/stores/auth';

	let error = '';

	onMount(async () => {
		const url = new URL($page.url);
		const code = url.searchParams.get('code');
		const state = url.searchParams.get('state');

		if (!code || !state) {
			error = 'Missing code or state';
			return;
		}

		try {
			const resp = await api.get<{ data: { user: any; token: { access_token: string } } }>(
				`/auth/github/callback?code=${encodeURIComponent(code)}&state=${encodeURIComponent(state)}`
			);
			authStore.setToken(resp.data.token.access_token);
			authStore.setUser(resp.data.user);
			goto('/projects');
		} catch (e: any) {
			error = e?.message ?? 'Login failed';
		}
	});
</script>

{#if error}
	<div class="text-red-600">{error}</div>
{:else}
	<div class="text-gray-700">Completing login...</div>
{/if}
