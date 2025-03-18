<!-- /src/routes/products/+page.svelte -->
<script lang="ts">
	import Pagination from '$lib/components/Pagination.svelte';
	import SearchFilters from '$lib/components/SearchFilters.svelte';
	import ProductCard from '$lib/components/ProductCard.svelte';

	// Mock product data with categories
	const allProducts = [
		{
			id: 1,
			name: 'NES Console',
			price: 59.99,
			image: '/assets/nes.webp',
			description: 'Classic 8-bit gaming console.',
			category: 'Consoles'
		},
		{
			id: 2,
			name: 'Polaroid SX-70',
			price: 129.99,
			image: '/assets/polaroid.jpg',
			description: 'Vintage instant camera.',
			category: 'Cameras'
		},
		{
			id: 3,
			name: 'Game Boy',
			price: 79.99,
			image: '/assets/gameboy.jpg',
			description: 'Portable retro gaming.',
			category: 'Consoles'
		},
		{
			id: 4,
			name: 'VHS Player',
			price: 49.99,
			image: '/assets/vhs.jpg',
			description: 'Rewind to the 80s.',
			category: 'Media Players'
		},
		{
			id: 5,
			name: 'Atari 2600',
			price: 39.99,
			image: '/assets/atari.webp',
			description: 'Classic arcade vibes.',
			category: 'Consoles'
		},
		{
			id: 6,
			name: 'SNES Controller',
			price: 29.99,
			image: '/assets/snes-controller.jpg',
			description: 'Super gaming control.',
			category: 'Accessories'
		},
		{
			id: 7,
			name: 'Walkman',
			price: 69.99,
			image: '/assets/walkman.jpg',
			description: 'Portable tunes, 80s style.',
			category: 'Media Players'
		},
		{
			id: 8,
			name: 'Floppy Disk Set',
			price: 19.99,
			image: '/assets/floppy.jpg',
			description: 'Retro data storage.',
			category: 'Accessories'
		},
		{
			id: 9,
			name: 'Tamagotchi',
			price: 24.99,
			image: '/assets/tamagotchi.jpg',
			description: 'Virtual pet nostalgia.',
			category: 'Toys'
		},
		{
			id: 10,
			name: 'CRT Monitor',
			price: 89.99,
			image: '/assets/crt.jpg',
			description: 'Chunky retro display.',
			category: 'Displays'
		}
	];

	// Filter state
	let searchQuery = '';
	let selectedCategory = 'All';
	let priceRange: [number, number] = [0, 200];
	const categories = [
		'All',
		'Consoles',
		'Cameras',
		'Media Players',
		'Accessories',
		'Toys',
		'Displays'
	];

	// Filtered products (updated only on submit)
	let filteredProducts = allProducts;

	function handleFilterSubmit(filters: {
		searchQuery: string;
		selectedCategory: string;
		priceRange: [number, number];
	}) {
		searchQuery = filters.searchQuery;
		selectedCategory = filters.selectedCategory;
		priceRange = filters.priceRange;

		filteredProducts = allProducts.filter((product) => {
			const matchesSearch =
				product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
				product.description.toLowerCase().includes(searchQuery.toLowerCase());
			const matchesCategory = selectedCategory === 'All' || product.category === selectedCategory;
			const matchesPrice = product.price >= priceRange[0] && product.price <= priceRange[1];
			return matchesSearch && matchesCategory && matchesPrice;
		});

		currentPage = 1; // Reset to first page on filter submit
	}

	// Pagination state
	let currentPage = 1;
	const itemsPerPage = 8;
	$: totalPages = Math.ceil(filteredProducts.length / itemsPerPage);
	$: displayedProducts = filteredProducts.slice(
		(currentPage - 1) * itemsPerPage,
		currentPage * itemsPerPage
	);

	// Scroll to keep pagination at the same bottom position
	function handlePageChange(page: number, event?: Event) {
		if (event) event.preventDefault();

		// Get the current bottom position relative to the document
		const viewportHeight = window.innerHeight;
		const currentBottom = window.scrollY + viewportHeight;

		// Change the page
		currentPage = page;

		// After render, adjust scroll to keep the bottom aligned
		requestAnimationFrame(() => {
			const newDocumentHeight = document.documentElement.scrollHeight;
			const newScrollPosition = newDocumentHeight - viewportHeight;

			// Only scroll if the new position keeps the bottom in view
			if (newScrollPosition >= 0 && currentBottom > window.scrollY) {
				window.scrollTo(0, newScrollPosition);
			}
		});
	}
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
	<!-- Page Header -->
	<section class="mb-12 sm:mb-16 text-center">
		<h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-4">All Retro Gear</h1>
		<p class="text-base sm:text-lg text-retroGray font-body max-w-2xl mx-auto">
			Browse our full collection of vintage tech treasures!
		</p>
	</section>

	<!-- Search and Filters -->
	<section class="mb-12 sm:mb-16">
		<SearchFilters
			{searchQuery}
			{selectedCategory}
			{priceRange}
			{categories}
			onSubmit={handleFilterSubmit}
		/>
	</section>

	<!-- Products Grid -->
	<section class="mb-12 sm:mb-16">
		{#if displayedProducts.length > 0}
			<div
				class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8 justify-items-center px-4"
			>
				{#each displayedProducts as product}
					<div class="group">
						<ProductCard
							{product}
							subClass="transform group-hover:scale-105 transition-transform duration-300 bg-white border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
						/>
					</div>
				{/each}
			</div>
		{:else}
			<p class="text-center font-pixel text-retroGray text-lg">
				No products found. Try adjusting your filters!
			</p>
		{/if}
	</section>

	<!-- Pagination -->
	<section class="mb-12 sm:mb-16">
		<Pagination {currentPage} {totalPages} onPageChange={handlePageChange} />
	</section>
</main>

<style>
	.pixelated {
		image-rendering: pixelated;
	}
</style>
