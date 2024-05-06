import { writable } from 'svelte/store'

export const journeys = writable<journey.IJourney[] | journey.IJourneyPrice[]>(
    []
)
export const selectedJourneyIndex = writable<[number, number]>([0, 0])
