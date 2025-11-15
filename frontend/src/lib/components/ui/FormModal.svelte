<script lang="ts">
	import Icon from '@iconify/svelte';
	import Button from './Button.svelte';
	import HelpText from './HelpText.svelte';

	interface Props {
		open: boolean;
		title: string;
		subtitle?: string;
		size?: 'sm' | 'md' | 'lg' | 'xl' | 'full';
		loading?: boolean;
		error?: string;
		showHelp?: boolean;
		helpTitle?: string;
		helpText?: string;
		helpType?: 'info' | 'tip' | 'warning';
		onClose: () => void;
		onSubmit?: (e: Event) => void;
		submitText?: string;
		submitDisabled?: boolean;
		submitColor?: 'primary' | 'success' | 'error';
		showCancel?: boolean;
		cancelText?: string;
		children?: any;
	}

	let {
		open = false,
		title,
		subtitle,
		size = 'lg',
		loading = false,
		error,
		showHelp = false,
		helpTitle,
		helpText,
		helpType = 'info',
		onClose,
		onSubmit,
		submitText = 'Save',
		submitDisabled = false,
		submitColor = 'primary',
		showCancel = true,
		cancelText = 'Cancel',
		children
	}: Props = $props();

	const sizeClasses = {
		sm: 'max-w-md',
		md: 'max-w-lg',
		lg: 'max-w-2xl',
		xl: 'max-w-4xl',
		full: 'max-w-7xl'
	};

	function handleSubmit(e: Event) {
		e.preventDefault();
		if (onSubmit) {
			onSubmit(e);
		}
	}

	function handleClose() {
		if (!loading) {
			onClose();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && !loading) {
			onClose();
		}
	}
</script>

{#if open}
	<div
		class="fixed inset-0 z-[100] overflow-y-auto"
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
	>
		<!-- Backdrop -->
		<div
			class="fixed inset-0 bg-black/60 backdrop-blur-sm transition-opacity"
			onclick={handleClose}
			onkeydown={handleKeydown}
			role="button"
			tabindex="-1"
			aria-label="Close modal"
		></div>

		<!-- Modal Container -->
		<div class="relative min-h-screen flex items-center justify-center p-4">
			<!-- Modal Content -->
			<div
				class="relative {sizeClasses[size]} w-full max-h-[90vh] overflow-hidden
					bg-white/95 dark:bg-slate-900/95 backdrop-blur-xl
					rounded-2xl shadow-2xl border border-slate-200/50 dark:border-slate-700/50"
				onclick={(e) => e.stopPropagation()}
				onkeydown={() => {}}
				role="document"
			>
				<!-- Header -->
				<div
					class="sticky top-0 z-10 flex items-start justify-between px-8 py-6
						bg-gradient-to-b from-white/90 to-white/50 dark:from-slate-900/90 dark:to-slate-900/50
						backdrop-blur-xl border-b border-slate-200/50 dark:border-slate-700/50"
				>
					<div class="flex-1 pr-4">
						<h2
							id="modal-title"
							class="text-2xl font-bold bg-gradient-to-r from-slate-800 to-slate-600 dark:from-slate-100 dark:to-slate-400 bg-clip-text text-transparent"
						>
							{title}
						</h2>
						{#if subtitle}
							<p class="text-sm text-slate-600 dark:text-slate-400 mt-1">
								{subtitle}
							</p>
						{/if}
					</div>
					<button
						onclick={handleClose}
						disabled={loading}
						class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
						aria-label="Close modal"
					>
						<Icon icon="mdi:close" width="24" class="text-slate-600 dark:text-slate-400" />
					</button>
				</div>

				<!-- Form -->
				<form onsubmit={handleSubmit} class="flex flex-col" style="max-height: calc(90vh - 140px);">
					<!-- Body -->
					<div class="flex-1 overflow-y-auto px-8 py-6 space-y-6">
						<!-- Error Alert -->
						{#if error}
							<div
								class="bg-red-50 dark:bg-red-900/20 border-2 border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 px-4 py-3 rounded-xl flex items-start gap-3"
							>
								<Icon icon="mdi:alert-circle" width="20" class="mt-0.5 flex-shrink-0" />
								<p class="text-sm font-medium">{error}</p>
							</div>
						{/if}

						<!-- Help Section -->
						{#if showHelp && helpText}
							<HelpText type={helpType} title={helpTitle} text={helpText} collapsible={true} />
						{/if}

						<!-- Content Slot -->
						{@render children?.()}
					</div>

					<!-- Footer -->
					<div
						class="sticky bottom-0 flex items-center justify-end gap-3 px-8 py-5
						bg-gradient-to-t from-white/90 to-white/50 dark:from-slate-900/90 dark:to-slate-900/50
						backdrop-blur-xl border-t border-slate-200/50 dark:border-slate-700/50"
					>
						{#if showCancel}
							<Button type="button" variant="ghost" color="neutral" onclick={handleClose} disabled={loading}>
								{cancelText}
							</Button>
						{/if}
						{#if onSubmit}
							<Button type="submit" variant="gradient" color={submitColor} disabled={submitDisabled || loading}>
								{#if loading}
					<Icon icon="mdi:loading" width="20" class="animate-spin" />
									Saving...
								{:else}
									<Icon icon="mdi:content-save" width="20" />
									{submitText}
								{/if}
							</Button>
						{/if}
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}
