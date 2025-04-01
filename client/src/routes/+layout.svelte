<script>
    import '../app.css';
    import Navbar from '$lib/components/Navbar.svelte';
    import ChatBubble from '$lib/components/ChatBubble.svelte';
</script>

<Navbar />
<slot />
<ChatBubble />
