<script lang="ts">
    import { goto } from '$app/navigation';
    import { GLOBAL } from '$lib';
    import { auth } from '$lib/auth';
    import Toast from '$lib/components/Toast.svelte';
    import { onMount } from 'svelte';
    export const prerender = false;
    export const ssr = false;

    interface OrderItem {
        productId: number;
        quantity: number;
        price: number;
        name: string;
    }

    interface Order {
        id: number;
        userId: number;
        createdAt: string;
        total: number;
        items: OrderItem[];
        line1: string;
        line2?: string | null;
        city: string;
        state: string;
        postalCode: string;
        country: string;
    }
    let orders: Order[] = [];
    let loading = true;
    let showToast = false;
    let toastMessage = '';
    let toastType: 'error' | 'success' | 'info' = 'error';

    // Redirect to signin if not authenticated
    if (!$auth.token || !$auth.user) {
        orders = [];
        goto('/signin');
    }

    async function fetchOrders() {
        try {
            console.log('Fetching orders...');
            loading = true;
            const response = await fetch(`${GLOBAL.SERVER_URL}/orders?userId=${$auth.user?.id}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${$auth.token}`
                }
            });
            if (!response.ok) {
                if (response.status === 401) {
                    auth.set({ token: null, user: null });
                    window.location.href = '/signin';
                    return;
                }
                throw new Error('Failed to fetch orders');
            }
            const data = await response.json();
            console.log('Orders fetched:', data);
            orders = data.orders || [];
            console.log(orders);
            loading = false;
        } catch (error) {
            console.error('Error fetching orders:', error);
            orders = [];
            showToast = true;
            toastMessage = 'Error fetching orders! Try again later';
            toastType = 'error';
        } finally {
            console.log('Finally: setting loading to false');
            loading = false;
        }
    }

    onMount(async () => {
        await fetchOrders();
        if (orders.length === 0 && !showToast) {
            showToast = true;
            toastMessage = 'No orders found.';
            toastType = 'info';
        } else {
            showToast = true;

            toastMessage = 'Fetch data successfuly';
            toastType = 'info';
        }
    });

    function formatDate(dateString: string): string {
        return new Date(dateString).toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <Toast bind:visible={showToast} message={toastMessage} type={toastType} duration={3000} />

    {#if !$auth.token || !$auth.user}
        <div class="flex justify-center items-center h-screen">
            <p class="font-pixel text-retroGray">Redirecting...</p>
        </div>
    {:else}
        <section class="mb-12 sm:mb-16">
            <h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-8 text-center">
                Your Orders
            </h1>

            {#if loading}
                <div class="flex justify-center items-center">
                    <div class="font-pixel text-retroGray text-lg">Loading your orders...</div>
                    <div
                        class="ml-4 animate-spin h-6 w-6 border-4 border-retroCoral border-t-transparent rounded-full"
                    ></div>
                </div>
            {:else if orders.length > 0}
                <div class="space-y-8">
                    {#each orders as order}
                        <div
                            class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
                        >
                            <h2 class="font-pixel text-2xl text-retroGray mb-4">
                                Order #{order.id} - {formatDate(order.createdAt)}
                            </h2>
                            <div class="mb-4">
                                <h3 class="font-pixel text-lg text-retroGray">Items:</h3>
                                {#each order.items as item}
                                    <div class="flex justify-between font-pixel text-retroGray">
                                        <span>
                                            {item.name || `Product ${item.productId}`} (x
                                            {item.quantity || 1})
                                        </span>
                                        <span>
                                            ${(item.price && item.quantity
                                                ? (item.price * item.quantity) / 100
                                                : 0
                                            ).toFixed(2)}
                                        </span>
                                    </div>
                                {/each}
                            </div>
                            <div class="mb-4">
                                <h3 class="font-pixel text-lg text-retroGray">Shipping Address:</h3>
                                <p class="font-body text-retroGray">
                                    Line: {order.line1}{order.line2 ? `, ${order.line2}` : ''}<br />
                                    City: {order.city}, {order.state}
                                    <br />
                                    Postal Code: {order.postalCode}<br />
                                    Country: {order.country}
                                </p>
                            </div>
                            <div
                                class="flex justify-between font-pixel text-xl text-retroGray border-t-2 border-retroBlack pt-3"
                            >
                                <span>Total</span>
                                <span class="text-retroCoral font-bold">
                                    ${(order.total / 100).toFixed(2)}
                                </span>
                            </div>
                        </div>
                    {/each}
                </div>
            {:else}
                <div
                    class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] text-center"
                >
                    <p class="font-pixel text-retroGray text-lg">
                        {showToast && toastMessage === 'Error fetching orders! Try again later'
                            ? 'Failed to load orders.'
                            : 'You haven’t placed any orders yet!'}
                    </p>
                    {#if showToast && toastMessage === 'Error fetching orders! Try again later'}
                        <button
                            on:click={fetchOrders}
                            class="font-pixel text-retroGray hover:text-retroCoral mt-4 inline-block"
                        >
                            Try Again
                        </button>
                    {:else}
                        <a
                            href="/products"
                            class="font-pixel text-retroGray hover:text-retroCoral mt-4 inline-block"
                        >
                            Start Shopping
                        </a>
                    {/if}
                </div>
            {/if}
        </section>
    {/if}
</main>

<style>
    .animate-spin {
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
