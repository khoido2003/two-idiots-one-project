<!-- /src/lib/components/ProductCard.svelte -->
<script lang="ts">
    import { goto } from '$app/navigation';
    import { GLOBAL } from '$lib';
    import { auth } from '$lib/auth';
    import Button from './Button.svelte';
    import Toast from './Toast.svelte';

    export let product = {
        id: 1,
        name: 'NES Console',
        price: 59.99,
        image: 'assets/nes.webp',
        description: 'Classic 8-bit gaming console.',
        category: 'Electron'
    };
    export let subClass = '';

    let showToast = false;
    let toastMessage = '';
    let toastType: 'error' | 'success' | 'info' = 'error';

    function handleClick(e: MouseEvent) {
        const target = e.target as HTMLElement;
        if (target.closest('button')) {
            console.log('Click was on button, skipping navigation');
            return;
        }
        console.log('Navigating to:', `/products/${product.id}`);
        goto(`/products/${product.id}`);
    }

    async function addToCart(e: MouseEvent) {
        e.preventDefault();
        e.stopPropagation();

        if (!$auth.token || !$auth.user?.id) {
            goto('/signin');
        }

        let token: string | null = $auth.token;
        let userId: number | undefined = $auth.user?.id;

        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/cart?userId=${userId}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${token}`
                },
                body: JSON.stringify({
                    productId: product.id,
                    quantity: 1
                }),
                cache: 'no-store'
            });

            if (!response.ok) {
                const errData = await response.json();
                toastMessage = 'Failed to add to cart!';
                toastType = 'error';
                showToast = true;
                console.log(errData);
            } else {
                const data = await response.json();
                toastMessage = 'Add to cart successfully!';
                toastType = 'success';
                showToast = true;
                console.log(data.message, data.cartItem);
            }
        } catch (error) {
            console.error('Error adding to cart', error);
            toastMessage = 'Failed to add to cart!';
            toastType = 'error';
            showToast = true;
        }
    }
</script>

<!-- svelte-ignore a11y_click_events_have_key_events -->
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
    class={`block bg-white w-72 p-3 sm:p-4 rounded-lg shadow-md hover:shadow-xl transition-shadow ${subClass}`}
    on:click={handleClick}
>
    <div
        class="w-40 h-40 sm:w-60 sm:h-48 mx-auto rounded-md mb-3 sm:mb-4 pixelated"
        style={`background-image: url('${product.image}'); background-size: cover; background-position: center;`}
    ></div>
    <h3 class="font-pixel text-retroGray text-base sm:text-lg truncate">{product.name}</h3>
    <p class="text-xs sm:text-sm text-gray-600 line-clamp-2">{product.description}</p>
    <p class="text-retroCoral font-bold mt-2 text-sm sm:text-base">${product.price}</p>
    <Button variant="primary" subClass="mt-2" onClick={addToCart}>Add to Cart</Button>
    <!-- Toast component -->
    <Toast bind:visible={showToast} message={toastMessage} type={toastType} duration={3000} />
</div>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
