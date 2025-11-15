<script lang="ts">
	interface Props {
		type?: string;
		value?: string | number;
		placeholder?: string;
		label?: string;
		error?: string;
		required?: boolean;
		disabled?: boolean;
		name?: string;
		step?: string;
		min?: string;
		max?: string;
		onChange?: (value: string) => void;
		oninput?: (event: Event) => void;
		class?: string;
	}

	let {
		type = 'text',
		value = $bindable(''),
		placeholder = '',
		label,
		error,
		required = false,
		disabled = false,
		name,
		step,
		min,
		max,
		onChange,
		oninput,
		class: className = ''
	}: Props = $props();

	const inputId = $derived(label ? `input-${label.replace(/\s+/g, '-').toLowerCase()}` : undefined);

	function handleInput(e: Event) {
		if (oninput) oninput(e);
		if (onChange) onChange((e.target as HTMLInputElement).value);
	}
</script>

<div class="input-wrapper {className}">
	{#if label}
		<label for={inputId} class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
			{label}
			{#if required}
				<span class="text-red-500 ml-1">*</span>
			{/if}
		</label>
	{/if}
	<input
		id={inputId}
		class="w-full px-4 py-2.5 rounded-xl border-2 transition-all duration-200
			{error
				? 'border-red-400 dark:border-red-500 focus:border-red-500 focus:ring-4 focus:ring-red-500/20'
				: 'border-slate-200 dark:border-slate-700 focus:border-blue-500 dark:focus:border-blue-400 focus:ring-4 focus:ring-blue-500/20 dark:focus:ring-blue-400/20'}
			bg-white/80 dark:bg-slate-800/80 backdrop-blur-sm
			text-slate-900 dark:text-slate-100
			placeholder:text-slate-400 dark:placeholder:text-slate-500
			disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-slate-100 dark:disabled:bg-slate-900
			shadow-sm hover:shadow-md focus:shadow-lg"
		{type}
		bind:value
		{placeholder}
		{required}
		{disabled}
		{name}
		{step}
		{min}
		{max}
		oninput={handleInput}
	/>
	{#if error}
		<p class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-center gap-1">
			<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
			</svg>
			{error}
		</p>
	{/if}
</div>
