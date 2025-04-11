<!-- /src/routes/cart/+page.svelte -->
<script lang="ts">
    import { goto } from '$app/navigation';
    import { GLOBAL } from '$lib';
    import { auth } from '$lib/auth';
    import Button from '$lib/components/Button.svelte';
    import Toast from '$lib/components/Toast.svelte';
    import { onMount } from 'svelte';

    interface CartItem {
        id: number;
        quantity: number;
        product: {
            price: number;
            name: string;
            images: { url: string }[];
        };
    }

    let cartItems: CartItem[] = [];
    let loading = true;

    if (!$auth.token || !$auth.user) {
        cartItems = [];
    }

    let showToast = false;
    let toastMessage = '';
    let toastType: 'error' | 'success' | 'info' = 'error';
    let originalItems: CartItem[] = [];
    $: isUpdating = false;

    async function fetchCart() {
        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/cart?userId=${$auth.user?.id}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${$auth.token}`
                }
            });
            if (!response.ok) throw new Error('Failed to fetch cart');
            const data = await response.json();
            cartItems = data.cartItems || [];
            originalItems = data.cartItems || [];
        } catch (error) {
            console.error('Error fetching cart:', error);
            cartItems = [];
            showToast = true;
            toastMessage = 'Error fetching cart! Try again later';
            toastType = 'error';
        } finally {
            loading = false;
            if (cartItems.length > 0) {
                showToast = true;
                toastMessage = 'Fetched cart successfully';
                toastType = 'success';
            }
        }
    }

    onMount(fetchCart);
    // Add a simple debounce utility
    // Add this to track button state
    function logButtonState() {
        console.log('[Button State] isUpdating:', isUpdating);
    }

    async function updateQuantity(id: number, delta: number) {
        if (!$auth.token) {
            console.log('[updateQuantity] No auth token, exiting');
            return;
        }

        const itemIndex = cartItems.findIndex((i) => i.id === id);
        if (itemIndex === -1) {
            console.log('[updateQuantity] Item not found, exiting');
            return;
        }

        const currentItem = cartItems[itemIndex];
        const newQuantity = Math.max(1, currentItem.quantity + delta);

        console.log(
            `[updateQuantity] Start - ID: ${id}, Current: ${currentItem.quantity}, New: ${newQuantity}`
        );

        cartItems = [
            ...cartItems.slice(0, itemIndex),
            { ...currentItem, quantity: newQuantity },
            ...cartItems.slice(itemIndex + 1)
        ];
        console.log(`[updateQuantity] Optimistic - UI Quantity: ${cartItems[itemIndex].quantity}`);

        try {
            console.log('[updateQuantity] Entering try block');
            isUpdating = true;
            console.log('[updateQuantity] isUpdating set to true');

            const response = await fetch(
                `${GLOBAL.SERVER_URL}/cart/${id}?userId=${$auth.user?.id}`,
                {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${$auth.token}`
                    },
                    body: JSON.stringify({ quantity: newQuantity })
                }
            );
            console.log('[updateQuantity] Fetch completed, status:', response.status);

            if (!response.ok) {
                const errorText = await response.text();
                console.log(
                    '[updateQuantity] Fetch failed, status:',
                    response.status,
                    'Error:',
                    errorText
                );
                throw new Error('Failed to update quantity');
            }

            const data = await response.json();
            console.log('[updateQuantity] Server response:', data);

            cartItems = [
                ...cartItems.slice(0, itemIndex),
                data.cartItem,
                ...cartItems.slice(itemIndex + 1)
            ];
            console.log(
                `[updateQuantity] Server update - UI Quantity: ${cartItems[itemIndex].quantity}`
            );

            showToast = true;
            toastMessage = 'Item updated successfully';
            toastType = 'success';
        } catch (error) {
            console.error('[updateQuantity] Error:', error);
            await fetchCart();
            showToast = true;
            toastMessage = 'Failed to update quantity';
            toastType = 'error';
        } finally {
            console.log('[updateQuantity] Finally block - Setting isUpdating to false');
            isUpdating = false;
            logButtonState(); // Log after reset
        }
    }
    async function removeItem(id: number) {
        if (!$auth.token) return;

        const originalItems = [...cartItems]; // Store original state

        // Optimistic update
        cartItems = cartItems.filter((i) => i.id !== id);

        try {
            isUpdating = true;
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/cart/${id}?userId=${$auth.user?.id}`,
                {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${$auth.token}`
                    }
                }
            );
            if (!response.ok) throw new Error('Failed to remove item');
            showToast = true;
            toastMessage = 'Item removed successfully';
            toastType = 'success';
        } catch (error) {
            console.error('Error removing item:', error);
            // Rollback on error
            cartItems = originalItems;
            showToast = true;
            toastMessage = 'Failed to remove item';
            toastType = 'error';
        } finally {
            isUpdating = false;
        }
    }

    $: subtotal = cartItems.reduce(
        (sum: number, item: CartItem) => sum + item.product.price * item.quantity,
        0
    );
    $: tax = subtotal * 0.08;
    $: shipping = cartItems.length > 0 ? 10.0 : 0;
    $: total = subtotal + tax + shipping;

    function checkout() {
        console.log('Proceeding to checkout with total:', total.toFixed(2));
        goto('/checkout');
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <Toast bind:visible={showToast} message={toastMessage} type={toastType} duration={3000} />

    <section class="mb-12 sm:mb-16">
        <h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-8 text-center">
            Your Cart
        </h1>

        {#if loading}
            <p class="font-pixel text-retroGray text-lg text-center">Loading...</p>
        {:else if cartItems.length > 0}
            <div class="space-y-6 mb-8">
                {#each cartItems as item (item.id)}
                    <div
                        class="bg-retroCream p-4 sm:p-6 rounded-xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] transition-all hover:shadow-[6px_6px_0_#2e2e2e]"
                    >
                        <div
                            class="grid grid-cols-1 md:grid-cols-[auto_1fr_auto] gap-4 sm:gap-6 items-center justify-center"
                        >
                            <div
                                class="w-40 h-40 sm:w-60 sm:h-48 rounded-md mb-3 sm:mb-4 pixelated"
                                style={`background-image: url('${item.product.images[0]?.url || '/placeholder.jpg'}'); background-size: cover; background-position: center;`}
                            ></div>
                            <div class="flex flex-col justify-center">
                                <h2
                                    class="font-pixel text-retroGray text-xl sm:text-2xl leading-tight"
                                >
                                    {item.product.name}
                                </h2>
                                <p class="font-body text-retroGray text-sm sm:text-base">
                                    Unit Price: ${item.product.price.toFixed(2)}
                                </p>
                            </div>
                            <div
                                class="flex flex-col xl:flex-row justify-center items-center gap-6"
                            >
                                <div class="flex items-center justify-center gap-6">
                                    <button
                                        on:click={() => updateQuantity(item.id, -1)}
                                        class={`w-10 h-10 bg-retroTeal text-retroCream rounded-full border-2 border-retroBlack hover:bg-retroCoral hover:scale-105 transition-all font-pixel text-lg cursor-pointer ${isUpdating && 'disabled:cursor-not-allowed'}`}
                                    >
                                        -
                                    </button>
                                    <span
                                        class="font-pixel text-retroGray text-lg min-w-[2rem] text-center font-bold"
                                    >
                                        {item.quantity}
                                    </span>
                                    <button
                                        on:click={() => updateQuantity(item.id, 1)}
                                        class={`w-10 h-10 bg-retroTeal text-retroCream rounded-full border-2 border-retroBlack hover:bg-retroCoral hover:scale-105 transition-all font-pixel text-lg cursor-pointer ${isUpdating && 'disabled:cursor-not-allowed'}`}
                                    >
                                        +
                                    </button>
                                </div>
                                <p class="font-pixel text-retroCoral text-xl sm:text-2xl">
                                    ${(item.product.price * item.quantity).toFixed(2)}
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
                    onClick={checkout}
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
