<!-- /src/lib/components/Navbar.svelte -->
<script>
  let isMenuOpen = false; // For mobile menu toggle
  const navItems = [
    { name: 'Home', href: '/' },
    { name: 'Products', href: '/products' },
    { name: 'Cart', href: '/cart' },
  ];
</script>

<nav class="bg-retroGray/95 text-retroCream py-4 shadow-xl sticky top-0 z-50">
  <div class="container mx-auto px-4 flex justify-between items-center">
    <h1 class="font-pixel text-lg sm:text-xl tracking-wider flex justify-center items-center gap-1.5">
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
            <a href={item.href} class="hover:text-retroCoral transition-colors text-sm sm:text-base">
              {item.name}
            </a>
          </li>
        {/each}
      </ul>
      <a
        href="/signin"
        class="bg-retroTeal text-retroCream px-4 py-2 rounded-md border-2 border-retroBlack hover:bg-retroCoral transition-all font-pixel text-sm sm:text-base"
      >
        Login
      </a>
    </div>
  </div>

  <!-- Mobile Menu (shown when toggled) -->
  {#if isMenuOpen}
    <ul class="sm:hidden bg-retroGray px-4 py-2">
      {#each navItems as item}
        <li class="py-2">
          <a href={item.href} class="hover:text-retroCoral transition-colors text-base">
            {item.name}
          </a>
        </li>
      {/each}
      <li class="py-2">
        <a
          href="/signin"
          class="block bg-retroTeal text-retroCream px-4 py-2 rounded-md border-2 border-retroBlack hover:bg-retroCoral transition-all font-pixel text-base"
        >
          Login
        </a>
      </li>
    </ul>
  {/if}
</nav>
