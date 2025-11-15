<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		open: boolean;
		title?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		onClose: () => void;
		children?: any;
	}

	let { open = false, title, size = 'md', onClose, children }: Props = $props();

	const sizeClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl'
	};
</script>

{#if open}
	<div class="fixed inset-0 z-[100] overflow-y-auto" role="dialog" aria-modal="true">
		<div
			class="modal-backdrop"
			onclick={onClose}
			onkeydown={(e) => e.key === 'Escape' && onClose()}
			role="button"
			tabindex="-1"
			aria-label="Close modal"
		></div>

		<div class="relative min-h-screen flex items-center justify-center p-4">
			<div
				class="modal card {sizeClasses[size]} w-full max-h-[90vh] overflow-y-auto"
				onclick={(e) => e.stopPropagation()}
				onkeydown={() => {}}
				role="document"
			>
				{#if title}
					<header class="modal-header flex items-center justify-between p-6 border-b border-surface-200 dark:border-surface-700">
						<h3 class="text-2xl font-bold">{title}</h3>
						<button onclick={onClose} class="p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded">
							<Icon icon="mdi:close" width="24" />
						</button>
					</header>
				{/if}

				<div class="modal-body">
					{@render children?.()}
				</div>
			</div>
		</div>
	</div>
{/if}

<style>
	.modal-backdrop {
		@apply fixed inset-0 bg-black/50 z-[90];
	}

	.modal {
		@apply relative bg-white dark:bg-surface-800 rounded-lg shadow-xl z-[100];
	}
</style>
