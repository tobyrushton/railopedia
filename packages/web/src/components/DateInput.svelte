<script>
    import { DatePicker } from 'date-picker-svelte'
    import { clickoutside } from '@svelte-put/clickoutside'
    import { Calendar } from 'lucide-svelte'
    import TimeInput from './TimeInput.svelte'
    import dayjs from 'dayjs'

    export let date = new Date()

    let inputSelected = false

    let formatedDate = ""

    $: formatedDate = dayjs(date).format('ddd D MMM')
</script>

<div class="flex w-80 gap-2">
    <div
    class="flex flex-col grow"
    use:clickoutside
    on:clickoutside={() => inputSelected = false}
>
    <label class="flex relative rounded focus:outline-none p-2 gap-2 border border-solid">
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
            <DatePicker bind:value={date} min={new Date()} max={dayjs().add(3, 'month').toDate()} />        
        </div>
    {/if}
    </div>
    <TimeInput bind:time={date} />
    <TimeInput bind:time={date} hours={false} />   
</div>