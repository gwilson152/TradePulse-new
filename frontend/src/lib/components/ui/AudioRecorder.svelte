<script lang="ts">
	import Icon from '@iconify/svelte';
	import Button from './Button.svelte';

	interface Props {
		onRecordingComplete: (blob: Blob) => void;
		maxDuration?: number; // in seconds
	}

	let { onRecordingComplete, maxDuration = 300 }: Props = $props();

	let recording = $state(false);
	let paused = $state(false);
	let duration = $state(0);
	let mediaRecorder: MediaRecorder | null = null;
	let audioChunks: Blob[] = [];
	let stream: MediaStream | null = null;
	let timerInterval: number | null = null;
	let error = $state('');

	async function startRecording() {
		try {
			error = '';
			stream = await navigator.mediaDevices.getUserMedia({ audio: true });

			mediaRecorder = new MediaRecorder(stream);
			audioChunks = [];

			mediaRecorder.ondataavailable = (event) => {
				if (event.data.size > 0) {
					audioChunks.push(event.data);
				}
			};

			mediaRecorder.onstop = () => {
				const audioBlob = new Blob(audioChunks, { type: 'audio/webm' });
				onRecordingComplete(audioBlob);
				cleanup();
			};

			mediaRecorder.start();
			recording = true;
			duration = 0;

			// Start timer
			timerInterval = window.setInterval(() => {
				duration++;
				// Auto-stop at max duration
				if (duration >= maxDuration) {
					stopRecording();
				}
			}, 1000);
		} catch (err) {
			error = 'Could not access microphone. Please check permissions.';
			console.error('Error accessing microphone:', err);
		}
	}

	function pauseRecording() {
		if (mediaRecorder && recording && !paused) {
			mediaRecorder.pause();
			paused = true;
			if (timerInterval) clearInterval(timerInterval);
		}
	}

	function resumeRecording() {
		if (mediaRecorder && recording && paused) {
			mediaRecorder.resume();
			paused = false;
			// Restart timer
			timerInterval = window.setInterval(() => {
				duration++;
				if (duration >= maxDuration) {
					stopRecording();
				}
			}, 1000);
		}
	}

	function stopRecording() {
		if (mediaRecorder && recording) {
			mediaRecorder.stop();
			recording = false;
			paused = false;
		}
	}

	function cancelRecording() {
		if (mediaRecorder) {
			mediaRecorder.stop();
			audioChunks = [];
		}
		cleanup();
	}

	function cleanup() {
		if (timerInterval) {
			clearInterval(timerInterval);
			timerInterval = null;
		}
		if (stream) {
			stream.getTracks().forEach((track) => track.stop());
			stream = null;
		}
		recording = false;
		paused = false;
		duration = 0;
	}

	function formatTime(seconds: number): string {
		const mins = Math.floor(seconds / 60);
		const secs = seconds % 60;
		return `${mins}:${secs.toString().padStart(2, '0')}`;
	}

	// Cleanup on unmount
	$effect(() => {
		return () => {
			cleanup();
		};
	});
</script>

<div class="audio-recorder bg-surface-50 dark:bg-surface-800 rounded-lg border border-surface-200 dark:border-surface-700 p-4">
	{#if error}
		<div class="mb-3 text-sm text-loss-600 dark:text-loss-400 flex items-center gap-2">
			<Icon icon="mdi:alert-circle" />
			{error}
		</div>
	{/if}

	{#if !recording}
		<!-- Not recording state -->
		<div class="text-center">
			<div class="mb-3">
				<Icon
					icon="mdi:microphone"
					class="text-6xl text-surface-400 dark:text-surface-500 mx-auto"
				/>
			</div>
			<p class="text-sm text-surface-600 dark:text-surface-400 mb-4">
				Record a voice note (max {Math.floor(maxDuration / 60)} min)
			</p>
			<Button onclick={startRecording} color="primary" size="md">
				<Icon icon="mdi:record" class="mr-2" />
				Start Recording
			</Button>
		</div>
	{:else}
		<!-- Recording state -->
		<div class="text-center">
			<!-- Animated recording indicator -->
			<div class="mb-4 flex items-center justify-center gap-3">
				<div class="recording-pulse">
					<Icon icon="mdi:microphone" class="text-3xl text-white" />
				</div>
				<div>
					<div class="text-2xl font-bold text-surface-900 dark:text-surface-100">
						{formatTime(duration)}
					</div>
					<div class="text-xs text-surface-600 dark:text-surface-400">
						{paused ? 'Paused' : 'Recording...'}
					</div>
				</div>
			</div>

			<!-- Progress bar -->
			<div class="mb-4 h-2 bg-surface-200 dark:bg-surface-700 rounded-full overflow-hidden">
				<div
					class="h-full bg-error-500 transition-all duration-1000"
					style="width: {(duration / maxDuration) * 100}%"
				></div>
			</div>

			<!-- Control buttons -->
			<div class="flex items-center justify-center gap-2">
				<Button onclick={cancelRecording} variant="ghost" size="sm">
					<Icon icon="mdi:close" class="mr-1" />
					Cancel
				</Button>

				{#if !paused}
					<Button onclick={pauseRecording} color="warning" size="sm">
						<Icon icon="mdi:pause" class="mr-1" />
						Pause
					</Button>
				{:else}
					<Button onclick={resumeRecording} color="primary" size="sm">
						<Icon icon="mdi:play" class="mr-1" />
						Resume
					</Button>
				{/if}

				<Button onclick={stopRecording} color="success" size="sm">
					<Icon icon="mdi:stop" class="mr-1" />
					Stop & Save
				</Button>
			</div>
		</div>
	{/if}
</div>

<style>
	.recording-pulse {
		width: 60px;
		height: 60px;
		border-radius: 50%;
		background-color: rgb(var(--color-error-500));
		display: flex;
		align-items: center;
		justify-content: center;
		animation: pulse 1.5s ease-in-out infinite;
	}

	@keyframes pulse {
		0%,
		100% {
			transform: scale(1);
			opacity: 1;
		}
		50% {
			transform: scale(1.1);
			opacity: 0.8;
		}
	}
</style>
