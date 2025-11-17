<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		open: boolean;
		title?: string;
		description?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl';
		showCloseButton?: boolean;
		onClose: () => void;
		children?: any;
		titleSlot?: any;
		actionsSlot?: any;
	}

	let {
		open = false,
		title,
		description,
		size = 'md',
		showCloseButton = true,
		onClose,
		children,
		titleSlot,
		actionsSlot
	}: Props = $props();

	const sizeClasses = {
		sm: 'max-w-sm',
		md: 'max-w-md',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl'
	};

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) {
			onClose();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			onClose();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<div
		class="fixed inset-0 z-[100] overflow-y-auto"
		role="dialog"
		aria-modal="true"
		onclick={handleBackdropClick}
	>
		<div class="modal-backdrop"></div>

		<div class="relative min-h-screen flex items-center justify-center p-4">
			<div
				class="modal {sizeClasses[size]} w-full max-h-[90vh] flex flex-col"
				role="dialog"
				aria-modal="true"
				onclick={(e) => e.stopPropagation()}
			>
				<!-- Header -->
				{#if title || titleSlot}
					<header class="modal-header flex items-center justify-between px-6 py-4 border-b border-surface-200 dark:border-surface-700 flex-shrink-0">
						<div class="flex-1 min-w-0">
							{#if titleSlot}
								{@render titleSlot()}
							{:else}
								<h3 class="text-xl font-bold text-surface-900 dark:text-surface-100 truncate">
									{title}
								</h3>
								{#if description}
									<p class="text-sm text-surface-600 dark:text-surface-400 mt-1">
										{description}
									</p>
								{/if}
							{/if}
						</div>
						{#if showCloseButton}
							<button
								onclick={onClose}
								class="ml-4 p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded-lg transition-colors flex-shrink-0"
								aria-label="Close modal"
							>
								<Icon icon="mdi:close" width="20" class="text-surface-600 dark:text-surface-400" />
							</button>
						{/if}
					</header>
				{/if}

				<!-- Body -->
				<div class="modal-body flex-1 overflow-y-auto px-6 py-5">
					{@render children?.()}
				</div>

				<!-- Footer/Actions -->
				{#if actionsSlot}
					<footer class="modal-footer px-6 py-4 border-t border-surface-200 dark:border-surface-700 flex-shrink-0">
						{@render actionsSlot()}
					</footer>
				{/if}
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
