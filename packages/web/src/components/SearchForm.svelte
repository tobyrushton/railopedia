<script lang="ts">
    import StationInput from './StationInput.svelte'
    import DateInput from './DateInput.svelte'
    import JourneyRadio from './JourneyRadio.svelte'
    import RailcardInput from './RailcardInput.svelte'
    import AddReturn from './AddReturn.svelte'
    import ErrorMessage from './ErrorMessage.svelte'
    import dayjs from 'dayjs'
    import z from 'zod'
    import { stationIsValid } from '../utils/station'
    import { type Railcard } from '../utils/railcard'

    interface IForm {
        journey: 'single' | 'return'
        outboundDate: Date
        returnDate: Date
        destination: journey.IStation
        origin: journey.IStation
        railcard: Railcard
    }

    const form: IForm = {
        journey: 'single',
        outboundDate: dayjs().set('second', 0).set('millisecond', 0).toDate(),
        returnDate: dayjs(new Date()).add(1, 'hour').set('second', 0).set('millisecond', 0).toDate(),
        destination: { name: '', code: '' },
        origin: { name: '', code: '' },
        railcard: 'N'
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

    const formatDate = (date: Date): string => dayjs(date).format('YYYY-MM-DDTHH:mm:ss')

    const handleSubmit = (): void => {
        const res = formSchema.safeParse(form)
        errorMessages = {}
        if(res.success) {
            window.location.href = 
                `/search?origin=${form.origin.code}&destination=${form.destination.code}&departure=${formatDate(form.outboundDate)}&return=${form.journey === 'return' ? formatDate(form.returnDate) : ''}&railcard=${form.railcard}`
        } else {
            res.error.errors.map(({ path, message }) => {
                errorMessages[path[0]] = message
            })
        }
    }
    
</script>

<form class="bg-white shadow-md rounded-lg p-6 md:p-8 lg:p-10 space-y-4">
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
    <span class="flex flex-col gap-3 w-full sm:justify-between sm:flex-row">
        <RailcardInput bind:value={form.railcard} />
        <button 
            class="w-80 bg-primary text-white font-bold p-2 rounded"
            on:click|preventDefault={handleSubmit}
        >
            Find tickets
        </button>
    </span>
</form>