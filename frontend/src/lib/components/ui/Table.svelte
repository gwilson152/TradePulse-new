<script lang="ts">
	interface Props {
		hover?: boolean;
		interactive?: boolean;
		children?: any;
	}

	let { hover = true, interactive = false, children }: Props = $props();

	const classes = `table ${hover ? 'table-hover' : ''} ${interactive ? 'table-interactive' : ''}`;
</script>

<div class="table-container">
	<table class={classes}>
		{@render children?.()}
	</table>
</div>
