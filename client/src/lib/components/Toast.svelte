<!-- Toast.svelte -->
<script lang="ts">
    import { onMount, onDestroy } from 'svelte';

    // Define the possible types as a union type
    type ToastType = 'error' | 'success' | 'info';

    // Props with defaults and explicit types
    export let message: string = 'An error occurred!';
    export let type: ToastType = 'error';
    export let duration: number = 3000; // milliseconds
    export let visible: boolean = false;

    let timeout: number | undefined;

    // Handle auto-hide
    function startTimer() {
        if (duration > 0 && visible) {
            clearTimeout(timeout);
            timeout = setTimeout(() => {
                visible = false;
            }, duration);
        }
    }

    // Watch for visibility changes
    $: if (visible) {
        startTimer();
    }

    onMount(() => {
        if (visible) startTimer();
    });

    onDestroy(() => {
        clearTimeout(timeout);
    });

    // Style classes based on type with explicit typing
    const typeStyles: Record<ToastType, string> = {
        error: 'border-red-600 bg-red-100 text-red-800',
        success: 'border-green-600 bg-green-100 text-green-800',
        info: 'border-blue-600 bg-blue-100 text-blue-800'
    };
</script>

{#if visible}
    <div
        class="fixed bottom-4 right-4 z-[99999] max-w-xs animate-toast-in
           border-4 {typeStyles[type]} font-mono text-sm
           shadow-[8px_8px_0px_0px_rgba(0,0,0,0.3)]"
        style="font-family: 'Press Start 2P', cursive;"
    >
        <div class="p-4">
            <div class="flex items-center justify-between">
                <span>{message}</span>
                <button
                    on:click={() => (visible = false)}
                    class="ml-2 text-lg hover:text-black focus:outline-none"
                >
                    ×
                </button>
            </div>
        </div>
    </div>
{/if}

<style>
    @keyframes toastIn {
        from {
            transform: translateY(100%);
            opacity: 0;
        }
        to {
            transform: translateY(0);
            opacity: 1;
        }
    }

    .animate-toast-in {
        animation: toastIn 0.3s ease-out;
    }
</style>
