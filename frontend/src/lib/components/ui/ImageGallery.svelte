<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		images: string[];
		alt?: string;
		onDelete?: (index: number) => void;
	}

	let { images = [], alt = 'Image', onDelete }: Props = $props();

	let lightboxOpen = $state(false);
	let currentIndex = $state(0);

	function openLightbox(index: number) {
		currentIndex = index;
		lightboxOpen = true;
	}

	function closeLightbox() {
		lightboxOpen = false;
	}

	function nextImage() {
		currentIndex = (currentIndex + 1) % images.length;
	}

	function prevImage() {
		currentIndex = (currentIndex - 1 + images.length) % images.length;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (!lightboxOpen) return;

		if (e.key === 'Escape') closeLightbox();
		if (e.key === 'ArrowRight') nextImage();
		if (e.key === 'ArrowLeft') prevImage();
	}

	function handleDelete(index: number, e: Event) {
		e.stopPropagation();
		if (onDelete) {
			onDelete(index);
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

<div class="image-gallery">
	{#if images.length === 0}
		<div
			class="empty-state text-center py-8 text-surface-600 dark:text-surface-400 bg-surface-50 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700"
		>
			<Icon icon="mdi:image-off" class="text-4xl mb-2 mx-auto" />
			<p class="text-sm">No images available</p>
		</div>
	{:else}
		<div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
			{#each images as image, index}
				<div class="image-card group relative">
					<button
						type="button"
						onclick={() => openLightbox(index)}
						class="block w-full aspect-square overflow-hidden rounded-lg border border-surface-200 dark:border-surface-700 hover:border-surface-300 dark:hover:border-surface-600 transition-all hover:shadow-lg"
					>
						<img
							src={image}
							{alt}
							class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
						/>
					</button>

					{#if onDelete}
						<button
							type="button"
							onclick={(e) => handleDelete(index, e)}
							class="absolute top-2 right-2 p-1.5 bg-loss-500 hover:bg-loss-600 text-white rounded-full opacity-0 group-hover:opacity-100 transition-opacity shadow-lg"
							title="Delete image"
						>
							<Icon icon="mdi:delete" class="text-lg" />
						</button>
					{/if}
				</div>
			{/each}
		</div>
	{/if}

	<!-- Lightbox -->
	{#if lightboxOpen}
		<div
			class="lightbox fixed inset-0 z-50 bg-black/90 flex items-center justify-center p-4"
			onclick={closeLightbox}
			onkeydown={(e) => e.key === 'Escape' && closeLightbox()}
			role="dialog"
			aria-modal="true"
			tabindex="-1"
		>
			<button
				type="button"
				onclick={closeLightbox}
				class="absolute top-4 right-4 p-2 text-white hover:bg-white/10 rounded-full transition-colors z-10"
				title="Close"
			>
				<Icon icon="mdi:close" class="text-3xl" />
			</button>

			{#if images.length > 1}
				<button
					type="button"
					onclick={(e) => {
						e.stopPropagation();
						prevImage();
					}}
					class="absolute left-4 p-3 text-white hover:bg-white/10 rounded-full transition-colors z-10"
					title="Previous"
				>
					<Icon icon="mdi:chevron-left" class="text-4xl" />
				</button>

				<button
					type="button"
					onclick={(e) => {
						e.stopPropagation();
						nextImage();
					}}
					class="absolute right-4 p-3 text-white hover:bg-white/10 rounded-full transition-colors z-10"
					title="Next"
				>
					<Icon icon="mdi:chevron-right" class="text-4xl" />
				</button>
			{/if}

			<div
				class="relative max-w-7xl max-h-full"
				onclick={(e) => e.stopPropagation()}
				onkeydown={() => {}}
				role="img"
				aria-label="Lightbox image"
			>
				<img
					src={images[currentIndex]}
					{alt}
					class="max-w-full max-h-[90vh] object-contain rounded-lg"
				/>

				{#if images.length > 1}
					<div
						class="absolute bottom-4 left-1/2 transform -translate-x-1/2 bg-black/50 text-white px-4 py-2 rounded-full text-sm"
					>
						{currentIndex + 1} / {images.length}
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

<style>
	.lightbox {
		animation: fadeIn 0.2s ease;
	}

	@keyframes fadeIn {
		from {
			opacity: 0;
		}
		to {
			opacity: 1;
		}
	}

	.image-card {
		transition: all 0.2s ease;
	}

	.image-card:hover {
		transform: translateY(-2px);
	}
</style>
