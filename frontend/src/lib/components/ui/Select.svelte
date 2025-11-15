<script lang="ts">
	interface Option {
		value: string;
		label: string;
		disabled?: boolean;
	}

	interface Props {
		label?: string;
		options: Option[];
		value: string;
		placeholder?: string;
		disabled?: boolean;
		required?: boolean;
		error?: string;
		onChange?: (value: string) => void;
	}

	let {
		label,
		options = [],
		value = $bindable(''),
		placeholder = 'Select an option',
		disabled = false,
		required = false,
		error,
		onChange
	}: Props = $props();

	function handleChange(e: Event) {
		const target = e.target as HTMLSelectElement;
		value = target.value;
		if (onChange) {
			onChange(value);
		}
	}
</script>

<div class="select-wrapper">
	{#if label}
		<label for="select-{label.replace(/\s+/g, '-').toLowerCase()}" class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
			{label}
			{#if required}
				<span class="text-red-500 ml-1">*</span>
			{/if}
		</label>
	{/if}

	<div class="relative">
		<select
			id="select-{label?.replace(/\s+/g, '-').toLowerCase() || 'default'}"
			{value}
			onchange={handleChange}
			{disabled}
			{required}
			class="w-full px-4 py-2.5 pr-10 rounded-xl border-2 appearance-none cursor-pointer transition-all duration-200
				{error
					? 'border-red-400 dark:border-red-500 focus:border-red-500 focus:ring-4 focus:ring-red-500/20'
					: 'border-slate-200 dark:border-slate-700 focus:border-blue-500 dark:focus:border-blue-400 focus:ring-4 focus:ring-blue-500/20 dark:focus:ring-blue-400/20'}
				bg-white/80 dark:bg-slate-800/80 backdrop-blur-sm
				text-slate-900 dark:text-slate-100
				focus:outline-none
				disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-slate-100 dark:disabled:bg-slate-900
				shadow-sm hover:shadow-md focus:shadow-lg"
		>
			{#if placeholder}
				<option value="" disabled selected={!value}>
					{placeholder}
				</option>
			{/if}

			{#each options as option}
				<option value={option.value} disabled={option.disabled}>
					{option.label}
				</option>
			{/each}
		</select>

		<!-- Custom dropdown arrow -->
		<div
			class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none text-slate-500 dark:text-slate-400"
		>
			<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" stroke-width="2.5">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					d="M19 9l-7 7-7-7"
				/>
			</svg>
		</div>
	</div>

	{#if error}
		<p class="mt-1.5 text-sm text-red-600 dark:text-red-400 flex items-center gap-1">
			<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
				<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
			</svg>
			{error}
		</p>
	{/if}
</div>

<style>
	select {
		background-image: none;
	}
</style>
