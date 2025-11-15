<script lang="ts">
	import Icon from '@iconify/svelte';
	import Input from '../ui/Input.svelte';
	import Select from '../ui/Select.svelte';
	import Textarea from '../ui/Textarea.svelte';
	import Button from '../ui/Button.svelte';
	import Badge from '../ui/Badge.svelte';
	import HelpText from '../ui/HelpText.svelte';
	import Tooltip from '../ui/Tooltip.svelte';
	import type { Trade, CostBasisMethod } from '$lib/types';

	interface Props {
		isOpen: boolean;
		onClose: () => void;
		onSave: (trade: Partial<Trade>) => void;
		trade?: Trade;
	}

	let { isOpen, onClose, onSave, trade }: Props = $props();

	let formData = $state({
		symbol: trade?.symbol || '',
		tradeType: trade?.trade_type || 'LONG',
		quantity: trade?.quantity || '',
		entryPrice: trade?.entry_price || '',
		entryDate: trade?.opened_at ? new Date(trade.opened_at).toISOString().split('T')[0] : new Date().toISOString().split('T')[0],
		entryTime: trade?.opened_at ? new Date(trade.opened_at).toTimeString().split(' ')[0].slice(0, 5) : '09:30',
		exitPrice: trade?.exit_price || '',
		exitDate: trade?.closed_at ? new Date(trade.closed_at).toISOString().split('T')[0] : '',
		exitTime: trade?.closed_at ? new Date(trade.closed_at).toTimeString().split(' ')[0].slice(0, 5) : '',
		fees: trade?.fees?.toString() || '0',
		notes: trade?.notes || '',
		tags: trade?.tags || [],
		costBasisMethod: (trade?.cost_basis_method || 'FIFO') as CostBasisMethod
	});

	let newTag = $state('');
	let error = $state('');

	const costBasisOptions = [
		{ value: 'FIFO', label: 'First In First Out (FIFO)' },
		{ value: 'LIFO', label: 'Last In First Out (LIFO)' },
		{ value: 'AVERAGE', label: 'Average Cost' }
	];

	function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		// Validation
		if (!formData.symbol || !formData.quantity || !formData.entryPrice) {
			error = 'Please fill in all required fields';
			return;
		}

		// Build trade data with position lifecycle structure
		const entryDateTime = new Date(`${formData.entryDate}T${formData.entryTime}`);
		const quantity = parseFloat(formData.quantity);
		const entryPrice = parseFloat(formData.entryPrice);
		const fees = parseFloat(formData.fees);

		const tradeData: Partial<Trade> = {
			symbol: formData.symbol.toUpperCase(),
			trade_type: formData.tradeType as 'LONG' | 'SHORT',
			// Legacy fields for backward compatibility
			quantity,
			entry_price: entryPrice,
			fees,
			// Position lifecycle fields
			entries: [
				{
					id: crypto.randomUUID(),
					price: entryPrice,
					quantity,
					timestamp: entryDateTime.toISOString(),
					fees: fees / 2 // Split fees between entry and potential exit
				}
			],
			exits: [],
			current_position_size: quantity,
			average_entry_price: entryPrice,
			total_fees: fees,
			realized_pnl: 0,
			unrealized_pnl: null,
			cost_basis_method: formData.costBasisMethod,
			opened_at: entryDateTime.toISOString(),
			notes: formData.notes,
			tags: formData.tags
		};

		// Add exit if provided
		if (formData.exitPrice && formData.exitDate) {
			const exitDateTime = new Date(`${formData.exitDate}T${formData.exitTime || '16:00'}`);
			const exitPrice = parseFloat(formData.exitPrice);

			// Calculate P&L
			const priceDiff = formData.tradeType === 'LONG'
				? (exitPrice - entryPrice)
				: (entryPrice - exitPrice);
			const pnl = (priceDiff * quantity) - fees;

			tradeData.exit_price = exitPrice;
			tradeData.pnl = pnl;
			tradeData.closed_at = exitDateTime.toISOString();
			tradeData.exits = [
				{
					id: crypto.randomUUID(),
					price: exitPrice,
					quantity,
					timestamp: exitDateTime.toISOString(),
					fees: fees / 2,
					pnl
				}
			];
			tradeData.current_position_size = 0;
			tradeData.realized_pnl = pnl;
			tradeData.unrealized_pnl = null;
		}

		onSave(tradeData);
		onClose();
	}

	function addTag() {
		if (newTag.trim() && !formData.tags.includes(newTag.trim())) {
			formData.tags = [...formData.tags, newTag.trim()];
			newTag = '';
		}
	}

	function removeTag(tag: string) {
		formData.tags = formData.tags.filter((t: string) => t !== tag);
	}

	function resetForm() {
		formData = {
			symbol: '',
			tradeType: 'LONG',
			quantity: '',
			entryPrice: '',
			entryDate: new Date().toISOString().split('T')[0],
			entryTime: '09:30',
			exitPrice: '',
			exitDate: '',
			exitTime: '',
			fees: '0',
			notes: '',
			tags: [],
			costBasisMethod: 'FIFO'
		};
		error = '';
	}

	function handleClose() {
		resetForm();
		onClose();
	}
</script>

{#if isOpen}
	<div class="fixed inset-0 z-[100] overflow-y-auto" role="dialog" aria-modal="true">
		<div
			class="fixed inset-0 bg-black/50"
			onclick={onClose}
			onkeydown={(e) => e.key === 'Escape' && onClose()}
			role="button"
			tabindex="-1"
			aria-label="Close modal"
		></div>

		<div class="relative min-h-screen flex items-center justify-center p-4">
			<div
				class="relative bg-white dark:bg-surface-800 rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
			>
				<!-- Header -->
				<div
					class="flex items-center justify-between p-6 border-b border-surface-200 dark:border-surface-700"
				>
					<div>
						<h2 class="text-2xl font-bold">{trade ? 'Edit Position' : 'Open New Position'}</h2>
						<p class="text-sm text-surface-600 dark:text-surface-400 mt-1">
							{trade ? 'Modify position details' : 'Enter initial position entry details'}
						</p>
					</div>
					<button onclick={handleClose} class="p-2 hover:bg-surface-100 dark:hover:bg-surface-700 rounded">
						<Icon icon="mdi:close" width="24" />
					</button>
				</div>

				<!-- Form -->
				<form onsubmit={handleSubmit} class="p-6 space-y-6">
					{#if error}
						<div class="bg-error-50 dark:bg-error-900/20 border border-error-200 dark:border-error-800 text-error-700 dark:text-error-300 px-4 py-3 rounded">
							<Icon icon="mdi:alert-circle" class="inline mr-2" />
							{error}
						</div>
					{/if}

					<!-- Help Section -->
					{#if !trade}
						<HelpText
							type="info"
							title="Opening a New Position"
							text="Record your initial position entry here. You can add additional entries or exits later to track scaling in/out of positions. Leave exit details blank to keep the position open."
							collapsible={true}
						/>
					{/if}
					<!-- Trade Details -->
					<div>
						<div class="flex items-center gap-2 mb-4">
							<h3 class="text-lg font-semibold">Trade Details</h3>
							<Tooltip text="Enter the ticker symbol and choose LONG (buy) or SHORT (sell)" position="right">
								<Icon icon="mdi:help-circle-outline" width="18" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help" />
							</Tooltip>
						</div>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<Input
								label="Symbol"
								bind:value={formData.symbol}
								placeholder="AAPL"
								required={true}
							/>
							<div>
								<fieldset>
									<legend class="block text-sm font-medium mb-2 text-surface-700 dark:text-surface-300">
										Type *
									</legend>
									<div class="flex gap-4 mt-2">
										<label class="flex items-center gap-2 cursor-pointer">
											<input
												type="radio"
												name="tradeType"
												value="LONG"
												bind:group={formData.tradeType}
												class="radio"
											/>
											<span>LONG</span>
										</label>
										<label class="flex items-center gap-2 cursor-pointer">
											<input
												type="radio"
												name="tradeType"
												value="SHORT"
												bind:group={formData.tradeType}
												class="radio"
											/>
											<span>SHORT</span>
										</label>
									</div>
								</fieldset>
							</div>
						</div>
					</div>

					<!-- Entry Details -->
					<div>
						<div class="flex items-center gap-2 mb-4">
							<h3 class="text-lg font-semibold">Initial Entry</h3>
							<Tooltip text="Record the price, size, and exact time you entered this position" position="right" maxWidth="250px">
								<Icon icon="mdi:help-circle-outline" width="18" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help" />
							</Tooltip>
						</div>
						<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
							<Input
								label="Entry Price *"
								type="number"
								step="0.01"
								bind:value={formData.entryPrice}
								placeholder="150.00"
								required={true}
							/>
							<Input
								label="Quantity *"
								type="number"
								step="1"
								bind:value={formData.quantity}
								placeholder="100"
								required={true}
							/>
							<Input label="Date *" type="date" bind:value={formData.entryDate} required={true} />
							<Input label="Time" type="time" bind:value={formData.entryTime} />
						</div>

						<HelpText
							type="tip"
							text="You can add more entries later to track averaging into positions. Just select 'Add Entry/Exit' from the trade details page."
						/>
					</div>

					<!-- Exit Details -->
					<div>
						<div class="flex items-center justify-between mb-4">
							<h3 class="text-lg font-semibold">Exit Details (Optional)</h3>
							<p class="text-xs text-surface-600 dark:text-surface-400">
								Leave blank to keep position open
							</p>
						</div>
						<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
							<Input
								label="Exit Price"
								type="number"
								step="0.01"
								bind:value={formData.exitPrice}
								placeholder="155.00"
							/>
							<Input
								label="Quantity"
								type="number"
								value={formData.quantity}
								disabled={true}
							/>
							<Input label="Date" type="date" bind:value={formData.exitDate} />
							<Input label="Time" type="time" bind:value={formData.exitTime} />
						</div>
					</div>

					<!-- Additional -->
					<div>
						<div class="flex items-center gap-2 mb-4">
							<h3 class="text-lg font-semibold">Additional Settings</h3>
							<Tooltip text="Configure fees and cost basis calculation method for this position" position="right" maxWidth="250px">
								<Icon icon="mdi:help-circle-outline" width="18" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help" />
							</Tooltip>
						</div>
						<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
							<Input
								label="Total Fees"
								type="number"
								step="0.01"
								bind:value={formData.fees}
								placeholder="2.50"
							/>
							<Select
								label="Cost Basis Method"
								options={costBasisOptions}
								bind:value={formData.costBasisMethod}
							/>
						</div>

						<div class="mt-4">
							<HelpText
								type="info"
								title="Cost Basis Methods"
								text="FIFO (First In, First Out): Sells oldest shares first. LIFO (Last In, First Out): Sells newest shares first. Average: Uses average price of all shares. This affects P&L calculation when you scale in/out of positions."
								collapsible={true}
							/>
						</div>
					</div>

					<!-- Notes and Tags -->
					<div>
						<div class="flex items-center gap-2 mb-4">
							<h3 class="text-lg font-semibold">Notes & Tags</h3>
							<Tooltip text="Add context about your trade setup, plan, and observations" position="right" maxWidth="250px">
								<Icon icon="mdi:help-circle-outline" width="18" class="text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 cursor-help" />
							</Tooltip>
						</div>
						<div class="space-y-4">

							<!-- Tags -->
							<div>
								<div class="flex gap-2 mb-2">
									<Input
										label="Tags"
										value={newTag}
										onChange={(val) => (newTag = val)}
										placeholder="Add tag..."
										class="flex-1"
									/>
									<Button type="button" variant="soft" size="md" onclick={addTag} class="mt-6">
										<Icon icon="mdi:plus" width="20" />
									</Button>
								</div>
								{#if formData.tags.length > 0}
									<div class="flex flex-wrap gap-2">
										{#each formData.tags as tag}
											<Badge color="primary" variant="soft" removable={true} onRemove={() => removeTag(tag)}>
												{tag}
											</Badge>
										{/each}
									</div>
								{/if}
							</div>

							<!-- Notes -->
							<Textarea
								label="Notes"
								bind:value={formData.notes}
								placeholder="Add notes about this trade setup, plan, or observations..."
								rows={4}
							/>

							<HelpText
								type="tip"
								text="Document WHY you took this trade. Note your setup, indicators, market conditions, and your trading plan. This helps with post-trade analysis and improving your strategy."
							/>
						</div>
					</div>

					<!-- Actions -->
					<div class="flex gap-4 justify-end pt-4 border-t border-surface-200 dark:border-surface-700">
						<Button type="button" variant="ghost" onclick={handleClose}>Cancel</Button>
						<Button type="submit" color="primary">
							<Icon icon="mdi:content-save" class="mr-2" />
							{trade ? 'Update Position' : 'Open Position'}
						</Button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}
