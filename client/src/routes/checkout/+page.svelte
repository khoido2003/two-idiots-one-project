<!-- /src/routes/checkout/+page.svelte -->
<script lang="ts">
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
    let showToast = false;
    let toastMessage = '';
    let toastType: 'error' | 'success' | 'info' = 'error';

    // Form data
    let address = {
        line1: '',
        line2: '',
        city: '',
        state: '',
        postalCode: '',
        country: 'US' // Default to US, could be a dropdown
    };

    // Stripe
    let stripe: Stripe | null = null;
    let elements: any; // You could define a more specific type if needed
    let clientSecret: string | null = null;
    let paymentProcessing = false;

    if (!$auth.token || !$auth.user) {
        cartItems = [];
        window.location.href = '/signin';
    }

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
        } catch (error) {
            console.error('Error fetching cart:', error);
            cartItems = [];
            showToast = true;
            toastMessage = 'Error fetching cart! Try again later';
            toastType = 'error';
        } finally {
            loading = false;
        }
    }

    onMount(async () => {
        await fetchCart();
        if (cartItems.length === 0) {
            window.location.href = '/cart';
            return;
        }

        // Load Stripe.js
        const stripeScript = document.createElement('script');
        stripeScript.src = 'https://js.stripe.com/v3/';
        stripeScript.onload = async () => {
            stripe = window.Stripe(
                'pk_test_51OfhClFk7bi00dHV71gF7kAAULWZPEtObnayx24VXzmDSYJCVKoZ1qQB9GtWe20QX3MVmsHyMDdw9euwCwMqco7l00dokhaMZO'
            );
            elements = stripe.elements();
            const card = elements.create('card', { style: { base: { fontSize: '16px' } } });
            card.mount('#card-element');
        };
        document.head.appendChild(stripeScript);
    });

    $: subtotal = cartItems.reduce(
        (sum: number, item: CartItem) => sum + item.product.price * item.quantity,
        0
    );
    $: tax = subtotal * 0.08;
    $: shipping = cartItems.length > 0 ? 10.0 : 0;
    $: total = subtotal + tax + shipping;

    async function handleCheckout() {
        if (!stripe || !elements || paymentProcessing) return;

        // Validate address
        if (
            !address.line1 ||
            !address.city ||
            !address.state ||
            !address.postalCode ||
            !address.country
        ) {
            showToast = true;
            toastMessage = 'Please fill out all required address fields';
            toastType = 'error';
            return;
        }

        paymentProcessing = true;

        // Create payment intent
        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/checkout?userId=${$auth.user?.id}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${$auth.token}`
                },
                body: JSON.stringify({
                    cartItems: cartItems.map((item) => ({
                        productId: item.product.id,
                        quantity: item.quantity
                    })),
                    address
                })
            });
            if (!response.ok) throw new Error('Failed to create payment intent');
            const data = await response.json();
            clientSecret = data.clientSecret;

            // Confirm payment
            const result = await stripe.confirmCardPayment(clientSecret, {
                payment_method: {
                    card: elements.getElement('card'),
                    billing_details: {
                        name: `${$auth.user?.firstName} ${$auth.user?.lastName}`,
                        address: {
                            line1: address.line1,
                            line2: address.line2,
                            city: address.city,
                            state: address.state,
                            postal_code: address.postalCode,
                            country: address.country
                        }
                    }
                }
            });

            if (result.error) {
                showToast = true;
                toastMessage = result.error.message;
                toastType = 'error';
            } else if (result.paymentIntent.status === 'succeeded') {
                showToast = true;
                toastMessage = 'Payment successful! Your order is being processed.';
                toastType = 'success';
                // Optionally clear cart or redirect
                setTimeout(() => (window.location.href = '/orders'), 2000);
            }
        } catch (error) {
            showToast = true;
            toastMessage = 'Checkout failed. Please try again.';
            toastType = 'error';
        } finally {
            paymentProcessing = false;
        }
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <Toast bind:visible={showToast} message={toastMessage} type={toastType} duration={3000} />

    <section class="mb-12 sm:mb-16">
        <h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-8 text-center">
            Checkout
        </h1>

        {#if loading}
            <p class="font-pixel text-retroGray text-lg text-center">Loading...</p>
        {:else if cartItems.length > 0}
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                <!-- Shipping Address Form -->
                <div
                    class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
                >
                    <h2 class="font-pixel text-2xl text-retroGray mb-4">Shipping Address</h2>
                    <form>
                        <div class="mb-4">
                            <label class="font-pixel text-retroGray" for="line1"
                                >Address Line 1 *</label
                            >
                            <input
                                id="line1"
                                type="text"
                                bind:value={address.line1}
                                class="w-full p-2 border-2 border-retroBlack rounded-md font-body"
                                required
                            />
                        </div>
                        <div class="mb-4">
                            <label class="font-pixel text-retroGray" for="line2"
                                >Address Line 2</label
                            >
                            <input
                                id="line2"
                                type="text"
                                bind:value={address.line2}
                                class="w-full p-2 border-2 border-retroBlack rounded-md font-body"
                            />
                        </div>
                        <div class="mb-4">
                            <label class="font-pixel text-retroGray" for="city">City *</label>
                            <input
                                id="city"
                                type="text"
                                bind:value={address.city}
                                class="w-full p-2 border-2 border-retroBlack rounded-md font-body"
                                required
                            />
                        </div>
                        <div class="mb-4">
                            <label class="font-pixel text-retroGray" for="state">State *</label>
                            <input
                                id="state"
                                type="text"
                                bind:value={address.state}
                                class="w-full p-2 border-2 border-retroBlack rounded-md font-body"
                                required
                            />
                        </div>
                        <div class="mb-4">
                            <label class="font-pixel text-retroGray" for="postalCode"
                                >Postal Code *</label
                            >
                            <input
                                id="postalCode"
                                type="text"
                                bind:value={address.postalCode}
                                class="w-full p-2 border-2 border-retroBlack rounded-md font-body"
                                required
                            />
                        </div>
                        <div class="mb-4">
                            <label class="font-pixel text-retroGray" for="country">Country *</label>
                            <input
                                id="country"
                                type="text"
                                bind:value={address.country}
                                class="w-full p-2 border-2 border-retroBlack rounded-md font-body"
                                required
                            />
                        </div>
                    </form>
                </div>

                <!-- Order Summary and Payment -->
                <div
                    class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
                >
                    <h2 class="font-pixel text-2xl text-retroGray mb-4">Order Summary</h2>
                    {#each cartItems as item}
                        <div class="flex justify-between mb-2">
                            <span class="font-pixel text-retroGray"
                                >{item.product.name} (x{item.quantity})</span
                            >
                            <span class="font-pixel text-retroCoral"
                                >${(item.product.price * item.quantity).toFixed(2)}</span
                            >
                        </div>
                    {/each}
                    <div class="flex justify-between font-pixel text-retroGray text-lg mt-4">
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
                        class="flex justify-between font-pixel text-retroGray text-xl border-t-2 border-retroBlack pt-3 mt-4"
                    >
                        <span>Total</span>
                        <span class="text-retroCoral font-bold">${total.toFixed(2)}</span>
                    </div>

                    <h2 class="font-pixel text-2xl text-retroGray mt-6 mb-4">Payment Details</h2>
                    <div
                        id="card-element"
                        class="p-2 border-2 border-retroBlack rounded-md mb-4"
                    ></div>
                    <Button
                        variant="primary"
                        subClass="bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 w-full"
                        onClick={handleCheckout}
                        disabled={paymentProcessing}
                    >
                        {paymentProcessing ? 'Processing...' : 'Pay Now'}
                    </Button>
                </div>
            </div>
        {:else}
            <div
                class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] text-center"
            >
                <p class="font-pixel text-retroGray text-lg">Your cart is empty!</p>
                <a
                    href="/products"
                    class="font-pixel text-retroGray hover:text-retroCoral mt-4 inline-block"
                    >Continue Shopping</a
                >
            </div>
        {/if}
    </section>
</main>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
