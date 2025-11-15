<script lang="ts">
	interface Props {
		padding?: 'none' | 'sm' | 'md' | 'lg';
		hover?: boolean;
		variant?: 'glass' | 'solid' | 'elevated';
		children?: any;
	}

	let { padding = 'md', hover = false, variant = 'glass', children }: Props = $props();

	const paddingClasses = {
		none: '',
		sm: 'p-4',
		md: 'p-6',
		lg: 'p-8'
	};

	const variantClasses = {
		glass: 'bg-white/60 dark:bg-slate-800/60 backdrop-blur-xl border border-slate-200/50 dark:border-slate-700/50',
		solid: 'bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700',
		elevated: 'bg-white dark:bg-slate-800 shadow-xl shadow-slate-200/50 dark:shadow-slate-950/50 border border-slate-100 dark:border-slate-900'
	};

	const hoverClass = hover ? 'hover:shadow-2xl hover:scale-[1.02] hover:border-slate-300 dark:hover:border-slate-600 transition-all duration-300 cursor-pointer' : '';

	const classes = `${variantClasses[variant]} rounded-2xl ${paddingClasses[padding]} ${hoverClass}`;
</script>

<div class={classes}>
	{@render children?.()}
</div>
