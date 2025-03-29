<!-- /src/lib/components/ProductCard.svelte -->
<script lang="ts">
	import Button from './Button.svelte';

	export let product = {
		id: 1,
		name: 'NES Console',
		price: 59.99,
		image: 'https://via.placeholder.com/150',
		description: 'Classic 8-bit gaming console.',
        category: "Electron"
	};

	export let subClass = '';
</script>

<a
    href={`/products/${product.id}`}
    class={`block bg-white w-72 p-3 sm:p-4 rounded-lg shadow-md hover:shadow-xl transition-shadow ${subClass}`}
>
    <div
        class="w-40 h-40 sm:w-60 sm:h-48 mx-auto rounded-md mb-3 sm:mb-4 pixelated"
        style={`background-image: url('${product.image}'); background-size: cover; background-position: center;`}
    ></div>
    <h3 class="font-pixel text-retroGray text-base sm:text-lg truncate">{product.name}</h3>
    <p class="text-xs sm:text-sm text-gray-600 line-clamp-2">{product.description}</p>
    <p class="text-retroCoral font-bold mt-2 text-sm sm:text-base">${product.price}</p>
    <Button
        variant="primary"
        subClass="mt-2"
        onClick={(e) => {
            e.preventDefault(); // Prevent navigation when clicking the button
            console.log(`Add ${product.name} to cart`);
        }}
    >
        Add to Cart
    </Button>
</a>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
