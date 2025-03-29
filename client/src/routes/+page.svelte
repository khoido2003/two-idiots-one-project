<!-- /src/routes/+page.svelte -->
<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import 'iconify-icon';

    export let data: {
        featuredProducts: any[];
        newArrivals: any[];
        limitedDeals: any[];
    };

    // Testimonials remain static for now
    const testimonials = [
        { name: 'Alex G.', quote: 'Found my dream retro console here—works like a charm!' },
        { name: 'Sam R.', quote: 'The cameras are pristine, and shipping was fast!' }
    ];
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <!-- Hero Section -->
    <section
        class="relative text-center py-20 sm:py-20 text-retroCream mb-8 sm:mb-12 overflow-hidden h-96 md:h-[520px]"
    >
        <div
            class="absolute inset-0 bg-[url('/assets/retro-bg.png')] bg-cover bg-center opacity-80"
        ></div>
        <div class="absolute inset-0 bg-pink-950 opacity-70"></div>
        <div
            class="absolute inset-0 opacity-50 pointer-events-none bg-[url('/assets/retro-car.gif')] bg-no-repeat bg-cover bg-center"
        ></div>
        <div class="relative max-w-4xl mx-auto z-10">
            <h1
                class="font-pixel text-4xl sm:text-4xl md:text-5xl mb-4 animate-pulse tracking-wide text-white"
            >
                RetroTech Marketplace
            </h1>
            <p
                class="text-base sm:text-lg md:text-xl max-w-2xl mx-auto mb-6 font-pixel text-retroCoral"
            >
                Unlock the past—score epic vintage tech today!
            </p>
            <Button variant="primary" subClass="hover:scale-105 transition-transform">
                <a href="/products">Grab Your Gear Now</a>
            </Button>
        </div>
    </section>

    <!-- Featured Products Section -->
    <section class="mb-12 sm:mb-16 relative">
        <h2 class="font-pixel text-2xl sm:text-3xl text-retroGray mb-6 text-center">
            Hot Retro Finds
        </h2>
        <div
            class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8 justify-items-center"
        >
            {#each data.featuredProducts as product}
                <div class="group">
                    <ProductCard
                        product={{
                            id: product.id,
                            name: product.name,
                            price: product.price,
                            image: product.images.find((img) => img.isPrimary)?.url || '',
                            description: product.description,
                            category: product.category.name
                        }}
                        subClass="transform group-hover:scale-105 transition-transform duration-300"
                    />
                </div>
            {/each}
        </div>
        <div class="text-center mt-6">
            <Button
                variant="secondary"
                subClass="hover:bg-retroBlue hover:text-retroGray transition-colors"
            >
                <a href="/products">More Retro Loot</a>
            </Button>
        </div>
    </section>

    <div class="bg-retroBlack h-1 w-full"></div>

    <!-- New Arrivals Section -->
    <section class="mb-12 sm:mb-16 py-12 bg-retroCream rounded-2xl">
        <h2 class="font-pixel text-2xl sm:text-3xl text-retroGray mb-6 text-center">
            New Arrivals
        </h2>
        <p class="text-center text-sm sm:text-base font-pixel text-retroGray mb-6">
            Fresh retro gear just dropped!
        </p>
        <div
            class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 px-4 sm:gap-8 justify-items-center"
        >
            {#each data.newArrivals as product}
                <div class="group">
                    <ProductCard
                        product={{
                            id: product.id,
                            name: product.name,
                            price: product.price,
                            image: product.images.find((img) => img.isPrimary)?.url || '',
                            description: product.description,
                            category: product.category.name
                        }}
                        subClass="transform group-hover:scale-105 transition-transform duration-300 bg-white border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
                    />
                </div>
            {/each}
        </div>
        <div class="text-center mt-6">
            <Button
                variant="secondary"
                subClass="hover:bg-retroBlue hover:text-retroGray transition-colors"
            >
                <a href="/products">Shop New Stock</a>
            </Button>
        </div>
    </section>

    <div class="bg-retroBlack h-1 w-full"></div>

    <!-- Limited Deals Section -->
    <section class="mt-10 mb-12 sm:mb-16 py-12 bg-stone-600 rounded-2xl text-retroCream">
        <h2 class="font-pixel text-2xl sm:text-3xl text-center mb-4 text-retroPeach">
            Limited Deals
        </h2>
        <p class="text-center text-sm sm:text-base font-pixel text-retroCream mb-6">
            Snag these before they vanish!
        </p>
        <div class="max-w-2xl mx-auto p-6">
            {#each data.limitedDeals as deal}
                <div
                    class="bg-white p-6 rounded-lg border-4 border-retroPlum shadow-[4px_4px_0_#ff6f61]"
                >
                    <img
                        src={deal.images.find((img) => img.isPrimary)?.url || ''}
                        alt={deal.name}
                        class="w-full h-40 object-cover rounded-md mb-4 pixelated"
                    />
                    <h3 class="font-pixel text-lg sm:text-xl text-retroGray mb-2">{deal.name}</h3>
                    <p class="text-sm font-pixel text-retroGray">{deal.description}</p>
                    <div class="flex justify-between items-center mt-4">
                        <div>
                            <span class="text-retroCoral font-pixel text-lg">${deal.price}</span>
                            <span class="text-retroGray line-through ml-2 opacity-70"
                                >${deal.originalPrice}</span
                            >
                        </div>
                        <Button
                            variant="primary"
                            subClass="text-retroCream hover:bg-retroCoral hover:text-retroCream transition-all duration-300 font-pixel"
                        >
                            <a href={`/products/${deal.id}`}>Buy Now</a>
                        </Button>
                    </div>
                </div>
            {/each}
        </div>
    </section>

    <!-- Why Shop With Us Section -->
    <section class="bg-retroBlack py-16 px-3 sm:py-16 mb-12 sm:mb-16 rounded-2xl">
        <h2 class="font-pixel text-2xl sm:text-3xl text-center mb-10 text-retroCream">
            Why RetroTech?
        </h2>
        <div class="container mx-auto px-4 grid grid-cols-1 sm:grid-cols-3 gap-6 sm:gap-8">
            <div
                class="text-center p-4 bg-retroGray text-retroCream rounded-md shadow-[4px_4px_0_#ff6f61]"
            >
                <iconify-icon class="h-7 w-7" icon="mdi:shield-crown" width="28" height="28"
                ></iconify-icon>
                <h3 class="font-pixel text-lg sm:text-xl mb-2 text-retroCoral">Authentic Gear</h3>
                <p class="text-sm sm:text-base font-pixel">Real retro, guaranteed.</p>
            </div>
            <div
                class="text-center p-4 bg-retroGray text-retroCream rounded-md shadow-[4px_4px_0_#ff6f61]"
            >
                <iconify-icon class="h-7 w-7" icon="mdi:truck-fast" width="28" height="28"
                ></iconify-icon>
                <h3 class="font-pixel text-lg sm:text-xl mb-2 text-retroCoral">Fast Shipping</h3>
                <p class="text-sm sm:text-base font-pixel">Quick to your door.</p>
            </div>
            <div
                class="text-center p-4 bg-retroGray text-retroCream rounded-md shadow-[4px_4px_0_#ff6f61]"
            >
                <iconify-icon class="h-7 w-7" icon="mdi:robot" width="28" height="28"
                ></iconify-icon>
                <h3 class="font-pixel text-lg sm:text-xl mb-2 text-retroCoral">AI Assist</h3>
                <p class="text-sm sm:text-base font-pixel">24/7 Tech help, 8-bit style.</p>
            </div>
        </div>
    </section>

    <!-- Testimonials Section -->
    <section class="mb-12 sm:mb-16">
        <h2 class="font-pixel text-2xl sm:text-3xl text-retroGray mb-6 text-center">
            Retro Fans Speak
        </h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 gap-6">
            {#each testimonials as testimonial}
                <div
                    class="relative bg-retroCream p-4 sm:p-6 rounded-lg shadow-md border-4 border-retroGray"
                >
                    <div class="absolute -top-4 left-4 w-8 h-8 bg-retroCoral rounded-full"></div>
                    <p class="text-sm sm:text-base text-retroGray font-pixel mb-4">
                        "{testimonial.quote}"
                    </p>
                    <p class="font-pixel text-retroCoral text-sm sm:text-base">
                        - {testimonial.name}
                    </p>
                </div>
            {/each}
        </div>
    </section>
</main>

<footer class="bg-retroGray text-retroCream py-6 sm:py-8 relative">
    <div class="absolute inset-0 opacity-10 bg-[url('/assets/pixel-border.png')] bg-repeat-x"></div>
    <div class="container mx-auto px-4 text-center relative z-10">
        <p class="font-pixel text-sm sm:text-base mb-2">
            RetroTech Marketplace © {new Date().getFullYear()}
        </p>
        <div class="flex justify-center space-x-4 sm:space-x-6">
            <a href="/about" class="hover:text-retroCoral transition-colors text-sm sm:text-base"
                >About</a
            >
            <a href="/contact" class="hover:text-retroCoral transition-colors text-sm sm:text-base"
                >Contact</a
            >
            <a href="/terms" class="hover:text-retroCoral transition-colors text-sm sm:text-base"
                >Terms</a
            >
        </div>
    </div>
</footer>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
    @keyframes bounce {
        0%,
        100% {
            transform: translateY(0);
        }
        50% {
            transform: translateY(-8px);
        }
    }
    .animate-bounce {
        animation: bounce 2s infinite;
    }
</style>
