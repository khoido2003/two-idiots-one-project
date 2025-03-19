<script lang="ts">
	import Button from '$lib/components/Button.svelte';

	let email = '';
	let password = '';

	function handleSignIn(event: Event) {
		event.preventDefault();
		console.log('Sign In:', { email, password });
		// Add authentication logic later
	}
</script>

<main class="container mx-auto px-4 py-6 sm:py-8">
	<section
		class="max-w-md mx-auto bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
	>
		<h1 class="font-pixel text-3xl text-retroGray mb-6 text-center">Sign In</h1>

		<form on:submit={handleSignIn} class="space-y-6">
			<!-- Email -->
			<div>
				<label for="email" class="font-pixel text-retroGray text-sm sm:text-base">Email</label>
				<input
					type="email"
					id="email"
					bind:value={email}
					class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
					placeholder="you@example.com"
					required
				/>
			</div>

			<!-- Password -->
			<div>
				<label for="password" class="font-pixel text-retroGray text-sm sm:text-base">Password</label
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

			<!-- Submit -->
			<Button
				variant="primary"
				type="submit"
				subClass="w-full bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3"
			>
				Sign In
			</Button>
		</form>

		<!-- Link to Sign Up -->
		<p class="text-center mt-4 font-pixel text-retroGray text-sm sm:text-base">
			Don’t have an account?
			<a href="/signup" class="text-retroTeal hover:text-retroCoral transition-colors">Sign Up</a>
		</p>
	</section>
</main>
