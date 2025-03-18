<!-- /src/lib/components/Button.svelte -->
<script lang="ts">
	export let variant = 'primary'; // primary, secondary
	export let onClick: (event: MouseEvent) => void = () => {};
	export let subClass = '';
	export let disabled = false;
	export let type: 'button' | 'submit' | 'reset' | null | undefined = 'button';
</script>

<button
	{disabled}
	on:click={onClick}
	class={`
    ${variant === 'primary' ? 'bg-retroCoral text-retroCream' : 'bg-retroBlue text-retroGray'} 
    ${subClass ? subClass : ''}
    px-3 py-1.5 cursor-pointer sm:px-4 sm:py-2 rounded-md font-pixel text-xs sm:text-sm hover:opacity-80 transition-opacity w-full sm:w-auto`}
	{type}
>
	<slot />
</button>
