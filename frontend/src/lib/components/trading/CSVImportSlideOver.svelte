<script lang="ts">
	import FormSlideOver from '../ui/FormSlideOver.svelte';
	import FormSection from '../ui/FormSection.svelte';
	import FileDropzone from '../ui/FileDropzone.svelte';
	import Select from '../ui/Select.svelte';
	import Input from '../ui/Input.svelte';
	import Button from '../ui/Button.svelte';
	import HelpText from '../ui/HelpText.svelte';
	import Badge from '../ui/Badge.svelte';
	import Icon from '@iconify/svelte';
	import { toast } from '$lib/stores/toast';
	import { platforms } from '$lib/utils/platforms';
	import { importCSV } from '$lib/utils/csvImport';
	import type { ImportResult } from '$lib/types/import';
	import type { Trade } from '$lib/types';

	interface Props {
		open: boolean;
		onClose: () => void;
		onImport: (trades: Partial<Trade>[]) => Promise<void>;
	}

	let { open, onClose, onImport }: Props = $props();

	let activeTab = $state<'upload' | 'preview' | 'results'>('upload');
	let selectedPlatformId = $state('das-trader');
	let tradingDate = $state(new Date().toISOString().split('T')[0]);
	let selectedFile: File | null = $state(null);
	let importResult: ImportResult | null = $state(null);
	let loading = $state(false);
	let error = $state('');
	let importing = $state(false);

	const tabs = $derived([
		{ id: 'upload', label: 'Upload', icon: 'mdi:upload' },
		{ id: 'preview', label: 'Preview', icon: 'mdi:eye', isComplete: importResult !== null },
		{ id: 'results', label: 'Results', icon: 'mdi:check-circle', isComplete: false }
	]);

	const platformOptions = $derived(
		platforms.map((p) => ({
			value: p.id,
			label: p.name
		}))
	);

	const selectedPlatform = $derived(platforms.find((p) => p.id === selectedPlatformId));

	function handleFileSelect(file: File) {
		selectedFile = file;
		error = '';
		importResult = null;
	}

	async function handleParseCSV() {
		if (!selectedFile || !selectedPlatform) {
			error = 'Please select a file and platform';
			return;
		}

		loading = true;
		error = '';

		try {
			const date = selectedPlatform.requiresDate ? new Date(tradingDate) : undefined;
			const result = await importCSV(selectedFile, selectedPlatform, date);
			importResult = result;

			if (result.errors.length > 0) {
				toast.warning(`Parsed with ${result.errors.length} errors`);
			} else {
				toast.success(`Successfully parsed ${result.statistics.validTrades} trades`);
			}

			activeTab = 'preview';
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to parse CSV';
			toast.error('Failed to parse CSV file');
		} finally {
			loading = false;
		}
	}

	async function handleConfirmImport() {
		if (!importResult || importResult.trades.length === 0) {
			error = 'No trades to import';
			return;
		}

		importing = true;
		error = '';

		try {
			await onImport(importResult.trades);
			toast.success(`Successfully imported ${importResult.trades.length} trades`);
			activeTab = 'results';
			setTimeout(() => {
				handleClose();
			}, 2000);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to import trades';
			toast.error('Failed to import trades');
		} finally {
			importing = false;
		}
	}

	function handleClose() {
		resetForm();
		onClose();
	}

	function resetForm() {
		activeTab = 'upload';
		selectedFile = null;
		importResult = null;
		error = '';
		loading = false;
		importing = false;
	}
</script>

<FormSlideOver
	{open}
	title="Import Trade History"
	subtitle="Import trades from CSV export"
	size="2xl"
	{error}
	loading={loading || importing}
	showHelp={true}
	helpTitle="CSV Import"
	helpText="Upload your trade history CSV export from your broker. We'll parse and validate the data before importing."
	{tabs}
	bind:activeTab
	onTabChange={(tabId) => (activeTab = tabId as typeof activeTab)}
	onClose={handleClose}
	onSubmit={activeTab === 'preview' ? handleConfirmImport : handleParseCSV}
	submitText={activeTab === 'preview' ? 'Import Trades' : 'Parse CSV'}
	submitColor={activeTab === 'preview' ? 'success' : 'primary'}
	submitDisabled={activeTab === 'upload' ? !selectedFile : false}
>
	{#if activeTab === 'upload'}
		<div class="space-y-6">
			<FormSection
				title="Select Platform"
				icon="mdi:application"
				helpText="Choose your trading platform to use the correct CSV format"
			>
				<Select
					label="Trading Platform"
					options={platformOptions}
					bind:value={selectedPlatformId}
					required={true}
				/>
				{#if selectedPlatform}
					<HelpText
						type="info"
						text={selectedPlatform.description}
					/>
				{/if}
			</FormSection>

			{#if selectedPlatform?.requiresDate}
				<FormSection
					title="Specify Trading Date"
					icon="mdi:calendar"
					helpText="DAS Trader exports don't include the date, only time"
				>
					<Input
						label="Trading Date"
						type="date"
						bind:value={tradingDate}
						required={true}
					/>
					<HelpText
						type="warning"
						text="Make sure to select the correct trading day for your CSV export. DAS Trader only exports the current day's trades."
					/>
				</FormSection>
			{/if}

			<FormSection
				title="Upload CSV File"
				icon="mdi:file-upload"
				helpText="Drag and drop your CSV file or click to browse"
			>
				<FileDropzone
					accept=".csv"
					maxSizeMB={10}
					onFileSelect={handleFileSelect}
				/>
				{#if selectedFile}
					<div class="mt-4 p-4 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl border border-emerald-200 dark:border-emerald-800">
						<div class="flex items-center gap-3">
							<Icon icon="mdi:file-check" width="24" class="text-emerald-600" />
							<div class="flex-1">
								<p class="font-semibold text-emerald-900 dark:text-emerald-100">
									{selectedFile.name}
								</p>
								<p class="text-sm text-emerald-700 dark:text-emerald-300">
									{(selectedFile.size / 1024).toFixed(2)} KB
								</p>
							</div>
							<Button
								type="button"
								variant="soft"
								color="error"
								size="sm"
								onclick={() => {
									selectedFile = null;
									importResult = null;
								}}
							>
								<Icon icon="mdi:close" width="16" />
								Remove
							</Button>
						</div>
					</div>
				{/if}
			</FormSection>

			{#if selectedPlatform?.id === 'das-trader'}
				<HelpText
					type="tip"
					title="How to Export from DAS Trader"
					text="1. Click 'Trade' menu → 'Trade Log'\n2. Right-click in the Trade Log window\n3. Select 'Export' to download CSV\n4. Make sure to export before 10 PM ET each day"
					collapsible={true}
				/>
			{/if}
		</div>
	{:else if activeTab === 'preview'}
		<div class="space-y-6">
			{#if importResult}
				<FormSection
					title="Import Statistics"
					icon="mdi:chart-bar"
				>
					<div class="grid grid-cols-2 md:grid-cols-4 gap-4">
						<div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200 dark:border-blue-800">
							<p class="text-sm text-blue-700 dark:text-blue-300 mb-1">Total Rows</p>
							<p class="text-3xl font-bold text-blue-900 dark:text-blue-100">
								{importResult.statistics.totalRows}
							</p>
						</div>
						<div class="p-4 bg-emerald-50 dark:bg-emerald-900/20 rounded-xl border border-emerald-200 dark:border-emerald-800">
							<p class="text-sm text-emerald-700 dark:text-emerald-300 mb-1">Valid Trades</p>
							<p class="text-3xl font-bold text-emerald-900 dark:text-emerald-100">
								{importResult.statistics.validTrades}
							</p>
						</div>
						<div class="p-4 bg-amber-50 dark:bg-amber-900/20 rounded-xl border border-amber-200 dark:border-amber-800">
							<p class="text-sm text-amber-700 dark:text-amber-300 mb-1">Warnings</p>
							<p class="text-3xl font-bold text-amber-900 dark:text-amber-100">
								{importResult.statistics.warnings}
							</p>
						</div>
						<div class="p-4 bg-red-50 dark:bg-red-900/20 rounded-xl border border-red-200 dark:border-red-800">
							<p class="text-sm text-red-700 dark:text-red-300 mb-1">Errors</p>
							<p class="text-3xl font-bold text-red-900 dark:text-red-100">
								{importResult.statistics.errors}
							</p>
						</div>
					</div>
				</FormSection>

				{#if importResult.errors.length > 0}
					<FormSection
						title="Errors"
						icon="mdi:alert-circle"
					>
						<div class="space-y-2 max-h-64 overflow-y-auto">
							{#each importResult.errors as err}
								<div class="p-3 bg-red-50 dark:bg-red-900/20 rounded-lg border border-red-200 dark:border-red-800">
									<div class="flex items-start gap-2">
										<Icon icon="mdi:alert" width="20" class="text-red-600 flex-shrink-0 mt-0.5" />
										<div class="flex-1">
											<p class="text-sm font-semibold text-red-900 dark:text-red-100">
												Row {err.row}: {err.message}
											</p>
											{#if err.column}
												<p class="text-xs text-red-700 dark:text-red-300 mt-1">
													Column: {err.column}
												</p>
											{/if}
										</div>
									</div>
								</div>
							{/each}
						</div>
					</FormSection>
				{/if}

				{#if importResult.warnings.length > 0}
					<FormSection
						title="Warnings"
						icon="mdi:alert-outline"
					>
						<div class="space-y-2 max-h-64 overflow-y-auto">
							{#each importResult.warnings as warn}
								<div class="p-3 bg-amber-50 dark:bg-amber-900/20 rounded-lg border border-amber-200 dark:border-amber-800">
									<div class="flex items-start gap-2">
										<Icon icon="mdi:information" width="20" class="text-amber-600 flex-shrink-0 mt-0.5" />
										<div class="flex-1">
											<p class="text-sm font-medium text-amber-900 dark:text-amber-100">
												{warn.message}
											</p>
											<div class="flex gap-2 mt-1">
												<Badge color="warning" variant="soft" size="sm">
													{warn.type}
												</Badge>
												<span class="text-xs text-amber-700 dark:text-amber-300">
													Row {warn.row}
												</span>
											</div>
										</div>
									</div>
								</div>
							{/each}
						</div>
					</FormSection>
				{/if}

				<FormSection
					title="Trade Preview"
					icon="mdi:table"
					helpText="First 10 trades that will be imported"
				>
					<div class="overflow-x-auto">
						<table class="w-full text-sm">
							<thead class="bg-slate-100 dark:bg-slate-800">
								<tr>
									<th class="px-3 py-2 text-left font-semibold">Symbol</th>
									<th class="px-3 py-2 text-left font-semibold">Type</th>
									<th class="px-3 py-2 text-right font-semibold">Entry</th>
									<th class="px-3 py-2 text-right font-semibold">Qty</th>
									<th class="px-3 py-2 text-right font-semibold">Exit</th>
									<th class="px-3 py-2 text-right font-semibold">P&L</th>
								</tr>
							</thead>
							<tbody>
								{#each importResult.trades.slice(0, 10) as trade}
									<tr class="border-b border-slate-200 dark:border-slate-700">
										<td class="px-3 py-2 font-medium">{trade.symbol}</td>
										<td class="px-3 py-2">
											<Badge
												color={trade.trade_type === 'LONG' ? 'success' : 'error'}
												variant="soft"
												size="sm"
											>
												{trade.trade_type}
											</Badge>
										</td>
										<td class="px-3 py-2 text-right">${trade.entry_price?.toFixed(2)}</td>
										<td class="px-3 py-2 text-right">{trade.quantity}</td>
										<td class="px-3 py-2 text-right">
											{trade.exit_price ? `$${trade.exit_price.toFixed(2)}` : '-'}
										</td>
										<td class="px-3 py-2 text-right font-medium {trade.pnl && trade.pnl > 0 ? 'text-emerald-600' : trade.pnl && trade.pnl < 0 ? 'text-red-600' : 'text-slate-600'}">
											{trade.pnl !== null && trade.pnl !== undefined ? `${trade.pnl >= 0 ? '+' : ''}$${trade.pnl.toFixed(2)}` : '-'}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
						{#if importResult.trades.length > 10}
							<p class="text-sm text-slate-500 dark:text-slate-400 mt-2 text-center">
								... and {importResult.trades.length - 10} more trades
							</p>
						{/if}
					</div>
				</FormSection>

				{#if importResult.errors.length === 0}
					<HelpText
						type="success"
						text="All trades parsed successfully! Click 'Import Trades' to add them to your account."
					/>
				{:else}
					<HelpText
						type="warning"
						text={`${importResult.errors.length} rows had errors and will be skipped. ${importResult.statistics.validTrades} valid trades will be imported.`}
					/>
				{/if}
			{/if}
		</div>
	{:else if activeTab === 'results'}
		<div class="space-y-6 text-center py-12">
			<div class="mb-6">
				<Icon icon="mdi:check-circle" width="96" class="mx-auto text-emerald-500" />
			</div>
			<h3 class="text-2xl font-bold text-slate-800 dark:text-slate-200">
				Import Complete!
			</h3>
			<p class="text-lg text-slate-600 dark:text-slate-400">
				Successfully imported {importResult?.statistics.validTrades} trades
			</p>
			<HelpText
				type="success"
				text="Your trades have been added to your account. You can now view, analyze, and journal them."
			/>
		</div>
	{/if}
</FormSlideOver>
