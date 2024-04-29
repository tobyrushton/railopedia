import { writable } from 'svelte/store'

export const outboundJourneys = writable<journey.IJourney[]>([])
export const inboundJourneys = writable<journey.IJourneyPrice[]>([])
export const selectedJourneyIndex = writable<[number, number]>([0, 0])
