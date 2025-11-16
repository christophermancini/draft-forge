<script lang="ts">
	let { children } = $props();
</script>

<div class="min-h-screen bg-base-100">
	<nav class="navbar bg-base-200 shadow-lg">
		<div class="flex-1">
			<a href="/" class="btn btn-ghost text-xl">
				<span class="font-bold">DraftForge</span>
			</a>
		</div>
		<div class="flex-none">
			<a href="/projects" class="btn btn-ghost">Projects</a>
			<a href="/login" class="btn btn-primary">Login with GitHub</a>
		</div>
	</nav>

	<main class="container mx-auto px-4 py-8">
		{@render children()}
	</main>

	<footer class="footer footer-center p-4 bg-base-200 text-base-content mt-auto">
		<aside>
			<p>© 2025 DraftForge - Forge your draft. Keep your voice.</p>
		</aside>
	</footer>
</div>
