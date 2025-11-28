<script lang="ts">
	import { onMount } from 'svelte';
	import { api } from '$lib/api/client';
	import type { ApiResponse, Project } from '$lib/types';

	let projects: Project[] = [];
	let loading = true;
	let error = '';

	let name = '';
	let projectType = 'novel';
	let template = 'novel';
	let description = '';
	let useGitHub = true;
	let githubOwner = '';
	let repoMeta = '';
	let scaffoldPath = '';

	async function loadProjects() {
		loading = true;
		error = '';
		try {
			const resp = await api.get<ApiResponse<Project[]>>('/projects');
			projects = resp.data;
		} catch (e: any) {
			error = e?.message ?? 'Failed to load projects';
		} finally {
			loading = false;
		}
	}

	async function createProject() {
		error = '';
		repoMeta = '';
		scaffoldPath = '';
		try {
			const resp = await api.post<ApiResponse<Project>>('/projects', {
				name,
				project_type: projectType,
				template,
				description,
				use_github: useGitHub,
				github_owner: githubOwner,
			});
			projects = [resp.data, ...projects];
			if (resp.meta?.repo_url) repoMeta = String(resp.meta.repo_url);
			if (resp.meta?.scaffold_path) scaffoldPath = String(resp.meta.scaffold_path);
			name = '';
			description = '';
		} catch (e: any) {
			error = e?.message ?? 'Failed to create project';
		}
	}

	onMount(loadProjects);
</script>

<div class="space-y-6">
	<section class="bg-white rounded-lg shadow-sm p-6">
		<h1 class="text-xl font-semibold mb-4">Create Project</h1>
		{#if error}
			<div class="text-sm text-red-600 mb-3">{error}</div>
		{/if}
		<div class="grid md:grid-cols-2 gap-4">
			<label class="space-y-1">
				<span class="text-sm text-gray-700">Name</span>
				<input class="input input-bordered w-full" bind:value={name} required />
			</label>
			<label class="space-y-1">
				<span class="text-sm text-gray-700">Project Type</span>
				<select class="select select-bordered w-full" bind:value={projectType}>
					<option value="novel">Novel</option>
					<option value="screenplay">Screenplay</option>
					<option value="technical-book">Technical Book</option>
					<option value="non-fiction">Non-Fiction</option>
				</select>
			</label>
			<label class="space-y-1">
				<span class="text-sm text-gray-700">Template</span>
				<select class="select select-bordered w-full" bind:value={template}>
					<option value="novel">Novel</option>
				</select>
			</label>
			<label class="space-y-1">
				<span class="text-sm text-gray-700">Description</span>
				<input class="input input-bordered w-full" bind:value={description} />
			</label>
			<label class="flex items-center gap-2">
				<input type="checkbox" class="checkbox" bind:checked={useGitHub} />
				<span class="text-sm text-gray-700">Create GitHub repo</span>
			</label>
			{#if useGitHub}
				<label class="space-y-1">
					<span class="text-sm text-gray-700">GitHub Owner (optional for org)</span>
					<input class="input input-bordered w-full" bind:value={githubOwner} />
				</label>
			{/if}
		</div>
		<button class="mt-4 px-4 py-2 bg-[#2F8F9D] text-white rounded hover:bg-[#287a86]" on:click|preventDefault={createProject}>
			Create
		</button>
		{#if repoMeta}
			<p class="text-sm text-gray-700 mt-2">Repo: <a class="text-[#2F8F9D]" href={repoMeta}>{repoMeta}</a></p>
		{/if}
		{#if scaffoldPath}
			<p class="text-sm text-gray-700">Scaffold path: {scaffoldPath}</p>
		{/if}
	</section>

	<section class="bg-white rounded-lg shadow-sm p-6">
		<div class="flex items-center justify-between mb-3">
			<h2 class="text-lg font-semibold">Projects</h2>
			{#if loading}
				<span class="text-sm text-gray-600">Loading...</span>
			{/if}
		</div>
		{#if projects.length === 0 && !loading}
			<p class="text-gray-700">No projects yet.</p>
		{:else}
			<div class="grid gap-3">
				{#each projects as project}
					<div class="border border-[#E3DFD7] rounded-lg p-4 flex justify-between items-center">
						<div>
							<div class="font-semibold">{project.name}</div>
							<div class="text-sm text-gray-700">{project.project_type}</div>
							{#if project.github_repo?.url}
								<a class="text-sm text-[#2F8F9D]" href={project.github_repo.url}>{project.github_repo.url}</a>
							{/if}
						</div>
						<a class="text-sm text-[#2F8F9D]" href={`/projects/${project.slug}`}>Open</a>
					</div>
				{/each}
			</div>
		{/if}
	</section>
</div>
