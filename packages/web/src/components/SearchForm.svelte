<script lang="ts">
    import StationInput from './StationInput.svelte'
    import DateInput from './DateInput.svelte'
    import JourneyRadio from './JourneyRadio.svelte'
    import AddReturn from './AddReturn.svelte'
    import ErrorMessage from './ErrorMessage.svelte'
    import dayjs from 'dayjs'
    import z from 'zod'
    import { stationIsValid } from '../utils/station'

    interface IForm {
        journey: 'single' | 'return'
        outboundDate: Date
        returnDate: Date
        destination: journey.IStation
        origin: journey.IStation
    }

    const form: IForm = {
        journey: 'single',
        outboundDate: new Date(),
        returnDate: dayjs(new Date()).add(1, 'hour').toDate(),
        destination: { name: '', code: '' },
        origin: { name: '', code: '' }
    }

    const formSchema = z.object({
        journey: z.enum(['single', 'return']),
        outboundDate: z.date().refine(value => dayjs(value).isAfter(dayjs()), {
            message: 'Date must be in the future',
        }),
        returnDate: z.date().refine(value => dayjs(value).isAfter(dayjs()) || form.journey === 'single', {
            message: 'Date must be in the future',
        }),
        destination: z.object({
            name: z.string(),
            code: z.string()
        }).refine((destination: journey.IStation) => stationIsValid(destination), { message: "Please enter a valid station"}),
        origin: z.object({
            name: z.string(),
            code: z.string()
        }).refine((origin: journey.IStation) => stationIsValid(origin), { message: "Please enter a valid station"})
    }).superRefine((data, context) => {
        if (data.outboundDate >= data.returnDate && data.journey === 'return') {
            context.addIssue({
                code: z.ZodIssueCode.custom,
                path: ['returnDate'],
                message: 'Return date must be after outbound date'
            })
        }
    })

    let errorMessages = {}

    const handleSubmit = (): void => {
        const res = formSchema.safeParse(form)
        errorMessages = {}
        if(res.success) {
            window.location.href = 
                `/search?origin=${form.origin.code}&destination=${form.destination.code}&outboundDate=${form.outboundDate.toISOString()}&returnDate=${form.journey === 'return' ? form.returnDate.toISOString() : ''}`
        } else {
            res.error.errors.map(({ path, message }) => {
                errorMessages[path[0]] = message
            })
        }
    }
    
</script>

<form class="flex flex-col gap-4 w-fit">
    <JourneyRadio bind:journey={form.journey} />
    <span class="flex flex-col gap-4 sm:flex-row">
        <span class="flex flex-col gap-2">
            <StationInput placeholder="Departing from" bind:value={form.origin}/>
            {#if errorMessages['origin']}
                <ErrorMessage>{errorMessages['origin']}</ErrorMessage>
            {/if}
            <DateInput bind:date={form.outboundDate}/>  
            {#if errorMessages['outboundDate']}
                <ErrorMessage>{errorMessages['outboundDate']}</ErrorMessage>
            {/if}  
        </span>
        <span class="flex flex-col gap-2">
            <StationInput placeholder="Arriving at" bind:value={form.destination}/>
            {#if errorMessages['destination']}
                <ErrorMessage>{errorMessages['destination']}</ErrorMessage>
            {/if}
            {#if form.journey === 'return'}
                <DateInput bind:date={form.returnDate}/>
                {#if errorMessages['returnDate']}
                    <ErrorMessage>{errorMessages['returnDate']}</ErrorMessage>
                {/if}
            {:else}
                <AddReturn bind:journey={form.journey} />
            {/if}
        </span>
    </span>
    <button 
        class="w-full self-end bg-primary text-white font-bold p-2 rounded sm:w-fit"
        on:click|preventDefault={handleSubmit}
    >
        Find tickets
    </button>
</form>