<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import AdminProductCard from '$lib/components/AdminProductCard.svelte';
    import 'iconify-icon';
    import { onMount } from 'svelte';
    import { goto } from '$app/navigation';
    import { GLOBAL } from '$lib';
    import { auth } from '$lib/auth';

    // Define types for dashboard data
    interface DashboardMetrics {
        totalUsers: number;
        totalOrders: number;
        totalRevenue: number;
        lowStockProducts: number;
    }

    interface Order {
        id: number;
        userId: number;
        total: number;
        createdAt: string;
        status: string;
        user: { email: string };
    }

    interface Product {
        id: number;
        name: string;
        price: number;
        stock: number;
        images: { url: string; isPrimary: boolean }[];
        category: { name: string };
        description: string;
    }

    interface User {
        id: number;
        email: string;
        firstName: string;
        lastName: string;
        role: 'user' | 'admin';
    }

    // State for dashboard data
    let metrics: DashboardMetrics = {
        totalUsers: 0,
        totalOrders: 0,
        totalRevenue: 0,
        lowStockProducts: 0
    };
    let recentOrders: Order[] = [];
    let products: Product[] = [];
    let users: User[] = [];
    let lowStockProducts: Product[] = [];
    let loading = true;
    let error: string | null = null;
    let showLowStockModal = false;

    // Pagination state
    let ordersPage = 1;
    let productsPage = 1;
    const limit = 10;

    // Fetch dashboard data
    async function fetchDashboardData() {
        try {
            loading = true;
            const userId = $auth?.user?.id;
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/dashboard?userId=${userId}&page=${ordersPage}&limit=${limit}`,
                {
                    headers: {
                        Authorization: `Bearer ${$auth.token}`
                    }
                }
            );
            if (!response.ok) throw new Error('Failed to fetch dashboard data');

            const data = await response.json();
            metrics = data.metrics;
            recentOrders = data.recentOrders;
            products = data.products;
            users = data.users;
            lowStockProducts = data.lowStockProductDetails;
        } catch (err) {
            error = 'Error loading dashboard data';
            console.error(err);
        } finally {
            loading = false;
        }
    }

    // Fetch data on mount and when pagination changes
    onMount(fetchDashboardData);
    $: if (ordersPage || productsPage) {
        fetchDashboardData();
    }

    // Handle product deletion
    async function deleteProduct(productId: number) {
        if (!confirm('Are you sure you want to delete this product?')) return;
        try {
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/products/${productId}?userId=${$auth?.user?.id}`,
                {
                    method: 'DELETE',
                    headers: {
                        Authorization: `Bearer ${$auth.token}`
                    }
                }
            );
            if (!response.ok) throw new Error('Failed to delete product');
            products = products.filter((p) => p.id !== productId);
            lowStockProducts = lowStockProducts.filter((p) => p.id !== productId);
        } catch (err) {
            error = 'Error deleting product';
            console.error(err);
        }
    }

    // Handle product edit
    function editProduct(productId: number) {
        goto(`/admin/products/${productId}`);
    }

    // Handle user role change
    async function updateUserRole(userId: number, newRole: 'user' | 'admin') {
        try {
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/users/${userId}?userId=${$auth?.user?.id}`,
                {
                    method: 'PATCH',
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${$auth.token}`
                    },
                    body: JSON.stringify({ role: newRole })
                }
            );
            if (!response.ok) throw new Error('Failed to update user role');
            users = users.map((u) => (u.id === userId ? { ...u, role: newRole } : u));
        } catch (err) {
            error = 'Error updating user role';
            console.error(err);
        }
    }

    // Handle order status update
    async function updateOrderStatus(orderId: number, newStatus: string) {
        try {
            const response = await fetch(
                `${GLOBAL.SERVER_URL}/admin/orders/${orderId}?userId=${$auth?.user?.id}`,
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
            recentOrders = recentOrders.map((o) =>
                o.id === orderId ? { ...o, status: newStatus } : o
            );
        } catch (err) {
            error = 'Error updating order status';
            console.error(err);
        }
    }

    // Navigate to order details
    function viewOrder(orderId: number) {
        goto(`/admin/orders/${orderId}`);
    }

    // Navigate to add product page
    function addProduct() {
        goto('/admin/products');
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8 bg-retroCream min-h-screen">
    <h1 class="font-pixel text-3xl sm:text-4xl text-retroGray mb-8 text-center">Admin Dashboard</h1>

    {#if loading}
        <p class="text-center font-pixel text-retroGray">Loading...</p>
    {:else if error}
        <p class="text-center font-pixel text-retroCoral">{error}</p>
    {:else}
        <!-- Metrics Section -->
        <section class="mb-12 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
            <div class="bg-retroGray p-6 rounded-lg shadow-[4px_4px_0_#ff6f61] text-center">
                <h2 class="font-pixel text-xl text-retroCoral mb-2">Total Users</h2>
                <p class="font-pixel text-2xl text-retroCream">{metrics.totalUsers}</p>
            </div>
            <div class="bg-retroGray p-6 rounded-lg shadow-[4px_4px_0_#ff6f61] text-center">
                <h2 class="font-pixel text-xl text-retroCoral mb-2">Total Orders</h2>
                <p class="font-pixel text-2xl text-retroCream">{metrics.totalOrders}</p>
            </div>
            <div class="bg-retroGray p-6 rounded-lg shadow-[4px_4px_0_#ff6f61] text-center">
                <h2 class="font-pixel text-xl text-retroCoral mb-2">Total Revenue</h2>
                <p class="font-pixel text-2xl text-retroCream">
                    ${(metrics.totalRevenue / 100).toFixed(2)}
                </p>
            </div>
            <div
                class="bg-retroGray p-6 rounded-lg shadow-[4px_4px_0_#ff6f61] text-center cursor-pointer"
                on:click={() => (showLowStockModal = true)}
            >
                <h2 class="font-pixel text-xl text-retroCoral mb-2">Low Stock Alerts</h2>
                <p class="font-pixel text-2xl text-retroCream">{metrics.lowStockProducts}</p>
            </div>
        </section>

        <!-- Recent Orders Section -->
        <section class="mb-12">
            <h2 class="font-pixel text-2xl text-retroGray mb-6">Recent Orders</h2>
            <div class="overflow-x-auto">
                <table
                    class="w-full bg-white border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
                >
                    <thead>
                        <tr class="bg-retroGray text-retroCream">
                            <th class="p-4 font-pixel text-left">Order ID</th>
                            <th class="p-4 font-pixel text-left">User</th>
                            <th class="p-4 font-pixel text-left">Total</th>
                            <th class="p-4 font-pixel text-left">Date</th>
                            <th class="p-4 font-pixel text-left">Status</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each recentOrders as order}
                            <tr
                                class="border-b-2 border-retroBlack cursor-pointer hover:bg-retroTeal hover:text-retroCream"
                                on:click={() => viewOrder(order.id)}
                            >
                                <td class="p-4 font-pixel">{order.id}</td>
                                <td class="p-4 font-pixel">{order.user.email}</td>
                                <td class="p-4 font-pixel">${(order.total / 100).toFixed(2)}</td>
                                <td class="p-4 font-pixel"
                                    >{new Date(order.createdAt).toLocaleDateString()}</td
                                >
                                <td class="p-4 font-pixel">
                                    <select
                                        class="font-pixel bg-retroCream border-2 border-retroBlack p-1"
                                        value={order.status}
                                        on:change={(e) => {
                                            e.stopPropagation(); // Prevent row click
                                            updateOrderStatus(order.id, e.target.value);
                                        }}
                                    >
                                        <option value="pending">Pending</option>
                                        <option value="processing">Processing</option>
                                        <option value="shipped">Shipped</option>
                                        <option value="delivered">Delivered</option>
                                        <option value="cancelled">Cancelled</option>
                                    </select>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
            <div class="flex justify-between mt-4">
                <Button
                    variant="secondary"
                    disabled={ordersPage === 1}
                    onClick={() => ordersPage--}
                >
                    Previous
                </Button>
                <span class="font-pixel text-retroGray">Page {ordersPage}</span>
                <Button
                    variant="secondary"
                    disabled={recentOrders.length < limit}
                    onClick={() => ordersPage++}
                >
                    Next
                </Button>
            </div>
        </section>

        <!-- Product Management Section -->
        <section class="mb-12">
            <div class="flex justify-between items-center mb-6">
                <h2 class="font-pixel text-2xl text-retroGray">Product Management</h2>
                <Button variant="primary" subClass="hover:bg-retroCoral" onClick={addProduct}>
                    Add Product
                </Button>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                {#each products as product}
                    <AdminProductCard
                        product={{
                            id: product.id,
                            name: product.name,
                            price: product.price,
                            image: product.images.find((img) => img.isPrimary)?.url || '',
                            description: product.description,
                            category: product.category.name,
                            stock: product.stock
                        }}
                        subClass="transform group-hover:scale-105 transition-transform duration-300"
                        onEdit={editProduct}
                        onDelete={deleteProduct}
                    />
                {/each}
            </div>
            <div class="flex justify-between mt-4">
                <Button
                    variant="secondary"
                    disabled={productsPage === 1}
                    onClick={() => productsPage--}
                >
                    Previous
                </Button>
                <span class="font-pixel text-retroGray">Page {productsPage}</span>
                <Button
                    variant="secondary"
                    disabled={products.length < limit}
                    onClick={() => productsPage++}
                >
                    Next
                </Button>
            </div>
        </section>

        <!-- User Management Section -->
        <section>
            <h2 class="font-pixel text-2xl text-retroGray mb-6">User Management</h2>
            <div class="overflow-x-auto">
                <table
                    class="w-full bg-white border-4 border-retroBlack shadow-[4px_4px_0_#26A69A]"
                >
                    <thead>
                        <tr class="bg-retroGray text-retroCream">
                            <th class="p-4 font-pixel text-left">User ID</th>
                            <th class="p-4 font-pixel text-left">Email</th>
                            <th class="p-4 font-pixel text-left">Name</th>
                            <th class="p-4 font-pixel text-left">Role</th>
                            <th class="p-4 font-pixel text-left">Actions</th>
                        </tr>
                    </thead>
                    <tbody>
                        {#each users as user}
                            <tr class="border-b-2 border-retroBlack">
                                <td class="p-4 font-pixel">{user.id}</td>
                                <td class="p-4 font-pixel">{user.email}</td>
                                <td class="p-4 font-pixel">{user.firstName} {user.lastName}</td>
                                <td class="p-4 font-pixel">{user.role}</td>
                                <td class="p-4 font-pixel">
                                    <select
                                        class="font-pixel bg-retroCream border-2 border-retroBlack p-1"
                                        value={user.role}
                                        on:change={(e) => updateUserRole(user.id, e.target.value)}
                                    >
                                        <option value="user">User</option>
                                        <option value="admin">Admin</option>
                                    </select>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        </section>
    {/if}

    <!-- Low Stock Modal -->
    {#if showLowStockModal}
        <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
            <div class="bg-retroCream p-6 rounded-lg max-w-4xl w-full max-h-[80vh] overflow-y-auto">
                <h2 class="font-pixel text-2xl text-retroGray mb-6">Low Stock Products</h2>
                {#if lowStockProducts.length === 0}
                    <p class="font-pixel text-retroGray">No low stock products.</p>
                {:else}
                    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                        {#each lowStockProducts as product}
                            <AdminProductCard
                                product={{
                                    id: product.id,
                                    name: product.name,
                                    price: product.price,
                                    image: product.images.find((img) => img.isPrimary)?.url || '',
                                    description: product.description,
                                    category: product.category.name,
                                    stock: product.stock
                                }}
                                subClass="transform group-hover:scale-105 transition-transform duration-300"
                                onEdit={editProduct}
                                onDelete={deleteProduct}
                            />
                        {/each}
                    </div>
                {/if}
                <Button
                    variant="secondary"
                    subClass="mt-6 w-full hover:bg-retroCoral"
                    onClick={() => (showLowStockModal = false)}
                >
                    Close
                </Button>
            </div>
        </div>
    {/if}
</main>

<style>
    .font-pixel {
        font-family: 'Jura', sans-serif;
    }
    .pixelated {
        image-rendering: pixelated;
    }
</style>
