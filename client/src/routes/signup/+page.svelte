<script lang="ts">
	import Button from '$lib/components/Button.svelte';

	let name = '';
	let email = '';
	let phone = '';
	let password = '';
	let confirmPassword = '';

	function handleSignUp(event: Event) {
		event.preventDefault();
		if (password !== confirmPassword) {
			console.log('Passwords do not match!');
			return;
		}
		console.log('Sign Up:', { email, password });
		// Add authentication logic later
	}
</script>

<main class="container mx-auto px-4 py-6 sm:py-7">
	<section
		class="max-w-md mx-auto bg-retroCream p-6 rounded-2xl border-4 border-retroGray shadow-[4px_4px_0_#2e2e2e]"
	>
		<h1 class="font-pixel text-3xl text-retroGray mb-6 text-center">Sign Up</h1>

		<form on:submit={handleSignUp} class="flex flex-col gap-2">
			<!-- FullName -->
			<div>
				<label for="name" class="font-pixel text-retroGray text-sm sm:text-base">Fullname</label>
				<input
					type="text"
					id="name"
					bind:value={name}
					class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
					placeholder="Jane Doe"
					required
				/>
			</div>

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

			<!-- PhoneNumber -->
			<div>
				<label for="phone" class="font-pixel text-retroGray text-sm sm:text-base">Your phone</label>
				<input
					type="tel"
					id="phone"
					bind:value={phone}
					class="w-full mt-2 px-3 py-2 bg-retroCream text-retroGray border-2 border-retroBlack rounded-md font-body focus:outline-none focus:border-retroTeal"
					placeholder="12345678"
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

			<!-- Confirm Password -->
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

			<!-- Submit -->
			<Button
				variant="primary"
				type="submit"
				subClass="w-full bg-retroTeal text-retroCream hover:bg-retroCoral transition-all duration-300 font-pixel border-2 border-retroBlack text-lg py-3 mt-2"
			>
				Sign Up
			</Button>
		</form>

		<!-- Link to Sign In -->
		<p class="text-center mt-4 font-pixel text-retroGray text-sm sm:text-base">
			Already have an account?
			<a href="/signin" class="text-retroTeal hover:text-retroCoral transition-colors">Sign In</a>
		</p>
	</section>
</main>
