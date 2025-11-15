<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import type { Trade, Entry, Exit } from '$lib/types';
	import Icon from '@iconify/svelte';

	interface Props {
		open: boolean;
		trade: Trade;
		action: 'entry' | 'exit';
		onClose: () => void;
		onSubmit: (data: Partial<Entry> | Partial<Exit>) => Promise<void>;
	}

	let { open = false, trade, action, onClose, onSubmit }: Props = $props();

	let formData = $state({
		price: '',
		quantity: '',
		timestamp: new Date().toISOString().slice(0, 16), // Format for datetime-local input
		notes: '',
		fees: ''
	});

	let loading = $state(false);
	let error = $state('');

	// Calculate estimated P&L for exits
	const estimatedPnl = $derived(() => {
		if (action === 'exit' && formData.price && formData.quantity) {
			const exitPrice = parseFloat(formData.price);
			const quantity = parseFloat(formData.quantity);
			const avgEntry = trade.average_entry_price || trade.entry_price;
			const fees = parseFloat(formData.fees) || 0;

			if (trade.trade_type === 'LONG') {
				return (exitPrice - avgEntry) * quantity - fees;
			} else {
				return (avgEntry - exitPrice) * quantity - fees;
			}
		}
		return null;
	});

	async function handleSubmit() {
		error = '';

		// Validation
		if (!formData.price || !formData.quantity) {
			error = 'Price and quantity are required';
			return;
		}

		const price = parseFloat(formData.price);
		const quantity = parseFloat(formData.quantity);

		if (price <= 0 || quantity <= 0) {
			error = 'Price and quantity must be greater than zero';
			return;
		}

		// For exits, check that we don't exit more than current position
		if (action === 'exit' && quantity > trade.current_position_size) {
			error = `Cannot exit ${quantity} shares. Current position: ${trade.current_position_size}`;
			return;
		}

		loading = true;

		try {
			const data: Partial<Entry> | Partial<Exit> = {
				price,
				quantity,
				timestamp: new Date(formData.timestamp).toISOString(),
				notes: formData.notes || undefined,
				fees: formData.fees ? parseFloat(formData.fees) : undefined
			};

			// Add P&L for exits
			if (action === 'exit' && estimatedPnl() !== null) {
				(data as Partial<Exit>).pnl = estimatedPnl()!;
			}

			await onSubmit(data);
			resetForm();
			onClose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to add to position';
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		formData = {
			price: '',
			quantity: '',
			timestamp: new Date().toISOString().slice(0, 16),
			notes: '',
			fees: ''
		};
		error = '';
	}

	function handleClose() {
		resetForm();
		onClose();
	}
</script>

<Modal {open} onClose={handleClose} title="{action === 'entry' ? 'Add to Position' : 'Reduce Position'} - {trade.symbol}">
	<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
		{#if error}
			<div class="bg-error-50 dark:bg-error-900/20 border border-error-200 dark:border-error-800 text-error-700 dark:text-error-300 px-4 py-3 rounded">
				<Icon icon="mdi:alert-circle" class="inline mr-2" />
				{error}
			</div>
		{/if}

		<div class="grid grid-cols-2 gap-4">
			<Input
				label="Price *"
				type="number"
				step="0.01"
				min="0"
				bind:value={formData.price}
				placeholder="0.00"
				required
			/>
			<Input
				label="Quantity *"
				type="number"
				step="1"
				min="1"
				bind:value={formData.quantity}
				placeholder="0"
				required
			/>
		</div>

		<Input
			label="Timestamp *"
			type="datetime-local"
			bind:value={formData.timestamp}
			required
		/>

		<Input
			label="Fees"
			type="number"
			step="0.01"
			min="0"
			bind:value={formData.fees}
			placeholder="0.00"
		/>

		<div>
			<label class="block text-sm font-medium mb-1 text-surface-700 dark:text-surface-300">
				Notes
			</label>
			<textarea
				bind:value={formData.notes}
				class="w-full px-3 py-2 border border-surface-300 dark:border-surface-600 rounded-lg bg-surface-50 dark:bg-surface-800 text-surface-900 dark:text-surface-100"
				rows="3"
				placeholder="Optional notes about this {action}..."
			></textarea>
		</div>

		{#if action === 'exit' && estimatedPnl() !== null}
			<div class="p-4 bg-surface-100 dark:bg-surface-800 rounded-lg border border-surface-300 dark:border-surface-700">
				<div class="flex justify-between items-center">
					<span class="text-sm font-medium text-surface-700 dark:text-surface-300">
						Estimated P&L:
					</span>
					<span
						class="text-lg font-bold {estimatedPnl()! >= 0
							? 'text-profit-600'
							: 'text-loss-600'}"
					>
						{estimatedPnl()! >= 0 ? '+' : ''}${estimatedPnl()!.toFixed(2)}
					</span>
				</div>
				<p class="text-xs text-surface-600 dark:text-surface-400 mt-1">
					Based on avg entry price: ${trade.average_entry_price?.toFixed(2) || trade.entry_price.toFixed(2)}
				</p>
			</div>
		{/if}

		{#if action === 'exit'}
			<div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded border border-blue-200 dark:border-blue-800">
				<p class="text-sm text-blue-700 dark:text-blue-300">
					<Icon icon="mdi:information" class="inline mr-1" />
					Current position: <strong>{trade.current_position_size} shares</strong>
				</p>
			</div>
		{/if}

		<div class="flex gap-2 pt-4">
			<Button type="button" variant="ghost" onclick={handleClose} class="flex-1">
				Cancel
			</Button>
			<Button type="submit" color="primary" disabled={loading} class="flex-1">
				{#if loading}
					<Icon icon="mdi:loading" class="animate-spin mr-2" />
					Processing...
				{:else}
					{action === 'entry' ? 'Add Entry' : 'Add Exit'}
				{/if}
			</Button>
		</div>
	</form>
</Modal>
