<script>
    import { DatePicker } from 'date-picker-svelte'
    import { clickoutside } from '@svelte-put/clickoutside'
    import { Calendar } from 'lucide-svelte'
    import dayjs from 'dayjs'

    export let date = new Date()

    let inputSelected = false

    let formatedDate = ""

    $: formatedDate = dayjs(date).format('HH:mm ddd D MMM')
</script>

<div
    class="flex flex-col w-80"
    use:clickoutside
    on:clickoutside={() => inputSelected = false}
>
    <label class="flex relative bg-ternary rounded focus:outline-none p-2 gap-2">
        <Calendar class="pointer-events-none size-6" />
        <input 
            class="focus:outline-none w-full bg-inherit"
            type="text"
            on:click={() => inputSelected = !inputSelected}
            value={formatedDate}
            readonly 
        />
    </label>
    {#if inputSelected}
        <div class="absolute z-10">
            <DatePicker bind:value={date} min={new Date()} timePrecision="minute"/>        
        </div>
    {/if}
</div>