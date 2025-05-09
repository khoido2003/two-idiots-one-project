<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import Toast from '$lib/components/Toast.svelte';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { GLOBAL } from '$lib';
    import { auth } from '$lib/auth';
    import { page } from '$app/stores';

    export let data: { id: string } = { id: $page.params.id };

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
        line1?: string;
        line2?: string | null;
        city?: string;
        state?: string;
        postalCode?: string;
        country?: string;
        user: { email: string };
        status: string;
    }

    let order: Order | null = null;
    let orders: Order[] = [];
    let loading = true;
    let showToast = false;
    let toastMessage = '';
    let toastType: 'error' | 'success' | 'info' = 'error';
    let currentPage = 1;
    const limit = 10;

    async function fetchOrder() {
        try {
            loading = true;
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/orders/${$page.params.id}?userId=${$auth.user?.id}`,
                {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${$auth.token}`
                    }
                }
            );
            if (!response.ok) {
                if (response.status === 401) {
                    auth.set({ token: null, user: null });
                    goto('/signin');
                    return;
                }
                throw new Error('Failed to fetch order');
            }
            const data = await response.json();
            order = data.order || null;
            if (!order) {
                throw new Error('Order not found');
            }
        } catch (error) {
            console.error('Error fetching order:', error);
            showToast = true;
            toastMessage = 'Error fetching order details! Try again later';
            toastType = 'error';
        } finally {
            loading = false;
        }
    }

    async function fetchOrders() {
        try {
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/orders?userId=${$auth.user?.id}&page=${currentPage}&limit=${limit}`,
                {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${$auth.token}`
                    }
                }
            );
            if (!response.ok) {
                throw new Error('Failed to fetch orders');
            }
            const data = await response.json();
            orders = data.orders || [];
        } catch (error) {
            console.error('Error fetching orders:', error);
            showToast = true;
            toastMessage = 'Error fetching orders for pagination';
            toastType = 'error';
        }
    }

    async function updateOrderStatus(newStatus: string) {
        try {
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/orders/${$page.params.id}?userId=${$auth.user?.id}`,
                {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${$auth.token}`
                    },
                    body: JSON.stringify({ status: newStatus })
                }
            );
            if (!response.ok) throw new Error('Failed to update order status');
            if (order) {
                order.status = newStatus;
                showToast = true;
                toastMessage = 'Order status updated successfully';
                toastType = 'success';
            }
        } catch (err) {
            showToast = true;
            toastMessage = 'Error updating order status';
            toastType = 'error';
            console.error(err);
        }
    }

    onMount(async () => {
        // Redirect to signin if not authenticated or not an admin
        if (!$auth.token || !$auth.user || $auth.user.role !== 'admin') {
            orders = [];
            goto('/signin');
        }
        await Promise.all([fetchOrder(), fetchOrders()]);
        if (!order && !showToast) {
            showToast = true;
            toastMessage = 'Order not found.';
            toastType = 'info';
        } else if (order) {
            showToast = true;
            toastMessage = 'Order details fetched successfully';
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

    function navigateToOrder(orderId: number) {
        goto(`/admin/orders/${orderId}`);
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8 bg-retroCream min-h-screen">
    <Toast bind:visible={showToast} message={toastMessage} type={toastType} duration={3000} />

    {#if !$auth.token || !$auth.user || $auth.user.role !== 'admin'}
        <div class="flex justify-center items-center h-screen">
            <p class="font-pixel text-retroGray">Redirecting...</p>
        </div>
    {:else}
        <section class="mb-12 sm:mb-16">
            <h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-8 text-center">
                Order Details
            </h1>

            {#if loading}
                <div class="flex justify-center items-center">
                    <div class="font-pixel text-retroGray text-lg">Loading order details...</div>
                    <div
                        class="ml-4 animate-spin h-6 w-6 border-4 border-retroCoral border-t-transparent rounded-full"
                    ></div>
                </div>
            {:else if order}
                <div class="space-y-8">
                    <div
                        class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
                    >
                        <h2 class="font-pixel text-2xl text-retroGray mb-4">
                            Order #{order.id} - {formatDate(order.createdAt)}
                        </h2>
                        <div class="mb-4">
                            <h3 class="font-pixel text-lg text-retroGray">Customer:</h3>
                            <p class="font-body text-retroGray">{order.user.email}</p>
                        </div>
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
                            {#if order.line1 && order.city && order.state && order.postalCode && order.country}
                                <p class="font-body text-retroGray">
                                    Line: {order.line1}{order.line2 ? `, ${order.line2}` : ''}<br />
                                    City: {order.city}, {order.state}<br />
                                    Postal Code: {order.postalCode}<br />
                                    Country: {order.country}
                                </p>
                            {:else}
                                <p class="font-body text-retroGray">No shipping address provided</p>
                            {/if}
                        </div>
                        <div class="mb-4">
                            <h3 class="font-pixel text-lg text-retroGray">Status:</h3>
                            <select
                                class="font-pixel bg-retroCream border-2 border-retroBlack p-1"
                                value={order.status}
                                on:change={(e) => updateOrderStatus(e.target.value)}
                            >
                                <option value="pending">Pending</option>
                                <option value="processing">Processing</option>
                                <option value="shipped">Shipped</option>
                                <option value="delivered">Delivered</option>
                                <option value="cancelled">Cancelled</option>
                            </select>
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
                </div>
                <!-- Pagination -->
                <div class="flex justify-between mt-4">
                    <Button
                        variant="secondary"
                        subClass="bg-retroBlue text-retroCream hover:bg-retroCoral"
                        disabled={currentPage === 1}
                        onClick={() => {
                            currentPage--;
                            fetchOrders();
                        }}
                    >
                        Previous
                    </Button>
                    <span class="font-pixel text-retroGray">Page {currentPage}</span>
                    <Button
                        variant="secondary"
                        subClass="bg-retroBlue text-retroCream hover:bg-retroCoral"
                        disabled={orders.length < limit}
                        onClick={() => {
                            currentPage++;
                            fetchOrders();
                        }}
                    >
                        Next
                    </Button>
                </div>
                <!-- Related Orders -->
                {#if orders.length > 0}
                    <div class="mt-8">
                        <h2 class="font-pixel text-2xl text-retroGray mb-4">Other Orders</h2>
                        <div class="space-y-4">
                            {#each orders as otherOrder}
                                {#if otherOrder.id !== order.id}
                                    <div
                                        class="bg-retroCream p-4 rounded-lg border-2 border-retroGray cursor-pointer hover:bg-retroTeal hover:text-retroCream"
                                        on:click={() => navigateToOrder(otherOrder.id)}
                                    >
                                        <p class="font-pixel text-lg">
                                            Order #{otherOrder.id} - {formatDate(
                                                otherOrder.createdAt
                                            )} -
                                            {otherOrder.user.email}
                                        </p>
                                        <p class="font-pixel">
                                            Total: ${(otherOrder.total / 100).toFixed(2)} - Status:
                                            {otherOrder.status}
                                        </p>
                                    </div>
                                {/if}
                            {/each}
                        </div>
                    </div>
                {/if}
            {:else}
                <div
                    class="bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e] text-center"
                >
                    <p class="font-pixel text-retroGray text-lg">
                        {showToast && toastMessage.includes('Error')
                            ? 'Failed to load order details.'
                            : 'Order not found.'}
                    </p>
                    <button
                        on:click={fetchOrder}
                        class="font-pixel text-retroGray hover:text-retroCoral mt-4 inline-block"
                    >
                        Try Again
                    </button>
                </div>
            {/if}
        </section>
    {/if}
</main>

<style>
    .font-pixel {
        font-family: 'Jura', sans-serif;
    }
    .animate-spin {
        animation: spin 1s linear infinite;
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
</style>
