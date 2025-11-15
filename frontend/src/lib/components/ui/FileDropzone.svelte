<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		accept?: string;
		maxSizeMB?: number;
		onFileSelect: (file: File) => void;
		disabled?: boolean;
	}

	let { accept = '.csv', maxSizeMB = 10, onFileSelect, disabled = false }: Props = $props();

	let isDragging = $state(false);
	let error = $state('');
	let fileInputElement: HTMLInputElement;

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		if (!disabled) {
			isDragging = true;
		}
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		error = '';

		if (disabled) return;

		const files = e.dataTransfer?.files;
		if (files && files.length > 0) {
			validateAndSelectFile(files[0]);
		}
	}

	function handleFileInput(e: Event) {
		error = '';
		const target = e.target as HTMLInputElement;
		const files = target.files;

		if (files && files.length > 0) {
			validateAndSelectFile(files[0]);
		}
	}

	function validateAndSelectFile(file: File) {
		// Check file type
		const extension = `.${file.name.split('.').pop()?.toLowerCase()}`;
		const acceptedExtensions = accept.split(',').map((ext) => ext.trim().toLowerCase());

		if (!acceptedExtensions.includes(extension)) {
			error = `Invalid file type. Please upload ${accept} file`;
			return;
		}

		// Check file size
		const sizeMB = file.size / (1024 * 1024);
		if (sizeMB > maxSizeMB) {
			error = `File too large. Maximum size is ${maxSizeMB}MB`;
			return;
		}

		onFileSelect(file);
	}

	function handleClick() {
		if (!disabled) {
			fileInputElement.click();
		}
	}
</script>

<div
	class="relative border-2 border-dashed rounded-xl transition-all duration-200 cursor-pointer
		{isDragging
			? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
			: error
				? 'border-red-400 dark:border-red-500 bg-red-50 dark:bg-red-900/10'
				: 'border-slate-300 dark:border-slate-600 hover:border-blue-400 dark:hover:border-blue-500'}
		{disabled ? 'opacity-50 cursor-not-allowed' : ''}"
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	ondrop={handleDrop}
	onclick={handleClick}
	onkeydown={(e) => e.key === 'Enter' && handleClick()}
	role="button"
	tabindex="0"
	aria-label="Upload file"
>
	<input
		type="file"
		{accept}
		{disabled}
		bind:this={fileInputElement}
		onchange={handleFileInput}
		class="hidden"
		aria-label="File input"
	/>

	<div class="p-12 text-center">
		<div class="mb-4">
			{#if isDragging}
				<Icon icon="mdi:file-upload" width="64" class="mx-auto text-blue-500 animate-bounce" />
			{:else if error}
				<Icon icon="mdi:alert-circle" width="64" class="mx-auto text-red-500" />
			{:else}
				<Icon icon="mdi:file-delimited" width="64" class="mx-auto text-slate-400" />
			{/if}
		</div>

		{#if error}
			<p class="text-red-600 dark:text-red-400 font-medium mb-2">{error}</p>
			<p class="text-sm text-slate-600 dark:text-slate-400">Try again with a valid file</p>
		{:else if isDragging}
			<p class="text-blue-600 dark:text-blue-400 font-semibold text-lg mb-2">
				Drop your CSV file here
			</p>
		{:else}
			<p class="text-lg font-semibold text-slate-800 dark:text-slate-200 mb-2">
				Drag & Drop CSV File
			</p>
			<p class="text-sm text-slate-600 dark:text-slate-400 mb-4">or click to browse</p>
			<div class="inline-flex items-center gap-2 px-4 py-2 bg-blue-100 dark:bg-blue-900/30 rounded-lg text-blue-700 dark:text-blue-300 text-sm font-medium">
				<Icon icon="mdi:file-find" width="20" />
				<span>Select {accept} file (max {maxSizeMB}MB)</span>
			</div>
		{/if}
	</div>
</div>
