<script lang="ts">
	import Button from './Button.svelte';

	export let searchQuery: string = '';
	export let selectedCategory: string = 'All';
	export let priceRange: [number, number] = [0, 200];
	export let categories: string[] = ['All'];
	export let onSubmit: (filters: {
		searchQuery: string;
		selectedCategory: string;
		priceRange: [number, number];
	}) => void;

	// Local state to hold input values until submission
	let localSearchQuery = searchQuery;
	let localSelectedCategory = selectedCategory;
	let localPriceRange = priceRange;

	function handleSubmit() {
		onSubmit({
			searchQuery: localSearchQuery,
			selectedCategory: localSelectedCategory,
			priceRange: localPriceRange
		});
	}
</script>

<div class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]">
	<form on:submit|preventDefault={handleSubmit} class="flex flex-col gap-2">
		<div class="grid grid-cols-1 sm:grid-cols-4 gap-6">
			<!-- Search Bar -->
			<div>
				<label for="search" class="font-pixel text-retroGray mb-2 block">Search</label>
				<input
					id="search"
					type="text"
					bind:value={localSearchQuery}
					placeholder="Type to search..."
					class="w-full p-2 rounded-md bg-white border-2 border-retroBlack font-body text-retroGray focus:outline-none focus:border-retroCoral"
				/>
			</div>

			<!-- Category Filter -->
			<div>
				<label for="category" class="font-pixel text-retroGray mb-2 block">Category</label>
				<select
					id="category"
					bind:value={localSelectedCategory}
					class="w-full p-2 rounded-md bg-white border-2 border-retroBlack font-body text-retroGray focus:outline-none focus:border-retroCoral"
				>
					{#each categories as category}
						<option value={category}>{category}</option>
					{/each}
				</select>
			</div>

			<!-- Price Range Filter -->
			<div>
				<label class="font-pixel text-retroGray mb-2 block">Price Range</label>
				<div class="flex items-center gap-2">
					<input
						type="number"
						bind:value={localPriceRange[0]}
						min="0"
						max="200"
						step="10"
						class="w-20 p-2 rounded-md bg-white border-2 border-retroBlack font-body text-retroGray focus:outline-none focus:border-retroCoral"
					/>
					<span class="font-pixel text-retroGray">-</span>
					<input
						type="number"
						bind:value={localPriceRange[1]}
						min="0"
						max="200"
						step="10"
						class="w-20 p-2 rounded-md bg-white border-2 border-retroBlack font-body text-retroGray focus:outline-none focus:border-retroCoral"
					/>
				</div>
			</div>
		</div>
		<!-- Submit Button -->
		<div class="flex items-end">
			<Button
				variant="primary"
				subClass="w-full bg-retroCoral text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack"
				type="submit"
			>
				Search
			</Button>
		</div>
	</form>
</div>
