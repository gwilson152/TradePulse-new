<script lang="ts">
	import Icon from '@iconify/svelte';

	interface Props {
		src: string;
		label?: string;
		onDelete?: () => void;
	}

	let { src, label, onDelete }: Props = $props();

	let audio: HTMLAudioElement;
	let playing = $state(false);
	let currentTime = $state(0);
	let duration = $state(0);
	let loading = $state(true);

	function togglePlay() {
		if (playing) {
			audio.pause();
		} else {
			audio.play();
		}
	}

	function handleTimeUpdate() {
		currentTime = audio.currentTime;
	}

	function handleLoadedMetadata() {
		duration = audio.duration;
		loading = false;
	}

	function handleEnded() {
		playing = false;
		currentTime = 0;
	}

	function handleSeek(e: Event) {
		const target = e.target as HTMLInputElement;
		const time = parseFloat(target.value);
		audio.currentTime = time;
		currentTime = time;
	}

	function formatTime(seconds: number): string {
		if (isNaN(seconds)) return '0:00';
		const mins = Math.floor(seconds / 60);
		const secs = Math.floor(seconds % 60);
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}

	function handlePlayPause() {
		playing = !audio.paused;
	}
</script>

<div class="audio-player bg-surface-50 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700 p-4">
	<audio
		bind:this={audio}
		{src}
		ontimeupdate={handleTimeUpdate}
		onloadedmetadata={handleLoadedMetadata}
		onended={handleEnded}
		onplay={handlePlayPause}
		onpause={handlePlayPause}
	></audio>

	<div class="flex items-center gap-3">
		<!-- Play/Pause button -->
		<button
			onclick={togglePlay}
			disabled={loading}
			class="play-button flex items-center justify-center w-12 h-12 rounded-full bg-primary-500 hover:bg-primary-600 disabled:bg-surface-400 disabled:cursor-not-allowed text-white transition-colors"
			title={playing ? 'Pause' : 'Play'}
		>
			{#if loading}
				<Icon icon="mdi:loading" class="animate-spin text-2xl" />
			{:else if playing}
				<Icon icon="mdi:pause" class="text-2xl" />
			{:else}
				<Icon icon="mdi:play" class="text-2xl ml-1" />
			{/if}
		</button>

		<!-- Progress and controls -->
		<div class="flex-1">
			{#if label}
				<div class="text-sm font-medium text-surface-900 dark:text-surface-100 mb-1">
					{label}
				</div>
			{/if}

			<div class="flex items-center gap-2">
				<span class="text-xs text-surface-600 dark:text-surface-400 tabular-nums min-w-[40px]">
					{formatTime(currentTime)}
				</span>

				<input
					type="range"
					min="0"
					max={duration || 0}
					value={currentTime}
					oninput={handleSeek}
					disabled={loading}
					class="flex-1 h-2 bg-surface-200 dark:bg-surface-700 rounded-full appearance-none cursor-pointer disabled:cursor-not-allowed slider"
				/>

				<span class="text-xs text-surface-600 dark:text-surface-400 tabular-nums min-w-[40px] text-right">
					{formatTime(duration)}
				</span>
			</div>
		</div>

		<!-- Delete button -->
		{#if onDelete}
			<button
				onclick={onDelete}
				class="p-2 hover:bg-loss-100 dark:hover:bg-loss-900/30 rounded-full transition-colors"
				title="Delete audio"
			>
				<Icon icon="mdi:delete" class="text-loss-600 dark:text-loss-400" />
			</button>
		{/if}
	</div>
</div>

<style>
	/* Custom range slider styling */
	.slider {
		-webkit-appearance: none;
		appearance: none;
		background: transparent;
		outline: none;
	}

	.slider::-webkit-slider-track {
		width: 100%;
		height: 8px;
		background: rgb(var(--color-surface-200));
		border-radius: 4px;
	}

	:global(.dark) .slider::-webkit-slider-track {
		background: rgb(var(--color-surface-700));
	}

	.slider::-webkit-slider-thumb {
		-webkit-appearance: none;
		appearance: none;
		width: 16px;
		height: 16px;
		background: rgb(var(--color-primary-500));
		border-radius: 50%;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.slider::-webkit-slider-thumb:hover {
		transform: scale(1.2);
		background: rgb(var(--color-primary-600));
	}

	.slider::-moz-range-track {
		width: 100%;
		height: 8px;
		background: rgb(var(--color-surface-200));
		border-radius: 4px;
	}

	:global(.dark) .slider::-moz-range-track {
		background: rgb(var(--color-surface-700));
	}

	.slider::-moz-range-thumb {
		width: 16px;
		height: 16px;
		background: rgb(var(--color-primary-500));
		border: none;
		border-radius: 50%;
		cursor: pointer;
		transition: all 0.2s ease;
	}

	.slider::-moz-range-thumb:hover {
		transform: scale(1.2);
		background: rgb(var(--color-primary-600));
	}

	.play-button {
		transition: all 0.2s ease;
	}

	.play-button:hover:not(:disabled) {
		transform: scale(1.05);
	}

	.play-button:active:not(:disabled) {
		transform: scale(0.95);
	}
</style>
