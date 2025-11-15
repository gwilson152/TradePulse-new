<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		accept?: string;
		multiple?: boolean;
		maxSize?: number; // in MB
		label?: string;
		preview?: boolean;
		value?: File[];
		onChange: (files: File[]) => void;
	}

	let {
		accept = '*',
		multiple = false,
		maxSize = 10,
		label = 'Upload Files',
		preview = true,
		value = [],
		onChange
	}: Props = $props();

	let files = $state<File[]>(value);
	let dragActive = $state(false);
	let error = $state('');
	let fileInput: HTMLInputElement;

	// Preview URLs for images
	let previewUrls = $state<string[]>([]);

	function handleFiles(newFiles: FileList | null) {
		if (!newFiles) return;

		error = '';
		const fileArray = Array.from(newFiles);

		// Validate file sizes
		const oversizedFiles = fileArray.filter((file) => file.size > maxSize * 1024 * 1024);
		if (oversizedFiles.length > 0) {
			error = `Some files exceed the ${maxSize}MB limit`;
			return;
		}

		if (multiple) {
			files = [...files, ...fileArray];
		} else {
			files = fileArray.slice(0, 1);
		}

		// Generate preview URLs for images
		generatePreviews();

		onChange(files);
	}

	function generatePreviews() {
		// Revoke old URLs to prevent memory leaks
		previewUrls.forEach((url) => URL.revokeObjectURL(url));

		previewUrls = files
			.filter((file) => file.type.startsWith('image/'))
			.map((file) => URL.createObjectURL(file));
	}

	function removeFile(index: number) {
		// Revoke preview URL if exists
		if (previewUrls[index]) {
			URL.revokeObjectURL(previewUrls[index]);
		}

		files = files.filter((_, i) => i !== index);
		previewUrls = previewUrls.filter((_, i) => i !== index);

		onChange(files);
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		dragActive = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		dragActive = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		dragActive = false;
		handleFiles(e.dataTransfer?.files || null);
	}

	function formatFileSize(bytes: number): string {
		if (bytes < 1024) return bytes + ' B';
		if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
		return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
	}

	function getFileIcon(file: File): string {
		if (file.type.startsWith('image/')) return 'mdi:image';
		if (file.type.startsWith('audio/')) return 'mdi:music';
		if (file.type.startsWith('video/')) return 'mdi:video';
		if (file.type.includes('pdf')) return 'mdi:file-pdf';
		return 'mdi:file';
	}
</script>

<div class="file-upload">
	{#if label}
		<label for="file-input-{label.replace(/\s+/g, '-').toLowerCase()}" class="block text-sm font-medium mb-2 text-surface-700 dark:text-surface-300">
			{label}
		</label>
	{/if}

	<!-- Drop zone -->
	<div
		class="drop-zone {dragActive ? 'drag-active' : ''}"
		ondragover={handleDragOver}
		ondragleave={handleDragLeave}
		ondrop={handleDrop}
		onclick={() => fileInput.click()}
		onkeydown={(e) => e.key === 'Enter' && fileInput.click()}
		role="button"
		tabindex="0"
		aria-label="File upload zone"
	>
		<Icon icon="mdi:cloud-upload" class="text-4xl text-surface-400 dark:text-surface-500 mb-2" />
		<p class="text-sm font-medium text-surface-700 dark:text-surface-300 mb-1">
			Click to upload or drag and drop
		</p>
		<p class="text-xs text-surface-500 dark:text-surface-400">
			{accept !== '*' ? `Accepts: ${accept}` : 'Any file type'} • Max {maxSize}MB
			{#if multiple}• Multiple files{/if}
		</p>

		<input
			bind:this={fileInput}
			id="file-input-{label?.replace(/\s+/g, '-').toLowerCase() || 'default'}"
			type="file"
			{accept}
			{multiple}
			onchange={(e) => handleFiles(e.currentTarget.files)}
			class="hidden"
		/>
	</div>

	{#if error}
		<div class="mt-2 text-sm text-loss-600 dark:text-loss-400 flex items-center gap-1">
			<Icon icon="mdi:alert-circle" />
			{error}
		</div>
	{/if}

	<!-- File list -->
	{#if files.length > 0}
		<div class="file-list mt-4 space-y-2">
			{#each files as file, index}
				<div
					class="file-item flex items-center gap-3 p-3 bg-surface-100 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700"
				>
					{#if preview && file.type.startsWith('image/') && previewUrls[index]}
						<img
							src={previewUrls[index]}
							alt={file.name}
							class="w-12 h-12 object-cover rounded"
						/>
					{:else}
						<div
							class="w-12 h-12 flex items-center justify-center bg-surface-200 dark:bg-surface-700 rounded"
						>
							<Icon icon={getFileIcon(file)} class="text-2xl text-surface-600 dark:text-surface-400" />
						</div>
					{/if}

					<div class="flex-1 min-w-0">
						<p class="text-sm font-medium text-surface-900 dark:text-surface-100 truncate">
							{file.name}
						</p>
						<p class="text-xs text-surface-600 dark:text-surface-400">
							{formatFileSize(file.size)}
						</p>
					</div>

					<button
						type="button"
						onclick={() => removeFile(index)}
						class="p-2 hover:bg-surface-200 dark:hover:bg-surface-700 rounded transition-colors"
						title="Remove file"
					>
						<Icon icon="mdi:close" class="text-surface-600 dark:text-surface-400" />
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>

<style>
	.drop-zone {
		border: 2px dashed rgb(var(--color-surface-300));
		border-radius: 0.5rem;
		padding: 2rem;
		text-align: center;
		cursor: pointer;
		transition: all 0.2s ease;
		background-color: rgb(var(--color-surface-50));
	}

	:global(.dark) .drop-zone {
		background-color: rgb(var(--color-surface-800));
		border-color: rgb(var(--color-surface-600));
	}

	.drop-zone:hover {
		border-color: rgb(var(--color-primary-500));
		background-color: rgb(var(--color-primary-50));
	}

	:global(.dark) .drop-zone:hover {
		background-color: rgb(var(--color-primary-900) / 0.2);
	}

	.drop-zone.drag-active {
		border-color: rgb(var(--color-primary-500));
		background-color: rgb(var(--color-primary-100));
		transform: scale(1.02);
	}

	:global(.dark) .drop-zone.drag-active {
		background-color: rgb(var(--color-primary-900) / 0.3);
	}

	.file-item {
		transition: all 0.2s ease;
	}

	.file-item:hover {
		transform: translateX(4px);
	}
</style>
