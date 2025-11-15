<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		color?: 'primary' | 'secondary' | 'success' | 'warning' | 'error' | 'info' | 'neutral';
		variant?: 'filled' | 'outlined' | 'soft';
		size?: 'sm' | 'md' | 'lg';
		icon?: string;
		removable?: boolean;
		onRemove?: () => void;
		children: any;
	}

	let {
		color = 'neutral',
		variant = 'soft',
		size = 'md',
		icon,
		removable = false,
		onRemove,
		children
	}: Props = $props();

	function getColorClasses(): string {
		const variants = {
			filled: {
				primary: 'bg-primary-500 text-white',
				secondary: 'bg-secondary-500 text-white',
				success: 'bg-success-500 text-white',
				warning: 'bg-warning-500 text-white',
				error: 'bg-error-500 text-white',
				info: 'bg-blue-500 text-white',
				neutral: 'bg-surface-500 text-white'
			},
			outlined: {
				primary:
					'border-2 border-primary-500 text-primary-700 dark:text-primary-300 bg-transparent',
				secondary:
					'border-2 border-secondary-500 text-secondary-700 dark:text-secondary-300 bg-transparent',
				success:
					'border-2 border-success-500 text-success-700 dark:text-success-300 bg-transparent',
				warning:
					'border-2 border-warning-500 text-warning-700 dark:text-warning-300 bg-transparent',
				error: 'border-2 border-error-500 text-error-700 dark:text-error-300 bg-transparent',
				info: 'border-2 border-blue-500 text-blue-700 dark:text-blue-300 bg-transparent',
				neutral:
					'border-2 border-surface-500 text-surface-700 dark:text-surface-300 bg-transparent'
			},
			soft: {
				primary: 'bg-primary-100 text-primary-700 dark:bg-primary-900/30 dark:text-primary-300',
				secondary:
					'bg-secondary-100 text-secondary-700 dark:bg-secondary-900/30 dark:text-secondary-300',
				success: 'bg-success-100 text-success-700 dark:bg-success-900/30 dark:text-success-300',
				warning: 'bg-warning-100 text-warning-700 dark:bg-warning-900/30 dark:text-warning-300',
				error: 'bg-error-100 text-error-700 dark:bg-error-900/30 dark:text-error-300',
				info: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
				neutral: 'bg-surface-100 text-surface-700 dark:bg-surface-800 dark:text-surface-300'
			}
		};

		return variants[variant][color];
	}

	function getSizeClasses(): string {
		const sizes = {
			sm: 'text-xs px-2 py-0.5',
			md: 'text-sm px-2.5 py-1',
			lg: 'text-base px-3 py-1.5'
		};

		return sizes[size];
	}
</script>

<span
	class="badge inline-flex items-center gap-1 font-medium rounded-full {getColorClasses()} {getSizeClasses()}"
>
	{#if icon}
		<Icon {icon} class="text-current" />
	{/if}

	{@render children()}

	{#if removable && onRemove}
		<button
			type="button"
			onclick={onRemove}
			class="ml-1 -mr-1 hover:opacity-70 transition-opacity"
			aria-label="Remove"
		>
			<Icon icon="mdi:close" class="text-current" />
		</button>
	{/if}
</span>

<style>
	.badge {
		transition: all 0.2s ease;
	}
</style>
