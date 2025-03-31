<!-- /src/lib/components/Navbar.svelte -->
<script lang="ts">
    import { auth, logout } from '$lib/auth';
    import { goto } from '$app/navigation';

    let isMenuOpen = false;

    const navItems = [
        { name: 'Home', href: '/' },
        { name: 'Products', href: '/products' },
        { name: 'Cart', href: '/cart' }
    ];

    // Random color generator for avatar
    const colors = ['#FF6F61', '#26A69A', '#AB47BC', '#42A5F5', '#FFCA28', '#8D6E63'];
    function getRandomColor() {
        return colors[Math.floor(Math.random() * colors.length)];
    }

       $: avatarLetter = $auth.user?.firstName ? $auth.user.firstName[0].toUpperCase() : '';
    $: avatarColor = $auth.user ? getRandomColor() : '#FFFFFF';
</script>

<nav class="bg-retroGray/95 text-retroCream py-4 shadow-xl sticky top-0 z-50">
    <div class="container mx-auto px-4 flex justify-between items-center">
        <h1
            class="font-pixel text-lg sm:text-xl tracking-wider flex justify-center items-center gap-1.5"
        >
            <img class="h-6 w-6 block" src="/assets/logo.svg" height="25" width="25" alt="Logo" />
            RetroTech
        </h1>

        <!-- Mobile Menu Button -->
        <button class="sm:hidden text-2xl" on:click={() => (isMenuOpen = !isMenuOpen)}>
            {#if isMenuOpen}✕{:else}☰{/if}
        </button>

        <!-- Desktop Menu -->
        <div class="hidden sm:flex items-center space-x-6">
            <ul class="flex space-x-6">
                {#each navItems as item}
                    <li>
                        <a
                            href={item.href}
                            class="hover:text-retroCoral transition-colors text-sm sm:text-base"
                        >
                            {item.name}
                        </a>
                    </li>
                {/each}
            </ul>
            {#if $auth.token}
                <!-- Profile Avatar -->
                <div
                    class="flex items-center justify-center w-10 h-10 rounded-full cursor-pointer"
                    style="background-color: {avatarColor}"
                    on:click={() => goto('/profile')}
                    title={`${$auth.user?.firstName} ${$auth.user?.lastName}`}
                >
                    <span class="font-pixel text-lg text-white">{avatarLetter}</span>
                </div>
            {:else}
                <!-- Login Button -->
                <a
                    href="/signin"
                    class="bg-retroTeal text-retroCream px-4 py-2 rounded-md border-2 border-retroBlack hover:bg-retroCoral transition-all font-pixel text-sm sm:text-base"
                >
                    Login
                </a>
            {/if}
        </div>
    </div>

    <!-- Mobile Menu -->
    {#if isMenuOpen}
        <ul class="sm:hidden bg-retroGray px-4 py-2">
            {#each navItems as item}
                <li class="py-2">
                    <a href={item.href} class="hover:text-retroCoral transition-colors text-base">
                        {item.name}
                    </a>
                </li>
            {/each}
            {#if $auth.token}
                <li class="py-2 flex items-center gap-2">
                    <div
                        class="flex items-center justify-center w-8 h-8 rounded-full"
                        style="background-color: {avatarColor}"
                    >
                        <span class="font-pixel text-base text-white">{avatarLetter}</span>
                    </div>
                    <span class="text-base text-retroCream"
                        >{$auth.user?.firstName} {$auth.user?.lastName}</span
                    >
                </li>
                <li class="py-2">
                    <button
                        on:click={() => {
                            logout();
                            isMenuOpen = false;
                            goto('/');
                        }}
                        class="block w-full text-left text-retroCoral hover:text-retroTeal transition-colors font-pixel text-base"
                    >
                        Logout
                    </button>
                </li>
            {:else}
                <li class="py-2">
                    <a
                        href="/signin"
                        class="block bg-retroTeal text-retroCream px-4 py-2 rounded-md border-2 border-retroBlack hover:bg-retroCoral transition-all font-pixel text-base"
                    >
                        Login
                    </a>
                </li>
            {/if}
        </ul>
    {/if}
</nav>
