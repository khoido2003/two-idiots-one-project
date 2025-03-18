<!-- /src/routes/products/[id]/+page.svelte -->
<script lang="ts">
	import Button from '$lib/components/Button.svelte';
	import ProductCard from '$lib/components/ProductCard.svelte';
	import { page } from '$app/stores';

	// Mock product data with additional images and specs
	const allProducts = [
		{
			id: 1,
			name: 'NES Console',
			price: 59.99,
			images: ['/assets/nes.webp', '/assets/nes-front.jpg', '/assets/nes-controllers.jpg'],
			description:
				'Classic 8-bit gaming console with iconic controllers and a library of legendary games. Perfect for reliving the 80s.',
			category: 'Consoles',
			specs: { condition: 'Used', releaseYear: '1983', includes: 'Console, 2 controllers' },
			related: [3, 5] // Game Boy, Atari 2600
		},
		{
			id: 2,
			name: 'Polaroid SX-70',
			price: 129.99,
			images: ['/assets/polaroid.jpg', '/assets/polaroid-open.jpg', '/assets/polaroid-side.jpg'],
			description:
				'Vintage instant camera with a sleek foldable design. Capture moments in true retro style.',
			category: 'Cameras',
			specs: { condition: 'Refurbished', releaseYear: '1972', includes: 'Camera, manual' },
			related: []
		},
		{
			id: 3,
			name: 'Game Boy',
			price: 79.99,
			images: ['/assets/gameboy.jpg', '/assets/gameboy-back.jpg', '/assets/gameboy-screen.jpg'],
			description:
				'Portable retro gaming powerhouse. Tetris not included, but the nostalgia is free!',
			category: 'Consoles',
			specs: { condition: 'Used', releaseYear: '1989', includes: 'Console only' },
			related: [1, 5]
		}
		// Add more as needed...
	];

	// Get product ID from URL
	const productId = Number($page.params.id);
	const product = allProducts.find((p) => p.id === productId) || null;

	// Carousel state
	let currentImageIndex = 0;
	$: currentImage = product?.images[currentImageIndex] || '';
	function prevImage() {
		if (product)
			currentImageIndex = (currentImageIndex - 1 + product.images.length) % product.images.length;
	}
	function nextImage() {
		if (product) currentImageIndex = (currentImageIndex + 1) % product.images.length;
	}

	// Handle "Add to Cart"
	function addToCart() {
		console.log(`Added ${product?.name} to cart!`);
	}
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
	{#if product}
		<!-- Product Overview -->
		<section
			class="mb-12 sm:mb-16 bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
		>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 sm:gap-8">
				<!-- Carousel -->
				<div class="relative">
					<div
						class="w-full h-64 sm:h-80 md:h-96 mx-auto rounded-md border-4 border-retroBlack pixelated bg-retroGray"
						style={`background-image: url('${currentImage}'); background-size: cover; background-position: center;`}
					></div>
					<!-- Navigation Arrows -->
					{#if product.images.length > 1}
						<button
							on:click={prevImage}
							class="absolute left-2 top-1/2 transform -translate-y-1/2 bg-retroBlack text-retroCream w-10 h-10 rounded-full border-2 border-retroBlack hover:bg-retroCoral transition-all font-pixel text-lg"
						>
							←
						</button>
						<button
							on:click={nextImage}
							class="absolute right-2 top-1/2 transform -translate-y-1/2 bg-retroBlack text-retroCream w-10 h-10 rounded-full border-2 border-retroBlack hover:bg-retroCoral transition-all font-pixel text-lg"
						>
							→
						</button>
					{/if}
					<!-- Retro Overlay -->
					<div
						class="absolute inset-0 opacity-20 pointer-events-none bg-[url('/assets/crt-overlay.png')] bg-cover"
					></div>
				</div>

				<!-- Product Details -->
				<div class="flex flex-col justify-between">
					<div>
						<h1 class="font-pixel text-2xl sm:text-3xl md:text-4xl text-retroGray mb-4">
							{product.name}
						</h1>
						<p class="text-retroCoral font-pixel text-xl sm:text-2xl mb-4">${product.price}</p>
						<p class="text-base sm:text-lg text-retroGray font-body mb-4">{product.description}</p>
						<p class="text-sm sm:text-base text-retroGray font-pixel mb-6">
							<span class="text-retroTeal">Category:</span>
							{product.category}
						</p>
					</div>
					<Button
						variant="primary"
						type="button"
						subClass="bg-retroCoral text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3"
						on:click={addToCart}
					>
						Add to Cart
					</Button>
				</div>
			</div>
		</section>

		<!-- Specs Section -->
		<section
			class="mb-12 sm:mb-16 bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
		>
			<h2 class="font-pixel text-2xl text-retroGray mb-4">Specifications</h2>
			<table class="w-full text-left font-body text-retroGray">
				<tbody>
					<tr class="border-b-2 border-retroBlack">
						<th class="py-2 text-retroTeal">Condition</th>
						<td class="py-2">{product.specs.condition}</td>
					</tr>
					<tr class="border-b-2 border-retroBlack">
						<th class="py-2 text-retroTeal">Release Year</th>
						<td class="py-2">{product.specs.releaseYear}</td>
					</tr>
					<tr class="border-b-2 border-retroBlack">
						<th class="py-2 text-retroTeal">Includes</th>
						<td class="py-2">{product.specs.includes}</td>
					</tr>
				</tbody>
			</table>
		</section>

		<!-- Reviews Section (Mock) -->
		<section
			class="mb-12 sm:mb-16 bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
		>
			<h2 class="font-pixel text-2xl text-retroGray mb-4">Customer Reviews</h2>
			<div class="space-y-4">
				<div class="border-2 border-retroBlack p-4 rounded-md">
					<p class="font-body text-retroGray">"Love this NES! Takes me back to my childhood."</p>
					<p class="font-pixel text-sm text-retroTeal mt-2">— RetroFan88</p>
				</div>
				<div class="border-2 border-retroBlack p-4 rounded-md">
					<p class="font-body text-retroGray">"Works great, but one controller was sticky."</p>
					<p class="font-pixel text-sm text-retroTeal mt-2">— PixelPete</p>
				</div>
			</div>
		</section>

		<!-- Related Products -->
		{#if product.related.length > 0}
			<section class="mb-12 sm:mb-16">
				<h2 class="font-pixel text-2xl text-retroGray mb-6 text-center">Related Products</h2>
				<div
					class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 sm:gap-8 justify-items-center"
				>
					{#each product.related.map( (id) => allProducts.find((p) => p.id === id) ) as relatedProduct}
						{#if relatedProduct}
							<ProductCard product={relatedProduct} subClass="bg-white border-4" />
						{/if}
					{/each}
				</div>
			</section>
		{/if}

		<!-- Back to Products Link -->
		<section class="text-center">
			<a
				href="/products"
				class="font-pixel text-retroGray hover:text-retroCoral transition-colors text-lg"
			>
				← Back to All Products
			</a>
		</section>
	{:else}
		<section class="text-center">
			<h1 class="font-pixel text-3xl text-retroGray mb-4">Product Not Found</h1>
			<p class="text-lg text-retroGray font-body mb-6">Sorry, we couldn’t find that retro gem!</p>
			<a
				href="/products"
				class="font-pixel text-retroGray hover:text-retroCoral transition-colors text-lg"
			>
				Back to Products
			</a>
		</section>
	{/if}
</main>

<style>
	.pixelated {
		image-rendering: pixelated;
	}
</style>
