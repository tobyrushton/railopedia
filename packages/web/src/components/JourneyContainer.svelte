<script lang="ts">
    import JourneyItem from './JourneyItem.svelte'
    import { inboundJourneys, outboundJourneys, selectedJourneyIndex } from '../stores/journey'
    import { onDestroy, onMount } from 'svelte'
    import { isIJourney } from '../utils/types'

    export let journeyListProp = ""

    let journeyList: journey.IJourney[] | journey.IJourneyPrice[] = []
    export let returnJourney: boolean = false

    let selectedIndex:number = 0

    let unsubscribeFromJourney = () => {}

    if(!returnJourney) {
        unsubscribeFromJourney = outboundJourneys.subscribe(value => {
            journeyList = value
        })
    } else {
        unsubscribeFromJourney = inboundJourneys.subscribe(value => {
            journeyList = value
        })
    }

    const unsubcribedFromIndex = selectedJourneyIndex.subscribe(value => {
        if(returnJourney) {
            selectedIndex = value[1]
        } else {
            selectedIndex = value[0]
        }
    })

    onMount(() => {
        if(journeyListProp) journeyList = JSON.parse(journeyListProp)
        if(!returnJourney && journeyList.length > 0 && isIJourney(journeyList)) {
            const outboundJourney = journeyList[0]
            inboundJourneys.set(outboundJourney.Prices)
        }
    })

    onDestroy(() => {
        unsubscribeFromJourney()
        unsubcribedFromIndex()
    })

    const getPrice = (index: number): number => {
        if(isIJourney(journeyList)) {
            const journey = journeyList[index]
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

            for(let i = 0; i < journeyList[index].Prices.length; i++) {
                if(journeyList[index].Prices[i].Price < cheapest) {
                    cheapest = journeyList[index].Prices[i].Price
                }
            }

            return cheapest
        }
    }

</script>

{#key journeyList}
    {#each journeyList as journey, index}
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
                        inboundJourneys.set(outboundJourney.Prices)
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
                        selectedJourneyIndex.update(() => [index, cheapeastReturnIdx])
                    }
                }
            }}
        />
    {/each}
{/key}

