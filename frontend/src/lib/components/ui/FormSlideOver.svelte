<script lang="ts">
	import Icon from '@iconify/svelte';
	import Button from './Button.svelte';
	import HelpText from './HelpText.svelte';

	interface Tab {
		id: string;
		label: string;
		icon: string;
		isComplete?: boolean;
	}

	interface Props {
		open: boolean;
		title: string;
		subtitle?: string;
		size?: 'md' | 'lg' | 'xl' | '2xl';
		loading?: boolean;
		error?: string;
		showHelp?: boolean;
		helpTitle?: string;
		helpText?: string;
		helpType?: 'info' | 'tip' | 'warning';
		tabs?: Tab[];
		activeTab?: string;
		onTabChange?: (tabId: string) => void;
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
		size = 'xl',
		loading = false,
		error,
		showHelp = false,
		helpTitle,
		helpText,
		helpType = 'info',
		tabs,
		activeTab = $bindable(),
		onTabChange,
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
		md: 'max-w-md',
		lg: 'max-w-lg',
		xl: 'max-w-xl',
		'2xl': 'max-w-2xl'
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

	function handleBackdropClick(e: MouseEvent) {
		if (!loading && e.target === e.currentTarget) {
			onClose();
		}
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' && !loading) {
			onClose();
		}
	}

	function handleTabClick(tabId: string) {
		if (activeTab !== undefined) {
			activeTab = tabId;
		}
		if (onTabChange) {
			onTabChange(tabId);
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if open}
	<!-- Backdrop -->
	<button
		class="fixed inset-0 z-[100] bg-black/20 backdrop-blur-sm transition-opacity duration-300 cursor-default"
		onclick={handleBackdropClick}
		onkeydown={(e) => e.key === 'Enter' && handleBackdropClick()}
		tabindex="-1"
		aria-label="Close panel"
	></button>

	<!-- Slide-over Panel -->
	<div
		class="fixed inset-y-0 right-0 z-[101] {sizeClasses[size]} w-full
			   bg-white/95 dark:bg-slate-900/95 backdrop-blur-xl
			   shadow-2xl border-l border-slate-200/50 dark:border-slate-700/50
			   transform transition-all duration-300 ease-out
			   {open ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'}"
		role="dialog"
		aria-modal="true"
		aria-labelledby="slideover-title"
	>
		<form onsubmit={handleSubmit} class="h-full flex flex-col">
			<!-- Header -->
			<div
				class="sticky top-0 z-10 flex flex-col
					   bg-gradient-to-b from-white/90 to-white/50 dark:from-slate-900/90 dark:to-slate-900/50
					   backdrop-blur-xl border-b border-slate-200/50 dark:border-slate-700/50"
			>
				<!-- Title and Close Button -->
				<div class="flex items-start justify-between px-6 py-5">
					<div class="flex-1 pr-4">
						<h2
							id="slideover-title"
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
						type="button"
						onclick={handleClose}
						disabled={loading}
						class="p-2 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-xl transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
						aria-label="Close panel"
					>
						<Icon icon="mdi:close" width="24" class="text-slate-600 dark:text-slate-400" />
					</button>
				</div>

				<!-- Tabs (if provided) -->
				{#if tabs && tabs.length > 0}
					<div class="px-6">
						<div class="flex gap-1 -mb-px overflow-x-auto">
							{#each tabs as tab}
								<button
									type="button"
									onclick={() => handleTabClick(tab.id)}
									class="group relative flex items-center gap-2.5 px-4 py-3 border-b-2 transition-all duration-200 whitespace-nowrap
										{activeTab === tab.id
											? 'border-blue-500 text-blue-600 dark:text-blue-400 bg-blue-50/50 dark:bg-blue-900/20'
											: 'border-transparent text-slate-600 dark:text-slate-400 hover:text-slate-800 dark:hover:text-slate-200 hover:bg-slate-50 dark:hover:bg-slate-800/30'}
										rounded-t-lg font-medium text-sm"
								>
									<Icon icon={tab.icon} width="18" />
									<span>{tab.label}</span>
									{#if tab.isComplete}
										<Icon icon="mdi:check-circle" width="16" class="text-emerald-500" />
									{/if}
								</button>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			<!-- Scrollable Content -->
			<div class="flex-1 overflow-y-auto px-6 py-6 space-y-6">
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
				class="sticky bottom-0 flex items-center justify-end gap-3 px-6 py-5
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
{/if}
