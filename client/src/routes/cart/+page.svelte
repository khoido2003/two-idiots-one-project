<!-- /src/routes/cart/+page.svelte -->
<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import { writable } from 'svelte/store';

    // Mock cart store (replace with persistent storage later)
    const cart = writable([
        { id: 1, name: 'NES Console', price: 59.99, image: '/assets/nes.webp', quantity: 1 },
        { id: 3, name: 'Game Boy', price: 79.99, image: '/assets/gameboy.jpg', quantity: 2 }
    ]);

    // Update quantity
    function updateQuantity(id: number, delta: number) {
        cart.update((items) => {
            const item = items.find((i) => i.id === id);
            if (item) {
                item.quantity = Math.max(1, item.quantity + delta); // Minimum 1
            }
            return items;
        });
    }

    // Remove item
    function removeItem(id: number) {
        cart.update((items) => items.filter((i) => i.id !== id));
    }

    // Calculate totals
    $: subtotal = $cart.reduce((sum, item) => sum + item.price * item.quantity, 0);
    $: tax = subtotal * 0.08; // Mock 8% tax
    $: shipping = $cart.length > 0 ? 10.0 : 0; // Mock flat shipping
    $: total = subtotal + tax + shipping;

    // Checkout placeholder
    function checkout() {
        console.log('Proceeding to checkout with total:', total.toFixed(2));
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <section class="mb-12 sm:mb-16">
        <h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-8 text-center">
            Your Cart
        </h1>

        {#if $cart.length > 0}
            <!-- Cart Items -->
            <div class="space-y-6 mb-8">
                {#each $cart as item (item.id)}
                    <div
                        class="bg-retroCream p-4 sm:p-6 rounded-xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] transition-all hover:shadow-[6px_6px_0_#2e2e2e]"
                    >
                        <div
                            class="grid grid-cols-1 sm:grid-cols-[auto_1fr_auto] gap-4 sm:gap-6 items-center"
                        >
                            <!-- Image -->
                            <div
                                class="w-40 h-40 sm:w-60 sm:h-48 rounded-md mb-3 sm:mb-4 pixelated"
                                style={`background-image: url('${item.image}'); background-size: cover; background-position: center;`}
                            ></div>

                            <!-- Name and Description -->
                            <div class="flex flex-col justify-center">
                                <h2
                                    class="font-pixel text-retroGray text-xl sm:text-2xl leading-tight"
                                >
                                    {item.name}
                                </h2>
                                <p class="font-body text-retroGray text-sm sm:text-base">
                                    Unit Price: ${item.price.toFixed(2)}
                                </p>
                            </div>

                            <!-- Controls and Price -->
                            <div class="flex flex-col xl:flex-row justify-center items-center">
                                <div class="flex items-center justify-center gap-6">
                                    <button
                                        on:click={() => updateQuantity(item.id, -1)}
                                        class="w-10 h-10 bg-retroTeal text-retroCream rounded-full border-2 border-retroBlack hover:bg-retroCoral hover:scale-105 transition-all font-pixel text-lg"
                                    >
                                        -
                                    </button>
                                    <span
                                        class="font-pixel text-retroGray text-lg min-w-[2rem] text-center font-bold"
                                        >{item.quantity}</span
                                    >
                                    <button
                                        on:click={() => updateQuantity(item.id, 1)}
                                        class="w-10 h-10 bg-retroTeal text-retroCream rounded-full border-2 border-retroBlack hover:bg-retroCoral hover:scale-105 transition-all font-pixel text-lg"
                                    >
                                        +
                                    </button>
                                </div>
                                <p class="font-pixel text-retroCoral text-xl sm:text-2xl">
                                    ${(item.price * item.quantity).toFixed(2)}
                                </p>
                                <button
                                    on:click={() => removeItem(item.id)}
                                    class="text-retroGray hover:text-retroCoral font-pixel text-sm sm:text-base transition-colors underline cursor-pointer"
                                >
                                    Remove
                                </button>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>

            <!-- Totals -->
            <div
                class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] max-w-lg mx-auto"
            >
                <div class="space-y-3">
                    <div class="flex justify-between font-pixel text-retroGray text-lg">
                        <span>Subtotal</span>
                        <span>${subtotal.toFixed(2)}</span>
                    </div>
                    <div class="flex justify-between font-pixel text-retroGray text-lg">
                        <span>Tax (8%)</span>
                        <span>${tax.toFixed(2)}</span>
                    </div>
                    <div class="flex justify-between font-pixel text-retroGray text-lg">
                        <span>Shipping</span>
                        <span>${shipping.toFixed(2)}</span>
                    </div>
                    <div
                        class="flex justify-between font-pixel text-retroGray text-xl border-t-2 border-retroBlack pt-3"
                    >
                        <span>Total</span>
                        <span class="text-retroCoral font-bold">${total.toFixed(2)}</span>
                    </div>
                </div>
                <Button
                    variant="primary"
                    type="button"
                    subClass="mt-6 bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 w-full"
                    on:click={checkout}
                >
                    Checkout
                </Button>
            </div>
        {:else}
            <div
                class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] text-center"
            >
                <p class="font-pixel text-retroGray text-lg">Your cart is empty!</p>
            </div>
        {/if}
    </section>

    <!-- Continue Shopping -->
    <section class="text-center">
        <a
            href="/products"
            class="font-pixel text-retroGray hover:text-retroCoral transition-colors text-lg"
        >
            ← Continue Shopping
        </a>
    </section>
</main>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
