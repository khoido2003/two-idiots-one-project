<script lang="ts">
    import { goto } from '$app/navigation';
    import Button from './Button.svelte';

    export let product = {
        id: 1,
        name: 'NES Console',
        price: 59.99,
        image: 'assets/nes.webp',
        description: 'Classic 8-bit gaming console.',
        category: 'Electronics',
        stock: 5
    };
    export let subClass = '';
    export let onEdit: (id: number) => void;
    export let onDelete: (id: number) => void;

    function handleClick(e: MouseEvent) {
        const target = e.target as HTMLElement;
        if (target.closest('button')) {
            return;
        }
        onEdit(product.id);
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
    <p class="text-gray-600 text-sm mt-1">Stock: {product.stock}</p>
    <div class="flex gap-2 mt-2">
        <Button variant="primary" subClass="flex-1" onClick={() => onEdit(product.id)}>Edit</Button>
        <Button
            variant="secondary"
            subClass="flex-1 hover:bg-retroCoral"
            onClick={() => onDelete(product.id)}>Delete</Button
        >
    </div>
</div>

<style>
    .pixelated {
        image-rendering: pixelated;
    }
</style>
