<script>
    import { Plus } from 'lucide-svelte'
    export let journey = 'single'
</script> 

<button 
    class="flex gap-2 p-2 border border-solid rounded hover:cursor-pointer" 
    tabindex="0"
    on:click|preventDefault={() =>  journey = 'return'}
>
    <Plus class="size-6" />
    <p>Add return</p>
</button>