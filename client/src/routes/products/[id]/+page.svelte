<!-- /src/routes/products/[id]/+page.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';
    import Button from '$lib/components/Button.svelte';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import { GLOBAL } from '$lib';

    interface Product {
        id: number;
        name: string;
        price: number;
        description: string;
        category: { name: string };
        images: { url: string; isPrimary: boolean; order: number }[];
        specs: { key: string; value: string }[];
        related?: number[];
    }

    let product: Product | null = null;
    let relatedProducts: Product[] = [];
    let error: string = '';
    let currentImageIndex = 0;

    // Reactive fetch when productId changes
    $: productId = Number($page.params.id);
    $: fetchProduct(productId);

    async function fetchProduct(id: number) {
        try {
            const productResponse = await fetch(`${GLOBAL.SERVER_URL}/products/${id}`, {
                cache: 'no-store'             });
            if (!productResponse.ok) throw new Error('Failed to fetch product');
            const productData = await productResponse.json();
            product = productData.product;

            if (product && product.related && product.related.length > 0) {
                const relatedPromises = product.related.map((relatedId: number) =>
                    fetch(`${GLOBAL.SERVER_URL}/products/${relatedId}`, { cache: 'no-store' }).then(
                        (res) => {
                            if (!res.ok)
                                throw new Error(`Failed to fetch related product ${relatedId}`);
                            return res.json();
                        }
                    )
                );
                const relatedResponses = await Promise.all(relatedPromises);
                relatedProducts = relatedResponses.map((res: any) => res.product);
                console.log(
                    'Related Products:',
                    relatedProducts.map((p) => ({ id: p.id, name: p.name }))
                );
            } else {
                relatedProducts = [];
            }
        } catch (err) {
            error = `Error: ${err}`;
            console.error(err);
        }
    }

    // Initial fetch on mount (optional, since $: will handle it)
    onMount(() => fetchProduct(productId));

    $: currentImage = product?.images[currentImageIndex]?.url || '';
    function prevImage() {
        if (product && product.images.length > 1) {
            currentImageIndex =
                (currentImageIndex - 1 + product.images.length) % product.images.length;
        }
    }
    function nextImage() {
        if (product && product.images.length > 1) {
            currentImageIndex = (currentImageIndex + 1) % product.images.length;
        }
    }

    // Add to Cart (placeholder)
    function addToCart() {
        if (product) console.log(`Added ${product.name} to cart!`);
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    {#if error}
        <section class="text-center">
            <h1 class="font-pixel text-3xl text-retroGray mb-4">Error</h1>
            <p class="text-lg text-retroCoral font-body mb-6">{error}</p>
            <a
                href="/products"
                class="font-pixel text-retroGray hover:text-retroCoral transition-colors text-lg"
            >
                Back to Products
            </a>
        </section>
    {:else if product}
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
                        <p class="text-retroCoral font-pixel text-xl sm:text-2xl mb-4">
                            ${product.price}
                        </p>
                        <p class="text-base sm:text-lg text-retroGray font-body mb-4">
                            {product.description}
                        </p>
                        <p class="text-sm sm:text-base text-retroGray font-pixel mb-6">
                            <span class="text-retroTeal">Category:</span>
                            {product.category.name}
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
                    {#each product.specs as spec}
                        <tr class="border-b-2 border-retroBlack">
                            <th class="py-2 text-retroTeal capitalize">{spec.key}</th>
                            <td class="py-2">{spec.value}</td>
                        </tr>
                    {/each}
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
                    <p class="font-body text-retroGray">
                        "Love this {product.name}! Takes me back."
                    </p>
                    <p class="font-pixel text-sm text-retroTeal mt-2">— RetroFan88</p>
                </div>
                <div class="border-2 border-retroBlack p-4 rounded-md">
                    <p class="font-body text-retroGray">"Works great, but could use a polish."</p>
                    <p class="font-pixel text-sm text-retroTeal mt-2">— PixelPete</p>
                </div>
            </div>
        </section>

        <!-- Related Products -->
        {#if relatedProducts.length > 0}
            <section class="mb-12 sm:mb-16">
                <h2 class="font-pixel text-2xl text-retroGray mb-6 text-center">
                    Related Products
                </h2>
                <div
                    class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-6 sm:gap-8 justify-items-center"
                >
                    {#each relatedProducts as relatedProduct}
                        <ProductCard
                            product={{
                                id: relatedProduct.id,
                                name: relatedProduct.name,
                                price: relatedProduct.price,
                                image:
                                    relatedProduct.images.find((img) => img.isPrimary)?.url || '',
                                description: relatedProduct.description,
                                category: relatedProduct.category.name
                            }}
                            subClass="bg-white border-4"
                        />
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
            <h1 class="font-pixel text-3xl text-retroGray mb-4">Loading...</h1>
            <p class="text-lg text-retroGray font-body mb-6">Fetching your retro gem...</p>
        </section>
    {/if}
</main>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
