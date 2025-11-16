<script lang="ts">
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/ui/Input.svelte';
	import Button from '$lib/components/ui/Button.svelte';
	import Icon from '@iconify/svelte';
	import { apiClient } from '$lib/api/client';
	import { toast } from '$lib/stores/toast';

	interface Props {
		open: boolean;
		hasPassword: boolean;
		onClose: () => void;
		onSuccess: () => void;
	}

	let { open = false, hasPassword = false, onClose, onSuccess }: Props = $props();

	let password = $state('');
	let confirmPassword = $state('');
	let loading = $state(false);
	let error = $state('');

	// Password strength indicator
	const passwordStrength = $derived(() => {
		if (!password) return { level: 0, text: '', color: '' };

		let strength = 0;
		if (password.length >= 8) strength++;
		if (password.length >= 12) strength++;
		if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength++;
		if (/\d/.test(password)) strength++;
		if (/[^a-zA-Z0-9]/.test(password)) strength++;

		if (strength <= 2) return { level: strength, text: 'Weak', color: 'text-red-600' };
		if (strength === 3) return { level: strength, text: 'Fair', color: 'text-orange-600' };
		if (strength === 4) return { level: strength, text: 'Good', color: 'text-blue-600' };
		return { level: strength, text: 'Strong', color: 'text-green-600' };
	});

	async function handleSubmit() {
		error = '';

		// Validation
		if (password.length < 8) {
			error = 'Password must be at least 8 characters';
			return;
		}

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		loading = true;

		try {
			await apiClient.setPassword(password);
			toast.success(hasPassword ? 'Password updated successfully' : 'Password set successfully');
			resetForm();
			onSuccess();
			onClose();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to set password';
			toast.error('Failed to set password');
		} finally {
			loading = false;
		}
	}

	function resetForm() {
		password = '';
		confirmPassword = '';
		error = '';
	}

	function handleClose() {
		resetForm();
		onClose();
	}
</script>

<Modal {open} onClose={handleClose} title={hasPassword ? 'Change Password' : 'Set Password'} size="md">
	<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="space-y-4">
		{#if !hasPassword}
			<div class="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
				<div class="flex items-start gap-3">
					<Icon icon="mdi:information" width="20" class="text-blue-600 dark:text-blue-400 mt-0.5 flex-shrink-0" />
					<div class="flex-1">
						<p class="text-sm text-blue-700 dark:text-blue-300">
							Setting a password allows you to sign in without waiting for a magic link email.
						</p>
					</div>
				</div>
			</div>
		{/if}

		{#if error}
			<div class="p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
				<div class="flex items-start gap-2">
					<Icon icon="mdi:alert-circle" width="20" class="text-red-600 dark:text-red-400 mt-0.5 flex-shrink-0" />
					<p class="text-sm text-red-700 dark:text-red-300">{error}</p>
				</div>
			</div>
		{/if}

		<div>
			<Input
				type="password"
				bind:value={password}
				label="New Password"
				placeholder="••••••••"
				required={true}
				disabled={loading}
			/>
			{#if password}
				<div class="mt-2">
					<div class="flex items-center justify-between mb-1">
						<span class="text-xs text-surface-600 dark:text-surface-400">Password strength:</span>
						<span class="text-xs font-medium {passwordStrength().color}">
							{passwordStrength().text}
						</span>
					</div>
					<div class="h-1.5 bg-surface-200 dark:bg-surface-700 rounded-full overflow-hidden">
						<div
							class="h-full transition-all duration-300 {
								passwordStrength().level <= 2 ? 'bg-red-500' :
								passwordStrength().level === 3 ? 'bg-orange-500' :
								passwordStrength().level === 4 ? 'bg-blue-500' :
								'bg-green-500'
							}"
							style="width: {(passwordStrength().level / 5) * 100}%"
						></div>
					</div>
				</div>
			{/if}
		</div>

		<Input
			type="password"
			bind:value={confirmPassword}
			label="Confirm Password"
			placeholder="••••••••"
			required={true}
			disabled={loading}
		/>

		<div class="pt-2 text-xs text-surface-600 dark:text-surface-400 space-y-1">
			<p>Password requirements:</p>
			<ul class="list-disc list-inside space-y-1">
				<li class={password.length >= 8 ? 'text-green-600' : ''}>At least 8 characters</li>
				<li class={/[a-z]/.test(password) && /[A-Z]/.test(password) ? 'text-green-600' : ''}>
					Mix of uppercase and lowercase
				</li>
				<li class={/\d/.test(password) ? 'text-green-600' : ''}>At least one number</li>
			</ul>
		</div>

		<div class="flex gap-3 pt-4">
			<Button type="button" variant="ghost" onclick={handleClose} disabled={loading} class="flex-1">
				Cancel
			</Button>
			<Button
				type="submit"
				color="primary"
				disabled={loading || !password || !confirmPassword || password !== confirmPassword}
				class="flex-1"
			>
				{#if loading}
					<Icon icon="mdi:loading" width="20" class="animate-spin mr-2" />
					{hasPassword ? 'Updating...' : 'Setting...'}
				{:else}
					{hasPassword ? 'Update Password' : 'Set Password'}
				{/if}
			</Button>
		</div>
	</form>
</Modal>
