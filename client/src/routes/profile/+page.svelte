<!-- /src/routes/profile/+page.svelte -->
<script lang="ts">
    import { auth, logout } from '$lib/auth';
    import { goto } from '$app/navigation';
    import Button from '$lib/components/Button.svelte';
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
    <section
        class="max-w-lg mx-auto bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
    >
        <h1 class="font-pixel text-3xl text-retroGray mb-6 text-center">Your Profile</h1>

        {#if $auth.user}
            <div class="space-y-4">
                <!-- First Name -->
                <div class="flex justify-between items-center">
                    <span class="font-pixel text-retroGray text-sm sm:text-base">First Name:</span>
                    <span class="font-pixel text-retroCoral text-sm sm:text-base"
                        >{$auth.user.firstName}</span
                    >
                </div>

                <!-- Last Name -->
                <div class="flex justify-between items-center">
                    <span class="font-pixel text-retroGray text-sm sm:text-base">Last Name:</span>
                    <span class="font-pixel text-retroCoral text-sm sm:text-base"
                        >{$auth.user.lastName}</span
                    >
                </div>

                <!-- Email -->
                <div class="flex justify-between items-center">
                    <span class="font-pixel text-retroGray text-sm sm:text-base">Email:</span>
                    <span class="font-pixel text-retroCoral text-sm sm:text-base"
                        >{$auth.user.email}</span
                    >
                </div>

                <!-- Phone -->
                <div class="flex justify-between items-center">
                    <span class="font-pixel text-retroGray text-sm sm:text-base">Phone:</span>
                    <span class="font-pixel text-retroCoral text-sm sm:text-base"
                        >{$auth.user.phone}</span
                    >
                </div>

                <!-- Role -->
                <div class="flex justify-between items-center">
                    <span class="font-pixel text-retroGray text-sm sm:text-base">Role:</span>
                    <span class="font-pixel text-retroCoral text-sm sm:text-base"
                        >{$auth.user.role}</span
                    >
                </div>

                <!-- Edit Profile Button -->
                <Button
                    variant="primary"
                    subClass="w-full bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 mt-6"
                    onClick={() => console.log('Edit profile clicked - implement later')}
                >
                    Edit Profile
                </Button>

                <!-- Logout Button -->
                <Button
                    variant="secondary"
                    subClass="w-full bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 mt-2"
                    onClick={() => {
                        logout();
                        goto('/');
                    }}
                >
                    Logout
                </Button>
            </div>
        {:else}
            <p class="text-center font-pixel text-retroGray text-lg">
                Please sign in to view your profile.
            </p>
            <Button
                variant="primary"
                subClass="w-full bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 mt-6"
                onClick={() => goto('/signin')}
            >
                Sign In
            </Button>
        {/if}

        <!-- Back to Home Link -->
        <div class="text-center mt-6">
            <a
                href="/"
                class="font-pixel text-retroGray hover:text-retroCoral transition-colors text-sm sm:text-base"
            >
                ← Back to Home
            </a>
        </div>
    </section>
</main>
