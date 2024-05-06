<script lang="ts">
    import dayjs from 'dayjs'
    import { onMount } from 'svelte';
    export let time = new Date()
    export let hours = true

    let options: number[] = []

    $: {
        options = []
        const date = dayjs(time)
        let currentDate = date.isSame(dayjs(), 'day')
    
        if(hours) {
            let currentHour = currentDate ? dayjs().hour() : 0
            if(currentDate && currentHour === dayjs().hour() && date.minute() > 45) currentHour++

            for(let i = currentHour; i < 24; i++) {
                options.push(i)
            }
        } else {
            const currentMinute = currentDate ? date.isSame(dayjs(), 'hour') ? date.minute() : 0 : 0
            for(let i = currentMinute; i < 60; i++) {
                if(i % 15 === 0 || i === 0) options.push(i)
            }
        }
    }

    let selected: number

    $: {
        const date = dayjs(time).set('millisecond', 0).set('second', 0)
        if(time && selected){
            if(hours) {
                time = date.hour(selected).toDate()
            } else {
                time = date.minute(selected).toDate()
            }
        }
    }

    onMount(() => {
        if(!selected) {
            selected = options[0]
        }
    })

</script>

<div>
    <select
        class="flex h-full bg-ternary p-2 rounded"
        name={hours ? "hours" : "minutes"}
        bind:value={selected}
    >
        {#each options as option}
            <option value={option}>{option.toString().padStart(2, "0")}</option>
        {/each}
    </select>
</div>