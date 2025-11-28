<script lang="ts">
	import '../app.css';
	import { authStore } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';

	let accessToken = '';
	let user = null;

	const unsubscribe = authStore.subscribe((state) => {
		accessToken = state.accessToken;
		user = state.user;
	});

	onMount(() => {
		const token = typeof window !== 'undefined' ? window.localStorage.getItem('df_access_token') : '';
		if (token) {
			authStore.setToken(token);
			authStore.fetchMe().catch(() => {
				authStore.reset();
			});
		}
		return () => unsubscribe();
	});

	function logout() {
		authStore.reset();
		goto('/');
	}
</script>

<div class="min-h-screen bg-[#F7F3EC] text-[#1C1C1C]">
	<nav class="flex items-center justify-between px-6 py-4 bg-white shadow-sm">
		<div class="flex items-center gap-3">
			<a href="/" class="text-lg font-semibold">DraftForge</a>
			<a href="/projects" class="text-sm text-[#2F8F9D] hover:underline">Projects</a>
		</div>
		<div class="flex items-center gap-3">
			{#if user}
				<span class="text-sm text-gray-700">Hello, {user.username}</span>
				<button class="text-sm text-[#B26E63] hover:underline" on:click={logout}>Logout</button>
			{:else}
				<a href="/login" class="text-sm text-[#2F8F9D] hover:underline">Login with GitHub</a>
			{/if}
		</div>
	</nav>

	<main class="max-w-5xl mx-auto px-6 py-8">
		<slot />
	</main>
</div>
