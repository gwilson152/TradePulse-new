<script lang="ts">
	import Icon from '@iconify/svelte';
	import Card from './Card.svelte';
	import Button from './Button.svelte';
	import BetaBadge from './BetaBadge.svelte';

	interface Props {
		planType: 'starter' | 'pro' | 'premium';
		title: string;
		description: string;
		price: string;
		originalPrice?: string;
		features: string[];
		recommended?: boolean;
		selected?: boolean;
		onSelect?: () => void;
	}

	let {
		planType,
		title,
		description,
		price,
		originalPrice,
		features,
		recommended = false,
		selected = false,
		onSelect
	}: Props = $props();

	const planColors = {
		starter: 'from-blue-500 to-cyan-500',
		pro: 'from-purple-500 to-pink-500',
		premium: 'from-amber-500 to-orange-500'
	};
</script>

<Card
	variant={recommended ? 'elevated' : 'glass'}
	padding="lg"
	hover={!selected}
	class="relative h-full flex flex-col {recommended ? 'ring-2 ring-primary-500' : ''}"
>
	{#if recommended}
		<div class="absolute -top-4 left-1/2 transform -translate-x-1/2">
			<span
				class="inline-flex items-center gap-1 px-4 py-1 rounded-full bg-gradient-to-r from-primary-500 to-secondary-500 text-white text-sm font-semibold shadow-lg"
			>
				<Icon icon="mdi:star" />
				Most Popular
			</span>
		</div>
	{/if}

	<div class="flex flex-col flex-grow">
		<!-- Header -->
		<div class="text-center mb-6">
			<div class="flex items-center justify-center gap-2 mb-2">
				<h3 class="text-2xl font-bold text-slate-900 dark:text-white">{title}</h3>
				<BetaBadge size="sm" />
			</div>
			<p class="text-sm text-slate-600 dark:text-slate-400">{description}</p>
		</div>

		<!-- Price -->
		<div class="text-center mb-6">
			<div class="flex items-center justify-center gap-3">
				{#if originalPrice}
					<span
						class="text-2xl text-slate-400 dark:text-slate-500 line-through font-semibold"
					>
						{originalPrice}
					</span>
				{/if}
				<div class="flex items-baseline">
					<span
						class="text-4xl font-bold bg-gradient-to-r {planColors[
							planType
						]} bg-clip-text text-transparent"
					>
						{price}
					</span>
					<span class="text-slate-500 dark:text-slate-400 ml-1">/month</span>
				</div>
			</div>
			<p class="text-xs text-success-600 dark:text-success-400 font-semibold mt-2">
				Free during Beta
			</p>
		</div>

		<!-- Features -->
		<div class="flex-grow mb-6">
			<ul class="space-y-3">
				{#each features as feature}
					<li class="flex items-start gap-2">
						<Icon
							icon="mdi:check-circle"
							class="text-success-500 mt-0.5 flex-shrink-0"
							width="20"
						/>
						<span class="text-sm text-slate-700 dark:text-slate-300">{feature}</span>
					</li>
				{/each}
			</ul>
		</div>

		<!-- CTA Button -->
		<div class="mt-auto">
			<Button
				color={selected ? 'secondary' : recommended ? 'primary' : 'neutral'}
				variant={selected ? 'soft' : 'filled'}
				onclick={onSelect}
				class="w-full justify-center"
				disabled={selected}
			>
				{#if selected}
					<Icon icon="mdi:check" width="20" class="mr-2" />
					Selected
				{:else}
					Choose {title}
				{/if}
			</Button>
		</div>
	</div>
</Card>

<style>
	.ring-2 {
		box-shadow:
			0 0 0 2px rgba(var(--color-primary-500), 0.5),
			0 20px 25px -5px rgba(0, 0, 0, 0.1);
	}
</style>
