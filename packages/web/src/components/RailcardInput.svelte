<script lang="ts">
    import { railcards, type Railcard } from '../utils/railcard'
    export let value: Railcard

</script>

<select 
    class="w-80 rounded border border-solid p-2"
    bind:value={value}
>
    {#each Object.entries(railcards) as [code, name]}
        <option value={code} selected={code === "N"}>
            {name}
        </option>
    {/each}
</select>