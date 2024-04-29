<script lang="ts">
    import { onDestroy } from 'svelte'
    import { ChevronRight } from 'lucide-svelte'
    import { inboundJourneys, outboundJourneys, selectedJourneyIndex } from '../stores/journey'

    let selectedIndex:[number, number] = [0, 0]
    let inboundJourneysList: journey.IJourneyPrice[] = []
    // let outboundJourneysList: journey.IJourneyPrice[] = []

    const unsubcribedFromIndex = selectedJourneyIndex.subscribe(value => {
        selectedIndex = value
    })
    // const unsubscribeFromOutboundJourneys = outboundJourneys.subscribe(value => {
    //     outboundJourneysList = value
    // })
    const unsubscribeFromInboundJourneys = inboundJourneys.subscribe(value => {
        inboundJourneysList = value
    })

    onDestroy(() => {
        unsubcribedFromIndex()
        // unsubscribeFromOutboundJourneys()
        unsubscribeFromInboundJourneys()
    })

    let price = 0
    
    $: {
        if(inboundJourneysList.length !== 0){
            const journey = inboundJourneysList[selectedIndex[1]]

            let cheapest = Infinity
            for (let i = 0; i < journey.Prices.length; i++) {
                const price = journey.Prices[i]
                if (price.Price < cheapest) {
                    cheapest = price.Price
                }
            }

            price = cheapest
        }
    }
</script>

<div class="flex flex-col rounded shadow h-fit p-3">
    <h3 class="text-xl font-semibold flex justify-between">
        TOTAL <span>£{price.toFixed(2)}</span>
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
    <!-- TODO: INSERT REAL LINK HERE -->
    <a 
        class="bg-primary rounded text-white p-2 text-center font-semibold text-xl mt-2 cursor-pointer flex"
        href="/"
        target="_blank"
    >
        CONTINUE <ChevronRight class="size-7"/>
    </a>
    <p class="text-sm">
        *Prices may be slightly inaccurate
    </p>
</div>