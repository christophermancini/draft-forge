<script lang="ts">
	import { api } from '$lib/api/client';

	async function startLogin() {
		const resp = await api.get<{ data: { auth_url: string } }>('/auth/github/start');
		window.location.href = resp.data.auth_url;
	}
</script>

<div class="space-y-4">
	<h1 class="text-xl font-semibold">Login with GitHub</h1>
	<p class="text-gray-700">You will be redirected to GitHub to authorize DraftForge.</p>
	<button class="px-4 py-2 bg-[#2F8F9D] text-white rounded hover:bg-[#287a86]" on:click={startLogin}>
		Continue with GitHub
	</button>
</div>
