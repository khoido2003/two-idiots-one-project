// /src/routes/+page.ts
import type { PageLoad } from './$types';
import { GLOBAL } from '$lib';

export const load: PageLoad = async ({ fetch }) => {
    // Fetch Featured Products (top 4 products)
    const featuredResponse = await fetch(`${GLOBAL.SERVER_URL}/products?limit=4`, {
        headers: { 'Cache-Control': 'no-cache' }
    });
    if (!featuredResponse.ok) throw new Error('Failed to fetch featured products');
    const featuredData = await featuredResponse.json();
    const featuredProducts = featuredData.products;

    // Fetch New Arrivals (sorted by createdAt, top 4)
    const newArrivalsResponse = await fetch(
        `${GLOBAL.SERVER_URL}/products?sort=createdAtDesc&limit=4`,
        {
            headers: { 'Cache-Control': 'no-cache' }
        }
    );
    if (!newArrivalsResponse.ok) throw new Error('Failed to fetch new arrivals');
    const newArrivalsData = await newArrivalsResponse.json();
    const newArrivals = newArrivalsData.products;

    // Fetch Limited Deals (low stock, limit 1)
    const limitedDealsResponse = await fetch(`${GLOBAL.SERVER_URL}/products?maxStock=5&limit=1`, {
        headers: { 'Cache-Control': 'no-cache' }
    });
    if (!limitedDealsResponse.ok) throw new Error('Failed to fetch limited deals');
    const limitedDealsData = await limitedDealsResponse.json();
    const limitedDeals = limitedDealsData.products.map((deal: any) => ({
        ...deal,
        originalPrice: deal.price * 1.5
    }));

    return {
        featuredProducts,
        newArrivals,
        limitedDeals
    };
};
