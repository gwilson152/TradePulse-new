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
	import { apiClient } from '$lib/api/client';
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
	let importMode = $state<'csv' | 'api'>('csv');
	let tradingDate = $state(new Date().toISOString().split('T')[0]);
	let selectedFile: File | null = $state(null);
	let importResult: ImportResult | null = $state(null);
	let loading = $state(false);
	let error = $state('');
	let importing = $state(false);
	let importComplete = $state(false);

	// API import fields
	let apiSite = $state('');
	let apiUsername = $state('');
	let apiPassword = $state('');
	let apiFromDate = $state('');
	let apiToDate = $state('');

	// Duplicate detection
	let existingTrades = $state<Trade[]>([]);
	let duplicates = $state<{
		importedTrade: Partial<Trade>;
		existingTrade: Trade;
		action: 'skip' | 'replace' | 'keep-both';
	}[]>([]);

	const tabs = $derived([
		{ id: 'upload', label: 'Upload', icon: 'mdi:upload', disabled: false },
		{ id: 'preview', label: 'Preview', icon: 'mdi:eye', disabled: importResult === null },
		{ id: 'results', label: 'Results', icon: 'mdi:check-circle', disabled: !importComplete }
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
		loading = true;
		error = '';

		try {
			let result: ImportResult;

			if (importMode === 'api') {
				// API Import Mode
				if (!apiSite || !apiUsername || !apiPassword || !selectedPlatform) {
					error = 'Please fill in all API credentials';
					return;
				}

				if (selectedPlatform.id === 'prop-reports') {
					const trades = await apiClient.fetchPropReportsTrades(
						apiSite,
						apiUsername,
						apiPassword,
						apiFromDate || undefined,
						apiToDate || undefined
					);
					// Convert API response to ImportResult format
					result = {
						success: true,
						trades: trades,
						errors: [],
						warnings: [],
						statistics: {
							totalRows: trades.length,
							validTrades: trades.length,
							duplicates: 0,
							errors: 0,
							warnings: 0
						}
					};
					toast.success(`Fetched ${trades.length} trades from PropReports`);
				} else {
					error = 'API import not supported for this platform';
					return;
				}
			} else {
				// CSV Import Mode
				if (!selectedFile || !selectedPlatform) {
					error = 'Please select a file and platform';
					return;
				}

				const date = selectedPlatform.requiresDate ? new Date(tradingDate) : undefined;
				result = await importCSV(selectedFile, selectedPlatform, date);

				if (result.errors.length > 0) {
					toast.warning(`Parsed with ${result.errors.length} errors`);
				} else {
					toast.success(`Successfully parsed ${result.statistics.validTrades} trades`);
				}
			}

			importResult = result;

			// Fetch existing trades for duplicate detection
			await detectDuplicates(result.trades);

			activeTab = 'preview';
		} catch (err) {
			error = err instanceof Error ? err.message : importMode === 'api' ? 'Failed to fetch from API' : 'Failed to parse CSV';
			toast.error(importMode === 'api' ? 'Failed to fetch trades from API' : 'Failed to parse CSV file');
		} finally {
			loading = false;
		}
	}

	async function detectDuplicates(importedTrades: Partial<Trade>[]) {
		try {
			// Get date range from imported trades
			const dates = importedTrades
				.map(t => t.opened_at)
				.filter(d => d)
				.map(d => new Date(d!));

			if (dates.length === 0) return;

			const minDate = new Date(Math.min(...dates.map(d => d.getTime())));
			const maxDate = new Date(Math.max(...dates.map(d => d.getTime())));

			// Fetch existing trades in this date range
			existingTrades = await apiClient.getTradesForDateRange(
				minDate.toISOString().split('T')[0],
				maxDate.toISOString().split('T')[0]
			);

			// Find duplicates
			const found: typeof duplicates = [];
			importedTrades.forEach(importedTrade => {
				const match = existingTrades.find(existing =>
					existing.symbol === importedTrade.symbol &&
					existing.opened_at === importedTrade.opened_at &&
					Math.abs(Number(existing.entry_price) - Number(importedTrade.entry_price || 0)) < 0.01
				);

				if (match) {
					found.push({
						importedTrade,
						existingTrade: match,
						action: 'skip' // Default action
					});
				}
			});

			duplicates = found;

			if (found.length > 0) {
				toast.warning(`Found ${found.length} potential duplicates`);
			}
		} catch (err) {
			console.error('Failed to detect duplicates:', err);
			// Don't fail the import if duplicate detection fails
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
			// Process duplicates based on user's choices
			let tradesToImport = [...importResult.trades];
			const tradesToReplace: string[] = []; // IDs of trades to delete before import

			if (duplicates.length > 0) {
				duplicates.forEach(dup => {
					const tradeIndex = tradesToImport.findIndex(
						t => t.symbol === dup.importedTrade.symbol &&
						     t.opened_at === dup.importedTrade.opened_at
					);

					if (tradeIndex === -1) return;

					switch (dup.action) {
						case 'skip':
							// Remove from import list
							tradesToImport.splice(tradeIndex, 1);
							break;
						case 'replace':
							// Keep in import list, mark existing for deletion
							tradesToReplace.push(dup.existingTrade.id);
							break;
						case 'keep-both':
							// Keep both - do nothing
							break;
					}
				});

				// Delete trades marked for replacement
				for (const tradeId of tradesToReplace) {
					try {
						await apiClient.deleteTrade(tradeId);
					} catch (err) {
						console.error('Failed to delete trade:', err);
					}
				}
			}

			// Import remaining trades
			if (tradesToImport.length > 0) {
				await onImport(tradesToImport);
			}

			importComplete = true;
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
		importComplete = false;
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
	submitText={activeTab === 'preview' ? 'Import Trades' : importMode === 'api' ? 'Fetch from API' : 'Parse CSV'}
	submitColor={activeTab === 'preview' ? 'success' : 'primary'}
	submitDisabled={activeTab === 'upload' ? (importMode === 'csv' ? !selectedFile : !apiSite || !apiUsername || !apiPassword) : false}
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

			{#if selectedPlatform?.supportsAPI}
				<FormSection
					title="Import Method"
					icon="mdi:swap-horizontal"
					helpText="Choose between CSV file upload or direct API import"
				>
					<div class="flex gap-3">
						<button
							type="button"
							onclick={() => {
								importMode = 'csv';
								apiSite = '';
								apiUsername = '';
								apiPassword = '';
								apiFromDate = '';
								apiToDate = '';
							}}
							class="flex-1 p-4 rounded-xl border-2 transition-all {importMode === 'csv'
								? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
								: 'border-slate-200 dark:border-slate-700 hover:border-slate-300 dark:hover:border-slate-600'}"
						>
							<Icon icon="mdi:file-upload" width="32" class="mx-auto mb-2 {importMode === 'csv' ? 'text-blue-600' : 'text-slate-400'}" />
							<p class="font-semibold text-sm {importMode === 'csv' ? 'text-blue-900 dark:text-blue-100' : 'text-slate-700 dark:text-slate-300'}">
								CSV Upload
							</p>
							<p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Upload exported file</p>
						</button>
						<button
							type="button"
							onclick={() => {
								importMode = 'api';
								selectedFile = null;
							}}
							class="flex-1 p-4 rounded-xl border-2 transition-all {importMode === 'api'
								? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
								: 'border-slate-200 dark:border-slate-700 hover:border-slate-300 dark:hover:border-slate-600'}"
						>
							<Icon icon="mdi:api" width="32" class="mx-auto mb-2 {importMode === 'api' ? 'text-blue-600' : 'text-slate-400'}" />
							<p class="font-semibold text-sm {importMode === 'api' ? 'text-blue-900 dark:text-blue-100' : 'text-slate-700 dark:text-slate-300'}">
								API Import
							</p>
							<p class="text-xs text-slate-500 dark:text-slate-400 mt-1">Connect directly</p>
						</button>
					</div>
				</FormSection>
			{/if}

			{#if importMode === 'api' && selectedPlatform?.apiConfig}
				<FormSection
					title="API Credentials"
					icon="mdi:key"
					helpText="Enter your {selectedPlatform.name} credentials"
				>
					{#if selectedPlatform.apiConfig.requiresSite}
						<Input
							label={selectedPlatform.apiConfig.siteLabel || 'Site URL'}
							type="text"
							bind:value={apiSite}
							placeholder={selectedPlatform.apiConfig.sitePlaceholder}
							required={true}
						/>
					{/if}
					{#if selectedPlatform.apiConfig.requiresAuth}
						<Input
							label="Username"
							type="text"
							bind:value={apiUsername}
							required={true}
						/>
						<Input
							label="Password"
							type="password"
							bind:value={apiPassword}
							required={true}
						/>
					{/if}

					<!-- Date Range -->
					<div class="grid grid-cols-2 gap-4 mt-4">
						<Input
							label="From Date (optional)"
							type="date"
							bind:value={apiFromDate}
							placeholder="Start date"
						/>
						<Input
							label="To Date (optional)"
							type="date"
							bind:value={apiToDate}
							placeholder="End date"
						/>
					</div>

					<HelpText
						type="info"
						text="Leave dates blank to fetch all available trades. Specify a date range to limit results."
					/>
					<HelpText
						type="warning"
						text="Your credentials are only used to fetch trades and are not stored."
					/>
				</FormSection>
			{/if}

			{#if selectedPlatform?.requiresDate && importMode === 'csv'}
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

			{#if importMode === 'csv'}
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
			{/if}

			{#if selectedPlatform?.id === 'das-trader' && importMode === 'csv'}
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

				{#if duplicates.length > 0}
					<FormSection
						title="Duplicate Trades Found"
						icon="mdi:file-multiple"
					>
						<HelpText
							type="warning"
							text="The following trades appear to already exist in your account. Choose how to handle each duplicate."
						/>
						<div class="space-y-3 mt-4 max-h-96 overflow-y-auto">
							{#each duplicates as dup, index}
								<div class="p-4 bg-orange-50 dark:bg-orange-900/20 rounded-lg border border-orange-200 dark:border-orange-800">
									<div class="flex items-start justify-between gap-4">
										<div class="flex-1">
											<div class="flex items-center gap-2 mb-2">
												<Badge color="warning" variant="soft">
													{dup.importedTrade.symbol}
												</Badge>
												<span class="text-xs text-slate-600 dark:text-slate-400">
													{new Date(dup.importedTrade.opened_at || '').toLocaleString()}
												</span>
											</div>
											<div class="grid grid-cols-2 gap-2 text-sm">
												<div>
													<p class="text-xs text-slate-500 dark:text-slate-400">Existing:</p>
													<p class="font-medium">${dup.existingTrade.entry_price} × {dup.existingTrade.quantity}</p>
													<p class="text-xs">P&L: <span class="{Number(dup.existingTrade.pnl) >= 0 ? 'text-emerald-600' : 'text-red-600'}">${dup.existingTrade.pnl?.toFixed(2)}</span></p>
												</div>
												<div>
													<p class="text-xs text-slate-500 dark:text-slate-400">Importing:</p>
													<p class="font-medium">${dup.importedTrade.entry_price} × {dup.importedTrade.quantity}</p>
													<p class="text-xs">P&L: <span class="{Number(dup.importedTrade.pnl) >= 0 ? 'text-emerald-600' : 'text-red-600'}">${dup.importedTrade.pnl?.toFixed(2)}</span></p>
												</div>
											</div>
										</div>
										<div class="flex flex-col gap-1">
											<button
												type="button"
												onclick={() => duplicates[index].action = 'skip'}
												class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors {dup.action === 'skip'
													? 'bg-slate-600 text-white'
													: 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'}"
											>
												Skip
											</button>
											<button
												type="button"
												onclick={() => duplicates[index].action = 'replace'}
												class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors {dup.action === 'replace'
													? 'bg-blue-600 text-white'
													: 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'}"
											>
												Replace
											</button>
											<button
												type="button"
												onclick={() => duplicates[index].action = 'keep-both'}
												class="px-3 py-1.5 rounded-lg text-xs font-medium transition-colors {dup.action === 'keep-both'
													? 'bg-emerald-600 text-white'
													: 'bg-slate-100 dark:bg-slate-800 text-slate-700 dark:text-slate-300 hover:bg-slate-200 dark:hover:bg-slate-700'}"
											>
												Keep Both
											</button>
										</div>
									</div>
								</div>
							{/each}
						</div>
						<div class="mt-4 flex gap-2">
							<Button
								type="button"
								variant="soft"
								color="secondary"
								size="sm"
								onclick={() => duplicates.forEach((_, i) => duplicates[i].action = 'skip')}
							>
								Skip All
							</Button>
							<Button
								type="button"
								variant="soft"
								color="primary"
								size="sm"
								onclick={() => duplicates.forEach((_, i) => duplicates[i].action = 'replace')}
							>
								Replace All
							</Button>
						</div>
					</FormSection>
				{/if}

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
