<!-- /src/lib/components/Button.svelte -->
<script lang="ts">
	export let variant = 'primary'; // primary, secondary
	export let onClick = () => {};
	export let subClass = '';
</script>

<button
	on:click={onClick}
	class={`
    ${variant === 'primary' ? 'bg-retroCoral text-retroCream' : 'bg-retroBlue text-retroGray' } 
    ${subClass ? subClass : ''}
    px-3 py-1.5 sm:px-4 sm:py-2 rounded-md font-pixel text-xs sm:text-sm hover:opacity-80 transition-opacity w-full sm:w-auto`}
>
	<slot />
</button>
