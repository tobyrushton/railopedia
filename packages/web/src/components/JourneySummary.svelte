<script lang="ts">
    import { onDestroy } from 'svelte'
    import { ChevronRight } from 'lucide-svelte'
    import { selectedJourneyIndex, journeys } from '../stores/journey'
    import { isIJourney } from '../utils/types'

    let selectedIndex:[number, number] = [0, 0]
    let journeyList: journey.IJourneyPrice[] | journey.IJourney[] = []
    let selectedJourneyList: journey.IJourneyPrice[] = []

    const unsubcribeFromIndex = selectedJourneyIndex.subscribe(value => {
        selectedIndex = value
    })

    $: {
        if (isIJourney(journeyList)) {
            selectedJourneyList = journeyList[selectedIndex[1]].Prices
        } else {
            selectedJourneyList = journeyList
        }
    }

    const unsubscribeFromJourneys = journeys.subscribe(value => {
        journeyList = value
    })

    onDestroy(() => {
        unsubcribeFromIndex()
        unsubscribeFromJourneys()
    })

    let cheapestJourney: journey.IPrice
    
    $: {
        if (journeyList.length !== 0){
            const journey = selectedJourneyList[selectedIndex[1]]

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
        {#if journeyList.length !== 0}
            {#each selectedJourneyList[selectedIndex[1]].Prices as price}
                <li class="flex justify-between">
                    <span>{price.Provider}</span>
                    <span>£{price.Price.toFixed(2)}</span>
                </li>
            {/each}
        {/if}
    </ul>
    <a 
        class="bg-primary rounded text-white p-2 justify-center font-semibold text-xl mt-2 cursor-pointer flex"
        href={cheapestJourney?.Link}
        target="_blank"
    >
        CONTINUE <ChevronRight class="size-7"/>
    </a>
    <p class="text-sm">
        *Prices may be slightly inaccurate
    </p>
</div>