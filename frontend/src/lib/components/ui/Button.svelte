<script lang="ts">
	interface Props {
		variant?: 'filled' | 'ghost' | 'soft' | 'gradient';
		color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'neutral';
		size?: 'sm' | 'md' | 'lg';
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: (event: MouseEvent) => void;
		class?: string;
		children?: any;
	}

	let {
		variant = 'filled',
		color = 'primary',
		size = 'md',
		disabled = false,
		type = 'button',
		onclick,
		class: className = '',
		children
	}: Props = $props();

	const baseClasses = 'inline-flex items-center justify-center font-semibold rounded-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed active:scale-95 hover:scale-[1.02]';

	const sizeClasses = {
		sm: 'text-xs px-3 py-1.5 gap-1.5',
		md: 'text-sm px-4 py-2.5 gap-2',
		lg: 'text-base px-6 py-3 gap-2.5'
	};

	const variantColorClasses = {
		filled: {
			primary: 'bg-blue-500 hover:bg-blue-600 text-white shadow-sm shadow-blue-500/30 hover:shadow-md hover:shadow-blue-500/40',
			secondary: 'bg-purple-500 hover:bg-purple-600 text-white shadow-sm shadow-purple-500/30 hover:shadow-md hover:shadow-purple-500/40',
			success: 'bg-emerald-500 hover:bg-emerald-600 text-white shadow-sm shadow-emerald-500/30 hover:shadow-md hover:shadow-emerald-500/40',
			warning: 'bg-orange-500 hover:bg-orange-600 text-white shadow-sm shadow-orange-500/30 hover:shadow-md hover:shadow-orange-500/40',
			error: 'bg-red-500 hover:bg-red-600 text-white shadow-sm shadow-red-500/30 hover:shadow-md hover:shadow-red-500/40',
			neutral: 'bg-slate-700 hover:bg-slate-800 dark:bg-slate-600 dark:hover:bg-slate-700 text-white shadow-sm'
		},
		gradient: {
			primary: 'bg-gradient-to-r from-blue-500 to-cyan-500 hover:from-blue-600 hover:to-cyan-600 text-white shadow-lg shadow-blue-500/30 hover:shadow-xl hover:shadow-blue-500/40',
			secondary: 'bg-gradient-to-r from-purple-500 to-pink-500 hover:from-purple-600 hover:to-pink-600 text-white shadow-lg shadow-purple-500/30 hover:shadow-xl hover:shadow-purple-500/40',
			success: 'bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-600 hover:to-teal-600 text-white shadow-lg shadow-emerald-500/30 hover:shadow-xl hover:shadow-emerald-500/40',
			warning: 'bg-gradient-to-r from-orange-500 to-yellow-500 hover:from-orange-600 hover:to-yellow-600 text-white shadow-lg shadow-orange-500/30 hover:shadow-xl hover:shadow-orange-500/40',
			error: 'bg-gradient-to-r from-red-500 to-rose-500 hover:from-red-600 hover:to-rose-600 text-white shadow-lg shadow-red-500/30 hover:shadow-xl hover:shadow-red-500/40',
			neutral: 'bg-gradient-to-r from-slate-700 to-slate-600 hover:from-slate-800 hover:to-slate-700 text-white shadow-lg'
		},
		soft: {
			primary: 'bg-blue-100 hover:bg-blue-200 dark:bg-blue-900/30 dark:hover:bg-blue-900/50 text-blue-700 dark:text-blue-400',
			secondary: 'bg-purple-100 hover:bg-purple-200 dark:bg-purple-900/30 dark:hover:bg-purple-900/50 text-purple-700 dark:text-purple-400',
			success: 'bg-emerald-100 hover:bg-emerald-200 dark:bg-emerald-900/30 dark:hover:bg-emerald-900/50 text-emerald-700 dark:text-emerald-400',
			warning: 'bg-orange-100 hover:bg-orange-200 dark:bg-orange-900/30 dark:hover:bg-orange-900/50 text-orange-700 dark:text-orange-400',
			error: 'bg-red-100 hover:bg-red-200 dark:bg-red-900/30 dark:hover:bg-red-900/50 text-red-700 dark:text-red-400',
			neutral: 'bg-slate-100 hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-300'
		},
		ghost: {
			primary: 'hover:bg-blue-100 dark:hover:bg-blue-900/30 text-blue-600 dark:text-blue-400',
			secondary: 'hover:bg-purple-100 dark:hover:bg-purple-900/30 text-purple-600 dark:text-purple-400',
			success: 'hover:bg-emerald-100 dark:hover:bg-emerald-900/30 text-emerald-600 dark:text-emerald-400',
			warning: 'hover:bg-orange-100 dark:hover:bg-orange-900/30 text-orange-600 dark:text-orange-400',
			error: 'hover:bg-red-100 dark:hover:bg-red-900/30 text-red-600 dark:text-red-400',
			neutral: 'hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-700 dark:text-slate-300'
		}
	};

	const classes = `${baseClasses} ${sizeClasses[size]} ${variantColorClasses[variant][color]} ${className}`;
</script>

<button class={classes} {type} {disabled} {onclick}>
	{@render children?.()}
</button>
