<!-- /src/routes/signin/+page.svelte -->
<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/auth';
    import { GLOBAL } from '$lib';

    let email = '';
    let password = '';
    let error = '';

    async function handleSignIn(event: Event) {
        event.preventDefault();

        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/signin`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password })
            });
            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.message || 'Signin failed');
            }

            const data = await response.json();
            auth.set({
                token: data.token,
                user: {
                    id: data.user.id,
                    email: data.user.email,
                    firstName: data.user.firstName,
                    lastName: data.user.lastName,
                    role: data.user.role,
                    phone: data.user.phone
                }
            });
            error = '';
            goto('/');         } catch (err: any) {
            error = err.message;
        }
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <section
        class="max-w-md mx-auto bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
    >
        <h1 class="font-pixel text-3xl text-retroGray mb-6 text-center">Sign In</h1>

        {#if error}
            <p class="text-retroCoral font-pixel text-center mb-4">{error}</p>
        {/if}

        <form on:submit={handleSignIn} class="space-y-6">
            <div>
                <label for="email" class="font-pixel text-retroGray text-sm sm:text-base"
                    >Email</label
                >
                <input
                    type="email"
                    id="email"
                    bind:value={email}
                    class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
                    placeholder="you@example.com"
                    required
                />
            </div>
            <div>
                <label for="password" class="font-pixel text-retroGray text-sm sm:text-base"
                    >Password</label
                >
                <input
                    type="password"
                    id="password"
                    bind:value={password}
                    class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
                    placeholder="••••••••"
                    required
                />
            </div>
            <Button
                variant="primary"
                type="submit"
                subClass="w-full bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3"
            >
                Sign In
            </Button>
        </form>

        <p class="text-center mt-4 font-pixel text-retroGray text-sm sm:text-base">
            Don’t have an account?
            <a href="/signup" class="text-retroTeal hover:text-retroCoral transition-colors"
                >Sign Up</a
            >
        </p>
    </section>
</main>
