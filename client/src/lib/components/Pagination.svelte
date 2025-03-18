<script lang="ts">
	import Button from './Button.svelte';

	export let currentPage: number;
	export let totalPages: number;
	export let onPageChange: (page: number, event?: Event) => void;

	function handlePageChange(page: number, event: Event) {
		event.preventDefault();
		onPageChange(page, event);
	}
</script>

{#if totalPages > 1}
	<div
		class="flex justify-center items-center gap-4 bg-retroGray p-4 rounded-lg border-4 border-retroBlack shadow-[4px_4px_0_#455A64]"
	>
		<!-- Previous Button -->
		<button
			class={`w-12 h-12 flex items-center justify-center bg-retroCoral text-retroCream rounded-full border-2 border-retroBlack hover:bg-retroCoral disabled:text-black disabled:bg-white disabled:opacity-50 transition-all duration-300 font-pixel text-lg cursor-pointer disabled:cursor-not-allowed`}
			disabled={currentPage === 1}
			on:click={(e) => handlePageChange(currentPage - 1, e)}
		>
			<span class="">-</span>
		</button>

		<!-- Page Indicator -->
		<div class="flex items-center bg-retroCream px-4 py-2 rounded-md border-4 border-retroBlack">
			<span class="font-pixel text-retroGray text-lg">
				{currentPage} / {totalPages}
			</span>
		</div>

		<!-- Next Button -->
		<button
			class={`w-12 h-12 flex items-center justify-center bg-retroCoral text-retroCream rounded-full border-2 border-retroBlack hover:bg-retroCoral disabled:text-black disabled:bg-white disabled:opacity-50 transition-all duration-300 font-pixel text-lg cursor-pointer disabled:cursor-not-allowed`}
			disabled={currentPage === totalPages}
			on:click={(e) => handlePageChange(currentPage + 1, e)}
		>
			<span class="">+</span>
		</button>
	</div>
{/if}
