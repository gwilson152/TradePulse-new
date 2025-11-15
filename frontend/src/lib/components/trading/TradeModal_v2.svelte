<script lang="ts">
	import Icon from '@iconify/svelte';
	import FormModal from '../ui/FormModal.svelte';
	import FormSection from '../ui/FormSection.svelte';
	import Input from '../ui/Input.svelte';
	import Select from '../ui/Select.svelte';
	import Textarea from '../ui/Textarea.svelte';
	import Badge from '../ui/Badge.svelte';
	import Button from '../ui/Button.svelte';
	import HelpText from '../ui/HelpText.svelte';
	import type { Trade, CostBasisMethod } from '$lib/types';
	import { toast } from '$lib/stores/toast';

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
		entryDate: trade?.opened_at
			? new Date(trade.opened_at).toISOString().split('T')[0]
			: new Date().toISOString().split('T')[0],
		entryTime: trade?.opened_at
			? new Date(trade.opened_at).toTimeString().split(' ')[0].slice(0, 5)
			: '09:30',
		exitPrice: trade?.exit_price || '',
		exitDate: trade?.closed_at ? new Date(trade.closed_at).toISOString().split('T')[0] : '',
		exitTime: trade?.closed_at
			? new Date(trade.closed_at).toTimeString().split(' ')[0].slice(0, 5)
			: '',
		fees: trade?.fees?.toString() || '0',
		notes: trade?.notes || '',
		tags: trade?.tags || [],
		costBasisMethod: (trade?.cost_basis_method || 'FIFO') as CostBasisMethod
	});

	let newTag = $state('');
	let error = $state('');
	let loading = $state(false);

	const costBasisOptions = [
		{ value: 'FIFO', label: 'First In First Out (FIFO)' },
		{ value: 'LIFO', label: 'Last In First Out (LIFO)' },
		{ value: 'AVERAGE', label: 'Average Cost' }
	];

	function handleSubmit() {
		error = '';

		// Validation
		if (!formData.symbol || !formData.quantity || !formData.entryPrice) {
			error = 'Please fill in all required fields';
			return;
		}

		loading = true;

		try {
			// Build trade data
			const entryDateTime = new Date(`${formData.entryDate}T${formData.entryTime}`);
			const quantity = parseFloat(formData.quantity.toString());
			const entryPrice = parseFloat(formData.entryPrice.toString());
			const fees = parseFloat(formData.fees.toString());

			const tradeData: Partial<Trade> = {
				symbol: formData.symbol.toUpperCase(),
				trade_type: formData.tradeType as 'LONG' | 'SHORT',
				quantity,
				entry_price: entryPrice,
				fees,
				entries: [
					{
						id: crypto.randomUUID(),
						price: entryPrice,
						quantity,
						timestamp: entryDateTime.toISOString(),
						fees: fees / 2
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
				const exitPrice = parseFloat(formData.exitPrice.toString());

				const priceDiff =
					formData.tradeType === 'LONG' ? exitPrice - entryPrice : entryPrice - exitPrice;
				const pnl = priceDiff * quantity - fees;

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
			toast.success(trade ? 'Trade updated successfully' : 'Trade created successfully');
			handleClose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to save trade';
		} finally {
			loading = false;
		}
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

	function handleClose() {
		resetForm();
		onClose();
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
		loading = false;
	}
</script>

<FormModal
	open={isOpen}
	title={trade ? 'Edit Position' : 'Open New Position'}
	subtitle={trade ? 'Modify position details' : 'Enter initial position entry details'}
	size="xl"
	{error}
	{loading}
	showHelp={!trade}
	helpTitle="Opening a New Position"
	helpText="Record your initial position entry here. You can add additional entries or exits later to track scaling in/out of positions. Leave exit details blank to keep the position open."
	onClose={handleClose}
	onSubmit={handleSubmit}
	submitText={trade ? 'Update Position' : 'Open Position'}
>
	<!-- Trade Details -->
	<FormSection
		title="Trade Details"
		icon="mdi:chart-line"
		helpText="Enter the ticker symbol and choose LONG (buy) or SHORT (sell)"
	>
		<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
			<Input label="Symbol" bind:value={formData.symbol} placeholder="AAPL" required={true} />
			<div>
				<fieldset>
					<legend class="block text-sm font-semibold mb-2 text-slate-700 dark:text-slate-300">
						Type *
					</legend>
					<div class="flex gap-4 mt-2">
						<label class="flex items-center gap-2 cursor-pointer">
							<input
								type="radio"
								name="tradeType"
								value="LONG"
								bind:group={formData.tradeType}
								class="w-4 h-4 text-blue-600 focus:ring-blue-500"
							/>
							<span class="text-sm font-medium text-slate-700 dark:text-slate-300">LONG</span>
						</label>
						<label class="flex items-center gap-2 cursor-pointer">
							<input
								type="radio"
								name="tradeType"
								value="SHORT"
								bind:group={formData.tradeType}
								class="w-4 h-4 text-blue-600 focus:ring-blue-500"
							/>
							<span class="text-sm font-medium text-slate-700 dark:text-slate-300">SHORT</span>
						</label>
					</div>
				</fieldset>
			</div>
		</div>
	</FormSection>

	<!-- Entry Details -->
	<FormSection
		title="Initial Entry"
		icon="mdi:login"
		helpText="Record the price, size, and exact time you entered this position"
	>
		<div class="grid grid-cols-1 md:grid-cols-4 gap-4">
			<Input
				label="Entry Price"
				type="number"
				step="0.01"
				bind:value={formData.entryPrice}
				placeholder="150.00"
				required={true}
			/>
			<Input
				label="Quantity"
				type="number"
				step="1"
				bind:value={formData.quantity}
				placeholder="100"
				required={true}
			/>
			<Input label="Date" type="date" bind:value={formData.entryDate} required={true} />
			<Input label="Time" type="time" bind:value={formData.entryTime} />
		</div>
		<HelpText
			type="tip"
			text="You can add more entries later to track averaging into positions. Just select 'Add Entry/Exit' from the trade details page."
		/>
	</FormSection>

	<!-- Exit Details -->
	<FormSection title="Exit Details (Optional)" icon="mdi:logout" collapsible={true}>
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
	</FormSection>

	<!-- Additional Settings -->
	<FormSection
		title="Additional Settings"
		icon="mdi:cog"
		helpText="Configure fees and cost basis calculation method for this position"
		collapsible={true}
	>
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
		<HelpText
			type="info"
			title="Cost Basis Methods"
			text="FIFO (First In, First Out): Sells oldest shares first. LIFO (Last In, First Out): Sells newest shares first. Average: Uses average price of all shares. This affects P&L calculation when you scale in/out of positions."
			collapsible={true}
		/>
	</FormSection>

	<!-- Notes & Tags -->
	<FormSection
		title="Notes & Tags"
		icon="mdi:note-text"
		helpText="Add context about your trade setup, plan, and observations"
		collapsible={true}
	>
		<div class="space-y-4">
			<!-- Tags -->
			<div>
				<div class="flex gap-2 mb-2">
					<Input
						label="Tags"
						bind:value={newTag}
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
	</FormSection>
</FormModal>
