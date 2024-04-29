<script lang="ts">
    import dayjs from 'dayjs'
    import { cn } from '../utils/tailwind'

    export let arrival: string
    export let departure: string
    export let price: number
    export let selected: boolean = false
    export let outbound: boolean = false
    export let onClick = () => {}

    const arrivalDate = dayjs(arrival)
    const departureDate = dayjs(departure)
    const duration = arrivalDate.diff(departureDate)
    const totalSeconds = Math.floor(duration / 1000)
    const totalMinutes = Math.floor(totalSeconds / 60)
    const totalHours = Math.floor(totalMinutes / 60)
</script>

<button class={cn("flex gap-2 p-2 px-4 cursor-pointer w-full", {
    "border-solid border-y-2 border-primary bg-primary-tint": selected,
})} on:click={onClick}>
    <span class="flex flex-col gap-2 grow items-start">
        <h3 class="font-semibold text-xl">
            {departureDate.format('HH:mm')} -> {arrivalDate.format('HH:mm')}
        </h3>
        <p>
            {totalHours}h {totalMinutes % 60}m
        </p>
    </span>
    <span>
        <p>{outbound ? 'From' : ''} £{price}</p>
    </span>
</button>