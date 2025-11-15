<script lang="ts">
	interface Props {
		label?: string;
		value: string;
		placeholder?: string;
		rows?: number;
		disabled?: boolean;
		required?: boolean;
		maxLength?: number;
		error?: string;
		helperText?: string;
		resize?: 'none' | 'vertical' | 'horizontal' | 'both';
		onChange?: (value: string) => void;
	}

	let {
		label,
		value = $bindable(''),
		placeholder,
		rows = 4,
		disabled = false,
		required = false,
		maxLength,
		error,
		helperText,
		resize = 'vertical',
		onChange
	}: Props = $props();

	function handleInput(e: Event) {
		const target = e.target as HTMLTextAreaElement;
		value = target.value;
		if (onChange) {
			onChange(value);
		}
	}

	const charCount = $derived(value.length);
</script>

<div class="textarea-wrapper">
	{#if label}
		<label for="textarea-{label.replace(/\s+/g, '-').toLowerCase()}" class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
			{label}
			{#if required}
				<span class="text-red-500 ml-1">*</span>
			{/if}
		</label>
	{/if}

	<textarea
		id="textarea-{label?.replace(/\s+/g, '-').toLowerCase() || 'default'}"
		{value}
		oninput={handleInput}
		{placeholder}
		{rows}
		{disabled}
		{required}
		maxlength={maxLength}
		style="resize: {resize};"
		class="w-full px-4 py-2.5 rounded-xl border-2 transition-all duration-200
			{error
				? 'border-red-400 dark:border-red-500 focus:border-red-500 focus:ring-4 focus:ring-red-500/20'
				: 'border-slate-200 dark:border-slate-700 focus:border-blue-500 dark:focus:border-blue-400 focus:ring-4 focus:ring-blue-500/20 dark:focus:ring-blue-400/20'}
			bg-white/80 dark:bg-slate-800/80 backdrop-blur-sm
			text-slate-900 dark:text-slate-100
			placeholder:text-slate-400 dark:placeholder:text-slate-500
			focus:outline-none
			disabled:opacity-50 disabled:cursor-not-allowed disabled:bg-slate-100 dark:disabled:bg-slate-900
			shadow-sm hover:shadow-md focus:shadow-lg"
	></textarea>

	<div class="flex items-center justify-between mt-1.5 text-xs">
		<div>
			{#if error}
				<p class="text-red-600 dark:text-red-400 flex items-center gap-1">
					<svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
					</svg>
					{error}
				</p>
			{:else if helperText}
				<p class="text-slate-500 dark:text-slate-400">
					{helperText}
				</p>
			{/if}
		</div>

		{#if maxLength}
			<p
				class="text-slate-500 dark:text-slate-400 font-medium {charCount > maxLength * 0.9
					? 'text-amber-600 dark:text-amber-400'
					: ''} {charCount === maxLength ? 'text-red-600 dark:text-red-400' : ''}"
			>
				{charCount}/{maxLength}
			</p>
		{/if}
	</div>
</div>
