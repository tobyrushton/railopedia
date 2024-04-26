<script>
    import { createEventDispatcher } from 'svelte'

    export let journey = 'single'
</script>

<span class="flex gap-2" role="radiogroup">
    <input 
        type="radio"
        class="accent-primary" 
        name="journey" 
        id="single" 
        value="single" 
        bind:group={journey} 
    />
    <label for="single">Single</label>
    <input 
        type="radio"
        class="accent-primary"
        name="journey" 
        id="return" 
        value="return" 
        bind:group={journey} 
    />
    <label for="return">Return</label>
</span>