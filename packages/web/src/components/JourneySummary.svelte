<script lang="ts">
    import { onDestroy } from 'svelte'
    import { ChevronRight } from 'lucide-svelte'
    import { inboundJourneys, selectedJourneyIndex } from '../stores/journey'

    let selectedIndex:[number, number] = [0, 0]
    let inboundJourneysList: journey.IJourneyPrice[] = []

    const unsubcribedFromIndex = selectedJourneyIndex.subscribe(value => {
        selectedIndex = value
    })
    const unsubscribeFromInboundJourneys = inboundJourneys.subscribe(value => {
        inboundJourneysList = value
    })

    onDestroy(() => {
        unsubcribedFromIndex()
        unsubscribeFromInboundJourneys()
    })

    let cheapestJourney: journey.IPrice
    
    $: {
        if(inboundJourneysList.length !== 0){
            const journey = inboundJourneysList[selectedIndex[1]]
            console.log(journey)

            let cheapest: journey.IPrice = journey.Prices[0]
            for (let i = 1; i < journey.Prices.length; i++) {
                const price = journey.Prices[i]
                if (price.Price < cheapest.Price) {
                    cheapest.Price = price.Price
                }
            }

            cheapestJourney = cheapest
        }
    }
</script>

<div class="flex flex-col rounded shadow h-fit p-3">
    <h3 class="text-xl font-semibold flex justify-between">
        TOTAL <span>£{cheapestJourney?.Price.toFixed(2)}</span>
    </h3>
    <ul>
        {#if inboundJourneysList.length !== 0}
            {#each inboundJourneysList[selectedIndex[1]].Prices as price}
                <li class="flex justify-between">
                    <span>{price.Provider}</span>
                    <span>£{price.Price.toFixed(2)}</span>
                </li>
            {/each}
        {/if}
    </ul>
    <a 
        class="bg-primary rounded text-white p-2 text-center font-semibold text-xl mt-2 cursor-pointer flex"
        href={cheapestJourney?.Link}
        target="_blank"
    >
        CONTINUE <ChevronRight class="size-7"/>
    </a>
    <p class="text-sm">
        *Prices may be slightly inaccurate
    </p>
</div>