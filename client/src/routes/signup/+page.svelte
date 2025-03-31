<!-- /src/routes/signup/+page.svelte -->
<script lang="ts">
    import Button from '$lib/components/Button.svelte';
    import { goto } from '$app/navigation';
    import { auth } from '$lib/auth';
    import { GLOBAL } from '$lib';

    let firstName = '';
    let lastName = '';
    let email = '';
    let phone = '';
    let password = '';
    let confirmPassword = '';
    let error = '';

    async function handleSignUp(event: Event) {
        event.preventDefault();
        if (password !== confirmPassword) {
            error = 'Passwords do not match!';
            return;
        }

        try {
            const response = await fetch(`${GLOBAL.SERVER_URL}/signup`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ email, password, firstName, lastName, phone })
            });
            if (!response.ok) {
                const data = await response.json();
                throw new Error(data.message || 'Signup failed');
            }

            const data = await response.json();
            auth.set({
                token: data.token,
                user: {
                    id: data.user.id,
                    email: data.user.email,
                    firstName: data.user.firstName,
                    lastName: data.user.lastName,
                    role: data.user.role
                }
            });
            error = '';
            goto('/');
        } catch (err: any) {
            error = err;
        }
    }
</script>

<main class="container mx-auto px-4 py-6 sm:py-7">
    <section
        class="max-w-md mx-auto bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
    >
        <h1 class="font-pixel text-3xl text-retroGray mb-6 text-center">Sign Up</h1>

        {#if error}
            <p class="text-retroCoral font-pixel text-center mb-4">{error}</p>
        {/if}

        <form on:submit={handleSignUp} class="flex flex-col gap-2">
            <div>
                <label for="firstName" class="font-pixel text-retroGray text-sm sm:text-base"
                    >First Name</label
                >
                <input
                    type="text"
                    id="firstName"
                    bind:value={firstName}
                    class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
                    placeholder="Jane"
                    required
                />
            </div>
            <div>
                <label for="lastName" class="font-pixel text-retroGray text-sm sm:text-base"
                    >Last Name</label
                >
                <input
                    type="text"
                    id="lastName"
                    bind:value={lastName}
                    class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
                    placeholder="Doe"
                    required
                />
            </div>
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
                <label for="phone" class="font-pixel text-retroGray text-sm sm:text-base"
                    >Phone</label
                >
                <input
                    type="tel"
                    id="phone"
                    bind:value={phone}
                    class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
                    placeholder="123-456-7890"
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
            <div>
                <label for="confirm-password" class="font-pixel text-retroGray text-sm sm:text-base"
                    >Confirm Password</label
                >
                <input
                    type="password"
                    id="confirm-password"
                    bind:value={confirmPassword}
                    class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
                    placeholder="••••••••"
                    required
                />
            </div>
            <Button
                variant="primary"
                type="submit"
                subClass="w-full bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 mt-2"
            >
                Sign Up
            </Button>
        </form>

        <p class="text-center mt-4 font-pixel text-retroGray text-sm sm:text-base">
            Already have an account?
            <a href="/signin" class="text-retroTeal hover:text-retroCoral transition-colors"
                >Sign In</a
            >
        </p>
    </section>
</main>
