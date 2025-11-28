<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { api } from '$lib/api/client';
	import type { ApiResponse, AgentRun, Project } from '$lib/types';

	let project: Project | null = null;
	let runs: AgentRun[] = [];
	let loading = true;
	let error = '';

	async function load() {
		loading = true;
		error = '';
		const slug = $page.params.slug;
		try {
			const projectsResp = await api.get<ApiResponse<Project[]>>('/projects');
			project = projectsResp.data.find((p) => p.slug === slug) || null;
			if (project) {
				const runsResp = await api.get<ApiResponse<AgentRun[]>>(`/projects/${project.id}/agents/runs`);
				runs = runsResp.data;
			} else {
				error = 'Project not found';
			}
		} catch (e: any) {
			error = e?.message ?? 'Failed to load project';
		} finally {
			loading = false;
		}
	}

	async function runAgent() {
		if (!project) return;
		error = '';
		try {
			const resp = await api.post<ApiResponse<AgentRun>>(`/projects/${project.id}/agents/run`, {
				agent_type: 'continuity',
				trigger: 'manual',
			});
			runs = [resp.data, ...runs];
		} catch (e: any) {
			error = e?.message ?? 'Failed to queue agent';
		}
	}

	onMount(load);
</script>

<div class="space-y-6">
	{#if loading}
		<p class="text-gray-700">Loading...</p>
	{:else if error}
		<p class="text-red-600">{error}</p>
	{:else if project}
		<section class="bg-white rounded-lg shadow-sm p-6">
			<div class="flex items-center justify-between">
				<div>
					<h1 class="text-xl font-semibold">{project.name}</h1>
					<p class="text-sm text-gray-700">{project.project_type}</p>
					{#if project.github_repo?.url}
						<a class="text-sm text-[#2F8F9D]" href={project.github_repo.url}>{project.github_repo.url}</a>
					{/if}
				</div>
				<button class="px-4 py-2 bg-[#2F8F9D] text-white rounded hover:bg-[#287a86]" on:click={runAgent}>
					Run ContinuityBot
				</button>
			</div>
		</section>

		<section class="bg-white rounded-lg shadow-sm p-6">
			<h2 class="text-lg font-semibold mb-3">Agent Runs</h2>
			{#if runs.length === 0}
				<p class="text-gray-700">No runs yet.</p>
			{:else}
				<div class="space-y-3">
					{#each runs as run}
						<div class="border border-[#E3DFD7] rounded p-3">
							<div class="flex justify-between text-sm">
								<span class="font-semibold">{run.agent_type}</span>
								<span class="text-gray-700">{run.status}</span>
							</div>
							{#if run.error_message}
								<p class="text-sm text-red-600 mt-1">{run.error_message}</p>
							{/if}
							{#if run.created_at}
								<p class="text-xs text-gray-600 mt-1">{run.created_at}</p>
							{/if}
						</div>
					{/each}
				</div>
			{/if}
		</section>
	{/if}
</div>
