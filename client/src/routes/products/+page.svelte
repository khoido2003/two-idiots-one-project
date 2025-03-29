<!-- /src/routes/products/+page.svelte -->
<script lang="ts">
    import { onMount } from 'svelte';
    import Pagination from '$lib/components/Pagination.svelte';
    import SearchFilters from '$lib/components/SearchFilters.svelte';
    import ProductCard from '$lib/components/ProductCard.svelte';
    import { GLOBAL } from '$lib';

    // Product interface based on your backend models
    interface Product {
        id: number;
        name: string;
        price: number;
        description: string;
        category: { name: string };
        images: { url: string; isPrimary: boolean; order: number }[];
    }

    // State for fetched data
    let displayedProducts: Product[] = [];
    let categories: string[] = ['All'];
    let totalPages = 1;
    let totalItems = 0;
    let error: string = '';

    // Filter and pagination state
    let searchQuery = '';
    let selectedCategory = 'All';
    let priceRange: [number, number] = [0, 500];
    let currentPage = 1;
    const itemsPerPage = 8;

    // Fetch products from server with filters and pagination
    async function fetchProducts() {
        try {
            const url = new URL(`${GLOBAL.SERVER_URL}/products`);
            url.searchParams.set('page', currentPage.toString());
            url.searchParams.set('limit', itemsPerPage.toString());
            if (searchQuery) url.searchParams.set('search', searchQuery);
            if (selectedCategory !== 'All') url.searchParams.set('category', selectedCategory);
            url.searchParams.set('minPrice', priceRange[0].toString());
            url.searchParams.set('maxPrice', priceRange[1].toString());

            const response = await fetch(url);
            if (!response.ok) throw new Error('Failed to fetch products');
            const data = await response.json();

            displayedProducts = data.products;
            totalPages = data.totalPages;
            totalItems = data.totalItems;
            currentPage = data.currentPage;
        } catch (err) {
            error = `Error fetching data: ${err}`;
        }
    }

    // Fetch categories on mount
    onMount(async () => {
        try {
            const categoriesResponse = await fetch(`${GLOBAL.SERVER_URL}/categories`);
            if (!categoriesResponse.ok) throw new Error('Failed to fetch categories');
            const categoriesData = await categoriesResponse.json();
            categories = [
                'All',
                ...categoriesData.categories.map((cat: { name: string }) => cat.name)
            ];

            // Initial fetch of products
            await fetchProducts();
        } catch (err) {
            error = `Error fetching categories: ${err}`;
        }
    });

    // Handle filter submission
    function handleFilterSubmit(filters: {
        searchQuery: string;
        selectedCategory: string;
        priceRange: [number, number];
    }) {
        searchQuery = filters.searchQuery;
        selectedCategory = filters.selectedCategory;
        priceRange = filters.priceRange;
        currentPage = 1; // Reset to first page
        fetchProducts();
    }

    // Handle page change with scroll adjustment
    function handlePageChange(page: number, event?: Event) {
        if (event) event.preventDefault();

        const viewportHeight = window.innerHeight;
        const currentBottom = window.scrollY + viewportHeight;

        currentPage = page;
        fetchProducts();

        requestAnimationFrame(() => {
            const newDocumentHeight = document.documentElement.scrollHeight;
            const newScrollPosition = newDocumentHeight - viewportHeight;

            if (newScrollPosition >= 0 && currentBottom > window.scrollY) {
                window.scrollTo(0, newScrollPosition);
            }
        });
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <!-- Page Header -->
    <section class="mb-12 sm:mb-16 text-center">
        <h1 class="font-pixel text-3xl sm:text-4xl md:text-5xl text-retroGray mb-4">
            All Retro Gear
        </h1>
        <p class="text-base sm:text-lg text-retroGray font-body max-w-2xl mx-auto">
            Browse our full collection of vintage tech treasures!
        </p>
    </section>

    <!-- Error Display -->
    {#if error}
        <section class="mb-12 sm:mb-16 text-center">
            <p class="text-retroCoral font-pixel text-lg">{error}</p>
        </section>
    {/if}

    <!-- Search and Filters -->
    <section class="mb-12 sm:mb-16">
        <SearchFilters
            {searchQuery}
            {selectedCategory}
            {priceRange}
            {categories}
            onSubmit={handleFilterSubmit}
        />
    </section>

    <!-- Products Grid -->
    <section class="mb-12 sm:mb-16">
        {#if displayedProducts.length > 0}
            <div
                class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 gap-6 sm:gap-8 justify-items-center px-4"
            >
                {#each displayedProducts as product}
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
        {:else}
            <p class="text-center font-pixel text-retroGray text-lg">
                No products found. Try adjusting your filters!
            </p>
        {/if}
    </section>

    <!-- Pagination -->
    <section class="mb-12 sm:mb-16">
        <Pagination {currentPage} {totalPages} onPageChange={handlePageChange} />
    </section>
</main>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
