<script lang="ts">
    import JourneyItem from './JourneyItem.svelte'
    import { journeys, selectedJourneyIndex } from '../stores/journey'
    import { onDestroy, onMount } from 'svelte'
    import { isIJourney } from '../utils/types'

    export let journeyListProp: string = ""
    export let returnJourney: boolean = false

    let journeyList: journey.IJourney[] | journey.IJourneyPrice[] = []
    let renderedJourneyList: journey.IJourney[] | journey.IJourneyPrice[] = []

    let selectedIndexes: [number, number] = [0, 0]

    let selectedIndex: number = 0

    let unsubscribeFromJourney = journeys.subscribe(value => {
        journeyList = value
    })

    const unsubcribedFromIndex = selectedJourneyIndex.subscribe(value => {
        if(returnJourney) {
            selectedIndex = value[1]
        } else {
            selectedIndex = value[0]
        }
        selectedIndexes = value
    })

    onMount(() => {
        if(journeyListProp) { 
            journeyList = JSON.parse(journeyListProp)
            journeys.set(journeyList)
        }
    })

    onDestroy(() => {
        unsubscribeFromJourney()
        unsubcribedFromIndex()
    })

    $: {
        if(isIJourney(journeyList) && returnJourney && journeyList.length ) {
            renderedJourneyList = journeyList[selectedIndexes[0]].Prices
        } else {
            renderedJourneyList = journeyList
        }
    }

    const getPrice = (index: number): number => {
        if(isIJourney(renderedJourneyList)) {
            const journey = renderedJourneyList[index]
            let cheapest: number = Infinity

            for(let i = 0; i < journey.Prices.length; i++) {
                for(let j = 0; j < journey.Prices[i].Prices.length; j++) {
                    if(journey.Prices[i].Prices[j].Price < cheapest) {
                        cheapest = journey.Prices[i].Prices[j].Price
                    }
                }
            }

            return cheapest
        } else {
            let cheapest: number = Infinity

            for(let i = 0; i < renderedJourneyList[index].Prices.length; i++) {
                if(renderedJourneyList[index].Prices[i].Price < cheapest) {
                    cheapest = renderedJourneyList[index].Prices[i].Price
                }
            }

            return cheapest
        }
    }

</script>

{#key renderedJourneyList}
    {#each renderedJourneyList as journey, index}
        <JourneyItem
            departure={journey.DepartureTime}
            arrival={journey.ArrivalTime}
            price={getPrice(index)}
            selected={index === selectedIndex}
            outbound={!returnJourney}
            onClick={() => {
                if(returnJourney) {
                    selectedJourneyIndex.update(value => [value[0], index])
                } else {
                    if(isIJourney(journeyList)) {
                        const outboundJourney = journeyList[index]
                        let cheapeastReturnIdx = 0
                        let cheapestReturn = Infinity
                        for(let i = 0; i < outboundJourney.Prices.length; i++) {
                            for(let j = 0; j < outboundJourney.Prices[i].Prices.length; j++) {
                                if(outboundJourney.Prices[i].Prices[j].Price < cheapestReturn) {
                                    cheapestReturn = outboundJourney.Prices[i].Prices[j].Price
                                    cheapeastReturnIdx = i
                                }
                            }
                        }
                        selectedJourneyIndex.set([index, cheapeastReturnIdx])
                    } else {
                        selectedJourneyIndex.set([index, index])
                    }
                }
            }}
        />
    {/each}
{/key}
